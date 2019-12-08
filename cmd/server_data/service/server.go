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

func (sm *serverManager) ServerAdd(addr string, conn network.ConnWriter) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if len(sm.servers) == 0 {
		sm.servers = make(map[string]*model.Server)
	}
	if len(sm.servers) >= conf.Server.MaxServerCount {
		return fmt.Errorf("[%s] reached maximum allowed connections", addr)
	}
	sm.servers[addr] = &model.Server{Addr: addr, Conn: conn}
	return nil
}

func (sm *serverManager) ServerGet(addr string) *model.Server {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.servers[addr]
}

func (*serverManager) ServerLogin(server *model.Server, msg *model.ServerLoginMsg) {
	server.Code = msg.Code
	server.Name = msg.Name
	server.Port = msg.Port
	server.Type = msg.Type
	server.Vip = msg.Vip
}

func (sm *serverManager) ServerDel(addr string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.servers, addr)
}
