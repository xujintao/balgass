package skill

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/class"
)

func init() {
	SkillManager.init()
}

const (
	SkillIndexPoison             = 1   // 毒咒
	SkillIndexMeteorite          = 2   // 陨石
	SkillIndexLightning          = 3   // 掌心雷
	SkillIndexFireBall           = 4   // 火球
	SkillIndexFlame              = 5   // 火龙
	SkillIndexTeleport           = 6   // 瞬间移动
	SkillIndexIce                = 7   // 冰封
	SkillIndexTwister            = 8   // 龙卷风
	SkillIndexEvilSpirit         = 9   // 黑龙波
	SkillIndexHellFire           = 10  // 地狱火
	SkillIndexPowerWave          = 11  // 真空波
	SkillIndexAquaBeam           = 12  // 极光
	SkillIndexCometFall          = 13  // 爆炎
	SkillIndexInferno            = 14  // 毁灭烈焰
	SkillIndexTeleportAlly       = 15  // 小挪移
	SkillIndexSoulBarrier        = 16  // 守护之魂
	SkillIndexEnergyBall         = 17  // 能量球
	SkillIndexDefense            = 18  // 圣盾防御
	SkillIndexFallingSlash       = 19  // 地裂斩(武器)
	SkillIndexLunge              = 20  // 牙突刺(武器)
	SkillIndexUppercut           = 21  // 升龙击(武器)
	SkillIndexCyclone            = 22  // 旋风斩(武器)
	SkillIndexSlash              = 23  // 天地十字剑(武器)
	SkillIndexTripleShot         = 24  // 多重箭(武器)
	SkillIndexHeal               = 26  // 治疗
	SkillIndexGreaterDefense     = 27  // 防御
	SkillIndexGreaterAttack      = 28  // 攻击
	SkillIndexSummonGoblin       = 30  // 召唤哥布林
	SkillIndexSummonStoneGolem   = 31  // 召唤石巨人
	SkillIndexSummonAssassin     = 32  // 召唤暗杀者
	SkillIndexSummonEliteYeti    = 33  // 召唤雪人王
	SkillIndexSummonDarkKnight   = 34  // 召唤暗黑骑士
	SkillIndexSummonBali         = 35  // 召唤巴里
	SkillIndexSummonSoldier      = 36  // 召唤黄金斗士
	SkillIndexDecay              = 38  // 单毒炎
	SkillIndexIceStorm           = 39  // 暴风雪
	SkillIndexNova               = 40  // 星辰一怒
	SkillIndexTwistingSlash      = 41  // 霹雳回旋斩
	SkillIndexRagefulBlow        = 42  // 雷霆裂闪
	SkillIndexDeathStab          = 43  // 袭风刺
	SkillIndexCrescentMoonSlash  = 44  // 半月斩(攻城)
	SkillIndexLance              = 45  // 回旋刃(攻城)
	SkillIndexStarfall           = 46  // 天堂之箭(攻城)
	SkillIndexImpale             = 47  // 钻云枪
	SkillIndexSwellHP            = 48  // 生命之光
	SkillIndexFireBreath         = 49  // 流星焰(彩云兽)
	SkillIndexDevilFire          = 50  // Flame of Evil (Monster)
	SkillIndexIceArrow           = 51  // 冰封箭
	SkillIndexPenetration        = 52  // 穿透箭
	SkillIndexFireSlash          = 55  // 玄月斩
	SkillIndexPowerSlash         = 56  // 天雷闪(武器)
	SkillIndexSpiralSlash        = 57  // 风舞回旋斩(攻城)
	SkillIndexForce              = 60  // 冲击
	SkillIndexFireBurst          = 61  // 星云火链
	SkillIndexEarthshake         = 62  // 地裂(黑王马)
	SkillIndexSummon             = 63  // 星云召唤
	SkillIndexAddCriticalDamage  = 64  // 致命圣印
	SkillIndexElectricSpike      = 65  // 圣极光
	SkillIndexForceWave          = 66  // 冲击波
	SkillIndexStun               = 67  // Stun
	SkillIndexCancelStun         = 68  // CancelStun
	SkillIndexSwellMP            = 69  // SwellMP
	SkillIndexInvisibility       = 70  // Invisibility
	SkillIndexCancelInvisibility = 71  // CancelInvisibility
	SkillIndexAbolishMagic       = 72  // AbolishMagic
	SkillIndexMPRays             = 73  // 幻魔光束(攻城)
	SkillIndexFireBlast          = 74  // 神圣火焰(攻城)
	SkillIndexPlasmaStorm        = 76  // 闪电链(炎狼兽)
	SkillIndexInfinityArrow      = 77  // 无影箭
	SkillIndexFireScream         = 78  // 火舞旋风
	SkillIndexExplosion          = 79  // Explosion
	SkillIndexSummonMonster      = 200 // Summon Monster
	SkillIndexMagicImmunity      = 201 // Magic Attack Immunity
	SkillIndexPhysicalImmunity   = 202 // Physical Attack Immunity
	SkillIndexPotionOfBless      = 203 // Potion of Bless
	SkillIndexPotionOfSoul       = 204 // Potion of Soul
	SkillIndexSpellOfProtection  = 210 // Spell of Protection
	SkillIndexSpellOfRestriction = 211 // Spell of Restriction
	SkillIndexSpellOfPursuit     = 212 // Spell of Pursuit
	SkillIndexShieldBurn         = 213 // Shied-Burn
	SkillIndexDrainLife          = 214 // 摄魂咒
	SkillIndexChainLightning     = 215 // 链雷咒
	SkillIndexDamageReflection   = 217 // 伤害反射
	SkillIndexBerserker          = 218 // 狂暴术
	SkillIndexSleep              = 219 // 昏睡
	SkillIndexWeakness           = 221 // 虚弱阵
	SkillIndexInnovation         = 222 // 破御阵
	SkillIndexSummonerExplosion  = 223 // 爆裂
	SkillIndexRequiem            = 224 // 刺袭
	SkillIndexPollution          = 225 // 污染
	SkillIndexLightningShock     = 230 // 烈光闪
	SkillIndexStrikeDestruction  = 232 // 破坏一击
	SkillIndexExpansionWizardry  = 233 // 法神附体
	SkillIndexRecovery           = 234 // 防护值恢复
	SkillIndexMultiShot          = 235 // 五重箭
	SkillIndexFlameStrike        = 236 // 火剑袭
	SkillIndexGiganticStorm      = 237 // 闪电轰顶
	SkillIndexChaoticDiseier     = 238 // 黑暗之力
	SkillIndexDoppelgangerSelf   = 239 // Doppelganger Self Explosion Skill
	SkillIndexKillingBlow        = 260 // 幽冥青狼拳
	SkillIndexBeastUppercut      = 261 // 斗气爆裂拳
	SkillIndexChainDrive         = 262 // 回旋踢
	SkillIndexDarkSide           = 263 // 幽冥光速拳
	SkillIndexDragonRoar         = 264 // 炎龙拳
	SkillIndexDragonSlasher      = 265 // 噬血之龙
	SkillIndexIgnoreDefense      = 266 // 斗神-破
	SkillIndexIncreaseHealth     = 267 // 斗神-命
	SkillIndexIncreaseBlock      = 268 // 斗神-御
	SkillIndexCharge             = 269 // 冲锋(攻城)
	SkillIndexPhoenixShot        = 270 // 神圣气旋
)

