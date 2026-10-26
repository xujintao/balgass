package item

const chaosBoxWidth = 8
const chaosBoxHeight = 4

type ChaosBox struct {
	items [chaosBoxWidth * chaosBoxHeight]*Item
}

func (box *ChaosBox) Item(position int) *Item {
	if position < 0 || position >= len(box.items) {
		return nil
	}
	return box.items[position]
}

func (box *ChaosBox) Items() []*Item {
	items := make([]*Item, len(box.items))
	copy(items, box.items[:])
	return items
}

func (box *ChaosBox) CanPlace(position int, it *Item, except int) bool {
	if position < 0 || position >= len(box.items) || it == nil || it.Width <= 0 || it.Height <= 0 {
		return false
	}
	x, y := position%chaosBoxWidth, position/chaosBoxWidth
	if x+it.Width > chaosBoxWidth || y+it.Height > chaosBoxHeight {
		return false
	}
	for at, other := range box.items {
		if at == except || other == nil {
			continue
		}
		ox, oy := at%chaosBoxWidth, at/chaosBoxWidth
		if x < ox+other.Width && ox < x+it.Width && y < oy+other.Height && oy < y+it.Height {
			return false
		}
	}
	return true
}

func (box *ChaosBox) AddItem(position int, it *Item) bool {
	if !box.CanPlace(position, it, -1) {
		return false
	}
	box.items[position] = it
	return true
}

func (box *ChaosBox) RemoveItem(position int) *Item {
	it := box.Item(position)
	if it == nil {
		return nil
	}
	box.items[position] = nil
	return it
}

func (box *ChaosBox) MoveItem(src, dst int) bool {
	it := box.Item(src)
	if src == dst || !box.CanPlace(dst, it, src) {
		return false
	}
	box.items[dst], box.items[src] = it, nil
	return true
}

func (box *ChaosBox) Clear() {
	*box = ChaosBox{}
}

func (box *ChaosBox) ReplaceItem(position int, it *Item) bool {
	if box.Item(position) == nil || it == nil {
		return false
	}
	box.items[position] = it
	return true
}
