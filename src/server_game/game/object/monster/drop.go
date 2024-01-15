package monster

import (
	"encoding/xml"
	"math/rand"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/item"
)

const MaxMagicBookCount int = 100
const MaxNormalItemCount int = 1000

var DropManager dropManager

type itemDropRate struct {
	magicBook       int
	jewelOfBless    int
	jewelOfSoul     int
	jewelOfLife     int
	jewelOfCreation int
	jewelOfChaos    int
	normalItem      int
}

type dropManager struct {
	itemDropRate    []*itemDropRate
	magicBook       [][]*item.Item
	jewel           [][]*item.Item
	jewelOfBless    *item.Item
	jewelOfSoul     *item.Item
	jewelOfChaos    *item.Item
	jewelOfLife     *item.Item
	jewelOfCreation *item.Item
	normalItem      [][]*item.Item
	excellentItem   [][]*item.Item
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
			magicBook:       int(lm.MagicBook * 10000000),
			jewelOfBless:    int(lm.JewelOfBless * 10000000),
			jewelOfSoul:     int(lm.JewelOfSoul * 10000000),
			jewelOfLife:     int(lm.JewelOfLife * 10000000),
			jewelOfCreation: int(lm.JewelOfCreation * 10000000),
			jewelOfChaos:    int(lm.JewelOfChaos * 10000000),
			normalItem:      int(lm.Items * 10000000),
		}
	}

	// make all
	m.magicBook = make([][]*item.Item, levelCount)
	m.jewel = make([][]*item.Item, levelCount)
	m.normalItem = make([][]*item.Item, levelCount)
	m.excellentItem = make([][]*item.Item, levelCount)
	for _, monster := range MonsterTable {
		// make magic book
		m.makeMagicBook(monster.Level)
		// make normal item
		m.makeNormalItem(monster.Level, monster.MaxItemLevel)
	}
	// make jewel
	m.makeJewel()

	// for level, its := range m.magicBook {
	// 	for _, it := range its {
	// 		log.Printf("[level]%d [magicbook]%s\n", level, it.Annotation)
	// 	}
	// }
	// for level, its := range m.normalItem {
	// 	for _, it := range its {
	// 		log.Printf("[level]%d [magicbook]%s [excellent]%t\n", level, it.Annotation, it.excellent)
	// 	}
	// }
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
			return
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
	for section, its := range item.ItemTable {
		for index := range its {
			if len(m.normalItem[monsterLevel]) >= MaxNormalItemCount {
				return
			}
			level := item.ItemTable.GetItemLevel(section, index, monsterLevel)
			if level >= 0 {
				if level > maxItemLevel {
					level = maxItemLevel
				}
				it := item.NewItem(section, index)
				it.Level = level
				it.Calc()
				it.Durability = it.MaxDurability
				excellent := false
				if it.Type == item.TypeRegular {
					excellent = true
				}
				m.normalItem[monsterLevel] = append(m.normalItem[monsterLevel], it)
				if excellent {
					m.excellentItem[monsterLevel] = append(m.excellentItem[monsterLevel], it)
				}
			}
		}
	}
}

func (m *dropManager) makeJewel() {
	type jewel struct {
		section int
		index   int
		it      *item.Item
	}
	jewels := []*jewel{
		{14, 13, nil}, // Jewel of Bless
		{14, 14, nil}, // Jewel of Soul
		{12, 15, nil}, // Jewel of Chaos
		{14, 16, nil}, // Jewel of Life
		{14, 22, nil}, // Jewel of Creation
	}
	for _, jewel := range jewels {
		it := item.NewItem(jewel.section, jewel.index)
		it.Level = 0
		it.Calc()
		it.Durability = it.MaxDurability
		jewel.it = it
	}
	m.jewelOfBless = jewels[0].it
	m.jewelOfSoul = jewels[1].it
	m.jewelOfChaos = jewels[2].it
	m.jewelOfLife = jewels[3].it
	m.jewelOfCreation = jewels[4].it
}

func (m *dropManager) isValid(monsterLevel int) bool {
	if monsterLevel < 0 || monsterLevel >= len(m.itemDropRate) {
		return false
	}
	return true
}

func (m *dropManager) DropItem(monsterLevel int) *item.Item {
	if !m.isValid(monsterLevel) {
		return nil
	}
	dropRate := m.itemDropRate[monsterLevel]
	book := dropRate.magicBook
	bless := dropRate.jewelOfBless + book
	soul := dropRate.jewelOfSoul + bless
	life := dropRate.jewelOfLife + soul
	creation := dropRate.jewelOfCreation + life
	chaos := dropRate.jewelOfChaos + creation
	items := dropRate.normalItem + chaos
	number := rand.Intn(10000000)
	switch {
	case number >= 0 && number < book:
		its := m.magicBook[monsterLevel]
		n := len(its)
		if n <= 0 {
			return nil
		}
		return its[rand.Intn(n)]
	case number >= book && number < bless:
		return m.jewelOfBless
	case number >= bless && number < soul:
		return m.jewelOfSoul
	case number >= soul && number < life:
		return m.jewelOfLife
	case number >= life && number < creation:
		return m.jewelOfCreation
	case number >= creation && number < chaos:
		return m.jewelOfChaos
	case number >= chaos && number < items:
		its := m.normalItem[monsterLevel]
		n := len(its)
		if n <= 0 {
			return nil
		}
		return its[rand.Intn(n)]
	}
	return nil
}

func (m *dropManager) DropExcellentItem(monsterLevel int) *item.Item {
	if !m.isValid(monsterLevel) {
		return nil
	}
	its := m.excellentItem[monsterLevel]
	n := len(its)
	if n <= 0 {
		return nil
	}
	return its[rand.Intn(n)]
}
