package object

import (
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/formula"
	"github.com/xujintao/balgass/src/server-game/game/math2"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

func (obj *Object) initSkill() {
	obj.Skills = make(skill.Skills)
	obj.skillUseTimes = make(map[int]time.Time)
}

func (obj *Object) clearSkill() {
	obj.Skills = nil
	obj.skillUseTimes = nil
}

func (obj *Object) LearnSkill(index int) (*skill.Skill, bool) {
	if _, ok := obj.Skills[index]; ok {
		slog.Error("LearnSkill obj.Skills[index] ok",
			"object", obj.Name, "skill", index)
		return nil, false
	}
	// obj.Skills[index] = skill.SkillManager.Get(index, level, obj.Skills)
	return obj.Skills.Get(index)
}

func (obj *Object) ForgetSkill(index int) (*skill.Skill, bool) {
	if _, ok := obj.Skills[index]; !ok {
		slog.Error("ForgetSkill obj.Skills[index] not ok",
			"object", obj.Name, "skill", index)
		return nil, false
	}
	return obj.Skills.Put(index)
}

func (obj *Object) UseSkill(msg *model.MsgUseSkill) {
	tobj, s, mp, ag, ok := obj.validateSkillUse(msg)
	if !ok {
		return
	}
	if !obj.useSkill(tobj, s) {
		return
	}
	obj.MP -= mp
	obj.AG -= ag
	obj.skillUseTimes[s.Index] = time.Now()
	obj.PushMPAG(obj.MP, obj.AG)
}

func (obj *Object) validateSkillUse(msg *model.MsgUseSkill) (*Object, *skill.Skill, int, int, bool) {
	tobj := ObjectManager.GetObject(msg.Target)
	if tobj == nil {
		slog.Error("UseSkill target is nil",
			"object", obj.Name, "target", msg.Target)
		return nil, nil, 0, 0, false
	}
	if !obj.Live || !tobj.Live {
		return nil, nil, 0, 0, false
	}
	if obj.MapNumber != tobj.MapNumber {
		return nil, nil, 0, 0, false
	}
	s, ok := obj.Skills[msg.Skill]
	if !ok {
		return nil, nil, 0, 0, false
	}
	mp, ag := obj.GetSkillMPAG(s)
	if obj.MP < mp || obj.AG < ag {
		return nil, nil, 0, 0, false
	}
	if !obj.checkSkillDistance(tobj, s) {
		return nil, nil, 0, 0, false
	}
	if !obj.checkSkillDelay(s) {
		return nil, nil, 0, 0, false
	}
	return tobj, s, mp, ag, true
}

func (obj *Object) checkSkillDistance(tobj *Object, s *skill.Skill) bool {
	distance := s.Distance
	if distance <= 0 {
		return obj.Index == tobj.Index || obj.CalcDistance(tobj) <= 1
	}
	return obj.CalcDistance(tobj) <= distance
}

func (obj *Object) checkSkillDelay(s *skill.Skill) bool {
	if s.Delay <= 0 {
		return true
	}
	last, ok := obj.skillUseTimes[s.Index]
	if !ok {
		return true
	}
	return time.Since(last) >= time.Duration(s.Delay)*time.Millisecond
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

func (obj *Object) useSkill(tobj *Object, s *skill.Skill) bool {
	switch s.Index {
	case skill.SkillIndexPoison, // 1毒咒
		skill.SkillIndexMeteorite,         // 2陨石
		skill.SkillIndexLightning,         // 3掌心雷
		skill.SkillIndexFireBall,          // 4火球
		skill.SkillIndexFlame,             // 5火龙
		skill.SkillIndexIce,               // 7冰封
		skill.SkillIndexTwister,           // 8龙卷风
		skill.SkillIndexEvilSpirit,        // 9黑龙波
		skill.SkillIndexPowerWave,         // 11真空波
		skill.SkillIndexAquaBeam,          // 12极光
		skill.SkillIndexCometFall,         // 13爆炎
		skill.SkillIndexEnergyBall,        // 17能量球(初始)
		skill.SkillIndexFallingSlash,      // 19地裂斩(武器)
		skill.SkillIndexLunge,             // 20牙突刺(武器)
		skill.SkillIndexUppercut,          // 21升龙击(武器)
		skill.SkillIndexCyclone,           // 22旋风斩(武器)
		skill.SkillIndexSlash,             // 23天地十字剑(武器)
		skill.SkillIndexTripleShot,        // 24多重箭(武器)
		skill.SkillIndexDecay,             // 38单毒炎
		skill.SkillIndexIceStorm,          // 39暴风雪
		skill.SkillIndexTwistingSlash,     // 41霹雳回旋斩
		skill.SkillIndexRagefulBlow,       // 42雷霆裂闪
		skill.SkillIndexCrescentMoonSlash, // 44半月斩(攻城)
		skill.SkillIndexImpale,            // 47钻云枪
		skill.SkillIndexFireBreath,        // 49流星焰(彩云兽)
		skill.SkillIndexIceArrow,          // 51冰封箭
		skill.SkillIndexPenetration,       // 52穿透箭
		skill.SkillIndexFireSlash,         // 55玄月斩
		skill.SkillIndexPowerSlash,        // 56天雷闪(武器)
		skill.SkillIndexSpiralSlash,       // 57风舞回旋斩(攻城)
		skill.SkillIndexForce,             // 60冲击(初始)
		skill.SkillIndexFireBurst,         // 61星云火链
		skill.SkillIndexElectricSpike,     // 65圣极光
		skill.SkillIndexForceWave:         // 66冲击波
		obj.useSkillAttack(tobj, s)
		switch s.Index {
		case skill.SkillIndexLightning, // 3掌心雷
			skill.SkillIndexFallingSlash, // 19地裂斩(武器)
			skill.SkillIndexLunge,        // 20牙突刺(武器)
			skill.SkillIndexUppercut,     // 21升龙击(武器)
			skill.SkillIndexCyclone,      // 22旋风斩(武器)
			skill.SkillIndexSlash:        // 23天地十字剑(武器)
			obj.AddDelayMsg(2, 0, 150, tobj.Index) // delay knockback target
		}
		return true
	case skill.SkillIndexDeathStab: // 43袭风刺
		obj.useSkillDeathStab(s, tobj)
		return true
	case skill.SkillIndexDefense: // 18圣盾防御
		obj.PushViewport(&model.MsgActionReply{
			Index:  obj.Index,
			Action: int(skill.SkillIndexDefense),
			Dir:    obj.Dir,
			Target: tobj.Index,
		})
		return true
	case skill.SkillIndexHeal: // 26治疗
		return obj.useSkillHeal(tobj, s)
	case skill.SkillIndexGreaterDefense: // 27防御
		return obj.useSkillGreaterDefense(tobj, s)
	case skill.SkillIndexGreaterAttack: // 28攻击
		return obj.useSkillGreaterAttack(tobj, s)
	case skill.SkillIndexSwellHP: // 48生命之光
		return obj.useSkillSwellHP(tobj, s)
	}
	return false
}

func (obj *Object) useSkillAttack(tobj *Object, s *skill.Skill) {
	obj.UseSkillReply(tobj, s, true)
	obj.attack(tobj, s, 0)
}

func (obj *Object) getAngle(tobj *Object) float32 {
	x := tobj.X - obj.X
	y := tobj.Y - obj.Y
	rad := float32(math.Atan2(float64(y), float64(x)))
	return rad*180/math.Pi + 90
}

func (obj *Object) CreateSkillFrustum(a, x, y float32) {
	p := [4][3]float32{
		{-x, y, 0.0},
		{x, y, 0.0},
		{1.0, 0.0, 0.0},
		{-1.0, 0.0, 0.0},
	}
	angle := [3]float32{0.0, 0.0, a}
	matrix := math2.Angle2Matrix(angle)
	var frustum [4][3]float32
	for i := 0; i < 4; i++ {
		frustum[i] = math2.VectorRotate(p[i], matrix)
		obj.SkillFrustumX[i] = int(frustum[i][0]) + obj.X
		obj.SkillFrustumY[i] = int(frustum[i][1]) + obj.Y
	}
}

func (obj *Object) CheckSkillFrustum(tobj *Object) bool {
	x := tobj.X
	y := tobj.Y
	for i, j := 0, 3; i < MaxArrayFrustum; j, i = i, i+1 {
		frustum := (obj.SkillFrustumX[i]-x)*(obj.SkillFrustumY[j]-y) -
			(obj.SkillFrustumX[j]-x)*(obj.SkillFrustumY[i]-y)
		if frustum < 0 {
			return false
		}
	}
	return true
}

func (obj *Object) useSkillDeathStab(s *skill.Skill, tobj *Object) {
	obj.UseSkillReply(tobj, s, true)
	obj.attack(tobj, s, 0)
	if rand.Intn(100)%3 == 0 {
		obj.attack(tobj, s, 0)
	}
	angle := obj.getAngle(tobj)
	obj.CreateSkillFrustum(angle, 1.5, 3.0)
	obj.ForEachViewportObject(func(vpobj *Object) {
		if vpobj != tobj &&
			vpobj.Live &&
			vpobj.Type != ObjectTypePlayer &&
			vpobj.Type != ObjectTypeNPC &&
			obj.CheckSkillFrustum(vpobj) {
			obj.attack(vpobj, s, 0)
		}
	})
}

func (obj *Object) useSkillHeal(tobj *Object, s *skill.Skill) bool {
	if obj.Class != int(class.Elf) || !tobj.canReceiveSupportSkill() {
		return false
	}
	addLife := 0
	formula.ElfHeal(tobj.Class, obj.Index, tobj.Index, obj.GetEnergy(), &addLife)
	if addLife <= 0 {
		return false
	}
	tobj.HP += addLife
	if tobj.HP > tobj.MaxHP {
		tobj.HP = tobj.MaxHP
	}
	obj.UseSkillReply(tobj, s, true)
	tobj.PushHPSD(tobj.HP, tobj.SD)
	return true
}

func (obj *Object) useSkillGreaterDefense(tobj *Object, s *skill.Skill) bool {
	if obj.Class != int(class.Elf) || !tobj.canReceiveSupportSkill() {
		return false
	}
	defense, duration := 0.0, 0.0
	formula.ElfDefense(tobj.Class, obj.Index, tobj.Index, obj.GetEnergy(), &defense, &duration)
	if defense <= 0 || duration <= 0 {
		return false
	}
	tobj.addSkillEffect(s.Index, &skillEffect{
		defense: int(defense),
		expire:  time.Now().Add(time.Duration(duration) * time.Second),
	})
	obj.UseSkillReply(tobj, s, true)
	return true
}

func (obj *Object) useSkillGreaterAttack(tobj *Object, s *skill.Skill) bool {
	if obj.Class != int(class.Elf) || !tobj.canReceiveSupportSkill() {
		return false
	}
	attack, duration := 0.0, 0.0
	formula.ElfAttack(tobj.Class, obj.Index, tobj.Index, obj.GetEnergy(), &attack, &duration)
	if attack <= 0 || duration <= 0 {
		return false
	}
	tobj.addSkillEffect(s.Index, &skillEffect{
		attack: int(attack),
		expire: time.Now().Add(time.Duration(duration) * time.Second),
	})
	obj.UseSkillReply(tobj, s, true)
	return true
}

func (obj *Object) useSkillSwellHP(tobj *Object, s *skill.Skill) bool {
	if (obj.Class != int(class.Knight) && obj.Class != int(class.Magumsa)) || !tobj.canReceiveSupportSkill() {
		return false
	}
	addLifeRate := 0.0
	duration := 0
	formula.KnightSkillAddLife(obj.GetVitality(), obj.GetEnergy(), 0, &addLifeRate, &duration)
	addHP := int(float64(tobj.MaxHP) * addLifeRate / 100)
	if addHP <= 0 || duration <= 0 {
		return false
	}
	tobj.addSkillEffect(s.Index, &skillEffect{
		maxHP:  addHP,
		expire: time.Now().Add(time.Duration(duration) * time.Second),
	})
	obj.UseSkillReply(tobj, s, true)
	tobj.PushMaxHPSD(tobj.MaxHP, tobj.MaxSD)
	tobj.PushHPSD(tobj.HP, tobj.SD)
	return true
}

func (obj *Object) canReceiveSupportSkill() bool {
	if obj.Type == ObjectTypePlayer {
		return true
	}
	return obj.Index >= ObjectManager.maxMonsterCount && obj.Index < ObjectManager.playerStartIndex
}
