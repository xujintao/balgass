package object

import (
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/effect"
	"github.com/xujintao/balgass/src/server-game/game/formula"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/math2"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

const (
	durationSkillMaxPackets = 5
	durationSkillMaxTargets = 5
	durationSkillMaxKey     = 60
)

func (obj *Object) initSkill() {
	obj.Skills = make(skill.Skills)
	obj.skillUseTimes = make(map[int]time.Time)
}

func (obj *Object) clearSkill() {
	obj.Skills = nil
	obj.skillUseTimes = nil
	obj.durationSkill = durationSkillState{}
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
	if obj.cannotAct() {
		return
	}
	if !obj.novaSkill.startedAt.IsZero() && msg.Skill != skill.SkillIndexNova {
		obj.releaseNovaSkill()
		return
	}
	if msg.Skill == skill.SkillIndexNova && obj.Type == ObjectTypePlayer {
		obj.useSkillNova(msg)
		return
	}
	if obj.Type == ObjectTypePlayer {
		switch msg.Skill {
		case skill.SkillIndexFlame, // 5火龙
			skill.SkillIndexTwister,           // 8龙卷风
			skill.SkillIndexEvilSpirit,        // 9黑龙波
			skill.SkillIndexHellFire,          // 10地狱火
			skill.SkillIndexAquaBeam,          // 12极光
			skill.SkillIndexCometFall,         // 13爆炎
			skill.SkillIndexInferno,           // 14毁灭烈焰
			skill.SkillIndexTripleShot,        // 24多重箭(武器)
			skill.SkillIndexLance,             // 45回旋刃(攻城)
			skill.SkillIndexImpale,            // 47钻云枪
			skill.SkillIndexPenetration,       // 52穿透箭
			skill.SkillIndexFireSlash,         // 55玄月斩
			skill.SkillIndexEarthshake,        // 62地裂(黑王马)
			skill.SkillIndexPlasmaStorm,       // 76闪电链(炎狼兽)
			skill.SkillIndexFireScream,        // 78火舞旋风
			skill.SkillIndexDrainLife,         // 214摄魂咒
			skill.SkillIndexWeakness,          // 221虚弱阵
			skill.SkillIndexInnovation,        // 222破御阵
			skill.SkillIndexSummonerExplosion, // 223爆裂
			skill.SkillIndexRequiem,           // 224刺袭
			skill.SkillIndexPollution:         // 225污染
			slog.Warn("UseSkill duration skill requires UseSkillDuration",
				"object", obj.Name, "index", obj.Index, "skill", msg.Skill)
			return
		}
	}
	tobj, s, mp, ag, ok := obj.validateSkillUse(msg)
	if !ok {
		return
	}
	var success bool
	switch s.Index {
	case skill.SkillIndexPoison, // 1毒咒
		skill.SkillIndexMeteorite,         // 2陨石
		skill.SkillIndexLightning,         // 3掌心雷
		skill.SkillIndexFireBall,          // 4火球
		skill.SkillIndexIce,               // 7冰封
		skill.SkillIndexEvilSpirit,        // 9黑龙波
		skill.SkillIndexPowerWave,         // 11真空波
		skill.SkillIndexEnergyBall,        // 17能量球(初始)
		skill.SkillIndexFallingSlash,      // 19地裂斩(武器)
		skill.SkillIndexLunge,             // 20牙突刺(武器)
		skill.SkillIndexUppercut,          // 21升龙击(武器)
		skill.SkillIndexCyclone,           // 22旋风斩(武器)
		skill.SkillIndexSlash,             // 23天地十字剑(武器)
		skill.SkillIndexCrescentMoonSlash, // 44半月斩(攻城)
		skill.SkillIndexStarfall,          // 46天堂之箭(攻城)
		skill.SkillIndexFireBreath,        // 49流星焰(彩云兽)
		skill.SkillIndexIceArrow,          // 51冰封箭
		skill.SkillIndexFireSlash,         // 55玄月斩
		skill.SkillIndexSpiralSlash,       // 57风舞回旋斩(攻城)
		skill.SkillIndexForceWave,         // 66冲击波
		skill.SkillIndexKillingBlow,       // 260幽冥青狼拳
		skill.SkillIndexBeastUppercut,     // 261斗气爆裂拳
		skill.SkillIndexChainDrive,        // 262回旋踢
		skill.SkillIndexDragonSlasher,     // 265噬血之龙
		skill.SkillIndexCharge:            // 269冲锋(攻城)
		switch s.Index {
		case skill.SkillIndexStarfall, // 46天堂之箭(攻城)
			skill.SkillIndexIceArrow: // 51冰封箭
			if !obj.UseSkillAmmo(s, false) {
				break
			}
			success = obj.useSkillAttack(tobj, s)
			if success {
				obj.UseSkillAmmo(s, true)
			}
		default:
			success = obj.useSkillAttack(tobj, s)
		}
	case skill.SkillIndexSoulBarrier: // 16守护之魂
		success = obj.useSkillSoulBarrier(tobj, s)
	case skill.SkillIndexDefense: // 18圣盾防御
		success = obj.useSkillDefense(tobj)
	case skill.SkillIndexHeal: // 26治疗
		success = obj.useSkillHeal(tobj, s)
	case skill.SkillIndexGreaterDefense: // 27防御
		success = obj.useSkillGreaterDefense(tobj, s)
	case skill.SkillIndexGreaterAttack: // 28攻击
		success = obj.useSkillGreaterAttack(tobj, s)
	case skill.SkillIndexSummonGoblin, // 30召唤哥布林
		skill.SkillIndexSummonStoneGolem, // 31召唤石巨人
		skill.SkillIndexSummonAssassin,   // 32召唤暗杀者
		skill.SkillIndexSummonEliteYeti,  // 33召唤雪人王
		skill.SkillIndexSummonDarkKnight, // 34召唤暗黑骑士
		skill.SkillIndexSummonBali,       // 35召唤巴里
		skill.SkillIndexSummonSoldier:    // 36召唤黄金斗士
		success = obj.useSkillSummon(s)
	case skill.SkillIndexDecay: // 38单毒炎
		success = obj.useSkillAttackArea(tobj, s, 3)
	case skill.SkillIndexIceStorm: // 39暴风雪
		success = obj.useSkillAttackArea(tobj, s, s.Distance)
	case skill.SkillIndexTwistingSlash: // 41霹雳回旋斩
		success = obj.useSkillAttackAreaSelf(tobj, s, s.Distance)
	case skill.SkillIndexRagefulBlow: // 42雷霆裂闪
		success = obj.useSkillAttackArea(tobj, s, 3)
	case skill.SkillIndexDeathStab: // 43袭风刺
		success = obj.useSkillDeathStab(s, tobj)
	case skill.SkillIndexSwellHP: // 48生命之光
		success = obj.useSkillSwellHP(tobj, s)
	case skill.SkillIndexPowerSlash: // 56天雷闪(武器)
		success = obj.useSkillAttackFrustum(tobj, s, 6.0, 6.0, 1.0, 0.0)
	case skill.SkillIndexForce: // 60冲击(初始)
		success = obj.useSkillAttackHitBox(tobj, s, skill.SkillManager.CheckSpearHitBox)
	case skill.SkillIndexFireBurst: // 61星云火链
		success = obj.useSkillAttackArea(tobj, s, 2)
	case skill.SkillIndexAddCriticalDamage: // 64致命圣印
		success = obj.useSkillAddCriticalDamage(s)
	case skill.SkillIndexElectricSpike: // 65圣极光
		success = obj.useSkillAttackHitBox(tobj, s, skill.SkillManager.CheckElectricHitBox)
	case skill.SkillIndexInfinityArrow: // 77无影箭
		success = obj.useSkillInfinityArrow(s)
	case skill.SkillIndexChainLightning: // 215链雷咒
		success = obj.useSkillChainLightning(tobj, s)
	case skill.SkillIndexDamageReflection: // 217伤害反射
		success = obj.useSkillDamageReflection(tobj, s)
	case skill.SkillIndexBerserker: // 218狂暴术
		success = obj.useSkillBerserker(s)
	case skill.SkillIndexSleep: // 219昏睡
		success = obj.useSkillSleep(tobj, s)
	case skill.SkillIndexLightningShock: // 230烈光闪
		success = obj.useSkillAttackAreaSelf(tobj, s, 6)
	case skill.SkillIndexStrikeDestruction: // 232破坏一击
		success = obj.useSkillAttackArea(tobj, s, 3)
	case skill.SkillIndexExpansionWizardry: // 233法神附体
		success = obj.useSkillExpansionWizardry(s)
	case skill.SkillIndexRecovery: // 234防护值恢复
		success = obj.useSkillRecovery(tobj, s)
	case skill.SkillIndexMultiShot: // 235五重箭
		if !obj.UseSkillAmmo(s, false) {
			break
		}
		success = obj.useSkillAttackFrustum(tobj, s, 6.0, 7.0, 1.0, 0.0)
		if success {
			obj.UseSkillAmmo(s, true)
		}
	case skill.SkillIndexFlameStrike: // 236火剑袭
		success = obj.useSkillAttackFrustum(tobj, s, 2.0, 4.0, 5.0, 0.0)
	case skill.SkillIndexGiganticStorm: // 237闪电轰顶
		success = obj.useSkillAttackArea(tobj, s, 6)
	case skill.SkillIndexChaoticDiseier: // 238黑暗之力
		success = obj.useSkillAttackFrustum(tobj, s, 1.5, 6.0, 1.5, 0.0)
	case skill.SkillIndexDarkSide: // 263幽冥光速拳
		success = obj.useSkillDarkSide(tobj, s)
	case skill.SkillIndexDragonRoar: // 264炎龙拳
		success = obj.useSkillAttackArea(tobj, s, 3)
	case skill.SkillIndexIgnoreDefense: // 266斗神-破
		success = obj.useSkillIgnoreDefense(s)
	case skill.SkillIndexIncreaseHealth: // 267斗神-命
		success = obj.useSkillIncreaseHealth(s)
	case skill.SkillIndexIncreaseBlock: // 268斗神-御
		success = obj.useSkillIncreaseBlock(s)
	case skill.SkillIndexPhoenixShot: // 270神圣气旋
		success = obj.useSkillPhoenixShot(tobj, s)
	default:
		return
	}
	if !success {
		return
	}
	obj.MP -= mp
	obj.AG -= ag
	obj.skillUseTimes[s.Index] = time.Now()
	obj.PushMPAG(obj.MP, obj.AG)
}

func (obj *Object) UseSkillDuration(msg *model.MsgUseSkillDuration) {
	if obj.Type != ObjectTypePlayer || obj.cannotAct() {
		return
	}
	switch msg.Skill {
	case skill.SkillIndexFlame, // 5火龙
		skill.SkillIndexTwister,           // 8龙卷风
		skill.SkillIndexEvilSpirit,        // 9黑龙波
		skill.SkillIndexHellFire,          // 10地狱火
		skill.SkillIndexAquaBeam,          // 12极光
		skill.SkillIndexCometFall,         // 13爆炎
		skill.SkillIndexInferno,           // 14毁灭烈焰
		skill.SkillIndexTripleShot,        // 24多重箭(武器)
		skill.SkillIndexLance,             // 45回旋刃(攻城)
		skill.SkillIndexImpale,            // 47钻云枪
		skill.SkillIndexPenetration,       // 52穿透箭
		skill.SkillIndexFireSlash,         // 55玄月斩
		skill.SkillIndexEarthshake,        // 62地裂(黑王马)
		skill.SkillIndexPlasmaStorm,       // 76闪电链(炎狼兽)
		skill.SkillIndexFireScream,        // 78火舞旋风
		skill.SkillIndexDrainLife,         // 214摄魂咒
		skill.SkillIndexWeakness,          // 221虚弱阵
		skill.SkillIndexInnovation,        // 222破御阵
		skill.SkillIndexSummonerExplosion, // 223爆裂
		skill.SkillIndexRequiem,           // 224刺袭
		skill.SkillIndexPollution:         // 225污染
	default:
		slog.Warn("UseSkillDuration skill requires UseSkill",
			"object", obj.Name, "index", obj.Index, "skill", msg.Skill)
		return
	}
	if msg.Skill == skill.SkillIndexEvilSpirit &&
		(msg.MagicKey <= 0 || msg.MagicKey > durationSkillMaxKey) {
		return
	}

	tobj, s, mp, ag, ok := obj.validateSkillUse(&model.MsgUseSkill{
		Target: msg.Target,
		Skill:  msg.Skill,
	})
	if !ok {
		return
	}
	obj.durationSkill = durationSkillState{}
	x, y := msg.X, msg.Y
	if x == 0 && y == 0 {
		x, y = tobj.X, tobj.Y
	}

	switch s.Index {
	case skill.SkillIndexLance: // 45回旋刃(攻城)
		if !obj.useSkillAttackAreaSelf(tobj, s, s.Distance) {
			return
		}
	case skill.SkillIndexEarthshake: // 62地裂(黑王马)
		if !obj.useSkillEarthshake(tobj, s) {
			return
		}
	case skill.SkillIndexPlasmaStorm: // 76闪电链(炎狼兽)
		if !obj.useSkillPlasmaStorm(tobj, s) {
			return
		}
	case skill.SkillIndexDrainLife: // 214摄魂咒
		if !obj.useSkillDrainLife(tobj, s) {
			return
		}
	case skill.SkillIndexWeakness, skill.SkillIndexInnovation: // 221虚弱阵,222破御阵
		if !obj.useSkillCurseArea(tobj, s) {
			return
		}
	case skill.SkillIndexSummonerExplosion, skill.SkillIndexRequiem: // 223爆裂,224刺袭
		if !obj.useSkillAttackAreaPoint(tobj, s, x, y, 2, 1000) {
			return
		}
	case skill.SkillIndexPollution: // 225污染
		if !obj.useSkillAttackAreaPoint(tobj, s, x, y, 3, 0) {
			return
		}
	default:
		switch s.Index {
		case skill.SkillIndexTripleShot, // 24多重箭(武器)
			skill.SkillIndexPenetration: // 52穿透箭
			if !obj.UseSkillAmmo(s, false) {
				return
			}
		}
		now := time.Now()
		state := durationSkillState{
			index:     s.Index,
			startedAt: now,
			x:         x,
			y:         y,
			dir:       msg.Dir,
		}
		if msg.MagicKey > 0 && msg.MagicKey <= durationSkillMaxKey {
			state.magicKeys[msg.MagicKey] = now
		}
		obj.durationSkill = state
		switch s.Index {
		case skill.SkillIndexTripleShot, // 24多重箭(武器)
			skill.SkillIndexPenetration: // 52穿透箭
			obj.UseSkillAmmo(s, true)
		}
	}

	now := time.Now()
	obj.MP -= mp
	obj.AG -= ag
	obj.skillUseTimes[s.Index] = now
	obj.PushViewport(&model.MsgUseSkillDurationReply{
		X:     tobj.X,
		Y:     tobj.Y,
		Dir:   msg.Dir,
		Skill: s.Index,
		Index: obj.Index,
	})
	obj.PushMPAG(obj.MP, obj.AG)
}

func (obj *Object) UseSkillAttackMultiTarget(msg *model.MsgUseSkillAttackMultiTarget) {
	if obj.Type != ObjectTypePlayer || !obj.Live {
		return
	}
	state := &obj.durationSkill
	if state.index != msg.Skill || state.startedAt.IsZero() {
		slog.Warn("UseSkillAttackMultiTarget requires matching UseSkillDuration",
			"object", obj.Name, "index", obj.Index, "skill", msg.Skill,
			"state_skill", state.index, "state_empty", state.startedAt.IsZero())
		return
	}
	now := time.Now()
	elapsed := now.Sub(state.startedAt)
	if elapsed < 0 || elapsed > 8*time.Second {
		obj.durationSkill = durationSkillState{}
		return
	}
	if state.count >= durationSkillMaxPackets || len(msg.Targets) == 0 {
		return
	}
	s, ok := obj.Skills[msg.Skill]
	if !ok {
		return
	}

	seen := make(map[int]struct{}, len(msg.Targets))
	for i, target := range msg.Targets {
		if i >= durationSkillMaxTargets {
			break
		}
		if _, ok := seen[target.Target]; ok {
			return
		}
		seen[target.Target] = struct{}{}
	}

	state.count++
	for i, target := range msg.Targets {
		if i >= durationSkillMaxTargets {
			break
		}
		if s.Index == skill.SkillIndexEvilSpirit {
			if target.MagicKey <= 0 || target.MagicKey > durationSkillMaxKey {
				continue
			}
			startedAt := state.magicKeys[target.MagicKey]
			if startedAt.IsZero() || now.Sub(startedAt) > 15*time.Second {
				continue
			}
		}

		tobj := ObjectManager.GetObject(target.Target)
		if tobj == nil ||
			!tobj.Live ||
			tobj.Type != ObjectTypeMonster ||
			tobj.MapNumber != obj.MapNumber ||
			obj.CalcDistance(tobj) > 13 {
			continue
		}
		if s.Index == skill.SkillIndexHellFire && obj.CalcDistance(tobj) > 4 {
			continue
		}
		if s.Index == skill.SkillIndexInferno && obj.CalcDistance(tobj) > 5 {
			continue
		}
		inViewport := false
		for _, vp := range obj.Viewports {
			if vp.State != 0 && vp.Type != 5 && vp.Number == tobj.Index {
				inViewport = true
				break
			}
		}
		if !inViewport {
			continue
		}
		if !obj.checkDurationTarget(tobj, s, state) {
			continue
		}
		obj.attack(tobj, s, 0, true)
		if s.Index == skill.SkillIndexFireScream &&
			rand.Intn(10000) < skill.Settings.FireScreamExplosionRate {
			obj.AddDelayMsg(8, s.Index, 300, tobj.Index)
		}
	}
}

func (obj *Object) checkDurationTarget(tobj *Object, s *skill.Skill, state *durationSkillState) bool {
	distance := s.Distance
	if s.Index == skill.SkillIndexHellFire {
		distance = 4
	} else if s.Index == skill.SkillIndexInferno {
		distance = 5
	}
	if distance > 0 && obj.CalcDistance(tobj) > distance {
		return false
	}
	if !maps.MapManager.CheckMapNoWall(obj.MapNumber, obj.X, obj.Y, tobj.X, tobj.Y) {
		return false
	}
	switch s.Index {
	case skill.SkillIndexCometFall: // 13爆炎
		return abs(tobj.X-state.x) <= 2 && abs(tobj.Y-state.y) <= 2
	case skill.SkillIndexFlame, // 5火龙
		skill.SkillIndexTwister,     // 8龙卷风
		skill.SkillIndexAquaBeam,    // 12极光
		skill.SkillIndexImpale,      // 47钻云枪
		skill.SkillIndexPenetration: // 52穿透箭
		angle := float32(math.Atan2(float64(state.y-obj.Y), float64(state.x-obj.X)))*180/math.Pi - 90
		obj.createSkillFrustum3(angle, 1.5, float32(max(distance, 1)), 1.0, 0.0)
		return obj.checkSkillFrustum(tobj)
	case skill.SkillIndexTripleShot: // 24多重箭
		angle := float32(math.Atan2(float64(state.y-obj.Y), float64(state.x-obj.X)))*180/math.Pi - 90
		obj.createSkillFrustum3(angle, 6.0, 7.0, 1.0, 0.0)
		return obj.checkSkillFrustum(tobj)
	case skill.SkillIndexFireSlash: // 55玄月斩
		angle := float32(math.Atan2(float64(state.y-obj.Y), float64(state.x-obj.X)))*180/math.Pi - 90
		obj.createSkillFrustum3(angle, 2.0, 2.0, 1.0, 0.0)
		return obj.checkSkillFrustum(tobj)
	}
	return true
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

func (obj *Object) useSkillReply(tobj *Object, s *skill.Skill, success bool) {
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

func (obj *Object) useSkillAttack(tobj *Object, s *skill.Skill) bool {
	obj.useSkillReply(tobj, s, true)
	obj.attack(tobj, s, 0, true)
	switch s.Index {
	case skill.SkillIndexLightning, // 3掌心雷
		skill.SkillIndexFallingSlash, // 19地裂斩(武器)
		skill.SkillIndexLunge,        // 20牙突刺(武器)
		skill.SkillIndexUppercut,     // 21升龙击(武器)
		skill.SkillIndexCyclone,      // 22旋风斩(武器)
		skill.SkillIndexSlash,        // 23天地十字剑(武器)
		skill.SkillIndexPhoenixShot:  // 270神圣气旋
		obj.AddDelayMsg(2, 0, 150, tobj.Index) // delay knockback target
	}
	return true
}

func (obj *Object) applySkillEffect(tobj *Object, s *skill.Skill, damage int) {
	if s == nil || !tobj.Live {
		return
	}
	now := time.Now()
	switch s.Index {
	case skill.SkillIndexPoison: // 1毒咒
		tobj.addEffect(&effect.Effect{
			BuffIndex: effect.BuffPoison, Dot: max(1, tobj.HP*3/100), Source: obj.Index,
			Expire: now.Add(20 * time.Second), NextTick: now.Add(time.Second),
		})
	case skill.SkillIndexIce: // 7冰封
		tobj.addEffect(&effect.Effect{BuffIndex: effect.BuffIce, Slow: true, Expire: now.Add(10 * time.Second)})
	case skill.SkillIndexDecay: // 38单毒炎
		tobj.addEffect(&effect.Effect{
			BuffIndex: effect.BuffPoison, Dot: max(1, tobj.HP*3/100), Source: obj.Index,
			Expire: now.Add(20 * time.Second), NextTick: now.Add(time.Second),
		})
	case skill.SkillIndexIceStorm: // 39暴风雪
		tobj.addEffect(&effect.Effect{BuffIndex: effect.BuffIce, Slow: true, Expire: now.Add(10 * time.Second)})
	case skill.SkillIndexIceArrow: // 51冰封箭
		tobj.addEffect(&effect.Effect{BuffIndex: effect.BuffIceArrow, Sleep: true, Expire: now.Add(time.Duration(skill.Settings.IceArrowTime) * time.Second)})
	case skill.SkillIndexSummonerExplosion: // 223爆裂
		dot, duration := 0, 0
		formula.ExplosionDotDamage(damage, 0, &dot, &duration)
		tobj.addEffect(&effect.Effect{
			BuffIndex: effect.BuffExplosion, Dot: dot, Source: obj.Index,
			Expire: now.Add(time.Duration(duration) * time.Second), NextTick: now.Add(time.Second),
		})
	case skill.SkillIndexRequiem: // 224刺袭
		dot, duration := 0, 0
		formula.RequiemDotDamage(damage, &dot, &duration)
		tobj.addEffect(&effect.Effect{
			BuffIndex: effect.BuffRequiem, Dot: dot, Source: obj.Index,
			Expire: now.Add(time.Duration(duration) * time.Second), NextTick: now.Add(time.Second),
		})
	case skill.SkillIndexPollution: // 225污染
		tobj.addEffect(&effect.Effect{BuffIndex: effect.BuffIce, Slow: true, Expire: now.Add(2 * time.Second)})
	case skill.SkillIndexLightningShock: // 230烈光闪
		tobj.addEffect(&effect.Effect{BuffIndex: effect.BuffLightningShock, Expire: now.Add(time.Second)})
	case skill.SkillIndexStrikeDestruction: // 232破坏一击
		tobj.addEffect(&effect.Effect{BuffIndex: effect.BuffCold, Slow: true, Expire: now.Add(10 * time.Second)})
	case skill.SkillIndexFlameStrike: // 236火剑袭
		tobj.addEffect(&effect.Effect{BuffIndex: effect.BuffFlameStrike, Expire: now.Add(time.Second)})
	case skill.SkillIndexGiganticStorm: // 237闪电轰顶
		tobj.addEffect(&effect.Effect{BuffIndex: effect.BuffGiganticStorm, Expire: now.Add(time.Second)})
	}
}

func (obj *Object) useSkillSoulBarrier(tobj *Object, s *skill.Skill) bool {
	if (obj.Class != int(class.Wizard) && obj.Class != int(class.Magumsa)) || !tobj.canReceiveSupportSkill() {
		return false
	}
	damageReduction := 0.0
	duration := 0
	formula.WizardMagicDefense(obj.Index, tobj.Index, obj.GetDexterity(), obj.GetEnergy(), &damageReduction, &duration)
	if damageReduction <= 0 || duration <= 0 {
		return false
	}
	if !tobj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffSoulBarrier, DamageReduction: int(damageReduction),
		ManaRate: skill.Settings.SoulBarrierManaRate[0],
		Expire:   time.Now().Add(time.Duration(duration) * time.Second),
	}) {
		return false
	}
	obj.useSkillReply(tobj, s, true)
	return true
}

func (obj *Object) useSkillDefense(tobj *Object) bool {
	obj.PushViewport(&model.MsgActionReply{
		Index:  obj.Index,
		Action: int(skill.SkillIndexDefense),
		Dir:    obj.Dir,
		Target: tobj.Index,
	})
	return true
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
	obj.useSkillReply(tobj, s, true)
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
	tobj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffDefensePower,
		Defense:   int(defense),
		Expire:    time.Now().Add(time.Duration(duration) * time.Second),
	})
	obj.useSkillReply(tobj, s, true)
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
	tobj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffAttackPower,
		Attack:    int(attack),
		MagicMin:  int(attack),
		MagicMax:  int(attack),
		CurseMin:  int(attack),
		CurseMax:  int(attack),
		Expire:    time.Now().Add(time.Duration(duration) * time.Second),
	})
	obj.useSkillReply(tobj, s, true)
	return true
}

