package player

import (
	"github.com/xujintao/balgass/src/server-game/game/exp"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object"
)

func (p *Player) GetAttackRatePVP() int {
	return p.attackRatePVP
}

func (p *Player) GetDefenseRatePVP() int {
	return p.defenseRatePVP
}

func (p *Player) GetIgnoreDefenseRate() int {
	return p.IgnoreDefenseRate
}

func (p *Player) GetCriticalAttackRate() int {
	return p.CriticalAttackRate
}

func (p *Player) GetCriticalAttackDamage() int {
	return p.CriticalAttackDamage
}

func (p *Player) GetExcellentAttackRate() int {
	return p.ExcellentAttackRate
}

func (p *Player) GetExcellentAttackDamage() int {
	return p.ExcellentAttackDamage
}

func (p *Player) GetAddDamage() int {
	return p.SetAddDamage
}

func (p *Player) GetArmorReduceDamage() int {
	return p.ArmorReduceDamage
}

func (p *Player) GetWingIncreaseDamage() int {
	return p.WingIncreaseDamage
}

func (p *Player) GetWingReduceDamage() int {
	return p.WingReduceDamage
}

func (p *Player) GetHelperReduceDamage() int {
	return p.HelperReduceDamage
}

func (p *Player) GetPetIncreaseDamage() int {
	return p.PetIncreaseDamage
}

func (p *Player) GetPetReduceDamage() int {
	return p.PetReduceDamage
}

func (p *Player) GetDoubleDamageRate() int {
	return p.DoubleDamageRate
}

func (p *Player) GetMonsterDieGetMoney() float64 {
	return p.MonsterDieGetMoney
}

func (p *Player) GetKnightGladiatorCalcSkillBonus() float64 {
	return p.KnightGladiatorCalcSkillBonus
}

func (p *Player) GetImpaleSkillCalc() float64 {
	return p.ImpaleSkillCalc
}

func (p *Player) Die(obj *object.Object, damage int) {

}

func (p *Player) MonsterDieGetExperience(tobj *object.Object, damage int) {
	level := p.Level + p.MasterLevel
	targetLevel := (tobj.Level + 25) * tobj.Level / 3
	if tobj.Level+10 < level {
		targetLevel = targetLevel * (tobj.Level + 10) / level
	}
	if tobj.Level >= 65 {
		targetLevel += (tobj.Level - 64) * tobj.Level / 4
	}
	addexp := 0
	maxexp := 0
	if targetLevel > 0 {
		maxexp = targetLevel / 2
	} else {
		targetLevel = 0
	}
	if maxexp < 1 {
		addexp = targetLevel
	} else {
		addexp = maxexp/2 + targetLevel
	}
	if addexp <= 0 {
		return
	}
	var mapBonus float64
	var baseBonus float64
	if !p.isMasterLevel() {
		mapBonus = maps.MapManager.GetExpBonus(p.MapNumber)
		baseBonus = exp.ExpManager.Normal
	} else {
		mapBonus = maps.MapManager.GetMasterExpBonus(p.MapNumber)
		baseBonus = exp.ExpManager.Master
	}
	addexp = int(float64(addexp) * (1 + mapBonus) * baseBonus)
	if !p.LevelUp(addexp) {
		reply := model.MsgExperienceReply{
			Number:     tobj.Index,
			Experience: addexp,
			Damage:     damage,
		}
		p.Push(&reply)
	}
}

func (p *Player) MonsterDieDropItem(*object.Object) {}

func (p *Player) MonsterDieRecoverHP() {
	if p.MonsterDieGetHP != 0 {
		p.HP += int(float64(p.MaxHP) * p.MonsterDieGetHP)
		if p.HP >= p.MaxHP {
			p.HP = p.MaxHP
		}
		p.pushHP(p.HP, p.SD)
	}
	if p.MonsterDieGetMP != 0 {
		p.MP += int(float64(p.MaxMP) * p.MonsterDieGetMP)
		if p.MP >= p.MaxMP {
			p.MP = p.MaxMP
		}
		p.PushMPAG(p.MP, p.AG)
	}
}
