package maps

import (
	"slices"
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/item"
)

func TestItemNewDropLifecycle(t *testing.T) {
	var manager mapManager
	for _, number := range []int{0, 1} {
		manager.maps[number] = &_map{width: 8, height: 8, inventory: make([]*mapItem, 3)}
		if !manager.AddItem(number, 2, 3, &item.Item{}) {
			t.Fatal("AddItem failed")
		}
	}
	check := func(number int, want []bool) {
		t.Helper()
		var got []bool
		manager.MapEachItem(number, func(it *item.Item, index, x, y int, newDrop bool) {
			if index != len(got) || x != 2 || y != 3 || manager.PeekItem(number, index) != it {
				t.Fatalf("map %d: item identity or location changed at index %d", number, index)
			}
			got = append(got, newDrop)
		})
		if !slices.Equal(got, want) {
			t.Fatalf("map %d: new drop flags = %v, want %v", number, got, want)
		}
	}

	// Reading for one viewer must not consume the flag for another viewer.
	check(0, []bool{true})
	check(0, []bool{true})
	manager.ClearItemNewDrop()
	check(0, []bool{false})
	check(1, []bool{false}) // Also clear drops that nobody saw this round.

	if !manager.AddItem(0, 2, 3, &item.Item{}) {
		t.Fatal("adding a later drop failed")
	}
	check(0, []bool{false, true})
	manager.ClearItemNewDrop()
	check(0, []bool{false, false})

	// Picking up and dropping again creates a new ground instance in the reused slot.
	it := manager.PeekItem(0, 0)
	manager.RemoveItem(0, 0)
	if !manager.AddItem(0, 2, 3, it) {
		t.Fatal("dropping the picked item again failed")
	}
	check(0, []bool{true, false})
}
