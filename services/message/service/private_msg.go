package service

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/app/model/mongo_model"
	"XtTalkServer/global"
	"XtTalkServer/modules/rabbit"
	"XtTalkServer/pb"
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/golang/protobuf/proto"
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
	if err := channel.CreateConsumer(ctx, 100, conts.MQ_Queue_PrivateMsg, func(ctx context.Context, d amqp.Delivery) error {
		//解析消息为私聊消息
		var msg pb.PacketPrivateMsg
		if err := proto.Unmarshal(d.Body, &msg); err != nil {
			//数据解析异常,无法解析
			return d.Reject(false)
		}
		var saveMsg = mongo_model.PrivateMsg{
			MsgId:      msg.GetMsgId(),
			FromId:     msg.GetFromId(),
			ReceiveId:  msg.GetReceiveId(),
			Payload:    msg.GetPayload(),
			MsgType:    msg.GetMsgType(),
			ClientTime: msg.GetClientTime(),
			ServerTime: msg.GetServerTime(),
		}
		saveMsg.SimpleContent = saveMsg.GetSimpleContent()

		if err := func() error {
			//进行消息保存
			resp, err := global.Mongo.Collection(mongo_model.TablePrivateMsg).InsertOne(ctx, &saveMsg)
			if err != nil {
				glog.Errorf(ctx, "消息保存数据库失败: %s", err.Error())
				return err
			}
			glog.Infof(ctx, "保存私聊消息成功: %s", resp.InsertedID)

			return nil
		}(); err != nil {
			//消息重发
			retrySize, _ := global.Redis.HGet(ctx, conts.RK_MqRetry_PrivateMsg, gvar.New(msg.MsgId).String()).Int()
			//判断是否还需要继续重发
			if 3 > retrySize+1 {
				glog.Warningf(ctx, "[%d]消息发送失败,等待下一次投递[%d/3]: %s", msg.MsgId, retrySize+1, err.Error())
				//消息重发
				global.Redis.HIncrBy(ctx, conts.RK_MqRetry_PrivateMsg, gvar.New(msg.MsgId).String(), 1)
				time.Sleep(time.Second) //1秒后延迟出去
				return d.Reject(true)
			}
			glog.Warningf(ctx, "[%d]消息发送失败,终止投递[%d/3]: %s", msg.MsgId, retrySize+1, err.Error())
			defer global.Redis.HDel(ctx, conts.RK_MqRetry_PrivateMsg, gvar.New(msg.MsgId).String())
			return d.Reject(false)
		}
		defer global.Redis.HDel(ctx, conts.RK_MqRetry_PrivateMsg, gvar.New(msg.MsgId).String())
		return d.Ack(false)
	}); err != nil {
		return err //启动失败
	}

	return nil
}
