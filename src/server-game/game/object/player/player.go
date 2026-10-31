package player

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log/slog"
	"math"
	"os"
	"sort"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/effect"
	"github.com/xujintao/balgass/src/server-game/game/exp"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/shop"
	"github.com/xujintao/balgass/src/server-game/game/skill"
	"gorm.io/gorm"
)

func init() {
	CharacterTable.init()
}

var CharacterTable characterTable

const autoSaveInterval = 10 * time.Minute

type characterTable map[int]model.Character

func (table *characterTable) init() {
	type ClassTable struct {
		XMLName xml.Name `xml:"ClassTable"`
		Class   []struct {
			Annotation string  `xml:"Annotation,attr"`
			Class      int     `xml:"Class,attr"`
			Level      int     `xml:"Level,attr"`
			MapNumber  int     `xml:"MapNumber,attr"`
			Strength   int     `xml:"Strength,attr"`
			Dexterity  int     `xml:"Dexterity,attr"`
			Vitality   int     `xml:"Vitality,attr"`
			Energy     int     `xml:"Energy,attr"`
			Leadership int     `xml:"Leadership,attr"`
			HP         int     `xml:"HP,attr"`
			MP         int     `xml:"MP,attr"`
			LevelHP    float32 `xml:"LevelHP,attr"`
			LevelMP    float32 `xml:"LevelMP,attr"`
			VitalityHP float32 `xml:"VitalityHP,attr"`
			EnergyMP   float32 `xml:"EnergyMP,attr"`
			Skills     struct {
				Skill []struct {
					Index int `xml:"Index,attr" json:"index"`
				} `xml:"Skill"`
			} `xml:"Skills"`
			Inventories struct {
				Inventory []struct {
					Position   int `xml:"Position,attr" json:"position"`
					Section    int `xml:"Section,attr" json:"section"`
					Index      int `xml:"Index,attr" json:"index"`
					Durability int `xml:"Durability,attr" json:"durability"`
				} `xml:"Inventory"`
			} `xml:"Inventories"`
		} `xml:"Class"`
	}
	var ct ClassTable
	conf.XML(conf.Server.PathCommon, "class/class.xml", &ct)
	t := make(characterTable)
	for _, v := range ct.Class {
		c := model.Character{
			Class:      v.Class,
			Level:      v.Level,
			MapNumber:  v.MapNumber,
			Strength:   v.Strength,
			Dexterity:  v.Dexterity,
			Vitality:   v.Vitality,
			Energy:     v.Energy,
			Leadership: v.Leadership,
			HP:         v.HP,
			MP:         v.MP,
			LevelHP:    v.LevelHP,
			LevelMP:    v.LevelMP,
			VitalityHP: v.VitalityHP,
			EnergyMP:   v.EnergyMP,
		}
		skills, err := json.Marshal(v.Skills.Skill)
		if err != nil {
			slog.Error("marshal character class skills", "class", v.Class, "err", err)
			os.Exit(1)
		}
		if err := c.Skills.UnmarshalJSON(skills); err != nil {
			slog.Error("unmarshal character class skills", "class", v.Class, "err", err)
			os.Exit(1)
		}
		inventory, err := json.Marshal(v.Inventories.Inventory)
		if err != nil {
			slog.Error("marshal character class inventory", "class", v.Class, "err", err)
			os.Exit(1)
		}
		if err := c.Inventory.UnmarshalJSON(inventory); err != nil {
			slog.Error("unmarshal character class inventory", "class", v.Class, "err", err)
			os.Exit(1)
		}
		t[c.Class] = c
	}
	*table = t
}

func NewPlayer(conn object.Conn, actioner object.Actioner) (int, error) {
	// register the new player to object manager
	return object.ObjectManager.AddPlayer(conn, func() *object.Object {
		return newPlayer(conn, actioner)
	})
}

func newPlayer(conn object.Conn, actioner object.Actioner) *object.Object {
	// create a new player
	p := Player{}
	p.Init()
	// p.LoginMsgSend = false
	// p.LoginMsgCount = 0
	p.conn = conn
	p.msgChan = make(chan any, 100)
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	// player.ConnectCheckTime = time.Now()
	// player.AutoSaveTime = player.ConnectCheckTime
	p.ConnectState = object.ConnectStateConnected
	// p.CheckSpeedHack = false
	// p.EnableCharacterCreate = false
	p.Type = object.ObjectTypePlayer
	p.actioner = actioner

	// new a new goroutine to reply message
	go func() {
		for {
			select {
			case msg := <-p.msgChan:
				err := p.conn.Write(msg)
				if err != nil {
					slog.Error("p.conn.Write",
						"err", err, "player", p.Name, "msg", msg)
				}
			case <-ctx.Done():
				close(p.msgChan)
				p.conn.Close()
				return // return ctx.Err()
			}
		}
	}()
	p.Objecter = &p
	return &p.Object
}