func (obj *Object) useSkillSummon(s *skill.Skill) bool {
	if obj.Class != int(class.Elf) {
		return false
	}
	monsterClass := map[int]int{
		skill.SkillIndexSummonGoblin:     26,
		skill.SkillIndexSummonStoneGolem: 32,
		skill.SkillIndexSummonAssassin:   21,
		skill.SkillIndexSummonEliteYeti:  20,
		skill.SkillIndexSummonDarkKnight: 10,
		skill.SkillIndexSummonBali:       150,
		skill.SkillIndexSummonSoldier:    151,
	}[s.Index]
	if monsterClass == 0 {
		return false
	}
	summon, err := ObjectManager.AddCallMonster(obj, monsterClass)
	if err != nil {
		slog.Error("useSkillSummon", "object", obj.Index, "skill", s.Index, "error", err)
		return false
	}
	obj.useSkillReply(summon, s, true)
	obj.Push(&model.MsgSummonHPReply{Percent: 100})
	return true
}

func (obj *Object) useSkillAttackArea(tobj *Object, s *skill.Skill, distance int) bool {
	obj.useSkillAttack(tobj, s)
	obj.ForEachViewportObject(func(vpobj *Object) {
		if vpobj != tobj &&
			vpobj.Live &&
			vpobj.Type == ObjectTypeMonster &&
			vpobj.MapNumber == obj.MapNumber &&
			tobj.CalcDistance(vpobj) <= distance {
			obj.attack(vpobj, s, 0, true)
		}
	})
	return true
}

