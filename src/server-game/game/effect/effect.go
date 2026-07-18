package effect

import (
	"log/slog"
	"os"
	"sort"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
)

const (
	MaxBuffEffect        = 32
	BuffAttackPower      = 1
	BuffDefensePower     = 2
	BuffSoulBarrier      = 4
	BuffCriticalDamage   = 5
	BuffInfinityArrow    = 6
	BuffSwellHP          = 8
	BuffPoison           = 55
	BuffIce              = 56
	BuffIceArrow         = 57
	BuffDamageReflection = 71
	BuffSleep            = 72
	BuffRequiem          = 74
	BuffExplosion        = 75
	BuffWeakness         = 76
	BuffInnovation       = 77
	BuffBerserker        = 81
	BuffExpansionWizard  = 82
	BuffFlameStrike      = 83
	BuffGiganticStorm    = 84
	BuffLightningShock   = 85
	BuffCold             = 86
	BuffIgnoreDefense    = 129
	BuffIncreaseHealth   = 130
	BuffIncreaseBlock    = 131
)

type BuffBase struct {
	Index       int
	EffectType  int
	Type        int
	NoticeType  int
	ClearType   int
	Name        string
	Description string
}

type Effect struct {
	BuffIndex        int
	Category         int
	EffectType       int
	Attack           int
	AttackMin        int
	AttackMax        int
	MagicMin         int
	MagicMax         int
	CurseMin         int
	CurseMax         int
	Defense          int
	DefenseRate      int
	MaxHP            int
	MaxMP            int
	CriticalDamage   int
	IgnoreDefense    int
	Reflect          int
	DamageReduction  int
	ManaRate         int
	AttackReduction  int
	DefenseReduction int
	Dot              int
	Source           int
	Slow             bool
	Sleep            bool
	Expire           time.Time
	NextTick         time.Time
}

type Effects map[int]*Effect

var BuffTable map[int]*BuffBase

func init() {
	initBuffTable()
}

func initBuffTable() {
	type buff struct {
		Index       int    `xml:"Index,attr"`
		EffectType  int    `xml:"EffectType,attr"`
		Type        int    `xml:"Type,attr"`
		NoticeType  int    `xml:"NoticeType,attr"`
		ClearType   int    `xml:"ClearType,attr"`
		Name        string `xml:"Name,attr"`
		Description string `xml:"Description,attr"`
	}
	type config struct {
		General struct {
			Buffs []buff `xml:"Buff"`
		} `xml:"General"`
	}

	var cfg config
	conf.XML(conf.PathCommon, "IGC_BuffEffectManager.xml", &cfg)
	BuffTable = make(map[int]*BuffBase, len(cfg.General.Buffs))
	for _, v := range cfg.General.Buffs {
		if v.Index <= 0 || v.Index > 255 || v.EffectType <= 0 || v.Name == "" ||
			(v.Type != 0 && v.Type != 1) {
			fatalEffectConfig("invalid buff definition", "index", v.Index)
		}
		if _, ok := BuffTable[v.Index]; ok {
			fatalEffectConfig("duplicate buff index", "index", v.Index)
		}
		BuffTable[v.Index] = &BuffBase{
			Index:       v.Index,
			EffectType:  v.EffectType,
			Type:        v.Type,
			NoticeType:  v.NoticeType,
			ClearType:   v.ClearType,
			Name:        v.Name,
			Description: v.Description,
		}
	}
}

func (effects Effects) Add(effect *Effect) ([]*Effect, bool) {
	if effect == nil || effect.BuffIndex == 0 {
		return nil, false
	}
	base := BuffTable[effect.BuffIndex]
	if base != nil {
		effect.Category = base.EffectType
	}
	removed := effects.removeConflict(effect.BuffIndex, effect.Category)
	if len(effects) >= MaxBuffEffect {
		return removed, false
	}
	effects[effect.BuffIndex] = effect
	return removed, true
}

func (effects Effects) Remove(index int) *Effect {
	effect := effects[index]
	if effect == nil {
		return nil
	}
	delete(effects, index)
	return effect
}

func (effects Effects) RemoveAll() []*Effect {
	removed := make([]*Effect, 0, len(effects))
	for index, effect := range effects {
		removed = append(removed, effect)
		delete(effects, index)
	}
	return removed
}

func (effects Effects) Get(index int) *Effect {
	return effects[index]
}

func (effects Effects) DefenseReduction() int {
	total := 0
	for _, effect := range effects {
		total += effect.DefenseReduction
	}
	return total
}

func (effects Effects) MagicMin() int {
	total := 0
	for _, effect := range effects {
		total += effect.MagicMin
	}
	return total
}

func (effects Effects) MagicMax() int {
	total := 0
	for _, effect := range effects {
		total += effect.MagicMax
	}
	return total
}

func (effects Effects) CurseMin() int {
	total := 0
	for _, effect := range effects {
		total += effect.CurseMin
	}
	return total
}

func (effects Effects) CurseMax() int {
	total := 0
	for _, effect := range effects {
		total += effect.CurseMax
	}
	return total
}

func (effects Effects) CriticalDamage() int {
	total := 0
	for _, effect := range effects {
		total += effect.CriticalDamage
	}
	return total
}

func (effects Effects) IgnoreDefense() int {
	total := 0
	for _, effect := range effects {
		total += effect.IgnoreDefense
	}
	return total
}

func (effects Effects) AttackReduction() int {
	total := 0
	for _, effect := range effects {
		total += effect.AttackReduction
	}
	return total
}

func (effects Effects) Reflect() int {
	total := 0
	for _, effect := range effects {
		total += effect.Reflect
	}
	return total
}

func (effects Effects) ActiveBuffIndexes() []int {
	buffs := make([]int, 0, len(effects))
	for _, effect := range effects {
		buffs = append(buffs, effect.BuffIndex)
	}
	sort.Ints(buffs)
	return buffs
}

func (effects Effects) removeConflict(index, category int) []*Effect {
	var removed []*Effect
	for oldIndex, old := range effects {
		if oldIndex == index || category != 0 && old.Category == category {
			removed = append(removed, old)
			delete(effects, oldIndex)
		}
	}
	return removed
}

func fatalEffectConfig(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
