package message

import (
	"XtTalkServer/core"
	"XtTalkServer/global"
	"XtTalkServer/services/message/service"
	"XtTalkServer/services/message/types"
	"XtTalkServer/utils"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"golang.org/x/net/context"
	"reflect"
	"time"
)

func RunApplication() {
	ctx := gctx.New()
	core.Log.Initialize(ctx)
	core.Viper.Initialize(ctx)
	core.Mysql.Initialize(ctx)
	core.Redis.Initialize(ctx)
	core.Mongo.Initialize(ctx)
	core.Rabbitmq.Initialize(ctx)
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
		&service.PrivateMsgConsumer{},
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
					if err := item.Start(ctx, ch); err != nil {
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
