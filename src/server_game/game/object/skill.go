package object

import (
	"log"

	"github.com/xujintao/balgass/src/server_game/game/skill"
)

func (obj *object) initSkill() {
	obj.skills = make(map[int]*skill.Skill)
}

// AddSkill  object add skill
func (obj *object) addSkill(index, level int) bool {
	if _, ok := obj.skills[index]; ok {
		log.Printf("[object]%s [skill]%d already exists", obj.Name, index)
		return false
	}
	// obj.skills[index] = skill.SkillManager.Get(index, level, obj.skills)
	obj.skills.Get(index, level)
	return true
}

func (obj *object) clearSkill() {
	obj.skills = nil
}