type Player struct {
	object.Object
	offline                        bool
	conn                           object.Conn
	msgChan                        chan any
	cancel                         context.CancelFunc
	actioner                       object.Actioner
	accountID                      int
	accountName                    string
	accountPassword                string
	characterID                    int
	authLevel                      int
	experience                     int
	levelPoint                     int
	masterExperience               int
	masterLevel                    int
	masterPoint                    int
	masterPointUsed                int
	fruitPoint                     int
	changeUp                       int     // 1=1转 2=2转 3=3转
	strength                       int     // 力量
	dexterity                      int     // 敏捷
	vitality                       int     // 体力
	energy                         int     // 智力
	leadership                     int     // 统力
	addStrength                    int     // 增加力量(套装+大师技能)
	addDexterity                   int     // 增加敏捷(套装+大师技能)
	addVitality                    int     // 增加体力(套装+大师技能)
	addEnergy                      int     // 增加智力(套装+洞装+大师技能)
	addLeadership                  int     // 增加统力(卓越翅膀+套装+大师技能)
	leftAttackMin                  int     // 左物攻min
	leftAttackMax                  int     // 左物攻max
	rightAttackMin                 int     // 右物攻min
	rightAttackMax                 int     // 右物攻max
	magic                          int     // 魔力提升
	magicAttackMin                 int     // 魔攻min
	magicAttackMax                 int     // 魔攻max
	curse                          int     // 诅咒力提升
	curseAttackMin                 int     // 诅咒攻击min
	curseAttackMax                 int     // 诅咒攻击max
	magicSpeed                     int     // 魔攻速度
	attackRatePVP                  int     // PVP攻击率(4维+380+大师技能)
	attackBonusPVP                 int     // PVP攻击加成(380)
	defenseRatePVP                 int     // PVP防御率(4维+380+大师技能)
	defenseBonusPVP                int     // PVP防御加成(380)
	harmonyIncSDRate               int     // PVP增加SD比率(再生)
	harmonyDecSDRate               int     // PVP降低SD比率(再生)
	harmonyIgnoreSDRate            int     // PVP无视SD概率(再生)
	knightGladiatorCalcSkillBonus  float64 // 战士/魔剑技能伤害加成(4维)
	impaleSkillCalc                float64 // 钻云枪技能伤害加成(4维)
	addSkillAttack                 int     // 增加技能攻击(套装+洞装+大师技能)
	setAddDamage                   int     // 套装增加伤害(套装)
	setIncreaseTwoHandWeaponDamage int     // 套装增加双手武器伤害(套装)
	criticalAttackRate             int     // 幸运一击概率(卓越+套装+洞装+大师技能)
	criticalAttackDamage           int     // 幸运一击伤害(套装+洞装)
	excellentAttackRate            int     // 卓越一击概率(卓越+套装+洞装+大师技能)
	excellentAttackDamage          int     // 卓越一击伤害(套装+洞装)
	doubleDamageRate               int     // 双倍伤害(套装+大师技能)
	ignoreDefenseRate              int     // 无视防御(套装+翅膀+大师技能)
	monsterDieGetHP                float64 // 杀怪回生(卓越+洞装+大师技能)
	monsterDieGetMP                float64 // 杀怪回蓝(卓越+洞装+大师技能)
	monsterDieGetMoney             float64 // 杀怪加钱(卓越+大师技能)
	recoverHP                      int     // 恢复HP(翅膀+项链+洞装+大师技能)
	recoverMP                      int     // 恢复MP(翅膀+项链+洞装+大师技能)
	recoverAG                      int     // 恢复AG(翅膀+项链+洞装+大师技能)
	agReduce                       int     // 减少AG消耗(翅膀+项链+洞装+大师技能)
	armorReduceDamage              int     // 防具减少伤害(卓越+洞装)
	armorReflectDamage             int     // 防具反射伤害(卓越+洞装)
	returnDamage                   int     // 反弹伤害(卓越翅膀+大师技能)
	recoverMaxHP                   int     // 恢复最大生命(卓越翅膀+大师技能)
	recoverMaxMP                   int     // 恢复最大魔法(卓越翅膀+大师技能)
	wingIncreaseDamage             int     // 翅膀增加伤害
	wingReduceDamage               int     // 翅膀减少伤害
	helperReduceDamage             int     // 天使减少伤害
	petIncreaseDamage              int     // pet增加伤害
	petReduceDamage                int     // pet减少伤害
	autoRecoverHPTick              int
	autoRecoverMPTick              int
	autoRecoverSD                  bool // 自动恢复SD(380)
	autoRecoverSDRate              int  // 自动恢复SD率(380)
	autoRecoverSDTime              time.Time
	delayRecoverHP                 int
	delayRecoverHPMax              int
	delayRecoverSD                 int
	delayRecoverSDMax              int
	// curseSpell           int
	autoSaveTime       time.Time
	inventory          item.Inventory
	chaosBoxOpen       bool
	chaosBox           item.ChaosBox
	inventoryExpansion int
	warehouse          item.Warehouse
	warehouseOpen      bool
	warehouseExpansion int
	warehouseMoney     int
	// PKCount                    int
	pkLevel int
	pet     *item.Item
	// PKTime                     int
	// PKTotalCount               int
	// guild                *guild.GuildInfo
	// guildName                     string
	// guildStatus                   int
	// guildUnionTimeStamp           int
	// guildNumber                   int
	// lastMoveTime                  time.Time
	// resets                        int
	// vipType                       uint8
	// vipEffect                     uint8
	// santaCount                    uint8
	// goblinTime                    time.Time
	// securityCheck                 bool
	// securityCode                  int
	// RegisterLMS                   uint8
	// registerLMSRoom               uint8
	// kanturuEntranceByNPC          bool
	// gensInfoLoad                  bool
	// questInfoLoad                 bool
	// wCoinP                        int
	// wCoinC                        int
	// goblinPoint                   int
	// periodItemEffectIndex         int
	// seedOptionList                [35]item.SocketOptionList
	// bonusOptionList               [7]item.SocketOptionList
	// setOptionList                 [2]item.SocketOptionList
	// refillHPSocketOption          uint16
	// refillMPSocketOption          uint16
	// socketOptionMonsterDieGetHP   uint16
	// socketOptionMonsterDieGetMana uint16
	// AGReduceRate                  uint8
	// muBotEnable                   bool
	// muBotTotalTime                time.Duration
	// muBotPayTime                  time.Duration
	// muBotTick                     time.Time
	// WarehouseExpansion            int
	// LastAuthTime                  time.Time
	// LastXorKey1                   [4]int
	// LaskXorKey2                   [4]int
	// bot                           bool
	// botIndex                      int
	// skillHellFire2State           int
	// skillHellFire2Count           int
	// skillHellFire2Time            time.Time
	// skillStrengthenHellFire2State int
	// skillStrengthenHellFire2Count int
	// skillStrengthenHellFire2Time  time.Time
	// reqWarehouseOpen              int
}

