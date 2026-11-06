package buff

import "testing"

func TestNewBuffUsesConfiguredBase(t *testing.T) {
	active := NewBuff(BuffTypeIncreaseAttackPower)
	if active.buffBase == nil || active.Index != BuffTypeIncreaseAttackPower ||
		active.EffectType != BuffEffectTypeIncreaseAttackPower || active.Name == "" {
		t.Fatalf("new buff did not embed configured base: %#v", active)
	}
}

func TestNewBuffPanicsForUnknownType(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("NewBuff did not panic for an unknown buff type")
		}
	}()
	NewBuff(BuffType(900))
}

func TestApplyBuffEffectAggregatesMultipleEffects(t *testing.T) {
	first, second := NewBuff(BuffTypeIncreaseAttackPower), NewBuff(BuffTypeIncreaseDefense)
	first.Effect = BuffEffect{Attack: 5, AttackMin: 2, AttackMax: 4, Defense: 3, MagicDefense: 9, DefenseRate: 6, AddMaxHP: 500, IncMaxHP: 10.5, IncMaxMP: 7.25, AttackSpeed: 15, Slow: true}
	second.Effect = BuffEffect{Attack: 8, Defense: 1, AddMaxMP: 300, IncMaxHP: -2.5, Dot: Dot{Damage: 4, Source: 7}}
	buffs := Buffs{first.Index: first, second.Index: second}
	var got BuffEffect
	calls := 0
	buffs.ApplyBuffEffect(func(effect BuffEffect) {
		calls++
		got = effect
	})
	want := BuffEffect{Attack: 13, AttackMin: 2, AttackMax: 4, Defense: 4, MagicDefense: 9, DefenseRate: 6, AddMaxHP: 500, AddMaxMP: 300, IncMaxHP: 8, IncMaxMP: 7.25, AttackSpeed: 15}
	if got != want {
		t.Fatalf("aggregated option = %+v, want %+v", got, want)
	}
	if calls != 1 {
		t.Fatalf("ApplyBuffEffect callback calls = %d, want 1", calls)
	}
}

func TestApplyBuffEffectEmptyInputs(t *testing.T) {
	var nilBuffs Buffs
	calls := 0
	nilBuffs.ApplyBuffEffect(func(BuffEffect) { calls++ })
	var got BuffEffect
	Buffs{}.ApplyBuffEffect(func(effect BuffEffect) {
		calls++
		got = effect
	})
	active := NewBuff(BuffTypeIncreaseAttackPower)
	Buffs{active.Index: active}.ApplyBuffEffect(func(effect BuffEffect) {
		calls++
		got = effect
	})
	Buffs{BuffTypeIncreaseAttackPower: NewBuff(BuffTypeIncreaseAttackPower)}.ApplyBuffEffect(nil)
	if calls != 2 || got != (BuffEffect{}) {
		t.Fatalf("zero-effect Buffs callback = %+v, calls = %d; want zero effect and two calls", got, calls)
	}
}

func TestBuffContainerAddReplaceRemove(t *testing.T) {
	buffs := make(Buffs)
	first := NewBuff(BuffTypeIncreaseAttackPower)
	if removed, ok := buffs.Add(first); !ok || len(removed) != 0 {
		t.Fatalf("add first buff = %v, %v", removed, ok)
	}
	second := NewBuff(BuffTypeBerserker81)
	if removed, ok := buffs.Add(second); !ok || len(removed) != 1 || removed[0] != first {
		t.Fatalf("replace same-category buff = %v, %v", removed, ok)
	}
	if buffs.Get(first.Index) != nil || buffs.Get(second.Index) != second {
		t.Fatal("container lookup after replacement is incorrect")
	}
	if got := buffs.Remove(second.Index); got != second || len(buffs) != 0 {
		t.Fatal("container remove did not return and remove the buff")
	}
}
