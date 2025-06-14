package skill

import (
	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/class"
)

func init() {
	SkillManager.init()
}

const (
	SkillIndexPoison             int        = 1 + iota // 毒咒
	SkillIndexMeteorite                                // 陨石
	SkillIndexLightning                                // 掌心雷
	SkillIndexFireBall                                 // 火球
	SkillIndexFlame                                    // 火龙
	SkillIndexTeleport                                 // 瞬间移动
	SkillIndexIce                                      // 冰封
	SkillIndexTwister                                  // 龙卷风
	SkillIndexEvilSpirit                               // 黑龙波
	SkillIndexHellFire                                 // 地狱火
	SkillIndexPowerWave                                // 真空波
	SkillIndexAquaBeam                                 // 极光
	SkillIndexCometFall                                // 爆炎
	SkillIndexInferno                                  // 毁灭烈焰
	SkillIndexTeleportAlly                             // 小挪移
	SkillIndexSoulBarrier                              // 守护之魂
	SkillIndexEnergyBall                               // 能量球
	SkillIndexDefense                                  // 圣盾防御
	SkillIndexFallingSlash                             // 地裂斩(武器)
	SkillIndexLunge                                    // 牙突刺(武器)
	SkillIndexUppercut                                 // 升龙击(武器)
	SkillIndexCyclone                                  // 旋风斩(武器)
	SkillIndexSlash                                    // 天地十字剑(武器)
	SkillIndexTripleShot                               // 多重箭(武器)
	SkillIndexHeal               = 2 + iota            // 治疗
	SkillIndexGreaterDefense                           // 防御
	SkillIndexGreaterAttack                            // 攻击
	SkillIndexSummonGoblin       = 3 + iota            // 召唤哥布林
	SkillIndexSummonStoneGolem                         // 召唤石巨人
	SkillIndexSummonAssassin                           // 召唤暗杀者
	SkillIndexSummonEliteYeti                          // 召唤雪人王
	SkillIndexSummonDarkKnight                         // 召唤暗黑骑士
	SkillIndexSummonBali                               // 召唤巴里
	SkillIndexSummonSoldier                            // 召唤黄金斗士
	SkillIndexDecay              = 4 + iota            // 单毒炎
	SkillIndexIceStorm                                 // 暴风雪
	SkillIndexNova                                     // 星辰一怒
	SkillIndexTwistingSlash                            // 霹雳回旋斩
	SkillIndexRagefulBlow                              // 雷霆裂闪
	SkillIndexDeathStab                                // 袭风刺
	SkillIndexCrescentMoonSlash                        // 半月斩(攻城)
	SkillIndexLance                                    // 回旋刃(攻城)
	SkillIndexStarfall                                 // 天堂之箭(攻城)
	SkillIndexImpale                                   // 钻云枪
	SkillIndexSwellHP                                  // 生命之光
	SkillIndexFireBreath                               // 流星焰(彩云兽)
	SkillIndexDevilFire                                // Flame of Evil (Monster)
	SkillIndexIceArrow                                 // 冰封箭
	SkillIndexPenetration                              // 穿透箭
	SkillIndexFireSlash          = 6 + iota            // 玄月斩
	SkillIndexPowerSlash                               // 天雷闪(武器)
	SkillIndexSpiralSlash                              // 风舞回旋斩(攻城)
	SkillIndexForce              = 8 + iota            // 冲击
	SkillIndexFireBurst                                // 星云火链
	SkillIndexEarthshake                               // 地裂(黑王马)
	SkillIndexSummon                                   // 星云召唤
	SkillIndexAddCriticalDamage                        // 致命圣印
	SkillIndexElectricSpike                            // 圣极光
	SkillIndexForceWave                                // 冲击波
	SkillIndexStun                                     // Stun
	SkillIndexCancelStun                               // CancelStun
	SkillIndexSwellMP                                  // SwellMP
	SkillIndexInvisibility                             // Invisibility
	SkillIndexCancelInvisibility                       // CancelInvisibility
	SkillIndexAbolishMagic                             // AbolishMagic
	SkillIndexMPRays                                   // 幻魔光束(攻城)
	SkillIndexFireBlast                                // 神圣火焰(攻城)
	SkillIndexPlasmaStorm        = 9 + iota            // 闪电链(炎狼兽)
	SkillIndexInfinityArrow                            // 无影箭
	SkillIndexFireScream                               // 火舞旋风
	SkillIndexExplosion                                // Explosion
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
		v.ReqClass[class.Wizard] = v.DarkWizard
		v.ReqClass[class.Knight] = v.DarkKnight
		v.ReqClass[class.Elf] = v.FairyElf
		v.ReqClass[class.Magumsa] = v.MagicGladiator
		v.ReqClass[class.DarkLord] = v.DarkLord
		v.ReqClass[class.Summoner] = v.Summoner
		v.ReqClass[class.RageFighter] = v.RageFighter
		// v.ReqClass[class.GrowLancer] = v.GrowLancer
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
	for _, class := range masterSkillTree.Class {
		for _, tree := range class.Tree {
			for _, skill := range tree.Skills {
				index := skill.Index%36 - 1
				rank := index >> 2
				pos := index % 4
				m.masterSkillTable[id2class[class.ID]][tree.Type][rank][pos] = skill
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

func (m *skillManager) getMasterSkillBase(class, index int) (*MasterSkillBase, bool) {
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