func (obj *Object) useSkillAttackAreaSelf(tobj *Object, s *skill.Skill, distance int) bool {
	obj.useSkillAttack(tobj, s)
	obj.ForEachViewportObject(func(vpobj *Object) {
		if vpobj != tobj &&
			vpobj.Live &&
			vpobj.Type == ObjectTypeMonster &&
			vpobj.MapNumber == obj.MapNumber &&
			obj.CalcDistance(vpobj) <= distance {
			obj.attack(vpobj, s, 0, true)
		}
	})
	return true
}

func (obj *Object) useSkillAttackFrustum(tobj *Object, s *skill.Skill, tx1, ty1, tx2, ty2 float32) bool {
	obj.useSkillAttack(tobj, s)
	x := tobj.X - obj.X
	y := tobj.Y - obj.Y
	angle := float32(math.Atan2(float64(y), float64(x)))*180/math.Pi - 90
	obj.createSkillFrustum3(angle, tx1, ty1, tx2, ty2)
	obj.ForEachViewportObject(func(vpobj *Object) {
		if vpobj != tobj &&
			vpobj.Live &&
			vpobj.Type == ObjectTypeMonster &&
			vpobj.MapNumber == obj.MapNumber &&
			obj.checkSkillFrustum(vpobj) {
			obj.attack(vpobj, s, 0, true)
		}
	})
	return true
}

