package service

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/internal/logic/logic_rpc"
	"XtTalkServer/modules/rabbit"
	"XtTalkServer/pb"
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

// FriendChangeConsumer 好友变动消费者
type FriendChangeConsumer struct {
	EventName string
	RouteKey  string
}

func (c *FriendChangeConsumer) Init(ctx context.Context, channel *rabbit.Channel) error {
	c.EventName = conts.MqFriendChange
	c.RouteKey = conts.MqKeyFriendChange

	return channel.InitDeclare(rabbit.RabbitInfo{
		ExchangeKind: rabbit.ExchangeTopic,
		ExchangeName: conts.GetExchangeName(c.EventName),
		QueueName:    conts.GetQueueName(c.EventName),
		RouteKey:     c.RouteKey,
	})
}
func (c *FriendChangeConsumer) ConsumerOptions(ctx context.Context) rabbit.ConsumerOptions {
	return rabbit.ConsumerOptions{
		QueueName: conts.GetQueueName(c.EventName),
		SyncLimit: 100,
		RetrySize: 3,
	}
}
func (c *FriendChangeConsumer) Callback(ctx context.Context, channel *rabbit.Channel, d amqp.Delivery) error {
	var mqMsg pb.MqMsg
	if err := proto.Unmarshal(d.Body, &mqMsg); err != nil {
		return d.Reject(false)
	}
	var msg pb.FriendChangeEvent
	if err := proto.Unmarshal(mqMsg.GetData(), &msg); err != nil {
		//数据解析异常,无法解析
		return d.Reject(false)
	}
	if err := logic_rpc.UserClient.SendUserBytes(ctx, logic_rpc.PacketSendInfo{
		UserId:         mqMsg.GetReceiveId(),
		ExcludeSession: []string{mqMsg.GetSessionId()}, //排除当前sessionid
	}, pb.Packet_EventFriendChange, mqMsg.GetData()); err != nil {
		return err
	}
	glog.Infof(ctx, "推送好友变动成功")
	return nil
}
