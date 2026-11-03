package item

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/lang"
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

func TestHarmonyItemJSONAndProtocolEncoding(t *testing.T) {
	it := NewItem(0, 0)
	it.HarmonyOption = 3
	it.HarmonyLevel = 6
	it.Level = 10
	data, err := json.Marshal(it)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"harmony_option":3`) || strings.Contains(string(data), `"harmony_effect":`) {
		t.Fatalf("Harmony option JSON key = %s, want harmony_option", data)
	}
	var restored Item
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.HarmonyOption != 3 || restored.HarmonyLevel != 6 || restored.Level != 10 {
		t.Fatalf("restored Harmony option = %d level %d at item level %d, want 3 level 6 at item level 10", restored.HarmonyOption, restored.HarmonyLevel, restored.Level)
	}
	if !strings.Contains(string(data), `"harmony_level":6`) {
		t.Fatalf("Harmony option level missing from JSON: %s", data)
	}
	var oldJSON Item
	if err := json.Unmarshal([]byte(`{"harmony_effect":3}`), &oldJSON); err != nil {
		t.Fatal(err)
	}
	if oldJSON.HarmonyOption != 0 {
		t.Fatalf("old harmony_effect JSON key unexpectedly restored Harmony option %d", oldJSON.HarmonyOption)
	}
	frame, err := it.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if frame[6] != 0x36 {
		t.Fatalf("Harmony protocol byte = %#02x, want 0x36", frame[6])
	}

	pet := NewItem(13, 2)
	pet.PetActive = 1
	petJSON, err := json.Marshal(pet)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(petJSON), `"pet_active":1`) {
		t.Fatalf("pet JSON does not contain pet_active: %s", petJSON)
	}
	var restoredPet Item
	if err := json.Unmarshal(petJSON, &restoredPet); err != nil {
		t.Fatal(err)
	}
	if restoredPet.PetActive != 1 {
		t.Fatalf("restored pet active state = %d, want 1", restoredPet.PetActive)
	}
	frame, err = pet.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if frame[6] != 1 {
		t.Fatalf("active pet protocol byte = %#02x, want 0x01", frame[6])
	}
	pet.Level = 6
	frame, err = pet.Marshal()
	if err != nil {
		t.Fatalf("Marshal active pet with a nonzero item level: %v", err)
	}
	if frame[6] != 1 {
		t.Fatalf("pet protocol byte with item level = %#02x, want 0x01", frame[6])
	}

	for _, outOfRange := range []Item{
		{ItemBase: it.ItemBase, Section: 0, Index: 0, HarmonyOption: 16, Level: 0},
		{ItemBase: it.ItemBase, Section: 0, Index: 0, HarmonyOption: 1, HarmonyLevel: 16},
	} {
		frame, err := outOfRange.Marshal()
		if err != nil {
			t.Fatalf("out-of-range Harmony value was rejected: %v", err)
		}
		want := byte(outOfRange.HarmonyOption&0x0F)<<4 | byte(outOfRange.HarmonyLevel&0x0F)
		if frame[6] != want {
			t.Fatalf("Harmony protocol byte = %#02x, want %#02x", frame[6], want)
		}
	}
}

func TestHarmonyEnhancementAndSmeltingUseSeparateOptionLevel(t *testing.T) {
	cfg := HarmonyManager.cfg
	options := append([]harmonyOption(nil), HarmonyManager.harmonyOptions[harmonyItemWeapon]...)
	defer func() {
		HarmonyManager.cfg = cfg
		HarmonyManager.harmonyOptions[harmonyItemWeapon] = options
	}()
	HarmonyManager.cfg.HarmonySystem.EnableStrengthenItem = true
	HarmonyManager.cfg.HarmonyMix.StrengthenItemSuccessRate = 100
	HarmonyManager.cfg.HarmonySystem.EnableSmeltItem = true
	HarmonyManager.cfg.HarmonyMix.SmeltItemLowerSuccessRate = 100
	for i := range HarmonyManager.harmonyOptions[harmonyItemWeapon] {
		HarmonyManager.harmonyOptions[harmonyItemWeapon][i].weight = 0
	}
	HarmonyManager.harmonyOptions[harmonyItemWeapon][harmonyOptionWeaponIncCriticalDamage].weight = 1

	it := NewItem(0, 0)
	it.Level = 10
	attempted, result := HarmonyManager.StrengthenItem(it)
	if !attempted || result != lang.MsgStrengthenSuccess || it.HarmonyOption == 0 {
		t.Fatalf("Harmony strengthening = (%t, %v), option %d", attempted, result, it.HarmonyOption)
	}
	if it.Level != 10 || it.HarmonyLevel != HarmonyManager.harmonyOptions[harmonyItemWeapon][it.HarmonyOption].reqLevel {
		t.Fatalf("strengthening levels = item %d, Harmony %d, want item 10 and configured requirement", it.Level, it.HarmonyLevel)
	}
	initialHarmonyLevel := it.HarmonyLevel
	attempted, success := HarmonyManager.SmeltItem(it, false)
	if !attempted || !success || it.Level != 10 || it.HarmonyLevel != initialHarmonyLevel+1 {
		t.Fatalf("smelting result = (%t, %t), item/Harmony levels %d/%d", attempted, success, it.Level, it.HarmonyLevel)
	}
	HarmonyManager.cfg.HarmonyMix.SmeltItemLowerSuccessRate = 0
	attempted, success = HarmonyManager.SmeltItem(it, false)
	if !attempted || success || it.HarmonyLevel != HarmonyManager.harmonyOptions[harmonyItemWeapon][it.HarmonyOption].reqLevel {
		t.Fatalf("failed smelting result = (%t, %t), Harmony level %d", attempted, success, it.HarmonyLevel)
	}
	if fee, ok := HarmonyManager.RestoreItemMoney(it); !ok || fee != HarmonyManager.harmonyOptions[harmonyItemWeapon][it.HarmonyOption].effects[it.HarmonyLevel].reqMoney {
		t.Fatalf("restore fee = %d, valid %t", fee, ok)
	}
}

func TestHarmonyStrengtheningFailureDoesNotSetOption(t *testing.T) {
	cfg := HarmonyManager.cfg
	options := append([]harmonyOption(nil), HarmonyManager.harmonyOptions[harmonyItemWeapon]...)
	defer func() {
		HarmonyManager.cfg = cfg
		HarmonyManager.harmonyOptions[harmonyItemWeapon] = options
	}()
	HarmonyManager.cfg.HarmonySystem.EnableStrengthenItem = true
	HarmonyManager.cfg.HarmonyMix.StrengthenItemSuccessRate = 0
	for i := range HarmonyManager.harmonyOptions[harmonyItemWeapon] {
		HarmonyManager.harmonyOptions[harmonyItemWeapon][i].weight = 0
	}
	HarmonyManager.harmonyOptions[harmonyItemWeapon][harmonyOptionWeaponIncCriticalDamage].weight = 1

	it := NewItem(0, 0)
	it.Level = 10
	attempted, err := HarmonyManager.StrengthenItem(it)
	if !attempted || err != lang.MsgStrengthenFailed || it.HarmonyOption != 0 || it.HarmonyLevel != 0 {
		t.Fatalf("failed strengthening = (%t, %v), option %d/%d", attempted, err, it.HarmonyOption, it.HarmonyLevel)
	}
}

func TestCanMakeRefiningStoneUsesConfiguredLevelList(t *testing.T) {
	it := NewItem(0, 0)
	it.Level = 4
	if valid, _ := HarmonyManager.CanMakeRefiningStone(it); !valid {
		t.Fatal("configured level 4 item should be eligible to make a refining stone")
	}
	it.Level = 3
	if valid, _ := HarmonyManager.CanMakeRefiningStone(it); valid {
		t.Fatal("item below configured refining stone level was accepted")
	}
}

func TestCanMakeRefiningStoneUsesEligibleMaterialCategories(t *testing.T) {
	for _, tc := range []struct {
		section, index int
		want           bool
	}{
		{section: 0, index: 0, want: true},  // weapon
		{section: 5, index: 0, want: true},  // staff
		{section: 6, index: 0, want: true},  // armor
		{section: 4, index: 7, want: false}, // excluded crossbow bolt
		{section: 4, index: 15, want: false},
		{section: 13, index: 0, want: false}, // non-equipment
	} {
		it := NewItem(tc.section, tc.index)
		it.Level = 15
		if got, _ := HarmonyManager.CanMakeRefiningStone(it); got != tc.want {
			t.Errorf("CanMakeRefiningStone(%d/%d) = %t, want %t", tc.section, tc.index, got, tc.want)
		}
	}
}

func TestCanMakeRefiningStoneTreatsSetItemsAsHigh(t *testing.T) {
	old := conf.Common.General.EnableUseSetHarmonyItem
	conf.Common.General.EnableUseSetHarmonyItem = true
	t.Cleanup(func() { conf.Common.General.EnableUseSetHarmonyItem = old })

	for _, tc := range []struct {
		name      string
		set       int
		excellent bool
		wantHigh  bool
	}{
		{name: "set", set: 1, wantHigh: true},
		{name: "excellent set", set: 1, excellent: true, wantHigh: true},
		{name: "normal item", wantHigh: false},
	} {
		it := NewItem(6, 0)
		it.Level, it.Set, it.ExcellentAttackHP = 15, tc.set, tc.excellent
		ok, high := HarmonyManager.CanMakeRefiningStone(it)
		if !ok || high != tc.wantHigh {
			t.Errorf("%s material eligibility = (%t, %t), want (true, %t)", tc.name, ok, high, tc.wantHigh)
		}
	}
}

func TestHarmonyItemEffectsUseConfiguredOptionAndLevel(t *testing.T) {
	weapon := NewItem(0, 0)
	weapon.HarmonyOption = int(harmonyOptionWeaponIncDamageMin)
	var option ItemHarmonyOption
	HarmonyManager.applyItemHarmonyOption(weapon, &option)
	if option.IncDamageMin != 2 {
		t.Fatalf("minimum damage contribution = %d, want 2", option.IncDamageMin)
	}

	weapon.HarmonyOption = int(harmonyOptionWeaponIncCriticalDamage)
	weapon.Level = 5
	weapon.HarmonyLevel = 6
	option = ItemHarmonyOption{}
	HarmonyManager.applyItemHarmonyOption(weapon, &option)
	if option != (ItemHarmonyOption{}) {
		t.Fatalf("inactive Harmony option contributed: %+v", option)
	}
	weapon.Level = 6
	HarmonyManager.applyItemHarmonyOption(weapon, &option)
	if option.IncCriticalDamage != 12 {
		t.Fatalf("critical damage contribution = %d, want 12", option.IncCriticalDamage)
	}

	inv := &Inventory{PositionedItems: PositionedItems{Items: make([]*Item, 237)}}
	inv.Items[0] = weapon
	rightWeapon := NewItem(0, 0)
	rightWeapon.HarmonyOption = int(harmonyOptionWeaponIncDamageMax)
	inv.Items[1] = rightWeapon
	inv.Items[12] = NewItem(0, 0)
	var positions []int
	HarmonyManager.ApplyHarmonyItemOption(inv, func(got ItemHarmonyOption, position int) {
		positions = append(positions, position)
		switch position {
		case 0:
			if got.IncCriticalDamage != 12 {
				t.Errorf("left weapon Harmony contribution = %+v", got)
			}
		case 1:
			if got.IncDamageMax != 3 {
				t.Errorf("right weapon Harmony contribution = %+v", got)
			}
		}
	})
	if len(positions) != 2 || positions[0] != 0 || positions[1] != 1 {
		t.Fatalf("Harmony callback positions = %v, want [0 1]", positions)
	}
}

func TestHarmonyItemCalcReducesRequirements(t *testing.T) {
	it := NewItem(0, 0)
	strength, agility := it.ReqStrength, it.ReqDexterity
	it.HarmonyOption = int(harmonyOptionWeaponDecReqStrength)
	if it.RequiredStrength() != max(0, strength-6) || it.RequiredDexterity() != agility {
		t.Fatalf("requirements after Harmony = %d/%d, want %d/%d", it.RequiredStrength(), it.RequiredDexterity(), max(0, strength-6), agility)
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
	var empty ItemExcellentOption
	ExcellentManager.applyItemExcellentOption(&plain, &empty)
	if empty != (ItemExcellentOption{}) {
		t.Fatalf("plain item contributed an effect: %+v", empty)
	}
	weapon := Item{
		ItemBase:             &ItemBase{KindA: KindAWeapon},
		ExcellentAttackRate:  true,
		ExcellentAttackLevel: true, ExcellentAttackPercent: true,
		ExcellentAttackSpeed: true, ExcellentAttackHP: true,
		ExcellentDefenseHP: true, ExcellentDefenseMP: true, ExcellentDefenseRate: true,
	}
	var got ItemExcellentOption
	ExcellentManager.applyItemExcellentOption(&weapon, &got)
	want := ItemExcellentOption{
		ExcellentAttackRate: 10,
		AttackLevel:         20,
		AttackPercent:       2,
		AttackSpeed:         7, MagicSpeed: 7, MonsterDieGetHP: 0.125,
		MaxHPPercent: 4, MaxMPPercent: 4, DefenseRatePercent: 10,
	}
	if got != want {
		t.Fatalf("excellent contributions = %+v, want %+v", got, want)
	}
	wing := Item{Level: 9, Wing2ExcellentHP: true, Wing2ExcellentMP: true, Wing2ExcellentAG: true, Wing2ExcellentSpeed: true, Wing3ExcellentIgnore: true}
	var wingEffect ItemExcellentOption
	ExcellentManager.applyItemExcellentOption(&wing, &wingEffect)
	if wingEffect.MaxHP != 95 || wingEffect.MaxMP != 95 || wingEffect.MaxAG != 50 || wingEffect.AttackSpeed != 5 || wingEffect.MagicSpeed != 5 || wingEffect.IgnoreDefenseRate != 5 {
		t.Fatalf("wing excellent contributions = %+v", wingEffect)
	}
}

func TestExcellentManagerCallbackIncludesWearingPosition(t *testing.T) {
	inventory := &Inventory{PositionedItems: PositionedItems{Items: make([]*Item, INVENTORY_WEAR_SIZE)}}
	inventory.Items[0] = &Item{ItemBase: &ItemBase{KindA: KindAWeapon}, ExcellentAttackLevel: true}
	inventory.Items[1] = &Item{ItemBase: &ItemBase{KindA: KindAWeapon}, ExcellentAttackLevel: true}

	var positions []int
	ExcellentManager.ApplyExcellentItemOption(inventory, func(effect ItemExcellentOption, position int) {
		positions = append(positions, position)
		if effect.AttackLevel != 20 {
			t.Errorf("position %d effect = %+v", position, effect)
		}
	})

	if len(positions) != 2 || positions[0] != 0 || positions[1] != 1 {
		t.Fatalf("callback positions = %v, want [0 1]", positions)
	}
}