func (obj *Object) useSkillAttackHitBox(tobj *Object, s *skill.Skill, hitCheck func(int, int, int, int, int) bool) bool {
	obj.useSkillAttack(tobj, s)
	angle := int(math.Atan2(float64(obj.Y-tobj.Y), float64(obj.X-tobj.X))*180/math.Pi + 90)
	if angle < 0 {
		angle += 360
	}
	obj.ForEachViewportObject(func(vpobj *Object) {
		if vpobj != tobj &&
			vpobj.Live &&
			vpobj.Type == ObjectTypeMonster &&
			vpobj.MapNumber == obj.MapNumber &&
			hitCheck(angle, obj.X, obj.Y, vpobj.X, vpobj.Y) {
			obj.attack(vpobj, s, 0, true)
		}
	})
	return true
}

func (obj *Object) getAngle(tobj *Object) float32 {
	x := tobj.X - obj.X
	y := tobj.Y - obj.Y
	rad := float32(math.Atan2(float64(y), float64(x)))
	return rad*180/math.Pi + 90
}

func (obj *Object) createSkillFrustum(a, x, y float32) {
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

func (obj *Object) createSkillFrustum3(a, tx1, ty1, tx2, ty2 float32) {
	p := [4][3]float32{
		{-tx1, ty1, 0.0},
		{tx1, ty1, 0.0},
		{tx2, ty2, 0.0},
		{-tx2, ty2, 0.0},
	}
	angle := [3]float32{0.0, 0.0, a}
	matrix := math2.Angle2Matrix(angle)
	for i := 0; i < 4; i++ {
		frustum := math2.VectorRotate(p[i], matrix)
		obj.SkillFrustumX[i] = int(frustum[0]) + obj.X
		obj.SkillFrustumY[i] = int(frustum[1]) + obj.Y
	}
}

func (obj *Object) checkSkillFrustum(tobj *Object) bool {
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

func (obj *Object) useSkillDeathStab(s *skill.Skill, tobj *Object) bool {
	obj.useSkillReply(tobj, s, true)
	obj.attack(tobj, s, 0, true)
	if rand.Intn(100)%3 == 0 {
		obj.attack(tobj, s, 0, true)
	}
	angle := obj.getAngle(tobj)
	obj.createSkillFrustum(angle, 1.5, 3.0)
	obj.ForEachViewportObject(func(vpobj *Object) {
		if vpobj != tobj &&
			vpobj.Live &&
			vpobj.Type != ObjectTypePlayer &&
			vpobj.Type != ObjectTypeNPC &&
			obj.checkSkillFrustum(vpobj) {
			obj.attack(vpobj, s, 0, true)
		}
	})
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
	tobj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffSwellHP,
		MaxHP:     addHP,
		Expire:    time.Now().Add(time.Duration(duration) * time.Second),
	})
	obj.useSkillReply(tobj, s, true)
	tobj.PushMaxHPSD(tobj.MaxHP, tobj.MaxSD)
	tobj.PushHPSD(tobj.HP, tobj.SD)
	return true
}

