package item

import (
	"math/rand"
	"path"

	"github.com/xujintao/balgass/cmd/server_game/conf"
)

const (
	maxOptionExcelCommon = 20
	maxOptionExcelWing   = 40
)

type OptionExcelCommon int

const (
	OptionExcelCommonIncMana OptionExcelCommon = iota
	OptionExcelCommonIncHP
	OptionExcelCommonIncAttackSpeed
	OptionExcelCommonIncAttackPercent
	OptionExcelCommonIncAttackLevel
	OptionExcelCommonIncExcellentDamage
	OptionExcelCommonIncZen
	OptionExcelCommonIncDefenseRate
	OptionExcelCommonReflectDamage
	OptionExcelCommonDecDamage
	OptionExcelCommonIncMaxMana
	OptionExcelCommonIncMaxHP
)

// (13,171) ~ (13,176)
// discard new pendant and ring
type optionExcelAccessory struct{}

type optionExcel struct {
	Common struct {
		Options []struct {
			ID         OptionExcelCommon `xml:"ID,attr"`
			Number     int               `xml:"Number,attr"`
			Value      int               `xml:"Value,attr"`
			ItemKindA1 itemKindA         `xml:"ItemKindA_1,attr"`
			ItemKindA2 itemKindA         `xml:"ItemKindA_2,attr"`
			ItemKindA3 itemKindA         `xml:"ItemKindA_3,attr"`
			Rate       int               `xml:"Rate,attr"`
			Name       string            `xml:"Name,attr"`
		} `xml:"Option"`
	} `xml:"Common"`
	Wings struct {
		Options []struct {
			ID        int       `xml:"ID,attr"`
			Number    int       `xml:"Number,attr"`
			Value     int       `xml:"Value,attr"`
			ItemKindA itemKindA `xml:"ItemKindA,attr"`
			ItemKindB itemKindB `xml:"ItemKindB,attr"`
			Name      string    `xml:"Name,attr"`
		} `xml:"Option"`
	} `xml:"Wings"`
	OptionDropRate struct {
		Common struct {
			One   int `xml:"One,attr"`
			Two   int `xml:"Two,attr"`
			Three int `xml:"Three,attr"`
			Four  int `xml:"Four,attr"`
			Five  int `xml:"Five,attr"`
			Six   int `xml:"Six,attr"`
		} `xml:"Common"`
	} `xml:"OptionDropRate"`
}

func (o *optionExcel) CommonRand(kindA itemKindA) (excel int) {
	var options [6]int
	var optionRates [6]int
	index := 0
	for _, v := range o.Common.Options {
		if v.ItemKindA1 == kindA || v.ItemKindA2 == kindA || v.ItemKindA3 == kindA {
			options[index] = v.Number
			optionRates[index] = v.Rate
			index++
		}
	}
	optionNum := dropRate(o.OptionDropRate.Common.One,
		o.OptionDropRate.Common.Two,
		o.OptionDropRate.Common.Three,
		o.OptionDropRate.Common.Four,
		o.OptionDropRate.Common.Five,
		o.OptionDropRate.Common.Six) + 1
	for optionNum > 0 {
		i := rand.Int() % index
		option := options[i]
		// optionRate := optionRates[i]
		if excel&option != option {
			excel |= option
			optionNum--
		}
	}
	return
}
func (o *optionExcel) CommonEffect(id int) {}

// func (o *optionExcel) WingRand() {}

// func (o *optionExcel) WingEffect(id int) {}

var OptionExcel optionExcel

func init() {
	conf.XML(path.Join(conf.PathCommon, "Items/IGC_ExcellentOptions"), &OptionExcel)
}

func dropRate(rates ...int) int {
	num := rand.Intn(10000)
	offset := 0
	for i, v := range rates {
		if num >= offset && num < v+offset {
			return i
		}
		offset += v
	}
	return 0
}
