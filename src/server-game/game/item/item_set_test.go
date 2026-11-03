package item

import "testing"

func TestRandomSetIndexWithHigherTiers(t *testing.T) {
	for _, tc := range []struct {
		name string
		tier int
		set  int
	}{
		{name: "third tier", tier: 0x10, set: 31},
		{name: "fourth tier", tier: 0x20, set: 42},
	} {
		t.Run(tc.name, func(t *testing.T) {
			manager := setManager{items: []map[int]*setItem{
				{7: {sets: map[int]int{tc.tier: tc.set}}},
			}}
			if !manager.HasSetItem(0, 7) {
				t.Fatal("configured set item was not recognized")
			}
			if got := manager.RandomSetIndex(0, 7); got != tc.set {
				t.Fatalf("RandomSetIndex() = %d, want %d", got, tc.set)
			}
		})
	}
}

func TestRandomItemUsesConfiguredItems(t *testing.T) {
	for _, tc := range []struct {
		tier int
		set  int
	}{{0x10, 31}, {0x20, 42}} {
		manager := setManager{
			items: []map[int]*setItem{
				{1: {sets: map[int]int{0x01: 11}}},
				{2: {sets: map[int]int{}}},
				{3: {sets: map[int]int{tc.tier: tc.set}}},
			},
			dropRates: []int{0, 5, 10},
		}
		section, index, setIndex, ok := manager.RandomItem()
		if !ok || section != 2 || index != 3 || setIndex != tc.set {
			t.Fatalf("RandomItem() = %d, %d, %d, %t", section, index, setIndex, ok)
		}

		manager.dropRates[2] = 0
		if _, _, _, ok := manager.RandomItem(); ok {
			t.Fatal("RandomItem() selected a section without positive rate and eligible items")
		}
	}
}

func TestSetManagerCountsWornItemsAndAppliesEffects(t *testing.T) {
	manager := setManager{sets: make([]*set, 4)}
	manager.sets[1] = &set{
		count: 2,
		options: [6]SetOption{
			{Index: SetOptionIncStrength, Value: 40},
			{Index: SetOptionIncVitality, Value: 20},
		},
		fullOptions: [5]SetOption{
			{Index: SetOptionIncAttackMin, Value: 20},
			{Index: SetOptionIncAttackMax, Value: 30},
			{Index: SetOptionIncMaxAG, Value: 10},
			{Index: SetOptionIncAG, Value: 5},
			{Index: SetOptionDoubleDamage, Value: 10},
		},
	}
	manager.sets[2] = &set{count: 1, options: [6]SetOption{{Index: SetOptionIncStrength, Value: 100}}}
	manager.sets[3] = &set{count: 1, options: [6]SetOption{{Index: SetOptionIncAgility, Value: 100}}}

	items := [12]*Item{
		0:  {Set: 1, Durability: 1},
		1:  {Set: 2, Durability: 1}, // ignored because slot 0 has a set
		2:  {Set: 1, Durability: 1},
		3:  {Set: 1, Durability: 1},
		4:  {Set: 1, Durability: 0}, // broken items do not count
		8:  {Set: 3, Durability: 1},
		10: {Set: 2, Durability: 1},
		11: {Set: 3, Durability: 1}, // ignored because slot 10 has a set
	}
	var inventory Inventory
	inventory.Items = items[:]
	var option ItemSetOption
	calls := 0
	manager.ApplySetItemOption(&inventory, func(got ItemSetOption) {
		calls++
		option = got
	})
	if calls != 1 {
		t.Fatalf("option callback called %d times, want once", calls)
	}

	want := ItemSetOption{Strength: 40, Vitality: 20, AttackMin: 20, AttackMax: 30, MaxAG: 15, DoubleDamageRate: 10}
	if option != want {
		t.Fatalf("set options = %+v, want %+v", option, want)
	}
}
