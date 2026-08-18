package object

import (
	"time"

	"github.com/xujintao/balgass/src/server-game/game/effect"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

func (obj *Object) initEffect() {
	obj.effects = make(effect.Effects)
}

func (obj *Object) clearEffect() {
	obj.removeAllEffects()
	obj.effects = nil
}

func (obj *Object) removeAllEffects() {
	for _, eff := range obj.effects.RemoveAll() {
		obj.unapplyEffect(eff)
	}
}

func (obj *Object) addEffect(eff *effect.Effect) bool {
	removed, ok := obj.effects.Add(eff)
	for _, old := range removed {
		obj.unapplyEffect(old)
	}
	if !ok {
		return false
	}
	obj.applyEffect(eff)
	return true
}

func (obj *Object) applyEffect(eff *effect.Effect) {
	obj.applyEffectValue(eff, 1)
	obj.pushEffect(eff, true)
}

func (obj *Object) unapplyEffect(eff *effect.Effect) {
	obj.applyEffectValue(eff, -1)
	obj.pushEffect(eff, false)
	if eff.MaxHP != 0 {
		obj.PushMaxHPSD(obj.MaxHP, obj.MaxSD)
		obj.PushHPSD(obj.HP, obj.SD)
	}
	if eff.MaxMP != 0 {
		obj.PushMaxMPAG(obj.MaxMP, obj.MaxAG)
		obj.PushMPAG(obj.MP, obj.AG)
	}
}

func (obj *Object) applyEffectValue(eff *effect.Effect, sign int) {
	obj.AttackMin += sign * eff.Attack
	obj.AttackMax += sign * eff.Attack
	obj.AttackMin += sign * eff.AttackMin
	obj.AttackMax += sign * eff.AttackMax
	obj.Defense += sign * eff.Defense
	obj.DefenseRate += sign * eff.DefenseRate
	obj.MaxHP += sign * eff.MaxHP
	obj.MaxMP += sign * eff.MaxMP
	if eff.Slow {
		obj.delayLevel += sign
		if obj.delayLevel < 0 {
			obj.delayLevel = 0
		}
	}
	if obj.MaxHP < 1 {
		obj.MaxHP = 1
	}
	if obj.MaxMP < 0 {
		obj.MaxMP = 0
	}
	if sign < 0 {
		if obj.HP > obj.MaxHP {
			obj.HP = obj.MaxHP
		}
		if obj.MP > obj.MaxMP {
			obj.MP = obj.MaxMP
		}
	}
}

func (obj *Object) removeEffect(index int) {
	eff := obj.effects.Remove(index)
	if eff == nil {
		return
	}
	obj.unapplyEffect(eff)
}

func (obj *Object) pushEffect(eff *effect.Effect, add bool) {
	state := 0
	option := 1
	left := 0
	if add {
		state = 1
		option = 0
		left = int(time.Until(eff.Expire).Seconds())
		if left < 0 {
			left = 0
		}
	}
	obj.PushViewport(&model.MsgEffectStateReply{
		State:     state,
		Index:     obj.Index,
		BuffIndex: eff.BuffIndex,
	})
	obj.Push(&model.MsgEffectReply{
		OptionType:  eff.Category,
		EffectType:  eff.EffectType,
		Option:      option,
		LeftTime:    left,
		BuffIndex:   eff.BuffIndex,
		EffectValue: effectValue(eff),
	})
}

func effectValue(eff *effect.Effect) int {
	switch {
	case eff.Attack != 0:
		return eff.Attack
	case eff.Defense != 0:
		return eff.Defense
	case eff.MaxHP != 0:
		return eff.MaxHP
	case eff.MaxMP != 0:
		return eff.MaxMP
	case eff.CriticalDamage != 0:
		return eff.CriticalDamage
	case eff.Reflect != 0:
		return eff.Reflect
	case eff.DamageReduction != 0:
		return eff.DamageReduction
	case eff.Dot != 0:
		return eff.Dot
	}
	return 0
}

func (obj *Object) processEffect() {
	now := time.Now()
	for index, eff := range obj.effects {
		if eff.Dot > 0 && !now.Before(eff.NextTick) {
			eff.NextTick = now.Add(time.Second)
			source := ObjectManager.GetObject(eff.Source)
			if source != nil {
				source.attack(obj, attackRequest{
					mode:   attackModeDOT,
					damage: eff.Dot,
				})
			}
		}
		if !now.Before(eff.Expire) {
			obj.removeEffect(index)
		}
	}
}

func (obj *Object) activeBuffIndexes() []int {
	return obj.effects.ActiveBuffIndexes()
}

func (obj *Object) effect(index int) *effect.Effect {
	return obj.effects.Get(index)
}

func (obj *Object) cannotAct() bool {
	for _, eff := range obj.effects {
		if eff.Sleep {
			return true
		}
	}
	return false
}

func (obj *Object) removeSleepEffect() {
	for index, eff := range obj.effects {
		if eff.Sleep {
			obj.removeEffect(index)
		}
	}
}
