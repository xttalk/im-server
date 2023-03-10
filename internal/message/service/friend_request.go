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

// FriendRequestConsumer 好友通知消费者
type FriendRequestConsumer struct {
	EventName string
	RouteKey  string
}

func (c *FriendRequestConsumer) Init(ctx context.Context, channel *rabbit.Channel) error {
	c.EventName = conts.MqFriendRequest
	c.RouteKey = conts.MqKeyFriendRequest

	return channel.InitDeclare(rabbit.RabbitInfo{
		ExchangeKind: rabbit.ExchangeTopic,
		ExchangeName: conts.GetExchangeName(c.EventName),
		QueueName:    conts.GetQueueName(c.EventName),
		RouteKey:     c.RouteKey,
	})
}
func (c *FriendRequestConsumer) ConsumerOptions(ctx context.Context) rabbit.ConsumerOptions {
	return rabbit.ConsumerOptions{
		QueueName: conts.GetQueueName(c.EventName),
		SyncLimit: 100,
		RetrySize: 3,
	}
}
func (c *FriendRequestConsumer) Callback(ctx context.Context, channel *rabbit.Channel, d amqp.Delivery) error {
	var mqMsg pb.MqMsg
	if err := proto.Unmarshal(d.Body, &mqMsg); err != nil {
		return d.Reject(false)
	}
	var msg pb.FriendRequestEvent
	if err := proto.Unmarshal(mqMsg.GetData(), &msg); err != nil {
		//数据解析异常,无法解析
		return d.Reject(false)
	}

	if err := logic_rpc.UserClient.SendUserBytes(ctx, logic_rpc.PacketSendInfo{
		UserId:         mqMsg.GetReceiveId(),
		ExcludeSession: []string{mqMsg.GetSessionId()}, //排除当前sessionid
	}, pb.Packet_EventFriendRequest, mqMsg.GetData()); err != nil {
		return err
	}
	glog.Infof(ctx, "推送好友验证请求成功: %d", msg.GetId())
	return nil
}
