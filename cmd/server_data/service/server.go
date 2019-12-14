package service

import (
	"fmt"
	"sync"

	"github.com/xujintao/balgass/cmd/server_data/conf"
	"github.com/xujintao/balgass/cmd/server_data/model"
	"github.com/xujintao/balgass/network"
)

type server struct {
	addr           string
	conn           network.ConnWriter
	name           string
	code           uint16
	typ            uint8
	port           uint16
	vip            uint8
	maxMIDUseCount int
}

type serverManager struct {
	mu      sync.RWMutex
	servers map[string]*server
}

func (m *serverManager) ServerAdd(addr string, c network.ConnWriter) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.servers) == 0 {
		m.servers = make(map[string]*server)
	}
	if len(m.servers) >= conf.Server.MaxServerCount {
		return "", fmt.Errorf("[%s] reached maximum allowed connections", addr)
	}
	index := addr // use addr(ip:port) as index
	m.servers[index] = &server{addr: addr, conn: c}
	return index, nil
}

func (m *serverManager) ServerLogin(index string, msg *model.ServerLoginReq) (*model.ServerLoginRes, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if server, ok := m.servers[index]; ok {
		server.name = msg.Name
		server.code = uint16(msg.Code)
		server.port = uint16(msg.Port)
		server.typ = uint8(msg.Type)
		server.vip = uint8(msg.Vip)
		server.maxMIDUseCount = int(msg.MaxHWIDUseCount)
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
	defer m.mu.RUnlock()
	if server, ok := m.servers[index]; ok {
		return server.conn.Write(res)
	}
	return fmt.Errorf("%s is off line", index)
}

func (m *serverManager) ServerGetMIDUseCount(index string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	server, ok := m.servers[index]
	if !ok {
		panic("server manager")
	}
	return server.maxMIDUseCount
}

func (m *serverManager) ServerGetName(index string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	server, ok := m.servers[index]
	if !ok {
		panic("server manager")
	}
	return server.name
}

func (m *serverManager) ServerGetAddr(index string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	server, ok := m.servers[index]
	if !ok {
		panic("server manager")
	}
	return server.addr
}
