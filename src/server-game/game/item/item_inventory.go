package item

import (
	"errors"
	"fmt"
)

const INVENTORY_WEAR_SIZE = 12

type Inventory struct {
	PositionedItems
}

func (inv *Inventory) WearingItems() []*Item {
	return inv.Items[:INVENTORY_WEAR_SIZE]
}

func (inv *Inventory) Item(position int) *Item {
	if position < 0 || position >= inv.Size || position >= len(inv.Items) {
		return nil
	}
	return inv.Items[position]
}

func (inv *Inventory) CanPlace(position int, it *Item) bool {
	if it == nil || position < 0 || position >= inv.Size {
		return false
	}
	if position < 12 {
		return !inv.Flags[position]
	}
	maxHeight1 := 8
	maxHeight2 := 12
	maxHeight3 := 16
	position -= 12
	flags := inv.Flags[12:]
	x := position % 8
	y := position / 8
	width := it.Width
	height := it.Height
	if width <= 0 || width > 8-x || height <= 0 ||
		height > (len(flags)-1-(x+width-1))/8-y+1 ||
		y < maxHeight1 && y+height > maxHeight1 ||
		y >= maxHeight1 && y < maxHeight2 && y+height > maxHeight2 ||
		y >= maxHeight2 && y < maxHeight3 && y+height > maxHeight3 {
		return false
	}
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			if flags[i+8*j] {
				return false
			}
		}
	}
	return true
}

func (inv *Inventory) FindFreePositionForItem(it *Item) int {
	for i := range inv.Flags {
		if i < 12 {
			continue
		}
		ok := inv.CanPlace(i, it)
		if ok {
			return i
		}
		it2 := inv.Items[i]
		if it2 == nil {
			continue
		}
		if it2.Overlap != 0 &&
			it2.Code == it.Code &&
			it2.Level == it.Level &&
			it2.Durability <= it.Overlap-it.Durability {
			return i
		}
	}
	return -1
}

func (inv *Inventory) setFlagsForItem(position int, it *Item, v bool) {
	if position < 12 {
		inv.Flags[position] = v
		return
	}
	position -= 12
	flags := inv.Flags[12:]
	x := position % 8
	y := position / 8
	width := it.Width
	height := it.Height
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			flags[i+8*j] = v
		}
	}
}

func (inv *Inventory) SetFlagsForItem(position int, it *Item) {
	inv.setFlagsForItem(position, it, true)
}

func (inv *Inventory) ClearFlagsForItem(position int, it *Item) {
	inv.setFlagsForItem(position, it, false)
}

func (inv *Inventory) AddItem(position int, it *Item) bool {
	if !inv.CanPlace(position, it) {
		return false
	}
	inv.Items[position] = it
	inv.SetFlagsForItem(position, it)
	return true
}

func (inv *Inventory) RemoveItem(position int, it *Item) {
	inv.Items[position] = nil
	inv.ClearFlagsForItem(position, it)
}

func (inv *Inventory) MoveItem(src, dst int) bool {
	it := inv.Item(src)
	if it == nil || src == dst {
		return false
	}
	inv.RemoveItem(src, it)
	if !inv.AddItem(dst, it) {
		inv.AddItem(src, it)
		return false
	}
	return true
}

func (inv *Inventory) UnmarshalJSON(buf []byte) error {
	oldSize := inv.Size
	inv.Size = 237
	err := inv.PositionedItems.UnmarshalJSON(
		buf,
		inv.CanPlace,
		inv.SetFlagsForItem,
	)
	if err != nil {
		inv.Size = oldSize
	}
	return err
}

func (inv *Inventory) Scan(value any) error {
	buf, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to Scan Inventory value:", value))
	}
	return inv.UnmarshalJSON(buf)
}
