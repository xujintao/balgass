package item

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestContainerItemPosition(t *testing.T) {
	it := &Item{}
	inv := Inventory{PositionedItems: PositionedItems{Size: 2, Items: []*Item{nil, it}}}
	warehouse := Warehouse{PositionedItems: PositionedItems{Size: 2, Items: []*Item{nil, it}}}
	if inv.Item(1) != it || inv.Item(-1) != nil || inv.Item(2) != nil {
		t.Fatal("inventory item lookup accepted an invalid position")
	}
	if warehouse.Item(1) != it || warehouse.Item(-1) != nil || warehouse.Item(2) != nil {
		t.Fatal("warehouse item lookup accepted an invalid position")
	}
}

func TestWingExcellentInventoryJSON(t *testing.T) {
	for _, tc := range []struct {
		name   string
		key    string
		oldKey string
		set    func(*Item)
		get    func(*Item) bool
	}{
		{"wing2", "wing2_excellent_speed", "excellent_wing2_speed", func(it *Item) { it.Wing2ExcellentSpeed = true }, func(it *Item) bool { return it.Wing2ExcellentSpeed }},
		{"wing3", "wing3_excellent_ignore", "excellent_wing3_ignore", func(it *Item) { it.Wing3ExcellentIgnore = true }, func(it *Item) bool { return it.Wing3ExcellentIgnore }},
		{"wing25", "wing25_excellent_hp", "excellent_wing25_hp", func(it *Item) { it.Wing25ExcellentHP = true }, func(it *Item) bool { return it.Wing25ExcellentHP }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			it := NewItem(12, 36)
			tc.set(it)
			inventory := Inventory{PositionedItems: PositionedItems{Size: 237, Items: make([]*Item, 237)}}
			inventory.Items[12] = it
			data, err := json.Marshal(inventory)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), `"`+tc.key+`":true`) || strings.Contains(string(data), tc.oldKey) {
				t.Fatalf("unexpected wing excellent JSON: %s", data)
			}
			var restored Inventory
			if err := json.Unmarshal(data, &restored); err != nil {
				t.Fatal(err)
			}
			if restored.Items[12] == nil || !tc.get(restored.Items[12]) {
				t.Fatalf("wing excellent option was not restored: %#v", restored.Items[12])
			}
		})
	}
}

func TestContainerRestoreRejectsInvalidItemsWithoutChangingState(t *testing.T) {
	for _, container := range []struct {
		name string
		size int
		new  func() (*PositionedItems, func([]byte) error, func(int, *Item) bool)
	}{
		{"inventory", 237, func() (*PositionedItems, func([]byte) error, func(int, *Item) bool) {
			inv := &Inventory{PositionedItems: PositionedItems{Size: 237, Items: make([]*Item, 237), Flags: make([]bool, 237)}}
			return &inv.PositionedItems, inv.UnmarshalJSON, inv.CanPlace
		}},
		{"warehouse", 240, func() (*PositionedItems, func([]byte) error, func(int, *Item) bool) {
			w := &Warehouse{PositionedItems: PositionedItems{Size: 240, Items: make([]*Item, 240), Flags: make([]bool, 240)}}
			return &w.PositionedItems, w.UnmarshalJSON, w.CanPlace
		}},
	} {
		t.Run(container.name, func(t *testing.T) {
			for _, tc := range []struct {
				name string
				data string
			}{
				{"negative position", `[ {"position":12,"section":4,"index":7}, {"position":-1,"section":4,"index":7} ]`},
				{"position at size", fmt.Sprintf(`[ {"position":12,"section":4,"index":7}, {"position":%d,"section":4,"index":7} ]`, container.size)},
				{"footprint beyond end", fmt.Sprintf(`[ {"position":12,"section":4,"index":7}, {"position":%d,"section":12,"index":36} ]`, container.size-1)},
				{"duplicate position", `[ {"position":12,"section":4,"index":7}, {"position":12,"section":4,"index":7} ]`},
				{"overlapping footprint", `[ {"position":12,"section":12,"index":36}, {"position":13,"section":4,"index":7} ]`},
				{"null item", `[ {"position":12,"section":4,"index":7}, null ]`},
				{"invalid section", `[ {"position":12,"section":4,"index":7}, {"position":13,"section":-1,"index":7} ]`},
				{"unknown item", `[ {"position":12,"section":4,"index":7}, {"position":13,"section":4,"index":511} ]`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					pi, restore, _ := container.new()
					original := NewItem(4, 7)
					pi.Items[0], pi.Flags[0] = original, true
					oldSize := pi.Size
					if err := restore([]byte(tc.data)); err == nil {
						t.Fatal("invalid container data was accepted")
					} else if tc.name != "null item" && !strings.Contains(err.Error(), "position") {
						t.Fatalf("error lacks position: %v", err)
					}
					if pi.Size != oldSize || pi.Items[0] != original || !pi.Flags[0] || pi.Items[12] != nil || pi.Flags[12] {
						t.Fatal("failed restore changed the existing container")
					}
				})
			}

			pi, restore, check := container.new()
			bolt := NewItem(4, 7)
			if check(-1, bolt) || check(container.size, bolt) || check(container.size-1, NewItem(12, 36)) {
				t.Fatal("invalid placement passed the occupancy check")
			}
			last := container.size - 1
			if err := restore([]byte(fmt.Sprintf(`[{"position":%d,"section":4,"index":7}]`, last))); err != nil {
				t.Fatalf("valid last position was rejected: %v", err)
			}
			if pi.Items[last] == nil || !pi.Flags[last] {
				t.Fatal("valid last position was not restored")
			}
		})
	}
}

