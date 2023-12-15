package object

import (
	"encoding/xml"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/item"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
	"github.com/xujintao/balgass/src/server_game/game/skill"
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
	Text                   string           `xml:",chardata"`
	Index                  int              `xml:"Index,attr"`
	ExpType                string           `xml:"ExpType,attr"`
	Level                  int              `xml:"Level,attr"`
	HP                     int              `xml:"HP,attr"`
	MP                     int              `xml:"MP,attr"`
	DamageMin              int              `xml:"DamageMin,attr"`
	DamageMax              int              `xml:"DamageMax,attr"`
	Defense                int              `xml:"Defense,attr"`
	MagicDefense           int              `xml:"MagicDefense,attr"`
	AttackRate             int              `xml:"AttackRate,attr"`
	BlockRate              int              `xml:"BlockRate,attr"`
	MoveRange              int              `xml:"MoveRange,attr"`
	AttackType             skill.SkillIndex `xml:"AttackType,attr"`
	AttackRange            int              `xml:"AttackRange,attr"`
	ViewRange              int              `xml:"ViewRange,attr"`
	MoveSpeed              int              `xml:"MoveSpeed,attr"`
	AttackSpeed            int              `xml:"AttackSpeed,attr"`
	RegenTime              int              `xml:"RegenTime,attr"`
	Attribute              int              `xml:"Attribute,attr"`
	ItemDropRate           int              `xml:"ItemDropRate,attr"`
	MoneyDropRate          int              `xml:"MoneyDropRate,attr"`
	MaxItemLevel           string           `xml:"MaxItemLevel,attr"`
	MonsterSkill           string           `xml:"MonsterSkill,attr"`
	IceRes                 int              `xml:"IceRes,attr"`
	PoisonRes              int              `xml:"PoisonRes,attr"`
	LightRes               int              `xml:"LightRes,attr"`
	FireRes                int              `xml:"FireRes,attr"`
	PentagramMainAttrib    int              `xml:"PentagramMainAttrib,attr"`
	PentagramAttribPattern int              `xml:"PentagramAttribPattern,attr"`
	PentagramDamageMin     int              `xml:"PentagramDamageMin,attr"`
	PentagramDamageMax     int              `xml:"PentagramDamageMax,attr"`
	PentagramAttackRate    int              `xml:"PentagramAttackRate,attr"`
	PentagramDefenseRate   int              `xml:"PentagramDefenseRate,attr"`
	PentagramDefense       int              `xml:"PentagramDefense,attr"`
	Name                   string           `xml:"Name,attr"`
	Annotation             string           `xml:"annotation,attr"`
}

func NewMonster(class, mapNumber, startX, startY, endX, endY, dir, dis, element int) *Monster {
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
	monster.MapNumber = mapNumber
	monster.spawnStartX = startX
	monster.spawnStartY = startY
	monster.spawnEndX = endX
	monster.spawnEndY = endY
	monster.spawnDir = dir
	monster.spawnPosition()
	monster.spawnDis = dis
	monster.pentagramMainAttribute = element
	monster.Name = mc.Name
	monster.Annotation = mc.Annotation
	monster.Level = mc.Level
	monster.attackMin = mc.DamageMin
	monster.attackMax = mc.DamageMax
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
		monster.learnSkill(monster.attackType - 100)
	case monster.attackType >= 1:
		monster.learnSkill(monster.attackType)
	}
	switch class {
	case 161, 181, 189, 197, 267, 275: // 昆顿
		monster.learnSkill(1)   // 毒咒
		monster.learnSkill(17)  // 能量球
		monster.learnSkill(55)  // 玄月斩
		monster.learnSkill(200) // 召唤怪
		monster.learnSkill(201) // 免疫魔攻
		monster.learnSkill(202) // 免疫物攻
	case 149, 179, 187, 195, 265, 273, 335: // 暗黑巫师
		monster.learnSkill(1) // 毒咒
		// monster.learnSkill(17) // 能量球
	case 66, 73, 77: // 诅咒之王 蓝魔龙 天魔菲尼斯
		// 163, 165, 167, 171, 173, 427: // 赤色要塞
		monster.learnSkill(17) // 能量球
	case 89, 95, 112, 118, 124, 130, 143: // 骷灵巫师
		monster.learnSkill(3)  // 掌心雷
		monster.learnSkill(17) // 能量球
	case 433: // 骷髅法师
		monster.learnSkill(3) // 掌心雷
	case 561: // 美杜莎
		monster.learnSkill(9)   // 黑龙波
		monster.learnSkill(38)  // 单毒炎
		monster.learnSkill(237) // 闪电轰顶
		monster.learnSkill(238) // 黑暗之力
	case 673: // 辛维斯特
		// monster.learnSkill(622) // ?
	}
	return &monster
}

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
}

func (m *Monster) addr() string {
	return fmt.Sprintf("%d", m.Class)
}

func (m *Monster) Offline() {}

func (m *Monster) push(msg any) {}

func (m *Monster) PushMPAG(int, int) {}

func (m *Monster) getPKLevel() int {
	return 0
}

func (m *Monster) GetSkillMPAG(s *skill.Skill) (int, int) {
	return 0, 0
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
	m.TX, m.TY = m.X, m.Y
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
	return dis > m.spawnDis
}

