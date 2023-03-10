package connect

import (
	"XtTalkServer/core"
	"XtTalkServer/global"
	"XtTalkServer/internal"
	"XtTalkServer/internal/connect/manager"
	"XtTalkServer/internal/connect/server_protocol"
	"XtTalkServer/internal/connect/service"
	"XtTalkServer/internal/connect/types"
	"XtTalkServer/utils"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"golang.org/x/net/context"
)

func RunApplication() {
	ctx := gctx.New()

	core.Entry.ExecStart(ctx, []internal.InitCtx{
		core.Log,
		core.Viper,
		core.Mysql,
		core.Redis,
	})
	conf := global.Config.Services.Connect
	glog.Debugf(ctx, "当前接入服务器ID: %d", conf.Id)
	core.Snowflake.Initialize(conf.Id, ctx)
	//启动RPC服务端
	go manager.InitRpcServer(ctx)
	//初始化socket客户端管理器
	manager.InitClientManager()
	//初始化socket服务发现
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
	go manager.CoreServer.OnStart(ctx)
	utils.WaitExit()
	StopApplication(ctx)
}

func StopApplication(ctx context.Context) {
	manager.CoreServer.OnShutdown(ctx)
	manager.StopRpcServer(ctx)
	service.ServiceDiscovery.UnRegister()
	glog.Infof(ctx, "连接层已停止运行")
}
