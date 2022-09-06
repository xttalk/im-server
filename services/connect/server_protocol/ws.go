package server_protocol

import (
	"XtTalkServer/services/connect/manager"
	"XtTalkServer/services/connect/service"
	"XtTalkServer/services/connect/types"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gorilla/websocket"
	"github.com/panjf2000/gnet/v2"
	"golang.org/x/net/context"
)

type WebsocketContext struct {
	IsWs bool
}
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
	conn.SetContext(new(WebsocketContext))
	//这里不能触发OnConnect事件,不然发消息会报错
	return nil, gnet.None
}

func (wss *WebSocket) OnClose(conn gnet.Conn, err error) (action gnet.Action) {
	client := manager.ClientManager.GetClientByFd(manager.WsClientMode, conn.Fd())
	if client != nil {
		manager.ClientManager.DelClient(client)
	}
	return gnet.None
}

func (c *WebSocket) OnTraffic(conn gnet.Conn) (action gnet.Action) {
	if !conn.Context().(*WebsocketContext).IsWs {
		_, err := ws.Upgrade(conn)
		if err != nil {
			return gnet.Close
		}
		conn.Context().(*WebsocketContext).IsWs = true
		var client = manager.CreateWsClient(conn)
		manager.ClientManager.AddClient(client)
		return gnet.None
	} else {
		dataBytes, _, err := wsutil.ReadClientData(conn)
		if err != nil {
			if _, ok := err.(wsutil.ClosedError); !ok {
				glog.Warningf(c.context, "读取Ws客户端消息失败: %s", err.Error())
			}
			return gnet.Close
		}
		//拆包
		head, data, err := c.codec.DecodeBytes(dataBytes)
		if err != nil {
			glog.Warningf(c.context, "读取数据包失败: %s", err.Error())
			return gnet.Close
		}
		client := manager.ClientManager.GetClientByFd(manager.WsClientMode, conn.Fd())
		if client != nil {
			manager.CoreServer.OnMessage(client, head, data)
			return gnet.None
		} else {
			//没有找到这个客户端
			glog.Warningf(c.context, "无法获取到WS客户端: %s", conn.Fd())
			return gnet.Close
		}
	}

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
