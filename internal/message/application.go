package message

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/core"
	"XtTalkServer/global"
	"XtTalkServer/internal"
	"XtTalkServer/internal/message/service"
	"XtTalkServer/internal/message/types"
	"XtTalkServer/utils"
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	"reflect"
	"time"
)

func RunApplication() {
	ctx := gctx.New()
	core.Entry.ExecStart(ctx, []internal.InitCtx{
		core.Log,
		core.Viper,
		core.Mysql,
		core.Redis,
		core.Mongo,
		core.Rabbitmq,
	})
	conf := global.Config.Services.Message
	_ = conf
	StartConsumer()
	glog.Infof(ctx, "消息中心已启动完成")
	utils.WaitExit()
	StopApplication(ctx)
}
func StopApplication(ctx context.Context) {
	glog.Infof(ctx, "消息中心已停止运行")
}
func StartConsumer() {
	consumerList := []types.IConsumer{
		&service.FriendChangeConsumer{},
		&service.PrivateMsgSaveConsumer{},
		&service.PrivateMsgConsumer{},
		&service.FriendRequestConsumer{},
	}
	for _, item := range consumerList {
		name := reflect.TypeOf(item).Name()

		ctx := gctx.New()
		channel, err := global.RabbitMQ.GetChannel()
		if err != nil {
			glog.Fatalf(ctx, "[%s] 获取RabbitMQ Channel失败: %s", name, err.Error())
		}
		defer channel.Release()
		if err := item.Init(ctx, channel); err != nil {
			glog.Fatalf(ctx, "[%s] 初始化消费者失败: %s", name, err.Error())
		}
		go func(name string, item types.IConsumer) {
			for {
				func() {
					ch, err := global.RabbitMQ.GetChannel()
					if err != nil {
						glog.Errorf(ctx, "[%s] 获取RabbitMQ Channel失败: %s ,等待重连", name, err.Error())
						return
					}
					defer ch.Release()
					options := item.ConsumerOptions(ctx)

					if err := channel.CreateConsumer(ctx, options, func(ctx context.Context, d amqp.Delivery) {
						//消息重发机制
						rkKey := fmt.Sprintf(conts.RkMqRetry, options.QueueName)
						rkField := d.MessageId
						if err := func(d amqp.Delivery) error {
							if execErr := item.Callback(ctx, channel, d); err != nil {
								//消息重发
								retrySize, _ := global.Redis.HGet(ctx, rkKey, rkField).Int()
								//判断是否还需要继续重发
								currRetrySize, _ := global.Redis.HIncrBy(ctx, rkKey, rkField, 1).Uint64()
								if options.RetrySize > int(currRetrySize) {
									glog.Warningf(ctx, "[%v-%v]消息发送失败,等待下一次投递[%d/%d]: %s", options.QueueName, rkField, currRetrySize, options.RetrySize, execErr.Error())
									//消息重发
									time.Sleep(time.Second) //1秒后延迟出去
									return d.Reject(true)   //消息重新排队
								}
								glog.Warningf(ctx, "[%v-%v]消息发送失败,终止投递[%d/%d]: %s", options.QueueName, rkField, retrySize+1, options.RetrySize, execErr.Error())
								//多次失败
								defer global.Redis.HDel(ctx, rkKey, rkField) //执行成功则删除记录
								return d.Reject(false)                       //删除消息
							} else {
								defer global.Redis.HDel(ctx, rkKey, rkField) //执行成功则删除记录
								return d.Ack(false)                          //确认消息消费
							}
						}(d); err != nil {
							glog.Warningf(ctx, "[%v-%v]消息处理失败: %s", options.QueueName, rkField, err.Error())
						}
					}); err != nil {
						glog.Errorf(ctx, "[%s] 启动消费者失败: %s,等待重试", name, err.Error())
					} else {
						glog.Warningf(ctx, "[%s] 消费者断开,等待重连", name)
					}
				}()
				time.Sleep(time.Second * 3)
			}
		}(name, item)
	}
}
