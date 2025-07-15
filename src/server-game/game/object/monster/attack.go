package monster

import (
	"log/slog"
	"math/rand"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/object"
)

func (*Monster) GetAttackRatePVP() int {
	return 0
}

func (*Monster) GetDefenseRatePVP() int {
	return 0
}

func (*Monster) GetIgnoreDefenseRate() int {
	return 0
}

func (*Monster) GetCriticalAttackRate() int {
	return 0
}

func (*Monster) GetCriticalAttackDamage() int {
	return 0
}

func (*Monster) GetExcellentAttackRate() int {
	return 0
}

func (*Monster) GetExcellentAttackDamage() int {
	return 0
}

func (*Monster) GetMonsterDieGetHP() float64 {
	return 0
}

func (*Monster) GetMonsterDieGetMP() float64 {
	return 0
}

func (*Monster) GetAddDamage() int {
	return 0
}

func (*Monster) GetArmorReduceDamage() int {
	return 0
}

func (*Monster) GetWingIncreaseDamage() int {
	return 0
}

func (*Monster) GetWingReduceDamage() int {
	return 0
}

func (*Monster) GetHelperReduceDamage() int {
	return 0
}

func (*Monster) GetPetIncreaseDamage() int {
	return 0
}

func (*Monster) GetPetReduceDamage() int {
	return 0
}

func (*Monster) GetDoubleDamageRate() int {
	return 0
}

func (*Monster) GetMonsterDieGetMoney() float64 {
	return 0.0
}

func (*Monster) GetKnightGladiatorCalcSkillBonus() float64 {
	return 1.0
}

func (*Monster) GetImpaleSkillCalc() float64 {
	return 1.0
}

func (m *Monster) Die(tobj *object.Object, damage int) {
	// give experience to target
	m.DieGiveExperience(tobj, damage)
	// delay drop item
	m.AddDelayMsg(1, 0, 800, tobj.Index)
	// delay recover target hp/mp/sd
	m.AddDelayMsg(3, 0, 2000, tobj.Index)
}

