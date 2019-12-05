package service

import (
	"log"
	"sync"
	"time"

	"github.com/xujintao/balgass/cmd/server_connect/conf"
	"github.com/xujintao/balgass/cmd/server_connect/model"
)

type serverInfo struct {
	sri *model.ServerRegisterInfo
	t   *time.Timer
}

type serverManager struct {
	mu  sync.RWMutex
	sis map[int]*serverInfo
}

func (s *serverManager) RegisterServer(sri *model.ServerRegisterInfo) {
	code := int(sri.Sli.Code)
	s.mu.Lock()
	if s.sis == nil {
		s.sis = make(map[int]*serverInfo)
	}
	if sit, ok := s.sis[code]; !ok {
		t := time.AfterFunc(5*time.Second, func() {
			s.mu.Lock()
			delete(s.sis, code)
			s.mu.Unlock()
			// log
			for _, sci := range conf.ServerList.Servers {
				if sci.Code == code {
					log.Printf("code:[%d] name:[%s] ip:[%s] port:[%d] visible:[%d] lost", sci.Code, sci.Name, sci.IP, sci.Port, sci.Visible)
					break
				}
			}
		})
		s.sis[code] = &serverInfo{sri, t}
		s.mu.Unlock()
		// log
		for _, sci := range conf.ServerList.Servers {
			if sci.Code == code {
				log.Printf("code:[%d] name:[%s] ip:[%s] port:[%d] visible:[%d] online", sci.Code, sci.Name, sci.IP, sci.Port, sci.Visible)
				break
			}
		}
	} else {
		sit.t.Reset(5 * time.Second)
		s.mu.Unlock()
	}
}

func (s *serverManager) GetServerList() (slis []*model.ServerListInfo) {
	s.mu.RLock()
	for _, si := range s.sis {
		sri := si.sri
		for _, sci := range conf.ServerList.Servers {
			if uint16(sci.Code) == sri.Sli.Code && sci.Visible == 1 {
				sli := sri.Sli // copy
				slis = append(slis, &sli)
			}
		}
	}
	s.mu.RUnlock()
	return
}

func (s *serverManager) GetServerInfo(code int) *model.ServerConnectInfo {
	for _, sci := range conf.ServerList.Servers {
		if sci.Code == code {
			return &model.ServerConnectInfo{
				IP:   sci.IP,
				Port: uint16(sci.Port),
			}
		}
	}
	return nil
}
