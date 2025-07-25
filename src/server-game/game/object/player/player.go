package player

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/exp"
	"github.com/xujintao/balgass/src/server-game/game/formula"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/shop"
	"github.com/xujintao/balgass/src/server-game/game/skill"
	"gorm.io/gorm"
)

func init() {
	var characterList []*model.Character
	conf.JSON(conf.PathCommon, "players/character.json", &characterList)
	CharacterTable = make(characterTable)
	for _, c := range characterList {
		CharacterTable[c.Class] = *c
	}
}

var CharacterTable characterTable

type characterTable map[int]model.Character

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
	offline           bool
	conn              object.Conn
	msgChan           chan any
	cancel            context.CancelFunc
	actioner          object.Actioner
	accountID         int
	accountName       string
	accountPassword   string
	characterID       int
	authLevel         int
	experience        int
	levelPoint        int
	masterExperience  int
	masterLevel       int
	masterPoint       int
	masterPointUsed   int
	fruitPoint        int
	strength          int
	dexterity         int
	vitality          int
	energy            int
	leadership        int
	addStrength       int
	addDexterity      int
	addVitality       int
	addEnergy         int
	addLeadership     int
	autoRecoverHPTick int
	autoRecoverMPTick int
	autoRecoverSDTime time.Time
	delayRecoverHP    int
	delayRecoverHPMax int
	delayRecoverSD    int
	delayRecoverSDMax int
	recoverHP         int // 恢复生命(翅膀+项链+大师技能)
	magic             int
	magicAttackMin    int // 魔攻min
	magicAttackMax    int // 魔攻max
	curse             int
	curseAttackMin    int // 诅咒min
	curseAttackMax    int // 诅咒max
	attackRatePVP     int
	defenseRatePVP    int
	magicSpeed        int // 魔攻速度
	// curseSpell           int
	inventory          item.Inventory
	inventoryExpansion int
	warehouse          item.Warehouse
	warehouseExpansion int
	warehouseMoney     int
	changeUp           int // 1=1转 2=2转 3=3转
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
	// jewelHarmonyEffect            item.JewelHarmonyItemEffect
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
	increaseAttackMin             int     // 增加最小攻击(套装+洞装+大师技能)
	increaseAttackMax             int     // 增加最大攻击(套装+洞装+大师技能)
	increaseMagicAttack           int     // 增加魔攻(套装+洞装+大师技能)
	increaseSkillAttack           int     // 增加技能攻击(套装+卓越+大师技能)
	setAddDamage                  int     // 增加伤害
	criticalAttackRate            int     // 幸运一击概率
	criticalAttackDamage          int     // 幸运一击伤害
	excellentAttackRate           int     // 卓越一击概率
	excellentAttackDamage         int     // 卓越一击伤害
	monsterDieGetHP               float64 // 杀怪回生(套装+卓越+大师技能)
	monsterDieGetMP               float64 // 杀怪回蓝(套装+卓越+大师技能)
	monsterDieGetMoney            float64 // 杀怪加钱(卓越)
	armorReduceDamage             int     // 防具减少伤害(卓越+洞装)
	wingIncreaseDamage            int     // 翅膀增加伤害
	wingReduceDamage              int     // 翅膀减少伤害
	helperReduceDamage            int     // 天使减少伤害
	petIncreaseDamage             int     // pet增加伤害
	petReduceDamage               int     // pet减少伤害
	armorReflectDamage            int     // 防具反射伤害(卓越+洞装)
	doubleDamageRate              int     // 双倍伤害(套装+大师技能)
	ignoreDefenseRate             int     // 无视防御(套装+翅膀+大师技能)
	returnDamage                  int     // 反弹伤害(翅膀+大师技能)
	recoverMaxHP                  int     // 恢复最大生命(翅膀+大师技能)
	recoverMaxMP                  int     // 恢复最大魔法(翅膀+大师技能)
	increaseTwoHandWeaponDamage   int     // 增加双手武器伤害
	item380Effect                 item.Item380Effect
	setFull                       bool
	knightGladiatorCalcSkillBonus float64
	impaleSkillCalc               float64
}

func (p *Player) Addr() string {
	return p.conn.Addr()
}

func (p *Player) Offline() {
	if p.offline {
		return
	}
	p.offline = true
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
	}
}

