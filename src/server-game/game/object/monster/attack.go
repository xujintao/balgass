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

func (m *Monster) GetIgnoreDefenseRate() int {
	return 0
}

func (m *Monster) GetCriticalAttackRate() int {
	return 0
}

func (m *Monster) GetCriticalAttackDamage() int {
	return 0
}

func (m *Monster) GetExcellentAttackRate() int {
	return 0
}

func (m *Monster) GetExcellentAttackDamage() int {
	return 0
}

func (m *Monster) GetAddDamage() int {
	return 0
}

func (*Monster) GetArmorReduceDamage() int {
	return 0
}

func (m *Monster) GetWingIncreaseDamage() int {
	return 0
}

func (m *Monster) GetWingReduceDamage() int {
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

func (m *Monster) GetDoubleDamageRate() int {
	return 0
}

func (m *Monster) GetMonsterDieGetMoney() float64 {
	return 0.0
}

func (*Monster) GetKnightGladiatorCalcSkillBonus() float64 {
	return 1.0
}

func (*Monster) GetImpaleSkillCalc() float64 {
	return 1.0
}

func (m *Monster) Die(obj *object.Object, damage int) {
	// give experience
	obj.MonsterDieGetExperience(&m.Object, damage)
	// drop item
	m.AddDelayMsg(1, 0, 800, obj.Index)
	// obj delay recover hp/mp/sd
	obj.AddDelayMsg(2, 0, 2000, m.Index)
}

func (*Monster) MonsterDieGetExperience(*object.Object, int) {}

func (m *Monster) MonsterDieDropItem(tobj *object.Object) {
	dropExcellentItem := false
	dropPlainItem := false
	dropMoney := false
	// roll excellent item
	excellentDropRate := conf.CommonServer.GameServerInfo.ExcelItemDropPercent
	if rand.Intn(10000) < excellentDropRate {
		// excellent item
		dropExcellentItem = true
		slog.Debug("MonsterDieDropItem",
			"item", "excellent",
			"monster", m.Annotation)
	} else {
		// roll normal item
		plainDropRate := conf.CommonServer.GameServerInfo.ItemDropPercent
		itemDropRate := m.ItemDropRate
		if itemDropRate < 1 {
			itemDropRate = 1
		}
		if rand.Intn(itemDropRate) < plainDropRate {
			// plain item
			slog.Debug("MonsterDieDropItem",
				"item", "plain",
				"monster", m.Annotation)
			dropPlainItem = true
		} else {
			// roll money
			moneyDropRate := m.MoneyDropRate
			if moneyDropRate < 1 {
				moneyDropRate = 1
			}
			if rand.Intn(moneyDropRate) < 10 {
				// money
				slog.Debug("MonsterDieDropItem",
					"item", "money",
					"monster", m.Annotation)
				dropMoney = true
			}
		}
	}

	if dropExcellentItem || dropPlainItem || dropMoney {
		var it *item.Item
		if dropExcellentItem || dropPlainItem {
			switch {
			case dropExcellentItem:
				it = DropManager.DropExcellentItem(m.Level - 25)
				if it == nil {
					return
				}
				drops := item.ExcellentDropManager.DropRegularExcellent(int(it.KindA))
				for _, v := range drops {
					switch it.KindA {
					case 1, 2:
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
					case 3, 4:
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
					}
				}
				it.Calc()
			case dropPlainItem:
				it = DropManager.DropItem(m.Level)
				if it == nil {
					return
				}
			}
			it.Durability = it.MaxDurability
			if it.Type == item.TypeRegular {
				skillRate := 0
				luckyRate := 0
				switch {
				case dropExcellentItem:
					skillRate = conf.CommonServer.GameServerInfo.ExcelItemSkillDropPercent
					luckyRate = conf.CommonServer.GameServerInfo.ExcelItemSkillDropPercent
				case dropPlainItem:
					skillRate = conf.CommonServer.GameServerInfo.ItemSkillDropPercent
					luckyRate = conf.CommonServer.GameServerInfo.ItemLuckyDropPercent
				}
				if !(it.KindA == item.KindAWeapon || it.KindB == item.KindBShield) {
					skillRate = 0
				}
				if !(it.KindA == item.KindAWeapon || it.KindA == item.KindAArmor) {
					luckyRate = 0
				}
				if rand.Intn(100) < skillRate {
					it.Skill = true
				}
				if rand.Intn(100) < luckyRate {
					it.Lucky = true
				}
				addtionRate := rand.Intn(3)
				switch addtionRate {
				case 0:
					it.Addition = 0
				case 1:
					it.Addition = 4
				case 2:
					it.Addition = 8
				}
			}
		} else if dropMoney {
			it = item.NewItem(14, 15)
			money := maps.MapManager.GetZen(m.MapNumber, m.Level)
			if money <= 0 {
				return
			}
			money = int(float64(money) * conf.Common.General.ZenDropMultiplier)
			money += int(float64(money) * tobj.GetMonsterDieGetMoney())
			it.Durability = money
		}
		if it != nil {
			maps.MapManager.PushItem(m.MapNumber, m.X, m.Y, it.Copy())
		}
	}
}

func (*Monster) MonsterDieRecoverHP() {}
