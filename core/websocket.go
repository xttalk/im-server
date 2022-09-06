package core

import (
	"context"
)

var Websocket = new(_websocket)

type _websocket struct {
}

func (_websocket) Initialize(ctx context.Context) {
	//	//初始化websocket服务
	//	websocket.InitCreateWebSocketServer()
	//	conf := global.Config.Websocket
	//	server_protocol := g.Server()
	//	server_protocol.SetAddr(conf.Addr)
	//	server_protocol.Use(middleware.Cors) //注册全局中间件
	//	server_protocol.Group("/", registerWsServerRouter)
	//	server_protocol.SetDumpRouterMap(false) //隐藏路由表打印
	//	go server_protocol.Run()
	//
}

//func registerWsServerRouter(group *ghttp.RouterGroup) {
//	group.ALL("/socket", websocket.Server.ClientEntry)
//}
