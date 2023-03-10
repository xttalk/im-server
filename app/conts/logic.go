package conts

const (
	RK_ClientAuth = "connect:%s" //设备鉴权,sessionId
	RK_UserDevice = "user:%d"    //用户设备列表,用户ID

	RK_MqRetry_PrivateMsgSave = "mq:retry:private_msg_save:%d" //私聊消息储存重发
	RK_MqRetry_PrivateMsg     = "mq:retry:private_msg:%d"      //私聊消息重发
	RkMqRetry                 = "mq:retry:%s"                  //MQ消息重发

	RK_MqRetry_FriendRequest = "mq:retry:friend_request" //好友验证消息重发
	RK_MqRetry_FriendChange  = "mq:retry:friend_change"  //好友变动消息重发

	RK_PrivateMsgSeq = "private:msg_seq:%s" //私聊消息的seq时序维护

)
