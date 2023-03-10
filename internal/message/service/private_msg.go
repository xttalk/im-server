package service

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/internal/logic/logic_rpc"
	"XtTalkServer/modules/rabbit"
	"XtTalkServer/pb"
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
)

// PrivateMsgConsumer 私聊消息消费者
type PrivateMsgConsumer struct {
	EventName string
	RouteKey  string
}

func (c *PrivateMsgConsumer) Init(ctx context.Context, channel *rabbit.Channel) error {
	c.EventName = conts.MqPrivateMsg
	c.RouteKey = conts.MqKeyPrivateMsg

	return channel.InitDeclare(rabbit.RabbitInfo{
		ExchangeKind: rabbit.ExchangeTopic,
		ExchangeName: conts.GetExchangeName(c.EventName),
		QueueName:    conts.GetQueueName(c.EventName),
		RouteKey:     c.RouteKey,
	})
}
func (c *PrivateMsgConsumer) ConsumerOptions(ctx context.Context) rabbit.ConsumerOptions {
	return rabbit.ConsumerOptions{
		QueueName: conts.GetQueueName(c.EventName),
		SyncLimit: 100,
		RetrySize: 3,
	}
}

func (c *PrivateMsgConsumer) Callback(ctx context.Context, channel *rabbit.Channel, d amqp.Delivery) error {
	var mqMsg pb.MqMsg
	if err := proto.Unmarshal(d.Body, &mqMsg); err != nil {
		return d.Reject(false)
	}
	var msg pb.PacketPrivateMsg
	if err := proto.Unmarshal(mqMsg.GetData(), &msg); err != nil {
		//数据解析异常,无法解析
		return d.Reject(false)
	}

	//向当前用户的其他端推送消息
	if err := logic_rpc.UserClient.SendUserBytes(ctx, logic_rpc.PacketSendInfo{
		UserId:         msg.GetFromId(),
		ExcludeSession: []string{mqMsg.GetSessionId()}, //排除当前发送端
	}, pb.Packet_PrivateMsg, mqMsg.GetData()); err != nil {
		glog.Warningf(ctx, "推送接收者消息失败: %s", err.Error())
		return err
	}

	//向消息接收者推送消息
	if err := logic_rpc.UserClient.SendUserBytes(ctx, logic_rpc.PacketSendInfo{
		UserId: msg.GetReceiveId(),
	}, pb.Packet_PrivateMsg, mqMsg.GetData()); err != nil {
		glog.Warningf(ctx, "推送接收者消息失败: %s", err.Error())
		return err
	}

	return nil
}
