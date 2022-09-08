package manager

import (
	"fmt"
	"sync"
)

var ClientManager *_ClientManager

// InitClientManager 初始化客户端socket管理器
func InitClientManager() {
	ClientManager = &_ClientManager{
		Clients:   make(map[string]*Client),
		FdClients: make(map[string]*Client),
	}
}

// ClientManager 客户端管理器
type _ClientManager struct {
	Clients    map[string]*Client //sessionId->Client
	FdClients  map[string]*Client //mode_fd -> Client
	ClientLock sync.RWMutex       //客户端锁
}

func (c *_ClientManager) AddClient(client *Client) {
	c.ClientLock.Lock()
	defer CoreServer.OnConnect(client)
	defer c.ClientLock.Unlock()
	c.Clients[client.SessionId] = client
	sessionId := fmt.Sprintf("%d_%d", client.ClientMode, client.conn.Fd())
	c.FdClients[sessionId] = client
}
func (c *_ClientManager) DelClient(client *Client) {
	c.ClientLock.Lock()
	defer CoreServer.OnClose(client)
	defer c.ClientLock.Unlock()
	sessionId := fmt.Sprintf("%d_%d", client.ClientMode, client.conn.Fd())
	delete(c.Clients, client.SessionId)
	delete(c.FdClients, sessionId)
}

// GetClientByFd 通过fd找到对应客户端
func (c *_ClientManager) GetClientByFd(mode ClientMode, fd int) *Client {
	c.ClientLock.Lock()
	defer c.ClientLock.Unlock()
	sessionKey := fmt.Sprintf("%d_%d", mode, fd)
	if client, has := c.FdClients[sessionKey]; has {
		return client
	}
	return nil
}

// GetClientBySession 通过SessionID找到对应客户端
func (c *_ClientManager) GetClientBySession(sessionId string) *Client {
	c.ClientLock.Lock()
	defer c.ClientLock.Unlock()
	if client, has := c.Clients[sessionId]; has {
		return client
	}
	return nil
}
