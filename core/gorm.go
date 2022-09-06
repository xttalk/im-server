package core

import (
	"XtTalkServer/global"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var Mysql = new(_Mysql)

type _Mysql struct {
}

func (_Mysql) Initialize(ctx context.Context) {
	conf := global.Config.Mysql
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?&parseTime=true&loc=%s&charset=%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Name, conf.Local, conf.Charset)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Warn, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.Prefix,
			SingularTable: true,
		},
		Logger: newLogger,
	}
	_db, err := gorm.Open(mysql.Open(dsn), &gormConfig)
	if err != nil {
		glog.Fatalf(ctx, "连接Mysql失败: %s", err.Error())
	}
	sqlDb, err := _db.DB()
	if err != nil {
		glog.Fatalf(ctx, "访问数据库失败: %s", err.Error())
	}
	sqlDb.SetMaxIdleConns(conf.MaxIdleConn)
	sqlDb.SetMaxOpenConns(conf.MaxOpenConn)
	sqlDb.SetConnMaxLifetime(time.Minute * time.Duration(conf.MaxLifeTime))
	global.Db = _db
	glog.Infof(ctx, "Mysql连接成功")
}
