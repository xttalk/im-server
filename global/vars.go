package global

import (
	"XtTalkServer/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// 全局公共
var (
	Config *config.Config = nil //配置文件
	Db     *gorm.DB       = nil //Mysql数据库
	Redis  *redis.Client  = nil //redis数据库
)
