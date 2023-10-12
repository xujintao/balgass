package server

import (
	"encoding/xml"
	"log"
	"time"

	"github.com/xujintao/balgass/src/server_connect/conf"
	"github.com/xujintao/balgass/src/server_connect/service/model"
)

func init() {
	type ServerListConfig struct {
		XMLName xml.Name              `xml:"ServerList"`
		Servers []*model.ServerConfig `xml:"Server"`
	}
	var serverListConfig ServerListConfig
	conf.XML("IGC_ServerList.xml", &serverListConfig)
	ConfigTable = make(configTable)
	for _, c := range serverListConfig.Servers {
		ConfigTable[c.Code] = c
	}

	ServerManager.init()
}

var ConfigTable configTable

type configTable map[int]*model.ServerConfig

var ServerManager serverManager

type server struct {
	Config *model.ServerConfig
	State  model.ServerState
	t      *time.Timer
}

type actioner interface {
	ServerAction(string, any)
}

type serverManager struct {
	servers  map[int]*server
	actioner actioner
}

func (m *serverManager) init() {
	m.servers = make(map[int]*server)
}

func (m *serverManager) RegisterActioner(actioner actioner) {
	m.actioner = actioner
}

func (m *serverManager) Register(msg *model.MsgRegister) {
	if s, ok := m.servers[msg.Code]; !ok {
		// register
		config, ok := ConfigTable[msg.Code]
		if !ok {
			return
		}
		t := time.AfterFunc(5*time.Second, func() {
			m.actioner.ServerAction("Unregister", &model.MsgUnregister{Code: config.Code})
		})
		m.servers[msg.Code] = &server{
			Config: config,
			State:  msg.ServerState,
			t:      t,
		}
		log.Printf("[code]%d [name]%s [ip]%s [port]%d [visible]%v online\n",
			config.Code, config.Name, config.IP, config.Port, config.Visible)
	} else {
		// refresh
		s.State = msg.ServerState
		s.t.Reset(5 * time.Second)
	}
}

func (m *serverManager) Unregister(msg *model.MsgUnregister) {
	if s, ok := m.servers[msg.Code]; ok {
		delete(m.servers, msg.Code)
		config := s.Config
		log.Printf("[code]%d [name]%s [ip]%s [port]%d [visible]%v offline\n",
			config.Code, config.Name, config.IP, config.Port, config.Visible)
	}
}

func (m *serverManager) GetServerList() []*model.ServerState {
	var list []*model.ServerState
	for _, s := range m.servers {
		if s.Config.Visible {
			list = append(list, &s.State)
		}
	}
	return list
}

func (m *serverManager) GetServer(code int) *model.ServerConfig {
	s, ok := m.servers[code]
	if !ok {
		return nil
	}
	return s.Config
}
