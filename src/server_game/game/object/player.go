package object

import (
	"log"
	"time"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/guild"
	"github.com/xujintao/balgass/src/server_game/game/item"
	"github.com/xujintao/balgass/src/server_game/game/model"
	"github.com/xujintao/balgass/src/server_game/game/skill"
)

type Pusher interface {
	Push(c1c2.ConnWriter, interface{})
}

type Player struct {
	Object
	Addr                          string
	Conn                          c1c2.ConnWriter
	pusher                        Pusher
	AccountID                     string
	AuthLevel                     int
	hwid                          string
	Experience                    uint
	ExperienceNext                uint
	ExperienceMaster              uint
	ExperienceMasterNext          uint
	masterLevel                   int
	LevelUpPoint                  int
	MasterPoint                   int
	MasterPointUsed               int
	FruitPoint                    int
	Money                         int
	Strength                      int
	Dexterity                     int
	Vitality                      int
	Energy                        int
	dbClass                       uint8
	ChangeUP                      int // 0=1转 1=2转 2=3转
	guild                         *guild.GuildInfo
	guildName                     string
	guildStatus                   int
	guildUnionTimeStamp           int
	guildNumber                   int
	lastMoveTime                  time.Time
	resets                        int
	vipType                       uint8
	vipEffect                     uint8
	santaCount                    uint8
	goblinTime                    time.Time
	securityCheck                 bool
	securityCode                  int
	RegisterLMS                   uint8
	registerLMSRoom               uint8
	jewelHarmonyEffect            item.JewelHarmonyItemEffect
	item380Effect                 item.Item380Effect
	kanturuEntranceByNPC          bool
	gensInfoLoad                  bool
	questInfoLoad                 bool
	wCoinP                        int
	wCoinC                        int
	goblinPoint                   int
	periodItemEffectIndex         int
	seedOptionList                [35]item.SocketOptionList
	bonusOptionList               [7]item.SocketOptionList
	setOptionList                 [2]item.SocketOptionList
	refillHPSocketOption          uint16
	refillMPSocketOption          uint16
	socketOptionMonsterDieGetHP   uint16
	socketOptionMonsterDieGetMana uint16
	AGReduceRate                  uint8
	muBotEnable                   bool
	muBotTotalTime                time.Duration
	muBotPayTime                  time.Duration
	muBotTick                     time.Time
	InventoryExpansion            int
	WarehouseExpansion            int
	LastAuthTime                  time.Time
	LastXorKey1                   [4]int
	LaskXorKey2                   [4]int
	bot                           bool
	botIndex                      int
	skillHellFire2State           int
	skillHellFire2Count           int
	skillHellFire2Time            time.Time
	skillStrengthenHellFire2State int
	skillStrengthenHellFire2Count int
	skillStrengthenHellFire2Time  time.Time
	reqWarehouseOpen              int
	// set
	setEffectIncSkillAttack        int
	setEffectIncExcelDamage        int
	setEffectIncExcelDamageRate    int
	setEffectIncCritiDamage        int
	setEffectIncCritiDamageRate    int
	setEffectIncAG                 int
	setEffectIncDamage             int
	setEffectIncAttackMin          int
	setEffectIncAttackMax          int
	setEffectIncAttack             int
	setEffectIncDefense            int
	setEffectIncDefenseRate        int
	setEffectIncMagicAttack        int
	setEffectIgnoreDefense         int
	setEffectDoubleDamage          int
	setEffectTwoHandSwordIncDamage int
	setEffectIncAttackRate         int
	setEffectReflectDamage         int
	setEffectIncShieldDefense      int
	setEffectDecAG                 int
	setEffectIncItemDropRate       int
	setFull                        bool
	// excel wing
	excelWingEffectIgnoreDefense int
	excelWingEffectReboundDamage int
	excelWingEffectRecoveryHP    int
	excelWingEffectRecoveryMP    int
	excelWingEffectDoubleDamage  int
}

func (player *Player) MasterLevel() bool {
	return player.ChangeUP == 2 && player.Level >= conf.Common.General.MaxLevelNormal
}

