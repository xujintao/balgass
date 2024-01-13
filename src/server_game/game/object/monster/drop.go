package monster

import (
	"encoding/xml"
	"math/rand"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/item"
)

const MaxMagicBookCount int = 100

var DropManager dropManager

type itemDropRate struct {
	MagicBook       float64
	JewelOfBless    float64
	JewelOfSoul     float64
	JewelOfLife     float64
	JewelOfCreation float64
	JewelOfChaos    float64
	NormalItem      float64
}

type dropManager struct {
	itemDropRate []*itemDropRate
	magicBook    [][]*item.Item
	jewel        [][]*item.Item
	normalItem   [][]*item.Item
}

func (m *dropManager) init() {
	// DropRate was generated 2024-01-13 17:14:47 by https://xml-to-go.github.io/ in Ukraine.
	type DropRate struct {
		XMLName xml.Name `xml:"DropRate"`
		Text    string   `xml:",chardata"`
		Monster []struct {
			Text            string  `xml:",chardata"`
			Level           int     `xml:"Level,attr"`
			MagicBook       float64 `xml:"MagicBook,attr"`
			JewelOfBless    float64 `xml:"JewelOfBless,attr"`
			JewelOfSoul     float64 `xml:"JewelOfSoul,attr"`
			JewelOfLife     float64 `xml:"JewelOfLife,attr"`
			JewelOfCreation float64 `xml:"JewelOfCreation,attr"`
			JewelOfChaos    float64 `xml:"JewelOfChaos,attr"`
			Items           float64 `xml:"Items,attr"`
		} `xml:"Monster"`
	}
	var dropRate DropRate
	conf.XML(conf.PathCommon, "IGC_MonsterItemDropRate.xml", &dropRate)
	levelCount := len(dropRate.Monster)
	m.itemDropRate = make([]*itemDropRate, levelCount)
	for _, lm := range dropRate.Monster {
		m.itemDropRate[lm.Level] = &itemDropRate{
			MagicBook:       lm.MagicBook,
			JewelOfBless:    lm.JewelOfBless,
			JewelOfSoul:     lm.JewelOfSoul,
			JewelOfLife:     lm.JewelOfLife,
			JewelOfCreation: lm.JewelOfCreation,
			JewelOfChaos:    lm.JewelOfChaos,
			NormalItem:      lm.Items,
		}
	}

	// make all
	m.magicBook = make([][]*item.Item, levelCount)
	m.jewel = make([][]*item.Item, levelCount)
	m.normalItem = make([][]*item.Item, levelCount)
	for _, monster := range MonsterTable {
		// make magic book
		m.makeMagicBook(monster.Level)
		// make normal item
		m.makeNormalItem(monster.Level, monster.MaxItemLevel)
	}
	// make jewel
	m.makeJewel()
}

func (m *dropManager) makeMagicBook(monsterLevel int) {
	type book struct {
		section int
		index   int
	}
	var books []*book
	for index, it := range item.ItemTable[12] {
		if it.KindA == item.KindASkill {
			books = append(books, &book{12, index})
		}
	}
	for index := range item.ItemTable[15] {
		books = append(books, &book{15, index})
	}
	count := 200
	for count > 0 {
		if len(m.magicBook[monsterLevel]) >= MaxMagicBookCount {
			break
		}
		n := rand.Intn(len(books))
		b := books[n]
		level := item.ItemTable.GetItemLevel(b.section, b.index, monsterLevel)
		if level >= 0 {
			it := item.NewItem(b.section, b.index)
			it.Level = 0
			if item.Code(b.section, b.index) == item.Code(12, 11) {
				it.Level = level
			}
			it.Calc()
			it.Durability = it.MaxDurability
			m.magicBook[monsterLevel] = append(m.magicBook[monsterLevel], it)
		}
		count--
	}
}

func (m *dropManager) makeNormalItem(monsterLevel, maxItemLevel int) {

}

func (m *dropManager) makeJewel() {

}

func (m *dropManager) DropMagicBook() {

}
