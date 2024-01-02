package object

import (
	"log"
	"math"
	"math/rand"

	"github.com/xujintao/balgass/src/server_game/game/math2"
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
	case skill.SkillIndexDefense: // 18圣盾防御
		obj.PushViewport(&model.MsgActionReply{
			Index:  obj.Index,
			Action: int(skill.SkillIndexDefense),
			Dir:    obj.Dir,
			Target: tobj.Index,
		})
	case skill.SkillIndexPoison, // 1毒咒
		skill.SkillIndexMeteorite,     // 2陨石
		skill.SkillIndexLightning,     // 3掌心雷
		skill.SkillIndexFireBall,      // 4火球
		skill.SkillIndexFlame,         // 5火龙
		skill.SkillIndexIce,           // 7冰封
		skill.SkillIndexEnergyBall,    // 17能量球(初始)
		skill.SkillIndexFallingSlash,  // 19地裂斩(武器)
		skill.SkillIndexLunge,         // 20牙突刺(武器)
		skill.SkillIndexUppercut,      // 21升龙击(武器)
		skill.SkillIndexCyclone,       // 22旋风斩(武器)
		skill.SkillIndexSlash,         // 23天地十字剑(武器)
		skill.SkillIndexTripleShot,    // 24多重箭(武器)
		skill.SkillIndexTwistingSlash, // 41霹雳回旋斩
		skill.SkillIndexRagefulBlow,   // 42雷霆裂闪
		skill.SkillIndexImpale,        // 47钻云枪
		skill.SkillIndexPenetration,   // 52穿透箭
		skill.SkillIndexPowerSlash:    // 56天雷闪(武器)
		obj.UseSkillReply(tobj, s, true)
		obj.attack(tobj, 0)
	case skill.SkillIndexDeathStab: // 43袭风刺
		obj.UseSkillDeathStab(s, tobj)
	}
}

func (obj *Object) getAngle(tobj *Object) float32 {
	x := tobj.X - obj.X
	y := tobj.Y - obj.Y
	rad := float32(math.Atan2(float64(y), float64(x)))
	return rad*180/math.Pi + 90
}

func (obj *Object) CreateSkillFrustrum(a, x, y float32) {
	p := [4][3]float32{
		{-x, y, 0.0},
		{x, y, 0.0},
		{1.0, 0.0, 0.0},
		{-1.0, 0.0, 0.0},
	}
	angle := [3]float32{0.0, 0.0, a}
	matrix := math2.Angle2Matrix(angle)
	var frustrum [4][3]float32
	for i := 0; i < 4; i++ {
		frustrum[i] = math2.VectorRotate(p[i], matrix)
		obj.SkillFrustrumX[i] = int(frustrum[i][0]) + obj.X
		obj.SkillFrustrumY[i] = int(frustrum[i][1]) + obj.Y
	}
}

func (obj *Object) CheckSkillFrustrum(tobj *Object) bool {
	x := tobj.X
	y := tobj.Y
	for i, j := 0, 3; i < MaxArrayFrustrum; j, i = i, i+1 {
		frustrum := (obj.SkillFrustrumX[i]-x)*(obj.SkillFrustrumY[j]-y) -
			(obj.SkillFrustrumX[j]-x)*(obj.SkillFrustrumY[i]-y)
		if frustrum < 0 {
			return false
		}
	}
	return true
}

func (obj *Object) UseSkillDeathStab(s *skill.Skill, tobj *Object) {
	obj.UseSkillReply(tobj, s, true)
	obj.attack(tobj, 0)
	if rand.Intn(100)%3 == 0 {
		obj.attack(tobj, 0)
	}
	angle := obj.getAngle(tobj)
	obj.CreateSkillFrustrum(angle, 1.5, 3.0)
	obj.ForEachViewportObject(func(vpobj *Object) {
		if vpobj != tobj && obj.CheckSkillFrustrum(vpobj) {
			obj.attack(vpobj, 0)
		}
	})
}
