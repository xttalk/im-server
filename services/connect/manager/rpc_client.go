package manager

import (
	"XtTalkServer/global"
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	zkClient "github.com/rpcxio/rpcx-zookeeper/client"
	"github.com/smallnest/rpcx/client"
)

var RpcClient = new(rpcClient)

type rpcClient struct {
	logic client.XClient
}

func (c *rpcClient) GetLogicClient(ctx context.Context) (client.XClient, error) {
	if c.logic == nil {
		conf := global.Config
		zkAddr := conf.Zookeeper.Servers
		d, err := zkClient.NewZookeeperDiscovery("rpc_logic", "LogicRpcService", zkAddr, nil)
		if err != nil {
			glog.Fatalf(ctx, "创建RPC客户端失败: %s", err.Error())
			return nil, gerror.Wrapf(err, "创建RPC客户端失败")
		}
		c.logic = client.NewXClient("LogicRpcService", client.Failtry, client.RoundRobin, d, client.DefaultOption)
	}
	return c.logic, nil
}
