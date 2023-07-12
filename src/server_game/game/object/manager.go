package object

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

func init() {
	ObjectManager.init()
}

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
	m.maxObjectCount = conf.Server.GameServerInfo.MaxMonsterCount + conf.Server.GameServerInfo.MaxSummonMonsterCount + conf.Server.GameServerInfo.MaxPlayerCount
	m.objects = make([]object, m.maxObjectCount)
	// objectBills = make([]bill, conf.Server.MaxPlayerCount)
	// 先有怪后有玩家
	m.playerStartIndex = conf.Server.GameServerInfo.MaxMonsterCount + conf.Server.GameServerInfo.MaxSummonMonsterCount
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

func (m *objectManager) AddPlayer(conn Conn) (int, error) {
	// limit max player count
	if m.playerCount >= conf.Server.GameServerInfo.MaxPlayerCount {
		// reply
		msg := model.MsgConnectFailed{Result: 4}
		conn.Write(&msg)
		return -1, fmt.Errorf("over max player count")
	}

	// get unified object index
	index := m.lastPlayerIndex
	cnt := conf.Server.GameServerInfo.MaxPlayerCount
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
	player.conn = conn
	player.msgChan = make(chan any, 100)
	ctx, cancel := context.WithCancel(context.Background())
	player.cancel = cancel
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
				player.conn.Write(msg)
			case <-ctx.Done():
				return // return ctx.Err()
			}
		}
	}()

	// register the new player to object manager
	m.objects[index] = player

	// reply
	msg := model.MsgConnectSuccess{
		Result:  1,
		ID:      index,
		Version: conf.MapServers.ServerInfo.Version,
	}
	player.Push(&msg)
	log.Printf("player online [id]%d [addr]%s", player.index, player.conn.Addr())
	return index, nil
}

func (m *objectManager) DeletePlayer(id int) {
	player := m.objects[id].(*Player)
	player.cancel()

	log.Printf("player offline [id]%d [addr]%s", player.index, player.conn.Addr())
	poolPlayer.Put(player)

	// unregister player from object manager
	m.objects[player.index] = nil
	m.playerCount--
}

func (m *objectManager) GetPlayer(id int) *Player {
	return m.objects[id].(*Player)
}