func (p *Player) Addr() string {
	return p.conn.Addr()
}

func (p *Player) Offline() {
	if p.offline {
		return
	}
	p.offline = true
	p.closeChaosBox()
	// todo
	p.saveCharacter()
	p.cancel()
}

func (p *Player) Push(msg any) {
	if p.offline {
		slog.Warn("Still pushing msg to offline player",
			"msg", msg, "player", p.Name)
		return
	}
	if len(p.msgChan) > 80 {
		p.Offline()
		return
	}
	p.msgChan <- msg
}

func (p *Player) ProcessAction() {}

func (p *Player) Process1000ms() {
	if p.ConnectState == object.ConnectStatePlaying {
		p.recoverHPSD()
		p.recoverMPAG()
		p.autoSave(time.Now())
	}
}

func (p *Player) recoverHPSD() {
	change := false
	if p.HP < p.MaxHP+p.AddHP {
		// auto recover HP
		p.autoRecoverHPTick++
		if p.autoRecoverHPTick >= 7 {
			p.autoRecoverHPTick = 0
			percent := 0.0
			// base item recover HP
			positions := []int{7, 9, 10, 11} // wing/pendant/ring
			for _, n := range positions {
				it := p.inventory.Items[n]
				if it != nil && it.Durability != 0 {
					percent += float64(it.AdditionRecoverHP)
				}
			}
			// master skill recover HP
			percent += 0.0
			if percent != 0.0 || p.recoverHP != 0 {
				hp := p.HP
				hp += int(float64(p.MaxHP+p.AddHP)*percent/100) + p.recoverHP
				switch {
				case hp < 1:
					hp = 1
				case hp > p.MaxHP+p.AddHP:
					hp = p.MaxHP + p.AddHP
				}
				p.HP = hp
				change = true
			}
		}
	} else {
		p.autoRecoverHPTick = 0
	}

	// auto recover SD
	if conf.CommonServer.GameServerInfo.ShieldSystemEnable && conf.CommonServer.GameServerInfo.SDAutoRefillEnable {
		now := time.Now()
		canRecover := true
		if conf.CommonServer.GameServerInfo.SDAutoRefillSafeZoneEnable {
			attr := maps.MapManager.GetMapAttr(p.MapNumber, p.X, p.Y)
			canRecover = attr&1 == 1 || p.autoRecoverSD
		}
		if p.autoRecoverSDTime.IsZero() || !canRecover || p.SD >= p.MaxSD {
			p.autoRecoverSDTime = now
		} else {
			elapsed := now.Sub(p.autoRecoverSDTime)
			if elapsed >= 10*time.Second {
				expressionA := p.MaxSD / 30
				expressionB := 100.0 + float64(p.autoRecoverSDRate)
				expression := int(float64(expressionA) * expressionB / 100 / 25)
				var sd int
				switch {
				case elapsed >= 25*time.Second:
					sd = expression * 3
				case elapsed >= 15*time.Second:
					sd = expression * 25 / 10
				default:
					sd = expression * 2
				}
				p.SD += sd
				if p.SD > p.MaxSD {
					p.SD = p.MaxSD
				}
				change = true
			}
		}
	}

	// posion delay recover HP
	if p.delayRecoverHP > 0 {
		hp := p.delayRecoverHPMax / 2
		if p.delayRecoverHP > hp {
			p.delayRecoverHP -= hp
		} else {
			hp = p.delayRecoverHP
			p.delayRecoverHP = 0
		}
		if p.HP < p.MaxHP+p.AddHP {
			hp += p.HP
			switch {
			case hp < 1:
				hp = 1
			case hp > p.MaxHP+p.AddHP:
				hp = p.MaxHP + p.AddHP
			}
			p.HP = hp
			change = true
		}
	}

	// posion delay recover SD
	if p.delayRecoverSD > 0 {
		sd := p.delayRecoverSDMax / 2
		if p.delayRecoverSD > sd {
			p.delayRecoverSD -= sd
		} else {
			sd = p.delayRecoverSD
			p.delayRecoverSD = 0
		}
		if p.SD < p.MaxSD {
			sd += p.SD
			if sd > p.MaxSD {
				sd = p.MaxSD
			}
			p.SD = sd
			change = true
		}
	}

	if change {
		p.PushHPSD(p.HP, p.SD)
	}
}

func (p *Player) recoverMPAG() {
	p.autoRecoverMPTick++
	if p.autoRecoverMPTick < 3 {
		return
	}
	p.autoRecoverMPTick = 0
	change := false
	if p.MP < p.MaxMP+p.AddMP {
		// base recover MP
		percent := 3.7
		// master skill recover MP
		percent += 0
		mp := p.MP
		mp += int(float64(p.MaxMP+p.AddMP)*percent/100) + p.recoverMP
		switch {
		case mp < 1:
			mp = 1
		case mp > p.MaxMP+p.AddMP:
			mp = p.MaxMP + p.AddMP
		}
		p.MP = mp
		change = true
	}
	if p.AG < p.MaxAG {
		// base recover AG
		percent := 3.0
		// master skill recover AG
		percent += 0
		if p.Class == int(class.Knight) {
			percent = 5
		}
		ag := p.AG
		ag += 5 + int(float64(p.MaxAG)*percent/100) + p.recoverAG
		switch {
		case ag < 1:
			ag = 1
		case ag > p.MaxAG:
			ag = p.MaxAG
		}
		p.AG = ag
		change = true
	}
	if change {
		p.PushMPAG(p.MP, p.AG)
	}
}

