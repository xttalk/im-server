package types

import (
	"XtTalkServer/modules/rabbit"
	"context"
)

type IConsumer interface {
	//Init 用于声明或其他初始化
	Init(context.Context, *rabbit.Channel) error
	//Start 用于启动消费者
	Start(context.Context, *rabbit.Channel) error
}
