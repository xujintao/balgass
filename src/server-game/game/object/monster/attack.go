package monster

import (
	"math/rand"

	"github.com/xujintao/balgass/src/server-game/game/drop"
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

func (m *Monster) GetAttackSpeedForDelay() int {
	return m.AttackSpeed
}

func (*Monster) GetDoubleDamageRate() int {
	return 0
}

func (*Monster) GetReturnDamageRate() int {
	return 0
}

func (*Monster) GetFullHPRecoveryRate() float64 {
	return 0
}

func (*Monster) GetFullMPRecoveryRate() float64 {
	return 0
}

func (*Monster) GetFullSDRecoveryRate() float64 {
	return 0
}

func (*Monster) GetMaceStunRate() float64 {
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

func (m *Monster) OnDamageTaken() {
	m.PassiveUntilHit = false
}

func (*Monster) GetMagicAttackMin() int {
	return 0
}

func (*Monster) GetMagicAttackMax() int {
	return 0
}

func (*Monster) GetCurseAttackMin() int {
	return 0
}

func (*Monster) GetCurseAttackMax() int {
	return 0
}

func (*Monster) GetCurse() int {
	return 0
}

func (*Monster) GetActivePetCode() int {
	return -1
}

func (*Monster) GetStrength() int {
	return 0
}

func (*Monster) GetDexterity() int {
	return 0
}

func (*Monster) GetEnergy() int {
	return 0
}

func (*Monster) GetVitality() int {
	return 0
}

func (*Monster) GetLeadership() int {
	return 0
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
	result := drop.DropManager.Drop(drop.Request{
		Trigger:       drop.TriggerMonsterDeath,
		MonsterClass:  m.Class,
		MonsterLevel:  m.Level,
		ItemDropRate:  m.ItemDropRate,
		MoneyDropRate: m.MoneyDropRate,
		Money:         m.MoneyDrop,
		MapNumber:     m.MapNumber,
		X:             m.X,
		Y:             m.Y,
		Player: drop.PlayerContext{
			Class:       tobj.Class,
			ChangeUp:    tobj.GetChangeUp(),
			Level:       tobj.Level,
			MasterLevel: tobj.GetMasterLevel(),
			MapNumber:   tobj.MapNumber,
			ZenBonus:    tobj.GetMonsterDieGetMoney(),
		},
	})
	for _, reward := range result.Rewards {
		it := reward.Item
		if it == nil && reward.Zen > 0 {
			it = item.NewItem(14, 15)
			it.Durability = reward.Zen
		}
		if it == nil {
			continue
		}
		mapNumber, x, y := m.MapNumber, m.X, m.Y
		switch reward.Placement {
		case drop.PlaceNearSource:
			x, y = itemBagDropPosition(mapNumber, x, y)
		case drop.PlaceAtPlayer:
			mapNumber, x, y = tobj.MapNumber, tobj.X, tobj.Y
		}
		maps.MapManager.AddItem(mapNumber, x, y, it)
	}
}

func itemBagDropPosition(mapNumber, x, y int) (int, int) {
	for i := 0; i < 10; i++ {
		dx, dy := rand.Intn(5), rand.Intn(5)
		if rand.Intn(2) == 0 {
			dx = -dx
		}
		if rand.Intn(2) == 0 {
			dy = -dy
		}
		dropX, dropY := x+dx, y+dy
		if maps.MapManager.GetMapAttr(mapNumber, dropX, dropY)&(4|8) == 0 {
			return dropX, dropY
		}
	}
	return x, y
}
