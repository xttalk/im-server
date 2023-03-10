package server_protocol

import (
	"XtTalkServer/internal/connect/manager"
	"XtTalkServer/internal/connect/service"
	"XtTalkServer/internal/connect/types"
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/panjf2000/gnet/v2"
	"golang.org/x/net/context"
)

func CreateTcpServer(port int) types.IServer {
	return &Tcp{
		addr:     fmt.Sprintf("tcp://:%d", port),
		port:     port,
		protocol: "tcp",
		context:  gctx.New(),
	}
}

type Tcp struct {
	addr     string
	protocol string
	port     int
	context  context.Context
	gnet.BuiltinEventEngine
	codec *types.XtTalkTcpCodec
}

func (c *Tcp) Start() error {
	return gnet.Run(c, c.addr,
		gnet.WithMulticore(false),
		gnet.WithReuseAddr(true),
	)
}
func (c *Tcp) GetPort() int {
	return c.port
}
func (c *Tcp) GetProtocol() string {
	return c.protocol
}

func (c *Tcp) OnBoot(eng gnet.Engine) gnet.Action {
	glog.Infof(c.context, "Tcp服务端已启动")
	c.codec = new(types.XtTalkTcpCodec)
	//注册服务发现
	if err := service.ServiceDiscovery.RegiterConnect(c.GetProtocol(), c.GetPort()); err != nil {
		glog.Warningf(c.context, "[%s]注册连接层服务发现失败: %s", c.GetProtocol(), err.Error())
		return gnet.Shutdown
	}
	return gnet.None
}
func (c *Tcp) OnOpen(conn gnet.Conn) (out []byte, action gnet.Action) {
	client := manager.CreateTcpClient(conn)
	manager.ClientManager.AddClient(client)
	return
}
func (c *Tcp) OnClose(conn gnet.Conn, err error) (action gnet.Action) {
	client := manager.ClientManager.GetClientByFd(manager.TcpClientMode, conn.Fd())
	if client != nil {
		manager.ClientManager.DelClient(client)
	}
	return
}
func (c *Tcp) OnTraffic(conn gnet.Conn) (action gnet.Action) {
	head, data, err := c.codec.Decode(conn)
	if err != nil {
		glog.Warningf(c.context, "读取数据包失败: %s", err.Error())
		return gnet.Close
	}
	client := manager.ClientManager.GetClientByFd(manager.TcpClientMode, conn.Fd())
	if client != nil {
		manager.CoreServer.OnMessage(client, head, data)
	} else {
		glog.Warningf(c.context, "无法获取到TCP客户端: %s", conn.Fd())
		return gnet.Close
	}
	return
}
