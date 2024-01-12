package monster

import (
	"log"
	"math/rand"

	"github.com/xujintao/balgass/src/server_game/conf"
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
	excellentDropRate := conf.CommonServer.GameServerInfo.ItemExcelDropPercent
	// roll excellent item
	if rand.Intn(10000) < excellentDropRate {
		// excellent item
		log.Printf("MonsterDieGiveItem %s\n", "excellent item")
	} else {
		plainDropRate := conf.CommonServer.GameServerInfo.ItemDropPercent
		// roll normal item
		itemDropRate := m.ItemDropRate
		if itemDropRate < 1 {
			itemDropRate = 1
		}
		if rand.Intn(itemDropRate) < plainDropRate {
			// plain item
			log.Printf("MonsterDieGiveItem %s\n", "plain item")
		} else {
			// roll money
			moneyDropRate := m.MoneyDropRate
			if moneyDropRate < 1 {
				moneyDropRate = 1
			}
			if rand.Intn(moneyDropRate) < 10 {
				// money
				log.Printf("MonsterDieGiveItem %s\n", "money")
			}
		}
	}
}

func (*Monster) MonsterDieRecoverHP() {}
