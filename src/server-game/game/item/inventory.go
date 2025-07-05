package item

import (
	"errors"
	"fmt"
)

type Inventory struct {
	PositionedItems
}

func (inv *Inventory) CheckFlagsForItem(position int, it *Item) bool {
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
	if x+width > 8 ||
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
		ok := inv.CheckFlagsForItem(i, it)
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

func (inv *Inventory) AddItem(position int, it *Item) {
	inv.Items[position] = it
	inv.SetFlagsForItem(position, it)
}

func (inv *Inventory) RemoveItem(position int, it *Item) {
	inv.Items[position] = nil
	inv.ClearFlagsForItem(position, it)
}

func (inv *Inventory) UnmarshalJSON(buf []byte) error {
	inv.Size = 237
	inv.PositionedItems.CheckFlagsForItem = inv.CheckFlagsForItem
	inv.PositionedItems.SetFlagsForItem = inv.SetFlagsForItem
	return inv.PositionedItems.UnmarshalJSON(buf)
}

func (inv *Inventory) Scan(value any) error {
	buf, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to Scan Inventory value:", value))
	}
	return inv.UnmarshalJSON(buf)
}