func (player *Player) addExcelCommonEffect(opt *item.ExcelCommon, wItem *item.Item, position int) {
	id := opt.ID
	value := opt.Value
	switch id {
	case item.ExcelCommonIncMPMonsterDie: // 杀怪回蓝
		player.MonsterDieGetMana++
	case item.ExcelCommonIncHPMonsterDie: // 杀怪回红
		player.MonsterDieGetLife++
	case item.ExcelCommonIncAttackSpeed: // 攻速
		player.attackSpeed += value
		player.magicSpeed += value
	case item.ExcelCommonIncAttackPercent: // 2%
		if wItem.Section == 5 || // 法杖
			wItem.Code == item.Code(13, 12) || // 雷链子
			wItem.Code == item.Code(13, 25) || // 冰链子
			wItem.Code == item.Code(13, 27) { // 水链子
			player.magicDamageMin += player.magicDamageMin * value / 100
			player.magicDamageMax += player.magicDamageMax * value / 100
		} else {
			if position == 0 || position == 9 {
				player.attackDamageLeftMin += player.attackDamageLeftMin * value / 100
				player.attackDamageLeftMax += player.attackDamageLeftMax * value / 100
			}
			if position == 1 || position == 9 {
				player.attackDamageRightMin += player.attackDamageRightMin * value / 100
				player.attackDamageRightMax += player.attackDamageRightMax * value / 100
			}
		}
	case item.ExcelCommonIncAttackLevel: // =20
		if wItem.Section == 5 || // 法杖
			wItem.Code == item.Code(13, 12) || // 雷链子
			wItem.Code == item.Code(13, 25) || // 冰链子
			wItem.Code == item.Code(13, 27) { // 水链子
			player.magicDamageMin += (player.Level + player.masterLevel) / value
			player.magicDamageMax += (player.Level + player.masterLevel) / value
		} else {
			if position == 0 || position == 9 {
				player.attackDamageLeftMin += (player.Level + player.masterLevel) / value
				player.attackDamageLeftMax += (player.Level + player.masterLevel) / value
			}
			if position == 1 || position == 9 {
				player.attackDamageRightMin += (player.Level + player.masterLevel) / value
				player.attackDamageRightMax += (player.Level + player.masterLevel) / value
			}
		}
	case item.ExcelCommonIncExcelDamage: // 一击
		player.excellentDamage += value
	case item.ExcelCommonIncZen: // 加钱
		player.MonsterDieGetMoney += value
	case item.ExcelCommonIncDefenseRate: // f10
		player.successfulBlocking += player.successfulBlocking * value / 100
	case item.ExcelCommonReflectDamage: // 反伤
		player.DamageReflect += value
	case item.ExcelCommonDecDamage: // 减伤
		player.DamageMinus += value
	case item.ExcelCommonIncMaxMP: // 加魔
		player.AddMP += player.MaxMP * value / 100
	case item.ExcelCommonIncMaxHP: // 加生
		player.AddHP += player.MaxHP * value / 100
	}
}

func (player *Player) addExcelWingEffect(opt *item.ExcelWing, wItem *item.Item) {
	id := opt.ID
	value := opt.Value
	switch id {
	case 0, 9, 13: // incHP
		player.AddHP += value + wItem.Level*5
	case 1, 10, 14: // incMP
		player.AddMP += value + wItem.Level*5
	case 2, 5, 11, 15, 16, 18, 23: // ingore
		player.excelWingEffectIgnoreDefense = value
	// case 3: // AG ?
	// 	player.AddAG += value
	// case 4: // speed ?
	// 	player.attackSpeed += value
	// 	player.magicSpeed += value
	case 6, 19:
		player.excelWingEffectReboundDamage = value
	case 7, 17, 20, 24:
		player.excelWingEffectRecoveryHP = value
	case 8, 21:
		player.excelWingEffectRecoveryMP = value
	case 12:
		player.AddLeadership += value + wItem.Level*5
	case 22:
		player.excelWingEffectDoubleDamage = value
	}
}

func (player *Player) addSetEffect(index item.SetEffectType, value int) {
	switch index {
	case item.SetEffectIncStrength:
		player.AddStrength += value
	case item.SetEffectIncAgility:
		player.AddDexterity += value
	case item.SetEffectIncEnergy:
		player.AddEnergy += value
	case item.SetEffectIncVitality:
		player.AddVitality += value
	case item.SetEffectIncLeadership:
		player.AddLeadership += value
	case item.SetEffectIncMaxHP:
		player.AddHP += value
	case item.SetEffectIncMaxMP:
		player.AddMP += value
	case item.SetEffectIncMaxAG:
		player.AddAG += value
	case item.SetEffectDoubleDamage:
		player.setEffectDoubleDamage += value
	case item.SetEffectIncShieldDefense:
		player.setEffectIncShieldDefense += value
	case item.SetEffectIncTwoHandSwordDamage:
		player.setEffectTwoHandSwordIncDamage += value
	case item.SetEffectIncAttackMin:
		player.setEffectIncAttackMin += value
	case item.SetEffectIncAttackMax:
		player.setEffectIncAttackMax += value
	case item.SetEffectIncMagicAttack:
		player.setEffectIncMagicAttack += value
	case item.SetEffectIncDamage:
		player.setEffectIncDamage += value
	case item.SetEffectIncAttackRate:
		player.setEffectIncAttackRate += value
	case item.SetEffectIncDefense:
		player.setEffectIncDefense += value
	case item.SetEffectIgnoreDefense:
		player.setEffectIgnoreDefense += value
	case item.SetEffectIncAG:
		player.setEffectIncAG += value
	case item.SetEffectIncCritiDamage:
		player.setEffectIncCritiDamage += value
	case item.SetEffectIncCritiDamageRate:
		player.setEffectIncCritiDamageRate += value
	case item.SetEffectIncExcelDamage:
		player.setEffectIncExcelDamage += value
	case item.SetEffectIncExcelDamageRate:
		player.setEffectIncExcelDamageRate += value
	case item.SetEffectIncSkillAttack:
		player.setEffectIncSkillAttack += value
	}
}

