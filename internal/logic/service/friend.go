package service

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/app/methods"
	"XtTalkServer/app/model/mongo_model"
	"XtTalkServer/app/model/mysql_model"
	"XtTalkServer/global"
	"XtTalkServer/pb"
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"time"
)

var FriendService = new(friendService)

type friendService struct {
}

// IncFriendMsgSeq 自增好友之间的时序ID
func (c *friendService) IncFriendMsgSeq(uid1, uid2 uint64) (int64, error) {
	//按照uid最小的排在前面
	rk := fmt.Sprintf(conts.RK_PrivateMsgSeq, methods.GetPrivateRelativeKey(uid1, uid2))
	ctx := gctx.New()
	//判断是否存在
	if global.Redis.Exists(ctx, rk).Val() != 0 {
		return global.Redis.IncrBy(context.TODO(), rk, 1).Val(), nil
	}
	//没有找到从mongodb中查询最大的聊天记录值,然后在缓存到redis
	table := fmt.Sprintf(mongo_model.TablePrivateMsg, methods.GetPrivateRelativeKey(uid1, uid2))
	var result mongo_model.PrivateMsg
	var seq int64 = -1
	if err := global.Mongo.Collection(table).FindOne(ctx, options.FindOne().SetSort(bson.M{
		"seq": -1,
	})).Decode(&result); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err.Error())
			return -1, nil
		}
		seq = 1
	} else {
		seq = result.Seq + 1
	}
	if err := global.Redis.Set(ctx, rk, seq, time.Duration(-1)).Err(); err != nil {
		glog.Warningf(ctx, "设置私聊时序SeqID失败: [%s] -> %s", rk, err.Error())
	}
	return seq, nil
}

// IsFriendRelative 验证是否是单向好友关系
func (c *friendService) IsFriendRelative(uid1, uid2 uint64) (bool, error) {
	var cnt int64 = 0
	if err := global.Db.Where("uid = ? AND fid = ?", uid1, uid2).Model(&mysql_model.UserFriend{}).Count(&cnt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return cnt > 0, nil
}

// IsFriendTwoWayRelative 验证是否是双向好友关系
func (c *friendService) IsFriendTwoWayRelative(uid1, uid2 uint64) (bool, error) {
	var cnt int64 = 0
	if err := global.Db.Where("(uid = ? AND fid = ?) OR (uid = ? AND fid = ?)", uid1, uid2, uid2, uid1).Model(&mysql_model.UserFriend{}).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt >= 2, nil
}

// CreateFriendApplyMsg 创建好友关系验证消息
func (c *friendService) CreateFriendApplyMsg(fromUid, targetUid uint64, reason string) (int, bool, error) {

	var isCompleteFriend = false //是否成为好友
	var requestId = 0

	//首次判断接收方是否已经对发送方发起过一次待验证的消息
	if err := global.Db.Transaction(func(tx *gorm.DB) error {
		var firstModel mysql_model.UserFriendRequest
		//先反向查询是否对方有申请自己为好友的记录
		if err := tx.Where("from_uid = ? AND to_uid = ? AND status = ?", targetUid, fromUid, mysql_model.UserFriendRequestStatusWait).First(&firstModel).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			//开始常规流程,
			model := mysql_model.UserFriendRequest{}
			//读取第一条
			if err = tx.Where("from_uid = ? AND to_uid = ?", fromUid, targetUid).First(&model).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				//以前没有申请过,新增一个
				model = mysql_model.UserFriendRequest{
					FromUid:     fromUid,
					ToUid:       targetUid,
					RequestTime: gtime.Now().Unix(),
					Reason:      reason,
					Status:      mysql_model.UserFriendRequestStatusWait,
					LastTime:    gtime.Now().Unix(),
				}
				//已经发起过验证,更新时间和消息
				if err = tx.Create(&model).Error; err != nil {
					return err
				}
				requestId = model.Id
				fmt.Println("已创建消息验证: ", model.Id)
			} else {
				//以前申请过,直接修改状态
				if err = tx.Where("from_uid = ? AND to_uid = ?", fromUid, targetUid).
					Model(&model).
					UpdateColumns(map[string]interface{}{
						"request_time": time.Now().Unix(),
						"reason":       reason,
						"status":       mysql_model.UserFriendRequestStatusWait,
						"last_time":    time.Now().Unix(),
					}).Error; err != nil {
					return err
				}
				requestId = model.Id
				fmt.Println("重复消息验证更新: ", model.Id)
			}
			//还不是好友
		} else {
			//1.创建好友关系(内置一个事务操作)
			if err := c.CompleteFriendRelative(firstModel.Id, fromUid, targetUid); err != nil {
				return err
			}
			isCompleteFriend = true
			requestId = firstModel.Id
		}
		return nil
	}); err != nil {
		return 0, false, err
	}
	return requestId, isCompleteFriend, nil
}

// CompleteFriendRelative 完成创建好友关系,并且下发消息
func (c *friendService) CompleteFriendRelative(requestId int, uid1, uid2 uint64) error {
	err := global.Db.Transaction(func(tx *gorm.DB) error {
		//1.修改验证消息状态
		if err := tx.Where("id = ?", requestId).
			Model(&mysql_model.UserFriendRequest{}).
			UpdateColumns(map[string]interface{}{
				"response_time": time.Now().Unix(),
				"status":        mysql_model.UserFriendRequestStatusAccept,
			}).Error; err != nil {
			return err
		}
		//添加双方好友关系
		//用户1创建好友关系
		addTime := time.Now().Unix()
		if err := tx.Create(&mysql_model.UserFriend{
			Uid:     uid1,
			Fid:     uid2,
			AddTime: addTime,
		}).Error; err != nil {
			return err
		}
		//用户2创建好友关系
		if err := tx.Create(&mysql_model.UserFriend{
			Uid:     uid2,
			Fid:     uid1,
			AddTime: addTime,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// RemoveFriend 删除好友关系
func (c *friendService) RemoveFriend(uid1, uid2 uint64) error {
	//1.移除mysql关系链(双向)
	//事务开
	if err := global.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("(uid = ? AND fid = ?) OR (uid = ? AND fid = ?)", uid1, uid2, uid2, uid1).Delete(&mysql_model.UserFriend{}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	//todo 2.发送通知事件
	go func() {
		//1.给对方发送好友删除事件

		//2.给自己其他在线设备发送该好友移除

	}()

	return nil
}

// GetFriend 获取好友信息
func (c *friendService) GetFriend(uid, fid uint64) (*pb.Friend, pb.RetCode) {
	var result mysql_model.UserFriend
	if err := global.Db.Preload("Friend", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,username,nickname")
	}).Where("uid = ? AND fid = ?", uid, fid).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pb.RetCode_FriendNotFound
		} else {
			return nil, pb.RetCode_Error
		}
	}
	return &pb.Friend{
		FriendId: result.ID,  //好友关系ID
		UserId:   result.Fid, //好友用户ID
		Nickname: result.Friend.Nickname,
		Username: result.Friend.Username,
		Remark:   result.RemarkName,
	}, pb.RetCode_Success
}
