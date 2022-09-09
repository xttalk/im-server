package logic_model

import (
	"golang.org/x/net/context"
)

type ConnDevice struct {
	Context    context.Context
	SessionId  string      //Socket设备标识
	ServerId   uint32      //服务器ID
	UserClient *UserClient //用户客户端鉴权信息
}
