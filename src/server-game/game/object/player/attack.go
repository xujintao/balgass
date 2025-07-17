package player

import (
	"log/slog"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/exp"
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
	return p.ignoreDefenseRate
}

func (p *Player) GetCriticalAttackRate() int {
	return p.criticalAttackRate
}

func (p *Player) GetCriticalAttackDamage() int {
	return p.criticalAttackDamage
}

func (p *Player) GetExcellentAttackRate() int {
	return p.excellentAttackRate
}

func (p *Player) GetExcellentAttackDamage() int {
	return p.excellentAttackDamage
}

func (p *Player) GetMonsterDieGetHP() float64 {
	return p.monsterDieGetHP
}

func (p *Player) GetMonsterDieGetMP() float64 {
	return p.monsterDieGetMP
}

func (p *Player) GetAddDamage() int {
	return p.setAddDamage
}

func (p *Player) GetArmorReduceDamage() int {
	return p.armorReduceDamage
}

func (p *Player) GetWingIncreaseDamage() int {
	return p.wingIncreaseDamage
}

func (p *Player) GetWingReduceDamage() int {
	return p.wingReduceDamage
}

func (p *Player) GetHelperReduceDamage() int {
	return p.helperReduceDamage
}

func (p *Player) GetPetIncreaseDamage() int {
	return p.petIncreaseDamage
}

func (p *Player) GetPetReduceDamage() int {
	return p.petReduceDamage
}

func (p *Player) GetDoubleDamageRate() int {
	return p.doubleDamageRate
}

func (p *Player) GetMonsterDieGetMoney() float64 {
	return p.monsterDieGetMoney
}

func (p *Player) GetKnightGladiatorCalcSkillBonus() float64 {
	return p.knightGladiatorCalcSkillBonus
}

func (p *Player) GetImpaleSkillCalc() float64 {
	return p.impaleSkillCalc
}

func (p *Player) Die(tobj *object.Object, damage int) {
	// drop experience
	p.DieDropExperience()
	// delay drop item
	p.AddDelayMsg(1, 0, 800, tobj.Index)
}

func (p *Player) DieDropItem(*object.Object) {
	slog.Debug("player DieDropItem placeholder")
}

func (p *Player) LevelUp(addexp int) bool {
	if !p.IsMasterLevel() {
		p.experience += addexp
		levelUpExp := exp.ExperienceTable[p.Level]
		if p.experience < levelUpExp {
			return false
		}
		p.experience = levelUpExp
		p.Level++
		switch class.Class(p.Class) {
		case class.Magumsa,
			class.DarkLord,
			class.RageFighter,
			class.GrowLancer:
			p.levelPoint += conf.CommonServer.GameServerInfo.LevelPoint7
		default:
			p.levelPoint += conf.CommonServer.GameServerInfo.LevelPoint5
		}
	} else {
		p.masterExperience += addexp
		levelUpExp := exp.MasterExperienceTable[p.masterLevel]
		if p.masterExperience < levelUpExp {
			return false
		}
		p.masterExperience = levelUpExp
		p.masterLevel++
		p.masterPoint += conf.Common.General.MasterPointPerLevel
	}
	p.calc()
	p.HP = p.MaxHP
	p.SD = p.MaxSD
	p.MP = p.MaxMP
	p.AG = p.MaxAG
	p.PushHPSD(p.HP, p.SD)
	p.PushMPAG(p.MP, p.AG)
	if !p.IsMasterLevel() {
		reply := model.MsgLevelUpReply{
			Level:      p.Level,
			LevelPoint: p.levelPoint,
			MaxHP:      p.MaxHP,
			MaxMP:      p.MaxMP,
			MaxSD:      p.MaxSD,
			MaxAG:      p.MaxAG,
		}
		p.Push(&reply)
	} else {
		reply := model.MsgMasterLevelUpReply{
			MasterLevel:         p.masterLevel,
			MasterPointPerLevel: conf.Common.General.MasterPointPerLevel,
			MasterPoint:         p.masterPoint,
			MaxMasterLevel:      conf.Common.General.MaxLevelMaster,
			MaxHP:               p.MaxHP,
			MaxMP:               p.MaxMP,
			MaxSD:               p.MaxSD,
			MaxAG:               p.MaxAG,
		}
		p.Push(&reply)
	}
	return true
}
