package logic_rpc

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/global"
	"XtTalkServer/pb"
	"XtTalkServer/services/logic/logic_model"
	"XtTalkServer/utils"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-zookeeper/serverplugin"
	"github.com/smallnest/rpcx/server"
	"google.golang.org/protobuf/proto"
	"time"
)

var LogicRpcServer *server.Server = nil

func InitRpcServer(ctx context.Context) {
	conf := global.Config.Services.Logic
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
		serverplugin.WithZKServiceAddress(fmt.Sprintf("tcp@%s:%d", localIp, conf.RpcPort)),
		serverplugin.WithZKServersAddress(zkConf.Servers),
		serverplugin.WithZKBasePath("rpc_logic"),
		serverplugin.WithZKMetrics(metrics.NewRegistry()),
		serverplugin.WithZkUpdateInterval(time.Minute),
	)
	if err := zkPlugin.Start(); err != nil {
		glog.Fatalf(ctx, "注册RPC Zookeeper服务失败: %s", err.Error())
	}
	LogicRpcServer.Plugins.Add(zkPlugin)
	defer zkPlugin.Stop()
	if err := LogicRpcServer.Register(new(LogicRpcService), ""); err != nil {
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

type LogicRpcService struct {
}

// LogicData 逻辑层RPC调用，返回error则客户端直接断开连接
func (c *LogicRpcService) LogicData(ctx context.Context, req *pb.LogicDataReq, res *pb.LogicDataRes) error {
	glog.Debugf(ctx, "收到来自接入层数据的消息")
	device := logic_model.ConnDevice{
		Context:   ctx,
		SessionId: req.SessionId,
		ServerId:  req.ServerId,
	}
	//读取客户端Redis信息
	if userValue, err := global.Redis.Get(ctx, fmt.Sprintf(conts.RK_ClientAuth, req.SessionId)).Bytes(); err == nil {
		if err := gjson.Unmarshal(userValue, &device.UserClient); err != nil {
			glog.Warningf(ctx, "解析用户Redis数据异常: %s", err.Error())
			return err //解析用户数据异常
		}
	}

	data, err := func() ([]byte, error) {
		//收到了来自通讯层发来的数据
		reply, err := routeService(device, req.GetCommand(), req.Data)
		if err != nil {
			return nil, err
		}
		if reply == nil {
			return nil, nil
		}
		bytes, err := proto.Marshal(reply)
		return bytes, err
	}()
	if err != nil {
		return err
	}
	if data != nil {
		res.IsSend = true
	}
	res.DataFormat = req.GetDataFormat()
	res.Data = data
	res.Command = req.Command
	return nil
}

func routeService(device logic_model.ConnDevice, commandId pb.Packet, data []byte) (replyMsg proto.Message, replyErr error) {
	//需要登录的数据包
	needLoginPacket := []pb.Packet{
		pb.Packet_GetProfile,
		pb.Packet_GetUser,
		pb.Packet_GetFriendList,
		pb.Packet_GetFriend,
		pb.Packet_RemoveFriend,
	}
	if utils.InArray(commandId, needLoginPacket) {
		if device.UserClient == nil {
			return nil, gerror.Newf("客户端未登录")
		}
	}

	switch commandId {

	case pb.Packet_Login: //登录
		{
			var params pb.PacketLoginReq
			if replyErr = proto.Unmarshal(data, &params); replyErr != nil {
				return
			}
			replyMsg, replyErr = Login.Login(device, &params)
			return
		}
	case pb.Packet_GetProfile: //获取当前账号信息
		{
			var params pb.PacketGetProfileReq
			if replyErr = proto.Unmarshal(data, &params); replyErr != nil {
				return
			}
			replyMsg, replyErr = Self.GetProfile(device, &params)
			return
		}
	case pb.Packet_ModifyProfile: //修改当前账号信息
		{
			var params pb.PacketModfiyProfileReq
			if replyErr = proto.Unmarshal(data, &params); replyErr != nil {
				return
			}
			replyMsg, replyErr = Self.ModifyProfile(device, &params)
			return
		}
	case pb.Packet_GetFriendList: //获取好友列表
		{
			var params pb.PacketGetFriendListReq
			if replyErr = proto.Unmarshal(data, &params); replyErr != nil {
				return
			}
			replyMsg, replyErr = Friend.GetFriendList(device, &params)
			return
		}
	case pb.Packet_GetFriend: //获取好友信息
		{
			var params pb.PacketGetFriendReq
			if replyErr = proto.Unmarshal(data, &params); replyErr != nil {
				return
			}
			replyMsg, replyErr = Friend.GetFriend(device, &params)
			return
		}
	case pb.Packet_PrivateMsg:
		{
			var params pb.PacketPrivateMsg
			if replyErr = proto.Unmarshal(data, &params); replyErr != nil {
				return
			}
			replyErr = Friend.SendMsg(device, &params)
			return
		}
	}
	return nil, gerror.Newf("未知操作命令")
}
