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

func (p *Player) GetMonsterDieGetHP() float64 {
	return p.MonsterDieGetHP
}

func (p *Player) GetMonsterDieGetMP() float64 {
	return p.MonsterDieGetMP
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
		p.Experience += addexp
		levelUpExp := exp.ExperienceTable[p.Level]
		if p.Experience < levelUpExp {
			return false
		}
		p.Experience = levelUpExp
		p.Level++
		switch class.Class(p.Class) {
		case class.Magumsa,
			class.DarkLord,
			class.RageFighter,
			class.GrowLancer:
			p.LevelPoint += conf.CommonServer.GameServerInfo.LevelPoint7
		default:
			p.LevelPoint += conf.CommonServer.GameServerInfo.LevelPoint5
		}
	} else {
		p.MasterExperience += addexp
		levelUpExp := exp.MasterExperienceTable[p.MasterLevel]
		if p.MasterExperience < levelUpExp {
			return false
		}
		p.MasterExperience = levelUpExp
		p.MasterLevel++
		p.MasterPoint += conf.Common.General.MasterPointPerLevel
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
			LevelPoint: p.LevelPoint,
			MaxHP:      p.MaxHP,
			MaxMP:      p.MaxMP,
			MaxSD:      p.MaxSD,
			MaxAG:      p.MaxAG,
		}
		p.Push(&reply)
	} else {
		reply := model.MsgMasterLevelUpReply{
			MasterLevel:         p.MasterLevel,
			MasterPointPerLevel: conf.Common.General.MasterPointPerLevel,
			MasterPoint:         p.MasterPoint,
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
