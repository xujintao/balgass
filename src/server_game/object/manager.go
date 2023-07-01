package object

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_game/conf"
)

var ObjectManager objectManager

type objectManager struct {
	maxObjectCount   int
	playerStartIndex int
	lastPlayerIndex  int
	objects          []object

	monsterCount       int
	summonMonsterCount int
	playerCount        int
}

type object any

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

func (m *objectManager) AddPlayer(ctx context.Context, cr *c1c2.ConnRequest, marshaller MsgMarshaller) (int, error) {
	// limit max player count
	if m.playerCount >= conf.Server.MaxPlayerCount {
		// reply
		// res := &network.Response{}
		// body := []byte{0x04}
		// res.WriteHead2(0xC1, 0xF1, 0x01).Write(body)
		// conn.Write(res)
		return -1, fmt.Errorf("over max player count")
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
	player.closeConn = cr.CloseConn
	player.writeConn = cr.WriteConn
	player.msgMarshaller = marshaller
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
			case msg := <-player.msgChan:
				player.Write(msg)
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
	return index, nil
}

func (m *objectManager) DeletePlayer(ctx context.Context) {
	player := ctx.Value(nil).(*Player)
	poolPlayer.Put(player)

	// unregister player from object manager
	m.objects[player.index] = nil
	m.playerCount--
}
