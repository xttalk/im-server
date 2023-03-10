package types

import (
	"XtTalkServer/modules/rabbit"
	"context"
	"github.com/streadway/amqp"
)

type IConsumer interface {
	//Init 用于声明或其他初始化
	Init(context.Context, *rabbit.Channel) error
	//ConsumerOptions 设置消费者参数
	ConsumerOptions(context.Context) rabbit.ConsumerOptions
	//Callback 用于处理数据
	Callback(context.Context, *rabbit.Channel, amqp.Delivery) error
}
