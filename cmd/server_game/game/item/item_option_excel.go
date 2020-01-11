package item

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

type optionExcelCommon struct {
	ID         OptionExcelCommon `xml:"ID,attr"`
	Number     int               `xml:"Number,attr"`
	Value      int               `xml:"Value,attr"`
	ItemKindA1 int               `xml:"ItemKindA_1,attr"`
	ItemKindA2 int               `xml:"ItemKindA_2,attr"`
	ItemKindA3 int               `xml:"ItemKindA_3,attr"`
	Rate       int               `xml:"Rate,attr"`
	Name       string            `xml:"Name,attr"`
}

type optionExcelWing struct {
	ID        string `xml:"ID,attr"`
	Number    string `xml:"Number,attr"`
	Value     string `xml:"Value,attr"`
	ItemKindA string `xml:"ItemKindA,attr"`
	ItemKindB string `xml:"ItemKindB,attr"`
	Name      string `xml:"Name,attr"`
}
