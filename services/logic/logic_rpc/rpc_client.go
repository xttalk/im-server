package logic_rpc

import (
	"XtTalkServer/global"
	"XtTalkServer/pb"
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"
	zkClient "github.com/rpcxio/rpcx-zookeeper/client"
	"github.com/smallnest/rpcx/client"
)

var RpcClient = new(rpcClient)

type rpcClient struct {
	connect client.XClient
}

func (c *rpcClient) GetConnectClient(ctx context.Context) (client.XClient, error) {
	if c.connect == nil {
		conf := global.Config
		zkAddr := conf.Zookeeper.Servers
		d, err := zkClient.NewZookeeperDiscovery("rpc_connect", "ConnectRpcService", zkAddr, nil)
		if err != nil {
			glog.Fatalf(ctx, "创建RPC客户端失败: %s", err.Error())
			return nil, gerror.Wrapf(err, "创建RPC客户端失败")
		}
		c.connect = client.NewXClient("ConnectRpcService", client.Failtry, client.RoundRobin, d, client.DefaultOption)
		c.connect.SetSelector(&UserSessionSelector{})
	}
	return c.connect, nil
}

type UserSessionSelector struct {
	servers map[uint32]string //服务器ID->服务器地址
}

func (s *UserSessionSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	serverId := func() uint32 {
		if req, has := args.(*pb.ConnectSendClientReq); has {
			return req.GetServerId()
		}
		return 0
	}()
	if serverId != 0 {
		if addr, has := s.servers[serverId]; has {
			return addr
		}
		//用户异常,存储的sid并未接入服务发现
		glog.Warningf(ctx, "RPC调用[%s]->[%s]目标服务器ID: %d 并未被服务发现注册!", servicePath, serviceMethod, serverId)
		return ""
	} else {
		//无法找到用户所在服务器
		glog.Warningf(ctx, "RPC调用[%s]->[%s]请求参数中无法找到对应接入层服务器", servicePath, serviceMethod)
		return ""
	}
}

func (s *UserSessionSelector) UpdateServer(servers map[string]string) {
	var connServer = make(map[uint32]string)
	for addr := range servers {
		item := gstr.Split(addr, "#")
		if len(item) >= 2 {
			serverId := gvar.New(item[0]).Uint32()
			serverAddr := gvar.New(item[1]).String()
			connServer[serverId] = serverAddr
		}
	}
	s.servers = connServer
}