func (player *Player) CalcExcelItem() {
	for i, wItem := range player.Inventory[0:InventoryWearSize] {
		if wItem.Durability == 0 {
			continue
		}
		if wItem.Excel == 0 {
			continue
		}
		if i == 7 {
			for _, opt := range item.ExcelManager.Wings.Options {
				if wItem.KindA == opt.ItemKindA && wItem.KindB == opt.ItemKindB {
					if wItem.Excel&opt.Number == opt.Number {
						player.addExcelWingEffect(opt, &wItem)
					}
				}
			}
		} else {
			for _, opt := range item.ExcelManager.Common.Options {
				switch wItem.KindA {
				case opt.ItemKindA1, opt.ItemKindA2, opt.ItemKindA3:
					if wItem.Excel&opt.Number == opt.Number {
						player.addExcelCommonEffect(opt, &wItem, i)
					}
				}
			}
		}
	}
}

func (player *Player) CalcSetItem() {
	type set struct {
		index int
		count int
	}
	var sets []set

	sameWeapon := 0
	sameRing := 0
	for i, wItem := range player.Inventory[0:InventoryWearSize] {
		if wItem.Durability == 0 {
			continue
		}
		tierIndex := wItem.GetSetTierIndex()
		if tierIndex == 0 {
			continue
		}
		index := item.SetManager.GetSetIndex(wItem.Section, wItem.Index, tierIndex)
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
				player.addSetEffect(setEffect.Index, setEffect.Value)
			}

			if count > item.SetManager.GetSetEffectCount(index) {
				player.setFull = true
				setEffects := item.SetManager.GetSetFull(index)
				for _, setEffect := range setEffects {
					player.addSetEffect(setEffect.Index, setEffect.Value)
				}
			}
		}
	}
}

func (player *Player) Calc380Item() {
	for _, wItem := range player.Inventory[0:InventoryWearSize] {
		if wItem.Durability == 0 {
			continue
		}
		if !wItem.Option380 || !item.Item380Manager.Is380Item(wItem.Section, wItem.Index) {
			continue
		}
		item.Item380Manager.Apply380ItemEffect(wItem.Section, wItem.Index, &player.item380Effect)
	}
	player.AddHP += player.item380Effect.Item380EffectIncMaxHP
	player.AddSD += player.item380Effect.Item380EffectIncMaxSD
}

func (player *Player) LimitUseItem(it *item.Item) bool {
	if player.Level < it.ReqLevel ||
		player.Strength+player.AddStrength < it.ReqStrength ||
		player.Dexterity+player.AddDexterity < it.ReqDexterity ||
		player.Vitality+player.AddVitality < it.ReqVitality ||
		player.Energy+player.AddEnergy < it.ReqEnergy ||
		player.Leadership+player.AddLeadership < it.ReqCommand {
		return true
	}
	reqClass := it.ReqClass[player.Class]
	if reqClass == 0 || (reqClass > player.ChangeUP+1) {
		return true
	}
	return false
}

// SkillLearn object learn skill from skill stone or master point
func (player *Player) SkillLearn(skillIndex int) bool {
	// validate skillIndex
	skillBase, ok := skill.SkillTable[skillIndex]
	if !ok {
		log.Printf("player[%s] learn invalid skill index[%d]", player.Name, skillIndex)
		return false
	}
	if skillBase.STID == 0 && skillBase.UseType == 0 {
		return player.SkillAdd(skillIndex, 0)
	}

	// validate player level
	if !player.MasterLevel() {
		return false // 2
	}

	// validate skill level
	level := 0
	if skill, ok := player.Skills[skillIndex]; ok {
		level = skill.Level
	}
	level += skillBase.ReqMLPoint
	skillMaster, _ := skill.SkillMasterTable[skillIndex]
	if level > skillMaster.MaxPoint {
		return false // 4
	}

	// validate master point
	if player.MasterPoint < skillBase.ReqMLPoint {
		return false // 4
	}

	// validate new skill
	if level == 1 {

	}

	return true
}

func (player *Player) PushSkillOne() {
	var msg model.MsgSkillList
	player.pusher.Push(player.Conn, &msg)
}

func (player *Player) PushSkillAll() {
	var msg model.MsgSkillList
	player.pusher.Push(player.Conn, &msg)
}
