package buff

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
)

func init() {
	BuffTable.init()
}

const MaxBuffCount = 32

type buffBase struct {
	Index       BuffType
	EffectType  BuffEffectType
	Type        int
	NoticeType  int
	ClearType   int
	Name        string
	Description string
}

type buffTable map[BuffType]*buffBase

var BuffTable buffTable

func (table *buffTable) init() {
	type buffConfig struct {
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
			Buffs []buffConfig `xml:"Buff"`
		} `xml:"General"`
	}

	var cfg config
	conf.XML(conf.PathCommon, "IGC_BuffEffectManager.xml", &cfg)
	buffs := make(buffTable, len(cfg.General.Buffs))
	for _, value := range cfg.General.Buffs {
		if value.Index <= 0 || value.Index > 255 || value.EffectType <= 0 || value.EffectType > 255 ||
			value.Name == "" || (value.Type != 0 && value.Type != 1) {
			slog.Error("invalid buff definition", "index", value.Index)
			os.Exit(1)
		}
		index := BuffType(value.Index)
		if _, exists := buffs[index]; exists {
			slog.Error("duplicate buff index", "index", value.Index)
			os.Exit(1)
		}
		buffs[index] = &buffBase{
			Index: index, EffectType: BuffEffectType(value.EffectType), Type: value.Type,
			NoticeType: value.NoticeType, ClearType: value.ClearType,
			Name: value.Name, Description: value.Description,
		}
	}
	*table = buffs
}

func NewBuff(index BuffType) *Buff {
	base := BuffTable[index]
	if base == nil {
		panic(fmt.Errorf("NewBuff(%d): undefined buff type", index))
	}
	return &Buff{buffBase: base}
}

// BuffEffect contains the attribute modifiers and behavior flags applied by a Buff.
type BuffEffect struct {
	Attack           int
	AttackMin        int
	AttackMax        int
	MagicMin         int
	MagicMax         int
	CurseMin         int
	CurseMax         int
	Defense          int
	MagicDefense     int
	DefenseRate      int
	AddMaxHP         int     // fixed maximum HP increase
	AddMaxMP         int     // fixed maximum MP increase
	IncMaxHP         float64 // maximum HP increase percentage
	IncMaxMP         float64 // maximum MP increase percentage
	AttackSpeed      int
	CriticalDamage   int
	IgnoreDefense    int
	Reflect          int
	DamageReduction  int
	ManaRate         int
	AttackReduction  int
	DefenseReduction int
	Slow             bool
	Sleep            bool
	Stun             bool
	Dot              Dot
}

type Dot struct {
	Damage   int
	Source   int
	NextTick time.Time
}

type Buff struct {
	*buffBase
	Effect BuffEffect
	Expire time.Time
}

type Buffs map[BuffType]*Buff

type BuffOption = BuffEffect

func (buffs *Buffs) Add(b *Buff) ([]*Buff, bool) {
	if buffs == nil || b == nil || b.buffBase == nil {
		return nil, false
	}
	removed := (*buffs).removeConflict(b.Index, b.EffectType)
	if len(*buffs) >= MaxBuffCount {
		return removed, false
	}
	if *buffs == nil {
		*buffs = make(Buffs)
	}
	(*buffs)[b.Index] = b
	return removed, true
}

func (buffs Buffs) Remove(index BuffType) *Buff {
	b := buffs[index]
	delete(buffs, index)
	return b
}

func (buffs Buffs) RemoveAll() []*Buff {
	removed := make([]*Buff, 0, len(buffs))
	for index, b := range buffs {
		removed = append(removed, b)
		delete(buffs, index)
	}
	return removed
}

func (buffs Buffs) Get(index BuffType) *Buff { return buffs[index] }

func (buffs Buffs) ActiveBuffIndexes() []int {
	indexes := make([]int, 0, len(buffs))
	for _, b := range buffs {
		indexes = append(indexes, int(b.Index))
	}
	sort.Ints(indexes)
	return indexes
}

func (buffs Buffs) ApplyBuffEffect(apply func(BuffEffect)) {
	if buffs == nil || apply == nil {
		return
	}
	var effect BuffEffect
	for _, b := range buffs {
		effect.Attack += b.Effect.Attack
		effect.AttackMin += b.Effect.AttackMin
		effect.AttackMax += b.Effect.AttackMax
		effect.MagicMin += b.Effect.MagicMin
		effect.MagicMax += b.Effect.MagicMax
		effect.CurseMin += b.Effect.CurseMin
		effect.CurseMax += b.Effect.CurseMax
		effect.Defense += b.Effect.Defense
		effect.MagicDefense += b.Effect.MagicDefense
		effect.DefenseRate += b.Effect.DefenseRate
		effect.AddMaxHP += b.Effect.AddMaxHP
		effect.AddMaxMP += b.Effect.AddMaxMP
		effect.IncMaxHP += b.Effect.IncMaxHP
		effect.IncMaxMP += b.Effect.IncMaxMP
		effect.AttackSpeed += b.Effect.AttackSpeed
		effect.CriticalDamage += b.Effect.CriticalDamage
		effect.IgnoreDefense += b.Effect.IgnoreDefense
		effect.Reflect += b.Effect.Reflect
		effect.DamageReduction += b.Effect.DamageReduction
		effect.ManaRate += b.Effect.ManaRate
		effect.AttackReduction += b.Effect.AttackReduction
		effect.DefenseReduction += b.Effect.DefenseReduction
	}
	apply(effect)
}

func (buffs Buffs) removeConflict(index BuffType, effectType BuffEffectType) []*Buff {
	var removed []*Buff
	for oldIndex, old := range buffs {
		if oldIndex == index || effectType != 0 && old.EffectType == effectType {
			removed = append(removed, old)
			delete(buffs, oldIndex)
		}
	}
	return removed
}