var SkillManager skillManager

type SkillBase struct {
	Index          int    `xml:"Index,attr"`
	Name           string `xml:"Name,attr"`
	ReqLevel       int    `xml:"ReqLevel,attr"`
	Damage         int    `xml:"Damage,attr"`
	STID           int    `xml:"STID,attr"`
	ManaUsage      int    `xml:"ManaUsage,attr"`
	BPUsage        int    `xml:"BPUsage,attr"`
	Distance       int    `xml:"Distance,attr"`
	Delay          int    `xml:"Delay,attr"`
	ReqStrength    int    `xml:"ReqStrength,attr"`
	ReqDexterity   int    `xml:"ReqDexterity,attr"`
	ReqEnergy      int    `xml:"ReqEnergy,attr"`
	ReqCommand     int    `xml:"ReqCommand,attr"`
	ReqMLPoint     int    `xml:"ReqMLPoint,attr"`
	Attribute      int    `xml:"Attribute,attr"`
	Type           int    `xml:"Type,attr"`
	UseType        int    `xml:"UseType,attr"`
	Brand          int    `xml:"Brand,attr"`
	KillCount      int    `xml:"KillCount,attr"`
	ReqStatus0     int    `xml:"ReqStatus0,attr"`
	ReqStatus1     int    `xml:"ReqStatus1,attr"`
	ReqStatus2     int    `xml:"ReqStatus2,attr"`
	DarkWizard     int    `xml:"DarkWizard,attr"`
	DarkKnight     int    `xml:"DarkKnight,attr"`
	FairyElf       int    `xml:"FairyElf,attr"`
	MagicGladiator int    `xml:"MagicGladiator,attr"`
	DarkLord       int    `xml:"DarkLord,attr"`
	Summoner       int    `xml:"Summoner,attr"`
	RageFighter    int    `xml:"RageFighter,attr"`
	GrowLancer     int    `xml:"GrowLancer,attr"`
	ReqClass       [8]int `xml:"-"`
	Rank           int    `xml:"Rank,attr"`
	Group          int    `xml:"Group,attr"`
	HP             int    `xml:"HP,attr"`
	SD             int    `xml:"SD,attr"`
	Duration       int    `xml:"Duration,attr"`
	IconNumber     int    `xml:"IconNumber,attr"`
	ItemSkill      bool   `xml:"ItemSkill,attr"`
	IsDamage       int    `xml:"isDamage,attr"`
	BuffIndex      int    `xml:"BuffIndex,attr"`
}

