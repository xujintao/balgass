package game

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
	"github.com/xujintao/balgass/src/server_game/game/object"
)

func init() {
	Game.init()
}

var Game game

type game struct {
	playerConnRequestChan      chan *connRequest
	playerCloseConnRequestChan chan *closeConnRequest
	playerActionChan           chan *actionRequest
	userConnRequestChan        chan *connRequest
	userCloseConnRequestChan   chan *closeConnRequest
	userActionChan             chan *actionRequest
	commandRequestChan         chan *commandRequest
	cancel                     context.CancelFunc
	bots                       map[string]*bot
	maxBotNumber               int
}

func (g *game) init() {
	g.playerConnRequestChan = make(chan *connRequest, 100)
	g.playerCloseConnRequestChan = make(chan *closeConnRequest, 100)
	g.playerActionChan = make(chan *actionRequest, 1000)
	g.userConnRequestChan = make(chan *connRequest, 100)
	g.userCloseConnRequestChan = make(chan *closeConnRequest, 100)
	g.userActionChan = make(chan *actionRequest, 1000)
	g.commandRequestChan = make(chan *commandRequest, 100)
	g.bots = make(map[string]*bot)
}

func (g *game) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel
	go func() {
		t100ms := time.NewTicker(time.Millisecond * 100)
		cnt := 0
		for {
			select {
			// player
			case connReq := <-g.playerConnRequestChan:
				conn := connReq.Conn
				id, err := object.ObjectManager.AddPlayer(conn, g)
				connResp := connResponse{id: id, err: err}
				connReq.connResponseChan <- &connResp
			case closeConnReq := <-g.playerCloseConnRequestChan:
				id := closeConnReq.id
				object.ObjectManager.DeletePlayer(id)
				closeConnReq.closeConnResponseChan <- struct{}{}
			case playerAction := <-g.playerActionChan:
				id := playerAction.id
				action := playerAction.action
				msg := playerAction.msg
				player := object.ObjectManager.GetPlayer(id)
				if player == nil {
					break
				}
				// player.Chat(msg)
				in := []reflect.Value{reflect.ValueOf(msg)}
				reflect.ValueOf(player).MethodByName(action).Call(in)
			// user
			case connReq := <-g.userConnRequestChan:
				conn := connReq.Conn
				id, err := object.ObjectManager.AddUser(conn)
				connResp := connResponse{id: id, err: err}
				connReq.connResponseChan <- &connResp
			case closeConnReq := <-g.userCloseConnRequestChan:
				id := closeConnReq.id
				object.ObjectManager.DeleteUser(id)
				closeConnReq.closeConnResponseChan <- struct{}{}
			case userAction := <-g.userActionChan:
				id := userAction.id
				action := userAction.action
				msg := userAction.msg
				user := object.ObjectManager.GetUser(id)
				if user == nil {
					break
				}
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
	g.addBot("bot1")
	g.addBot("bot2")
	g.addBot("bot3")
	g.addBot("bot4")
	g.Command("AddBot", &model.MsgAddBot{Name: "bot5"})
	g.Command("DeleteBot", &model.MsgDeleteBot{Name: "bot5"})
}

func (g *game) Close() {
	// use command to kick all players and users
	g.Command("OfflineAllObjects", &model.MsgOfflineAllObjects{})
	// check
	for {
		v, _ := g.Command("GetOnlineObjectsNumber", &model.MsgGetOnlineObjectNumber{})
		reply := v.(*model.MsgGetOnlineObjectNumberReply)
		playerNumber := reply.PlayerNumber
		userNumber := reply.UserNumber
		number := playerNumber + userNumber
		log.Printf("online objects [player]%d [user]%d\n", playerNumber, userNumber)
		if number == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	g.deleteAllBots()
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

type closeConnRequest struct {
	id                    int
	closeConnResponseChan chan struct{}
}

type actionRequest struct {
	id     int
	action string
	msg    any
}

func (g *game) PlayerConn(conn object.Conn) (int, error) {
	connReq := connRequest{
		Conn:             conn,
		connResponseChan: make(chan *connResponse),
	}
	g.playerConnRequestChan <- &connReq
	connResp := <-connReq.connResponseChan
	return connResp.id, connResp.err
}

func (g *game) PlayerCloseConn(id int) {
	closeConnReq := closeConnRequest{
		id:                    id,
		closeConnResponseChan: make(chan struct{}),
	}
	g.playerCloseConnRequestChan <- &closeConnReq
	<-closeConnReq.closeConnResponseChan
}

func (g *game) PlayerAction(id int, action string, msg any) {
	actionReq := actionRequest{
		id:     id,
		action: action,
		msg:    msg,
	}
	g.playerActionChan <- &actionReq
}

func (g *game) UserConn(conn object.Conn) (int, error) {
	connReq := connRequest{
		Conn:             conn,
		connResponseChan: make(chan *connResponse),
	}
	g.userConnRequestChan <- &connReq
	connResp := <-connReq.connResponseChan
	return connResp.id, connResp.err
}

func (g *game) UserCloseConn(id int) {
	closeConnReq := closeConnRequest{
		id:                    id,
		closeConnResponseChan: make(chan struct{}),
	}
	g.userCloseConnRequestChan <- &closeConnReq
	<-closeConnReq.closeConnResponseChan
}

func (g *game) UserAction(id int, action string, msg any) {
	actionReq := actionRequest{
		id:     id,
		action: action,
		msg:    msg,
	}
	g.userActionChan <- &actionReq
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

func (g *game) AddBot(msg *model.MsgAddBot) (any, error) {
	g.addBot(msg.Name)
	return nil, nil
}

func (g *game) DeleteBot(msg *model.MsgDeleteBot) (any, error) {
	g.deleteBot(msg.Name)
	return nil, nil
}

func (g *game) OfflineAllObjects(msg any) (any, error) {
	object.ObjectManager.OfflineAllObjects()
	return nil, nil
}

func (g *game) GetOnlineObjectsNumber(msg any) (*model.MsgGetOnlineObjectNumberReply, error) {
	return object.ObjectManager.GetOnlineObjectsNumber(), nil
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
