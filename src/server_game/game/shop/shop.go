package shop

import (
	"encoding/xml"
	"path"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/item"
)

func init() {
	ShopManager.init()
}

type Shop struct {
	NPCIndex     int    `xml:"NPCIndex,attr"`
	MapNumber    int    `xml:"MapNumber,attr"`
	PosX         int    `xml:"PosX,attr"`
	PosY         int    `xml:"PosY,attr"`
	Dir          int    `xml:"Dir,attr"`
	VipType      int    `xml:"VipType,attr"`
	GMShop       bool   `xml:"GMShop,attr"`
	FileName     string `xml:"FileName,attr"`
	Inventory    []*item.Item
	MossMerchant bool `xml:"MossMerchant,attr"`
	VIPType      int  `xml:"VIPType,attr"`
	BattleCore   bool `xml:"BattleCore,attr"`
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
	shopInventory := make(map[string][]*item.Item)
	m.shopTable = make(map[int]map[int]*Shop)
	for _, shop := range shopList.Shop {
		file := shop.FileName
		v, ok := shopInventory[file]
		if !ok {
			file = path.Join("Shops", file)
			var ShopInventory ShopInventory
			conf.XML(conf.PathCommon, file, &ShopInventory)
			var inventory []*item.Item
			for _, sitem := range ShopInventory.Item {
				item := item.NewItem(sitem.Cat, sitem.Index)
				item.Level = sitem.Level
				item.Durability = sitem.Durability
				item.Skill = sitem.Skill
				item.Lucky = sitem.Luck
				item.Addition = sitem.Option << 2
				inventory = append(inventory, item)
			}
			shopInventory[file] = inventory
			v = inventory
		}
		shop.Inventory = v
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
