package service

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/app/methods"
	"XtTalkServer/app/model/mongo_model"
	"XtTalkServer/global"
	"XtTalkServer/modules/rabbit"
	"XtTalkServer/pb"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"time"
)

// PrivateMsgSaveConsumer 私聊消息保存消费者
type PrivateMsgSaveConsumer struct {
	EventName string
	RouteKey  string
}

func (c *PrivateMsgSaveConsumer) Init(ctx context.Context, channel *rabbit.Channel) error {
	c.EventName = conts.MqPrivateMsgSave
	c.RouteKey = conts.MqKeyPrivateMsgSave

	return channel.InitDeclare(rabbit.RabbitInfo{
		ExchangeKind: rabbit.ExchangeTopic,
		ExchangeName: conts.GetExchangeName(c.EventName),
		QueueName:    conts.GetQueueName(c.EventName),
		RouteKey:     c.RouteKey,
	})
}

func (c *PrivateMsgSaveConsumer) ConsumerOptions(ctx context.Context) rabbit.ConsumerOptions {
	return rabbit.ConsumerOptions{
		QueueName: conts.GetQueueName(c.EventName),
		SyncLimit: 100,
		RetrySize: 3,
	}
}

func (c *PrivateMsgSaveConsumer) Callback(ctx context.Context, channel *rabbit.Channel, d amqp.Delivery) error {
	var mqMsg pb.MqMsg
	if err := proto.Unmarshal(d.Body, &mqMsg); err != nil {
		return d.Reject(false)
	}
	//解析消息为私聊消息
	var msg pb.PacketPrivateMsg
	if err := proto.Unmarshal(mqMsg.GetData(), &msg); err != nil {
		//数据解析异常,无法解析
		return d.Reject(false)
	}
	var saveMsg = mongo_model.PrivateMsg{
		Seq:        msg.GetSeq(),
		MsgId:      msg.GetMsgId(),
		MsgSeq:     msg.GetMsgSeq(),
		MsgRand:    msg.GetMsgRand(),
		FromId:     msg.GetFromId(),
		ReceiveId:  msg.GetReceiveId(),
		Payload:    msg.GetPayload(),
		MsgType:    msg.GetMsgType(),
		ClientTime: msg.GetClientTime(),
		ServerTime: msg.GetServerTime(),
	}
	saveMsg.SimpleContent = saveMsg.GetSimpleContent()
	//1.进行消息保存
	table := fmt.Sprintf(mongo_model.TablePrivateMsg, methods.GetPrivateRelativeKey(msg.GetFromId(), msg.GetReceiveId()))
	resp, err := global.Mongo.Collection(table).InsertOne(ctx, &saveMsg)
	if err != nil {
		glog.Errorf(ctx, "消息保存数据库失败: %s", err.Error())
		return err
	}
	glog.Infof(ctx, "保存私聊消息成功: [%s] %s", table, resp.InsertedID)
	//2.投递消息
	routeKey := conts.GetRouteKey(conts.MqKeyPrivateMsg, msg.MsgId)
	if err := channel.CreatePublisher(conts.GetExchangeName(conts.MqPrivateMsg), routeKey, amqp.Publishing{
		Timestamp:    time.Now(),
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         d.Body,
	}); err != nil { //忽略错误
		glog.Warningf(ctx, "投递私聊消息推送失败: %s", err.Error())
	}

	return nil
}
