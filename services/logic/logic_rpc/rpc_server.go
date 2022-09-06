package logic_rpc

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/global"
	"XtTalkServer/pb"
	"XtTalkServer/services/connect/types"
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
		reply, err := routeService(device, types.DataFormat(req.GetDataFormat()), req.GetCommand(), req.Data)
		fmt.Println("调用结果:", reply, err)
		if err != nil {
			return nil, err
		}
		if reply == nil {
			return nil, nil
		}

		var dataBytes []byte = nil
		switch types.DataFormat(req.GetDataFormat()) {
		//pb格式转换
		case types.DataFormatPb:
			{
				if pb, has := reply.(proto.Message); has {
					dataBytes, err = proto.Marshal(pb)
				} else {
					err = gerror.Newf("数据格式无法解析")
				}
			}
		//json格式转换
		case types.DataFormatJson:
			{
				dataBytes, err = gjson.Encode(reply)
			}
		}
		if err != nil {
			return nil, err
		}
		return dataBytes, nil
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

func routeService(device logic_model.ConnDevice, format types.DataFormat, commandId pb.Packet, data []byte) (interface{}, error) {
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
			if err := parseData(format, data, &params); err != nil {
				return nil, err
			}
			return Login.Login(device, &params)
		}
	case pb.Packet_GetProfile: //获取当前账号信息
		{
			var params pb.PacketGetProfileReq
			if err := parseData(format, data, &params); err != nil {
				return nil, err
			}
			return Self.GetProfile(device, &params)
		}
	case pb.Packet_ModifyProfile: //修改当前账号信息
		{
			var params pb.PacketModfiyProfileReq
			if err := parseData(format, data, &params); err != nil {
				return nil, err
			}
			return Self.ModifyProfile(device, &params)
		}
	case pb.Packet_GetFriendList: //获取好友列表
		{
			var params pb.PacketGetFriendListReq
			if err := parseData(format, data, &params); err != nil {
				return nil, err
			}
			return Friend.GetFriendList(device, &params)
		}
	case pb.Packet_GetFriend: //获取好友信息
		{
			var params pb.PacketGetFriendReq
			if err := parseData(format, data, &params); err != nil {
				return nil, err
			}
			return Friend.GetFriend(device, &params)
		}
	case pb.Packet_PrivateMsg:
		{
			var params pb.PacketPrivateMsg
			if err := parseData(format, data, &params); err != nil {
				return nil, err
			}
			if err := Friend.SendMsg(device, &params); err != nil {
				return nil, err
			}
			return nil, nil
		}
	}

	return nil, gerror.Newf("未知操作命令")
}
func parseData(format types.DataFormat, bytes []byte, v proto.Message) error {
	switch format {
	case types.DataFormatJson:
		if err := gjson.Unmarshal(bytes, v); err != nil {
			return err
		}
	case types.DataFormatPb:
		if err := proto.Unmarshal(bytes, v); err != nil {
			return err
		}
	default:
		return gerror.Newf("不支持的数据协议")
	}
	return nil
}
