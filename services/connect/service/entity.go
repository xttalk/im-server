package service

type ConnectEntity struct {
	Id       uint32 `json:"id"`       //服务器Id
	Protocol string `json:"protocol"` //客户端连接协议
	Host     string `json:"host"`     //服务端IP
	Port     int    `json:"port"`     //服务端端口
}
