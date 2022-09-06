package core

//
//import (
//	"context"
//	"XtTalkServer/global"
//	"github.com/gogf/gf/v2/os/glog"
//	"github.com/smallnest/rpcx/server_protocol"
//)
//
//var RpcServer = new(_rpcServer)
//
//type _rpcServer struct {
//}
//
//func (_rpcServer) Initialize(ctx context.Context) {
//	conf := global.Config.Rpc
//	s := server_protocol.NewServer()
//	if err := s.Register(new(websocket.RpcDeviceService), ""); err != nil {
//		glog.Fatalf(ctx, "注册RPC服务失败: %s", err.Error())
//	}
//	go func() {
//		if err := s.Serve("tcp", conf.Addr); err != nil {
//			glog.Fatalf(ctx, "启动RPC服务失败: %s", err.Error())
//		}
//	}()
//	glog.Infof(ctx, "RPC启动成功,端口启动在: %s", conf.Addr)
//}
