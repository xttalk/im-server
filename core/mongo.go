package core

import (
	"XtTalkServer/global"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

var Mongo = new(_mongo)

type _mongo struct {
}

// Initialize 初始化Mongo
func (_mongo) Initialize(ctx context.Context) {
	conf := global.Config.Mongo
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", conf.User, conf.Password, conf.Host, conf.Port)
	fmt.Println(uri)
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		glog.Fatalf(ctx, "连接Mongo失败: %s", err.Error())
	}
	if err := client.Ping(ctx, nil); err != nil {
		glog.Fatalf(ctx, "连接Mongo异常: %s", err.Error())
	}
	global.Mongo = client.Database(conf.Name)
	glog.Infof(ctx, "MongoDB连接成功")
}
