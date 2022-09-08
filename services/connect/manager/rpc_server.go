package manager

import (
	"XtTalkServer/global"
	"XtTalkServer/pb"
	"XtTalkServer/utils"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-zookeeper/serverplugin"
	"github.com/smallnest/rpcx/server"
	"time"
)

var LogicRpcServer *server.Server = nil

func InitRpcServer(ctx context.Context) {
	conf := global.Config.Services.Connect
	LogicRpcServer = server.NewServer()
	localIp := utils.GetLocalIp()
	if localIp == nil {
		glog.Fatalf(ctx, "无法获取到当前系统IP地址")
	}
	zkConf := global.Config.Zookeeper
	if 0 >= len(zkConf.Servers) {
		glog.Fatalf(ctx, "未配置Zookeeper服务信息")
	}
	//注册zk插件
	zkPlugin := serverplugin.NewZooKeeperRegisterPlugin(
		serverplugin.WithZKServiceAddress(fmt.Sprintf("%d#tcp@%s:%d", conf.Id, localIp, conf.RpcPort)),
		serverplugin.WithZKServersAddress(zkConf.Servers),
		serverplugin.WithZKBasePath("rpc_connect"),
		serverplugin.WithZKMetrics(metrics.NewRegistry()),
		serverplugin.WithZkUpdateInterval(time.Minute),
	)
	if err := zkPlugin.Start(); err != nil {
		glog.Fatalf(ctx, "注册RPC Zookeeper服务失败: %s", err.Error())
	}
	LogicRpcServer.Plugins.Add(zkPlugin)
	defer zkPlugin.Stop()
	if err := LogicRpcServer.Register(new(ConnectRpcService), ""); err != nil {
		glog.Fatalf(ctx, "注册RPC服务失败: %s", err.Error())
	}
	if err := LogicRpcServer.Serve("tcp", fmt.Sprintf(":%d", conf.RpcPort)); err != nil {
		glog.Fatalf(ctx, "启动RPC服务端失败: %s", err.Error())
	}
}

func StopRpcServer(ctx context.Context) {
	if LogicRpcServer != nil {
		for _, plugin := range LogicRpcServer.Plugins.All() {
			if p, has := plugin.(*serverplugin.ZooKeeperRegisterPlugin); has {
				if err := p.Stop(); err != nil {
					glog.Warningf(ctx, "停止Rpcx -> Zookeeper插件失败: %s", err.Error())
				}
			}
		}
	}
}

type ConnectRpcService struct {
}

// SendClientPacket 向指定Session客户端发送消息
func (c *ConnectRpcService) SendClientPacket(ctx context.Context, req *pb.ConnectSendClientReq, res *pb.ConnectSendClientRes) error {
	glog.Infof(ctx, "收到来自: %s 发来的消息: %s", req.GetSessionId(), req.GetPayload())
	client := ClientManager.GetClientBySession(req.GetSessionId())
	if client == nil {
		//	//设备不在线
		res = &pb.ConnectSendClientRes{
			RetCode: pb.ConnectRetCode_CR_Offline,
		}
		return nil
	}
	if err := client.SendServerBytes(req.GetCommand(), req.GetPayload()); err != nil {
		glog.Warningf(ctx, "接入层发送客户端数据失败: %s", err.Error())
		res = &pb.ConnectSendClientRes{
			RetCode: pb.ConnectRetCode_CR_Error,
		}
		return nil
	}
	res = &pb.ConnectSendClientRes{
		RetCode: pb.ConnectRetCode_CR_Success,
	}
	return nil
}

// KickClientPacket 踢出客户端
func (c *ConnectRpcService) KickClientPacket(ctx context.Context, req *pb.ConnectKickClientReq, res *pb.ConnectKickClientRes) error {
	client := ClientManager.GetClientBySession(req.GetSessionId())
	if client != nil {
		return client.Close()
	}
	return nil
}
