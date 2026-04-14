package handle

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"

	"github.com/gorilla/websocket"
	"github.com/xujintao/balgass/src/server-game/game"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

func init() {
	WSHandle.init()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("upgrader.Upgrade", "err", err)
		return
	}
	addr := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		addr = realIP
	}
	id, err := WSHandle.OnConn(c, addr)
	if err != nil {
		c.Close()
		return
	}
	defer WSHandle.OnClose(id)
	for {
		var req map[string]any
		err := c.ReadJSON(&req)
		if err != nil {
			slog.Error("conn.ReadJSON", "err", err, "addr", addr)
			return
		}
		actionField, ok := req["action"]
		if !ok {
			s := "websocket request has no action field"
			slog.Error(s, "addr", addr)
			continue
		}
		action, ok := actionField.(string)
		if !ok {
			s := "websocket request action field is not a string"
			slog.Error(s, "addr", addr)
			continue
		}
		inField, ok := req["in"]
		if !ok {
			s := fmt.Sprintf("websocket request [action]%s has no in field", action)
			slog.Error(s, "addr", addr)
			continue
		}
		data, err := json.Marshal(inField)
		if err != nil {
			s := fmt.Sprintf("websocket request [action]%s in field is not a valid json", action)
			slog.Error(s, "err", err, "addr", addr)
			continue
		}
		WSHandle.Handle(id, action, data)
	}
}

var WSHandle wsHandle

type wsHandle struct {
	wsIns  map[string]*wsIn
	wsOuts map[any]*wsOut
}

func (h *wsHandle) init() {
	// ingress
	h.wsIns = make(map[string]*wsIn)
	for _, v := range wsIns {
		h.wsIns[v.action] = v
	}
	// egress
	h.wsOuts = make(map[any]*wsOut)
	for _, v := range wsOuts {
		t := reflect.TypeOf(v.msg)
		if t.Kind() != reflect.Ptr {
			slog.Error("wsOut msg field must be a pointer", "action", v.action)
			os.Exit(1)
		}
		h.wsOuts[t] = v
	}
}

type wsconn struct {
	*websocket.Conn
	addr string
}

func (c *wsconn) Addr() string {
	return c.addr
}

func (c *wsconn) Write(msg any) error {
	reply, err := WSHandle.marshal(msg)
	if err != nil {
		slog.Error("(*wsconn).Write WSHandle.marshal", "err", err)
		return err
	}
	err = c.WriteJSON(reply)
	if err != nil {
		slog.Error("(*wsconn).Write c.WriteJSON", "err", err)
		return err
	}
	return nil
}

func (h *wsHandle) OnConn(c *websocket.Conn, addr string) (int, error) {
	conn := wsconn{Conn: c, addr: addr}
	return game.Game.UserConn(&conn)
}

func (h *wsHandle) OnClose(id int) {
	game.Game.UserCloseConn(id)
}

func (h *wsHandle) Handle(id int, action string, data []byte) {
	var api *wsIn
	var ok bool
	if api, ok = h.wsIns[action]; !ok {
		slog.Error("invalid websocket request", "action", action)
		return
	}
	msg := reflect.New(reflect.TypeOf(api.msg).Elem()).Interface()
	if err := json.Unmarshal(data, msg); err != nil {
		slog.Error("json.Unmarshal", "err", err, "action", action)
		return
	}
	game.Game.UserAction(id, api.action, msg)
}

func (h *wsHandle) marshal(msg any) (map[string]any, error) {
	v := reflect.ValueOf(msg)
	t := v.Type()
	api, ok := h.wsOuts[t]
	if !ok {
		err := fmt.Errorf("%s has not yet be registered to wsOuts table", t.String())
		return nil, err
	}
	reply := map[string]any{
		"action": api.action,
		"out":    msg,
	}
	return reply, nil
}

type wsIn struct {
	action string
	msg    any
}
type wsOut struct {
	action string
	msg    any
}

var wsIns = [...]*wsIn{
	{"SubscribeMap", (*model.MsgSubscribeMap)(nil)},
	// {"SubscribePlayer", nil},
}

var wsOuts = [...]*wsOut{
	{"HandleErrorReply", (*model.MsgHandleErrorReply)(nil)},
	{"SubscribeMapReply", (*model.MsgSubscribeMapReply)(nil)},
}
