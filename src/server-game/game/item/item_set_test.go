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
