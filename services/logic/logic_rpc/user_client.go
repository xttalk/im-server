package logic_rpc

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/global"
	"XtTalkServer/pb"
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/smallnest/rpcx/client"
)

var UserClient = new(_UserClient)

type _UserClient struct {
}

// SendUserPacket 向指定用户发送数据
func (_UserClient) SendUserPacket(ctx context.Context, userId uint64, command pb.Packet, bytes []byte) error {
	glog.Infof(ctx, "向用户: %d 发送数据", userId)
	//找到这个用户的所有session
	userSessionList := global.Redis.HGetAll(ctx, fmt.Sprintf(conts.RK_UserDevice, userId)).Val()
	rpcClient, err := RpcClient.GetConnectClient(ctx)
	if err != nil {
		return err
	}
	go func() {
		for session, sid := range userSessionList {
			serverId := gvar.New(sid).Uint32()
			var req = pb.ConnectSendClientReq{
				ServerId:  serverId,
				SessionId: session,
				Command:   command,
				Payload:   bytes,
			}
			var isDeleteSession = false
			var res = pb.ConnectSendClientRes{}
			if err := rpcClient.Call(ctx, "SendClientPacket", &req, &res); err != nil {
				if errors.Is(err, client.ErrXClientNoServer) {
					isDeleteSession = true
				} else {
					glog.Warningf(ctx, "推送消息失败: 服务器ID:[%d][%s] %s", serverId, session, err.Error())
				}
			}
			if res.RetCode == pb.ConnectRetCode_CR_Offline { //session离线,也删除
				isDeleteSession = true
			}
			if isDeleteSession {
				//目标服务器不存在,移除这个session,并且移除user
				global.Redis.Del(ctx, fmt.Sprintf(conts.RK_ClientAuth, session))
				global.Redis.HDel(ctx, fmt.Sprintf(conts.RK_UserDevice, userId), session)
			}
		}
	}()
	return nil
}