func (p *Player) autoSave(now time.Time) {
	if p.ConnectState != object.ConnectStatePlaying || p.Name == "" {
		return
	}
	if p.autoSaveTime.IsZero() {
		p.autoSaveTime = now
		return
	}
	if now.Sub(p.autoSaveTime) < autoSaveInterval {
		return
	}
	p.autoSaveTime = now
	if p.warehouseOpen {
		p.saveWarehouse()
	}
	p.saveCharacter()
}

func (p *Player) saveWarehouse() {
	account := model.Account{
		ID:             p.accountID,
		Warehouse:      p.warehouse,
		WarehouseMoney: p.warehouseMoney,
	}
	err := model.DB.UpdateAccountWarehouse(&account)
	if err != nil {
		slog.Error("saveWarehouse model.DB.UpdateAccountWarehouse",
			"err", err, "player", p.Name)
	}
}

func (p *Player) saveCharacter() {
	if p.Name == "" {
		return
	}
	c := model.Character{
		ID:                 p.characterID,
		ChangeUp:           p.changeUp,
		MapNumber:          p.MapNumber,
		X:                  p.X,
		Y:                  p.Y,
		Dir:                p.Dir,
		Money:              p.Money,
		Strength:           p.strength,
		Dexterity:          p.dexterity,
		Vitality:           p.vitality,
		Energy:             p.energy,
		Leadership:         p.leadership,
		Level:              p.Level,
		LevelPoint:         p.levelPoint,
		Experience:         p.experience,
		MasterLevel:        p.masterLevel,
		MasterPoint:        p.masterPoint,
		MasterExperience:   p.masterExperience,
		HP:                 p.HP,
		MP:                 p.MP,
		InventoryExpansion: p.inventoryExpansion,
		Inventory:          p.inventory,
		Skills:             p.Skills,
	}
	err := model.DB.UpdateCharacter(&c)
	if err != nil {
		slog.Error("saveCharacter model.DB.UpdateCharacter",
			"err", err, "player", p.Name)
	}
}

func (p *Player) spawnPosition() {
	gate := 0
	switch p.MapNumber {
	case maps.Arena, // 古战场
		maps.DuelArena, // 竞技场
		maps.Exile,     // 流放地
		maps.SantaTown: // 圣诞之地

	case maps.LorenMarket, // 罗兰市场
		maps.ImperialGuardian1, // 帝国要塞1
		maps.ImperialGuardian2, // 帝国要塞2
		maps.ImperialGuardian3, // 帝国要塞3
		maps.ImperialGuardian4, // 帝国要塞4
		maps.IllusionTemple1,   // 幻影寺院1
		maps.IllusionTemple2,   // 幻影寺院2
		maps.IllusionTemple3,   // 幻影寺院3
		maps.IllusionTemple4,   // 幻影寺院4
		maps.IllusionTemple5,   // 幻影寺院5
		maps.IllusionTemple6,   // 幻影寺院6
		maps.IllusionTemple7,   // 幻影寺院7
		maps.IllusionTemple8,   // 幻影寺院8
		maps.Doppelganger1,     //生魂广场1
		maps.Doppelganger2,     //生魂广场2
		maps.Doppelganger3,     //生魂广场3
		maps.Doppelganger4:     //生魂广场4
		gate = 333
	case maps.Lorencia, // 勇者大陆
		maps.Dungeon, // 地下城1~3
		maps.Kalima1, // 卡利玛1
		maps.Kalima2, // 卡利玛2
		maps.Kalima3, // 卡利玛3
		maps.Kalima4, // 卡利玛4
		maps.Kalima5, // 卡利玛5
		maps.Kalima6, // 卡利玛6
		maps.Kalima7: // 卡利玛7
		gate = 17 // Lorencia
	case maps.Devias, // 冰风谷1~4
		maps.Icarus,           // 天空之城
		maps.DevilSquare,      // 恶魔广场1~4
		maps.DevilSquare2,     // 恶魔广场5~7
		maps.DevilSquareFinal, // 恶魔广场
		maps.BloodCastle1,     // 血色城堡1
		maps.BloodCastle2,     // 血色城堡2
		maps.BloodCastle3,     // 血色城堡3
		maps.BloodCastle4,     // 血色城堡4
		maps.BloodCastle5,     // 血色城堡5
		maps.BloodCastle6,     // 血色城堡6
		maps.BloodCastle7,     // 血色城堡7
		maps.BloodCastle8,     // 血色城堡8
		maps.ChaosCastle1,     // 赤色要塞1
		maps.ChaosCastle2,     // 赤色要塞2
		maps.ChaosCastle3,     // 赤色要塞3
		maps.ChaosCastle4,     // 赤色要塞4
		maps.ChaosCastle5,     // 赤色要塞5
		maps.ChaosCastle6,     // 赤色要塞6
		maps.ChaosCastle7:     // 赤色要塞7
		gate = 22 // Devias
	case maps.Noria: // 仙踪林
		gate = 27 // Noria
	case maps.Elbeland: // 幻术园
		gate = 267 // Elbeland
	case maps.LostTower: // 失落之塔1~7
		gate = 42 // LostTower
	case maps.Atlans: // 亚特兰蒂斯1~3
		gate = 49 // Atlans
	case maps.Tarkan: // 死亡沙漠1~2
		gate = 57 // Tarkan
	case maps.Aida: // 幽暗森林1~2
		gate = 119 // Aida
	case maps.Kanturu: // 坎特鲁废墟1~3
		gate = 138 // Kanturu
	case maps.KanturuRemain: // 坎特鲁遗址
		gate = 139 // KanturuRemain
	case maps.Karutan1, // 卡伦特1
		maps.Karutan2: // 卡伦特2
		gate = 335 // Karutan1
	case maps.Raklion, // 冰霜之城
		maps.RaklionBoss: // 冰霜之城Boss
		gate = 287
	case maps.SwampOfCalmness: // 安宁池
		gate = 273 // SwampOfCalmness
	case maps.Acheron, // 阿卡伦
		maps.AcheronArcaWar, // 阿卡伦战役
		maps.Debenter,
		maps.DebenterArcaWar,
		maps.UrkMontain,
		maps.UrkMontainArcaWar:
		gate = 417
	case maps.Vulcanus: // 囚禁之岛
		gate = 294
	case maps.ValleyOfLoren: // 罗兰峡谷

	case maps.Crywolf:
		gate = 114 // Crywolf
	case maps.BalgassBarracks, // 巴卡斯兵营
		maps.BalgassRefuge: // 巴卡斯休息室
		gate = 256
	}
	maps.GateManager.Handle(gate, func(mapNumber, x, y, dir int) {
		if p.MapNumber != mapNumber {
			p.MapNumber = mapNumber
			p.LoadMiniMap()
		}
		p.X, p.Y = x, y
		p.TX, p.TY = x, y
		p.Dir = dir
	})
}

