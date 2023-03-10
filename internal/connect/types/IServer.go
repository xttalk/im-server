package types

type IServer interface {
	//Start 启动服务器
	Start() error
	//GetPort 获取服务端口
	GetPort() int
	//GetProtocol 获取连接协议
	GetProtocol() string
}
