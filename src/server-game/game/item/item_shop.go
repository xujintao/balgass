package item

import (
	"encoding/xml"
	"log/slog"
	"path"

	"github.com/xujintao/balgass/src/server-game/conf"
)

const MaxShopItemCount = 120 // 8*15

type Shop struct {
	*PositionedItems
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

func (s *Shop) init(file string) {
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
	for _, sit := range shopInventory.Item {
		it := NewItem(sit.Cat, sit.Index)
		it.Level = sit.Level
		it.Durability = sit.Durability
		it.Skill = sit.Skill
		it.Lucky = sit.Luck
		it.Addition = sit.Option << 2
		if it.KindA == KindAWing {
			it.RandWingAdditionKind()
		}
		it.DecodeExcellent(sit.Exc)
		it.Calc()
		if it.Durability == 0 {
			it.Durability = it.MaxDurability
		}
		if !s.AddItem(it) {
			slog.Error("Shop.AddItem",
				"shop", s.FileName, "item", it.Name)
			continue
		}
	}
}

func (s *Shop) AddItem(it *Item) bool {
	for i := 0; i < s.Size; i++ {
		if s.CanPlace(i, it) {
			s.Items[i] = it
			x := i % 8
			y := i / 8
			for column := x; column < x+it.Width; column++ {
				for row := y; row < y+it.Height; row++ {
					s.Flags[column+8*row] = true
				}
			}
			return true
		}
	}
	return false
}

func (s *Shop) CanPlace(position int, it *Item) bool {
	if it == nil {
		return false
	}
	if position < 0 || position >= s.Size || it.Width <= 0 || it.Height <= 0 {
		return false
	}
	x := position % 8
	y := position / 8
	if x+it.Width > 8 || y+it.Height > s.Size/8 {
		return false
	}
	for i := x; i < x+it.Width; i++ {
		for j := y; j < y+it.Height; j++ {
			if s.Flags[i+8*j] {
				return false
			}
		}
	}
	return true
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
	shopItems := make(map[string]*PositionedItems)
	m.shopTable = make(map[int]map[int]*Shop)
	for _, shop := range shopList.Shop {
		pi, ok := shopItems[shop.FileName]
		if !ok {
			pi = &PositionedItems{
				Size:  MaxShopItemCount,
				Items: make([]*Item, MaxShopItemCount),
				Flags: make([]bool, MaxShopItemCount),
			}
			shop.PositionedItems = pi
			file := path.Join("Shops", shop.FileName)
			shop.init(file)
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

func (m *shopManager) GetShopInventory(npcIndex, mapNumber int) []*Item {
	if m == nil {
		return nil
	}
	shop := m.shopTable[npcIndex][mapNumber]
	if shop == nil || shop.PositionedItems == nil {
		return nil
	}
	return shop.Items
}

func (m *shopManager) GetShopItem(npcIndex, mapNumber, position int) *Item {
	inventory := m.GetShopInventory(npcIndex, mapNumber)
	if position < 0 || position >= len(inventory) {
		return nil
	}
	sit := inventory[position]
	if sit == nil {
		return nil
	}
	return sit.Copy()
}
