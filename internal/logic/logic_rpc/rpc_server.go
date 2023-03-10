package logic_rpc

import (
	"XtTalkServer/global"
	"XtTalkServer/internal/logic/logic_model"
	"XtTalkServer/pb"
	"XtTalkServer/utils"
	"context"
	"fmt"
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
		SessionId: req.GetSessionId(), //来自sessionID
		ServerId:  req.GetServerId(),  //来自服务器
		UserId:    req.GetUserId(),
	}
	//通过sessionID 读取客户端Redis信息
	//if userValue, err := global.Redis.Get(ctx, fmt.Sprintf(conts.RK_ClientAuth, req.SessionId)).Bytes(); err == nil {
	//	if err := gjson.Unmarshal(userValue, &device.UserClient); err != nil {
	//		glog.Warningf(ctx, "解析用户Redis数据异常: %s", err.Error())
	//		return err //解析用户数据异常
	//	}
	//}

	data, err := func() ([]byte, error) {
		//收到了来自通讯层发来的数据
		var reply proto.Message
		var err error
		if req.GetCommand() == pb.Packet_Login {
			loginReply, loginErr := loginService(device, req.GetData())
			if loginErr != nil {
				return nil, loginErr
			}
			device.UserId = loginReply.GetUid() //返回给连接层用户ID
			reply = loginReply
		} else {
			if device.UserId == 0 {
				//没有登录
				return nil, gerror.Newf("用户未登录")
			}
			reply, err = routeService(device, req.GetCommand(), req.GetData())
		}
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
	res.UserId = device.UserId
	res.Data = data
	res.Command = req.Command
	return nil
}

// 单独处理登录
func loginService(device logic_model.ConnDevice, data []byte) (res *pb.PacketLoginRes, replyErr error) {
	var params pb.PacketLoginReq
	if replyErr = proto.Unmarshal(data, &params); replyErr == nil {
		return Login.Login(device, &params)
	}
	return
}

var execTables = map[pb.Packet]func(logic_model.ConnDevice, []byte) (msg proto.Message, err error){
	//基础&账号
	pb.Packet_GetProfile: func(device logic_model.ConnDevice, data []byte) (msg proto.Message, err error) { //获取当前账号信息
		var params pb.PacketGetProfileReq
		if err = proto.Unmarshal(data, &params); err == nil {
			msg, err = User.GetProfile(device, &params)
		}
		return
	},
	pb.Packet_ModifyProfile: func(device logic_model.ConnDevice, data []byte) (msg proto.Message, err error) {
		var params pb.PacketModfiyProfileReq
		if err = proto.Unmarshal(data, &params); err == nil {
			msg, err = User.ModifyProfile(device, &params)
		}
		return
	},
	//用户相关
	pb.Packet_GetUser: func(device logic_model.ConnDevice, data []byte) (msg proto.Message, err error) {
		var params pb.PacketGetUserInfoReq
		if err = proto.Unmarshal(data, &params); err == nil {
			msg, err = User.GetUser(device, &params)
		}
		return
	},
	//好友相关
	pb.Packet_GetFriendList: func(device logic_model.ConnDevice, data []byte) (msg proto.Message, err error) {
		var params pb.PacketGetFriendListReq
		if err = proto.Unmarshal(data, &params); err == nil {
			msg, err = Friend.GetFriendList(device, &params)
		}
		return
	},
	pb.Packet_GetFriend: func(device logic_model.ConnDevice, data []byte) (msg proto.Message, err error) {
		var params pb.PacketGetFriendReq
		if err = proto.Unmarshal(data, &params); err == nil {
			msg, err = Friend.GetFriend(device, &params)
		}
		return
	},
	pb.Packet_RemoveFriend: func(device logic_model.ConnDevice, data []byte) (msg proto.Message, err error) {
		var params pb.PacketRemoveFriendReq
		if err = proto.Unmarshal(data, &params); err == nil {
			msg, err = Friend.RemoveFriend(device, &params)
		}
		return
	},
	//发起好友申请
	pb.Packet_FriendApply: func(device logic_model.ConnDevice, data []byte) (msg proto.Message, err error) {
		var params pb.PacketFriendApplyReq
		if err = proto.Unmarshal(data, &params); err == nil {
			msg, err = Friend.FriendApply(device, &params)
		}
		return
	},
	//处理好友申请
	pb.Packet_FriendHandle: func(device logic_model.ConnDevice, data []byte) (msg proto.Message, err error) {
		var params pb.PacketFriendHandleReq
		if err := proto.Unmarshal(data, &params); err == nil {
			msg, err = Friend.FriendHandle(device, &params)
		}
		return
	},
	//群组相关
	//私聊
	pb.Packet_PrivateMsg: func(device logic_model.ConnDevice, data []byte) (msg proto.Message, err error) {
		var params pb.PacketPrivateMsg
		if err = proto.Unmarshal(data, &params); err == nil {
			err = Friend.SendMsg(device, &params)
		}
		return
	},
	pb.Packet_PrivateMsgList: func(device logic_model.ConnDevice, data []byte) (msg proto.Message, err error) {
		var params pb.PacketPrivateMsgListReq
		if err = proto.Unmarshal(data, &params); err == nil {
			msg, err = Friend.GetMessageList(device, &params)
		}
		return
	},
}

