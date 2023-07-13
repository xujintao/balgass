package game

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/object"
)

func init() {
	Game.init()
}

var Game game

type game struct {
	connRequestChan      chan *connRequest
	closeConnRequestChan chan *closeConnRequest
	playerActionChan     chan *playerAction
	cancel               context.CancelFunc
}

func (g *game) init() {
	g.connRequestChan = make(chan *connRequest, 100)
	g.closeConnRequestChan = make(chan *closeConnRequest, 100)
	g.playerActionChan = make(chan *playerAction, 1000)
}

func (g *game) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel
	go func() {
		for {
			select {
			case connReq := <-g.connRequestChan:
				conn := connReq.Conn
				id, err := object.ObjectManager.AddPlayer(conn)
				connResp := connResponse{id: id, err: err}
				connReq.connResponseChan <- &connResp
			case closeConnReq := <-g.closeConnRequestChan:
				id := closeConnReq.id
				object.ObjectManager.DeletePlayer(id)
				closeConnReq.closeConnResponseChan <- struct{}{}
			case playerAction := <-g.playerActionChan:
				id := playerAction.id
				action := playerAction.action
				msg := playerAction.msg
				player := object.ObjectManager.GetPlayer(id)
				// player.Chat(msg)
				in := []reflect.Value{reflect.ValueOf(msg)}
				reflect.ValueOf(player).MethodByName(action).Call(in)
			case <-time.Tick(time.Second):
				maps.MapManager.ProcessWeather(g)
			case <-ctx.Done():
				// todo
				return
			}
		}
	}()
}

func (g *game) Close() {
	g.cancel()
}

type connRequest struct {
	object.Conn
	connResponseChan chan *connResponse
}

type connResponse struct {
	id  int
	err error
}

func (g *game) Conn(conn object.Conn) (int, error) {
	connReq := connRequest{
		Conn:             conn,
		connResponseChan: make(chan *connResponse),
	}
	g.connRequestChan <- &connReq
	connResp := <-connReq.connResponseChan
	return connResp.id, connResp.err
}

type closeConnRequest struct {
	id                    int
	closeConnResponseChan chan struct{}
}

func (g *game) CloseConn(id int) {
	closeConnReq := closeConnRequest{
		id:                    id,
		closeConnResponseChan: make(chan struct{}),
	}
	g.closeConnRequestChan <- &closeConnReq
	<-closeConnReq.closeConnResponseChan
}

type playerAction struct {
	id     int
	action string
	msg    any
}

func (g *game) PlayerAction(id int, action string, msg any) {
	playerAction := playerAction{
		id:     id,
		action: action,
		msg:    msg,
	}
	g.playerActionChan <- &playerAction
}

func (g *game) SendWeather(number, weather int) {
	if number == 0 {
		log.Println(number, weather)
	}
}
