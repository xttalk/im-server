package rabbit

import (
	"errors"
	"time"
)

var (
	ErrInvalidConfig     = errors.New("invalid pool config")
	ErrFailedConnection  = errors.New("failed to establish connection")
	ErrConnectionMaximum = errors.New("the number of connections exceeds the maximum")
	ErrChannelMaximum    = errors.New("the number of channels exceeds the maximum")
	ErrGetChannelTimeOut = errors.New("get channel timeout")
)

type Config struct {
	Host              string // MQ的地址
	MinConn           int    // 最少建立的连接数
	MaxConn           int    // 最大建立的连接数
	MaxChannelPerConn int    // 每个连接最多建立的信道数量
	MaxLifetime       time.Duration
}