func (obj *Object) useSkillAddCriticalDamage(s *skill.Skill) bool {
	if obj.Class != int(class.DarkLord) {
		return false
	}
	value, duration := 0, 0
	formula.DarkLordCriticalDamage(obj.GetLeadership(), obj.GetEnergy(), &value, &duration)
	if value <= 0 || duration <= 0 {
		return false
	}
	if !obj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffCriticalDamage, CriticalDamage: value,
		Expire: time.Now().Add(time.Duration(duration) * time.Second),
	}) {
		return false
	}
	obj.useSkillReply(obj, s, true)
	return true
}

func (obj *Object) useSkillInfinityArrow(s *skill.Skill) bool {
	if obj.Class != int(class.Elf) || obj.Level < skill.Settings.InfinityArrowUseLevel {
		return false
	}
	if !obj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffInfinityArrow,
		Expire:    time.Now().Add(time.Duration(skill.Settings.InfinityArrowSkillTime) * time.Second),
	}) {
		return false
	}
	obj.useSkillReply(obj, s, true)
	return true
}

func (obj *Object) useSkillDamageReflection(tobj *Object, s *skill.Skill) bool {
	if obj.Class != int(class.Summoner) || !tobj.canReceiveSupportSkill() {
		return false
	}
	value, duration := 0, 0
	formula.SummonerDamageReflect(obj.GetEnergy(), &value, &duration)
	if value <= 0 || duration <= 0 {
		return false
	}
	if !tobj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffDamageReflection, Reflect: value,
		Expire: time.Now().Add(time.Duration(duration) * time.Second),
	}) {
		return false
	}
	obj.useSkillReply(tobj, s, true)
	return true
}

