package object

import (
	"encoding/xml"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

func init() {
	// MonsterList was generated 2023-07-17 11:34:17 by https://xml-to-go.github.io/ in Ukraine.
	type MonsterList struct {
		XMLName  xml.Name         `xml:"MonsterList"`
		Text     string           `xml:",chardata"`
		Monsters []*MonsterConfig `xml:"Monster"`
	}
	var monsterList MonsterList
	conf.XML(conf.PathCommon, "Monsters/IGC_MonsterList.xml", &monsterList)
	MonsterTable = make(monsterTable)
	for _, monster := range monsterList.Monsters {
		MonsterTable[monster.Index] = monster
	}
}

var MonsterTable monsterTable

type monsterTable map[int]*MonsterConfig

type MonsterConfig struct {
	Text                   string `xml:",chardata"`
	Index                  int    `xml:"Index,attr"`
	ExpType                string `xml:"ExpType,attr"`
	Level                  int    `xml:"Level,attr"`
	HP                     int    `xml:"HP,attr"`
	MP                     int    `xml:"MP,attr"`
	DamageMin              int    `xml:"DamageMin,attr"`
	DamageMax              int    `xml:"DamageMax,attr"`
	Defense                int    `xml:"Defense,attr"`
	MagicDefense           int    `xml:"MagicDefense,attr"`
	AttackRate             int    `xml:"AttackRate,attr"`
	BlockRate              int    `xml:"BlockRate,attr"`
	MoveRange              int    `xml:"MoveRange,attr"`
	AttackType             int    `xml:"AttackType,attr"`
	AttackRange            int    `xml:"AttackRange,attr"`
	ViewRange              int    `xml:"ViewRange,attr"`
	MoveSpeed              int    `xml:"MoveSpeed,attr"`
	AttackSpeed            int    `xml:"AttackSpeed,attr"`
	RegenTime              int    `xml:"RegenTime,attr"`
	Attribute              int    `xml:"Attribute,attr"`
	ItemDropRate           int    `xml:"ItemDropRate,attr"`
	MoneyDropRate          int    `xml:"MoneyDropRate,attr"`
	MaxItemLevel           string `xml:"MaxItemLevel,attr"`
	MonsterSkill           string `xml:"MonsterSkill,attr"`
	IceRes                 int    `xml:"IceRes,attr"`
	PoisonRes              int    `xml:"PoisonRes,attr"`
	LightRes               int    `xml:"LightRes,attr"`
	FireRes                int    `xml:"FireRes,attr"`
	PentagramMainAttrib    int    `xml:"PentagramMainAttrib,attr"`
	PentagramAttribPattern int    `xml:"PentagramAttribPattern,attr"`
	PentagramDamageMin     int    `xml:"PentagramDamageMin,attr"`
	PentagramDamageMax     int    `xml:"PentagramDamageMax,attr"`
	PentagramAttackRate    int    `xml:"PentagramAttackRate,attr"`
	PentagramDefenseRate   int    `xml:"PentagramDefenseRate,attr"`
	PentagramDefense       int    `xml:"PentagramDefense,attr"`
	Name                   string `xml:"Name,attr"`
	Annotation             string `xml:"annotation,attr"`
}

func NewMonster(class int) *Monster {
	mc, ok := MonsterTable[class]
	if !ok {
		panic(fmt.Sprintf("monster invalid [class]%d", class))
	}
	monster := Monster{}
	monster.init()
	monster.ConnectState = ConnectStatePlaying
	monster.Live = true
	monster.State = 1
	switch class {
	case 240: // 仓库使者塞弗特
		monster.NpcType = NpcTypeWarehouse
	case 238, 368, 369, 370, 452, 453, 478, 450:
		monster.NpcType = NpcTypeChaosMix
	case 236:
		monster.NpcType = NpcTypeGoldarcher
	case 582:
		monster.NpcType = NpcTypePentagramMix
	default:
		monster.NpcType = NpcTypeNone
	}
	switch {
	case class >= 204 && class <= 259 ||
		class >= 367 && class <= 385 ||
		class >= 406 && class <= 408 ||
		class >= 414 && class <= 417 ||
		class >= 450 && class <= 453 ||
		class >= 464 && class <= 475 && class != 466 ||
		class == 478 || class == 479 ||
		class == 492 ||
		class == 522 ||
		class >= 540 && class <= 547 ||
		class >= 566 && class <= 568 ||
		class >= 577 && class <= 584 ||
		class == 603 || class == 604 ||
		class == 643 ||
		class == 651 ||
		class >= 658 && class <= 668 ||
		class >= 682 && class <= 688:
		monster.Type = ObjectTypeNPC
	default:
		monster.Type = ObjectTypeMonster
	}
	monster.Class = class
	monster.Name = mc.Name
	monster.Annotation = mc.Annotation
	monster.Level = mc.Level
	monster.attackDamageMin = mc.DamageMin
	monster.attackDamageMax = mc.DamageMax
	monster.attackRate = mc.AttackRate
	monster.attackSpeed = mc.AttackSpeed
	monster.defense = mc.Defense
	monster.magicDefense = mc.MagicDefense
	monster.defenseRate = mc.BlockRate
	monster.HP = mc.HP
	monster.MaxHP = mc.HP
	monster.MP = mc.MP
	monster.MaxMP = mc.MP
	monster.moveRange = mc.MoveRange
	monster.moveSpeed = mc.MoveSpeed
	monster.attackRange = mc.AttackRange
	monster.attackType = mc.AttackType
	monster.viewRange = mc.ViewRange
	monster.attribute = mc.Attribute
	monster.itemDropRate = mc.ItemDropRate
	monster.moneyDropRate = mc.MoneyDropRate
	monster.maxRegenTime = time.Duration(mc.RegenTime)
	monster.pentagramAttributePattern = mc.PentagramAttribPattern
	monster.pentagramAttackMin = mc.PentagramDamageMin
	monster.pentagramAttackMax = mc.PentagramDamageMax
	monster.pentagramAttackRate = mc.PentagramAttackRate
	monster.pentagramDefense = mc.PentagramDefense
	switch {
	case monster.attackType == 150:
		monster.addSkill(monster.attackType-100, 1)
	case monster.attackType >= 1:
		monster.addSkill(monster.attackType, 1)
	}
	switch class {
	case 161, 181, 189, 197, 267, 275: // 昆顿
		monster.addSkill(1, 1)   // 毒咒
		monster.addSkill(17, 1)  // 能量球
		monster.addSkill(55, 1)  // 玄月斩
		monster.addSkill(200, 1) // 召唤怪
		monster.addSkill(201, 1) // 免疫魔攻
		monster.addSkill(202, 1) // 免疫物攻
	case 149, 179, 187, 195, 265, 273, 335: // 暗黑巫师
		monster.addSkill(1, 1) // 毒咒
		// monster.addSkill(17, 1) // 能量球
	case 66, 73, 77: // 诅咒之王 蓝魔龙 天魔菲尼斯
		// 163, 165, 167, 171, 173, 427: // 赤色要塞
		monster.addSkill(17, 1) // 能量球
	case 89, 95, 112, 118, 124, 130, 143: // 骷灵巫师
		monster.addSkill(3, 1)  // 掌心雷
		monster.addSkill(17, 1) // 能量球
	case 433: // 骷髅法师
		monster.addSkill(3, 1) // 掌心雷
	case 561: // 美杜莎
		monster.addSkill(9, 1)   // 黑龙波
		monster.addSkill(38, 1)  // 单毒炎
		monster.addSkill(237, 1) // 闪电轰顶
		monster.addSkill(238, 1) // 黑暗之力
	case 673: // 辛维斯特
		// monster.addSkill(622, 1) // ?
	}
	return &monster
}

type NpcType int

const (
	NpcTypeNone = iota
	NpcTypeShop
	NpcTypeWarehouse
	NpcTypeChaosMix
	NpcTypeGoldarcher
	NpcTypePentagramMix
)

type actionState struct {
	rest         bool
	attack       bool
	move         bool
	escape       int
	emotion      int
	emotionCount int
}

type Monster struct {
	object
	NpcType         int
	moveRange       int // 移动范围
	spawnStartX     int
	spawnStartY     int
	spawnEndX       int
	spawnEndY       int
	spawnDir        int
	spawnDis        int
	actionState     actionState
	curActionTime   int64
	nextActionTime  int64
	delayActionTime int64
	// MTX             int
	// MTY             int
}

func (m *Monster) randPosition(number, x1, y1, x2, y2 int) (int, int) {
	w := x2 - x1
	if w <= 0 {
		w = 1
	}
	h := y2 - y1
	if h <= 0 {
		h = 1
	}
	if w == 1 && h == 1 {
		return x1, y1
	}
	for i := 0; i < 100; i++ {
		x := x1 + rand.Intn(w)
		y := y1 + rand.Intn(h)
		attr := maps.MapManager.GetMapAttr(number, x, y)
		if attr&1 == 0 && attr&4 == 0 && attr&8 == 0 {
			return x, y
		}
	}
	// panic(fmt.Sprintf("randPosition failed [number]%d", number))
	log.Printf("randPosition failed [map]%d [start](%d,%d) [end](%d,%d)\n", number, x1, y1, x2, y2)
	return x1, y1
}

func (m *Monster) spawnPosition() {
	// // wrong
	// if _map.Number == maps.Atlans && spawn.StartX == 251 && spawn.StartY == 51 ||
	// 	_map.Number == maps.Atlans && spawn.StartX == 7 && spawn.StartY == 52 ||
	// 	_map.Number == maps.LandOfTrial && spawn.StartX == 14 && spawn.StartY == 43 ||
	// 	_map.Number == maps.KanturuBoss && spawn.Index == 106 {
	// 	continue
	// }
	m.StartX, m.StartY = m.randPosition(m.MapNumber, m.spawnStartX, m.spawnStartY, m.spawnEndX, m.spawnEndY)
	maps.MapManager.SetMapAttrStand(m.MapNumber, m.StartX, m.StartY)
	m.X, m.Y = m.StartX, m.StartY
	m.Dir = m.spawnDir
	if m.Dir < 0 {
		m.Dir = rand.Intn(8)
	}
	m.createFrustrum()
}

func (m *Monster) overDis(tx, ty int) bool {
	if m.spawnDis < 1 {
		return false
	}
	x := tx - m.StartX
	y := ty - m.StartY
	dis := int(math.Sqrt(float64(x*x + y*y)))
	return dis < m.spawnDis
}

func (m *Monster) roamMove() {
	maxMoveRange := m.moveRange << 1
	m.nextActionTime = 1000
	x := 0
	y := 0
	cnt := 10
	for cnt > 0 {
		cnt--
		x = m.X + rand.Intn(maxMoveRange+1) - m.moveRange
		y = m.Y + rand.Intn(maxMoveRange+1) - m.moveRange
		if !m.overDis(x, y) {
			continue
		}
		attr := maps.MapManager.GetMapAttr(m.MapNumber, x, y)
		if ((m.Class == 249 || m.Class == 247) && attr&2 == 0) || // Guard
			attr&15 == 0 {
			m.TX = x
			m.TY = y
			m.actionState.move = true
			m.nextActionTime = 500
			return
		}
	}
}

func (m *Monster) searchEnemy() int {
	mindis := m.viewRange
	target := -1
	for _, vp := range m.viewports2 {
		tnum := vp.number
		if tnum < 0 {
			continue
		}
		om := m.objectManager
		tobj := om.object(om.objects[tnum])
		if (m.Class == 247 || m.Class == 249) &&
			tobj.PKLevel <= 4 {
			continue
		}
		attr := maps.MapManager.GetMapAttr(tobj.MapNumber, tobj.X, tobj.Y)
		if attr&1 == 0 {
			x := m.X - tobj.X
			y := m.Y - tobj.Y
			dis := int(math.Sqrt(float64(x*x + y*y)))
			vp.dis = dis
			if dis < mindis {
				mindis = dis
				target = tnum
			}
		}
	}
	return target
}

func (m *Monster) chaseMove(tobj *object) bool {
	mtx := tobj.X
	mty := tobj.Y
	tx := mtx
	ty := mty
	dis := 0
	if m.attackType >= 100 {
		dis = m.attackRange + 2
	} else {
		dis = m.attackRange
	}
	if m.X < mtx {
		tx -= dis
	} else {
		tx += dis
	}
	if m.Y < mty {
		ty -= dis
	} else {
		ty += dis
	}
	if maps.MapManager.CheckMapAttrStand(m.MapNumber, tx, ty) {
		dir := maps.CalcDir(tobj.X, tobj.Y, m.X, m.Y)
		cnt := len(maps.Dirs)
		for cnt > 0 {
			cnt--
			mtx = tobj.X + maps.Dirs[dir].X
			mty = tobj.Y + maps.Dirs[dir].Y
			attr := maps.MapManager.GetMapAttr(m.MapNumber, mtx, mty)
			if ((m.Class == 247 || m.Class == 249) && attr&2 == 0) ||
				attr&15 == 0 {
				m.TX = mtx
				m.TY = mty
				return true
			}
			if dir == len(maps.Dirs) {
				dir = 0
			}
		}
	}
	attr := maps.MapManager.GetMapAttr(m.MapNumber, tx, ty)
	if ((m.Class == 247 || m.Class == 249) && attr&2 == 0) ||
		attr&15 == 0 {
		m.TX = tx
		m.TY = ty
		return true
	}
	return false
}

func (m *Monster) baseAction() {
	if m.attribute == 0 {
		// attribute为0的怪物没有行为
		return
	}
	var tobj *object
	if m.targetNumber >= 0 {
		tnum := m.targetNumber
		om := m.objectManager
		tobj = om.object(om.objects[tnum])
	} else {
		m.actionState.emotion = 0
	}
	switch m.actionState.emotion {
	case 0: // 寻找目标
		// if m.actionState.attack {
		// 	m.actionState.attack = false
		// 	m.targetNumber = -1
		// 	m.nextActionTime = 500
		// }
		// rn := rand.Intn(2)
		// if rn == 0 {
		// 	m.actionState.rest = true
		// 	m.nextActionTime = 500
		// }
		m.targetNumber = m.searchEnemy()
		if m.targetNumber >= 0 {
			m.actionState.emotion = 1
			m.actionState.emotionCount = 30 // 30*500ms=15s
		} else if m.moveRange > 0 {
			m.roamMove()
		}
	case 1: // 移动及攻击
		if m.actionState.emotionCount > 0 {
			m.actionState.emotionCount--
		} else {
			m.actionState.emotion = 0
		}
		if m.targetNumber < 0 {
			return
		}
		dis := m.calcDistance(tobj)
		attackRange := m.attackRange
		if m.attackType >= 100 {
			attackRange = m.attackRange + 2
		}
		if dis <= attackRange {
			if maps.MapManager.CheckMapNoWall(m.MapNumber, m.X, m.Y, tobj.X, tobj.Y) {
				attr := maps.MapManager.GetMapAttr(m.MapNumber, tobj.X, tobj.Y)
				if attr&1 == 0 {
					m.actionState.attack = true
				} else {
					// 目标在安全区，傻看15s
					m.targetNumber = -1
					m.actionState.emotion = 1
					m.actionState.emotionCount = 30 // 30*500ms=15s
				}
				m.Dir = maps.CalcDir(tobj.X, tobj.Y, m.X, m.Y)
				m.nextActionTime = int64(m.attackSpeed)
			} else {
				// 隔着障碍物，傻看最多15s
			}
		} else {
			// 目标不在攻击范围
			if m.chaseMove(tobj) {
				// 可以寻路
				if maps.MapManager.CheckMapNoWall(m.MapNumber, m.X, m.Y, m.TX, m.TY) {
					m.actionState.move = true
					m.nextActionTime = 400
					m.Dir = maps.CalcDir(tobj.X, tobj.Y, m.X, m.Y)
				} else {
					// 隔着障碍物，进入状态3，原地傻等10s
					m.roamMove()
					m.actionState.emotion = 3
					m.actionState.emotionCount = 10
				}
			} else {
				// 不可以寻路
				m.roamMove()
			}
		}
	case 2:
		if m.actionState.emotionCount > 0 {
			m.actionState.emotionCount--
		} else {
			m.actionState.emotion = 0
		}
		m.actionState.move = false
		m.actionState.attack = false
		m.nextActionTime = 800
	case 3:
		if m.actionState.emotionCount > 0 {
			m.actionState.emotionCount--
		} else {
			m.actionState.emotion = 0
		}
		m.actionState.move = false
		m.actionState.attack = false
		m.nextActionTime = 400
	}
}

// 模拟怪物基本行为
func (m *Monster) process500ms() {
	// if m.Class != 249 {
	// 	return
	// }
	if m.ConnectState < ConnectStatePlaying ||
		!m.Live {
		return
	}
	curActionTime := time.Now().UnixMilli()
	if curActionTime-m.curActionTime+1 < m.nextActionTime+m.delayActionTime {
		return
	}
	m.curActionTime = curActionTime
	m.baseAction()
	if m.actionState.move {
		m.actionState.move = false
		// start := time.Now()
		path, ok := maps.MapManager.FindMapPath(m.MapNumber, m.X, m.Y, m.TX, m.TY)
		// fmt.Println("0500ms", time.Since(start).Microseconds())
		if !ok {
			return
		}
		dirs := make([]int, len(path))
		x, y := 0, 0
		for i := range path {
			if i == 0 {
				x, y = m.X, m.Y
			} else {
				x, y = path[i-1].X, path[i-1].Y
			}
			dirs[i] = maps.CalcDir(x, y, path[i].X, path[i].Y)
		}
		msg := model.MsgMove{Dirs: dirs, Path: path}
		m.Move(&msg)
		return
	}
	if m.actionState.attack {
		m.actionState.attack = false
		// 将普通攻击等效成一种技能
		// 技能数越多则普通攻击的概率越低
		n := len(m.skills)
		if rand.Intn(n+1) == 0 {
			msg := model.MsgAttack{
				Target: m.targetNumber,
			}
			m.Attack(&msg)
		} else {
			// 从map中随机获取一个元素
			cnt := rand.Intn(n) + 1
			skillNumber := 0
			for i := range m.skills {
				cnt--
				if cnt == 0 {
					skillNumber = i
					break
				}
			}
			msg := model.MsgSkillAttack{
				Target: m.targetNumber,
				Skill:  skillNumber,
			}
			m.SkillAttack(&msg)
		}
	}
}

func (m *Monster) processViewport() {
	m.destoryViewport()
	m.createViewport()
	if m.State == 1 {
		m.State = 2
	}
}

func (m *Monster) processRegen() {
	if !m.dieRegen {
		return
	}
	if m.ConnectState < ConnectStatePlaying {
		return
	}
	if time.Now().Unix()-int64(m.regenTime) < int64(m.maxRegenTime) {
		return
	}
	m.HP = m.MaxHP + m.AddHP
	m.MP = m.MaxMP + m.AddMP
	m.Live = true
	m.spawnPosition()
	m.dieRegen = false
	m.State = 1
}

func (m *Monster) process1000ms() {
	m.processViewport() // 1->2
	m.processRegen()    // 4->1
}
