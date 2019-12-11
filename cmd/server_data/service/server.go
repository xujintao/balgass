package service

import (
	"fmt"
	"sync"

	"github.com/xujintao/balgass/cmd/server_data/conf"
	"github.com/xujintao/balgass/cmd/server_data/model"
	"github.com/xujintao/balgass/network"
)

type serverManager struct {
	mu      sync.RWMutex
	servers map[string]*model.Server
}

func (m *serverManager) ServerAdd(addr string, conn network.ConnWriter) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.servers) == 0 {
		m.servers = make(map[string]*model.Server)
	}
	if len(m.servers) >= conf.Server.MaxServerCount {
		return "", fmt.Errorf("[%s] reached maximum allowed connections", addr)
	}
	index := addr // use addr(ip:port) as index
	m.servers[index] = &model.Server{Addr: addr, Conn: conn}
	return index, nil
}

func (m *serverManager) ServerGet(index string) *model.Server {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.servers[index]
}

func (m *serverManager) ServerLogin(index string, msg *model.ServerLoginReq) (*model.ServerLoginRes, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if server, ok := m.servers[index]; ok {
		server.Code = msg.Code
		server.Name = msg.Name
		server.Port = msg.Port
		server.Type = msg.Type
		server.Vip = msg.Vip
		return &model.ServerLoginRes{Result: 1}, nil
	}
	return nil, fmt.Errorf("%s is off line", index)
}

func (m *serverManager) ServerDel(index string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.servers, index)
}

// Send send package to server with index
func (m *serverManager) Send(index string, res *network.Response) error {
	m.mu.RLock()
	m.mu.RUnlock()
	if server, ok := m.servers[index]; ok {
		return server.Conn.Write(res)
	}
	return fmt.Errorf("%s is off line", index)
}
