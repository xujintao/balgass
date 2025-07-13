package item

import (
	"errors"
	"fmt"
)

type Warehouse struct {
	PositionedItems
}

func (w *Warehouse) CheckFlagsForItem(position int, it *Item) bool {
	maxHeight1 := 15
	maxHeight2 := 30
	x := position % 8
	y := position / 8
	width := it.Width
	height := it.Height
	if x+width > 8 ||
		y < maxHeight1 && y+height > maxHeight1 ||
		y >= maxHeight1 && y < maxHeight2 && y+height > maxHeight2 {
		return false
	}
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			if w.Flags[i+8*j] {
				return false
			}
		}
	}
	return true
}

func (w *Warehouse) FindFreePositionForItem(it *Item) int {
	for i, v := range w.Flags {
		if v {
			continue
		}
		ok := w.CheckFlagsForItem(i, it)
		if ok {
			return i
		}
	}
	return -1
}

func (w *Warehouse) setFlagsForItem(position int, it *Item, v bool) {
	x := position % 8
	y := position / 8
	width := it.Width
	height := it.Height
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			w.Flags[i+8*j] = v
		}
	}
}

func (w *Warehouse) SetFlagsForItem(position int, it *Item) {
	w.setFlagsForItem(position, it, true)
}

func (w *Warehouse) ClearFlagsForItem(position int, it *Item) {
	w.setFlagsForItem(position, it, false)
}

func (w *Warehouse) AddItem(position int, it *Item) {
	w.Items[position] = it
	w.SetFlagsForItem(position, it)
}

func (w *Warehouse) RemoveItem(position int, it *Item) {
	w.Items[position] = nil
	w.ClearFlagsForItem(position, it)
}

func (w *Warehouse) UnmarshalJSON(buf []byte) error {
	w.Size = 240
	return w.PositionedItems.UnmarshalJSON(
		buf,
		w.CheckFlagsForItem,
		w.SetFlagsForItem,
	)
}

func (w *Warehouse) Scan(value any) error {
	buf, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to Scan Inventory value:", value))
	}
	return w.UnmarshalJSON(buf)
}
