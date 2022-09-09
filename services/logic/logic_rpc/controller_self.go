package logic_rpc

import (
	"XtTalkServer/app/model/mysql_model"
	"XtTalkServer/global"
	"XtTalkServer/pb"
	"XtTalkServer/services/logic/logic_model"
	"XtTalkServer/utils"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
)

var Self = new(_SelfController)

type _SelfController struct {
}

func (_SelfController) GetProfile(device logic_model.ConnDevice, req *pb.PacketGetProfileReq) (res *pb.PacketGetProfileRes, fail error) {
	fmt.Println("收到了来自客户端GetProfile请求")
	var user mysql_model.User
	if err := global.Db.Where("id = ?", device.UserClient.Uid).First(&user).Error; err != nil {
		return nil, err //账号不存在,异常,都退出登录
	}

	res = &pb.PacketGetProfileRes{
		Sex:      1,
		Age:      2,
		NickName: user.Nickname,
		Email:    user.Email,
		UserId:   user.ID,
		Note:     "默认签名",
	}
	return
}

func (_SelfController) ModifyProfile(device logic_model.ConnDevice, req *pb.PacketModfiyProfileReq) (res *pb.PacketModfiyProfileRes, fail error) {
	fmt.Println("收到了来自客户端ModifyProfile请求")
	userUpdate := map[string]interface{}{}
	if !utils.IsEmpty(req.GetNickName()) {
		userUpdate["nickname"] = req.GetNickName()
	}
	if err := global.Db.Model(&mysql_model.User{}).Where("id = ?", device.UserClient.Uid).Updates(userUpdate).Error; err != nil {
		glog.Warningf(device.Context, "更新用户信息失败: %s", err.Error())
		res = &pb.PacketModfiyProfileRes{
			RetCode: pb.RetCode_Error,
		}
	} else {
		res = &pb.PacketModfiyProfileRes{
			RetCode: pb.RetCode_Success,
		}
	}

	return
}