func (m *Monster) DieDropItem(tobj *object.Object) {
	dropExcellentItem := false
	dropPlainItem := false
	dropMoney := false
	// roll excellent item
	excellentDropRate := conf.CommonServer.GameServerInfo.ExcelItemDropPercent
	if rand.Intn(10000) < excellentDropRate {
		// excellent item
		dropExcellentItem = true
		slog.Debug("DieDropItem", "item", "excellent", "monster", m.Name)
	} else {
		// roll plain item
		plainDropRate := conf.CommonServer.GameServerInfo.ItemDropPercent
		itemDropRate := m.ItemDropRate
		if itemDropRate < 1 {
			itemDropRate = 1
		}
		if rand.Intn(itemDropRate) < plainDropRate {
			// plain item
			slog.Debug("DieDropItem", "item", "plain", "monster", m.Name)
			dropPlainItem = true
		} else {
			// roll money
			moneyDropRate := m.MoneyDropRate
			if moneyDropRate < 1 {
				moneyDropRate = 1
			}
			if rand.Intn(moneyDropRate) < 10 {
				// money
				slog.Debug("DieDropItem", "item", "money", "monster", m.Name)
				dropMoney = true
			}
		}
	}

	var it *item.Item
	switch {
	case dropExcellentItem || dropPlainItem:
		switch {
		case dropExcellentItem:
			it = DropManager.DropItemExcellent(m.Level - 25)
			if it == nil {
				return
			}
			slog.Debug("DieDropItem", "item", "excellent", "monster", m.Name, "item", it.Name)
			drops := item.ExcellentDropManager.DropExcellent(it.KindA, it.KindB)
			for _, v := range drops {
				switch it.KindA {
				case item.KindAWeapon, item.KindAPendant:
					switch v {
					case 1:
						it.ExcellentAttackMP = true
					case 2:
						it.ExcellentAttackHP = true
					case 4:
						it.ExcellentAttackSpeed = true
					case 8:
						it.ExcellentAttackPercent = true
					case 16:
						it.ExcellentAttackLevel = true
					case 32:
						it.ExcellentAttackRate = true
					}
				case item.KindAArmor, item.KindARing:
					switch v {
					case 1:
						it.ExcellentDefenseMoney = true
					case 2:
						it.ExcellentDefenseRate = true
					case 4:
						it.ExcellentDefenseReflect = true
					case 8:
						it.ExcellentDefenseReduce = true
					case 16:
						it.ExcellentDefenseMP = true
					case 32:
						it.ExcellentDefenseHP = true
					}
				case item.KindAWing:
					switch it.KindB {
					case item.KindBWing2nd, item.KindBCapeLord, item.KindBCapeFighter:
						// roll wing addition
						wingAdd := rand.Intn(2)
						switch wingAdd {
						case 0:
							if it.Code == item.Code(12, 4) || // Wings of Soul 魔魂之翼
								it.Code == item.Code(12, 5) { // Wing of Despair 绝望之翼
								it.WingAdditionMagicAttack = true
							} else {
								it.WingAdditionAttack = true
							}
						case 1:
							it.WingAdditionRecoverHP = true
						}
						// add wing excellent options
						switch v {
						case 1:
							it.ExcellentWing2HP = true
						case 2:
							it.ExcellentWing2MP = true
						case 4:
							it.ExcellentWing2Ignore = true
						case 8:
							switch it.KindB {
							case item.KindBCapeLord:
								it.ExcellentWing2Leadership = true
							case item.KindBCapeFighter:
							default:
								it.ExcellentWing2AG = true
							}
						case 16:
							switch it.KindB {
							case item.KindBCapeLord:
							case item.KindBCapeFighter:
							default:
								it.ExcellentWing2Speed = true
							}
						}
					case item.KindBWing3rd:
						// roll wing addition
						wingAdd := rand.Intn(3)
						switch wingAdd {
						case 0:
							if it.Code == item.Code(12, 39) { // Wing of Ruin 破灭之翼
								it.WingAdditionMagicAttack = true
							} else if it.Code == item.Code(12, 43) { // Wing of Dimension 次元之翼
								it.WingAdditionCurseAttack = true
							} else {
								it.WingAdditionDefense = true
							}
						case 1:
							if it.Code == item.Code(12, 37) || // Wing of Eternal 时空之翼
								it.Code == item.Code(12, 43) { // Wing of Dimension 次元之翼
								it.WingAdditionMagicAttack = true
							} else {
								it.WingAdditionAttack = true
							}
						case 2:
							it.WingAdditionRecoverHP = true
						}
						// add wing excellent options
						switch v {
						case 1:
							it.ExcellentWing3Ignore = true
						case 2:
							it.ExcellentWing3Return = true
						case 4:
							it.ExcellentWing3HP = true
						case 8:
							it.ExcellentWing3MP = true
						}
					case item.KindBWingMonster:
						// roll wing addition
						wingAdd := rand.Intn(2)
						switch wingAdd {
						case 0:
							if it.Code == item.Code(12, 264) { // Wings of Magic 魔力之翼
								it.WingAdditionCurseAttack = true
							} else {
								it.WingAdditionRecoverHP = true
							}
						case 1:
							if it.Code == item.Code(12, 264) { // Wings of Magic 魔力之翼
								it.WingAdditionMagicAttack = true
							} else {
								it.WingAdditionAttack = true
							}
						}
						// add wing excellent options
						switch v {
						case 1:
							it.ExcellentWing25Ignore = true
						case 2:
							it.ExcellentWing25HP = true
						}
					}
				}
			}
		case dropPlainItem:
			it = DropManager.DropItem(m.Level)
			if it == nil {
				return
			}
			slog.Debug("DieDropItem", "item", "plain", "monster", m.Name, "item", it.Name)
		}
		if it.ItemBase.Durability <= 5 {
			it.Durability = it.ItemBase.Durability
		} else {
			it.Durability = rand.Intn(it.ItemBase.Durability)
		}
		skillRate := 0
		luckyRate := 0
		switch {
		case dropExcellentItem:
			skillRate = conf.CommonServer.GameServerInfo.ExcelItemSkillDropPercent
			luckyRate = conf.CommonServer.GameServerInfo.ExcelItemLuckyDropPercent
		case dropPlainItem:
			skillRate = conf.CommonServer.GameServerInfo.ItemSkillDropPercent
			luckyRate = conf.CommonServer.GameServerInfo.ItemLuckyDropPercent
		}
		if it.SkillIndex == 0 || it.Type == item.TypeCommon {
			skillRate = 0
		}
		if !(it.KindA == item.KindAWeapon ||
			it.KindA == item.KindAArmor ||
			it.KindA == item.KindAWing) {
			luckyRate = 0
		}
		if rand.Intn(100) < skillRate {
			it.Skill = true
		}
		if rand.Intn(100) < luckyRate {
			it.Lucky = true
		}
		addtionRate := rand.Intn(3)
		if it.Type == item.TypeCommon {
			addtionRate = 0
		}
		switch addtionRate {
		case 0:
			it.Addition = 0
		case 1:
			it.Addition = 4
		case 2:
			it.Addition = 8
		}
		it.Calc()
	case dropMoney:
		it = item.NewItem(14, 15)
		// money := maps.MapManager.GetZen(m.MapNumber, m.Level)
		money := m.MoneyDrop
		if money <= 0 {
			return
		}
		money = int(float64(money) * conf.Common.General.ZenDropMultiplier)
		money += int(float64(money) * tobj.GetMonsterDieGetMoney())
		it.Durability = money
	}
	if it != nil {
		maps.MapManager.AddItem(m.MapNumber, m.X, m.Y, it)
	}
}
