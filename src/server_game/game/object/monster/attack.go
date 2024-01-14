package monster

import (
	"log"
	"math/rand"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/item"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/object"
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

func (*Monster) GetKnightGladiatorCalcSkillBonus() float64 {
	return 1.0
}

func (*Monster) GetImpaleSkillCalc() float64 {
	return 1.0
}

func (m *Monster) Die(obj *object.Object) {
	// give experience
	obj.MonsterDieGetExperience(&m.Object)
	// give item
	m.AddDelayMsg(1, 0, 800, obj.Index)
	// obj delay recover hp/mp/sd
	obj.AddDelayMsg(2, 0, 2000, m.Index)
}

func (*Monster) MonsterDieGetExperience(*object.Object) {}

func (m *Monster) MonsterDieGiveItem(id int) {
	dropExcellentItem := false
	dropPlainItem := false
	dropMoney := false
	// roll excellent item
	excellentDropRate := conf.CommonServer.GameServerInfo.ItemExcelDropPercent
	if rand.Intn(10000) < excellentDropRate {
		// excellent item
		dropExcellentItem = true
		log.Printf("MonsterDieGiveItem %s\n", "excellent item")
	} else {
		// roll normal item
		plainDropRate := conf.CommonServer.GameServerInfo.ItemDropPercent
		itemDropRate := m.ItemDropRate
		if itemDropRate < 1 {
			itemDropRate = 1
		}
		if rand.Intn(itemDropRate) < plainDropRate {
			// plain item
			log.Printf("MonsterDieGiveItem %s\n", "plain item")
			dropPlainItem = true
		} else {
			// roll money
			moneyDropRate := m.MoneyDropRate
			if moneyDropRate < 1 {
				moneyDropRate = 1
			}
			if rand.Intn(moneyDropRate) < 10 {
				// money
				log.Printf("MonsterDieGiveItem %s\n", "money")
				dropMoney = true
			}
		}
	}

	if dropExcellentItem || dropPlainItem || dropMoney {
		var it *item.Item
		if dropExcellentItem || dropPlainItem {
			switch {
			case dropExcellentItem:
				// it = DropManager.DropExcellentItem(m.Level - 25)
				// // option
				// it.Calc()
			case dropPlainItem:
				it = DropManager.DropItem(m.Level)
			}
			if it == nil {
				return
			}
			it.Durability = it.MaxDurability
			if it.Type == item.TypeRegular {
				skillRate := 0
				luckyRate := 0
				switch {
				case dropExcellentItem:
					skillRate = conf.CommonServer.GameServerInfo.ItemExcelSkillDropPercent
					luckyRate = conf.CommonServer.GameServerInfo.ItemExcelSkillDropPercent
				case dropPlainItem:
					skillRate = conf.CommonServer.GameServerInfo.ItemSkillDropPercent
					luckyRate = conf.CommonServer.GameServerInfo.ItemLuckyDropPercent
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
			it.Durability = 2000
		}
		if it != nil {
			maps.MapManager.PushItem(m.MapNumber, m.X, m.Y, it.Copy())
		}
	}
}

func (*Monster) MonsterDieRecoverHP() {}
