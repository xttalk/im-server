package service

import (
	"XtTalkServer/global"
	"XtTalkServer/utils"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/samuel/go-zookeeper/zk"
	"golang.org/x/net/context"
	"time"
)

var ServiceDiscovery *_ServiceDiscovery = nil

type _ServiceDiscovery struct {
	context   context.Context
	conn      *zk.Conn
	paths     []string
	root_path string //根路径
	localIp   string //本地IP
}

var (
	ZkAcls = zk.WorldACL(zk.PermAll)
)

func InitDiscovery() error {
	ctx := gctx.New()
	ip := utils.GetLocalIp()
	if ip == nil {
		return gerror.Newf("无法获取到当前系统IP地址")
	}
	//初始化zk服务
	conf := global.Config.Zookeeper
	if 0 >= len(conf.Servers) {
		return gerror.Newf("未配置Zookeeper服务信息")
	}
	conn, _, err := zk.Connect(conf.Servers, time.Second*5)
	if err != nil {
		return gerror.Newf("连接Zookeeper服务失败: %s", err.Error())
	}
	zkRootDir := conf.Root
	if utils.IsEmpty(zkRootDir) {
		zkRootDir = "/XtTalk"
	}
	//创建目录
	if has, _, _ := conn.Exists(zkRootDir); !has {
		if _, err := conn.Create(zkRootDir, nil, 0, ZkAcls); err != nil {
			return gerror.Newf("注册Zookeeper目录失败: %s", err.Error())
		}
	}

	ServiceDiscovery = &_ServiceDiscovery{
		context:   ctx,
		conn:      conn,
		root_path: zkRootDir,
	}
	return nil
}

// RegiterConnect 注册连接服务
func (c *_ServiceDiscovery) RegiterConnect(protocol string, port int) error {
	conf := global.Config.Services.Connect

	//写入ws地址
	data := utils.DataToJsonBytes(ConnectEntity{
		Id:       conf.Id,
		Protocol: protocol,
		Host:     c.localIp,
		Port:     port,
	})
	path := fmt.Sprintf("%s/%s_%d", c.root_path, protocol, conf.Id)
	if _, err := c.conn.Create(path, data, zk.FlagEphemeral, ZkAcls); err != nil {
		return err
	}
	c.paths = append(c.paths, path)
	return nil
}
func (c *_ServiceDiscovery) UnRegister() {
	for _, path := range c.paths {
		_, sate, err := c.conn.Get(path)
		if err != nil {
			glog.Warningf(c.context, "获取Zookeeper路径失败: %s -> %s", path, err.Error())
			continue
		}
		if err := c.conn.Delete(path, sate.Version); err != nil {
			glog.Warningf(c.context, "移除Zookeeper路径失败: %s", err.Error())
		}
	}
}
