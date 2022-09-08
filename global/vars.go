package global

import (
	"XtTalkServer/config"
	"XtTalkServer/modules/rabbit"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// 全局公共
var (
	Config   *config.Config = nil //配置文件
	Db       *gorm.DB       = nil //Mysql数据库
	Redis    *redis.Client  = nil //redis数据库
	RabbitMQ *rabbit.Pool   = nil //RabbitMQ连接,需要使用这个创建新的channel去调用
)
