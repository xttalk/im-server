package server_protocol

import (
	"XtTalkServer/internal/connect/manager"
	"XtTalkServer/internal/connect/service"
	"XtTalkServer/internal/connect/types"
	"bytes"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gorilla/websocket"
	"github.com/panjf2000/gnet/v2"
	"golang.org/x/net/context"
	"io"
)

type WebSocket struct {
	addr     string
	port     int
	protocol string
	context  context.Context
	upgrader websocket.Upgrader
	gnet.BuiltinEventEngine
	codec *types.XtTalkTcpCodec
}

func (c *WebSocket) Start() error {
	return gnet.Run(c, c.addr,
		gnet.WithMulticore(false),
		gnet.WithReuseAddr(true),
	)
}
func (c *WebSocket) GetPort() int {
	return c.port
}
func (c *WebSocket) GetProtocol() string {
	return c.protocol
}

func (c *WebSocket) OnBoot(eng gnet.Engine) gnet.Action {
	glog.Infof(c.context, "Websocket服务端已建立")
	c.codec = new(types.XtTalkTcpCodec)

	//注册服务发现
	if err := service.ServiceDiscovery.RegiterConnect(c.GetProtocol(), c.GetPort()); err != nil {
		glog.Warningf(c.context, "[%s]注册连接层服务发现失败: %s", c.GetProtocol(), err.Error())
		return gnet.Shutdown
	}
	return gnet.None
}

func (c *WebSocket) OnOpen(conn gnet.Conn) ([]byte, gnet.Action) {
	conn.SetContext(new(wsCodec))
	//这里不能触发OnConnect事件,不然发消息会报错
	return nil, gnet.None
}

func (c *WebSocket) OnClose(conn gnet.Conn, err error) (action gnet.Action) {
	client := manager.ClientManager.GetClientByFd(manager.WsClientMode, conn.Fd())
	if client != nil {
		manager.ClientManager.DelClient(client)
	}
	return gnet.None
}
func (c *WebSocket) OnTraffic(conn gnet.Conn) (action gnet.Action) {
	ws := conn.Context().(*wsCodec)
	if ws.readBufferBytes(conn) == gnet.Close {
		return gnet.Close
	}
	ok, action := ws.upgrade(conn)
	if !ok {
		return gnet.None
	}
	if !ws.isConnect { //初始化
		ws.isConnect = true
		client := manager.CreateWsClient(conn)
		manager.ClientManager.AddClient(client)
		return gnet.None
	}

	if ws.buf.Len() <= 0 {
		return gnet.None
	}
	messages, err := ws.Decode(conn)
	if err != nil {
		return gnet.Close
	}
	if messages == nil {
		return gnet.None
	}

	client := manager.ClientManager.GetClientByFd(manager.WsClientMode, conn.Fd())
	if client == nil {
		glog.Warningf(c.context, "无法获取到WS客户端: %v", conn.Fd())
		return gnet.Close
	}
	for _, message := range messages {
		//拆包
		head, data, err := c.codec.DecodeBytes(message.Payload)
		if err != nil {
			glog.Warningf(c.context, "读取数据包失败: %s", err.Error())
			return gnet.Close
		}
		manager.CoreServer.OnMessage(client, head, data)
	}
	return gnet.None
}

// CreateWsServer 创建Ws服务器
func CreateWsServer(port int) types.IServer {
	return &WebSocket{
		port:     port,
		protocol: "ws",
		addr:     fmt.Sprintf("tcp://:%d", port),
		context:  gctx.New(),
		upgrader: websocket.Upgrader{},
	}
}

type wsCodec struct {
	upgraded  bool         // 链接是否升级
	buf       bytes.Buffer // 从实际socket中读取到的数据缓存
	wsMsgBuf  wsMessageBuf // ws 消息缓存
	isConnect bool         //是否连接
}
type wsMessageBuf struct {
	firstHeader *ws.Header
	curHeader   *ws.Header
	cachedBuf   bytes.Buffer
}
type readWrite struct {
	io.Reader
	io.Writer
}

