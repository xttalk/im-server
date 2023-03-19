package logic_rpc

import (
	"XtTalkServer/app/methods"
	"XtTalkServer/app/model"
	"XtTalkServer/app/model/mongo_model"
	"XtTalkServer/app/model/mysql_model"
	"XtTalkServer/global"
	"XtTalkServer/internal/logic/logic_model"
	"XtTalkServer/internal/logic/service"
	"XtTalkServer/pb"
	"XtTalkServer/utils/snowflake"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"time"
)

var Friend = new(_FriendController)

type _FriendController struct {
}

// GetFriendList 获取好友列表
func (_FriendController) GetFriendList(device logic_model.ConnDevice, req *pb.PacketGetFriendListReq) (res *pb.PacketGetFriendListRes, fail error) {
	fmt.Println("收到了来自客户端GetFriendList请求")

	nav := model.NavPageReq{
		Size: req.GetSize(),
		Page: req.GetPage(),
	}

	//获取当前账号的好友列表
	var resultList []mysql_model.UserFriend
	var resultTotal int64
	if err := global.Db.Model(&resultList).Where("uid = ?", device.UserId).Count(&resultTotal).Error; err != nil {
		res = &pb.PacketGetFriendListRes{
			RetCode: pb.RetCode_Error,
		}
		return
	}
	if err := global.Db.Preload("Friend", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname,username,note,age,sex")
	}).Scopes(nav.UseNav).Where("uid = ?", device.UserId).Find(&resultList).Error; err != nil {
		res = &pb.PacketGetFriendListRes{
			RetCode: pb.RetCode_Error,
		}
		return
	}
	res = &pb.PacketGetFriendListRes{
		List:    make([]*pb.Friend, 0),
		Nav:     model.BuildNavPb(nav, resultTotal),
		RetCode: pb.RetCode_Success,
	}
	for _, item := range resultList {
		res.List = append(res.List, &pb.Friend{
			FriendId: item.ID,  //好友关系ID
			UserId:   item.Fid, //好友的用户ID
			Nickname: item.Friend.Nickname,
			Username: item.Friend.Username,
			Remark:   item.RemarkName,
			Note:     item.Friend.Note,
			Age:      item.Friend.Age,
			Sex:      item.Friend.Sex,
		})
	}
	return
}

// GetFriend 获取好友信息
func (_FriendController) GetFriend(device logic_model.ConnDevice, req *pb.PacketGetFriendReq) (res *pb.PacketGetFriendRes, fail error) {
	friendInfo, retCode := service.FriendService.GetFriend(device.UserId, req.GetUserId())
	if retCode != pb.RetCode_Success {
		res = &pb.PacketGetFriendRes{
			RetCode: retCode,
		}
		return
	}

	res = &pb.PacketGetFriendRes{
		RetCode: pb.RetCode_Success,
		Friend:  friendInfo,
	}
	return
}

// RemoveFriend 删除好友
func (_FriendController) RemoveFriend(device logic_model.ConnDevice, req *pb.PacketRemoveFriendReq) (res *pb.PacketRemoveFriendResp, fail error) {
	//1.验证是否是好友关系
	res = &pb.PacketRemoveFriendResp{}
	has, err := service.FriendService.IsFriendRelative(device.UserId, req.GetUserId())
	if err != nil {
		glog.Errorf(device.Context, "验证好友关系失败: %s", err.Error())
		res.RetCode = pb.RetCode_Error
		return
	}
	if !has {
		res.RetCode = pb.RetCode_FriendNotFound //不是好友关系
		return
	}

	//移除好友关系操作
	if err := service.FriendService.RemoveFriend(device.UserId, req.GetUserId()); err != nil {
		glog.Errorf(device.Context, "删除好友失败: %s", err.Error())
		res.RetCode = pb.RetCode_Error
		return
	} else {
		res.RetCode = pb.RetCode_Success
	}
	res.RetCode = pb.RetCode_Success
	//推送好友关系已建立,推送好友添加事件

	go func() {
		if err := service.PublisherService.FriendChangeEvent(device.UserId, req.GetUserId(), device.SessionId); err != nil {
			glog.Warningf(device.Context, fmt.Sprintf("[消息投递] FriendChangeEvent 向[%v]投递[%v]好友变动失败: %s", device.UserId, req.GetUserId(), err.Error()))
		} else {
			glog.Infof(device.Context, fmt.Sprintf("[消息投递] FriendChangeEvent 向[%v]投递[%v]好友变动成功", device.UserId, req.GetUserId()))
		}
		if err := service.PublisherService.FriendChangeEvent(req.GetUserId(), device.UserId); err != nil {
			glog.Warningf(device.Context, fmt.Sprintf("[消息投递] FriendChangeEvent 向[%v]投递[%v]好友变动失败: %s", req.GetUserId(), device.UserId, err.Error()))
		} else {
			glog.Infof(device.Context, fmt.Sprintf("[消息投递] FriendChangeEvent 向[%v]投递[%v]好友变动成功", req.GetUserId(), device.UserId))
		}
	}()

	return
}

