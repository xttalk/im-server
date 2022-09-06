package core

//
//import (
//	"context"
//	"fmt"
//	"XtTalkServer/global"
//	"github.com/gogf/gf/v2/os/glog"
//	"github.com/smallnest/rpcx/manager"
//)
//
//var RpcClient = new(_rpcClient)
//
//type _rpcClient struct {
//}
//
//func (_rpcClient) Initialize(ctx context.Context) {
//	conf := global.Config.Rpc
//	d, err := manager.NewPeer2PeerDiscovery(fmt.Sprintf("tcp@%s", conf.Addr), "")
//	if err != nil {
//		glog.Fatalf(ctx, "连接RPC服务端失败: %s", err.Error())
//	}
//	global.RpcClient = manager.NewXClient("RpcDeviceService", manager.Failtry, manager.RandomSelect, d, manager.DefaultOption)
//}
