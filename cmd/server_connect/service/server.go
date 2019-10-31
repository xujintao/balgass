package service

import (
	"sync"
	"time"

	"github.com/xujintao/balgass/cmd/server_connect/conf"
)

// ServerInfo server register to connect server with detail info
type ServerInfo struct {
	code         uint16
	percent      byte
	playType     byte
	userCount    uint16
	accountCount uint16
	maxUserCount uint16
}

type serverInfot struct {
	si ServerInfo
	t  *time.Timer
}

type server struct {
	mu   sync.RWMutex
	sits map[int]*serverInfot
}

func (s *server) UpdateServerInfo(si *ServerInfo) {
	code := int(si.code)
	s.mu.Lock()
	if sit, ok := s.sits[code]; !ok {
		t := time.AfterFunc(5*time.Second, func() {
			s.mu.Lock()
			delete(s.sits, code)
			s.mu.Unlock()
		})
		s.sits[code] = &serverInfot{*si, t}
	} else {
		sit.t.Reset(5 * time.Second)
	}
	s.mu.Unlock()
}

// ServerInfoBrief server list use brief info
type ServerInfoBrief struct {
	Code     uint16
	Percent  byte
	PlayType byte
}

func (s *server) GetServerList() (sibs []*ServerInfoBrief) {
	s.mu.Lock()
	for _, sit := range s.sits {
		si := sit.si
		for _, sic := range conf.Config.Servers {
			if uint16(sic.Code) == si.code && sic.Visible == 1 {
				sib := &ServerInfoBrief{
					Code:     si.code,
					Percent:  si.percent,
					PlayType: si.playType,
				}
				sibs = append(sibs, sib)
				break
			}
		}
	}
	s.mu.Unlock()
	return
}

// ServerInfoConnect contains ip and port
type ServerInfoConnect struct {
	IP   string
	Port uint16
}

func (s *server) GetServerInfo(code int) *ServerInfoConnect {
	for _, v := range conf.Config.Servers {
		if v.Code == code {
			return &ServerInfoConnect{
				IP:   v.IP,
				Port: uint16(v.Port),
			}
		}
	}
	return nil
}