func (p *Player) Regen() {
	p.spawnPosition()
	maps.MapManager.SetMapAttrStand(p.MapNumber, p.TX, p.TY)
	p.CreateFrustum()
	p.HP = p.MaxHP + p.AddHP
	p.SD = p.MaxSD
	p.MP = p.MaxMP + p.AddMP
	p.AG = p.MaxAG
	reply := model.MsgReloadCharacterReply{
		X:          p.X,
		Y:          p.Y,
		MapNumber:  p.MapNumber,
		Dir:        p.Dir,
		HP:         p.HP,
		MP:         p.MP,
		SD:         p.SD,
		AG:         p.AG,
		Experience: int(p.experience),
		Money:      p.Money,
	}
	p.Push(&reply)
	p.PushSelfViewport()
}

func (p *Player) GetPKLevel() int {
	return p.pkLevel
}

func (p *Player) GetMasterLevel() int {
	return p.masterLevel
}

func (p *Player) IsMasterLevel() bool {
	return p.changeUp >= 2
}

func (p *Player) GetSkillMPAG(s *skill.Skill) (int, int) {
	mp := s.ManaUsage
	if p.HasBuff(effect.BuffInfinityArrow) {
		switch s.Index {
		case skill.SkillIndexTripleShot, skill.SkillIndexIceArrow,
			skill.SkillIndexPenetration, skill.SkillIndexMultiShot:
			for _, position := range []int{0, 1} {
				it := p.inventory.Items[position]
				if it == nil || (it.Code != item.Code(4, 7) && it.Code != item.Code(4, 15)) {
					continue
				}
				level := it.Level
				if level < 0 {
					level = 0
				}
				if level > 3 {
					level = 3
				}
				mp += skill.Settings.InfinityArrowMPConsumption[level]
				break
			}
		}
	}
	ag := s.BPUsage * (100 - p.agReduce) / 100
	return mp, ag
}

func (p *Player) UseSkillAmmo(s *skill.Skill, consume bool) bool {
	switch s.Index {
	case skill.SkillIndexTripleShot, skill.SkillIndexStarfall,
		skill.SkillIndexIceArrow, skill.SkillIndexPenetration, skill.SkillIndexMultiShot:
	default:
		return true
	}
	if p.HasBuff(effect.BuffInfinityArrow) {
		return true
	}
	for _, position := range []int{0, 1} {
		it := p.inventory.Items[position]
		if it == nil || it.Durability <= 0 ||
			(it.Code != item.Code(4, 7) && it.Code != item.Code(4, 15)) {
			continue
		}
		if consume {
			it.Durability--
			p.Push(&model.MsgItemDurabilityReply{Position: position, Durability: it.Durability, Flag: 0})
		}
		return true
	}
	return false
}

func (p *Player) GetChangeUp() int {
	return p.changeUp
}

func (p *Player) GetInventory() *item.Inventory {
	return &p.inventory
}

func (p *Player) GetInventoryItem(position int) *item.Item {
	return p.inventory.Item(position)
}

func (p *Player) SetDelayRecoverHP(hp, hpMax int) {
	p.delayRecoverHP = hp
	p.delayRecoverHPMax = hpMax
}

func (p *Player) SetDelayRecoverSD(sd, sdMax int) {
	p.delayRecoverSD = sd
	p.delayRecoverSDMax = sdMax
}

func (p *Player) Login(msg *model.MsgLogin) {
	// validate msg
	resp := model.MsgLoginReply{Result: 1}
	defer p.Push(&resp)
	account, err := model.DB.GetAccountByName(msg.Account)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.Result = 2
			return
		}
		slog.Error("Login model.DB.GetAccountByName",
			"err", err, "account", msg.Account)
		resp.Result = 7
		return
	}
	if account.Password != msg.Password {
		resp.Result = 0
		return
	}
	p.accountID = account.ID
	p.accountName = account.Name
	p.accountPassword = account.Password
	p.warehouseExpansion = account.WarehouseExpansion
	p.ConnectState = object.ConnectStateLogged

	// async
	// go func() {
	// 	account, err := model.DB.GetAccountByName(msg.Account)
	// 	p.actioner.PlayerAction(p.index, "SetAccount", &model.MsgSetAccount{
	// 		MsgLogin: msg,
	// 		Account:  account,
	// 		Err:      err,
	// 	})
	// }()
}

