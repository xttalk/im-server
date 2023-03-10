package main

import (
	"XtTalkServer/cmd/client/sdk"
	"XtTalkServer/core"
	"XtTalkServer/internal"
	"XtTalkServer/internal/connect/types"
	"XtTalkServer/pb"
	"fmt"
	inf "github.com/fzdwx/infinite"
	"github.com/fzdwx/infinite/components/selection/singleselect"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/liushuochen/gotable"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"os"
	"reflect"
	"strings"
	"time"
)

var uid uint64 = 0
var name string = ""

// main 测试客户端
func main() {
	core.Entry.ExecStart(gctx.New(), []internal.InitCtx{})
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
var clientController *ImClientController

func ClientExec(addr string) error {
	client = sdk.CreateClient()
	if err := client.Connect(addr); err != nil {
		return err
	}
	fmt.Println("连接服务器成功: ", addr)
	clientController = &ImClientController{
		client: client,
	}
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
	commandList := make([]string, 0)
	ref := reflect.TypeOf(clientController)
	for i := 0; i < ref.NumMethod(); i++ {
		commandList = append(commandList, ref.Method(i).Name)
	}
	for {
		ctx := gctx.New()
		selected, err := inf.NewSingleSelect(
			commandList,
			singleselect.WithDisableFilter(), //禁止输入
			singleselect.WithPageSize(100),
		).Display("请选择需要操作的命令")

		if err != nil {
			fmt.Println("选择失败: ", err.Error())
			continue
		}
		if selected >= len(commandList) {
			fmt.Println("没有这个命令: ", err.Error())
			continue
		}
		cmd := commandList[selected]
		_, err = callback(ctx, cmd, []string{})
		if err != nil {
			glog.Errorf(ctx, "调用%s失败: %s", cmd, err.Error())
		}
	}

}
func dumpTable(msg proto.Message) {
	if msg == nil {
		fmt.Println("无数据可显示")
		return
	}
	titleList := make([]string, 0)
	valueList := make([]string, 0)
	ref := reflect.TypeOf(msg)
	for i := 0; i < ref.NumMethod(); i++ {
		rex, _ := gregex.MatchString("^Get(.*)", ref.Method(i).Name)
		if len(rex) >= 2 {
			ret := ref.Method(i).Func.Call([]reflect.Value{reflect.ValueOf(msg)})
			if len(ret) >= 1 {
				titleList = append(titleList, rex[1])
				valueList = append(valueList, fmt.Sprintf("%v", gvar.New(ret[0]).Interface()))
			}
		}
	}
	if 0 >= len(titleList) {
		fmt.Println("无详细数据展示")
		return
	}
	bodyTable, err := gotable.Create(titleList...)
	if err != nil {
		fmt.Println("详细数据构建失败")
		return
	}
	if err := bodyTable.AddRow(valueList); err != nil {
		fmt.Println("详细数据构建内容失败: ", err.Error())
		return
	}
	fmt.Println(bodyTable)
	//fmt.Println("field", .NumMethod())
	//g.Dump(msg)
	//
	//if err != nil {
	//	glog.Warningf(ctx, "解析数据Body失败: %s", err.Error())
	//	return
	//}
	//fmt.Println(bodyTable)
}
func OnMessage(ctx context.Context, pack types.FixedHeader, bytes []byte) error {
	//包头输出
	headTable, err := gotable.Create("协议版本", "命令名称", "命令码", "序列号", "数据长度")
	if err != nil {
		glog.Warningf(ctx, "解析数据Header失败: %s", err.Error())
		return nil
	}
	headTable.AddRow([]string{gvar.New(pack.Version).String(), gvar.New(pb.Packet(pack.Command)).String(), gvar.New(pack.Command).String(), gvar.New(pack.Sequence).String(), gvar.New(pack.Length).String()})
	fmt.Println(headTable)

	//判断是否是seq序列码
	fmt.Println("收到seq", pack.Sequence)
	if fn, has := client.SeqHandler.LoadAndDelete(pack.Sequence); has {
		if f, has := fn.(func(types.FixedHeader, []byte)); has {
			f(pack, bytes)
		}
		return nil
	}
	switch pb.Packet(pack.Command) {
	case pb.Packet_PrivateMsgAck:
		var msg pb.PacketPrivateMsgAck
		proto.Unmarshal(bytes, &msg)
		PrivateMsgAckEvent(ctx, &msg)
	case pb.Packet_PrivateMsg:
		var msg pb.PacketPrivateMsg
		proto.Unmarshal(bytes, &msg)
		PrivateMsgEvent(ctx, &msg)
	case pb.Packet_EventFriendRequest:
		var msg pb.FriendRequestEvent
		proto.Unmarshal(bytes, &msg)
		FriendRequestEvent(ctx, &msg)
	case pb.Packet_EventFriendChange:
		var msg pb.FriendChangeEvent
		proto.Unmarshal(bytes, &msg)
		FriendRequestChange(ctx, &msg)
	default:
		glog.Warningf(ctx, "暂未支持的命令码: %s", bytes)
	}

	return nil
}

func ListenCommander(ctx context.Context, command string, args []string) (isTrigger bool, err error) {
	values := []reflect.Value{
		reflect.ValueOf(clientController),
		reflect.ValueOf(args),
	}
	ref := reflect.TypeOf(clientController)
	for i := 0; i < ref.NumMethod(); i++ {
		if strings.ToLower(ref.Method(i).Name) == strings.ToLower(command) {
			glog.Infof(ctx, "调用控制器: 【%s】调用参数: %v", ref.Method(i).Name, args)
			retValues := ref.Method(i).Func.Call(values)
			if er, has := retValues[0].Interface().(error); has {
				return true, er
			}
			return true, nil
		}
	}
	return false, gerror.Newf("没有找到这个功能: %s", command)
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
