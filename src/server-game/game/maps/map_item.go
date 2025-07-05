package maps

import (
	"time"

	"github.com/xujintao/balgass/src/server-game/game/item"
)

type mapItem struct {
	*item.Item
	x           int
	y           int
	expiredTime time.Time
}

func (m *mapManager) AddItem(number, x, y int, item *item.Item) bool {
	return m.maps[number].addItem(x, y, item)
}

func (m *mapManager) PeekItem(number, index int) *item.Item {
	return m.maps[number].peekItem(index)
}

func (m *mapManager) RemoveItem(number, index int) {
	m.maps[number].removeItem(index)
}

func (m *mapManager) ExpireItem(now time.Time) {
	for _, v := range m.maps {
		v.expireItem(now)
	}
}

func (m *mapManager) MapEachItem(number int, f func(item *item.Item, index, x, y int)) {
	m.maps[number].eachItem(f)
}

func (m *mapManager) MapItem(number, index int, f func(item *item.Item, index, x, y int)) {
	m.maps[number].item(index, f)
}

func (m *_map) addItem(x, y int, item *item.Item) bool {
	if !m.valid(x, y) {
		return false
	}

	for i := range m.inventory {
		if m.inventory[i] == nil {
			mapItem := mapItem{
				Item:        item,
				x:           x,
				y:           y,
				expiredTime: time.Now().Add(time.Minute),
			}
			m.inventory[i] = &mapItem
			return true
		}
	}
	return false
}

func (m *_map) peekItem(index int) *item.Item {
	if index < 0 || index >= len(m.inventory) {
		return nil
	}
	mapItem := m.inventory[index]
	if mapItem == nil {
		return nil
	}
	return mapItem.Item
}

func (m *_map) removeItem(index int) {
	item := m.peekItem(index)
	if item == nil {
		return
	}
	m.inventory[index] = nil
}

func (m *_map) expireItem(now time.Time) {
	for i, mapItem := range m.inventory {
		if mapItem == nil {
			continue
		}
		if now.After(mapItem.expiredTime) {
			m.inventory[i] = nil
		}
	}
}

func (m *_map) eachItem(f func(item *item.Item, index, x, y int)) {
	for i, mapItem := range m.inventory {
		if mapItem == nil {
			continue
		}
		f(mapItem.Item, i, mapItem.x, mapItem.y)
	}
}

func (m *_map) item(index int, f func(item *item.Item, index, x, y int)) {
	mapItem := m.inventory[index]
	if mapItem == nil {
		f(nil, 0, 0, 0)
		return
	}
	f(mapItem.Item, index, mapItem.x, mapItem.y)
}
