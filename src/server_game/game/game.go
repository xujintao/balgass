package game

import (
	"context"
	"reflect"

	"github.com/xujintao/balgass/src/server_game/game/object"
)

var (
	connRequestChan      chan *connRequest
	closeConnRequestChan chan *closeConnRequest
	playerActionChan     chan *playerAction
	cancel               context.CancelFunc
)

func init() {
	connRequestChan = make(chan *connRequest, 100)
	closeConnRequestChan = make(chan *closeConnRequest, 100)
	playerActionChan = make(chan *playerAction, 1000)
}

func Start() {
	ctx, cancel2 := context.WithCancel(context.Background())
	cancel = cancel2
	go func() {
		for {
			select {
			case connReq := <-connRequestChan:
				conn := connReq.Conn
				id, err := object.ObjectManager.AddPlayer(conn)
				connResp := connResponse{id: id, err: err}
				connReq.connResponseChan <- &connResp
			case closeConnReq := <-closeConnRequestChan:
				id := closeConnReq.id
				object.ObjectManager.DeletePlayer(id)
				closeConnReq.closeConnResponseChan <- struct{}{}
			case playerAction := <-playerActionChan:
				id := playerAction.id
				action := playerAction.action
				msg := playerAction.msg
				player := object.ObjectManager.GetPlayer(id)
				// player.Chat(msg)
				in := []reflect.Value{reflect.ValueOf(msg)}
				reflect.ValueOf(player).MethodByName(action).Call(in)
			case <-ctx.Done():
				// todo
				return
			}
		}
	}()
}

func Close() {
	cancel()
}

type connRequest struct {
	object.Conn
	connResponseChan chan *connResponse
}

type connResponse struct {
	id  int
	err error
}

func Conn(conn object.Conn) (int, error) {
	connReq := connRequest{
		Conn:             conn,
		connResponseChan: make(chan *connResponse),
	}
	connRequestChan <- &connReq
	connResp := <-connReq.connResponseChan
	return connResp.id, connResp.err
}

type closeConnRequest struct {
	id                    int
	closeConnResponseChan chan struct{}
}

func CloseConn(id int) {
	closeConnReq := closeConnRequest{
		id:                    id,
		closeConnResponseChan: make(chan struct{}),
	}
	closeConnRequestChan <- &closeConnReq
	<-closeConnReq.closeConnResponseChan
}

type playerAction struct {
	id     int
	action string
	msg    any
}

func PlayerAction(id int, action string, msg any) {
	playerAction := playerAction{
		id:     id,
		action: action,
		msg:    msg,
	}
	playerActionChan <- &playerAction
}