func (obj *Object) useSkillBerserker(s *skill.Skill) bool {
	if obj.Class != int(class.Summoner) {
		return false
	}
	up, down, duration := 0, 0, 0
	formula.SummonerBerserker(obj.GetEnergy(), &up, &down, &duration)
	attackMin, attackMax := 0, 0
	formula.SummonerBerserkerAttackDamage(obj.GetStrength(), obj.GetDexterity(), &attackMin, &attackMax)
	magicMin, magicMax := 0.0, 0.0
	curseMin, curseMax := 0.0, 0.0
	formula.SummonerBerserkerMagicDamage(up, obj.GetEnergy(), &magicMin, &magicMax)
	formula.SummonerBerserkerCurseDamage(up, obj.GetEnergy(), &curseMin, &curseMax)
	hpRate := 40 - down
	if hpRate < 10 {
		hpRate = 10
	}
	if !obj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffBerserker, AttackMin: attackMin, AttackMax: attackMax,
		MagicMin: int(magicMin), MagicMax: int(magicMax), CurseMin: int(curseMin), CurseMax: int(curseMax),
		MaxMP: obj.MaxMP * up / 100, MaxHP: -obj.MaxHP * hpRate / 100,
		Defense: -obj.Defense * hpRate / 100,
		Expire:  time.Now().Add(time.Duration(duration) * time.Second),
	}) {
		return false
	}
	obj.useSkillReply(obj, s, true)
	obj.PushMaxHPSD(obj.MaxHP, obj.MaxSD)
	obj.PushMaxMPAG(obj.MaxMP, obj.MaxAG)
	return true
}

func (obj *Object) useSkillSleep(tobj *Object, s *skill.Skill) bool {
	if obj.Class != int(class.Summoner) || tobj.Type != ObjectTypeMonster {
		return false
	}
	rate, duration := 0, 0
	formula.SleepMonster(obj.GetEnergy(), obj.GetCurse(), tobj.Level, &rate, &duration)
	if duration <= 0 || rand.Intn(100) >= rate {
		return false
	}
	if !tobj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffSleep, Sleep: true,
		Expire: time.Now().Add(time.Duration(duration) * time.Second),
	}) {
		return false
	}
	obj.useSkillReply(tobj, s, true)
	return true
}

func (obj *Object) useSkillExpansionWizardry(s *skill.Skill) bool {
	if obj.Class != int(class.Wizard) {
		return false
	}
	value := int(float64(obj.GetEnergy()/9) * 0.20)
	if value <= 0 {
		return false
	}
	if !obj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffExpansionWizard, MagicMin: value,
		Expire: time.Now().Add(1800 * time.Second),
	}) {
		return false
	}
	obj.useSkillReply(obj, s, true)
	return true
}

