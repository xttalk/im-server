package core

//import (
//	"github.com/MikuAdmin/MikuAdmin-Golang-GoFrame/app/websocket"
//	"github.com/MikuAdmin/MikuAdmin-Golang-GoFrame/global"
//	"github.com/gogf/gf/v2/frame/g"
//	"github.com/gogf/gf/v2/os/gctx"
//	"github.com/gogf/gf/v2/os/glog"
//	"os"
//	"os/signal"
//	"syscall"
//)
//
//func init() {
//	g.Log().SetAsync(true)
//}
//
////Initialize 初始化所有初始服务
//func Initialize() {
//	ctx := gctx.New()
//	Log.Initialize(ctx)
//	glog.Infof(ctx, "Api版本号: %s", global.ApiVersion)
//	Viper.Initialize(ctx)
//	Mysql.Initialize(ctx)
//	Redis.Initialize(ctx)
//	AuthSystem.Initialize(ctx, global.ApiVersion, "api")
//	RpcClient.Initialize(ctx)
//	HttpServer.Initialize(ctx)
//}
//
//func InitializeServer() {
//	ctx := gctx.New()
//	Log.Initialize(ctx)
//	glog.Infof(ctx, "服务端版本号: %s", global.ServerVersion)
//	Viper.Initialize(ctx)
//	Mysql.Initialize(ctx)
//	Redis.Initialize(ctx)
//	AuthSystem.Initialize(ctx, global.ServerVersion, "server_protocol")
//	RpcServer.Initialize(ctx)
//	Websocket.Initialize(ctx)
//	sigc := make(chan os.Signal, 1)
//	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
//	_ = <-sigc
//	//统计当前在线的客户端在线时长
//	glog.Infof(ctx, "正在停止服务端,正在进行数据更新....(请勿强制中断该操作!!!)")
//	//更新当前在线设备的所有在线时间
//	websocket.Server.GetClientManager().UpdateOnlineTime(false)
//	glog.Infof(ctx, "服务端已停止")
//}