// func (p *Player) SetAccount(msg *model.MsgSetAccount) {
// 	resp := model.MsgLoginReply{Result: 1}
// 	defer p.Push(&resp)
// 	login := msg.MsgLogin
// 	account := msg.Account
// 	err := msg.Err
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			resp.Result = 2
// 			return
// 		}
// 		slog.Error("model.DB.GetAccountByName",
// 			"err", err, "account", account.Name)
// 		resp.Result = 7
// 		return
// 	}
// 	if account.Password != login.Password {
// 		resp.Result = 0
// 		return
// 	}
// }

func (p *Player) Logout(msg *model.MsgLogout) {
	p.closeChaosBox()
	maps.MapManager.ClearMapAttrStand(p.MapNumber, p.TX, p.TY)
	defer p.Push(&model.MsgLogoutReply{Flag: msg.Flag})
	switch msg.Flag {
	case 0: // close game
	case 1: // back to pick character
		// offline to login state
		p.ConnectState = object.ConnectStateLogged
		p.saveCharacter()
		p.Reset()
		p.pet = nil
	case 2: // back to pick server
		// offline to init state.
		// Do not close connection
	default:
		slog.Error("Logout", "flag", msg.Flag, "player", p.Name)
	}
}

func (p *Player) GetCharacterList(msg *model.MsgEmpty) {
	reply := model.MsgGetCharacterListReply{}

	// get account
	reply.EnableCharacterClass = 0xFF
	reply.WarehouseExpansion = p.warehouseExpansion
	p.Push(&model.MsgEnableCharacterClassReply{
		Class: reply.EnableCharacterClass,
	})
	p.Push(&model.MsgResetCharacterReply{
		Reset: "012345678901234567",
	})

	// get character list
	chars, err := model.DB.GetCharacterList(p.accountID)
	if err != nil {
		slog.Error("GetCharacterList model.DB.GetCharacterList",
			"err", err, "account", p.accountName)
		return
	}

	reply.CharacterList = make([]*model.MsgCharacter, len(chars))
	for i, c := range chars {
		reply.CharacterList[i] = &model.MsgCharacter{
			Index: c.Position,
			Name:  c.Name,
			Level: c.Level,
			// Level:       c.Level + c.MasterLevel,
			Class:       c.Class,
			ChangeUp:    c.ChangeUp,
			Inventory:   [9]*item.Item(c.Inventory.Items[:9]),
			GuildStatus: 0xFF,
			PKLevel:     0,
		}
	}
	p.Push(&reply)
}

func (p *Player) CreateCharacter(msg *model.MsgCreateCharacter) {
	reply := model.MsgCreateCharacterReply{Result: 0}
	defer p.Push(&reply)

	// validate msg
	if msg.Name == "" || msg.Class > 6 {
		slog.Error("CreateCharacter validate msg",
			"msg", msg, "account", p.accountName)
		return
	}

	// try to get an empty postion
	chars, err := model.DB.GetCharacterList(p.accountID)
	if err != nil {
		slog.Error("CreateCharacter model.DB.GetCharacterList",
			"err", err, "account", p.accountName)
		return
	}
	position := 0
	for i, char := range chars {
		if i != char.Position {
			position = i
			break
		}
		position++
	}
	if position > 4 {
		slog.Error("CreateCharacter over max character count",
			"account", p.accountName)
		return
	}
	// create character
	c := CharacterTable[msg.Class]
	c.AccountID = p.accountID
	c.Position = position
	c.Name = msg.Name
	if err := model.DB.CreateCharacter(&c); err != nil {
		slog.Error("CreateCharacter model.DB.CreateCharacter",
			"err", err, "account", p.accountName)
		return
	}
	// reply
	reply.Result = 1
	reply.Name = c.Name
	reply.Index = c.Position
	reply.Level = c.Level
	reply.Class = c.Class
}

func (p *Player) DeleteCharacter(msg *model.MsgDeleteCharacter) {
	reply := model.MsgDeleteCharacterReply{Result: 0}
	defer p.Push(&reply)

	if p.ConnectState == object.ConnectStatePlaying {
		return
	}

	// validate msg
	if msg.Name == "" || msg.Password == "" {
		slog.Error("DeleteCharacter validate msg",
			"msg", msg, "account", p.accountName)
		return
	}

	// check password
	if msg.Password != "1234567" {
		reply.Result = 2
		return
	}

	// delete character
	if err := model.DB.DeleteCharacterByName(p.accountID, msg.Name); err != nil {
		slog.Error("DeleteCharacter model.DB.DeleteCharacterByName",
			"err", err, "account", p.accountName, "character", msg.Name)
		return
	}

	// reply
	reply.Result = 1
}

func (p *Player) CheckCharacter(msg *model.MsgCheckCharacter) {
	p.Push(&model.MsgCheckCharacterReply{
		Result: 0,
	})
}

