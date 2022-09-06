package logic_rpc

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/global"
	"XtTalkServer/pb"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/glog"
	"google.golang.org/protobuf/proto"
)

var UserClient = new(_UserClient)

type _UserClient struct {
}

// SendUserPacket 向指定用户发送数据
func (_UserClient) SendUserPacket(ctx context.Context, userId uint64, command pb.Packet, msg proto.Message) error {
	glog.Infof(ctx, "向用户: %d 发送数据", userId)
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	//找到这个用户的所有session
	userSessionList := global.Redis.HGetAll(ctx, fmt.Sprintf(conts.RK_UserDevice, userId)).Val()
	for session, sid := range userSessionList {
		serverId := gvar.New(sid).Uint32()
		var req = pb.ConnectSendClientReq{
			ServerId:  serverId,
			SessionId: session,
			Command:   command,
			Payload:   bytes,
		}
		var res = pb.ConnectSendClientRes{}
		if err := ConnectRpcClient.Call(ctx, "SendClientPacket", &req, &res); err != nil {
			glog.Warningf(ctx, "推送消息失败: 服务器ID:[%d][%s] %s", serverId, session, err.Error())
		}
	}
	return nil
}