func (p *Player) recoverHPSD() {
	change := false
	if p.HP < p.MaxHP {
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
			if percent != 0.0 {
				hp := p.HP
				hp += int(float64(p.MaxHP) * percent / 100)
				switch {
				case hp < 1:
					hp = 1
				case hp > p.MaxHP:
					hp = p.MaxHP
				}
				p.HP = hp
				change = true
			}
		}
	} else {
		p.autoRecoverHPTick = 0
	}

	// auto recover SD
	if p.SD < p.MaxSD {
		if conf.CommonServer.GameServerInfo.SDAutoRefillSafeZoneEnable {
			attr := maps.MapManager.GetMapAttr(p.MapNumber, p.X, p.Y)
			if attr&1 == 1 {
				now := time.Now()
				if now.After(p.autoRecoverSDTime.Add(time.Second * 1)) {
					p.autoRecoverSDTime = now
					expressionA := float64(p.MaxSD) / 30
					expressionB := 100.0 // 380 option
					sd := p.SD
					sd += int(expressionA * expressionB / 100)
					switch {
					case sd < 1:
						sd = 1
					case sd > p.MaxSD:
						sd = p.MaxSD
					}
					p.SD = sd
					change = true
				}
			}
		}
	}

	// posion delay recover hp
	if p.delayRecoverHP > 0 {
		hp := p.delayRecoverHPMax / 2
		if p.delayRecoverHP > hp {
			p.delayRecoverHP -= hp
		} else {
			hp = p.delayRecoverHP
			p.delayRecoverHP = 0
		}
		if p.HP < p.MaxHP {
			hp += p.HP
			switch {
			case hp < 1:
				hp = 1
			case hp > p.MaxHP:
				hp = p.MaxHP
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
	if p.MP < p.MaxMP {
		// base recover MP
		percent := 3.7
		// master skill recover MP
		percent += 0
		mp := p.MP
		mp += int(float64(p.MaxMP) * percent / 100)
		switch {
		case mp < 1:
			mp = 1
		case mp > p.MaxMP:
			mp = p.MaxMP
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
		ag += 5 + int(float64(p.MaxAG)*percent/100)
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

func (p *Player) SpawnPosition() {
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
	maps.GateMoveManager.Move(gate, func(mapNumber, x, y, dir int) {
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
	p.HP = p.MaxHP
	p.SD = p.MaxSD
	p.MP = p.MaxMP
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
}

func (p *Player) calc() {
	leftHand := p.inventory.Items[0]
	rightHand := p.inventory.Items[1]
	glove := p.inventory.Items[5]
	boot := p.inventory.Items[6]
	wing := p.inventory.Items[7]
	helper := p.inventory.Items[8]

	p.addStrength = 0
	p.addDexterity = 0
	p.addVitality = 0
	p.addEnergy = 0
	p.addLeadership = 0
	p.increaseAttackMin = 0
	p.increaseAttackMax = 0
	p.increaseMagicAttack = 0
	p.increaseSkillAttack = 0
	p.setAddDamage = 0
	p.criticalAttackRate = 0
	p.criticalAttackDamage = 0
	p.excellentAttackRate = 0
	p.excellentAttackDamage = 0
	p.monsterDieGetHP = 0
	p.monsterDieGetMP = 0
	p.monsterDieGetMoney = 0
	p.armorReduceDamage = 0
	p.armorReflectDamage = 0
	p.ignoreDefenseRate = 0
	p.returnDamage = 0
	p.recoverMaxHP = 0
	p.recoverMaxMP = 0
	p.wingIncreaseDamage = 0
	p.wingReduceDamage = 0
	p.helperReduceDamage = 0
	p.petIncreaseDamage = 0
	p.petReduceDamage = 0

	// wing item contribution
	if wing != nil && wing.ExcellentWing2Leadership {
		p.addLeadership += 10 + 5*wing.Level
	}
	// set item contribution
	p.calcSetItem(true)
	// master skill contribution

	strength := p.strength + p.addStrength
	dexterity := p.dexterity + p.addDexterity
	vitality := p.vitality + p.addVitality
	energy := p.energy + p.addEnergy
	leadership := p.leadership + p.addLeadership
	level := p.Level + p.masterLevel

	// base attack
	leftAttackMin, leftAttackMax := 0, 0
	rightAttackMin, rightAttackMax := 0, 0
	switch class.Class(p.Class) {
	case class.Wizard:
		formula.WizardDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.WizardMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.Knight:
		formula.KnightDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.KnightMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.Elf:
		if leftHand != nil && leftHand.KindB == item.KindBCrossbow ||
			rightHand != nil && rightHand.KindB == item.KindBBow {
			formula.ElfWithBowDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		} else {
			formula.ElfWithoutBowDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		}
		formula.ElfMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.Magumsa:
		formula.GladiatorDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.GladiatorMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.DarkLord:
		formula.LordDamageCalc(strength, dexterity, vitality, energy, leadership, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.LordMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	case class.Summoner:
		formula.SummonerDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.SummonerMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax, &p.curseAttackMin, &p.curseAttackMax)
	case class.RageFighter:
		formula.RageFighterDamageCalc(strength, dexterity, vitality, energy, &leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		formula.RageFighterMagicDamageCalc(energy, &p.magicAttackMin, &p.magicAttackMax)
	}

	// base defense
	formula.CalcDefense(p.Class, dexterity, &p.Defense)

	// base attack/defense success rate
	formula.CalcAttackSuccessRate_PvM(p.Class, strength, dexterity, leadership, p.Level, &p.AttackRate)
	formula.CalcAttackSuccessRate_PvP(p.Class, dexterity, p.Level, &p.attackRatePVP)
	formula.CalcDefenseSuccessRate_PvM(p.Class, dexterity, &p.DefenseRate)
	formula.CalcDefenseSuccessRate_PvP(p.Class, dexterity, p.Level, &p.defenseRatePVP)

	// base speed
	formula.CalcAttackSpeed(p.Class, dexterity, &p.AttackSpeed, &p.magicSpeed)

	// Stat Specialization
	// calc Stat Specialization: increase attack power
	// http://muonline.webzen.com/guides/219/1976/season-9/season-9-character-renewal
	// - (Bonus stat generated by equipping items like weapon, armor, wing, or master skill,
	// etc is not applied to specialization calculation.)
	var options []*model.MsgStatSpec
	var percent float64

	// STAT_OPTION_INC_ATTACK_POWER
	formula.StatSpec_GetPercent(p.Class, 1, strength, dexterity, vitality, energy, leadership, &percent)
	min := leftAttackMin * int(percent) / 100
	max := leftAttackMax * int(percent) / 100
	leftAttackMin += min
	leftAttackMax += max
	rightAttackMin += min
	rightAttackMax += max
	options = append(options, &model.MsgStatSpec{ID: 1, Min: min, Max: max})

	// STAT_OPTION_INC_MAGIC_DAMAGE
	formula.StatSpec_GetPercent(p.Class, 9, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.magicAttackMin * int(percent) / 100
	max = p.magicAttackMax * int(percent) / 100
	p.magicAttackMin += min
	p.magicAttackMax += max
	options = append(options, &model.MsgStatSpec{ID: 9, Min: min, Max: max})

	// STAT_OPTION_INC_CURSE_DAMAGE
	formula.StatSpec_GetPercent(p.Class, 10, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.curseAttackMin * int(percent) / 100
	max = p.curseAttackMax * int(percent) / 100
	p.curseAttackMin += min
	p.curseAttackMax += max
	options = append(options, &model.MsgStatSpec{ID: 10, Min: min, Max: max})

	// STAT_OPTION_INC_DEFENSE
	formula.StatSpec_GetPercent(p.Class, 4, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.Defense * int(percent) / 100
	p.Defense += min
	options = append(options, &model.MsgStatSpec{ID: 4, Min: min})

	// STAT_OPTION_INC_ATTACK_RATE
	formula.StatSpec_GetPercent(p.Class, 2, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.AttackRate * int(percent) / 100
	p.AttackRate += min
	options = append(options, &model.MsgStatSpec{ID: 2, Min: min})

	// STAT_OPTION_INC_ATTACK_RATE_PVP
	formula.StatSpec_GetPercent(p.Class, 3, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.attackRatePVP * int(percent) / 100
	p.attackRatePVP += min
	options = append(options, &model.MsgStatSpec{ID: 3, Min: min})

	// STAT_OPTION_INC_DEFENSE_RATE
	formula.StatSpec_GetPercent(p.Class, 6, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.DefenseRate * int(percent) / 100
	p.DefenseRate += min
	options = append(options, &model.MsgStatSpec{ID: 6, Min: min})

	// STAT_OPTION_INC_DEFENSE_RATE_PVP
	formula.StatSpec_GetPercent(p.Class, 7, strength, dexterity, vitality, energy, leadership, &percent)
	min = p.defenseRatePVP * int(percent) / 100
	p.defenseRatePVP += min
	options = append(options, &model.MsgStatSpec{ID: 7, Min: min})

	// weapon and weapon addition attack
	if leftHand != nil {
		leftAttackMin += leftHand.AttackMin + leftHand.AdditionAttack
		leftAttackMax += leftHand.AttackMax + leftHand.AdditionAttack
		p.magicAttackMin += leftHand.AdditionMagicAttack
		p.magicAttackMax += leftHand.AdditionMagicAttack
		p.curseAttackMin += leftHand.AdditionCurseAttack
		p.curseAttackMax += leftHand.AdditionCurseAttack
		p.curse = leftHand.Magic
	}
	if rightHand != nil {
		rightAttackMin += rightHand.AttackMin + rightHand.AdditionAttack
		rightAttackMax += rightHand.AttackMax + rightHand.AdditionAttack
		p.magicAttackMin += rightHand.AdditionMagicAttack
		p.magicAttackMax += rightHand.AdditionMagicAttack
		p.curseAttackMin += rightHand.AdditionCurseAttack
		p.curseAttackMax += rightHand.AdditionCurseAttack
		p.magic = rightHand.Magic
	}

	// wing addition attack
	if wing != nil {
		leftAttackMin += wing.AdditionAttack
		leftAttackMax += wing.AdditionAttack
		rightAttackMin += wing.AdditionAttack
		rightAttackMax += wing.AdditionAttack
		p.magicAttackMin += wing.AdditionMagicAttack
		p.magicAttackMax += wing.AdditionMagicAttack
		p.curseAttackMin += wing.AdditionCurseAttack
		p.curseAttackMax += wing.AdditionCurseAttack
	}

	// glove speed
	if glove != nil {
		p.AttackSpeed += glove.AttackSpeed
		p.magicSpeed += glove.AttackSpeed
	}

	// armor(shield|armor|wing) addition defense(rate)
	for i := 1; i <= 7; i++ {
		it := p.inventory.Items[i]
		if it != nil {
			p.Defense += it.Defense
			p.Defense += it.AdditionDefense
			p.DefenseRate += it.SuccessfulBlocking
			p.DefenseRate += it.AdditionDefenseRate
			p.recoverHP = it.AdditionRecoverHP
		}
	}

	// armor bonus defense and defense rate
	// defense level>=10 bonus of item the same type contributed
	// defense success rate bonus of item the same type contributed
	sameCount := 0
	level10Count := 0
	level11Count := 0
	level12Count := 0
	level13Count := 0
	level14Count := 0
	level15Count := 0
	if boot != nil {
		refer := boot.Code % item.MaxItemIndex
		for i := 2; i <= 6; i++ {
			it := p.inventory.Items[i]
			if it != nil && it.Code%item.MaxItemIndex == refer {
				sameCount++
				if it.Level > 9 {
					level10Count++
				}
				if it.Level > 10 {
					level11Count++
				}
				if it.Level > 11 {
					level12Count++
				}
				if it.Level > 12 {
					level13Count++
				}
				if it.Level > 13 {
					level14Count++
				}
				if it.Level > 14 {
					level15Count++
				}
			}
		}
		if p.Class == int(class.Magumsa) || p.Class == int(class.RageFighter) {
			sameCount++
			level10Count++
			level11Count++
			level12Count++
			level13Count++
			level14Count++
			level15Count++
		}
		if sameCount == 5 {
			n := 0
			switch {
			case level15Count == 5:
				n = 30
			case level14Count == 5:
				n = 25
			case level13Count == 5:
				n = 20
			case level12Count == 5:
				n = 15
			case level11Count == 5:
				n = 10
			case level10Count == 5:
				n = 5
			}
			p.Defense += p.Defense * n / 100
			p.DefenseRate += p.DefenseRate / 10
		}
	}

	// hp/mp
	c := CharacterTable[p.Class]
	p.MaxHP = c.HP + int(float32(level-1)*c.LevelHP) + int(float32(vitality-c.Vitality)*c.VitalityHP)
	p.MaxMP = c.MP + int(float32(level-1)*c.LevelMP) + int(float32(energy-c.Energy)*c.EnergyMP)

	// sd
	sdGageConstA := conf.CommonServer.GameServerInfo.SDGageConstA
	sdGageConstB := conf.CommonServer.GameServerInfo.SDGageConstB
	expressionA := strength + dexterity + vitality + energy
	if p.Class == int(class.DarkLord) {
		expressionA += leadership
	}
	expressionB := level * level / sdGageConstB
	p.MaxSD = expressionA*sdGageConstA/10 + expressionB + p.Defense

	// ag
	f := func(s, d, v, e, l float32) int {
		return int(float32(strength)*s + float32(dexterity)*d + float32(vitality)*v + float32(energy)*e + float32(leadership)*l)
	}
	switch class.Class(p.Class) {
	case class.Wizard:
		p.MaxAG = f(0.2, 0.4, 0.3, 0.2, 0)
	case class.Knight:
		p.MaxAG = f(0.15, 0.2, 0.3, 1.0, 0)
	case class.Elf:
		p.MaxAG = f(0.3, 0.2, 0.3, 0.2, 0)
	case class.Magumsa:
		p.MaxAG = f(0.2, 0.25, 0.3, 0.15, 0)
	case class.DarkLord:
		p.MaxAG = f(0.3, 0.2, 0.1, 0.15, 0.3)
	case class.Summoner:
		p.MaxAG = f(0.2, 0.25, 0.3, 0.15, 0)
	case class.RageFighter:
		p.MaxAG = f(0.15, 0.2, 0.3, 1.0, 0)
	}

	p.calcSetItem(false)
	// set item contribution
	leftAttackMin += p.increaseAttackMin
	leftAttackMax += p.increaseAttackMax
	rightAttackMin += p.increaseAttackMin
	rightAttackMax += p.increaseAttackMax

	// excellent item contribution
	for i := 0; i < 12; i++ {
		it := p.inventory.Items[i]
		if it != nil {
			if it.Lucky {
				p.criticalAttackRate += 5
			}
			if it.ExcellentAttackRate {
				p.excellentAttackRate += 10
			}
			if it.ExcellentAttackLevel {
				value := level / 20
				switch i {
				case 0:
					leftAttackMin += value
					leftAttackMax += value
				case 1:
					rightAttackMin += value
					rightAttackMax += value
				}
				p.magicAttackMin += value
				p.magicAttackMax += value
			}
			if it.ExcellentAttackPercent {
				switch i {
				case 0:
					leftAttackMin += leftAttackMin * 2 / 100
					leftAttackMax += leftAttackMax * 2 / 100
				case 1:
					rightAttackMin += rightAttackMin * 2 / 100
					rightAttackMax += rightAttackMax * 2 / 100
				}
				p.magicAttackMin += p.magicAttackMin * 2 / 100
				p.magicAttackMax += p.magicAttackMax * 2 / 100
			}
			if it.ExcellentAttackSpeed {
				p.AttackSpeed += 7
				p.magicSpeed += 7
			}
			if it.ExcellentAttackHP {
				p.monsterDieGetHP += 1.0 / 8
			}
			if it.ExcellentAttackMP {
				p.monsterDieGetMP += 1.0 / 8
			}
			if it.ExcellentDefenseHP {
				p.MaxHP += p.MaxHP * 4 / 100
			}
			if it.ExcellentDefenseMP {
				p.MaxMP += p.MaxMP * 4 / 100
			}
			if it.ExcellentDefenseReduce {
				p.armorReduceDamage += 4
			}
			if it.ExcellentDefenseReflect {
				p.armorReflectDamage += 5
			}
			if it.ExcellentDefenseRate {
				p.DefenseRate += p.DefenseRate * 10 / 100
			}
			if it.ExcellentDefenseMoney {
				p.monsterDieGetMoney += float64(30) / 100
			}
			if it.ExcellentWing2Speed {
				p.AttackSpeed += 5
				p.magicSpeed += 5
			}
			if it.ExcellentWing2AG {
				p.MaxAG += 50
			}
			if it.ExcellentWing2Ignore {
				p.ignoreDefenseRate += 3
			}
			if it.ExcellentWing2MP {
				p.MaxMP += 50 + 5*it.Level
			}
			if it.ExcellentWing2HP {
				p.MaxHP += 50 + 5*it.Level
			}
			if it.ExcellentWing3Ignore {
				p.ignoreDefenseRate += 5
			}
			if it.ExcellentWing3Return {
				p.returnDamage += 5
			}
			if it.ExcellentWing3HP {
				p.recoverMaxHP += 5
			}
			if it.ExcellentWing3MP {
				p.recoverMaxMP += 5
			}
		}
	}

	// convert left/right attack to attack
	left, right := false, false
	if leftHand != nil &&
		leftHand.KindA == item.KindAWeapon &&
		leftHand.Code != item.Code(4, 7) &&
		leftHand.Code != item.Code(4, 15) {
		left = true
	}
	if rightHand != nil &&
		rightHand.KindA == item.KindAWeapon &&
		rightHand.Code != item.Code(4, 7) &&
		rightHand.Code != item.Code(4, 15) {
		right = true
	}
	switch {
	case left && right:
		switch class.Class(p.Class) {
		case class.Knight, class.Magumsa:
			if leftHand.Code == rightHand.Code {
				formula.CalcTwoSameWeaponBonus(
					leftAttackMin, leftAttackMax, rightAttackMin, rightAttackMax,
					&leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
			} else {
				formula.CalcTwoDifferentWeaponBonus(
					leftAttackMin, leftAttackMax, rightAttackMin, rightAttackMax,
					&leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
			}
		case class.RageFighter:
			formula.CalcRageFighterTwoWeaponBonus(
				leftAttackMin, leftAttackMax, rightAttackMin, rightAttackMax,
				&leftAttackMin, &rightAttackMin, &leftAttackMax, &rightAttackMax)
		}
		p.AttackMin = leftAttackMin + rightAttackMin
		p.AttackMax = leftAttackMax + rightAttackMax
		p.AttackSpeed += (leftHand.AttackSpeed + rightHand.AttackSpeed) / 2
		p.magicSpeed += (leftHand.AttackSpeed + rightHand.AttackSpeed) / 2
	case left:
		p.AttackMin = leftAttackMin
		p.AttackMax = leftAttackMax
		p.AttackSpeed += leftHand.AttackSpeed
		p.magicSpeed += leftHand.AttackSpeed
	case right:
		p.AttackMin = rightAttackMin
		p.AttackMax = rightAttackMax
		p.AttackSpeed += rightHand.AttackSpeed
		p.magicSpeed += rightHand.AttackSpeed
	default:
		p.AttackMin = (leftAttackMin + rightAttackMin) / 2
		p.AttackMax = (leftAttackMax + rightAttackMax) / 2
	}

	// convert magic/curse to magic/curse attack
	p.magicAttackMin += p.magicAttackMax * p.magic / 100
	p.magicAttackMax += p.magicAttackMax * p.magic / 100
	p.curseAttackMin += p.curseAttackMin * p.curseAttackMax / 100
	p.curseAttackMax += p.curseAttackMin * p.curseAttackMax / 100

	// wing
	if wing != nil {
		wingIncreaseDamage := 0
		formula.Wings_CalcIncAttack(wing.Code, wing.Level, &wingIncreaseDamage)
		p.wingIncreaseDamage = wingIncreaseDamage - 100
		wingReduceDamage := 0
		formula.Wings_CalcAbsorb(wing.Code, wing.Level, &wingReduceDamage)
		p.wingReduceDamage = 100 - wingReduceDamage
	}

	// helper item contribution
	switch {
	case p.equippedItem(helper, item.Code(13, 1)): // 小恶魔
		percent := conf.PetRing.Pets.Imp.AddAttackPercent
		p.AttackMin += p.AttackMin * percent / 100
		p.AttackMax += p.AttackMax * percent / 100
		p.magicAttackMin += p.magicAttackMin * percent / 100
		p.magicAttackMax += p.magicAttackMax * percent / 100
		p.curseAttackMin += p.curseAttackMin * percent / 100
		p.curseAttackMax += p.curseAttackMax * percent / 100
	case p.equippedItem(helper, item.Code(13, 64)): // 强化恶魔
		percent := conf.PetRing.Pets.Demon.AddAttackPercent
		p.AttackMin += p.AttackMin * percent / 100
		p.AttackMax += p.AttackMax * percent / 100
		p.magicAttackMin += p.magicAttackMin * percent / 100
		p.magicAttackMax += p.magicAttackMax * percent / 100
		p.curseAttackMin += p.curseAttackMin * percent / 100
		p.curseAttackMax += p.curseAttackMax * percent / 100
		p.AttackSpeed += conf.PetRing.Pets.Demon.AddAttackSpeed
		p.magicSpeed += conf.PetRing.Pets.Demon.AddAttackSpeed
	case p.equippedItem(helper, item.Code(13, 0)): // 守护天使
		p.MaxHP += conf.PetRing.Pets.Angel.AddHP
		p.helperReduceDamage = conf.PetRing.Pets.Angel.ReduceDamagePercent
	case p.equippedItem(helper, item.Code(13, 65)): // 强化天使
		p.MaxHP += conf.PetRing.Pets.SpiritAngel.AddHP
		p.helperReduceDamage = conf.PetRing.Pets.SpiritAngel.ReduceDamagePercent
	case p.equippedItem(helper, item.Code(13, 80)): // 熊猫
		p.Defense += conf.PetRing.Pets.Panda.AddDefenseValue
	}

	// pet item contribution
	if p.pet != nil {
		switch p.pet.Code {
		case item.Code(13, 3): // Horn of Dinorant 彩云兽
			p.petIncreaseDamage = 15
			p.petReduceDamage = 10
		case item.Code(13, 4): // Dark Horse 黑王马之角
			p.petReduceDamage = (30 + p.pet.Level) / 2
		case item.Code(13, 37): // Horn of Fenrir 炎狼兽之角
			p.petIncreaseDamage = 33
			p.petReduceDamage = 10
		}
	}

	// ...

	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
	if p.MP > p.MaxMP {
		p.MP = p.MaxMP
	}
	if p.SD > p.MaxSD {
		p.SD = p.MaxSD
	}
	if p.AG > p.MaxAG {
		p.AG = p.MaxAG
	}

	// skill bonus
	formula.Knight_Gladiator_CalcSkillBonus(p.Class, energy, &p.knightGladiatorCalcSkillBonus)
	formula.ImpaleSkillCalc(p.Class, energy, &p.impaleSkillCalc)

	// Push
	p.Push(&model.MsgStatSpecReply{
		Options: options,
	})
	p.Push(&model.MsgAttackSpeedReply{
		AttackSpeed: p.AttackSpeed,
		MagicSpeed:  p.magicSpeed,
	})
	p.PushMaxHPSD(p.MaxHP, p.MaxSD)
	p.PushMaxMPAG(p.MaxMP, p.MaxAG)
	p.PushHPSD(p.HP, p.SD)
	p.PushMPAG(p.MP, p.AG)
}

func (p *Player) EquipmentChanged() {
	// 1, change skill
	newItemSkills := make(map[int]struct{})
	primaryHandWeapon := p.inventory.Items[0]
	if primaryHandWeapon != nil && primaryHandWeapon.SkillIndex != 0 {
		newItemSkills[primaryHandWeapon.SkillIndex] = struct{}{}
	}
	secondaryHandWeapon := p.inventory.Items[1]
	if secondaryHandWeapon != nil && secondaryHandWeapon.SkillIndex != 0 {
		newItemSkills[secondaryHandWeapon.SkillIndex] = struct{}{}
	}
	oldItemSkills := make(map[int]struct{})
	for _, s := range p.Skills {
		if s.Index < 300 && s.SkillBase.ItemSkill {
			oldItemSkills[s.Index] = struct{}{}
		}
	}
	var needLearnSkills []int
	for newSkill := range newItemSkills {
		if _, ok := oldItemSkills[newSkill]; !ok {
			needLearnSkills = append(needLearnSkills, newSkill)
		}
	}
	var needForgetSkills []int
	for oldSkill := range oldItemSkills {
		if _, ok := newItemSkills[oldSkill]; !ok {
			needForgetSkills = append(needForgetSkills, oldSkill)
		}
	}
	for _, index := range needLearnSkills {
		if s, ok := p.LearnSkill(index); ok {
			p.Push(&model.MsgSkillOneReply{
				Flag:  -2,
				Skill: s,
			})
		}
	}
	for _, index := range needForgetSkills {
		if s, ok := p.ForgetSkill(index); ok {
			p.Push(&model.MsgSkillOneReply{
				Flag:  -1,
				Skill: s,
			})
		}
	}
	// 2, calculate player
	p.calc()
}

func (p *Player) equippedItem(it *item.Item, code int) bool {
	if it == nil || it.Durability == 0 {
		return false
	}
	return it.Code == code
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
			if it.HarmonyOption == 1 {
				p.pet = it
			}
		}
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
		return
	}
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
	return s.ManaUsage, s.BPUsage
}

func (p *Player) addSetEffect(index item.SetEffectType, value int, base bool) {
	if base {
		switch index {
		case item.SetEffectIncStrength:
			p.addStrength += value
		case item.SetEffectIncAgility:
			p.addDexterity += value
		case item.SetEffectIncEnergy:
			p.addEnergy += value
		case item.SetEffectIncVitality:
			p.addVitality += value
		case item.SetEffectIncLeadership:
			p.addLeadership += value
		}
	} else {
		switch index {
		case item.SetEffectIncAttackMin:
			p.increaseAttackMin += value
		case item.SetEffectIncAttackMax:
			p.increaseAttackMax += value
		case item.SetEffectIncMagicAttack:
			p.magicAttackMin += value
			p.magicAttackMax += value
		case item.SetEffectIncDamage:
			p.setAddDamage += value
		case item.SetEffectIncAttackRate:
			p.AttackRate += value
		case item.SetEffectIncDefense:
			p.Defense += value
		case item.SetEffectIncMaxHP:
			p.MaxHP += value
		case item.SetEffectIncMaxMP:
			p.MaxMP += value
		case item.SetEffectIncMaxAG:
			p.MaxAG += value
		case item.SetEffectIncAG:
			p.MaxAG += value
		case item.SetEffectIncCritiDamageRate:
			p.criticalAttackRate += value
		case item.SetEffectIncCritiDamage:
			p.criticalAttackDamage += value
		case item.SetEffectIncExcelDamageRate:
			p.excellentAttackRate += value
		case item.SetEffectIncExcelDamage:
			p.excellentAttackDamage += value
		case item.SetEffectIncSkillAttack:
			p.increaseSkillAttack += value
		case item.SetEffectDoubleDamage:
			p.doubleDamageRate += value
		case item.SetEffectIgnoreDefense:
			p.ignoreDefenseRate += value
		case item.SetEffectIncShieldDefense:
			p.MaxSD += value
		case item.SetEffectIncTwoHandSwordDamage:
			p.increaseTwoHandWeaponDamage += value
		}
	}
}

func (p *Player) calcSetItem(base bool) {
	type set struct {
		index int
		count int
	}
	var sets []set

	sameWeapon := 0
	sameRing := 0
	for i, it := range p.inventory.WearingItems() {
		if it == nil {
			continue
		}
		if it.Durability == 0 {
			continue
		}
		index := it.Set
		if index <= 0 {
			continue
		}
		if i == 0 {
			sameWeapon = index
		}
		if i == 1 && sameWeapon > 0 {
			continue
		}
		if i == 10 {
			sameRing = index
		}
		if i == 11 && sameRing > 0 {
			continue
		}
		ok := false
		for i := range sets {
			if sets[i].index == index {
				sets[i].count++
				ok = true
				break
			}
		}
		if !ok {
			sets = append(sets, set{index, 1})
		}
	}

	for _, v := range sets {
		index := v.index
		count := v.count
		if count >= 2 {
			for i := 0; i < count-1; i++ {
				setEffect := item.SetManager.GetSet(index, i)
				p.addSetEffect(setEffect.Index, setEffect.Value, base)
			}

			if count > item.SetManager.GetSetEffectCount(index) {
				p.setFull = true
				setEffects := item.SetManager.GetSetFull(index)
				for _, setEffect := range setEffects {
					p.addSetEffect(setEffect.Index, setEffect.Value, base)
				}
			}
		}
	}
}

func (player *Player) calc380Item() {
	for _, wItem := range player.inventory.WearingItems() {
		if wItem.Durability == 0 {
			continue
		}
		if !wItem.Option380 || !item.Item380Manager.Is380Item(wItem.Section, wItem.Index) {
			continue
		}
		item.Item380Manager.Apply380ItemEffect(wItem.Section, wItem.Index, &player.item380Effect)
	}
	player.MaxHP += player.item380Effect.Item380EffectIncMaxHP
	player.MaxSD += player.item380Effect.Item380EffectIncMaxSD
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

func (p *Player) GetChangeUp() int {
	return p.changeUp
}

func (p *Player) pushChangedEquipment(it *item.Item) {
	reply := model.MsgChangeEquipmentReply{
		Number: p.Index,
		Item:   it,
	}
	p.PushViewport(&reply)
}

func (p *Player) limitUseItem(it *item.Item) bool {
	if p.Level < it.ReqLevel ||
		p.strength+p.addStrength < it.ReqStrength ||
		p.dexterity+p.addDexterity < it.ReqDexterity ||
		p.vitality+p.addVitality < it.ReqVitality ||
		p.energy+p.addEnergy < it.ReqEnergy ||
		p.leadership+p.addLeadership < it.ReqCommand {
		return true
	}
	reqClass := it.ReqClass[p.Class]
	if reqClass == 0 || (reqClass > p.changeUp+1) {
		return true
	}
	return false
}

func (p *Player) GetInventory() *item.Inventory {
	return &p.inventory
}

func (p *Player) GetInventoryItem(position int) *item.Item {
	return p.inventory.Item(position)
}

func (p *Player) GetWarehouse() *item.Warehouse {
	return &p.warehouse
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
	maps.MapManager.ClearMapAttrStand(p.MapNumber, p.TX, p.TY)
	defer p.Push(&model.MsgLogoutReply{Flag: msg.Flag})
	switch msg.Flag {
	case 0: // close game
	case 1: // back to pick character
		// offline to login state.\
		p.ConnectState = object.ConnectStateLogged
		p.saveCharacter()
		p.Reset()
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
		p.SpawnPosition()
	}
	maps.MapManager.SetMapAttrStand(p.MapNumber, p.TX, p.TY)
	p.CreateFrustum()
	p.MoveSpeed = 1000
	p.MaxRegenTime = 4 * time.Second
	p.ConnectState = object.ConnectStatePlaying
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
	p.Push(&model.MsgCreateViewportPlayerReply{
		Players: []*model.CreateViewportPlayer{
			{
				Index:                  p.Index,
				X:                      p.X,
				Y:                      p.Y,
				Class:                  p.Class,
				ChangeUp:               p.changeUp,
				Inventory:              [9]*item.Item(p.GetInventory().Items[:9]),
				Name:                   p.Name,
				TX:                     p.TX,
				TY:                     p.TY,
				Dir:                    p.Dir,
				PKLevel:                p.pkLevel,
				PentagramMainAttribute: p.PentagramAttributePattern,
				MuunItem:               -1,
				MuunSubItem:            -1,
				MuunRideItem:           -1,
				Level:                  p.Level,
				MaxHP:                  p.MaxHP,
				HP:                     p.HP,
				ServerCode:             0,
			},
		},
	})

	// client will calculate character after receiving inventory msg and master msg
	// calculate
	p.calc()

	// go func() {
	// 	time.Sleep(100 * time.Millisecond) // get character info
	// 	p.actioner.PlayerAction(p.index, "SetCharacter", &model.MsgSetCharacter{Name: msg.Name})
	// }()
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
	if conf.ServerEnv.Debug {
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
		p.warehouse = account.Warehouse
		p.warehouseMoney = account.WarehouseMoney
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
	case object.NpcTypeGoldarcher:
	case object.NpcTypePentagramMix:
	}
}

func (p *Player) CloseTalkWindow(msg *model.MsgEmpty) {
	p.TargetNumber = -1
}

func (p *Player) CloseWarehouseWindow(msg *model.MsgEmpty) {
	account := model.Account{
		ID:             p.accountID,
		Warehouse:      p.warehouse,
		WarehouseMoney: p.warehouseMoney,
	}
	err := model.DB.UpdateAccountWarehouse(&account)
	if err != nil {
		slog.Error("CloseWarehouseWindow model.DB.UpdateAccountWarehouse",
			"err", err, "player", p.Name)
		return
	}
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
		reply.MaxHPMP = p.MaxHP
	case 3:
		reply.MaxHPMP = p.MaxMP
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

func (p *Player) UsePet(msg *model.MsgUsePet) {
	reply := model.MsgUsePetReply{
		Position: msg.Position,
		Result:   0,
	}
	defer p.Push(&reply)

	if msg.Position < 12 || msg.Position > 204 {
		return
	}
	it := p.inventory.Items[msg.Position]
	if it == nil || it.Durability <= 0 {
		return
	}
	switch it.Code {
	case item.Code(13, 2), // Horn of Uniria 兽角
		item.Code(13, 3),  // Horn of Dinorant 彩云兽
		item.Code(13, 4),  // Dark Horse 黑王马之角
		item.Code(13, 37): // Horn of Fenrir 炎狼兽之角
		if p.pet != nil {
			it.HarmonyOption = 0
			p.pet = nil
			reply.Result = -1
		} else {
			it.HarmonyOption = 1
			p.pet = it
			reply.Result = -2
		}
		p.EquipmentChanged()
		p.pushChangedEquipment(it)
	}
}

func (p *Player) MuunSystem(msg *model.MsgMuunSystem)      {}
func (p *Player) StartPartyNumberPosition(*model.MsgEmpty) {}
func (p *Player) StopPartyNumberPosition(*model.MsgEmpty)  {}