// FriendApply 发起好友申请
func (_FriendController) FriendApply(device logic_model.ConnDevice, req *pb.PacketFriendApplyReq) (res *pb.PacketFriendApplyResp, fail error) {
	res = &pb.PacketFriendApplyResp{}
	//1.查询是否已经是好友关系
	isFriend, err := service.FriendService.IsFriendRelative(device.UserId, req.GetUserId())
	if err != nil {
		res.RetCode = pb.RetCode_Error
		glog.Errorf(device.Context, "查询好友关系失败: %s", err.Error())
		return
	}
	if isFriend {
		res.RetCode = pb.RetCode_FriendAlready //已经是好友关系
		return
	}

	//2.添加好友申请记录
	/**
	情况1:没有创建过记录,直接添加记录发送好友申请消息
	情况2:对方已经申请自己为好友验证,自己又再次发起好友申请,这时候直接成为好友即可
	*/
	requestId, isCompleteFriend, err := service.FriendService.CreateFriendApplyMsg(device.UserId, req.GetUserId(), req.GetReason())
	if err != nil {
		res.RetCode = pb.RetCode_Error //创建失败
		return
	}

	res.RetCode = pb.RetCode_Success
	res.Id = int32(requestId)
	res.IsCompleteFriend = isCompleteFriend
	//推送好友请求通知到消息中心

	go func() {
		if isCompleteFriend {
			//向双方推送好友关系建立
			if err := service.PublisherService.FriendChangeEvent(device.UserId, req.GetUserId(), device.SessionId); err != nil {
				glog.Warningf(device.Context, fmt.Sprintf("[消息投递] FriendChangeEvent 向[%v]投递[%v]好友变动失败: %s", device.UserId, req.GetUserId(), err.Error()))
			} else {
				glog.Infof(device.Context, fmt.Sprintf("[消息投递] FriendChangeEvent 向[%v]投递[%v]好友变动成功", device.UserId, req.GetUserId()))
			}
			if err := service.PublisherService.FriendChangeEvent(req.GetUserId(), device.UserId); err != nil {
				glog.Warningf(device.Context, fmt.Sprintf("[消息投递] FriendChangeEvent 向[%v]投递[%v]好友变动失败: %s", req.GetUserId(), device.UserId, err.Error()))
			} else {
				glog.Infof(device.Context, fmt.Sprintf("[消息投递] FriendChangeEvent 向[%v]投递[%v]好友变动成功", req.GetUserId(), device.UserId))
			}
		} else {
			//推送好友申请事件
			if err := service.PublisherService.FriendRequestEvent(req.GetUserId(), &pb.FriendRequestEvent{
				Id:          int32(requestId),
				FromUid:     device.UserId,
				ToUid:       req.GetUserId(),
				RequestTime: gtime.Now().Unix(),
				Reason:      req.GetReason(),
				Status:      int32(mysql_model.UserFriendRequestStatusWait),
			}); err != nil {
				glog.Warningf(device.Context, fmt.Sprintf("[消息投递] FriendRequestEvent 向[%v]投递好友验证请求失败: %s", req.GetUserId(), err.Error()))
			} else {
				glog.Infof(device.Context, fmt.Sprintf("[消息投递] FriendRequestEvent 向[%v]投递好友验证请求成功", req.GetUserId()))
			}
		}

	}()

	return
}

