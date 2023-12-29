package object

import (
	"log"

	"github.com/xujintao/balgass/src/server_game/game/model"
	"github.com/xujintao/balgass/src/server_game/game/skill"
)

func (obj *Object) initSkill() {
	obj.Skills = make(skill.Skills)
}

func (obj *Object) clearSkill() {
	obj.Skills = nil
}

func (obj *Object) LearnSkill(index int) (*skill.Skill, bool) {
	if _, ok := obj.Skills[index]; ok {
		log.Printf("[object]%s [skill]%d already exists", obj.Name, index)
		return nil, false
	}
	// obj.Skills[index] = skill.SkillManager.Get(index, level, obj.Skills)
	return obj.Skills.Get(index)
}

func (obj *Object) ForgetSkill(index int) (*skill.Skill, bool) {
	if _, ok := obj.Skills[index]; !ok {
		log.Printf("[object]%s [skill]%d doesn't exist", obj.Name, index)
		return nil, false
	}
	return obj.Skills.Put(index)
}

func (obj *Object) UseSkill(msg *model.MsgUseSkill) {
	tobj := ObjectManager.objects[msg.Target]
	if tobj == nil {
		log.Printf("UseSkill target is invalid [index]%d->[index]%d\n",
			obj.Index, msg.Target)
		return
	}
	s, ok := obj.Skills[msg.Skill]
	if !ok {
		return
	}
	mp, ag := obj.GetSkillMPAG(s)
	if obj.MP < mp || obj.AG < ag {
		return
	}
	obj.canUseSkill(tobj, s)
	obj.MP -= mp
	obj.AG -= ag
	obj.PushMPAG(obj.MP, obj.AG)
}

func (obj *Object) UseSkillReply(tobj *Object, s *skill.Skill, success bool) {
	target := tobj.Index
	if success {
		target |= 0x8000
	}
	reply := model.MsgUseSkillReply{
		Index:  obj.Index,
		Skill:  s.Index,
		Target: target,
	}
	obj.PushViewport(&reply)
}

func (obj *Object) canUseSkill(tobj *Object, s *skill.Skill) {
	switch s.Index {
	case skill.SkillIndexDefense: // 圣盾防御
		obj.PushViewport(&model.MsgActionReply{
			Index:  obj.Index,
			Action: int(skill.SkillIndexDefense),
			Dir:    obj.Dir,
			Target: tobj.Index,
		})
	case skill.SkillIndexPoison, // 毒咒
		skill.SkillIndexMeteorite,     // 陨石
		skill.SkillIndexLightning,     // 掌心雷
		skill.SkillIndexFireBall,      // 火球
		skill.SkillIndexFlame,         // 火龙
		skill.SkillIndexIce,           // 冰封
		skill.SkillIndexEnergyBall,    // 能量球(初始)
		skill.SkillIndexFallingSlash,  // 地裂斩(武器)
		skill.SkillIndexLunge,         // 牙突刺(武器)
		skill.SkillIndexUppercut,      // 升龙击(武器)
		skill.SkillIndexCyclone,       // 旋风斩(武器)
		skill.SkillIndexSlash,         // 天地十字剑(武器)
		skill.SkillIndexTripleShot,    // 多重箭(武器)
		skill.SkillIndexTwistingSlash, // 霹雳回旋斩
		skill.SkillIndexRagefulBlow,   // 雷霆裂闪
		skill.SkillIndexImpale,        // 钻云枪
		skill.SkillIndexPenetration,   // 穿透箭
		skill.SkillIndexPowerSlash:    // 天雷闪(武器)
		obj.attack(tobj, 0)
	case skill.SkillIndexDeathStab: // 袭风刺
		obj.UseSkillDeathStab(s, tobj)
	}
}

func (obj *Object) UseSkillDeathStab(s *skill.Skill, tobj *Object) {
	obj.UseSkillReply(tobj, s, true)
	obj.attack(tobj, 0)
}
