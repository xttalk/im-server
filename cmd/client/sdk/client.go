package sdk

import (
	"XtTalkServer/internal/connect/types"
	"XtTalkServer/pb"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type XtTalkClient struct {
	ctx    context.Context
	conn   net.Conn
	isConn bool

	seq          atomic.Uint32 //全局Seq消息ID
	privateSeqId atomic.Uint32 //私聊Seq消息ID

	SeqHandler sync.Map
}

func (c *XtTalkClient) NextPrivateSeq() uint32 {
	return c.privateSeqId.Add(1)
}
func CreateClient() *XtTalkClient {
	client := &XtTalkClient{
		ctx:        gctx.New(),
		SeqHandler: sync.Map{},
	}
	//初始化seq
	client.seq.Store(0)
	client.privateSeqId.Store(2000)

	return client
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

type SendWaitPacket struct {
	Head types.FixedHeader
	Body []byte
}

func (c *XtTalkClient) SendAndWaitPacket(commandId pb.Packet, msg proto.Message) (*SendWaitPacket, error) {
	seq, data, err := c.BuildPacket(commandId, msg)
	if err != nil {
		return nil, err
	}

	var wait = make(chan *SendWaitPacket, 1)
	fmt.Println("存入seq", seq)
	c.SeqHandler.Store(seq, func(head types.FixedHeader, value []byte) {
		wait <- &SendWaitPacket{
			Head: head,
			Body: value,
		}
	})
	_, err = c.conn.Write(data)
	if err != nil {
		return nil, err
	}
	//监听
	for {
		select {
		case v := <-wait:
			return v, nil
		case <-time.After(time.Second * 5):
			return nil, gerror.Newf("请求超时")
		}
	}
}
func (c *XtTalkClient) BuildPacket(commandId pb.Packet, msg proto.Message) (uint32, []byte, error) {
	_bytes, err := proto.Marshal(msg)
	if err != nil {
		return 0, nil, err
	}
	sendSeq := c.seq.Add(1)
	head := types.FixedHeader{
		Version:  0x01,
		Command:  uint16(commandId),
		Sequence: sendSeq,
		Length:   uint32(len(_bytes)),
	}
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, head.Version); err != nil {
		return 0, nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, head.Command); err != nil {
		return 0, nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, head.Sequence); err != nil {
		return 0, nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, head.Length); err != nil {
		return 0, nil, err
	}
	resultBytes := append(buf.Bytes(), _bytes...)
	return sendSeq, resultBytes, nil
}
func (c *XtTalkClient) SendPacket(commandId pb.Packet, msg proto.Message) error {
	_, data, err := c.BuildPacket(commandId, msg)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(data)
	return err
}

// ListenReader 读取服务端数据
func (c *XtTalkClient) ListenReader(callback func(context.Context, types.FixedHeader, []byte) error) error {
	defer func() {
		fmt.Println("结束监听")
		c.isConn = false
	}()
	fmt.Println("开始监听服务器消息...")
	for c.isConn {
		var headBytes = make([]byte, types.DataPackHeaderLength)
		if _, err := c.conn.Read(headBytes); err != nil {
			return gerror.Wrapf(err, "读取数据Head失败")
		}
		var imHead types.FixedHeader
		buffer := bytes.NewBuffer(headBytes)
		if err := binary.Read(buffer, binary.LittleEndian, &imHead.Version); err != nil {
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
	return nil
}
