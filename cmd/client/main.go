package main

import (
	"XtTalkServer/pb"
	"XtTalkServer/services/connect/types"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"sync/atomic"
	"time"
)

var conn net.Conn = nil
var seq uint32 = 0 //序列ID
var uid uint64 = 0

// main 测试客户端
func main() {
	ctx := gctx.New()

	uid = gvar.New(os.Args[1]).Uint64()
	addr := gvar.New(os.Args[2]).String()

	var err error = nil
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		glog.Fatalf(ctx, "连接服务器失败: %s", err.Error())
	}
	glog.Infof(ctx, "连接服务器成功")
	go func() {
		for {
			var headBytes = make([]byte, types.DataPackHeaderLength)
			if _, err := conn.Read(headBytes); err != nil {
				glog.Fatalf(ctx, "读取数据失败: %s", err.Error())
				break
			}
			var imHead types.ImHeadDataPack
			buffer := bytes.NewBuffer(headBytes)
			if err := binary.Read(buffer, binary.LittleEndian, &imHead.ProtocolVersion); err != nil {
				glog.Fatalf(ctx, "解包失败ProtocolVersion: %s", err.Error())
			}
			if err := binary.Read(buffer, binary.LittleEndian, &imHead.Command); err != nil {
				glog.Fatalf(ctx, "解包失败Command: %s", err.Error())
			}
			if err := binary.Read(buffer, binary.LittleEndian, &imHead.Sequence); err != nil {
				glog.Fatalf(ctx, "解包失败Command: %s", err.Error())
			}
			if err := binary.Read(buffer, binary.LittleEndian, &imHead.Length); err != nil {
				glog.Fatalf(ctx, "解包失败Length: %s", err.Error())
			}
			var dataBytes = make([]byte, imHead.Length)
			if _, err := conn.Read(dataBytes); err != nil {
				glog.Fatalf(ctx, "读取内容数据失败: %s", err.Error())
				break
			}
			onMessage(imHead, dataBytes)
		}
	}()
	fmt.Println("请输入操作命令")
	for {
		reader := bufio.NewReader(os.Stdin)

		res, _, err := reader.ReadLine()
		if nil != err {
			fmt.Println("请重新输入:", err.Error())
		}
		var msg = string(res)
		params := gstr.Split(msg, " ")
		if 0 >= len(params) {
			fmt.Println("请重新输入!")
			continue
		}
		fmt.Println("接收参数:", params[0], params[1:])
		input(params[0], params[1:])
	}
}
func input(msg string, params []string) {
	switch msg {
	case "login": //登录
		fmt.Println("发送结果:", sendPack(pb.Packet_Login, &pb.PacketLoginReq{
			Token: gvar.New(uid).String(),
		}))
	case "getprofile": //获取个人信息
		fmt.Println("发送结果:", sendPack(pb.Packet_GetProfile, &pb.PacketGetProfileReq{}))
	case "modfiyprofile": //修改个人信息
		fmt.Println("发送结果:", sendPack(pb.Packet_ModifyProfile, &pb.PacketModfiyProfileReq{
			NickName: fmt.Sprintf("幻音%d", time.Now().Unix()),
			Age:      2,
			Sex:      3,
			Note:     fmt.Sprintf("签名%d", time.Now().Unix()),
		}))
	case "getfriendlist": //获取好友列表
		fmt.Println("发送结果:", sendPack(pb.Packet_GetFriendList, &pb.PacketGetFriendListReq{
			Page: 1,
			Size: 2,
		}))
	case "getfriend": //获取好友信息
		sendUid := 2
		if uid == 2 {
			sendUid = 1
		}
		fmt.Println("发送结果:", sendPack(pb.Packet_GetFriend, &pb.PacketGetFriendReq{
			UserId: uint64(sendUid),
		}))
	case "sendmsg":
		textMsg, _ := proto.Marshal(&pb.TextMsg{
			Content: fmt.Sprintf("我是 %d,当前时间是: %s", uid, gtime.Now().Format("Y-m-d H:i:s")),
		})
		sendUid := 2
		if uid == 2 {
			sendUid = 1
		}
		fmt.Println("发送结果:", sendPack(pb.Packet_PrivateMsg, &pb.PacketPrivateMsg{
			MsgId:      time.Now().UnixNano(), //消息ID
			FromId:     uid,
			ReceiveId:  uint64(sendUid),
			MsgType:    pb.PacketMsgType_Text,
			Payload:    textMsg,
			ClientTime: time.Now().Unix(),
		}))
	}
}
func sendPack(commandId pb.Packet, pb proto.Message) error {
	_bytes, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	sendSeq := atomic.AddUint32(&seq, 1)
	head := types.ImHeadDataPack{
		ProtocolVersion: types.ProtocolVersion,
		Command:         uint16(commandId),
		Sequence:        sendSeq,
		Length:          uint32(len(_bytes)),
	}
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, head.ProtocolVersion); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, head.Command); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, head.Sequence); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, head.Length); err != nil {
		return err
	}
	resultBytes := append(buf.Bytes(), _bytes...)
	_, err = conn.Write(resultBytes)
	return err
}
func onMessage(head types.ImHeadDataPack, data []byte) {
	ctx := gctx.New()
	glog.Infof(ctx, "协议版本: %d", head.ProtocolVersion)
	glog.Infof(ctx, "数据命令: %d", head.Command)
	glog.Infof(ctx, "请求序列: %d", head.Sequence)
	glog.Infof(ctx, "数据长度: %d", head.Length)
	glog.Infof(ctx, "PB数据内容  : %v", data)
	glog.Debugf(ctx, "-------------------------------------------------")

	switch pb.Packet(head.Command) {
	case pb.Packet_Login:
		{
			glog.Infof(ctx, "登录成功")
		}
	}

}

/**
****主动事件
登录
---个人信息
获取当前登录账号信息
修改当前登录账号信息

---用户
获取陌生用户信息


---好友
获取好友列表
获取好友信息
删除好友
修改好友信息

---私聊
发送私聊消息

---群聊
获取群组列表
获取群信息
获取群成员列表
获取群成员信息
退出群组
踢出群成员

*/
