package item

import "testing"

func TestChaosBoxItemsAndPositions(t *testing.T) {
	var box ChaosBox
	first := &Item{ItemBase: &ItemBase{Width: 2, Height: 2}}
	second := &Item{ItemBase: &ItemBase{Width: 2, Height: 2}}
	if box.AddItem(7, first) || box.AddItem(24, first) {
		t.Fatal("item crossed a Chaos Box edge")
	}
	if !box.AddItem(0, first) || box.AddItem(1, second) || !box.AddItem(2, second) {
		t.Fatal("Chaos Box placement did not respect occupied cells")
	}
	if !box.MoveItem(0, 8) || box.Item(0) != nil || box.Item(8) != first {
		t.Fatal("moving an item did not preserve it")
	}
	replacement := &Item{ItemBase: &ItemBase{Width: 1, Height: 1}}
	if !box.ReplaceItem(8, replacement) || box.Item(8) != replacement {
		t.Fatal("replacing an item failed")
	}
	items := box.Items()
	items[8] = nil
	if box.Item(8) != replacement {
		t.Fatal("item snapshot changed the Chaos Box")
	}
	if removed := box.RemoveItem(8); removed != replacement || box.Item(8) != nil {
		t.Fatal("removing an item failed")
	}
	box.Clear()
	if box.Item(2) != nil {
		t.Fatal("clearing the Chaos Box left an item")
	}
}
