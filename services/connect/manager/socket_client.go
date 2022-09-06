package manager

import (
	"XtTalkServer/pb"
	"XtTalkServer/services/connect/types"
	"XtTalkServer/utils"
	"XtTalkServer/utils/snowflake"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/panjf2000/gnet/v2"
	"google.golang.org/protobuf/proto"
	"sync/atomic"
	"time"
)

type ClientMode int

const (
	WsClientMode  ClientMode = 1
	TcpClientMode ClientMode = 2
)

// CreateWsClient 创建一个客户端
func CreateWsClient(conn gnet.Conn) *Client {
	client := &Client{
		Context:    gctx.New(),
		ClientMode: WsClientMode,
		conn:       conn,
	}
	client.init()
	return client
}

// CreateTcpClient 创建一个客户端
func CreateTcpClient(conn gnet.Conn) *Client {
	client := &Client{
		Context:    gctx.New(),
		ClientMode: TcpClientMode,
		conn:       conn,
	}
	client.init()
	return client
}

type ClientUserData struct {
	//用户额外信息
	Uid       uint64                 //用户ID
	LoginTime int64                  //登录成功时间
	Params    map[string]interface{} //更多参数存放
}
type Client struct {
	Context        context.Context
	ClientMode     ClientMode       //客户端连接模式
	ConnectTime    int64            //连接时间
	UserData       *ClientUserData  //没有登录则为nil
	SessionId      string           //设备当前sessionID
	HeartTime      int64            //上次心跳时间
	isClose        bool             //是否已经关闭
	ServerSeq      uint32           //服务端Seq
	LastDataFormat types.DataFormat //最后一次发送消息时数据类型

	conn gnet.Conn
}

func (c *Client) init() {
	//生成客户端唯一标识
	rk := fmt.Sprintf("%d_%d_%s", snowflake.GetNextID(), c.conn.Fd(), utils.RandomStr(16))
	v, err := gmd5.Encrypt(rk)
	if err != nil {
		c.Warningf("客户端生成唯一链接表示失败: %s", err.Error())
		c.Close()
		return
	}
	c.SessionId = v
}

func (c *Client) SendClientPbPack(seq uint32, command pb.Packet, pb proto.Message) error {
	_bytes, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	resultBytes, err := c.buildPacket(seq, command, types.DataFormatPb, _bytes)
	if err != nil {
		return nil
	}
	return c.SendBytes(resultBytes)
}
func (c *Client) SendClientJsonPack(seq uint32, command pb.Packet, v interface{}) error {
	_bytes, err := gjson.Marshal(v)
	if err != nil {
		return err
	}
	resultBytes, err := c.buildPacket(seq, command, types.DataFormatJson, _bytes)
	if err != nil {
		return nil
	}
	return c.SendBytes(resultBytes)
}
func (c *Client) SendServerBytes(command pb.Packet, _bytes []byte) error {
	seq := atomic.AddUint32(&c.ServerSeq, 1)
	return c.SendClientPacket(seq, command, c.LastDataFormat, _bytes)
}
func (c *Client) SendClientPacket(seq uint32, command pb.Packet, dataType types.DataFormat, _bytes []byte) error {
	resultBytes, err := c.buildPacket(seq, command, dataType, _bytes)
	if err != nil {
		return nil
	}
	return c.SendBytes(resultBytes)
}
func (c *Client) buildPacket(seq uint32, command pb.Packet, dataType types.DataFormat, _bytes []byte) ([]byte, error) {
	head := &types.ImHeadDataPack{
		ProtocolVersion: types.ProtocolVersion,
		DataFormat:      dataType,
		Command:         uint16(command),
		Sequence:        seq,
		Length:          uint32(len(_bytes)),
	}
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, head); err != nil {
		return nil, err
	}
	resultBytes := append(buf.Bytes(), _bytes...)
	return resultBytes, nil
}

// SendBytes 发送原始数据,不含包头
func (c *Client) SendBytes(bytes []byte) (err error) {
	if c.isClose {
		return nil
	}
	switch c.ClientMode {
	case TcpClientMode:
		_, err = c.conn.Write(bytes)
	case WsClientMode:
		err = wsutil.WriteServerMessage(c.conn, ws.OpText, bytes)
	}
	return
}

// Close 主动断开连接
func (c *Client) Close() (err error) {
	if c.conn != nil {
		err = c.conn.Close()
		if err == nil {
			c.conn = nil
		}
	}
	c.isClose = true
	return
}

// IsHeartTimeout 检测心跳是否超时
func (c *Client) IsHeartTimeout() bool {
	s := 120 //todo 120秒心跳时间
	if time.Now().Unix() >= c.HeartTime+int64(s) {
		return true
	}
	return false
}

// Infof 普通日志输出
func (c *Client) Infof(format string, args ...interface{}) {
	glog.Infof(c.Context, "[%s]%s", c.SessionId, fmt.Sprintf(format, args...))
}

// Debugf 调试日志输出
func (c *Client) Debugf(format string, args ...interface{}) {
	glog.Debugf(c.Context, "[%s]%s", c.SessionId, fmt.Sprintf(format, args...))
}

// Warningf 警告日志输出
func (c *Client) Warningf(format string, args ...interface{}) {
	glog.Warningf(c.Context, "[%s]%s", c.SessionId, fmt.Sprintf(format, args...))
}