type MasterSkillBase struct {
	Index        int    `xml:"Index,attr"`
	ReqMinPoint  int    `xml:"ReqMinPoint,attr"`
	MaxPoint     int    `xml:"MaxPoint,attr"`
	ParentSkill1 int    `xml:"ParentSkill1,attr"`
	ParentSkill2 int    `xml:"ParentSkill2,attr"`
	SkillID      int    `xml:"MagicNumber,attr"`
	Name         string `xml:"Name,attr"`
}

type valueType int

const (
	valueTypeNormal = iota
	valueTypeDamage
	valueTypeManaInc
)

type masterSkillValue struct {
	valueType valueType
	values    [21]float32
}

type skillManager struct {
	skillTable            map[int]*SkillBase
	masterSkillTable      [8][3][9][4]*MasterSkillBase
	masterSkillValueTable [30]masterSkillValue
}

func (m *skillManager) init() {
	// array -> map
	type skillListConfig struct {
		Skills []*SkillBase `xml:"Skill"`
	}
	var skillList skillListConfig
	conf.XML(conf.PathCommon, "Skills/IGC_SkillList.xml", &skillList)
	m.skillTable = make(map[int]*SkillBase)
	for _, v := range skillList.Skills {
		if v == nil {
			m.fatalf("nil skill entry in Skills/IGC_SkillList.xml")
		}
		if v.Index <= 0 {
			m.fatalf("invalid skill index %d", v.Index)
		}
		if _, ok := m.skillTable[v.Index]; ok {
			m.fatalf("duplicate skill index %d", v.Index)
		}
		v.ReqClass[class.Wizard] = v.DarkWizard
		v.ReqClass[class.Knight] = v.DarkKnight
		v.ReqClass[class.Elf] = v.FairyElf
		v.ReqClass[class.Magumsa] = v.MagicGladiator
		v.ReqClass[class.DarkLord] = v.DarkLord
		v.ReqClass[class.Summoner] = v.Summoner
		v.ReqClass[class.RageFighter] = v.RageFighter
		v.ReqClass[class.GrowLancer] = v.GrowLancer
		m.skillTable[v.Index] = v
	}

	// array -> map
	type MasterSkillTree struct {
		Class []struct {
			ID   int `xml:"ID,attr"`
			Tree []struct {
				Type   int                `xml:"Type,attr"`
				Skills []*MasterSkillBase `xml:"Skill"`
			} `xml:"Tree"`
		} `xml:"Class"`
	}
	var masterSkillTree MasterSkillTree
	conf.XML(conf.PathCommon, "IGC_MasterSkillTree.xml", &masterSkillTree)
	id2class := map[int]class.Class{
		1:   class.Knight,
		2:   class.Wizard,
		4:   class.Elf,
		8:   class.Summoner,
		16:  class.Magumsa,
		32:  class.DarkLord,
		64:  class.RageFighter,
		128: class.GrowLancer,
	}
	// m.masterSkillTable = make(map[int]*MasterSkillBase)
	for _, classNode := range masterSkillTree.Class {
		classID, ok := id2class[classNode.ID]
		if !ok {
			m.fatalf("invalid master skill class id %d", classNode.ID)
		}
		for _, tree := range classNode.Tree {
			if tree.Type < 0 || tree.Type >= 3 {
				m.fatalf("invalid master skill tree type %d for class id %d", tree.Type, classNode.ID)
			}
			for _, skill := range tree.Skills {
				if skill == nil {
					m.fatalf("nil master skill entry for class id %d tree type %d", classNode.ID, tree.Type)
				}
				if _, ok := m.skillTable[skill.SkillID]; !ok {
					m.fatalf("master skill %d references missing skill id %d", skill.Index, skill.SkillID)
				}
				if skill.ReqMinPoint <= 0 {
					m.fatalf("master skill %d has invalid ReqMinPoint %d", skill.Index, skill.ReqMinPoint)
				}
				if skill.MaxPoint <= 0 || skill.ReqMinPoint > skill.MaxPoint {
					m.fatalf("master skill %d has invalid MaxPoint %d", skill.Index, skill.MaxPoint)
				}
				if skill.ParentSkill1 != 0 {
					if _, ok := m.skillTable[skill.ParentSkill1]; !ok {
						m.fatalf("master skill %d references missing parent skill %d", skill.SkillID, skill.ParentSkill1)
					}
				}
				if skill.ParentSkill2 != 0 {
					if _, ok := m.skillTable[skill.ParentSkill2]; !ok {
						m.fatalf("master skill %d references missing parent skill %d", skill.SkillID, skill.ParentSkill2)
					}
				}
				index := skill.Index%36 - 1
				if index < 0 {
					m.fatalf("master skill %d has invalid index %d", skill.SkillID, skill.Index)
				}
				rank := index >> 2
				pos := index % 4
				if rank < 0 || rank >= 9 || pos < 0 || pos >= 4 {
					m.fatalf("master skill %d maps to invalid rank/pos %d/%d", skill.SkillID, rank, pos)
				}
				if m.masterSkillTable[classID][tree.Type][rank][pos] != nil {
					m.fatalf("duplicate master skill slot class %d tree %d rank %d pos %d", classID, tree.Type, rank, pos)
				}
				m.masterSkillTable[classID][tree.Type][rank][pos] = skill
			}
		}
	}
	// for t := 0; t < 3; t++ {
	// 	for rank := 0; rank < 9; rank++ {
	// 		for pos := 0; pos < 4; pos++ {
	// 			skill := m.table[class.Knight][t][rank][pos]
	// 			if skill == nil {
	// 				fmt.Print("[null]")
	// 				fmt.Print("\t")
	// 				continue
	// 			}
	// 			fmt.Print(skill.Name)
	// 			fmt.Print("\t")
	// 		}
	// 		fmt.Println()
	// 	}
	// }
	// fmt.Println(1)

	// fulfill masterSkillVauleTable by lua script
}

func (m *skillManager) fatalf(format string, args ...any) {
	err := fmt.Errorf(format, args...)
	slog.Error("load skill config", "err", err)
	os.Exit(1)
}

func (m *skillManager) getMasterSkillBase(class, index int) (*MasterSkillBase, bool) {
	if class < 0 || class >= len(m.masterSkillTable) {
		return nil, false
	}
	for t := 0; t < 3; t++ {
		for rank := 0; rank < 9; rank++ {
			for pos := 0; pos < 4; pos++ {
				base := m.masterSkillTable[class][t][rank][pos]
				if base != nil && base.SkillID == index {
					return base, true
				}
			}
		}
	}
	return nil, false
}