func routeService(device logic_model.ConnDevice, commandId pb.Packet, data []byte) (proto.Message, error) {
	fn, has := execTables[commandId]
	if has {
		fmt.Println("路由控制器指向: ", commandId)
		return fn(device, data)
	}
	fmt.Println("未知的操作指向: ", commandId)
	return nil, gerror.Newf("未知操作命令")
	//
	//switch commandId {
	//case pb.Packet_GetProfile: //获取当前账号信息
	//	{
	//		var params pb.PacketGetProfileReq
	//		if replyErr = proto.Unmarshal(data, &params); replyErr == nil {
	//			replyMsg, replyErr = User.GetProfile(device, &params)
	//		}
	//		return
	//	}
	//case pb.Packet_ModifyProfile: //修改当前账号信息
	//	{
	//		var params pb.PacketModfiyProfileReq
	//		if replyErr = proto.Unmarshal(data, &params); replyErr == nil {
	//			replyMsg, replyErr = User.ModifyProfile(device, &params)
	//		}
	//		return
	//	}
	//case pb.Packet_GetFriendList: //获取好友列表
	//	{
	//		var params pb.PacketGetFriendListReq
	//		if replyErr = proto.Unmarshal(data, &params); replyErr == nil {
	//			replyMsg, replyErr = Friend.GetFriendList(device, &params)
	//		}
	//		return
	//	}
	//case pb.Packet_GetFriend: //获取好友信息
	//	{
	//		var params pb.PacketGetFriendReq
	//		if replyErr = proto.Unmarshal(data, &params); replyErr == nil {
	//			replyMsg, replyErr = Friend.GetFriend(device, &params)
	//		}
	//		return
	//	}
	//case pb.Packet_PrivateMsg: //发送私聊消息
	//	{
	//		var params pb.PacketPrivateMsg
	//		if replyErr = proto.Unmarshal(data, &params); replyErr == nil {
	//			replyErr = Friend.SendMsg(device, &params)
	//		}
	//		return
	//	}
	//case pb.Packet_PrivateMsgList: //获取私聊消息列表
	//	{
	//		var params pb.PacketPrivateMsgListReq
	//		if replyErr = proto.Unmarshal(data, &params); replyErr == nil {
	//			replyMsg, replyErr = Friend.GetMessageList(device, &params)
	//		}
	//		return
	//	}
	//case pb.Packet_GetUser:
	//	{
	//		var params pb.PacketGetUserInfoReq
	//		if replyErr = proto.Unmarshal(data, &params); replyErr == nil {
	//			replyMsg, replyErr = User.GetUser(device, &params)
	//		}
	//		return
	//	}
	//}
	return nil, gerror.Newf("未知操作命令")
}
