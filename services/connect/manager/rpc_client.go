package manager

import (
	"XtTalkServer/global"
	"context"
	"github.com/gogf/gf/v2/os/glog"
	zkClient "github.com/rpcxio/rpcx-zookeeper/client"
	"github.com/smallnest/rpcx/client"
)

var LogicRpcClient client.XClient

func InitLogicRpcClient(ctx context.Context) {
	conf := global.Config
	zkAddr := conf.Zookeeper.Servers
	d, err := zkClient.NewZookeeperDiscovery("rpc_logic", "LogicRpcService", zkAddr, nil)
	if err != nil {
		glog.Fatalf(ctx, "初始化Logic_Rpc失败: %s", err.Error())
	}
	LogicRpcClient = client.NewXClient("LogicRpcService", client.Failtry, client.RandomSelect, d, client.DefaultOption)
}
