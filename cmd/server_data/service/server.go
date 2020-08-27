package service

import (
	"fmt"
	"strings"
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
	userCount      int
	userCountMax   int
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

func (m *serverManager) ServerLogin(index interface{}, msg *model.ServerLoginReq) (*model.ServerLoginRes, error) {
	sid := index.(string)
	m.mu.RLock()
	defer m.mu.RUnlock()
	if server, ok := m.servers[sid]; ok {
		server.name = msg.Name
		server.code = uint16(msg.Code)
		server.port = uint16(msg.Port)
		server.typ = uint8(msg.Type)
		server.vip = uint8(msg.Vip)
		server.maxMIDUseCount = int(msg.MaxHWIDUseCount)
		return &model.ServerLoginRes{Result: 1}, nil
	}
	return nil, fmt.Errorf("%s is off line", sid)
}

func (m *serverManager) ServerExit(index interface{}) {
	UserManager.UserDelete(index)

	sid := index.(string)
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.servers, sid)
}

func (m *serverManager) ServerCodeValidate(code uint16) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range m.servers {
		if s.code == code {
			return true
		}
	}
	return false
}

func (m *serverManager) ServerUserCountLimit(code uint16) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range m.servers {
		if s.code == code {
			return s.userCount >= s.userCountMax
		}
	}
	return true
}

func (m *serverManager) ServerUserCountSet(index interface{}, req *model.ServerUserCountSetReq) error {
	sid := index.(string)
	m.mu.Lock()
	defer m.mu.Unlock()
	server, ok := m.servers[sid]
	if !ok {
		panic("server manager")
	}
	server.userCount = int(req.UserCount)
	server.userCountMax = int(req.UserCountMax)
	return nil
}

// Send send package to server with index
func (m *serverManager) Send(index interface{}, res *network.Response) error {
	sid := index.(string)
	m.mu.RLock()
	defer m.mu.RUnlock()
	if server, ok := m.servers[sid]; ok {
		return server.conn.Write(res)
	}
	return fmt.Errorf("%s is off line", index)
}

func (m *serverManager) ServerGetMIDUseCount(index interface{}) int {
	sid := index.(string)
	m.mu.RLock()
	defer m.mu.RUnlock()
	server, ok := m.servers[sid]
	if !ok {
		panic("server manager")
	}
	return server.maxMIDUseCount
}

func (m *serverManager) ServerGetName(index interface{}) string {
	sid := index.(string)
	m.mu.RLock()
	defer m.mu.RUnlock()
	server, ok := m.servers[sid]
	if !ok {
		panic("server manager")
	}
	return server.name
}

func (m *serverManager) ServerGetAddr(index interface{}) string {
	sid := index.(string)
	m.mu.RLock()
	defer m.mu.RUnlock()
	server, ok := m.servers[sid]
	if !ok {
		panic("server manager")
	}
	return server.addr
}

func (m *serverManager) ServerGetIP(index interface{}) string {
	sid := index.(string)
	m.mu.RLock()
	defer m.mu.RUnlock()
	server, ok := m.servers[sid]
	if !ok {
		panic("server manager")
	}
	addr := server.addr
	ip := addr[:strings.Index(addr, ":")]
	return ip
}
