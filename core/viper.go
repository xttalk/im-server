package core

import (
	"XtTalkServer/global"
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/spf13/viper"
	"os"
)

var Viper = new(_viper)

type _viper struct {
	Viper *viper.Viper
}

func (c *_viper) Initialize(ctx context.Context) {
	confName := "config"
	if len(os.Args) >= 2 {
		confName = os.Args[1]
	}
	c.Viper = viper.New()
	c.Viper.SetConfigName(confName)
	c.Viper.SetConfigType("yaml")
	c.Viper.AddConfigPath(".")
	if err := c.Viper.ReadInConfig(); err != nil {
		glog.Fatalf(ctx, "[配置文件] 加载失败: %s", err.Error())
	}
	c.Viper.WatchConfig()
	c.Viper.OnConfigChange(func(in fsnotify.Event) {
		if err := c.Viper.Unmarshal(&global.Config); err != nil {
			glog.Warningf(ctx, "配置文件重载失败: %s", err.Error())
		} else {
			glog.Infof(ctx, "配置文件重载完成: %s", in.Name)
		}
	})

	if err := c.Viper.Unmarshal(&global.Config); err != nil {
		glog.Fatalf(ctx, "解析配置文件失败: %s", err.Error())
	}
	glog.Info(ctx, "配置文件加载成功")
}
