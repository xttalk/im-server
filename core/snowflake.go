package core

import (
	"XtTalkServer/utils/snowflake"
	"github.com/gogf/gf/v2/os/glog"
	"golang.org/x/net/context"
)

var Snowflake = new(_snowflake)

type _snowflake struct {
}

func (_snowflake) Initialize(serverId uint32, ctx context.Context) {
	if err := snowflake.InitSnowFlake(int64(serverId)); err != nil {
		glog.Fatalf(ctx, "初始化雪花算法失败: %s", err.Error())
	}

}
