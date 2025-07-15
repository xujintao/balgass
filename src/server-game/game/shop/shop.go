package shop

import (
	"encoding/xml"
	"log/slog"
	"path"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/item"
)

const MaxShopItemCount = 120 // 8*15

func init() {
	ShopManager.init()
}

type Shop struct {
	*item.PositionedItems
	NPCIndex     int    `xml:"NPCIndex,attr"`
	MapNumber    int    `xml:"MapNumber,attr"`
	PosX         int    `xml:"PosX,attr"`
	PosY         int    `xml:"PosY,attr"`
	Dir          int    `xml:"Dir,attr"`
	VipType      int    `xml:"VipType,attr"`
	GMShop       bool   `xml:"GMShop,attr"`
	FileName     string `xml:"FileName,attr"`
	MossMerchant bool   `xml:"MossMerchant,attr"`
	VIPType      int    `xml:"VIPType,attr"`
	BattleCore   bool   `xml:"BattleCore,attr"`
}

func (s *Shop) Scan(file string) {
	// Shop was generated 2023-11-17 14:35:07 by https://xml-to-go.github.io/ in Ukraine.
	type ShopInventory struct {
		XMLName xml.Name `xml:"Shop"`
		Text    string   `xml:",chardata"`
		Item    []*struct {
			Text        string `xml:",chardata"`
			Cat         int    `xml:"Cat,attr"`
			Index       int    `xml:"Index,attr"`
			Level       int    `xml:"Level,attr"`
			Durability  int    `xml:"Durability,attr"`
			Skill       bool   `xml:"Skill,attr"`
			Luck        bool   `xml:"Luck,attr"`
			Option      int    `xml:"Option,attr"`
			Exc         int    `xml:"Exc,attr"`
			SetItem     int    `xml:"SetItem,attr"`
			SocketCount int    `xml:"SocketCount,attr"`
			Elemental   int    `xml:"Elemental,attr"`
			Serial      int    `xml:"Serial,attr"`
		} `xml:"Item"`
	}
	var shopInventory ShopInventory
	conf.XML(conf.PathCommon, file, &shopInventory)
	for _, sitem := range shopInventory.Item {
		it := item.NewItem(sitem.Cat, sitem.Index)
		it.Level = sitem.Level
		it.Durability = sitem.Durability
		it.Skill = sitem.Skill
		it.Lucky = sitem.Luck
		it.Addition = sitem.Option << 2
		it.Calc()
		if it.Durability == 0 {
			it.Durability = it.MaxDurability
		}
		i := s.FindFreePositionForItem(it)
		if i == -1 {
			slog.Error("FindFreePositionForItem",
				"shop", s.FileName, "item", it.Name)
			continue
		}
		s.SetFlagsForItem(i, it)
		s.Items[i] = it
	}
}

func (s *Shop) FindFreePositionForItem(it *item.Item) int {
	maxHeight := len(s.Flags) / 8
outer:
	for i, v := range s.Flags {
		if v {
			continue
		}
		x := i % 8
		y := i / 8
		width := it.Width
		height := it.Height
		if x+width > 8 ||
			y+height > maxHeight {
			continue
		}
		for i := x; i < x+width; i++ {
			for j := y; j < y+height; j++ {
				if s.Flags[i+8*j] {
					continue outer
				}
			}
		}
		return i
	}
	return -1
}

func (s *Shop) SetFlagsForItem(position int, it *item.Item) {
	x := position % 8
	y := position / 8
	width := it.Width
	height := it.Height
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			s.Flags[i+8*j] = true
		}
	}
}

var ShopManager shopManager

type shopManager struct {
	shopTable map[int]map[int]*Shop
}

func (m *shopManager) init() {
	// ShopList was generated 2023-11-17 14:15:26 by https://xml-to-go.github.io/ in Ukraine.
	type ShopList struct {
		XMLName xml.Name `xml:"ShopList"`
		Text    string   `xml:",chardata"`
		Shop    []*Shop  `xml:"Shop"`
	}
	var shopList ShopList
	conf.XML(conf.PathCommon, "IGC_ShopList.xml", &shopList)
	shopItems := make(map[string]*item.PositionedItems)
	m.shopTable = make(map[int]map[int]*Shop)
	for _, shop := range shopList.Shop {
		pi, ok := shopItems[shop.FileName]
		if !ok {
			pi = &item.PositionedItems{
				Size:  MaxShopItemCount,
				Items: make([]*item.Item, MaxShopItemCount),
				Flags: make([]bool, MaxShopItemCount),
			}
			shop.PositionedItems = pi
			file := path.Join("Shops", shop.FileName)
			shop.Scan(file)
			shopItems[shop.FileName] = pi
		} else {
			shop.PositionedItems = pi
		}
		_, ok = m.shopTable[shop.NPCIndex]
		if !ok {
			m.shopTable[shop.NPCIndex] = make(map[int]*Shop)
		}
		m.shopTable[shop.NPCIndex][shop.MapNumber] = shop
	}
}

func (m *shopManager) ForEachShop(f func(int, int, int, int, int)) {
	for _, v1 := range m.shopTable {
		for _, v2 := range v1 {
			f(v2.NPCIndex, v2.MapNumber, v2.PosX, v2.PosY, v2.Dir)
		}
	}
}

func (m *shopManager) GetShopInventory(npcIndex, mapNumber int) []*item.Item {
	return m.shopTable[npcIndex][mapNumber].Items
}

func (m *shopManager) GetShopItem(npcIndex, mapNumber, position int) *item.Item {
	inventory := m.GetShopInventory(npcIndex, mapNumber)
	if inventory == nil {
		return nil
	}
	shopItem := inventory[position]
	if shopItem == nil {
		return nil
	}
	return shopItem.Copy()
}