func (obj *Object) useSkillRecovery(tobj *Object, s *skill.Skill) bool {
	if obj.Class != int(class.Elf) || tobj.Type != ObjectTypePlayer || tobj.SD >= tobj.MaxSD {
		return false
	}
	value := 0.0
	formula.ElfShieldRecovery(obj.GetEnergy(), obj.Level, &value)
	if value <= 0 {
		return false
	}
	tobj.SD += int(value)
	if tobj.SD > tobj.MaxSD {
		tobj.SD = tobj.MaxSD
	}
	obj.useSkillReply(tobj, s, true)
	tobj.PushHPSD(tobj.HP, tobj.SD)
	return true
}

func (obj *Object) useSkillIgnoreDefense(s *skill.Skill) bool {
	if obj.Class != int(class.RageFighter) {
		return false
	}
	value := int(float64(obj.GetEnergy()-404)/100 + 3)
	if value < 0 {
		value = 0
	}
	if value > 10 {
		value = 10
	}
	if !obj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffIgnoreDefense, IgnoreDefense: value,
		Expire: time.Now().Add(time.Duration(obj.GetEnergy()/5+60) * time.Second),
	}) {
		return false
	}
	obj.useSkillReply(obj, s, true)
	return true
}

func (obj *Object) useSkillIncreaseHealth(s *skill.Skill) bool {
	if obj.Class != int(class.RageFighter) {
		return false
	}
	value := int(float64(obj.GetEnergy()-132)/10 + 30)
	if value <= 0 {
		return false
	}
	if !obj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffIncreaseHealth, MaxHP: value,
		Expire: time.Now().Add(time.Duration(obj.GetEnergy()/5+60) * time.Second),
	}) {
		return false
	}
	obj.useSkillReply(obj, s, true)
	obj.PushMaxHPSD(obj.MaxHP, obj.MaxSD)
	return true
}

func (obj *Object) useSkillIncreaseBlock(s *skill.Skill) bool {
	if obj.Class != int(class.RageFighter) {
		return false
	}
	value := int(float64(obj.GetEnergy()-80)/10 + 10)
	if value > 100 {
		value = 100
	}
	if value <= 0 {
		return false
	}
	if !obj.addEffect(&effect.Effect{
		BuffIndex: effect.BuffIncreaseBlock, DefenseRate: value,
		Expire: time.Now().Add(time.Duration(obj.GetEnergy()/5+60) * time.Second),
	}) {
		return false
	}
	obj.useSkillReply(obj, s, true)
	return true
}

func (obj *Object) useSkillEarthshake(tobj *Object, s *skill.Skill) bool {
	if obj.GetActivePetCode() != item.Code(13, 4) ||
		tobj.Type != ObjectTypeMonster ||
		obj.CalcDistance(tobj) > 5 {
		return false
	}
	obj.useSkillReply(obj, s, true)
	obj.AddDelayMsg(4, s.Index, 500, tobj.Index)
	obj.ForEachViewportObject(func(target *Object) {
		if target != tobj &&
			target.Live &&
			target.Type == ObjectTypeMonster &&
			obj.CalcDistance(target) <= 5 {
			obj.AddDelayMsg(4, s.Index, 500, target.Index)
		}
	})
	return true
}

func (obj *Object) useSkillPlasmaStorm(tobj *Object, s *skill.Skill) bool {
	if obj.Level < 300 ||
		obj.GetActivePetCode() != item.Code(13, 37) ||
		tobj.Type != ObjectTypeMonster ||
		obj.CalcDistance(tobj) > 6 {
		return false
	}
	obj.useSkillReply(tobj, s, true)
	obj.AddDelayMsg(4, s.Index, 300, tobj.Index)
	hit := 1
	obj.ForEachViewportObject(func(target *Object) {
		if hit >= 6 ||
			target == tobj ||
			!target.Live ||
			target.Type != ObjectTypeMonster ||
			obj.CalcDistance(target) > 6 {
			return
		}
		hit++
		obj.AddDelayMsg(4, s.Index, 300, target.Index)
	})
	return hit > 0
}

func (obj *Object) useSkillDrainLife(tobj *Object, s *skill.Skill) bool {
	if obj.Class != int(class.Summoner) || tobj.Type != ObjectTypeMonster {
		return false
	}
	obj.useSkillReply(tobj, s, true)
	obj.AddDelayMsg(4, s.Index, 700, tobj.Index)
	return true
}

func (obj *Object) useSkillCurseArea(tobj *Object, s *skill.Skill) bool {
	if obj.Class != int(class.Summoner) ||
		tobj.Type != ObjectTypeMonster ||
		obj.CalcDistance(tobj) > 6 {
		return false
	}
	obj.useSkillReply(tobj, s, true)
	obj.AddDelayMsg(7, s.Index, 700, tobj.Index)
	count := 1
	obj.ForEachViewportObject(func(target *Object) {
		if count >= 5 ||
			target == tobj ||
			!target.Live ||
			target.Type != ObjectTypeMonster ||
			obj.CalcDistance(target) > 6 {
			return
		}
		count++
		obj.AddDelayMsg(7, s.Index, 700, target.Index)
	})
	return true
}

func (obj *Object) useSkillAttackAreaPoint(tobj *Object, s *skill.Skill, x, y, distance, delay int) bool {
	if tobj.Type != ObjectTypeMonster || obj.CalcDistance(tobj) > s.Distance {
		return false
	}
	obj.useSkillReply(tobj, s, true)
	if delay > 0 {
		obj.AddDelayMsg(4, s.Index, delay, tobj.Index)
	} else {
		obj.attack(tobj, s, 0, true)
	}
	obj.ForEachViewportObject(func(target *Object) {
		if target == tobj ||
			!target.Live ||
			target.Type != ObjectTypeMonster ||
			abs(target.X-x) > distance || abs(target.Y-y) > distance {
			return
		}
		if delay > 0 {
			obj.AddDelayMsg(4, s.Index, delay, target.Index)
		} else {
			obj.attack(target, s, 0, true)
		}
	})
	return true
}

