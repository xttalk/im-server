package connect

import (
	"XtTalkServer/core"
	"XtTalkServer/global"
	"XtTalkServer/services/connect/manager"
	"XtTalkServer/services/connect/server_protocol"
	"XtTalkServer/services/connect/service"
	"XtTalkServer/services/connect/types"
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
	conf := global.Config.Services.Connect
	glog.Debugf(ctx, "当前接入服务器ID: %d", conf.Id)
	core.Snowflake.Initialize(conf.Id, ctx)
	//初始化对接连接层rpc
	manager.InitLogicRpcClient(ctx)
	//初始化连接层rpc服务器
	go manager.InitRpcServer(ctx)
	//初始化客户端管理器
	manager.InitClientManager()

	//服务发现初始化
	if err := service.InitDiscovery(); err != nil {
		glog.Fatalf(ctx, "初始化服务发现失败: %s", err.Error())
	}
	//启动socket服务器
	servers := []types.IServer{
		server_protocol.CreateWsServer(conf.WsPort),   //创建ws服务端
		server_protocol.CreateTcpServer(conf.TcpPort), //创建tcp服务端
	}

	for index, s := range servers {
		go func(index int, s types.IServer) {
			if err := s.Start(); err != nil {
				glog.Fatalf(ctx, "启动服务端失败: %v-> %s", s.GetProtocol(), err.Error())
			}
		}(index, s)
	}
	go manager.CoreServer.OnStart()
	glog.Infof(ctx, "连接层已启动完成")
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigc
	manager.StopRpcServer(ctx)
	//解除注册zk路径
	service.ServiceDiscovery.UnRegister()
	glog.Infof(ctx, "连接层已停止运行")
}