// FriendHandle 处理好友申请
func (_FriendController) FriendHandle(device logic_model.ConnDevice, req *pb.PacketFriendHandleReq) (res *pb.PacketFriendHandleResp, fail error) {
	res = &pb.PacketFriendHandleResp{}
	//1.查询该记录touid是否是自己才可以处理消息,并且是未处理状态
	var request mysql_model.UserFriendRequest
	if err := global.Db.Where("id = ? AND to_uid = ? AND status = ?", req.GetId(), device.UserId, mysql_model.UserFriendRequestStatusWait).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.RetCode = pb.RetCode_NotRecord //没有找到数据
		} else {
			res.RetCode = pb.RetCode_Error
			glog.Warningf(device.Context, "查询好友验证消息失败: %s", err.Error())
		}
		return
	}

	if req.GetFlag() {
		//成为好友
		if err := service.FriendService.CompleteFriendRelative(request.Id, request.FromUid, request.ToUid); err != nil {
			glog.Warningf(device.Context, "[好友验证]成为好友失败: %s", err.Error())
			res.RetCode = pb.RetCode_Error
			return
		}
		//推送成为好友信息
		go func() {
			if err := service.PublisherService.FriendChangeEvent(request.ToUid, request.FromUid, device.SessionId); err != nil {
				glog.Warningf(device.Context, fmt.Sprintf("[消息投递] FriendChangeEvent 向[%v]投递[%v]好友变动失败: %s", request.ToUid, request.FromUid, err.Error()))
			} else {
				glog.Infof(device.Context, fmt.Sprintf("[消息投递] FriendChangeEvent 向[%v]投递[%v]好友变动成功", request.ToUid, request.FromUid))
			}
			if err := service.PublisherService.FriendChangeEvent(request.FromUid, request.ToUid); err != nil {
				glog.Warningf(device.Context, fmt.Sprintf("[消息投递] FriendChangeEvent 向[%v]投递[%v]好友变动失败: %s", request.FromUid, request.ToUid, err.Error()))
			} else {
				glog.Infof(device.Context, fmt.Sprintf("[消息投递] FriendChangeEvent 向[%v]投递[%v]好友变动成功", request.FromUid, request.ToUid))
			}
		}()
	} else {
		//拒绝好友,更新记录
		if err := global.Db.Where("id = ? AND to_uid = ? AND status = ?", req.GetId(), device.UserId, mysql_model.UserFriendRequestStatusWait).
			Model(&mysql_model.UserFriendRequest{}).
			UpdateColumns(map[string]interface{}{
				"reject_reason": req.GetRejectReason(),
				"response_time": gtime.Now().Unix(),
			}).Error; err != nil {
			res.RetCode = pb.RetCode_Error
			glog.Warningf(device.Context, "[好友验证]拒绝好友失败: %s", err.Error())
			return
		}

		//推送好友拒绝
		go func() {
			//推送好友申请事件
			if err := service.PublisherService.FriendRequestEvent(request.FromUid, &pb.FriendRequestEvent{
				Id:           int32(request.Id),
				FromUid:      request.FromUid,
				ToUid:        request.ToUid,
				RequestTime:  request.RequestTime,
				Reason:       request.Reason,
				Status:       int32(mysql_model.UserFriendRequestStatusReject), //拒绝状态
				RejectReason: request.RejectReason,                             //拒绝原因
			}); err != nil {
				glog.Warningf(device.Context, fmt.Sprintf("[消息投递] FriendRequestEvent 向[%v]投递好友拒绝失败: %s", request.FromUid, err.Error()))
			} else {
				glog.Infof(device.Context, fmt.Sprintf("[消息投递] FriendRequestEvent 向[%v]投递好友拒绝成功", request.FromUid))
			}
		}()
	}
	res.RetCode = pb.RetCode_Success
	return
}

