package core

import (
	"XtTalkServer/global"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/os/glog"
)

var Redis = new(_redis)

type _redis struct {
}

func (_redis) Initialize(ctx context.Context) {
	conf := global.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.Db,
	})
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		glog.Fatalf(ctx, "连接Redis失败: %s", err.Error())
	}
	glog.Infof(ctx, "Redis连接成功: %s", pong)
	global.Redis = client
}
