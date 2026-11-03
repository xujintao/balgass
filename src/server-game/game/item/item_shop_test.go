package item

import "testing"

func TestShopAddItem(t *testing.T) {
	newItem := func(width, height int) *Item {
		return &Item{ItemBase: &ItemBase{Width: width, Height: height}}
	}
	newShop := func() *Shop {
		return &Shop{PositionedItems: &PositionedItems{
			Size:  MaxShopItemCount,
			Items: make([]*Item, MaxShopItemCount),
			Flags: make([]bool, MaxShopItemCount),
		}}
	}

	shop := newShop()
	item := newItem(2, 2)
	if !shop.AddItem(item) {
		t.Fatal("AddItem rejected an item that fits at the first position")
	}
	if shop.Items[0] != item || !shop.Flags[0] || !shop.Flags[1] || !shop.Flags[8] || !shop.Flags[9] {
		t.Fatal("AddItem did not store the item and mark all occupied cells")
	}

	shop = newShop()
	shop.Flags[9] = true
	if !shop.AddItem(item) {
		t.Fatal("AddItem rejected an item that fits after an occupied position")
	}
	if shop.Items[2] != item || !shop.Flags[2] || !shop.Flags[3] || !shop.Flags[10] || !shop.Flags[11] {
		t.Fatal("AddItem did not use the first available position after overlap")
	}

	shop = newShop()
	if !shop.CanPlace(MaxShopItemCount-1, newItem(1, 1)) {
		t.Fatal("CanPlace rejected the last one-cell position")
	}
	if shop.CanPlace(MaxShopItemCount-1, newItem(2, 1)) {
		t.Fatal("CanPlace accepted an item extending past the grid")
	}
	if shop.CanPlace(0, nil) || shop.CanPlace(-1, item) || shop.CanPlace(MaxShopItemCount, item) ||
		shop.CanPlace(0, newItem(0, 1)) {
		t.Fatal("CanPlace accepted an invalid item or position")
	}

	fullShop := newShop()
	for i := range fullShop.Flags {
		fullShop.Flags[i] = true
	}
	fullItem := newItem(1, 1)
	if fullShop.AddItem(fullItem) {
		t.Fatal("AddItem succeeded when the shop had no free space")
	}
	if fullShop.Items[0] != nil {
		t.Fatal("AddItem changed inventory after failing to place an item")
	}
	for i, occupied := range fullShop.Flags {
		if !occupied {
			t.Fatalf("AddItem changed occupied flag %d after failing", i)
		}
	}

	shop = newShop()
	if shop.AddItem(nil) {
		t.Fatal("AddItem accepted a nil item")
	}
	for i, it := range shop.Items {
		if it != nil || shop.Flags[i] {
			t.Fatalf("AddItem changed shop at position %d for nil item", i)
		}
	}
}

func TestShopManagerQueriesInvalidShopAndPosition(t *testing.T) {
	shopItem := &Item{}
	shop := &Shop{PositionedItems: &PositionedItems{
		Size:  1,
		Items: []*Item{shopItem},
		Flags: []bool{true},
	}}
	manager := shopManager{shopTable: map[int]map[int]*Shop{1: {2: shop}}}

	if manager.GetShopInventory(99, 2) != nil || manager.GetShopInventory(1, 99) != nil {
		t.Fatal("GetShopInventory returned inventory for an unknown shop")
	}
	if manager.GetShopItem(99, 2, 0) != nil || manager.GetShopItem(1, 2, -1) != nil ||
		manager.GetShopItem(1, 2, 1) != nil {
		t.Fatal("GetShopItem returned an item for an invalid shop or position")
	}
	got := manager.GetShopItem(1, 2, 0)
	if got == nil || got == shopItem {
		t.Fatal("GetShopItem did not return a copy of the configured item")
	}
}