func TestInventoryAddAndMove(t *testing.T) {
	inv := Inventory{PositionedItems: PositionedItems{Size: 237, Items: make([]*Item, 237), Flags: make([]bool, 237)}}
	large := &Item{ItemBase: &ItemBase{Width: 2, Height: 2}}
	blocker := &Item{ItemBase: &ItemBase{Width: 1, Height: 1}}
	if inv.AddItem(-1, large) || inv.AddItem(19, large) || inv.AddItem(237, large) {
		t.Fatal("invalid inventory position was accepted")
	}
	if !inv.AddItem(12, large) || inv.AddItem(13, blocker) || !inv.AddItem(14, blocker) {
		t.Fatal("inventory placement did not respect occupied cells")
	}
	if inv.MoveItem(12, 13) || inv.Item(12) != large || inv.Item(14) != blocker || !inv.Flags[12] || !inv.Flags[20] {
		t.Fatal("blocked inventory move did not restore the source")
	}
	inv.RemoveItem(14, blocker)
	if !inv.MoveItem(12, 13) || inv.Item(12) != nil || inv.Item(13) != large ||
		inv.Flags[12] || inv.Flags[20] || !inv.Flags[13] || !inv.Flags[14] || !inv.Flags[21] || !inv.Flags[22] {
		t.Fatal("overlapping inventory move left incorrect items or flags")
	}
}

func TestWarehouseAddAndMove(t *testing.T) {
	w := Warehouse{PositionedItems: PositionedItems{Size: 240, Items: make([]*Item, 240), Flags: make([]bool, 240)}}
	large := &Item{ItemBase: &ItemBase{Width: 2, Height: 2}}
	blocker := &Item{ItemBase: &ItemBase{Width: 1, Height: 1}}
	if w.AddItem(-1, large) || w.AddItem(7, large) || w.AddItem(240, large) {
		t.Fatal("invalid warehouse position was accepted")
	}
	if !w.AddItem(0, large) || w.AddItem(1, blocker) || !w.AddItem(2, blocker) {
		t.Fatal("warehouse placement did not respect occupied cells")
	}
	if w.MoveItem(0, 1) || w.Item(0) != large || w.Item(2) != blocker || !w.Flags[0] || !w.Flags[8] {
		t.Fatal("blocked warehouse move did not restore the source")
	}
	w.RemoveItem(2, blocker)
	if !w.MoveItem(0, 1) || w.Item(0) != nil || w.Item(1) != large ||
		w.Flags[0] || w.Flags[8] || !w.Flags[1] || !w.Flags[2] || !w.Flags[9] || !w.Flags[10] {
		t.Fatal("overlapping warehouse move left incorrect items or flags")
	}
}

func TestExcellentContributionsDescribeBonusesAndDivisors(t *testing.T) {
	var plain Item
	var empty ItemExcellentEffect
	ExcellentManager.applyExcellentEffect(&plain, &empty)
	if empty != (ItemExcellentEffect{}) {
		t.Fatalf("plain item contributed an effect: %+v", empty)
	}
	weapon := Item{
		ItemBase: &ItemBase{KindA: KindAWeapon},
		ExcellentAttackRate: true,
		ExcellentAttackLevel: true, ExcellentAttackPercent: true,
		ExcellentAttackSpeed: true, ExcellentAttackHP: true,
		ExcellentDefenseHP: true, ExcellentDefenseMP: true, ExcellentDefenseRate: true,
	}
	var got ItemExcellentEffect
	ExcellentManager.applyExcellentEffect(&weapon, &got)
	want := ItemExcellentEffect{
		ExcellentAttackRate: 10,
		AttackLevel:   20,
		AttackPercent: 2,
		AttackSpeed:   7, MagicSpeed: 7, MonsterDieGetHP: 0.125,
		MaxHPPercent: 4, MaxMPPercent: 4, DefenseRatePercent: 10,
	}
	if got != want {
		t.Fatalf("excellent contributions = %+v, want %+v", got, want)
	}
	wing := Item{Level: 9, Wing2ExcellentHP: true, Wing2ExcellentMP: true, Wing2ExcellentAG: true, Wing2ExcellentSpeed: true, Wing3ExcellentIgnore: true}
	var wingEffect ItemExcellentEffect
	ExcellentManager.applyExcellentEffect(&wing, &wingEffect)
	if wingEffect.MaxHP != 95 || wingEffect.MaxMP != 95 || wingEffect.MaxAG != 50 || wingEffect.AttackSpeed != 5 || wingEffect.MagicSpeed != 5 || wingEffect.IgnoreDefenseRate != 5 {
		t.Fatalf("wing excellent contributions = %+v", wingEffect)
	}
}

func TestExcellentManagerCallbackIncludesWearingPosition(t *testing.T) {
	inventory := &Inventory{PositionedItems: PositionedItems{Items: make([]*Item, INVENTORY_WEAR_SIZE)}}
	inventory.Items[0] = &Item{ItemBase: &ItemBase{KindA: KindAWeapon}, ExcellentAttackLevel: true}
	inventory.Items[1] = &Item{ItemBase: &ItemBase{KindA: KindAWeapon}, ExcellentAttackLevel: true}

	var positions []int
	ExcellentManager.ApplyExcellentItemEffect(inventory, func(effect ItemExcellentEffect, position int) {
		positions = append(positions, position)
		if effect.AttackLevel != 20 {
			t.Errorf("position %d effect = %+v", position, effect)
		}
	})

	if len(positions) != 2 || positions[0] != 0 || positions[1] != 1 {
		t.Fatalf("callback positions = %v, want [0 1]", positions)
	}
}