func (p *Player) LoadCharacter(msg *model.MsgLoadCharacter) {
	// validate msg
	if msg.Name == "" || msg.Position < 0 || msg.Position > 4 {
		slog.Error("LoadCharacter validate msg",
			"msg", msg, "account", p.accountName)
		return
	}

	// load character data from db
	c, err := model.DB.GetCharacterByName(p.accountID, msg.Name)
	if err != nil {
		slog.Error("LoadCharacter model.DB.GetCharacterByName",
			"err", err, "account", p.accountName)
		return
	}

	// set player with character data
	p.characterID = c.ID
	p.Name = c.Name
	p.Class = c.Class
	p.changeUp = c.ChangeUp
	p.Level = c.Level
	p.levelPoint = c.LevelPoint
	p.experience = c.Experience
	p.strength = c.Strength
	p.dexterity = c.Dexterity
	p.vitality = c.Vitality
	p.energy = c.Energy
	p.leadership = c.Leadership
	p.masterLevel = c.MasterLevel
	p.masterPoint = c.MasterPoint
	p.masterExperience = c.MasterExperience
	p.HP = c.HP
	p.MP = c.MP
	p.Skills = c.Skills
	p.Skills.FillSkillData(p.Class)
	p.inventory = c.Inventory
	p.inventoryExpansion = c.InventoryExpansion
	p.Money = c.Money
	p.MapNumber = c.MapNumber
	p.X, p.TX = c.X, c.X
	p.Y, p.TY = c.Y, c.Y
	p.Dir = c.Dir
	if p.Level <= 10 {
		p.spawnPosition()
	}
	maps.MapManager.SetMapAttrStand(p.MapNumber, p.TX, p.TY)
	p.CreateFrustum()
	p.MoveSpeed = 1000
	p.MaxRegenTime = 4 * time.Second
	p.ConnectState = object.ConnectStatePlaying
	p.autoSaveTime = time.Now()
	p.warehouseOpen = false
	p.Live = true
	p.State = 1

	experience := p.experience
	if p.IsMasterLevel() {
		experience = p.masterExperience
	}
	nextExperience := exp.ExperienceTable[p.Level]
	if p.IsMasterLevel() {
		experience = exp.MasterExperienceTable[p.masterLevel]
	}

	// p.Push(&model.MsgResetGameReply{})
	// reply
	p.Push(&model.MsgLoadCharacterReply{
		X:                  p.X,
		Y:                  p.Y,
		MapNumber:          p.MapNumber,
		Dir:                p.Dir,
		Experience:         experience,
		NextExperience:     nextExperience,
		LevelPoint:         p.levelPoint,
		Strength:           p.strength,
		Dexterity:          p.dexterity,
		Vitality:           p.vitality,
		Energy:             p.energy,
		Leadership:         p.leadership,
		Money:              p.Money,
		PKLevel:            p.pkLevel,
		CtlCode:            0,
		AddPoint:           0,
		MaxAddPoint:        122,
		MinusPoint:         0,
		MaxMinusPoint:      122,
		InventoryExpansion: p.inventoryExpansion,
	})
	p.LoadMiniMap()
	p.findAndUsePet()

	// reply inventory
	p.Push(&model.MsgItemListReply{
		Items: p.inventory.Items,
	})
	// reply master data
	p.Push(&model.MsgMasterDataReply{
		MasterLevel:          p.masterLevel,
		MasterExperience:     p.masterExperience,
		MasterNextExperience: exp.MasterExperienceTable[p.masterLevel],
		MasterPoint:          p.masterPoint,
	})
	p.pushSkillList()
	p.pushMasterSkillList()

	p.Push(&model.MsgMuKeyReply{
		MsgMuKey: c.MuKey,
	})
	p.Push(&model.MsgMuBotReply{
		MsgMuBot: c.MuBot,
	})
	p.PushSelfViewport()
	// reply weather
	p.PushWeather(maps.MapManager.GetMapWeather(p.MapNumber))

	// client will calculate character after receiving inventory msg and master msg
	// calculate
	p.calc()

	// go func() {
	// 	time.Sleep(100 * time.Millisecond) // get character info
	// 	p.actioner.PlayerAction(p.index, "SetCharacter", &model.MsgSetCharacter{Name: msg.Name})
	// }()
}

func (p *Player) findAndUsePet() {
	for _, it := range p.inventory.Items[12:204] {
		if it == nil {
			continue
		}
		switch it.Code {
		case item.Code(13, 2), // Horn of Uniria 兽角
			item.Code(13, 3),  // Horn of Dinorant 彩云兽
			item.Code(13, 4),  // Dark Horse 黑王马之角
			item.Code(13, 37): // Horn of Fenrir 炎狼兽之角
			if it.PetActive == 1 {
				if p.pet == nil && it.Durability > 0 {
					p.pet = it
				} else {
					it.PetActive = 0
				}
			}
		}
	}
}

func (p *Player) pushSkillList() {
	var skills []*skill.Skill
	p.Skills.ForEachActiveSkill(func(s *skill.Skill) {
		skills = append(skills, s)
	})
	sort.Sort(skill.SortedSkillSlice(skills))
	p.Push(&model.MsgSkillListReply{Skills: skills})
}

func (p *Player) pushMasterSkillList() {
	var skills []*model.MsgMasterSkill
	p.Skills.ForEachMasterSkill(func(index int, level int, curValue float32, nextValue float32) {
		skills = append(skills, &model.MsgMasterSkill{
			MasterSkillUIIndex:   index,
			MasterSkillLevel:     level,
			MasterSkillCurValue:  curValue,
			MasterSkillNextValue: nextValue,
		})
	})
	p.Push(&model.MsgMasterSkillListReply{Skills: skills})
}

