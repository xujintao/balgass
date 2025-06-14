package service

import (
	"context"
	"reflect"

	"github.com/xujintao/balgass/src/server-connect/service/object"
	"github.com/xujintao/balgass/src/server-connect/service/server"
)

func init() {
	Service.init()
}

var Service service

type service struct {
	playerConnRequestChan      chan *connRequest
	playerCloseConnRequestChan chan *closeConnRequest
	playerActionChan           chan *actionRequest
	serverActionChan           chan *actionRequest
	cancel                     context.CancelFunc
}

func (s *service) init() {
	s.playerConnRequestChan = make(chan *connRequest, 100)
	s.playerCloseConnRequestChan = make(chan *closeConnRequest, 100)
	s.playerActionChan = make(chan *actionRequest, 1000)
	s.serverActionChan = make(chan *actionRequest, 100)
	server.ServerManager.RegisterActioner(s)
}

func (s *service) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	go func() {
		for {
			select {
			// player
			case connReq := <-s.playerConnRequestChan:
				conn := connReq.Conn
				id, err := object.ObjectManager.AddPlayer(conn)
				connResp := connResponse{id: id, err: err}
				connReq.connResponseChan <- &connResp
			case closeConnReq := <-s.playerCloseConnRequestChan:
				id := closeConnReq.id
				object.ObjectManager.DeletePlayer(id)
				closeConnReq.closeConnResponseChan <- struct{}{}
			case playerAction := <-s.playerActionChan:
				id := playerAction.id
				action := playerAction.action
				msg := playerAction.msg
				player := object.ObjectManager.GetPlayer(id)
				if player == nil {
					break
				}
				// player.CheckVersion(msg)
				in := []reflect.Value{reflect.ValueOf(msg)}
				reflect.ValueOf(player).MethodByName(action).Call(in)
			case serverAction := <-s.serverActionChan:
				action := serverAction.action
				msg := serverAction.msg
				// server.ServerManager.Register(msg)
				in := []reflect.Value{reflect.ValueOf(msg)}
				reflect.ValueOf(&server.ServerManager).MethodByName(action).Call(in)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *service) Close() {
	s.cancel()
}

type connRequest struct {
	object.Conn
	connResponseChan chan *connResponse
}

type connResponse struct {
	id  int
	err error
}

func (s *service) PlayerConn(conn object.Conn) (int, error) {
	connReq := connRequest{
		Conn:             conn,
		connResponseChan: make(chan *connResponse),
	}
	s.playerConnRequestChan <- &connReq
	connResp := <-connReq.connResponseChan
	return connResp.id, connResp.err
}

type closeConnRequest struct {
	id                    int
	closeConnResponseChan chan struct{}
}

func (s *service) PlayerCloseConn(id int) {
	closeConnReq := closeConnRequest{
		id:                    id,
		closeConnResponseChan: make(chan struct{}),
	}
	s.playerCloseConnRequestChan <- &closeConnReq
	<-closeConnReq.closeConnResponseChan
}

type actionRequest struct {
	id     int
	action string
	msg    any
}

func (s *service) PlayerAction(id int, action string, msg any) {
	actionReq := actionRequest{
		id:     id,
		action: action,
		msg:    msg,
	}
	s.playerActionChan <- &actionReq
}

func (s *service) ServerAction(action string, msg any) {
	// log.Println(action, msg)
	actionReq := actionRequest{
		action: action,
		msg:    msg,
	}
	s.serverActionChan <- &actionReq
}
