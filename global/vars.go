package global

import (
	"XtTalkServer/config"
	"XtTalkServer/modules/rabbit"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// 全局公共
var (
	Config   *config.Config  = nil //配置文件
	Db       *gorm.DB        = nil //Mysql数据库
	Redis    *redis.Client   = nil //Redis数据库
	Mongo    *mongo.Database = nil //Mongo数据库,已经选好的数据库
	RabbitMQ *rabbit.Pool    = nil //RabbitMQ连接池
)
