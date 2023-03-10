package service

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/global"
	"XtTalkServer/pb"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	"time"
)

var PublisherService = new(publisherService)

type publisherService struct {
}

// FriendChangeEvent 投递好友变动事件
func (publisherService) FriendChangeEvent(uid, friendId uint64, sessionIds ...string) error {
	sessionId := ""
	if len(sessionIds) > 0 {
		sessionId = sessionIds[0]
	}
	friend, retCode := FriendService.GetFriend(uid, friendId)
	bytes, err := proto.Marshal(&pb.FriendChangeEvent{
		FriendId: friendId,
		IsFriend: retCode == pb.RetCode_Success, //好友关系状态
		Friend:   friend,                        //如果不是好友这里会是nil
		Time:     gtime.Now().Unix(),
	})
	if err != nil {
		return err
	}
	mqBytes, err := proto.Marshal(&pb.MqMsg{
		Data:      bytes,
		ReceiveId: uid,
		SessionId: sessionId,
	})
	if err != nil {
		return err
	}
	channel, err := global.RabbitMQ.GetChannel()
	if err != nil {
		return err
	}
	defer channel.Release()
	routeKey := conts.GetRouteKey(conts.MqKeyFriendChange, uid)
	if err := channel.CreatePublisher(conts.GetExchangeName(conts.MqFriendChange), routeKey, amqp.Publishing{
		Timestamp:    time.Now(),
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         mqBytes,
		MessageId:    uuid.New().String(),
	}); err != nil {
		return err
	}
	return nil
}

// FriendRequestEvent 投递好友请求/处理事件
func (publisherService) FriendRequestEvent(receiveId uint64, event *pb.FriendRequestEvent) error {
	bytes, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	mqBytes, err := proto.Marshal(&pb.MqMsg{
		Data:      bytes,
		ReceiveId: receiveId, //消息发送给对方
	})
	channel, err := global.RabbitMQ.GetChannel()
	if err != nil {
		return err
	}
	defer channel.Release()
	routeKey := conts.GetRouteKey(conts.MqKeyFriendRequest, event.GetId())
	if err := channel.CreatePublisher(conts.GetExchangeName(conts.MqFriendRequest), routeKey, amqp.Publishing{
		Timestamp:    time.Now(),
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         mqBytes,
		MessageId:    uuid.New().String(),
	}); err != nil {
		return err
	}
	return nil
}

// PrivateMsgSaveEvent 投递私聊消息保存
func (publisherService) PrivateMsgSaveEvent(event *pb.PacketPrivateMsg, sessionId string) error {
	//数据处理
	bytes, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	mqMsg := pb.MqMsg{
		Data:      bytes,
		SessionId: sessionId,
	}
	mqBytes, err := proto.Marshal(&mqMsg)
	if err != nil {
		return err
	}

	//消息投递到消息中心
	channel, err := global.RabbitMQ.GetChannel()
	if err != nil {
		return err
	}
	defer channel.Release()
	routeKey := conts.GetRouteKey(conts.MqKeyPrivateMsgSave, event.GetMsgId())
	if err := channel.CreatePublisher(conts.GetExchangeName(conts.MqPrivateMsgSave), routeKey, amqp.Publishing{
		Timestamp:    time.Now(),
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         mqBytes,
	}); err != nil {
		return err
	}
	return nil
}
