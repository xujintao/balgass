package object

import (
	"fmt"
	"log"

	"github.com/xujintao/balgass/src/server-connect/service/model"
)

func init() {
	ObjectManager.init()
}

type Conn interface {
	Addr() string
	Write(any) error
	Close() error
}

var ObjectManager objectManager

type objectManager struct {
	maxPlayerCount   int
	playerStartIndex int
	lastPlayerIndex  int
	playerCount      int
	players          []*player
}

func (m *objectManager) init() {
	m.maxPlayerCount = 1000
	m.playerStartIndex = 0
	m.lastPlayerIndex = m.playerStartIndex - 1
	m.players = make([]*player, m.maxPlayerCount)
}

func (m *objectManager) AddPlayer(conn Conn) (int, error) {
	if m.playerCount >= m.maxPlayerCount {
		return -1, fmt.Errorf("over max player count")
	}
	index := m.lastPlayerIndex
	cnt := m.maxPlayerCount
	for cnt > 0 {
		cnt--
		index++
		if index >= m.maxPlayerCount {
			index = m.playerStartIndex
		}
		if m.players[index] == nil {
			break
		}
	}
	if cnt == 0 {
		panic(fmt.Errorf("have no free player index"))
	}
	m.lastPlayerIndex = index
	m.playerCount++
	p := NewPlayer(conn)
	p.objectManager = m
	p.index = index
	m.players[index] = p

	// reply
	msg := model.MsgConnectReply{
		Result: 1,
	}
	p.Push(&msg)
	log.Printf("player online [id]%d [addr]%s\n", p.index, p.conn.Addr())
	return index, nil
}

func (m *objectManager) GetPlayer(id int) *player {
	return m.players[id]
}

func (m *objectManager) DeletePlayer(id int) {
	p := m.players[id]
	if p == nil {
		return
	}
	p.Offline()
	log.Printf("player offline [id]%d [addr]%s\n", p.index, p.conn.Addr())

	// unregister player from object manager
	m.players[id] = nil
	m.playerCount--
}
