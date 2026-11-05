package object

import (
	"time"

	"github.com/xujintao/balgass/src/server-game/game/buff"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

func (obj *Object) initBuff() {
	obj.buffs = make(buff.Buffs)
}

func (obj *Object) clearBuff() {
	obj.removeAllBuffs()
	obj.buffs = nil
}

func (obj *Object) removeAllBuffs() {
	removed := obj.buffs.RemoveAll()
	for _, eff := range removed {
		obj.pushBuff(eff, false)
	}
	if len(removed) > 0 {
		obj.Objecter.Calc()
	}
}

func (obj *Object) AddBuff(index buff.BuffType, duration int, effect buff.BuffEffect) bool {
	eff := buff.NewBuff(index)
	eff.Expire = time.Now().Add(time.Duration(duration) * time.Second)
	eff.Effect = effect
	removed, ok := obj.buffs.Add(eff)
	for _, old := range removed {
		obj.pushBuff(old, false)
	}
	if !ok {
		if len(removed) > 0 {
			obj.Objecter.Calc()
		}
		return false
	}
	obj.pushBuff(eff, true)
	obj.Objecter.Calc()
	return true
}

func (obj *Object) RemoveBuff(index buff.BuffType) {
	eff := obj.buffs.Remove(index)
	if eff == nil {
		return
	}
	obj.pushBuff(eff, false)
	obj.Objecter.Calc()
}

func (obj *Object) pushBuff(eff *buff.Buff, add bool) {
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
		BuffIndex: int(eff.Index),
	})
	obj.Push(&model.MsgEffectReply{
		OptionType:  int(eff.EffectType),
		EffectType:  0,
		Option:      option,
		LeftTime:    left,
		BuffIndex:   int(eff.Index),
		EffectValue: buffValue(eff),
	})
}

func buffValue(eff *buff.Buff) int {
	set := buff.Buffs{eff.Index: eff}
	var option buff.BuffEffect
	set.ApplyBuffEffect(func(value buff.BuffEffect) { option = value })
	switch {
	case option.Attack != 0:
		return option.Attack
	case option.Defense != 0:
		return option.Defense
	case option.IncMaxHP != 0:
		return int(option.IncMaxHP)
	case option.IncMaxMP != 0:
		return int(option.IncMaxMP)
	case option.CriticalDamage != 0:
		return option.CriticalDamage
	case option.Reflect != 0:
		return option.Reflect
	case option.DamageReduction != 0:
		return option.DamageReduction
	case option.AttackReduction != 0:
		return option.AttackReduction
	case option.DefenseReduction != 0:
		return option.DefenseReduction
	case eff.Effect.Dot.Damage != 0:
		return eff.Effect.Dot.Damage
	}
	return 0
}

func (obj *Object) processBuffs() {
	now := time.Now()
	for index, eff := range obj.buffs {
		if eff.Effect.Dot.Damage > 0 && !now.Before(eff.Effect.Dot.NextTick) {
			eff.Effect.Dot.NextTick = now.Add(time.Second)
			source := ObjectManager.GetObject(eff.Effect.Dot.Source)
			if source != nil {
				source.attack(obj, attackRequest{
					mode:   attackModeDOT,
					damage: eff.Effect.Dot.Damage,
				})
			}
		}
		if !now.Before(eff.Expire) {
			obj.RemoveBuff(index)
		}
	}
}

func (obj *Object) activeBuffIndexes() []int {
	return obj.buffs.ActiveBuffIndexes()
}

func (obj *Object) Buff(index buff.BuffType) *buff.Buff {
	return obj.buffs.Get(index)
}

func (obj *Object) Buffs() buff.Buffs {
	return obj.buffs
}

func (obj *Object) SetEffectDelayLevel(level int) {
	obj.delayLevel = level
}

func (obj *Object) cannotAct() bool {
	for _, eff := range obj.buffs {
		if eff.Effect.Sleep || eff.Effect.Stun {
			return true
		}
	}
	return false
}

func (obj *Object) removeSleepBuff() {
	for index, eff := range obj.buffs {
		if eff.Effect.Sleep {
			obj.RemoveBuff(index)
		}
	}
}
