package logic

import (
	"XtTalkServer/core"
	"XtTalkServer/services/logic/logic_rpc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"os"
	"os/signal"
	"syscall"
)

func RunApplication() {
	ctx := gctx.New()
	core.Log.Initialize(ctx)
	core.Viper.Initialize(ctx)
	core.Mysql.Initialize(ctx)
	core.Redis.Initialize(ctx)

	go logic_rpc.InitRpcServer(ctx)
	//初始化RPC调用连接层客户端
	logic_rpc.InitConnectRpcClient(ctx)
	glog.Infof(ctx, "逻辑层已启动完成")
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigc

	logic_rpc.StopRpcServer(ctx)
	glog.Infof(ctx, "逻辑层已停止运行")
}
