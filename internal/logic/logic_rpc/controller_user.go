package logic_rpc

import (
	"XtTalkServer/app/model/mysql_model"
	"XtTalkServer/global"
	"XtTalkServer/internal/logic/logic_model"
	"XtTalkServer/pb"
	"XtTalkServer/utils"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"gorm.io/gorm"
)

var User = new(_UserController)

type _UserController struct {
}

func (_UserController) GetProfile(device logic_model.ConnDevice, req *pb.PacketGetProfileReq) (res *pb.PacketGetProfileRes, fail error) {
	fmt.Println("收到了来自客户端GetProfile请求")
	var user mysql_model.User
	if err := global.Db.Where("id = ?", device.UserId).First(&user).Error; err != nil {
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

func (_UserController) ModifyProfile(device logic_model.ConnDevice, req *pb.PacketModfiyProfileReq) (res *pb.PacketModfiyProfileRes, fail error) {
	fmt.Println("收到了来自客户端ModifyProfile请求")
	userUpdate := map[string]interface{}{}
	if !utils.IsEmpty(req.GetNickName()) {
		userUpdate["nickname"] = req.GetNickName()
	}
	if err := global.Db.Model(&mysql_model.User{}).Where("id = ?", device.UserId).Updates(userUpdate).Error; err != nil {
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

func (_UserController) GetUser(device logic_model.ConnDevice, req *pb.PacketGetUserInfoReq) (res *pb.PacketGetUserInfoResp, fail error) {
	res = &pb.PacketGetUserInfoResp{}
	var user mysql_model.User
	if err := global.Db.
		Select("id,username,nickname").
		Where("username = ?", req.GetUsername()).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.RetCode = pb.RetCode_UserNotFound //没有找到用户
		} else {
			res.RetCode = pb.RetCode_Error //系统错误
		}
		return
	}

	res.RetCode = pb.RetCode_Success
	res.User = &pb.User{
		UserId:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	}
	return
}
