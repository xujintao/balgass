package object

import (
	"context"
	"log/slog"

	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

func NewUser(conn Conn) *user {
	ctx, cancel := context.WithCancel(context.Background())
	u := &user{
		conn:    conn,
		msgChan: make(chan any, 100),
		cancel:  cancel,
	}
	go func() {
		for {
			select {
			case msg := <-u.msgChan:
				u.conn.Write(msg)
			case <-ctx.Done():
				close(u.msgChan)
				u.conn.Close()
				return // return ctx.Err()
			}
		}
	}()
	return u
}

type user struct {
	index   int
	offline bool
	conn    Conn
	msgChan chan any
	cancel  context.CancelFunc
}

func (u *user) Offline() {
	if u.offline {
		return
	}
	u.offline = true
	// todo
	u.unsubscribeMap()
	u.cancel()
}

func (u *user) Push(msg any) {
	if u.offline {
		slog.Warn("Still pushing msg to offline user", "msg", msg, "user", u.index)
		return
	}
	if len(u.msgChan) > 80 {
		u.Offline()
		return
	}
	u.msgChan <- msg
}

func (u *user) Test(msg *model.MsgTest) {
	u.Push(msg)
}

func (u *user) SubscribeMap(msg *model.MsgSubscribeMap) {
	number := msg.Number
	if number < 0 || number >= maps.MAX_MAP_NUMBER {
		resp := model.MsgSubscribeMapReply{
			Err: "invalid map number",
		}
		u.Push(&resp)
		return
	}
	// remove
	table := ObjectManager.mapSubscribeTable
	for n, users := range table {
		if _, ok := users[u]; ok {
			if n == number {
				return
			}
			delete(users, u)
			if len(users) == 0 {
				delete(table, n)
			}
		}
	}
	// insert
	users, ok := table[number]
	if !ok {
		users = make(map[*user]struct{})
		table[number] = users
	}
	users[u] = struct{}{}

	// reply map data
	reply := model.MsgSubscribeMapReply{
		Name: "map",
		Data: maps.MapManager.GetMapPots(number),
	}
	u.Push(&reply)
}

func (u *user) unsubscribeMap() {
	table := ObjectManager.mapSubscribeTable
	for n, users := range table {
		if _, ok := users[u]; ok {
			delete(users, u)
			if len(users) == 0 {
				delete(table, n)
			}
		}
	}
}

func (u *user) publishMap(objects any) {
	msg := model.MsgSubscribeMapReply{
		Name: "object",
		Data: objects,
	}
	u.Push(&msg)
}
