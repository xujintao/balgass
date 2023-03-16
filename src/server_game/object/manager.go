package object

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/xujintao/balgass/src/server_game/conf"
)

type object interface{}

// AddMonster add a monster
func AddMonster(class int) {
	// objectManagerDefault.monsterAdd(class)
}

// DeleteMonster delete a monster
func DeleteMonster(id int) {
	// objectManagerDefault.monsterDelete(id)
}

type pusher interface {
	Push(context.Context, interface{})
}

// AddPlayer add a player
func AddPlayer(ctx context.Context, pusher pusher) {
	objectManagerDefault.AddPlayer(ctx, pusher)
}

// DeletePlayer delete a player
func DeletePlayer(ctx context.Context, force bool) {
	objectManagerDefault.DeletePlayer(ctx)
	if force {
		// close connection
	}
}

var objectManagerDefault objectManager

type objectManager struct {
	maxObjectCount   int
	playerStartIndex int
	lastPlayerIndex  int
	objects          []object

	monsterCount       int
	summonMonsterCount int
	playerCount        int
}

func (m *objectManager) init() {
	m.maxObjectCount = conf.Server.MaxMonsterCount + conf.Server.MaxSummonMonsterCount + conf.Server.MaxPlayerCount
	m.objects = make([]object, m.maxObjectCount)
	// objectBills = make([]bill, conf.Server.MaxPlayerCount)
	// 先有怪后有玩家
	m.playerStartIndex = conf.Server.MaxMonsterCount + conf.Server.MaxSummonMonsterCount
	m.lastPlayerIndex = m.playerStartIndex
}

func (m *objectManager) objectMaxRange(index int) bool {
	if index < 0 || index >= m.maxObjectCount {
		return false
	}
	return true
}

var poolPlayer = sync.Pool{
	New: func() interface{} {
		return &Player{}
	},
}

func (m *objectManager) AddPlayer(ctx context.Context, pusher pusher) {
	// limit max player count
	if m.playerCount >= conf.Server.MaxPlayerCount {
		// reply
		// res := &network.Response{}
		// body := []byte{0x04}
		// res.WriteHead2(0xC1, 0xF1, 0x01).Write(body)
		// conn.Write(res)
		return
	}

	// get unified object index
	index := m.lastPlayerIndex
	cnt := conf.Server.MaxPlayerCount
	for cnt > 0 {
		if m.objects[index] == nil {
			break
		}
		index++
		if index >= m.maxObjectCount {
			index = m.playerStartIndex
		}
		cnt--
	}
	if cnt == 0 {
		panic(fmt.Errorf("have no free index"))
	}
	m.lastPlayerIndex = index
	m.playerCount++

	// create a new player
	player := poolPlayer.Get().(*Player)
	player.LoginMsgSend = false
	player.LoginMsgCount = 0
	player.index = index
	// player.conn = conn
	player.ConnectCheckTime = time.Now()
	player.AutoSaveTime = player.ConnectCheckTime
	player.Connected = PlayerConnected
	player.CheckSpeedHack = false
	player.EnableCharacterCreate = false
	player.Type = ObjectUser
	// player.Addr = addr
	// player.Conn = conn
	// player.pusher = pusher.(Pusher)

	// new a new goroutine to reply message
	go func() {
		for {
			select {
			case msg := <-player.pushChan:
				pusher.Push(ctx, msg)
			case <-ctx.Done():
				return // return ctx.Err()
			}
		}
	}()

	// register the new player to object manager
	m.objects[index] = player

	// reply
	// msg := model.MsgConnectResult{}
	// ctx, err = game.OnConn(addr, conn, h)
	// if err != nil {
	// 	msg.Result = 0
	// } else {
	// 	msg.Result = ctx.(int)
	// }
	// h.Push(conn, &msg)
	// return
}

func (m *objectManager) DeletePlayer(ctx context.Context) {
	player := ctx.Value(nil).(*Player)
	poolPlayer.Put(player)

	// unregister player from object manager
	m.objects[player.index] = nil
	m.playerCount--
}
