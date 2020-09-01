package skill

import (
	"sync"
)

type Skill struct {
	*SkillBase
	*SkillMasterBase
	Level     int
	DamageMin int
	DamageMax int
}

var poolSkill = sync.Pool{
	New: func() interface{} {
		return &Skill{}
	},
}

// Get get a skill from pool
func Get(skillIndex, level int) *Skill {
	skill := poolSkill.Get().(*Skill)
	skill.SkillBase = SkillTable[skillIndex]
	skill.Level = level
	skill.DamageMin = skill.SkillBase.Damage
	skill.DamageMax = skill.SkillBase.Damage + skill.SkillBase.Damage/2
	return skill
}

// Put put a skill to pool
func Put(skill *Skill) {
	skill.SkillBase = nil
	poolSkill.Put(skill)
}
