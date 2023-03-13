package logic_rpc

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/global"
	"XtTalkServer/internal/logic/logic_model"
	"XtTalkServer/pb"
	"XtTalkServer/utils"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
)

var Login = new(_LoginController)

type _LoginController struct {
}

func (_LoginController) Login(device logic_model.ConnDevice, req *pb.PacketLoginReq) (res *pb.PacketLoginRes, fail error) {
	fmt.Println("收到了来自客户端Login请求")
	var uid uint64 = gvar.New(req.Token).Uint64()
	userClient := logic_model.UserClient{
		Uid:       uid,
		SessionId: device.SessionId,
		ServerId:  device.ServerId,
	}
	//登录成功,设置设备标识各项信息
	if err := global.Redis.Set(device.Context, fmt.Sprintf(conts.RK_ClientAuth, device.SessionId), utils.DataToJson(userClient), 0).Err(); err != nil {
		return nil, err //断开
	}
	//设置用户登录端,记录sessionID->ServerID
	if err := global.Redis.HSetNX(device.Context, fmt.Sprintf(conts.RK_UserDevice, uid), device.SessionId, device.ServerId).Err(); err != nil {
		return nil, err
	}
	fmt.Println("用户登录成功", uid)
	res = &pb.PacketLoginRes{
		RetCode: pb.RetCode_Success,
		Uid:     uid,
	}
	return
}
