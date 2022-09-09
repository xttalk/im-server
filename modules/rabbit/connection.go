package rabbit

import (
	"github.com/gogf/gf/v2/os/glog"
	"github.com/streadway/amqp"
	"sync"
	"sync/atomic"
	"time"
)

type Connection struct {
	mu                  *sync.Mutex
	conn                *amqp.Connection
	pool                *Pool
	expireTime          time.Duration
	isExpire            bool
	connIdentity        int64 // 连接标记
	channelNum          int32 // 该连接的信道数量
	channelSerialNumber int64 // 第几个channel
	isClose             bool  //是否已经断开
}

func (c *Connection) NewChannel() (*Channel, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if atomic.AddInt32(&c.channelNum, 1) > int32(c.pool.conf.MaxChannelPerConn) {
		atomic.AddInt32(&c.channelNum, -1)
		return nil, ErrChannelMaximum
	}
	ch, err := c.conn.Channel()
	if err != nil {
		atomic.AddInt32(&c.channelNum, -1)
		return nil, err
	}
	channel := &Channel{
		Channel:      ch,
		conn:         c,
		chanIdentity: atomic.AddInt64(&c.channelSerialNumber, 1),
	}
	chCloseEvent := make(chan *amqp.Error)
	ch.NotifyClose(chCloseEvent)
	go func() {
		select {
		case err := <-chCloseEvent:
			if c != nil {
				c.isClose = true
			}
			glog.Warningf(c.pool.ctx, err.Error())
		}
	}()
	return channel, nil
}

func (c *Connection) ReleaseChannel(ch *Channel) error {
	if c.CheckExpire() {
		return c.CloseChannel(ch)
	}
	if c.isClose { //已经断开,不能再释放了
		return c.CloseChannel(ch)
	}
	return c.pool.ReleaseChannel(ch)
}

func (c *Connection) CloseChannel(ch *Channel) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	atomic.AddInt32(&c.channelNum, -1)
	var err = ch.Channel.Close()
	if atomic.LoadInt32(&c.channelNum) <= 0 && (c.CheckExpire() || c.isClose) {
		//if atomic.LoadInt32(&c.channelNum) <= 0 && c.CheckExpire() {
		return c.pool.CloseConnection(c)
	}
	return err
}

// 检查是否过期
func (c *Connection) CheckExpire() bool {
	if c.isExpire {
		return true
	}
	if time.Duration(time.Now().Unix()) > c.expireTime {
		c.isExpire = true
	}
	return c.isExpire
}
