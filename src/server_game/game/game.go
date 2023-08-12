package game

import (
	"context"
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
	connRequestChan        chan *connRequest
	closeConnRequestChan   chan *closeConnRequest
	playerActionChan       chan *playerAction
	wsConnRequestChan      chan *wsConnRequest
	closeWSConnRequestChan chan *closeWSConnRequest
	userActionChan         chan *userAction
	commandRequestChan     chan *commandRequest
	cancel                 context.CancelFunc
}

func (g *game) init() {
	g.connRequestChan = make(chan *connRequest, 100)
	g.closeConnRequestChan = make(chan *closeConnRequest, 100)
	g.playerActionChan = make(chan *playerAction, 1000)
	g.wsConnRequestChan = make(chan *wsConnRequest, 100)
	g.closeWSConnRequestChan = make(chan *closeWSConnRequest, 100)
	g.userActionChan = make(chan *userAction, 1000)
	g.commandRequestChan = make(chan *commandRequest, 100)
}

func (g *game) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel
	go func() {
		t100ms := time.NewTicker(time.Millisecond * 100)
		cnt := 0
		for {
			select {
			// c1c2
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
			// websocket
			case wsConnReq := <-g.wsConnRequestChan:
				conn := wsConnReq.Conn
				id, err := object.ObjectManager.AddUser(conn)
				wsConnResp := wsConnResponse{id: id, err: err}
				wsConnReq.wsConnResponseChan <- &wsConnResp
			case closeWSConnReq := <-g.closeWSConnRequestChan:
				id := closeWSConnReq.id
				object.ObjectManager.DeleteUser(id)
				closeWSConnReq.closeWSConnResponseChan <- struct{}{}
			case userAction := <-g.userActionChan:
				id := userAction.id
				action := userAction.action
				msg := userAction.msg
				user := object.ObjectManager.GetUser(id)
				// user.Subscribe(msg)
				in := []reflect.Value{reflect.ValueOf(msg)}
				reflect.ValueOf(user).MethodByName(action).Call(in)
			// command
			case commandReq := <-g.commandRequestChan:
				name := commandReq.name
				msg := commandReq.msg
				in := []reflect.Value{reflect.ValueOf(msg)}
				out := reflect.ValueOf(g).MethodByName(name).Call(in)
				commandResp := commandResponse{
					data: out[0].Interface(),
				}
				if ierr := out[1].Interface(); ierr != nil {
					commandResp.err = ierr.(error)
				}
				commandReq.commandResponseChan <- &commandResp
			case <-t100ms.C:
				cnt++
				// fmt.Println(time.Now().Format("2006-01-02 15:04:05.999999"))
				// start := time.Now()
				object.ObjectManager.Process100ms()
				// fmt.Println("100ms", time.Since(start).Milliseconds())
				if cnt%10 == 0 { // 1s
					// start := time.Now()
					maps.MapManager.ProcessWeather(g)
					object.ObjectManager.Process1000ms()
					// fmt.Println("1000ms", time.Since(start).Milliseconds())
				}
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

// websocket
type wsConnRequest struct {
	object.Conn
	wsConnResponseChan chan *wsConnResponse
}

type wsConnResponse struct {
	id  int
	err error
}

func (g *game) WSConn(conn object.Conn) (int, error) {
	wsConnReq := wsConnRequest{
		Conn:               conn,
		wsConnResponseChan: make(chan *wsConnResponse),
	}
	g.wsConnRequestChan <- &wsConnReq
	wsConnResp := <-wsConnReq.wsConnResponseChan
	return wsConnResp.id, wsConnResp.err
}

type closeWSConnRequest struct {
	id                      int
	closeWSConnResponseChan chan struct{}
}

func (g *game) CloseWSConn(id int) {
	closeWSConnReq := closeWSConnRequest{
		id:                      id,
		closeWSConnResponseChan: make(chan struct{}),
	}
	g.closeWSConnRequestChan <- &closeWSConnReq
	<-closeWSConnReq.closeWSConnResponseChan
}

type userAction struct {
	id     int
	action string
	msg    any
}

func (g *game) UserAction(id int, action string, msg any) {
	userAction := userAction{
		id:     id,
		action: action,
		msg:    msg,
	}
	g.userActionChan <- &userAction
}

type commandRequest struct {
	name                string
	msg                 any
	commandResponseChan chan *commandResponse
}

type commandResponse struct {
	data any
	err  error
}

func (g *game) Command(name string, msg any) (any, error) {
	commandReq := commandRequest{
		name:                name,
		msg:                 msg,
		commandResponseChan: make(chan *commandResponse),
	}
	g.commandRequestChan <- &commandReq
	commandResp := <-commandReq.commandResponseChan
	return commandResp.data, commandResp.err
}

// func (g *game) GetObjectsByMapNumber(msg *model.MsgSubscribeMap) (*model.MsgSubscribeMapReply, error) {
// 	number := msg.Number
// 	pots := object.ObjectManager.GetObjectsByMapNumber(number)
// 	return &model.MsgSubscribeMapReply{
// 		Name: "object",
// 		Data: pots,
// 	}, nil
// }

func (g *game) SendWeather(number, weather int) {
	// if number == 0 {
	// 	log.Println(number, weather)
	// }
}