func (m *Monster) searchEnemy() int {
	mindis := m.viewRange
	target := -1
	for _, vp := range m.viewports {
		tnum := vp.number
		if tnum < 0 {
			continue
		}
		tobj := ObjectManager.objects[tnum]
		if tobj == nil {
			continue
		}
		if (m.Class == 247 || m.Class == 249) && tobj.getPKLevel() <= 4 {
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

func (m *Monster) roamMove() bool {
	maxMoveRange := m.moveRange << 1
	cnt := 10
	for cnt > 0 {
		cnt--
		x := m.X + rand.Intn(maxMoveRange+1) - m.moveRange
		y := m.Y + rand.Intn(maxMoveRange+1) - m.moveRange
		if m.overDis(x, y) {
			continue
		}
		attr := maps.MapManager.GetMapAttr(m.MapNumber, x, y)
		if ((m.Class == 249 || m.Class == 247) && attr&2 == 0) || // Guard
			attr&15 == 0 {
			m.TX = x
			m.TY = y
			return true
		}
	}
	return false
}

func (m *Monster) chaseMove(tobj *object) bool {
	dir := maps.CalcDir(m.X, m.Y, tobj.X, tobj.Y)
	cnt := len(maps.Dirs)
	for cnt > 0 {
		tx := m.X + maps.Dirs[dir].X
		ty := m.Y + maps.Dirs[dir].Y
		attr := maps.MapManager.GetMapAttr(m.MapNumber, tx, ty)
		if ((m.Class == 247 || m.Class == 249) && attr&2 == 0) ||
			attr&15 == 0 {
			m.TX = tx
			m.TY = ty
			return true
		}
		cnt--
		dir++
		if dir == len(maps.Dirs) {
			dir = 0
		}
	}
	return false
}

func (m *Monster) baseAction() {
	if m.attribute == 0 {
		// attribute为0的怪物没有行为
		return
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
		} else if m.moveRange > 0 && !m.pathMoving {
			m.nextActionTime = 500
			if m.roamMove() {
				m.actionState.move = true
			}
		}
	case 1: // 移动及攻击
		if m.targetNumber < 0 {
			m.actionState.emotion = 3
			m.actionState.emotionCount = 10 // 10*500ms=5s
			return
		}
		if m.pathMoving {
			return
		}
		tobj := ObjectManager.objects[m.targetNumber]
		if tobj == nil || !tobj.Live || tobj.MapNumber != m.MapNumber {
			m.targetNumber = -1
			m.actionState.emotion = 3
			m.actionState.emotionCount = 2 // 2*500ms=1s
			return
		}
		dis := m.calcDistance(tobj)
		attackRange := m.attackRange
		if m.attackType >= 100 {
			attackRange = m.attackRange + 2
		}
		if dis > m.viewRange<<1 {
			// target is too far
		} else if dis > attackRange {
			// 目标不在攻击范围内
			if m.chaseMove(tobj) {
				// 可以寻路
				if maps.MapManager.CheckMapNoWall(m.MapNumber, m.X, m.Y, m.TX, m.TY) {
					attr := maps.MapManager.GetMapAttr(m.MapNumber, tobj.X, tobj.Y)
					if attr&1 == 0 {
						m.actionState.move = true
						m.Dir = maps.CalcDir(m.X, m.Y, tobj.X, tobj.Y)
						m.nextActionTime = 200
						return
					}
				}
			}
		} else {
			// 目标在攻击范围内
			if maps.MapManager.CheckMapNoWall(m.MapNumber, m.X, m.Y, tobj.X, tobj.Y) {
				attr := maps.MapManager.GetMapAttr(m.MapNumber, tobj.X, tobj.Y)
				if attr&1 == 0 {
					m.actionState.attack = true
					m.Dir = maps.CalcDir(m.X, m.Y, tobj.X, tobj.Y)
					m.nextActionTime = int64(m.attackSpeed)
					return
				}
			}
		}
		// 进入状态3，原地傻等5s
		m.targetNumber = -1
		m.actionState.emotion = 3
		m.actionState.emotionCount = 2 // 2*500ms=1s
	case 2:
		if m.actionState.emotionCount > 0 {
			m.actionState.emotionCount--
		} else {
			m.actionState.emotion = 0
		}
		m.actionState.move = false
		m.actionState.attack = false
		m.nextActionTime = 500
	case 3:
		if m.actionState.emotionCount > 0 {
			m.actionState.emotionCount--
		} else {
			if m.overDis(m.X, m.Y) {
				m.StartX = m.X
				m.StartY = m.Y
			}
			m.actionState.emotion = 0
		}
		m.actionState.move = false
		m.actionState.attack = false
		m.nextActionTime = 500
	}
}

func (m *Monster) move() {
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
}

func (m *Monster) attack() {
	// 将普通攻击等效成一种技能
	// 技能数越多则普通攻击的概率越低
	n := len(m.skills)
	if rand.Intn(n+1) == 0 {
		msg := model.MsgAttack{
			Target: m.targetNumber,
			Action: 120,
			Dir:    m.Dir,
		}
		m.Attack(&msg)
	} else {
		// 从map中随机获取一个元素
		cnt := rand.Intn(n) + 1
		skillNumber := 0
		for i := range m.skills {
			cnt--
			if cnt == 0 {
				skillNumber = int(i)
				break
			}
		}
		msg := model.MsgUseSkill{
			Target: m.targetNumber,
			Skill:  skillNumber,
		}
		m.UseSkill(&msg)
	}
}

// 模拟怪物基本行为
func (m *Monster) ProcessAction() {
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
		m.move()
		return
	}
	if m.actionState.attack {
		m.actionState.attack = false
		m.attack()
		return
	}
}

func (m *Monster) Process1000ms() {

}

func (m *Monster) Die(obj *object) {

}

func (m *Monster) Regen() {
	m.HP = m.MaxHP + m.AddHP
	m.MP = m.MaxMP + m.AddMP
}

func (m *Monster) GetChangeUp() int {
	return 0
}

func (m *Monster) GetInventory() [9]*item.Item {
	return [9]*item.Item{}
}
