package main

import (
	"XtTalkServer/cmd/client/sdk"
	"XtTalkServer/pb"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"time"
)

type ImClientController struct {
	client *sdk.XtTalkClient
}

// 通过反射调用
func (c *ImClientController) Login(args []string) error {
	data, err := c.client.SendAndWaitPacket(pb.Packet_Login, &pb.PacketLoginReq{
		Token: gvar.New(uid).String(),
	})
	if err != nil {
		return err
	}
	var msg pb.PacketLoginRes
	proto.Unmarshal(data.Body, &msg)
	dumpTable(&msg)
	return nil
}
func (c *ImClientController) GetProfile(args []string) error {
	data, err := c.client.SendAndWaitPacket(pb.Packet_GetProfile, &pb.PacketGetProfileReq{})
	if err != nil {
		return err
	}
	var msg pb.PacketGetProfileRes
	proto.Unmarshal(data.Body, &msg)
	dumpTable(&msg)
	return nil
}
func (c *ImClientController) ModfiyProfile(args []string) error {
	data, err := c.client.SendAndWaitPacket(pb.Packet_ModifyProfile, &pb.PacketModfiyProfileReq{
		NickName: fmt.Sprintf("%s%d", name, time.Now().Unix()),
		Age:      2,
		Sex:      3,
		Note:     fmt.Sprintf("签名%d", time.Now().Unix()),
	})
	if err != nil {
		return err
	}
	var msg pb.PacketModfiyProfileRes
	proto.Unmarshal(data.Body, &msg)
	dumpTable(&msg)
	return nil
}
func (c *ImClientController) GetFriendList(args []string) error {
	data, err := c.client.SendAndWaitPacket(pb.Packet_GetFriendList, &pb.PacketGetFriendListReq{
		Page: 1,
		Size: 2,
	})
	if err != nil {
		return err
	}
	var msg pb.PacketGetFriendListRes
	proto.Unmarshal(data.Body, &msg)
	dumpTable(&msg)
	return nil
}
func (c *ImClientController) GetFriend(args []string) error {
	friendId := 2
	if uid == 2 {
		friendId = 1
	}
	data, err := c.client.SendAndWaitPacket(pb.Packet_GetFriend, &pb.PacketGetFriendReq{
		UserId: uint64(friendId),
	})
	if err != nil {
		return err
	}
	var msg pb.PacketGetFriendRes
	proto.Unmarshal(data.Body, &msg)
	dumpTable(&msg)
	return nil
}
func (c *ImClientController) SendPrivateMsg(args []string) error {
	seq := c.client.NextPrivateSeq()
	textMsg, _ := proto.Marshal(&pb.TextMsg{
		Content: fmt.Sprintf("我是 %d -> %s,当前时间是: %s", uid, name, gtime.Now().Format("Y-m-d H:i:s")),
	})
	friendId := 2
	if uid == 2 {
		friendId = 1
	}
	return c.client.SendPacket(pb.Packet_PrivateMsg, &pb.PacketPrivateMsg{
		MsgSeq:     int64(seq),
		MsgRand:    rand.Int63(), //随机id
		FromId:     uid,
		ReceiveId:  uint64(friendId),
		MsgType:    pb.PacketMsgType_Text,
		Payload:    textMsg,
		ClientTime: time.Now().Unix(),
	})
}

// GetUser 获取用户,搜索,精准
func (c *ImClientController) GetUser(args []string) error {
	data, err := c.client.SendAndWaitPacket(pb.Packet_GetUser, &pb.PacketGetUserInfoReq{
		Username: "abc",
	})
	if err != nil {
		return err
	}
	var msg pb.PacketGetUserInfoResp
	proto.Unmarshal(data.Body, &msg)
	dumpTable(&msg)
	return nil
}

// FriendApply 发送好友申请
func (c *ImClientController) FriendApply(args []string) error {
	var friendId uint64 = 2
	if uid == 2 {
		friendId = 1
	}
	data, err := c.client.SendAndWaitPacket(pb.Packet_FriendApply, &pb.PacketFriendApplyReq{
		UserId: friendId,
		Reason: "加我好友啊 尼玛的",
	})
	if err != nil {
		return err
	}
	var msg pb.PacketFriendApplyResp
	proto.Unmarshal(data.Body, &msg)
	dumpTable(&msg)
	return nil
}

// FriendReject 发送好友拒绝申请
func (c *ImClientController) FriendReject(args []string) error {
	data, err := c.client.SendAndWaitPacket(pb.Packet_FriendHandle, &pb.PacketFriendHandleReq{
		Id:           37,
		Flag:         false,
		RejectReason: "拒绝添加",
	})
	if err != nil {
		return err
	}
	var msg pb.PacketFriendHandleResp
	proto.Unmarshal(data.Body, &msg)
	dumpTable(&msg)
	return nil
}

// FriendAccept 发送好友同意申请
func (c *ImClientController) FriendAccept(args []string) error {
	data, err := c.client.SendAndWaitPacket(pb.Packet_FriendHandle, &pb.PacketFriendHandleReq{
		Id:   37,
		Flag: true,
	})
	if err != nil {
		return err
	}
	var msg pb.PacketFriendHandleResp
	proto.Unmarshal(data.Body, &msg)
	dumpTable(&msg)
	return nil
}

// RemoveFriend 删除好友
func (c *ImClientController) RemoveFriend(args []string) error {
	var friendId uint64 = 2
	if uid == 2 {
		friendId = 1
	}
	data, err := c.client.SendAndWaitPacket(pb.Packet_RemoveFriend, &pb.PacketRemoveFriendReq{
		UserId: friendId,
	})
	if err != nil {
		return err
	}
	var msg pb.PacketRemoveFriendResp
	proto.Unmarshal(data.Body, &msg)
	dumpTable(&msg)
	return nil
}

func (c *ImClientController) GetPrivateMsg(args []string) error {
	friendId := 2
	if uid == 2 {
		friendId = 1
	}
	data, err := c.client.SendAndWaitPacket(pb.Packet_PrivateMsgList, &pb.PacketPrivateMsgListReq{
		UserId: uint64(friendId),
		Size:   2,
		//LastMsgId: 595109068935270400,
	})
	if err != nil {
		return err
	}
	var msg pb.PacketPrivateMsgListResp
	proto.Unmarshal(data.Body, &msg)
	fmt.Println("消息列表:")
	for index, item := range msg.GetList() {
		fmt.Println(fmt.Sprintf("[%d / %d]  [%v]  [%s]  [%v]", index+1, len(msg.GetList()), item.MsgId, item.Payload, item.ServerTime))
		fmt.Println("-------------------------------------------------------------")
	}

	return nil
}
