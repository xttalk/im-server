package conts

import (
	"fmt"
	"github.com/gogf/gf/v2/text/gstr"
)

const (
	MqQueueFormat    = "queue_"    //队列
	MqExchangeFormat = "exchange_" //交换机

	//私聊消息存储
	MqPrivateMsgSave    = "private_msg_save"
	MqKeyPrivateMsgSave = "private_msg_save.*"
	//私聊消息推送
	MqPrivateMsg    = "private_msg"
	MqKeyPrivateMsg = "private_msg.*"

	//好友请求
	MqFriendRequest    = "friend_request"
	MqKeyFriendRequest = "friend_request.*"
	//好友变动
	MqFriendChange    = "friend_change"
	MqKeyFriendChange = "friend_change.*"
)

// GetQueueName 获取队列名称
func GetQueueName(eventName string) string {
	return fmt.Sprintf("%s%s", MqQueueFormat, eventName)
}

// GetExchangeName 获取交换机名称
func GetExchangeName(eventName string) string {
	return fmt.Sprintf("%s%s", MqExchangeFormat, eventName)
}
func GetRouteKey(keyFormat string, args ...interface{}) string {
	return fmt.Sprintf(gstr.Replace(keyFormat, "*", "%v"), args...)
}
