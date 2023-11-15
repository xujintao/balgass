package game

import (
	"context"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/bot"
	"github.com/xujintao/balgass/src/server_game/game/cmd"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
	"github.com/xujintao/balgass/src/server_game/game/object"
)

func init() {
	Game.init()
}

var Game game

type game struct {
	serverRegisterChan         chan *model.MsgServerRegister
	playerConnRequestChan      chan *connRequest
	playerCloseConnRequestChan chan *closeConnRequest
	playerActionChan           chan *actionRequest
	userConnRequestChan        chan *connRequest
	userCloseConnRequestChan   chan *closeConnRequest
	userActionChan             chan *actionRequest
	commandRequestChan         chan *commandRequest
	cancel                     context.CancelFunc
}

func (g *game) init() {
	g.serverRegisterChan = make(chan *model.MsgServerRegister, 100)
	g.playerConnRequestChan = make(chan *connRequest, 100)
	g.playerCloseConnRequestChan = make(chan *closeConnRequest, 100)
	g.playerActionChan = make(chan *actionRequest, 1000)
	g.userConnRequestChan = make(chan *connRequest, 100)
	g.userCloseConnRequestChan = make(chan *closeConnRequest, 100)
	g.userActionChan = make(chan *actionRequest, 1000)
	g.commandRequestChan = make(chan *commandRequest, 100)
	bot.BotManager.Register(g)
}

func (g *game) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel
	go func() {
		addr := fmt.Sprintf("%s:%d",
			conf.Server.GameServerInfo.ConnectServerIP,
			conf.Server.GameServerInfo.ConnectServerPort)
		c, err := net.Dial("udp", addr)
		if err != nil {
			log.Fatalf("net.Dial failed [err]%v", err)
		}
		for {
			select {
			case serverRegister := <-g.serverRegisterChan:
				serverRegister.Code = conf.Server.GameServerInfo.Code
				nonPVP := 0
				if conf.Server.GameServerInfo.NonPVP {
					nonPVP = 1
				}
				serverRegister.Type_ = nonPVP
				data, err := serverRegister.Marshal()
				if err != nil {
					log.Printf("serverRegister.Marshal failed [err]%v\n", err)
				}
				frame := []byte{0xC1, 0x08, 0x01}
				frame = append(frame, data...)
				c.Write(frame)
			case <-ctx.Done():
				// todo
				return
			}
		}
	}()
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
				t := reflect.TypeOf(player)
				if _, ok := t.MethodByName(action); !ok {
					log.Printf("player has no [action]%s\n", action)
					break
				}
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
				out := reflect.ValueOf(&cmd.Command).MethodByName(name).Call(in)
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
					maps.MapManager.ExpireItem(time.Now())
					g.serverRegisterChan <- &model.MsgServerRegister{
						Percent: object.ObjectManager.GetPlayerPercent(),
					}
					// fmt.Println("1000ms", time.Since(start).Milliseconds())
				}
			case <-ctx.Done():
				// todo
				return
			}
		}
	}()
	// g.Command("AddBot", &model.MsgAddBot{Name: "bot1"})
	// g.Command("AddBot", &model.MsgAddBot{Name: "bot2"})
	// g.Command("AddBot", &model.MsgAddBot{Name: "bot3"})
	// g.Command("AddBot", &model.MsgAddBot{Name: "bot4"})
	// g.Command("AddBot", &model.MsgAddBot{Name: "bot5"})
	// g.Command("DeleteBot", &model.MsgDeleteBot{Name: "bot5"})
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
	bot.BotManager.DeleteAllBots()
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

func (g *game) SendWeather(number, weather int) {
	// if number == 0 {
	// 	log.Println(number, weather)
	// }
}
