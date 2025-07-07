package monster

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/shop"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

func SpawnMonster() {
	MonsterTable.init()
	DropManager.init()

	// MonsterSpawn was generated 2023-07-17 16:05:41 by https://xml-to-go.github.io/ in Ukraine.
	type MonsterSpawn struct {
		XMLName xml.Name `xml:"MonsterSpawn"`
		Text    string   `xml:",chardata"`
		Map     []*struct {
			Text       string `xml:",chardata"`
			Number     int    `xml:"Number,attr"`
			Name       string `xml:"Name,attr"`
			Annotation string `xml:"annotation,attr"`
			Spot       []*struct {
				Text        string `xml:",chardata"`
				Type        int    `xml:"Type,attr"`
				Description string `xml:"Description,attr"`
				Spawn       []*struct {
					Text     string `xml:",chardata"`
					Index    int    `xml:"Index,attr"`
					Distance int    `xml:"Distance,attr"`
					StartX   int    `xml:"StartX,attr"`
					StartY   int    `xml:"StartY,attr"`
					Dir      int    `xml:"Dir,attr"`
					EndX     int    `xml:"EndX,attr"`
					EndY     int    `xml:"EndY,attr"`
					Count    int    `xml:"Count,attr"`
					Element  int    `xml:"Element,attr"`
				} `xml:"Spawn"`
			} `xml:"Spot"`
		} `xml:"Map"`
	}
	var monsterSpawn MonsterSpawn
	conf.XML(conf.PathCommon, "Monsters/IGC_MonsterSpawn.xml", &monsterSpawn)
	for _, _map := range monsterSpawn.Map {
		for _, spot := range _map.Spot {
			for _, spawn := range spot.Spawn {
				cnt := spawn.Count
				if cnt == 0 {
					cnt = 1
				}
				for i := 0; i < cnt; i++ {
					spawnClass := spawn.Index
					spawnMapNumber := _map.Number
					spawnStartX := 0
					spawnStartY := 0
					spawnEndX := 0
					spawnEndY := 0
					switch spot.Type {
					case 0: // npc
						spawnStartX = spawn.StartX
						spawnStartY = spawn.StartY
						spawnEndX = spawn.StartX
						spawnEndY = spawn.StartY
					case 1, 3: // multiple
						spawnStartX = spawn.StartX
						spawnStartY = spawn.StartY
						spawnEndX = spawn.EndX
						spawnEndY = spawn.EndY
					case 2: // single
						spawnStartX = spawn.StartX - 3
						spawnStartY = spawn.StartY - 3
						spawnEndX = spawn.StartX + 3
						spawnEndY = spawn.StartY + 3
					}
					spawnDir := spawn.Dir
					spawnDis := spawn.Distance
					spawnElement := spawn.Element
					// register the new monster to object manager
					_, err := object.ObjectManager.AddMonster(func() *object.Object {
						return newMonster(
							spawnClass,
							spawnMapNumber,
							spawnStartX,
							spawnStartY,
							spawnEndX,
							spawnEndY,
							spawnDir,
							spawnDis,
							spawnElement,
						)
					})
					if err != nil {
						slog.Error("range monsterSpawn.Map object.ObjectManager.AddMonster", "err", err)
						os.Exit(1)
					}
				}
			}
		}
	}

	shop.ShopManager.ForEachShop(func(class, mapNumber, x, y, dir int) {
		obj, err := object.ObjectManager.AddMonster(func() *object.Object {
			return newMonster(
				class,
				mapNumber,
				x,
				y,
				x,
				y,
				dir,
				0,
				0,
			)
		})
		if err != nil {
			slog.Error("shop.ShopManager.ForEachShop object.ObjectManager.AddMonster", "err", err)
			os.Exit(1)
		}
		obj.NpcType = object.NpcTypeShop
	})
}

var MonsterTable monsterTable

type monsterTable map[int]*MonsterConfig

func (m *monsterTable) init() {
	// MonsterList was generated 2023-07-17 11:34:17 by https://xml-to-go.github.io/ in Ukraine.
	type MonsterList struct {
		XMLName  xml.Name         `xml:"MonsterList"`
		Text     string           `xml:",chardata"`
		Monsters []*MonsterConfig `xml:"Monster"`
	}
	var monsterList MonsterList
	conf.XML(conf.PathCommon, "Monsters/IGC_MonsterList.xml", &monsterList)
	tm := make(monsterTable)
	for _, monster := range monsterList.Monsters {
		tm[monster.Index] = monster
	}
	*m = tm
}