func (obj *Object) useSkillCurseDebuffTarget(tobj *Object, s *skill.Skill) {
	if tobj == nil || !tobj.Live || tobj.Type != ObjectTypeMonster {
		return
	}
	rate, value, duration := 0, 0, 0
	switch s.Index {
	case skill.SkillIndexWeakness:
		formula.SummonerWeaknessMonster(obj.GetEnergy(), obj.GetCurse(), tobj.Level, &rate, &value, &duration)
	case skill.SkillIndexInnovation:
		formula.SummonerInnovationMonster(obj.GetEnergy(), obj.GetCurse(), tobj.Level, &rate, &value, &duration)
	}
	if value <= 0 || duration <= 0 || rand.Intn(100) >= rate {
		return
	}
	se := &effect.Effect{
		BuffIndex:       effect.BuffWeakness,
		AttackReduction: value,
		Expire:          time.Now().Add(time.Duration(duration) * time.Second),
	}
	if s.Index == skill.SkillIndexInnovation {
		se.BuffIndex = effect.BuffInnovation
		se.AttackReduction = 0
		se.DefenseReduction = value
	}
	tobj.addEffect(se)
}

func (obj *Object) useSkillNova(msg *model.MsgUseSkill) {
	s, ok := obj.Skills[skill.SkillIndexNova]
	if !ok || !obj.Live {
		return
	}
	if msg.Target == 58 {
		if !obj.novaSkill.startedAt.IsZero() || !obj.checkSkillDelay(s) {
			return
		}
		mp, ag := obj.GetSkillMPAG(s)
		if obj.MP < mp || obj.AG < ag {
			return
		}
		now := time.Now()
		obj.novaSkill = novaSkillState{startedAt: now, lastTick: now}
		obj.MP -= mp
		obj.AG -= ag
		obj.skillUseTimes[s.Index] = now
		obj.PushViewport(&model.MsgUseSkillReply{Index: obj.Index, Skill: 58, Target: obj.Index | 0x8000})
		obj.PushMPAG(obj.MP, obj.AG)
		return
	}
	if obj.novaSkill.startedAt.IsZero() {
		return
	}
	obj.releaseNovaSkill()
}

func (obj *Object) processNovaSkill() {
	if obj.novaSkill.startedAt.IsZero() {
		return
	}
	now := time.Now()
	if now.Sub(obj.novaSkill.startedAt) >= 6*time.Second {
		obj.releaseNovaSkill()
		return
	}
	if now.Sub(obj.novaSkill.lastTick) < 500*time.Millisecond {
		return
	}
	s := obj.Skills[skill.SkillIndexNova]
	useMP := s.ManaUsage * 20 / 100
	if obj.MP < useMP {
		obj.releaseNovaSkill()
		return
	}
	obj.MP -= useMP
	obj.novaSkill.lastTick = now
	obj.novaSkill.count++
	if obj.novaSkill.count > 12 {
		obj.novaSkill.count = 12
	}
	obj.PushMPAG(obj.MP, obj.AG)
	obj.PushViewport(&model.MsgNovaCountReply{Type: skill.SkillIndexNova, Count: obj.novaSkill.count, Index: obj.Index})
}

func (obj *Object) releaseNovaSkill() {
	if obj.novaSkill.startedAt.IsZero() {
		return
	}
	s := obj.Skills[skill.SkillIndexNova]
	obj.novaSkill.startedAt = time.Time{}
	obj.useSkillReply(obj, s, true)
	obj.ForEachViewportObject(func(target *Object) {
		if target.Live && target.Type == ObjectTypeMonster && obj.CalcDistance(target) <= s.Distance {
			obj.AddDelayMsg(4, s.Index, 600, target.Index)
		}
	})
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (obj *Object) useSkillChainLightning(tobj *Object, s *skill.Skill) bool {
	obj.chainLightningTarget = 1
	obj.useSkillAttack(tobj, s)
	hit := map[int]struct{}{tobj.Index: {}}
	for distance := 1; distance <= 2 && len(hit) < 3; distance++ {
		obj.ForEachViewportObject(func(vpobj *Object) {
			if len(hit) >= 3 {
				return
			}
			if _, ok := hit[vpobj.Index]; ok {
				return
			}
			dx := vpobj.X - tobj.X
			dy := vpobj.Y - tobj.Y
			if vpobj.Live &&
				vpobj.Type == ObjectTypeMonster &&
				vpobj.MapNumber == obj.MapNumber &&
				dx >= -distance && dx <= distance &&
				dy >= -distance && dy <= distance {
				hit[vpobj.Index] = struct{}{}
				obj.chainLightningTarget = len(hit)
				obj.attack(vpobj, s, 0, true)
			}
		})
	}
	obj.chainLightningTarget = 0
	return true
}

func (obj *Object) useSkillDarkSide(tobj *Object, s *skill.Skill) bool {
	obj.useSkillAttack(tobj, s)
	hitCount := 1
	obj.ForEachViewportObject(func(vpobj *Object) {
		if hitCount >= 5 {
			return
		}
		if vpobj != tobj &&
			vpobj.Live &&
			vpobj.Type == ObjectTypeMonster &&
			vpobj.MapNumber == obj.MapNumber &&
			obj.CalcDistance(vpobj) < s.Distance {
			hitCount++
			obj.attack(vpobj, s, 0, true)
		}
	})
	return true
}

func (obj *Object) useSkillPhoenixShot(tobj *Object, s *skill.Skill) bool {
	obj.useSkillReply(tobj, s, true)
	attack := func(target *Object) {
		for i := 0; i < 4; i++ {
			obj.attack(target, s, 0, true)
			if !target.Live {
				break
			}
		}
		obj.AddDelayMsg(2, 0, 150, target.Index)
	}
	attack(tobj)
	obj.ForEachViewportObject(func(target *Object) {
		if target != tobj && target.Live && target.Type == ObjectTypeMonster &&
			tobj.CalcDistance(target) <= 2 {
			attack(target)
		}
	})
	return true
}

func (obj *Object) canReceiveSupportSkill() bool {
	if obj.Type == ObjectTypePlayer {
		return true
	}
	return obj.Index >= ObjectManager.maxMonsterCount && obj.Index < ObjectManager.playerStartIndex
}
