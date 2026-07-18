package skill

import (
	"log/slog"
	"os"

	"github.com/xujintao/balgass/src/server-game/conf"
)

type SkillSettings struct {
	CanWizardUseSkillsWhileTeleport bool
	InfinityArrowSkillTime          int
	InfinityArrowUseLevel           int
	InfinityArrowMPConsumption      [4]int
	FireScreamExplosionDistance     int
	FireScreamExplosionDamage       int
	FireScreamExplosionRate         int
	FireScreamMaxAttackCount        int
	SoulBarrierManaRate             [21]int
	IceArrowTime                    int
}

var Settings SkillSettings

func init() {
	initSkillSettings()
}

func initSkillSettings() {
	type skillInfo struct {
		CanWizardUseSkillsWhileTeleport int `ini:"CanDarkWizardUseSkillsWhileTeleport"`
		InfinityArrowSkillTime          int `ini:"InfinityArrowSkillTime"`
		InfinityArrowUseLevel           int `ini:"InfinityArrowUseLevel"`
		InfinityArrowMP0                int `ini:"InfinityArraowMPConsumptionPlus0"`
		InfinityArrowMP1                int `ini:"InfinityArraowMPConsumptionPlus1"`
		InfinityArrowMP2                int `ini:"InfinityArraowMPConsumptionPlus2"`
		InfinityArrowMP3                int `ini:"InfinityArraowMPConsumptionPlus3"`
		FireScreamExplosionDistance     int `ini:"FireScreamExplosionAttackDistance"`
		FireScreamExplosionDamage       int `ini:"FireScreamExplosionDamage"`
		FireScreamExplosionRate         int `ini:"FireScreamExplosionRate"`
		FireScreamMaxAttackCount        int `ini:"FireScreamMaxAttackCountSameSerial"`
		IceArrowTime                    int `ini:"IceArrowTime"`
		SoulBarrierManaRate0            int `ini:"SoulBarrierManaRate_Level0"`
		SoulBarrierManaRate1            int `ini:"SoulBarrierManaRate_Level1"`
		SoulBarrierManaRate2            int `ini:"SoulBarrierManaRate_Level2"`
		SoulBarrierManaRate3            int `ini:"SoulBarrierManaRate_Level3"`
		SoulBarrierManaRate4            int `ini:"SoulBarrierManaRate_Level4"`
		SoulBarrierManaRate5            int `ini:"SoulBarrierManaRate_Level5"`
		SoulBarrierManaRate6            int `ini:"SoulBarrierManaRate_Level6"`
		SoulBarrierManaRate7            int `ini:"SoulBarrierManaRate_Level7"`
		SoulBarrierManaRate8            int `ini:"SoulBarrierManaRate_Level8"`
		SoulBarrierManaRate9            int `ini:"SoulBarrierManaRate_Level9"`
		SoulBarrierManaRate10           int `ini:"SoulBarrierManaRate_Level10"`
		SoulBarrierManaRate11           int `ini:"SoulBarrierManaRate_Level11"`
		SoulBarrierManaRate12           int `ini:"SoulBarrierManaRate_Level12"`
		SoulBarrierManaRate13           int `ini:"SoulBarrierManaRate_Level13"`
		SoulBarrierManaRate14           int `ini:"SoulBarrierManaRate_Level14"`
		SoulBarrierManaRate15           int `ini:"SoulBarrierManaRate_Level15"`
		SoulBarrierManaRate16           int `ini:"SoulBarrierManaRate_Level16"`
		SoulBarrierManaRate17           int `ini:"SoulBarrierManaRate_Level17"`
		SoulBarrierManaRate18           int `ini:"SoulBarrierManaRate_Level18"`
		SoulBarrierManaRate19           int `ini:"SoulBarrierManaRate_Level19"`
		SoulBarrierManaRate20           int `ini:"SoulBarrierManaRate_Level20"`
	}
	type config struct {
		SkillInfo skillInfo `ini:"SkillInfo"`
	}

	var cfg config
	conf.INI(conf.PathCommon, "Skills/IGC_SkillSettings.ini", &cfg)
	v := cfg.SkillInfo
	Settings = SkillSettings{
		CanWizardUseSkillsWhileTeleport: v.CanWizardUseSkillsWhileTeleport != 0,
		InfinityArrowSkillTime:          v.InfinityArrowSkillTime,
		InfinityArrowUseLevel:           v.InfinityArrowUseLevel,
		InfinityArrowMPConsumption:      [4]int{v.InfinityArrowMP0, v.InfinityArrowMP1, v.InfinityArrowMP2, v.InfinityArrowMP3},
		FireScreamExplosionDistance:     v.FireScreamExplosionDistance,
		FireScreamExplosionDamage:       v.FireScreamExplosionDamage,
		FireScreamExplosionRate:         v.FireScreamExplosionRate,
		FireScreamMaxAttackCount:        v.FireScreamMaxAttackCount,
		IceArrowTime:                    v.IceArrowTime,
		SoulBarrierManaRate: [21]int{
			v.SoulBarrierManaRate0, v.SoulBarrierManaRate1, v.SoulBarrierManaRate2,
			v.SoulBarrierManaRate3, v.SoulBarrierManaRate4, v.SoulBarrierManaRate5,
			v.SoulBarrierManaRate6, v.SoulBarrierManaRate7, v.SoulBarrierManaRate8,
			v.SoulBarrierManaRate9, v.SoulBarrierManaRate10, v.SoulBarrierManaRate11,
			v.SoulBarrierManaRate12, v.SoulBarrierManaRate13, v.SoulBarrierManaRate14,
			v.SoulBarrierManaRate15, v.SoulBarrierManaRate16, v.SoulBarrierManaRate17,
			v.SoulBarrierManaRate18, v.SoulBarrierManaRate19, v.SoulBarrierManaRate20,
		},
	}
	if Settings.InfinityArrowSkillTime <= 0 || Settings.InfinityArrowUseLevel <= 0 ||
		Settings.IceArrowTime <= 0 ||
		Settings.FireScreamExplosionDistance <= 0 ||
		Settings.FireScreamExplosionDamage <= 0 ||
		Settings.FireScreamExplosionRate <= 0 ||
		Settings.FireScreamExplosionRate > 10000 ||
		Settings.FireScreamMaxAttackCount <= 0 {
		fatalSkillConfig("invalid skill settings")
	}
	for level, mp := range Settings.InfinityArrowMPConsumption {
		if mp < 0 {
			fatalSkillConfig("invalid infinity arrow MP consumption", "level", level, "mp", mp)
		}
	}
	for i, rate := range Settings.SoulBarrierManaRate {
		if rate <= 0 {
			fatalSkillConfig("invalid soul barrier mana rate", "level", i, "rate", rate)
		}
	}
}

func fatalSkillConfig(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
