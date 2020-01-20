package item

import (
	"path"

	"github.com/xujintao/balgass/cmd/server_game/conf"
)

type SetItem struct {
	tiers    [4]int
	mixLevel [2]int
}

type SetEffect struct {
	Index int
	Value int
}

type SetKind struct {
	name        string
	count       int
	effects     [6]SetEffect
	effectsFull [5]SetEffect
	class       [8]bool
}

type set struct {
	items []map[int]*SetItem
	kinds map[int]*SetKind
}

func (o *set) IsSetItem(section, index int) bool {
	_, ok := o.items[section][index]
	return ok
}

func (o *set) GetSetKindIndex(secion, index int, optIndex int) int {
	if !o.IsSetItem(secion, index) {
		return -1
	}
	return o.items[secion][index].tiers[optIndex]
}

func (o *set) GetSet(kindIndex int) {

}

var Set set

func init() {
	type SetItemXml struct {
		DropRate struct {
			Sections []struct {
				Index    int `xml:"Index,attr"`
				DropRate int `xml:"DropRate,attr"`
			} `xml:"Section"`
		} `xml:"DropRate"`
		Sections []struct {
			Index int    `xml:"Index,attr"`
			Name  string `xml:"Name,attr"`
			Items []struct {
				Index     int `xml:"Index,attr"`
				TierI     int `xml:"TierI,attr"`
				TierII    int `xml:"TierII,attr"`
				TierIII   int `xml:"TierIII,attr"`
				TierIV    int `xml:"TierIV,attr"`
				MixLevelA int `xml:"MixLevelA,attr"`
				MixLevelB int `xml:"MixLevelB,attr"`
			} `xml:"Item"`
		} `xml:"Section"`
	}
	var setItem SetItemXml
	conf.XML(path.Join(conf.PathCommon, "Items/IGC_ItemSetType.xml"), &setItem)
	// convert
	Set.items = make([]map[int]*SetItem, len(setItem.Sections))
	for _, section := range setItem.Sections {
		items := make(map[int]*SetItem)
		for _, item := range section.Items {
			var sItem SetItem
			sItem.tiers[0] = item.TierI
			sItem.tiers[1] = item.TierII
			sItem.tiers[2] = item.TierIII
			sItem.tiers[3] = item.TierIV
			sItem.mixLevel[0] = item.MixLevelA
			sItem.mixLevel[1] = item.MixLevelB
			items[item.Index] = &sItem
		}
		Set.items[section.Index] = items
	}

	type SetKindXml struct {
		Sets []struct {
			Index          int    `xml:"Index,attr"`
			Name           string `xml:"Name,attr"`
			OptIdx11       int    `xml:"OptIdx1_1,attr"`
			OptVal11       int    `xml:"OptVal1_1,attr"`
			OptIdx21       int    `xml:"OptIdx2_1,attr"`
			OptVal21       int    `xml:"OptVal2_1,attr"`
			OptIdx12       int    `xml:"OptIdx1_2,attr"`
			OptVal12       int    `xml:"OptVal1_2,attr"`
			OptIdx22       int    `xml:"OptIdx2_2,attr"`
			OptVal22       int    `xml:"OptVal2_2,attr"`
			OptIdx13       int    `xml:"OptIdx1_3,attr"`
			OptVal13       int    `xml:"OptVal1_3,attr"`
			OptIdx23       int    `xml:"OptIdx2_3,attr"`
			OptVal23       int    `xml:"OptVal2_3,attr"`
			OptIdx14       int    `xml:"OptIdx1_4,attr"`
			OptVal14       int    `xml:"OptVal1_4,attr"`
			OptIdx24       int    `xml:"OptIdx2_4,attr"`
			OptVal24       int    `xml:"OptVal2_4,attr"`
			OptIdx15       int    `xml:"OptIdx1_5,attr"`
			OptVal15       int    `xml:"OptVal1_5,attr"`
			OptIdx25       int    `xml:"OptIdx2_5,attr"`
			OptVal25       int    `xml:"OptVal2_5,attr"`
			OptIdx16       int    `xml:"OptIdx1_6,attr"`
			OptVal16       int    `xml:"OptVal1_6,attr"`
			OptIdx26       int    `xml:"OptIdx2_6,attr"`
			OptVal26       int    `xml:"OptVal2_6,attr"`
			SpecialOptIdx1 int    `xml:"SpecialOptIdx1,attr"`
			SpecialOptVal1 int    `xml:"SpecialOptVal1,attr"`
			SpecialOptIdx2 int    `xml:"SpecialOptIdx2,attr"`
			SpecialOptVal2 int    `xml:"SpecialOptVal2,attr"`
			FullOptIdx1    int    `xml:"FullOptIdx1,attr"`
			FullOptVal1    int    `xml:"FullOptVal1,attr"`
			FullOptIdx2    int    `xml:"FullOptIdx2,attr"`
			FullOptVal2    int    `xml:"FullOptVal2,attr"`
			FullOptIdx3    int    `xml:"FullOptIdx3,attr"`
			FullOptVal3    int    `xml:"FullOptVal3,attr"`
			FullOptIdx4    int    `xml:"FullOptIdx4,attr"`
			FullOptVal4    int    `xml:"FullOptVal4,attr"`
			FullOptIdx5    int    `xml:"FullOptIdx5,attr"`
			FullOptVal5    int    `xml:"FullOptVal5,attr"`
			DarkWizard     bool   `xml:"DarkWizard,attr"`
			DarkKnight     bool   `xml:"DarkKnight,attr"`
			FairyElf       bool   `xml:"FairyElf,attr"`
			MagicGladiator bool   `xml:"MagicGladiator,attr"`
			DarkLord       bool   `xml:"DarkLord,attr"`
			Summoner       bool   `xml:"Summoner,attr"`
			RageFighter    bool   `xml:"RageFighter,attr"`
		} `xml:"SetItem"`
	}
	var setKind SetKindXml
	conf.XML(path.Join(conf.PathCommon, "Items/IGC_ItemSetOption.xml"), &setKind)
	// convert
	Set.kinds = make(map[int]*SetKind)
	for _, set := range setKind.Sets {
		var kind SetKind
		kind.effects[0].Index = set.OptIdx11
		kind.effects[0].Value = set.OptVal11
		if set.OptIdx11 != -1 {
			kind.count++
		}
		kind.effects[1].Index = set.OptIdx12
		kind.effects[1].Value = set.OptVal12
		if set.OptIdx12 != -1 {
			kind.count++
		}
		kind.effects[2].Index = set.OptIdx13
		kind.effects[2].Value = set.OptVal13
		if set.OptIdx13 != -1 {
			kind.count++
		}
		kind.effects[3].Index = set.OptIdx14
		kind.effects[3].Value = set.OptVal14
		if set.OptIdx14 != -1 {
			kind.count++
		}
		kind.effects[4].Index = set.OptIdx15
		kind.effects[4].Value = set.OptVal15
		if set.OptIdx15 != -1 {
			kind.count++
		}
		kind.effects[5].Index = set.OptIdx16
		kind.effects[5].Value = set.OptVal16
		if set.OptIdx16 != -1 {
			kind.count++
		}

		kind.effectsFull[0].Index = set.FullOptIdx1
		kind.effectsFull[0].Value = set.FullOptVal1
		kind.effectsFull[1].Index = set.FullOptIdx2
		kind.effectsFull[1].Value = set.FullOptVal2
		kind.effectsFull[2].Index = set.FullOptIdx3
		kind.effectsFull[2].Value = set.FullOptVal3
		kind.effectsFull[3].Index = set.FullOptIdx4
		kind.effectsFull[3].Value = set.FullOptVal4
		kind.effectsFull[4].Index = set.FullOptIdx5
		kind.effectsFull[4].Value = set.FullOptVal5

		Set.kinds[set.Index] = &kind
	}
}
