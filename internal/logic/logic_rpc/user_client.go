package logic_rpc

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/global"
	"XtTalkServer/pb"
	"XtTalkServer/utils"
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/smallnest/rpcx/client"
	"google.golang.org/protobuf/proto"
)

type PacketSendInfo struct {
	UserId         uint64   //目标发送用户
	SendSession    []string //指定发送的session列表,nil则代表全部设备
	ExcludeSession []string //不发送的session列表
}

var UserClient = new(_UserClient)

type _UserClient struct {
}

// SendUserPacket 向指定用户发送数据
func (c *_UserClient) SendUserPacket(ctx context.Context, info PacketSendInfo, command pb.Packet, msg proto.Message) error {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return c.SendUserBytes(ctx, info, command, bytes)
}
func (c *_UserClient) SendUserBytes(ctx context.Context, info PacketSendInfo, command pb.Packet, bytes []byte) error {
	glog.Infof(ctx, "向用户: %d 发送数据: 允许session: %#v   ->    排除session: %#v", info.UserId, info.SendSession, info.ExcludeSession)
	//找到这个用户的所有session
	userSessionList := global.Redis.HGetAll(ctx, fmt.Sprintf(conts.RK_UserDevice, info.UserId)).Val()
	rclient, err := RpcClient.GetConnectClient(ctx)

	if err != nil {
		return err
	}
	go func() {
		for session, sid := range userSessionList {
			serverId := gvar.New(sid).Uint32()
			if info.SendSession != nil && len(info.SendSession) > 0 {
				//判断是否是需要发送的
				if !utils.InArray(session, info.SendSession) {
					//glog.Debugf(ctx, "[%d]->[%s]->[%d] 没有处于允许session中本次跳过发送", info.UserId, session, serverId)
					continue //跳过发送
				}
			}
			if info.ExcludeSession != nil && utils.InArray(session, info.ExcludeSession) {
				//glog.Debugf(ctx, "[%d]->[%s]->[%d] 处于排除session中本次跳过发送", info.UserId, session, serverId)
				continue //跳过发送
			}

			var req = pb.ConnectSendClientReq{
				ServerId:  serverId,
				SessionId: session,
				Command:   command,
				Payload:   bytes,
			}
			var res pb.ConnectSendClientRes
			var isDeleteSession = false
			if err := rclient.Call(ctx, "SendClientPacket", &req, &res); err != nil {
				if errors.Is(err, client.ErrXClientNoServer) {
					isDeleteSession = true
				} else {
					glog.Warningf(ctx, "推送消息失败: 服务器ID:[%d][%s] %s", serverId, session, err.Error())
				}
			}
			if res.RetCode == pb.ConnectRetCode_CR_Offline { //session离线,也删除
				isDeleteSession = true
			} else {
				//glog.Debugf(ctx, "[%d]->[%s]->[%d] 推送完成: %s", info.UserId, session, serverId, gvar.New(res.RetCode).String())
			}
			if isDeleteSession {
				//目标服务器不存在,移除这个session,并且移除user
				global.Redis.Del(ctx, fmt.Sprintf(conts.RK_ClientAuth, session))
				global.Redis.HDel(ctx, fmt.Sprintf(conts.RK_UserDevice, info.UserId), session)
			}
		}
	}()
	return nil
}
