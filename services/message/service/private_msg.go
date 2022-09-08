package service

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/modules/rabbit"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/streadway/amqp"
	"time"
)

// PrivateMsgConsumer 私聊消息消费者
type PrivateMsgConsumer struct {
}

func (c *PrivateMsgConsumer) Init(ctx context.Context, channel *rabbit.Channel) error {
	return channel.InitDeclare(rabbit.RabbitInfo{
		ExchangeKind: rabbit.ExchangeTopic,
		ExchangeName: conts.MQ_Exchange_PrivateMsg,
		QueueName:    conts.MQ_Queue_PrivateMsg,
		RouteKey:     conts.MQ_Key_PrivateMsg,
	})
}

func (c *PrivateMsgConsumer) Start(ctx context.Context, channel *rabbit.Channel) error {
	for {
		if err := channel.CreateConsumer(conts.MQ_Queue_PrivateMsg, func(delivery amqp.Delivery) {
			delivery.Ack(false) //消息确认
			fmt.Println(fmt.Sprintf("收到消息: %s", delivery.Body))
		}); err != nil {
			return err //启动失败
		}
		glog.Warningf(ctx, "消费者断开,3秒后重连")
		time.Sleep(time.Second * 3)
	}
}
