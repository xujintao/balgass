package object

import (
	"log"

	"github.com/xujintao/balgass/src/server_game/game/model"
	"github.com/xujintao/balgass/src/server_game/game/skill"
)

func (obj *object) initSkill() {
	obj.skills = make(skill.Skills)
}

func (obj *object) clearSkill() {
	obj.skills = nil
}

func (obj *object) learnSkill(index skill.SkillIndex) (*skill.Skill, bool) {
	if _, ok := obj.skills[index]; ok {
		log.Printf("[object]%s [skill]%d already exists", obj.Name, index)
		return nil, false
	}
	// obj.skills[index] = skill.SkillManager.Get(index, level, obj.skills)
	return obj.skills.Get(index)
}

func (obj *object) forgetSkill(index skill.SkillIndex) (*skill.Skill, bool) {
	if _, ok := obj.skills[index]; !ok {
		log.Printf("[object]%s [skill]%d doesn't exist", obj.Name, index)
		return nil, false
	}
	return obj.skills.Put(index)
}

func (obj *object) UseSkill(msg *model.MsgUseSkill) {
	tobj := ObjectManager.objects[msg.Target]
	if tobj == nil {
		log.Printf("UseSkill target is invalid [index]%d->[index]%d\n",
			obj.index, msg.Target)
		return
	}
	s, ok := obj.skills[skill.SkillIndex(msg.Skill)]
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

func (obj *object) canUseSkill(tobj *object, s *skill.Skill) {
	switch s.Index {
	case skill.SkillIndexDefense:
		obj.pushViewport(&model.MsgActionReply{
			Index:  obj.index,
			Action: int(skill.SkillIndexDefense),
			Dir:    obj.Dir,
			Target: tobj.index,
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
		obj.attack(tobj)
	}
}
