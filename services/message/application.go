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
)

func RunApplication() {
	ctx := gctx.New()
	core.Log.Initialize(ctx)
	core.Viper.Initialize(ctx)
	core.Mysql.Initialize(ctx)
	core.Redis.Initialize(ctx)
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
		ctx := gctx.New()
		channel, err := global.RabbitMQ.GetChannel()
		if err != nil {
			glog.Fatalf(ctx, "获取RabbitMQ Channel失败: %s", err.Error())
		}
		defer channel.Release()
		if err := item.Init(ctx, channel); err != nil {
			glog.Fatalf(ctx, "初始化消费者失败: %s", err.Error())
		}
		go func(item types.IConsumer) {
			if err := item.Start(ctx, channel); err != nil {
				glog.Fatalf(ctx, "启动消费者失败: %s", err.Error())
			}
		}(item)
	}
}
