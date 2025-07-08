package monster

import (
	"encoding/xml"
	"log/slog"
	"math/rand"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/item"
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

type dropItem struct {
	section    int
	index      int
	annotation string // for debug
	level      int
}

type dropManager struct {
	itemDropRate    []*itemDropRate
	magicBook       [][]dropItem
	jewelOfBless    dropItem
	jewelOfSoul     dropItem
	jewelOfChaos    dropItem
	jewelOfLife     dropItem
	jewelOfCreation dropItem
	normalItem      [][]dropItem
	excellentItem   [][]dropItem
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
	m.magicBook = make([][]dropItem, levelCount)
	m.normalItem = make([][]dropItem, levelCount)
	m.excellentItem = make([][]dropItem, levelCount)
	for _, monster := range MonsterTable {
		// make magic book
		m.makeMagicBook(monster.Level)
		// make item
		m.makeItem(monster.Level, monster.MaxItemLevel)
	}
	// make jewel
	m.makeJewel()

	// debug
	if conf.ServerEnv.Debug {
		for mlevel := range m.itemDropRate {
			for _, dit := range m.magicBook[mlevel] {
				slog.Debug("magic book", "monster level", mlevel, "annotation", dit.annotation, "level", dit.level)
			}
			for _, dit := range m.normalItem[mlevel] {
				slog.Debug("normal item", "monster level", mlevel, "annotation", dit.annotation, "level", dit.level)
			}
			for _, dit := range m.excellentItem[mlevel] {
				slog.Debug("excellent item", "monster level", mlevel, "annotation", dit.annotation, "level", dit.level)
			}
		}
	}
}

func (m *dropManager) makeMagicBook(monsterLevel int) {
	type book struct {
		section    int
		index      int
		annotation string // for debug
	}
	var books []*book
	for index, it := range item.ItemTable[12] {
		if it.KindA == item.KindASkill {
			books = append(books, &book{12, index, it.Annotation})
		}
	}
	for index, it := range item.ItemTable[15] {
		books = append(books, &book{15, index, it.Annotation})
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
			if item.Code(b.section, b.index) != item.Code(12, 11) {
				level = 0
			}
			dit := dropItem{
				section:    b.section,
				index:      b.index,
				annotation: b.annotation,
				level:      level,
			}
			m.magicBook[monsterLevel] = append(m.magicBook[monsterLevel], dit)
		}
		count--
	}
}

func (m *dropManager) makeItem(monsterLevel, maxItemLevel int) {
	for section, its := range item.ItemTable {
		for index, it := range its {
			if len(m.normalItem[monsterLevel]) >= MaxNormalItemCount {
				return
			}
			level := item.ItemTable.GetItemLevel(section, index, monsterLevel)
			if level >= 0 {
				if level > maxItemLevel {
					level = maxItemLevel
				}
				dit := dropItem{
					section:    section,
					index:      index,
					annotation: it.Annotation,
					level:      level,
				}
				m.normalItem[monsterLevel] = append(m.normalItem[monsterLevel], dit)
				if it.Option {
					m.excellentItem[monsterLevel] = append(m.excellentItem[monsterLevel], dit)
				}
			}
		}
	}
}

func (m *dropManager) makeJewel() {
	make := func(section, index int, annotation string) dropItem {
		return dropItem{
			section:    section,
			index:      index,
			annotation: annotation,
			level:      0,
		}
	}
	m.jewelOfBless = make(14, 13, "Jewel of Bless")
	m.jewelOfSoul = make(14, 14, "Jewel of Soul")
	m.jewelOfLife = make(14, 16, "Jewel of Life")
	m.jewelOfCreation = make(14, 22, "Jewel of Creation")
	m.jewelOfChaos = make(12, 15, "Jewel of Chaos")
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
	var dit dropItem
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
	case number >= 0 && number < book: // magic book
		its := m.magicBook[monsterLevel]
		n := len(its)
		if n <= 0 {
			return nil
		}
		dit = its[rand.Intn(n)]
	case number >= book && number < bless: // jewel of bless
		dit = m.jewelOfBless
	case number >= bless && number < soul: // jewel of soul
		dit = m.jewelOfSoul
	case number >= soul && number < life: // jewel of life
		dit = m.jewelOfLife
	case number >= life && number < creation: // jewel of creation
		dit = m.jewelOfCreation
	case number >= creation && number < chaos: // jewel of chaos
		dit = m.jewelOfChaos
	case number >= chaos && number < items: // normal item
		its := m.normalItem[monsterLevel]
		n := len(its)
		if n <= 0 {
			return nil
		}
		dit = its[rand.Intn(n)]
	}
	it := item.NewItem(dit.section, dit.index)
	if it == nil {
		slog.Error("DropItem", "error", "item not found", "section", dit.section, "index", dit.index)
		return nil
	}
	it.Level = dit.level
	return nil
}

func (m *dropManager) DropItemExcellent(monsterLevel int) *item.Item {
	if !m.isValid(monsterLevel) {
		return nil
	}
	its := m.excellentItem[monsterLevel]
	n := len(its)
	if n <= 0 {
		return nil
	}
	dit := its[rand.Intn(n)]
	it := item.NewItem(dit.section, dit.index)
	if it == nil {
		slog.Error("DropItemExcellent", "error", "item not found", "section", dit.section, "index", dit.index)
		return nil
	}
	it.Level = dit.level
	return it
}
