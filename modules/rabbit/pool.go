package rabbit

import (
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	"sync"
	"sync/atomic"
	"time"
)

// Pool 连接池
type Pool struct {
	ctx                    context.Context
	mu                     *sync.Mutex
	conf                   *Config
	connectionNum          int32
	connections            map[int64]*Connection
	connectionSerialNumber int64
	idleChannels           chan *Channel
}

func NewPool(conf *Config) (*Pool, error) {
	if conf.MaxConn <= 0 || conf.MinConn > conf.MaxConn {
		return nil, ErrInvalidConfig
	}
	p := &Pool{
		ctx:          gctx.New(),
		mu:           new(sync.Mutex),
		connections:  make(map[int64]*Connection),
		idleChannels: make(chan *Channel, conf.MaxConn*conf.MaxChannelPerConn),
	}

	if conf.MaxLifetime == 0 {
		conf.MaxLifetime = time.Duration(3600)
	}

	p.conf = conf

	var conn *Connection
	var err error
	// 建立最少连接数
	for i := 0; i < conf.MinConn; i++ {
		conn, err = p.NewConnection()
		if err != nil {
			glog.Warningf(p.ctx, ErrFailedConnection.Error())
			return nil, ErrFailedConnection
		}
		p.connections[conn.connIdentity] = conn
	}
	return p, nil
}

func (p *Pool) SetConfig(conf *Config) *Pool {
	p.conf = conf
	return p
}

func (p *Pool) GetConfig() *Config {
	return p.conf
}

func (p *Pool) NewConnection() (*Connection, error) {
	// 判断连接是否达到最大值
	if atomic.AddInt32(&p.connectionNum, 1) > int32(p.conf.MaxConn) {
		atomic.AddInt32(&p.connectionNum, -1)
		return nil, ErrConnectionMaximum
	}
	conn, err := amqp.Dial(p.conf.Host)
	if err != nil {
		atomic.AddInt32(&p.connectionNum, -1)
		return nil, err
	}

	return &Connection{
		mu:           new(sync.Mutex),
		conn:         conn,
		pool:         p,
		channelNum:   0,
		expireTime:   time.Duration(time.Now().Unix()) + p.conf.MaxLifetime,
		connIdentity: atomic.AddInt64(&p.connectionSerialNumber, 1),
		isClose:      false,
	}, nil
}

func (p *Pool) CloseConnection(c *Connection) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	atomic.AddInt32(&p.connectionNum, -1)
	delete(p.connections, c.connIdentity)
	return c.conn.Close()
}

func (p *Pool) GetChannel() (*Channel, error) {
	ch, _ := p.getOrCreate()
	if ch != nil {
		return ch, nil
	}

	C := time.After(time.Second * 60)
	for {
		ch, _ := p.getOrCreate()
		if ch != nil {
			return ch, nil
		}
		select {
		case <-C:
			glog.Warningf(p.ctx, ErrGetChannelTimeOut.Error())
			return nil, ErrGetChannelTimeOut
		default:
			time.Sleep(time.Second)
		}
	}
}

func (p *Pool) getOrCreate() (*Channel, error) {
	// 池中是否有空闲channel
	var (
		ch  *Channel
		err error
	)
	select {
	case ch = <-p.idleChannels:
		return ch, nil
	default:
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	// 池中已有连接是否可以建立新的channel
	for _, conn := range p.connections {
		if conn.CheckExpire() {
			continue
		}
		ch, err = conn.NewChannel()
		if ch != nil {
			return ch, nil
		}
	}
	// 新建连接获取新的channel
	var conn *Connection
	conn, err = p.NewConnection()
	if err != nil {
		return nil, err
	}
	p.connections[conn.connIdentity] = conn
	ch, err = conn.NewChannel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (p *Pool) ReleaseChannel(ch *Channel) error {
	p.idleChannels <- ch
	return nil
}
