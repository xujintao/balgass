package object

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/xujintao/balgass/src/server-connect/conf"
	"github.com/xujintao/balgass/src/server-connect/service/model"
	"github.com/xujintao/balgass/src/server-connect/service/server"
)

func init() {
	conf.INI(conf.PathConfig, "IGCCS.ini", "AutoUpdate", &AutoUpdate)
	_, err := fmt.Sscanf(AutoUpdate.VerStr, "%d.%d.%d",
		&AutoUpdate.Ver.Major, &AutoUpdate.Ver.Minor, &AutoUpdate.Ver.Patch)
	if err != nil {
		slog.Error("fmt.Sscanf Failed", "err", err)
		os.Exit(1)
	}
}

func NewPlayer(conn Conn) *player {
	ctx, cancel := context.WithCancel(context.Background())
	p := player{
		conn:    conn,
		msgChan: make(chan any, 100),
		cancel:  cancel,
	}
	go func() {
		for {
			select {
			case msg := <-p.msgChan:
				err := p.conn.Write(msg)
				if err != nil {
					slog.Error("p.conn.Write Failed",
						"err", err, "msg", msg)
				}
			case <-ctx.Done():
				close(p.msgChan)
				p.conn.Close()
				return
			}
		}
	}()
	return &p
}

type player struct {
	objectManager *objectManager
	index         int
	offline       bool
	conn          Conn
	msgChan       chan any
	cancel        context.CancelFunc
}

func (p *player) Offline() {
	if p.offline {
		return
	}
	p.offline = true
	// todo
	p.cancel()
}

func (p *player) Push(msg any) {
	if p.offline {
		slog.Warn("Still pushing",
			"msg", msg, "player", p.index)
		return
	}
	if len(p.msgChan) > 80 {
		p.Offline()
		return
	}
	p.msgChan <- msg
}

var AutoUpdate model.AutoUpdateConfig

func (p *player) CheckVersion(msg *model.MsgCheckVersion) {
	if msg.Version == AutoUpdate.Ver {
		resp := model.MsgCheckVersionSuccess{Result: true}
		p.Push(&resp)
	} else {
		resp := model.MsgCheckVersionFailed{AutoUpdateConfig: &AutoUpdate}
		p.Push(&resp)
	}
}

func (p *player) GetServerList(msg *model.MsgGetServerList) {
	servers := server.ServerManager.GetServerList()
	resp := model.MsgGetServerListReply{
		ServerList: servers,
	}
	p.Push(&resp)
}

func (p *player) GetServer(msg *model.MsgGetServer) {
	server := server.ServerManager.GetServer(msg.Code)
	resp := model.MsgGetServerReply{
		Server: *server,
	}
	p.Push(&resp)
}
