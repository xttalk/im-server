package core

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
)

var Log = new(_log)

type _log struct {
}

//Initialize 初始化Log
func (_log) Initialize(ctx context.Context) {
	if err := glog.SetConfigWithMap(map[string]interface{}{
		"path":     "./logs",
		"file":     "{Y-m-d}.log",
		"level":    "all",
		"stdout":   true,
		"StStatus": 0,
	}); err != nil {
		glog.Fatalf(ctx, "日志组件初始化失败: %s", err.Error())
	}
	glog.Info(ctx, "日志组件初始化成功")
}
