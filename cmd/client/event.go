package main

import (
	"XtTalkServer/pb"
	"context"
	"github.com/gogf/gf/v2/os/glog"
)

// 收到私聊消息回执
func PrivateMsgAckEvent(ctx context.Context, msg *pb.PacketPrivateMsgAck) {
	glog.Infof(ctx, "收到消息回执: 【%v】【%v】【%v】", msg.MsgSeq, msg.MsgRand, msg.RetCode)
	dumpTable(msg)
}

// 收到私聊消息
func PrivateMsgEvent(ctx context.Context, msg *pb.PacketPrivateMsg) {
	glog.Infof(ctx, "收到私聊消息:【%v】【%v】【%s】", msg.MsgSeq, msg.MsgRand, msg.Payload)
}

// 收到好友验证消息
func FriendRequestEvent(ctx context.Context, msg *pb.FriendRequestEvent) {
	switch msg.Status {
	case 1:
		glog.Infof(ctx, "收到好友验证消息: [%v] [%v]->[%v] : %s", msg.GetId(), msg.GetFromUid(), msg.GetToUid(), msg.GetReason())
	case 2:
		glog.Infof(ctx, "收到好友验证消息: [%v] [%v]同意了[%v]的好友申请", msg.GetId(), msg.GetToUid(), msg.GetFromUid())
	case 3:
		glog.Infof(ctx, "收到好友验证消息: [%v] [%v]拒绝了[%v]的好友申请: %s", msg.GetId(), msg.GetToUid(), msg.GetFromUid(), msg.GetRejectReason())
	}
}

// 收到好友变动
func FriendRequestChange(ctx context.Context, msg *pb.FriendChangeEvent) {
	if msg.GetIsFriend() {
		glog.Infof(ctx, "收到好友验证变动: 好友新增: %v -> %s", msg.FriendId, msg.Friend.Nickname)
	} else {
		glog.Infof(ctx, "收到好友验证变动: 好友删除: %v", msg.FriendId)
	}

}
