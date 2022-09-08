package core

import (
	"XtTalkServer/global"
	"XtTalkServer/modules/rabbit"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"time"
)

var Rabbitmq = new(_rabbitmq)

type _rabbitmq struct {
}

// Initialize 初始化RabbitMQ
func (_rabbitmq) Initialize(ctx context.Context) {
	conf := global.Config.RabbitMq
	//连接rabbitmq
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d", conf.User, conf.Password, conf.Host, conf.Port)

	pool, err := rabbit.NewPool(&rabbit.Config{
		Host:              dsn,
		MinConn:           1,  // 最少连接数量
		MaxConn:           10, // 最大连接数量
		MaxChannelPerConn: 10, // 每个连接的最大信道数量
		MaxLifetime:       time.Duration(3600),
	}) // 建立连接池
	if err != nil {
		fmt.Println(err)
		return
	}
	global.RabbitMQ = pool
	glog.Infof(ctx, "RabbitMQ连接成功")
}
