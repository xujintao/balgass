package item

import (
	"fmt"
	"math/rand"

	"github.com/xujintao/balgass/src/server-game/conf"
)

var SocketManager socketManager

type socketManager struct {
	rates   [12][5]int
	twoHand [5]int
}

func (m *socketManager) init() {
	type rateConfig struct {
		CatID  int `xml:"CatID,attr"`
		Count1 int `xml:"Count1,attr"`
		Count2 int `xml:"Count2,attr"`
		Count3 int `xml:"Count3,attr"`
		Count4 int `xml:"Count4,attr"`
		Count5 int `xml:"Count5,attr"`
	}
	type config struct {
		SocketRateSettings struct {
			Sockets []rateConfig `xml:"Socket"`
			TwoHand *rateConfig  `xml:"TwoHandWeapons"`
		} `xml:"SocketRateSettings"`
	}
	var cfg config
	conf.XML(conf.PathCommon, "SocketSystem/IGC_SocketSystem.xml", &cfg)

	weights := func(label string, rate rateConfig) [5]int {
		values := [5]int{rate.Count1, rate.Count2, rate.Count3, rate.Count4, rate.Count5}
		total := 0
		for _, value := range values {
			if value < 0 || value > 10000 {
				panic(fmt.Errorf("SocketSystem/IGC_SocketSystem.xml: invalid %s slot weight %d", label, value))
			}
			total += value
		}
		if total == 0 {
			panic(fmt.Errorf("SocketSystem/IGC_SocketSystem.xml: %s has no positive slot weights", label))
		}
		return values
	}

	var seen [12]bool
	for _, rate := range cfg.SocketRateSettings.Sockets {
		if rate.CatID < 0 || rate.CatID >= len(m.rates) {
			panic(fmt.Errorf("SocketSystem/IGC_SocketSystem.xml: invalid section %d", rate.CatID))
		}
		if seen[rate.CatID] {
			panic(fmt.Errorf("SocketSystem/IGC_SocketSystem.xml: duplicate section %d", rate.CatID))
		}
		m.rates[rate.CatID] = weights(fmt.Sprintf("section %d", rate.CatID), rate)
		seen[rate.CatID] = true
	}
	for section, present := range seen {
		if !present {
			panic(fmt.Errorf("SocketSystem/IGC_SocketSystem.xml: missing section %d", section))
		}
	}
	if cfg.SocketRateSettings.TwoHand == nil {
		panic(fmt.Errorf("SocketSystem/IGC_SocketSystem.xml: missing TwoHandWeapons"))
	}
	m.twoHand = weights("TwoHandWeapons", *cfg.SocketRateSettings.TwoHand)
}

func (m *socketManager) RandomSlotCount(it *Item) int {
	if it == nil || it.ItemBase == nil || it.Type != TypeSocket {
		return 0
	}
	var rates [5]int
	if it.TwoHand {
		rates = m.twoHand
	} else if it.Section >= 0 && it.Section < len(m.rates) {
		rates = m.rates[it.Section]
	}
	total := 0
	for _, rate := range rates {
		total += rate
	}
	if total == 0 {
		return 0
	}
	n := rand.Intn(total)
	for i, rate := range rates {
		if n < rate {
			return i + 1
		}
		n -= rate
	}
	return 0
}

func (it *Item) MakeSocketSlots(count int) {
	if it == nil ||
		it.ItemBase == nil ||
		it.Type != TypeSocket ||
		count < 0 || count > 5 {
		return
	}
	it.SocketBonus = 0
	for i, slot := range []*int{
		&it.SocketSlot1,
		&it.SocketSlot2,
		&it.SocketSlot3,
		&it.SocketSlot4,
		&it.SocketSlot5} {
		*slot = 0
		if i < count {
			*slot = 1
		}
	}
}
