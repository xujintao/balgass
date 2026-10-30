package item

import "testing"

func TestConfigured380Effects(t *testing.T) {
	inventory := Inventory{}
	inventory.Items = make([]*Item, INVENTORY_WEAR_SIZE)
	for position, code := range [][2]int{{0, 22}, {7, 29}, {8, 29}, {9, 29}, {10, 29}, {11, 29}} {
		section, index := code[0], code[1]
		if !Item380Manager.Is380Item(section, index) {
			t.Fatalf("configured 380 item %d/%d was not loaded", section, index)
		}
		inventory.Items[position] = &Item{
			Section: section, Index: index, Durability: 1, Option380: true,
		}
	}
	var got Item380Effect
	calls := 0
	Item380Manager.Apply380ItemEffect(&inventory, func(effect Item380Effect) {
		calls++
		got = effect
	})
	if calls != 1 {
		t.Fatalf("effect callback called %d times, want once", calls)
	}
	want := Item380Effect{
		IncPVPAttackRate:  10,
		IncPVPAttack:      200,
		IncPVPDefenseRate: 50,
		IncPVPDefense:     200,
		IncMaxHP:          200,
		IncMaxSD:          700,
		AutoRecoverySD:    true,
		IncRecoverySDRate: 20,
	}
	if got != want {
		t.Fatalf("configured 380 effects = %+v, want %+v", got, want)
	}
}

func TestApply380ItemEffectIgnoresIneligibleItems(t *testing.T) {
	inventory := Inventory{}
	inventory.Items = make([]*Item, INVENTORY_WEAR_SIZE)
	inventory.Items[1] = &Item{Section: 0, Index: 22, Durability: 1}
	inventory.Items[2] = &Item{Section: 0, Index: 22, Option380: true}
	inventory.Items[3] = &Item{Section: -1, Index: -1, Durability: 1, Option380: true}

	var got Item380Effect
	calls := 0
	Item380Manager.Apply380ItemEffect(&inventory, func(effect Item380Effect) {
		calls++
		got = effect
	})
	if calls != 1 {
		t.Fatalf("effect callback called %d times, want once", calls)
	}
	if got != (Item380Effect{}) {
		t.Fatalf("ineligible item effects = %+v, want zero value", got)
	}
}

func Test380MixSuccessRateUsesItemLevelAndCharm(t *testing.T) {
	mix := Item380Manager.MixConfig()
	for _, tc := range []struct{ level, rate int }{
		{6, 50}, {7, 60}, {9, 60}, {10, 70}, {13, 70}, {14, 80},
	} {
		if got := mix.SuccessRate(tc.level, 0); got != tc.rate {
			t.Errorf("SuccessRate(%d, 0) = %d, want %d", tc.level, got, tc.rate)
		}
		if got := mix.SuccessRate(tc.level, 10); got != tc.rate+10 {
			t.Errorf("SuccessRate(%d, 10) = %d, want %d", tc.level, got, tc.rate+10)
		}
	}
}