// SendMsg 发送好友消息
func (_FriendController) SendMsg(device logic_model.ConnDevice, req *pb.PacketPrivateMsg) (fail error) {
	//判断双方好友关系
	msg := req

	//补全信息
	msg.ServerTime = time.Now().Unix()
	msg.FromId = device.UserId
	msg.MsgId = snowflake.GetNextIdByServer(int64(device.ServerId)) //通过服务器ID生成唯一ID
	// todo 校验好友关系

	//构建消息
	seq, seqErr := service.FriendService.IncFriendMsgSeq(device.UserId, msg.GetReceiveId())
	if seqErr != nil {
		return seqErr
	}
	msg.Seq = seq

	if err := service.PublisherService.PrivateMsgSaveEvent(msg, device.SessionId); err != nil {
		glog.Warningf(device.Context, "投递私聊消息失败: %s", err.Error())
		return
	}
	//1.向客户端发送ACK消息确认
	ackMsg := pb.PacketPrivateMsgAck{
		MsgSeq:  msg.GetMsgSeq(),
		MsgRand: msg.GetMsgRand(),
		RetCode: pb.PacketMsgStatus_MsgAck,
		MsgId:   msg.MsgId,
		Seq:     msg.Seq,
	}
	//向当前发送端推送ack确认消息
	if err := UserClient.SendUserPacket(device.Context, PacketSendInfo{
		UserId:      device.UserId,
		SendSession: []string{device.SessionId},
	}, pb.Packet_PrivateMsgAck, &ackMsg); err != nil {
		glog.Warningf(device.Context, "私聊ack消息推送失败: %s", err.Error())
		return
	}
	return
}

// GetMessageList 读取私聊消息
func (_FriendController) GetMessageList(device logic_model.ConnDevice, req *pb.PacketPrivateMsgListReq) (res *pb.PacketPrivateMsgListResp, fail error) {
	res = &pb.PacketPrivateMsgListResp{
		IsCompleted: false,
		RetCode:     pb.RetCode_Success,
		List:        make([]*pb.PacketPrivateMsg, 0),
	}

	fromUid := device.UserId
	receiveUid := req.GetUserId() //与目标用户
	var limitSize = req.GetSize()
	if 0 >= limitSize || limitSize > 100 { //最小1 最大100条
		limitSize = 50 //默认50条
	}
	where := bson.D{} //查询条件
	sort := options.Find().SetLimit(limitSize).SetSort(bson.M{
		"seq": -1, //按照时序倒序排序
	})

	table := fmt.Sprintf(mongo_model.TablePrivateMsg, methods.GetPrivateRelativeKey(fromUid, receiveUid))
	//需要查询指定消息ID前的消息列表
	if req.GetLastMsgId() > 0 {
		//消息分割拉取
		//1.查询这条消息的位置对应ObjectID
		var item mongo_model.PrivateMsg
		lastMsgResult := global.Mongo.Collection(table).FindOne(device.Context, bson.D{
			{"msgid", req.GetLastMsgId()},
		})
		if err := lastMsgResult.Decode(&item); err != nil {
			//if lastMsgResult.Err() == mongo.ErrNoDocuments {
			//	fail = gerror.Newf("没有找到这条消息")
			//} else {
			//	fail = gerror.Wrapf(err, "查询失败")
			//}
			fmt.Println("没有找到消息", req.GetLastMsgId())
			res.RetCode = pb.RetCode_Error
			return
		}
		where = append(where, bson.E{"seq", bson.M{"$lt": item.Seq}})
	}
	result, err := global.Mongo.Collection(table).Find(device.Context, where, sort)
	if err != nil {
		//fail = gerror.Wrapf(err, "查询消息记录失败")
		res.RetCode = pb.RetCode_Error
		return
	}
	defer result.Close(device.Context)

	for result.Next(device.Context) {
		var item mongo_model.PrivateMsg
		if err := result.Decode(&item); err != nil {
			fmt.Println("解析消息失败:", err.Error())
			continue
		}
		res.List = append([]*pb.PacketPrivateMsg{
			{
				Seq:        item.Seq,
				MsgSeq:     item.MsgSeq,
				MsgRand:    item.MsgRand,
				FromId:     item.FromId,
				ReceiveId:  item.ReceiveId,
				ClientTime: item.ClientTime,
				ServerTime: item.ServerTime,
				Payload:    item.Payload,
				Extends:    item.Extends,
				MsgType:    item.MsgType,
				MsgId:      item.MsgId,
			},
		}, res.List...)
	}

	return
}
