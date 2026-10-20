package item

import (
	"encoding/json"
	"strings"
	"testing"
)

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