func (p *Player) Talk(msg *model.MsgTalk) {
	reply := model.MsgTalkReply{
		Result: 0,
	}
	tobj := object.ObjectManager.GetObject(msg.Target)
	if tobj == nil {
		return
	}
	if math.Abs(float64(p.X-tobj.X)) > 5 ||
		math.Abs(float64(p.Y-tobj.Y)) > 5 {
		return
	}
	if conf.Server.Debug {
		slog.Debug("Talk",
			"object", p.Name, "position", fmt.Sprintf("(%d,%d)", p.X, p.Y),
			"target", tobj.Name, "position", fmt.Sprintf("(%d,%d)", tobj.X, tobj.Y))
	}
	// 1. Far: Move(Start) and Move(Stop) then Talk
	// 2. Close: Move(Start) then Talk
	// If object is close to the NPC, for example, object at (119,127)
	// wants to talk to target(NPC578) at (117,126), client would send
	// Talk request immediately instead of Move(Stop) request first,
	// which means the object is still moving while talking.
	// 3. Just right here: Talk directly
	if p.PathMoving {
		p.SetPosition(&model.MsgSetPosition{
			X: p.TX,
			Y: p.TY,
		})
	}
	p.TargetNumber = tobj.Index
	if tobj.NpcType != object.NpcTypeChaosMix {
		p.closeChaosBox()
	}
	switch tobj.NpcType {
	case object.NpcTypeShop:
		inventory := shop.ShopManager.GetShopInventory(tobj.Class, tobj.MapNumber)
		if tobj.Class == 492 {
			reply.Result = 34
		}
		p.Push(&reply)
		shopItemListReply := model.MsgTypeItemListReply{Type: 0}
		shopItemListReply.Items = inventory
		p.Push(&shopItemListReply)
	case object.NpcTypeWarehouse:
		account, err := model.DB.GetAccountByID(p.accountID)
		if err != nil {
			slog.Error("Talk model.DB.GetAccountByID",
				"err", err, "player", p.Name)
			return
		}
		reply.Result = 2
		p.Push(&reply)
		p.loadWarehouse(account)
		warehouseItemListReply := model.MsgTypeItemListReply{Type: 0}
		warehouseItemListReply.Items = account.Warehouse.Items
		p.Push(&warehouseItemListReply)
		warehouseMoneyReply := model.MsgWarehouseMoneyReply{
			Result:         1,
			WarehouseMoney: p.warehouseMoney,
			InventoryMoney: p.Money,
		}
		p.Push(&warehouseMoneyReply)
	case object.NpcTypeChaosMix:
		switch tobj.Class {
		case 452:
			reply.Result = 0x17 // Seed Master
		case 453:
			reply.Result = 0x18 // Seed Researcher
		default:
			return
		}
		p.chaosBoxOpen = true
		p.Push(&reply)
		p.Push(&model.MsgTypeItemListReply{Type: 3, MsgItemListReply: model.MsgItemListReply{Items: p.chaosBox.Items()}})
	case object.NpcTypeGoldarcher:
	case object.NpcTypePentagramMix:
	}
}

func (p *Player) loadWarehouse(account *model.Account) {
	p.warehouse = account.Warehouse
	p.warehouseMoney = account.WarehouseMoney
	p.warehouseOpen = true
}

func (p *Player) CloseTalkWindow(msg *model.MsgEmpty) {
	p.closeChaosBox()
	p.TargetNumber = -1
}

func (p *Player) CloseWarehouseWindow(msg *model.MsgEmpty) {
	p.saveWarehouse()
	p.warehouseOpen = false
	p.saveCharacter()
	reply := model.MsgCloseWarehouseWindowReply{}
	p.Push(&reply)
}

func (p *Player) KeepLive(msg *model.MsgKeepLive)      {}
func (p *Player) Hack(msg *model.MsgHack)              {}
func (p *Player) BattleCoreNotice(*model.MsgEmpty)     {}
func (p *Player) MapDataLoadingOK(msg *model.MsgEmpty) {}

func (p *Player) AddLevelPoint(msg *model.MsgAddLevelPoint) {
	reply := model.MsgAddLevelPointReply{}
	defer p.Push(&reply)
	if p.levelPoint < 1 {
		return
	}
	switch msg.Type {
	case 0:
		p.strength++
	case 1:
		p.dexterity++
	case 2:
		p.vitality++
	case 3:
		p.energy++
	case 4:
		p.leadership++
	default:
		return
	}
	p.levelPoint--
	p.calc()
	reply.Type = 0x10 + msg.Type
	switch msg.Type {
	case 2:
		reply.MaxHPMP = p.MaxHP + p.AddHP
	case 3:
		reply.MaxHPMP = p.MaxMP + p.AddMP
	}
	reply.MaxSD = p.MaxSD
	reply.MaxAG = p.MaxAG
}

func (p *Player) LearnMasterSkill(msg *model.MsgLearnMasterSkill) {
	reply := model.MsgLearnMasterSkillReply{
		Result:           0,
		MasterPoint:      p.masterPoint,
		MasterSkillIndex: -1,
	}
	defer p.Push(&reply)

	if p.masterLevel <= 0 ||
		p.masterPoint <= 0 {
		return
	}

	p.Skills.GetMaster(p.Class, msg.SkillIndex, p.masterPoint, func(point, uiIndex, index, level int, curValue, NextValue float32) {
		p.masterPoint -= point
		reply.Result = 1
		reply.MasterPoint -= point
		reply.MasterSkillUIIndex = uiIndex
		reply.MasterSkillIndex = index
		reply.MasterSkillLevel = level
		reply.MasterSkillCurValue = curValue
		reply.MasterSkillNextValue = NextValue
	})
}

func (p *Player) DefineMuKey(msg *model.MsgDefineMuKey) {
	if p.ConnectState != object.ConnectStatePlaying {
		return
	}
	err := model.DB.UpdateCharacterMuKey(p.characterID, msg)
	if err != nil {
		slog.Error("model.DB.UpdateCharacterMuKey",
			"err", err, "player", p.Name)
		return
	}
}

func (p *Player) DefineMuBot(msg *model.MsgDefineMuBot) {
	err := model.DB.UpdateCharacterMuBot(p.characterID, msg)
	if err != nil {
		slog.Error("model.DB.UpdateCharacterMuBot",
			"err", err, "player", p.Name)
		return
	}
}

func (p *Player) EnableMuBot(msg *model.MsgEnableMuBot) {
	reply := model.MsgEnableMuBotReply{Flag: msg.Flag}
	defer p.Push(&reply)
	switch msg.Flag {
	case 0:
	case 1:
	}
}

func (p *Player) MuunSystem(msg *model.MsgMuunSystem)      {}
func (p *Player) StartPartyNumberPosition(*model.MsgEmpty) {}
func (p *Player) StopPartyNumberPosition(*model.MsgEmpty)  {}
