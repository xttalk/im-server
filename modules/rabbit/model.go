package rabbit

import (
	"github.com/streadway/amqp"
)

type ExchangeKind string

const (
	ExchangeDirect ExchangeKind = "direct" //直连交换机
	ExchangeFanout ExchangeKind = "fanout" //扇形交换机
	ExchangeTopic  ExchangeKind = "topic"  //主题交换机
)

type RabbitInfo struct {
	ExchangeKind   ExchangeKind //交换机类型
	ExchangeName   string       //交换机名称
	ExchangeParams amqp.Table   //交换机参数

	QueueName   string     //队列名称
	QueueParams amqp.Table //队列参数
	RouteKey    string     //路由key
}
