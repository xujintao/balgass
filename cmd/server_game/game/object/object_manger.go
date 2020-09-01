package object

import (
	"fmt"
	"sync"
	"time"

	"github.com/xujintao/balgass/cmd/server_game/conf"
	"github.com/xujintao/balgass/network"
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

func ObjectAdd(addr string, conn network.ConnWriter) (int, error) {
	return objectManagerDefault.objectAdd(addr, conn)
}

func ObjectFind(id int) *Object {
	return objectManagerDefault.objectFind(id)
}

func ObjectDelete(id int) {

}

var poolObject = sync.Pool{
	New: func() interface{} {
		return &Object{}
	},
}

var objectManagerDefault objectManager

type objectManager struct {
	mu      sync.Mutex
	objects map[int]*Object
}

// ObjectAdd search a avaliable object from pool
// and return the object index
func (m *objectManager) objectAdd(addr string, conn network.ConnWriter) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

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
	return 0, nil
}

func (m *objectManager) objectFind(id int) *Object {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.objects[id]
}

func (m *objectManager) objectDelete(id int) {
	m.mu.Lock()
	obj := m.objectFind(id)
	delete(m.objects, id)
	m.mu.Unlock()
	poolObject.Put(obj)
}
