package sdk

import (
	"XtTalkServer/pb"
	"XtTalkServer/services/connect/types"
	"bytes"
	"encoding/binary"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"net"
	"sync/atomic"
)

type XtTalkClient struct {
	ctx    context.Context
	conn   net.Conn
	seq    uint32
	isConn bool
}

func CreateClient() *XtTalkClient {
	return &XtTalkClient{
		ctx: gctx.New(),
	}
}

// Connect 连接服务器
func (c *XtTalkClient) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return gerror.Wrapf(err, "连接服务器失败")
	}
	c.isConn = true
	c.conn = conn
	return nil
}
func (c *XtTalkClient) SendPacket(commandId pb.Packet, pb proto.Message) error {
	_bytes, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	sendSeq := atomic.AddUint32(&c.seq, 1)
	head := types.ImHeadDataPack{
		ProtocolVersion: types.ProtocolVersion,
		Command:         uint16(commandId),
		Sequence:        sendSeq,
		Length:          uint32(len(_bytes)),
	}
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, head.ProtocolVersion); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, head.Command); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, head.Sequence); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, head.Length); err != nil {
		return err
	}
	resultBytes := append(buf.Bytes(), _bytes...)
	_, err = c.conn.Write(resultBytes)
	return err
}

// ListenReader 读取服务端数据
func (c *XtTalkClient) ListenReader(callback func(context.Context, types.ImHeadDataPack, []byte) error) error {
	defer func() {
		c.isConn = false
	}()
	for {
		var headBytes = make([]byte, types.DataPackHeaderLength)
		if _, err := c.conn.Read(headBytes); err != nil {
			return gerror.Wrapf(err, "读取数据Head失败")
		}
		var imHead types.ImHeadDataPack
		buffer := bytes.NewBuffer(headBytes)
		if err := binary.Read(buffer, binary.LittleEndian, &imHead.ProtocolVersion); err != nil {
			return gerror.Wrapf(err, "解包失败ProtocolVersion")
		}
		if err := binary.Read(buffer, binary.LittleEndian, &imHead.Command); err != nil {
			return gerror.Wrapf(err, "解包失败Command")
		}
		if err := binary.Read(buffer, binary.LittleEndian, &imHead.Sequence); err != nil {
			return gerror.Wrapf(err, "解包失败Sequence")
		}
		if err := binary.Read(buffer, binary.LittleEndian, &imHead.Length); err != nil {
			return gerror.Wrapf(err, "解包失败Length")
		}
		var dataBytes = make([]byte, imHead.Length)
		if _, err := c.conn.Read(dataBytes); err != nil {
			return gerror.Wrapf(err, "读取内容数据失败")
		}
		ctx := gctx.New()
		if err := callback(ctx, imHead, dataBytes); err != nil {
			return gerror.Wrapf(err, "客户端处理失败")
		}
	}
}
