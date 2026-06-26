package object

import "time"

type skillEffect struct {
	attack  int
	defense int
	maxHP   int
	expire  time.Time
}

func (obj *Object) initSkillEffect() {
	obj.skillEffects = make(map[int]*skillEffect)
}

func (obj *Object) clearSkillEffect() {
	obj.skillEffects = nil
}

func (obj *Object) addSkillEffect(index int, effect *skillEffect) {
	if obj.skillEffects[index] != nil {
		obj.removeSkillEffect(index)
	}
	obj.skillEffects[index] = effect
	obj.AttackMin += effect.attack
	obj.AttackMax += effect.attack
	obj.Defense += effect.defense
	obj.MaxHP += effect.maxHP
}

func (obj *Object) removeSkillEffect(index int) {
	effect := obj.skillEffects[index]
	if effect == nil {
		return
	}
	delete(obj.skillEffects, index)
	obj.AttackMin -= effect.attack
	obj.AttackMax -= effect.attack
	obj.Defense -= effect.defense
	obj.MaxHP -= effect.maxHP
	if obj.HP > obj.MaxHP {
		obj.HP = obj.MaxHP
	}
	if effect.maxHP != 0 {
		obj.PushMaxHPSD(obj.MaxHP, obj.MaxSD)
		obj.PushHPSD(obj.HP, obj.SD)
	}
}

func (obj *Object) processSkillEffect() {
	now := time.Now()
	for index, effect := range obj.skillEffects {
		if now.Before(effect.expire) {
			continue
		}
		obj.removeSkillEffect(index)
	}
}
