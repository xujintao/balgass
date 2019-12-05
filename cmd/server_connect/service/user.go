package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/xujintao/balgass/cmd/server_connect/conf"
	"github.com/xujintao/balgass/cmd/server_connect/model"
	"github.com/xujintao/balgass/network"
)

type userManager struct {
	mu    sync.Mutex
	users map[string]*model.User
}

func (um *userManager) AddUser(addr string, conn network.ConnWriter) error {
	um.mu.Lock()
	defer um.mu.Unlock()
	if len(um.users) == 0 {
		um.users = make(map[string]*model.User)
	}

	ipnew := addr[0:strings.LastIndex(addr, ":")]
	sum := 0
	for addr := range um.users {
		ip := addr[0:strings.LastIndex(addr, ":")]
		if ip == ipnew {
			sum++
		}
		if sum >= conf.Net.MaxConnectionsPerIP {
			return fmt.Errorf("[%s] reached maximum allowed connections", ipnew)
		}
	}
	um.users[addr] = &model.User{Addr: addr, Conn: conn}
	return nil
}

func (um *userManager) GetUser(addr string) *model.User {
	um.mu.Lock()
	defer um.mu.Unlock()
	return um.users[addr]
}

func (um *userManager) DelUser(addr string) {
	um.mu.Lock()
	defer um.mu.Unlock()
	delete(um.users, addr)
}
