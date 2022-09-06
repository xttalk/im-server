package utils

import (
	"fmt"
	"github.com/gogf/gf/v2/util/grand"
	"log"
	"net"
)

// GetLocalIp 获取当前内网网卡Ip,传参代表为指定网卡
func GetLocalIp(names ...string) net.IP {
	name := ""
	if len(names) > 0 {
		name = names[0]
	}
	nets, err := net.Interfaces()
	if err != nil {
		log.Fatalln("获取失败", err.Error())
	}
	for _, n := range nets {
		addrs, err := n.Addrs()
		if err != nil {
			continue
		}
		if !IsEmpty(name) {
			if name != n.Name {
				fmt.Println("jump")
				continue
			}
		}
		if 0 > len(addrs) {
			continue
		}
		for _, addr := range addrs {
			switch ip := addr.(type) {
			case *net.IPNet:
				if ip.IP.DefaultMask() != nil && ip.IP.To4() != nil && !ip.IP.IsLoopback() {
					return ip.IP
				}
			}
		}
	}
	return nil
}

// RandomStr 随机字符串
func RandomStr(size int) string {
	return grand.Str("0123456789qwertyuiopasdfghjklzxcvbnm", size)
}