func (w *wsCodec) readBufferBytes(c gnet.Conn) gnet.Action {
	size := c.InboundBuffered()
	buf := make([]byte, size, size)
	read, err := c.Read(buf)
	if err != nil {
		//logging.Infof("read err! %w", err)
		return gnet.Close
	}
	if read < size {
		return gnet.Close
	}
	w.buf.Write(buf)
	return gnet.None
}
func (w *wsCodec) upgrade(c gnet.Conn) (ok bool, action gnet.Action) {
	if w.upgraded {
		ok = true
		return
	}
	buf := &w.buf
	tmpReader := bytes.NewReader(buf.Bytes())
	oldLen := tmpReader.Len()
	hs, err := ws.Upgrade(readWrite{tmpReader, c})
	_ = hs
	skipN := oldLen - tmpReader.Len()
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF { //数据不完整
			return
		}
		buf.Next(skipN)
		//logging.Infof("conn[%v] [err=%v]", c.RemoteAddr().String(), err.Error())
		action = gnet.Close
		return
	}
	buf.Next(skipN)
	//logging.Infof("conn[%v] upgrade websocket protocol! Handshake: %v", c.RemoteAddr().String(), hs)
	if err != nil {
		//logging.Infof("conn[%v] [err=%v]", c.RemoteAddr().String(), err.Error())
		action = gnet.Close
		return
	}
	ok = true
	w.upgraded = true
	return
}
func (w *wsCodec) Decode(c gnet.Conn) (outs []wsutil.Message, err error) {
	//fmt.Println("do Decode")
	messages, err := w.readWsMessages()
	if err != nil {
		//logging.Infof("Error reading message! %v", err)
		return nil, err
	}
	if messages == nil || len(messages) <= 0 { //没有读到完整数据 不处理
		return
	}
	for _, message := range messages {
		if message.OpCode.IsControl() {
			err = wsutil.HandleClientControlMessage(c, message)
			if err != nil {
				return
			}
			continue
		}
		if message.OpCode == ws.OpText || message.OpCode == ws.OpBinary {
			outs = append(outs, message)
		}
	}
	return
}
func (w *wsCodec) readWsMessages() (messages []wsutil.Message, err error) {
	msgBuf := &w.wsMsgBuf
	in := &w.buf
	for {
		if msgBuf.curHeader == nil {
			if in.Len() < ws.MinHeaderSize { //头长度至少是2
				return
			}
			var head ws.Header
			if in.Len() >= ws.MaxHeaderSize {
				head, err = ws.ReadHeader(in)
				if err != nil {
					return messages, err
				}
			} else { //有可能不完整，构建新的 reader 读取 head 读取成功才实际对 in 进行读操作
				tmpReader := bytes.NewReader(in.Bytes())
				oldLen := tmpReader.Len()
				head, err = ws.ReadHeader(tmpReader)
				skipN := oldLen - tmpReader.Len()
				if err != nil {
					if err == io.EOF || err == io.ErrUnexpectedEOF { //数据不完整
						return messages, nil
					}
					in.Next(skipN)
					return nil, err
				}
				in.Next(skipN)
			}

			msgBuf.curHeader = &head
			err = ws.WriteHeader(&msgBuf.cachedBuf, head)
			if err != nil {
				return nil, err
			}
		}
		dataLen := (int)(msgBuf.curHeader.Length)
		if dataLen > 0 {
			if in.Len() >= dataLen {
				_, err = io.CopyN(&msgBuf.cachedBuf, in, int64(dataLen))
				if err != nil {
					return
				}
			} else { //数据不完整
				//fmt.Println(in.Len(), dataLen)
				//logging.Infof("incomplete data")
				return
			}
		}
		if msgBuf.curHeader.Fin { //当前 header 已经是一个完整消息
			messages, err = wsutil.ReadClientMessage(&msgBuf.cachedBuf, messages)
			if err != nil {
				return nil, err
			}
			msgBuf.cachedBuf.Reset()
		} else {
			//logging.Infof("The data is split into multiple frames")
		}
		msgBuf.curHeader = nil
	}
}
