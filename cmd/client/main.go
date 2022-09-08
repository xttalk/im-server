package main

import (
	"XtTalkServer/cmd/client/sdk"
	"XtTalkServer/pb"
	"XtTalkServer/services/connect/types"
	"bufio"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/liushuochen/gotable"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"os"
	"time"
)

var uid uint64 = 0
var name string = ""

// main 测试客户端
func main() {
	uid = gvar.New(os.Args[1]).Uint64()
	addr := gvar.New(os.Args[2]).String()
	if uid == 1 {
		name = "幻音"
	} else {
		name = "夏花"
	}
	for {
		ctx := gctx.New()
		if err := ClientExec(addr); err != nil {
			glog.Errorf(ctx, "执行失败: %s", err.Error())
		}
		time.Sleep(time.Second * 3)
	}
}

var client *sdk.XtTalkClient
var isStartCommander = false

func ClientExec(addr string) error {
	client = sdk.CreateClient()
	if err := client.Connect(addr); err != nil {
		return err
	}
	fmt.Println("连接服务器成功: ", addr)
	if !isStartCommander {
		isStartCommander = true
		go ListenCommand(ListenCommander)
	}
	if err := client.ListenReader(OnMessage); err != nil {
		return err
	}
	return nil
}
func ListenCommand(callback func(context.Context, string, []string) (bool, error)) {
	fmt.Println("请输入命令...")
	for {
		reader := bufio.NewReader(os.Stdin)
		bytes, _, err := reader.ReadLine()
		if nil != err {
			fmt.Println("请重新输入:", err.Error())
		}
		params := gstr.Split(string(bytes), " ")
		ctx := gctx.New()
		command := params[0]
		callback(ctx, command, params[1:])
	}
}
func OnMessage(ctx context.Context, pack types.ImHeadDataPack, bytes []byte) error {
	//包头输出
	headTable, err := gotable.Create("协议版本", "命令名称", "命令码", "序列号", "数据长度")
	if err != nil {
		glog.Warningf(ctx, "解析数据Header失败: %s", err.Error())
		return nil
	}
	headTable.AddRow([]string{gvar.New(pack.ProtocolVersion).String(), gvar.New(pb.Packet(pack.Command)).String(), gvar.New(pack.Command).String(), gvar.New(pack.Sequence).String(), gvar.New(pack.Length).String()})
	fmt.Println(headTable)
	dumpTable := func(msg interface{}) {
		if msg == nil {
			fmt.Println("无数据可显示")
			return
		}
		g.Dump(msg)
		//bodyTable, err := gotable.CreateByStruct(msg)
		//if err != nil {
		//	glog.Warningf(ctx, "解析数据Body失败: %s", err.Error())
		//	return
		//}
		//fmt.Println(bodyTable)
	}
	switch pb.Packet(pack.Command) {
	case pb.Packet_Login:
		var msg pb.PacketLoginRes
		proto.Unmarshal(bytes, &msg)
		dumpTable(&msg)
	case pb.Packet_GetProfile:
		var msg pb.PacketGetProfileRes
		proto.Unmarshal(bytes, &msg)
		dumpTable(&msg)
	case pb.Packet_ModifyProfile:
		var msg pb.PacketModfiyProfileRes
		proto.Unmarshal(bytes, &msg)
		dumpTable(&msg)
	case pb.Packet_GetFriend:
		var msg pb.PacketGetFriendRes
		proto.Unmarshal(bytes, &msg)
		dumpTable(&msg)
	case pb.Packet_GetFriendList:
		var msg pb.PacketGetFriendListRes
		proto.Unmarshal(bytes, &msg)
		dumpTable(&msg)
	default:
		glog.Warningf(ctx, "暂未支持的命令码: %s", bytes)
	}

	return nil
}
func ListenCommander(ctx context.Context, command string, args []string) (isTrigger bool, err error) {
	switch command {
	case "login": //登录
		err = client.SendPacket(pb.Packet_Login, &pb.PacketLoginReq{
			Token: gvar.New(uid).String(),
		})
	case "getprofile": //获取个人信息
		err = client.SendPacket(pb.Packet_GetProfile, &pb.PacketGetProfileReq{})
	case "modfiyprofile": //修改信息
		err = client.SendPacket(pb.Packet_ModifyProfile, &pb.PacketModfiyProfileReq{
			NickName: fmt.Sprintf("%s%d", name, time.Now().Unix()),
			Age:      2,
			Sex:      3,
			Note:     fmt.Sprintf("签名%d", time.Now().Unix()),
		})
	case "getfriendlist": //获取好友列表
		err = client.SendPacket(pb.Packet_GetFriendList, &pb.PacketGetFriendListReq{
			Page: 1,
			Size: 2,
		})
	case "getfriend": //获取好友信息
		friendId := 2
		if uid == 2 {
			friendId = 1
		}
		err = client.SendPacket(pb.Packet_GetFriend, &pb.PacketGetFriendReq{
			UserId: uint64(friendId),
		})
	case "sendmsg":
		textMsg, _ := proto.Marshal(&pb.TextMsg{
			Content: fmt.Sprintf("我是 %d -> %s,当前时间是: %s", uid, name, gtime.Now().Format("Y-m-d H:i:s")),
		})
		friendId := 2
		if uid == 2 {
			friendId = 1
		}
		err = client.SendPacket(pb.Packet_PrivateMsg, &pb.PacketPrivateMsg{
			MsgId:      time.Now().UnixNano(), //消息ID0
			FromId:     uid,
			ReceiveId:  uint64(friendId),
			MsgType:    pb.PacketMsgType_Text,
			Payload:    textMsg,
			ClientTime: time.Now().Unix(),
		})
	default:
		glog.Warningf(ctx, "未知命令: %s -> %v", command, args)
		return false, nil
	}
	glog.Infof(ctx, "执行命令: %s -> 参数: %v", command, args)
	return
}

//
///**
//****主动事件
//登录
//---个人信息
//获取当前登录账号信息
//修改当前登录账号信息
//
//---用户
//获取陌生用户信息
//
//
//---好友
//获取好友列表
//获取好友信息
//删除好友
//修改好友信息
//
//---私聊
//发送私聊消息
//
//---群聊
//获取群组列表
//获取群信息
//获取群成员列表
//获取群成员信息
//退出群组
//踢出群成员
//
//*/
