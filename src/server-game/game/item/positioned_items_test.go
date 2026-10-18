package item

import (
	"fmt"
	"strings"
	"testing"
)

func TestContainerRestoreRejectsInvalidItemsWithoutChangingState(t *testing.T) {
	for _, container := range []struct {
		name string
		size int
		new  func() (*PositionedItems, func([]byte) error, func(int, *Item) bool)
	}{
		{"inventory", 237, func() (*PositionedItems, func([]byte) error, func(int, *Item) bool) {
			inv := &Inventory{PositionedItems: PositionedItems{Size: 237, Items: make([]*Item, 237), Flags: make([]bool, 237)}}
			return &inv.PositionedItems, inv.UnmarshalJSON, inv.CheckFlagsForItem
		}},
		{"warehouse", 240, func() (*PositionedItems, func([]byte) error, func(int, *Item) bool) {
			w := &Warehouse{PositionedItems: PositionedItems{Size: 240, Items: make([]*Item, 240), Flags: make([]bool, 240)}}
			return &w.PositionedItems, w.UnmarshalJSON, w.CheckFlagsForItem
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
