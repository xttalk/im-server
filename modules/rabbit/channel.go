package rabbit

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/streadway/amqp"
)

type Channel struct {
	*amqp.Channel
	conn         *Connection
	chanIdentity int64 // 该连接的第几个channel
}

func (ch *Channel) Release() error {
	return ch.conn.ReleaseChannel(ch)
}

func (ch *Channel) Close() error {
	return ch.conn.CloseChannel(ch)
}

func (c *Channel) CreateConsumer(ctx context.Context, syncLimit int, queueName string, callback func(context.Context, amqp.Delivery) error) error {
	if 0 >= syncLimit {
		syncLimit = 1
	}
	deliverys, err := c.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	glog.Infof(ctx, "消费者已建立,当前读取队列: %s", queueName)
	ch := make(chan bool, syncLimit) //允许同时执行的数量
	for msg := range deliverys {
		ch <- true
		go func(msg amqp.Delivery) {
			defer func() {
				<-ch
			}()
			if err := callback(ctx, msg); err != nil {
				glog.Warningf(ctx, "消息处理失败: %s", err)
			}
		}(msg)
	}
	return nil
}
func (c *Channel) InitDeclare(data RabbitInfo) error {
	//创建延迟投递的交换机
	if err := c.ExchangeDeclare(data.ExchangeName, string(data.ExchangeKind), true, false, false, false, data.ExchangeParams); err != nil {
		return gerror.Newf("交换机声明失败: %s", err.Error())
	}
	//创建延迟接收队列
	queue, err := c.QueueDeclare(data.QueueName, true, false, false, false, data.QueueParams)
	if err != nil {
		return gerror.Newf("队列声明失败: %s", err.Error())
	}
	//将交换机和队列绑定在一起
	if err := c.QueueBind(queue.Name, data.RouteKey, data.ExchangeName, false, nil); err != nil {
		return gerror.Newf("绑定失败: %s", err.Error())
	}
	return nil
}
func (c *Channel) CreatePublisher(exchangeName, routeKey string, msg amqp.Publishing) error {
	return c.Publish(exchangeName, routeKey, true, false, msg)
}
