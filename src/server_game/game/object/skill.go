package object

import (
	"log"

	"github.com/xujintao/balgass/src/server_game/game/skill"
)

func (obj *object) initSkill() {
	obj.skills = make(skill.Skills)
}

func (obj *object) learnSkill(index, level int) (*skill.Skill, bool) {
	if _, ok := obj.skills[index]; ok {
		log.Printf("[object]%s [skill]%d already exists", obj.Name, index)
		return nil, false
	}
	// obj.skills[index] = skill.SkillManager.Get(index, level, obj.skills)
	return obj.skills.Get(index, level)
}

func (obj *object) forgetSkill(index int) (*skill.Skill, bool) {
	if _, ok := obj.skills[index]; !ok {
		log.Printf("[object]%s [skill]%d doesn't exist", obj.Name, index)
		return nil, false
	}
	return obj.skills.Put(index)
}

func (obj *object) clearSkill() {
	obj.skills = nil
}
