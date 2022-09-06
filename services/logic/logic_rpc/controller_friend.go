package logic_rpc

import (
	"XtTalkServer/app/model"
	"XtTalkServer/global"
	"XtTalkServer/pb"
	"XtTalkServer/services/logic/logic_model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var Friend = new(_FriendController)

type _FriendController struct {
}

func (_FriendController) GetFriendList(device logic_model.ConnDevice, req *pb.PacketGetFriendListReq) (res *pb.PacketGetFriendListRes, fail error) {
	fmt.Println("收到了来自客户端GetFriendList请求")

	nav := model.NavPageReq{
		Size: req.GetSize(),
		Page: req.GetPage(),
	}

	//获取当前账号的好友列表
	var resultList []model.UserFriend
	var resultTotal int64
	if err := global.Db.Model(&resultList).Where("uid = ?", device.UserClient.Uid).Count(&resultTotal).Error; err != nil {
		res = &pb.PacketGetFriendListRes{
			RetCode: pb.RetCode_Error,
		}
		return
	}
	if err := global.Db.Preload("Friend", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,username")
	}).Scopes(nav.UseNav).Where("uid = ?", device.UserClient.Uid).Find(&resultList).Error; err != nil {
		res = &pb.PacketGetFriendListRes{
			RetCode: pb.RetCode_Error,
		}
		return
	}
	res = &pb.PacketGetFriendListRes{
		List: make([]*pb.Friend, 0),
		Nav:  model.BuildNavPb(nav, resultTotal),
	}
	for _, item := range resultList {
		res.List = append(res.List, &pb.Friend{
			FriendId: item.ID,  //好友关系ID
			UserId:   item.Fid, //好友的用户ID
			Nickname: item.Friend.Nickname,
			Username: item.Friend.Username,
			Remark:   item.RemarkName,
		})
	}
	return
}

func (_FriendController) GetFriend(device logic_model.ConnDevice, req *pb.PacketGetFriendReq) (res *pb.PacketGetFriendRes, fail error) {
	var result model.UserFriend
	if err := global.Db.Preload("Friend", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,username,nickname")
	}).Where("uid = ? AND fid = ?", device.UserClient.Uid, req.GetUserId()).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res = &pb.PacketGetFriendRes{
				RetCode: pb.RetCode_FriendNotFound,
			}
		} else {
			res = &pb.PacketGetFriendRes{
				RetCode: pb.RetCode_Error,
			}
		}
		return
	}
	res = &pb.PacketGetFriendRes{
		RetCode: pb.RetCode_Success,
		Friend: &pb.Friend{
			FriendId: result.ID,  //好友关系ID
			UserId:   result.Fid, //好友用户ID
			Nickname: result.Friend.Nickname,
			Username: result.Friend.Username,
			Remark:   result.RemarkName,
		},
	}

	return
}

func (_FriendController) SendMsg(device logic_model.ConnDevice, req *pb.PacketPrivateMsg) (fail error) {
	//判断双方好友关系
	fmt.Println("发送私聊消息 -> ", req.ReceiveId)
	UserClient.SendUserPacket(device.Context, req.ReceiveId, pb.Packet_PrivateMsg, req)
	return
}