type MonsterConfig struct {
	Text                   string `xml:",chardata"`
	Index                  int    `xml:"Index,attr"`
	ExpType                int    `xml:"ExpType,attr"`
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
	MaxItemLevel           int    `xml:"MaxItemLevel,attr"`
	MonsterSkill           int    `xml:"MonsterSkill,attr"`
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

func newMonster(class, mapNumber, startX, startY, endX, endY, dir, dis, element int) *object.Object {
	mc, ok := MonsterTable[class]
	if !ok {
		panic(fmt.Sprintf("monster invalid [class]%d", class))
	}
	m := Monster{}
	m.Init()
	m.ConnectState = object.ConnectStatePlaying
	m.Live = true
	m.State = 1
	switch class {
	case 240: // 仓库使者塞弗特
		m.NpcType = object.NpcTypeWarehouse
	case 238, 368, 369, 370, 452, 453, 478, 450:
		m.NpcType = object.NpcTypeChaosMix
	case 236:
		m.NpcType = object.NpcTypeGoldarcher
	case 582:
		m.NpcType = object.NpcTypePentagramMix
	default:
		m.NpcType = object.NpcTypeNone
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
		m.Type = object.ObjectTypeNPC
	default:
		m.Type = object.ObjectTypeMonster
	}
	m.Class = class
	m.MapNumber = mapNumber
	m.spawnStartX = startX
	m.spawnStartY = startY
	m.spawnEndX = endX
	m.spawnEndY = endY
	m.spawnDir = dir
	m.SpawnPosition()
	m.spawnDis = dis
	m.PentagramMainAttribute = element
	m.Name = mc.Name
	m.Annotation = mc.Annotation
	m.Level = mc.Level
	m.AttackMin = mc.DamageMin
	m.AttackMax = mc.DamageMax
	m.AttackRate = mc.AttackRate
	m.AttackSpeed = mc.AttackSpeed
	m.Defense = mc.Defense
	m.MagicDefense = mc.MagicDefense
	m.DefenseRate = mc.BlockRate
	m.HP = mc.HP
	m.MaxHP = mc.HP
	m.MP = mc.MP
	m.MaxMP = mc.MP
	m.moveRange = mc.MoveRange
	m.MoveSpeed = mc.MoveSpeed
	m.AttackRange = mc.AttackRange
	m.AttackType = mc.AttackType
	m.ViewRange = mc.ViewRange
	m.Attribute = mc.Attribute
	m.ItemDropRate = mc.ItemDropRate
	m.MoneyDropRate = mc.MoneyDropRate
	m.MaxRegenTime = time.Duration(mc.RegenTime) * time.Second
	m.PentagramAttributePattern = mc.PentagramAttribPattern
	m.PentagramAttackMin = mc.PentagramDamageMin
	m.PentagramAttackMax = mc.PentagramDamageMax
	m.PentagramAttackRate = mc.PentagramAttackRate
	m.PentagramDefense = mc.PentagramDefense
	switch {
	case m.AttackType == 150:
		m.LearnSkill(m.AttackType - 100)
	case m.AttackType >= 1:
		m.LearnSkill(m.AttackType)
	}
	switch class {
	case 161, 181, 189, 197, 267, 275: // 昆顿
		m.LearnSkill(1)   // 毒咒
		m.LearnSkill(17)  // 能量球
		m.LearnSkill(55)  // 玄月斩
		m.LearnSkill(200) // 召唤怪
		m.LearnSkill(201) // 免疫魔攻
		m.LearnSkill(202) // 免疫物攻
	case 149, 179, 187, 195, 265, 273, 335: // 暗黑巫师
		m.LearnSkill(1) // 毒咒
		// m.LearnSkill(17) // 能量球
	case 66, 73, 77: // 诅咒之王 蓝魔龙 天魔菲尼斯
		// 163, 165, 167, 171, 173, 427: // 赤色要塞
		m.LearnSkill(17) // 能量球
	case 89, 95, 112, 118, 124, 130, 143: // 骷灵巫师
		m.LearnSkill(3)  // 掌心雷
		m.LearnSkill(17) // 能量球
	case 433: // 骷髅法师
		m.LearnSkill(3) // 掌心雷
	case 561: // 美杜莎
		m.LearnSkill(9)   // 黑龙波
		m.LearnSkill(38)  // 单毒炎
		m.LearnSkill(237) // 闪电轰顶
		m.LearnSkill(238) // 黑暗之力
	case 673: // 辛维斯特
		// m.LearnSkill(622) // ?
	}
	m.Objecter = &m
	return &m.Object
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
	object.Object
	moveRange           int // 移动范围
	spawnStartX         int
	spawnStartY         int
	spawnEndX           int
	spawnEndY           int
	spawnDir            int
	spawnDis            int
	actionState         actionState
	lastActionTime      time.Time
	nextActionInterval  int
	delayActionInterval int
}

func (m *Monster) Addr() string {
	return fmt.Sprintf("%d", m.Class)
}

func (m *Monster) Offline() {}

func (m *Monster) Push(msg any) {}

func (m *Monster) PushMPAG(int, int) {}

func (m *Monster) GetPKLevel() int {
	return 0
}

func (*Monster) GetMasterLevel() int {
	return 0
}

func (*Monster) IsMasterLevel() bool {
	return false
}

func (m *Monster) GetSkillMPAG(s *skill.Skill) (int, int) {
	return 0, 0
}

func (m *Monster) SpawnPosition() {
	// // wrong
	// if _map.Number == maps.Atlans && spawn.StartX == 251 && spawn.StartY == 51 ||
	// 	_map.Number == maps.Atlans && spawn.StartX == 7 && spawn.StartY == 52 ||
	// 	_map.Number == maps.LandOfTrial && spawn.StartX == 14 && spawn.StartY == 43 ||
	// 	_map.Number == maps.KanturuBoss && spawn.Index == 106 {
	// 	continue
	// }
	m.StartX, m.StartY = m.RandPosition(m.MapNumber, m.spawnStartX, m.spawnStartY, m.spawnEndX, m.spawnEndY)
	m.X, m.Y = m.StartX, m.StartY
	m.TX, m.TY = m.X, m.Y
	m.Dir = m.spawnDir
	if m.Dir < 0 {
		m.Dir = rand.Intn(8)
	}
	maps.MapManager.SetMapAttrStand(m.MapNumber, m.TX, m.TY)
	m.CreateFrustum()
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
	mindis := m.ViewRange
	target := -1
	for _, vp := range m.Viewports {
		tnum := vp.Number
		if tnum < 0 {
			continue
		}
		tobj := object.ObjectManager.GetObject(tnum)
		if tobj == nil {
			continue
		}
		if (m.Class == 247 || m.Class == 249) && tobj.GetPKLevel() <= 4 {
			continue
		}
		attr := maps.MapManager.GetMapAttr(tobj.MapNumber, tobj.X, tobj.Y)
		if attr&1 == 0 {
			x := m.X - tobj.X
			y := m.Y - tobj.Y
			dis := int(math.Sqrt(float64(x*x + y*y)))
			vp.Dis = dis
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

func (m *Monster) chaseMove(tobj *object.Object) bool {
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
	if m.Attribute == 0 {
		// attribute为0的怪物没有行为
		return
	}

	switch m.actionState.emotion {
	case 0: // 寻找目标
		// if m.actionState.attack {
		// 	m.actionState.attack = false
		// 	m.TargetNumber = -1
		// 	m.nextActionInterval = 500
		// }
		// rn := rand.Intn(2)
		// if rn == 0 {
		// 	m.actionState.rest = true
		// 	m.nextActionInterval = 500
		// }
		m.TargetNumber = m.searchEnemy()
		if m.TargetNumber >= 0 {
			m.actionState.emotion = 1
		} else if m.moveRange > 0 && !m.PathMoving {
			m.nextActionInterval = 500
			if m.roamMove() {
				m.actionState.move = true
			}
		}
	case 1: // 移动及攻击
		if m.TargetNumber < 0 {
			m.actionState.emotion = 3
			m.actionState.emotionCount = 10 // 10*500ms=5s
			return
		}
		if m.PathMoving {
			return
		}
		tobj := object.ObjectManager.GetObject(m.TargetNumber)
		if tobj == nil || !tobj.Live || tobj.MapNumber != m.MapNumber {
			m.TargetNumber = -1
			m.actionState.emotion = 3
			m.actionState.emotionCount = 2 // 2*500ms=1s
			return
		}
		dis := m.CalcDistance(tobj)
		attackRange := m.AttackRange
		if m.AttackType >= 100 {
			attackRange = m.AttackRange + 2
		}
		if dis > m.ViewRange<<1 {
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
						m.nextActionInterval = 400
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
					m.nextActionInterval = m.AttackSpeed
					return
				}
			}
		}
		// 进入状态3，原地傻等5s
		m.TargetNumber = -1
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
		m.nextActionInterval = 500
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
		m.nextActionInterval = 500
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
	msg := model.MsgMove{X: m.X, Y: m.Y, Dir: m.Dir, Dirs: dirs, Path: path}
	m.Move(&msg)
}

func (m *Monster) attack() {
	// 将普通攻击等效成一种技能
	// 技能数越多则普通攻击的概率越低
	n := len(m.Skills)
	if rand.Intn(n+1) == 0 {
		msg := model.MsgAttack{
			Target: m.TargetNumber,
			Action: 120,
			Dir:    m.Dir,
		}
		m.Attack(&msg)
	} else {
		// 从map中随机获取一个元素
		cnt := rand.Intn(n) + 1
		skillNumber := 0
		for i := range m.Skills {
			cnt--
			if cnt == 0 {
				skillNumber = int(i)
				break
			}
		}
		msg := model.MsgUseSkill{
			Target: m.TargetNumber,
			Skill:  skillNumber,
		}
		m.UseSkill(&msg)
	}
}

// 模拟怪物基本行为
func (m *Monster) ProcessAction() {
	if m.ConnectState < object.ConnectStatePlaying ||
		!m.Live {
		return
	}
	now := time.Now()
	interval := m.nextActionInterval + m.delayActionInterval
	if now.Before(m.lastActionTime.Add(time.Duration(interval) * time.Millisecond)) {
		return
	}
	m.lastActionTime = now
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

func (m *Monster) Regen() {
	m.HP = m.MaxHP
	m.MP = m.MaxMP
}

func (m *Monster) GetChangeUp() int {
	return 0
}
