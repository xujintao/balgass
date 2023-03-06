package object

import (
	"sync"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_game/conf"
)

var (
	maxObjectCount            int
	objectUserCountStartIndex int
	objects                   []Object
	objectCount               int
	objectUserCount           int
	objectMonsterCount        int
	objectSummonMonsterCount  int
)

func init() {
	maxObjectCount = conf.Server.MaxObjectMonsterCount + conf.Server.MaxObjectSummonMonsterCount + conf.Server.MaxObjectUserCount
	// objects = make([]Object, maxObjectCount)
	// objectBills = make([]bill, conf.Server.MaxObjectUserCount)
	// 先有怪后有玩家
	objectUserCountStartIndex = maxObjectCount - conf.Server.MaxObjectUserCount
	objectCount = objectUserCountStartIndex

}

func objectMaxRange(index int) bool {
	if index < 0 || index >= maxObjectCount {
		return false
	}
	return true
}

// Find find a object from object-manager
func Find(id int) interface{} {
	return objectManagerDefault.find(id)
}

// MonsterAdd add a monster
func MonsterAdd(class int) {
	// return objectManagerDefault.monsterAdd(class)
}

// PlayerAdd add a player
func PlayerAdd(addr string, conn c1c2.ConnWriter, pusher interface{}) (int, error) {
	return objectManagerDefault.playerAdd(addr, conn, pusher)
}

// MonsterDelete delete a monster
func MonsterDelete(id int) {
	// objectManagerDefault.monsterDelete(id)
}

// PlayerDelete delete a player
func PlayerDelete(id int) {
	objectManagerDefault.playerDelete(id)
}

var poolPlayer = sync.Pool{
	New: func() interface{} {
		return &Player{}
	},
}

var objectManagerDefault objectManager

type objectManager struct {
	mu      sync.Mutex
	objects map[int]interface{}
}

func (m *objectManager) find(id int) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.objects[id]
}

func (m *objectManager) playerAdd(addr string, conn c1c2.ConnWriter, pusher interface{}) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 判断当前玩家数
	player := poolPlayer.Get().(*Player)
	player.Addr = addr
	player.Conn = conn
	player.pusher = pusher.(Pusher)
	/*
		if objectUserCount > conf.Server.MaxObjectUserCount {
			// 响应
			res := &network.Response{}
			body := []byte{0x04}
			res.WriteHead2(0xC1, 0xF1, 0x01).Write(body)
			conn.Write(res)
			return -1, fmt.Errorf("current user number: [%d], over maximum number of users: [%d]", objectUserCount, conf.Server.MaxObjectUserCount)
		}

		index := objectCount
		cnt := conf.Server.MaxObjectUserCount
		for cnt > 0 {
			if objects[index].Connected == PlayerEmpty {
				break
			}
			index++
			if index >= maxObjectCount {
				index = objectUserCountStartIndex
			}
			cnt--
		}
		if cnt == 0 {
			return 0, fmt.Errorf("have no free index")
		}

		o := &objects[index]
		o.Reset()
		o.LoginMsgSend = false
		o.LoginMsgCount = 0
		o.index = index
		o.conn = conn
		o.ConnectCheckTime = time.Now()
		o.AutoSaveTime = o.ConnectCheckTime
		o.Connected = PlayerConnected
		o.CheckSpeedHack = false
		o.EnableCharacterCreate = false
		o.Type = ObjectUser
	*/
	return 0, nil
}

func (m *objectManager) playerDelete(id int) {
	m.mu.Lock()
	obj := m.find(id)
	delete(m.objects, id)
	m.mu.Unlock()
	poolPlayer.Put(obj)
}
