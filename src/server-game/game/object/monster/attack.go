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
			excel := item.ExcellentDropManager.DropExcellent(it.KindA, it.KindB)
			it.DecodeExcellent(excel)
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
		if it.KindA == item.KindAWing {
			it.RandWingAdditionKind()
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
