package logic_model

type UserClient struct {
	Uid       uint64 //用户ID,当用户ID为0时则代表没有登录
	SessionId string //设备标识SessionID
	ServerId  uint32 //所在服务器ID
}
