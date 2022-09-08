package logic

import (
	"XtTalkServer/core"
	"XtTalkServer/global"
	"XtTalkServer/services/logic/logic_rpc"
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
	conf := global.Config.Services.Logic
	glog.Debugf(ctx, "当前业务服务器ID: %d", conf.Id)
	go logic_rpc.InitRpcServer(ctx)
	glog.Infof(ctx, "逻辑层已启动完成")
	utils.WaitExit()
	StopApplication(ctx)
}
func StopApplication(ctx context.Context) {
	logic_rpc.StopRpcServer(ctx)
	glog.Infof(ctx, "逻辑层已停止运行")
}
