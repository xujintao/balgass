package item

import (
	"encoding/xml"
	"math/rand"

	"github.com/xujintao/balgass/src/server-game/conf"
)

func init() {
	ExcellentDropManager.init()
}

var ExcellentDropManager excellentDropManager

type excellentDropManager struct {
	regulars map[int][]int
	wings    map[int][]int
	multiple []int
}

func (m *excellentDropManager) init() {
	// ExcellentOptions was generated 2024-01-15 17:48:16 by https://xml-to-go.github.io/ in Ukraine.
	type ExcellentOptions struct {
		XMLName xml.Name `xml:"ExcellentOptions"`
		Text    string   `xml:",chardata"`
		Common  struct {
			Text   string `xml:",chardata"`
			Option []struct {
				Text       string `xml:",chardata"`
				ID         int    `xml:"ID,attr"`
				Number     int    `xml:"Number,attr"`
				Value      int    `xml:"Value,attr"`
				ItemKindA1 int    `xml:"ItemKindA_1,attr"`
				ItemKindA2 int    `xml:"ItemKindA_2,attr"`
				ItemKindA3 int    `xml:"ItemKindA_3,attr"`
				Rate       int    `xml:"Rate,attr"`
				Name       string `xml:"Name,attr"`
			} `xml:"Option"`
		} `xml:"Common"`
		Wings struct {
			Text   string `xml:",chardata"`
			Option []struct {
				Text      string `xml:",chardata"`
				ID        int    `xml:"ID,attr"`
				Number    int    `xml:"Number,attr"`
				Value     int    `xml:"Value,attr"`
				ItemKindA int    `xml:"ItemKindA,attr"`
				ItemKindB int    `xml:"ItemKindB,attr"`
				Name      string `xml:"Name,attr"`
			} `xml:"Option"`
		} `xml:"Wings"`
		OptionDropRate struct {
			Text   string `xml:",chardata"`
			Common struct {
				Text  string `xml:",chardata"`
				One   int    `xml:"One,attr"`
				Two   int    `xml:"Two,attr"`
				Three int    `xml:"Three,attr"`
				Four  int    `xml:"Four,attr"`
				Five  int    `xml:"Five,attr"`
				Six   int    `xml:"Six,attr"`
			} `xml:"Common"`
		} `xml:"OptionDropRate"`
	}
	var excellentOptions ExcellentOptions
	conf.XML(conf.PathCommon, "Items/IGC_ExcellentOptions.xml", &excellentOptions)
	// regular
	m.regulars = make(map[int][]int)
	for _, v := range excellentOptions.Common.Option {
		m.regulars[v.ItemKindA1] = append(m.regulars[v.ItemKindA1], v.Number)
		m.regulars[v.ItemKindA2] = append(m.regulars[v.ItemKindA2], v.Number)
	}
	// wing
	m.wings = make(map[int][]int)
	for _, v := range excellentOptions.Wings.Option {
		m.wings[v.ItemKindB] = append(m.wings[v.ItemKindB], v.Number)
	}
	// multiple
	rate := excellentOptions.OptionDropRate.Common
	m.multiple = append(m.multiple,
		rate.One,
		rate.Two,
		rate.Three,
		rate.Four,
		rate.Five,
		rate.Six,
	)
}

func (m *excellentDropManager) dropExcellentCount() int {
	num := rand.Intn(10000)
	offset := 0
	for i, v := range m.multiple {
		if num >= offset && num < v+offset {
			return i + 1
		}
		offset += v
	}
	return 0
}

func (m *excellentDropManager) DropExcellent(kindA itemKindA, kindB itemKindB) []int {
	var pool []int
	switch kindA {
	case KindAWeapon, KindAPendant, KindAArmor, KindARing:
		pool = m.regulars[int(kindA)]
	case KindAWing:
		pool = m.wings[int(kindB)]
	}
	n := len(pool)
	if n <= 0 {
		return nil
	}
	cnt := m.dropExcellentCount()
	s1 := make(map[int]struct{})
	for cnt > 0 {
		i := rand.Intn(n)
		v := pool[i]
		_, ok := s1[v]
		if !ok {
			s1[v] = struct{}{}
			cnt--
		}
	}
	var s2 []int
	for k := range s1 {
		s2 = append(s2, k)
	}
	return s2
}
