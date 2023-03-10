package mongo_model

import (
	"XtTalkServer/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"
)

const (
	TablePrivateMsg = "private_msg_%s"
)

type PrivateMsg struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Seq           int64              //时序ID,由服务端生成
	MsgId         int64              //消息ID,由服务端生成
	MsgSeq        int64              //消息序列
	MsgRand       int64              //消息随机值
	FromId        uint64             //发送方UID
	ReceiveId     uint64             //接收方UID
	MsgType       pb.PacketMsgType   //消息类型
	Payload       []byte             //数据内容
	SimpleContent string             //用于简略查看的文字
	ClientTime    int64              //客户端发送时间
	ServerTime    int64              //服务端接收时间,查询按照这个时间查询
	Extends       []byte             //额外数据
}

func (c *PrivateMsg) GetSimpleContent() string {
	switch c.MsgType {
	case pb.PacketMsgType_Text:
		if c.Payload != nil {
			var msg pb.TextMsg
			if err := proto.Unmarshal(c.Payload, &msg); err == nil {
				return msg.GetContent()
			}
		}
		return ""
	case pb.PacketMsgType_Image:
		return "[图片]"
	case pb.PacketMsgType_Audio:
		return "[语音]"
	case pb.PacketMsgType_Video:
		return "[视频]"
	case pb.PacketMsgType_File:
		return "[文件]"
	}
	return "暂未支持的消息格式"
}
