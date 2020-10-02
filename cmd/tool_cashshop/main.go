package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

var conf config

type config struct {
	Version string `json:"version"`
	Short   bool   `json:"short"`
	Unit    string `json:"unit"`
}

type ItemInfo struct {
	GUID        int    `xml:"GUID,attr" dat:"guid"`
	ID          int    `xml:"ID,attr" dat:"id"`
	Cat         int    `xml:"Cat,attr" dat:"group"`
	Index       int    `xml:"Index,attr" dat:"index"`
	Level       int    `xml:"Level,attr" dat:"level"`
	Durability  int    `xml:"Durability,attr" dat:"durability"`
	Skill       int    `xml:"Skill,attr" dat:"skill"`
	Luck        int    `xml:"Luck,attr" dat:"luck"`
	Option      int    `xml:"Option,attr" dat:"option"`
	Exc         int    `xml:"Exc,attr" dat:"excel"`
	Set         int    `xml:"Set,attr" dat:"set"`
	SocketCount int    `xml:"SocketCount,attr" dat:"socket"`
	Element     int    `xml:"Element,attr" dat:"-"`
	Type        int    `xml:"Type,attr" dat:"type"`
	Period      int    `xml:"Duration,attr" dat:"period"`
	Comment     string `xml:"-" dat:"description"`
}

type CashItemInfo struct {
	// self-closing current is milestone
	// https://github.com/golang/go/issues/21399
	Infos []ItemInfo `xml:"Item,allowempty"`
}
type Item struct {
	GUID             int    `xml:"GUID,attr" dat:"guid"`
	Index            int    `xml:"iIndex,attr" dat:"index"`        // be consist with client
	SubIndex         int    `xml:"iSubIndex,attr" dat:"sub_index"` // depends on CoinType
	OptionSelect     int    `xml:"OptionSelect,attr" dat:"option_select"`
	PackageID        int    `xml:"PackageID,attr" dat:"package_id"`
	CoinType         int    `xml:"CoinType,attr" dat:"coin_type"`
	CoinValue        int    `xml:"CoinValue,attr" dat:"coin_value"`
	UniqueID1        int    `xml:"UniqueID1,attr" dat:"unique_id_1"`
	UniqueID2        int    `xml:"UniqueID2,attr" dat:"unique_id_2"`
	ShopCategory     int    `xml:"ShopCategory,attr" dat:"shop_category"`
	GPRewardValue    int    `xml:"GPRewardValue,attr" dat:"gp_reward_value"`
	CanBuy           int    `xml:"CanBuy,attr" dat:"can_buy"`
	CanGift          int    `xml:"CanGift,attr" dat:"can_gift"`
	RandomItemSelect int    `xml:"RandomItemSelect,attr" dat:"random_item_select"`
	Comment          string `xml:"-" dat:"description"`
}
type CashItemList struct {
	Items []Item `xml:"Item"`
}

type Package struct {
	GUID         int    `xml:"GUID,attr" dat:"guid"`
	ID           int    `xml:"ID,attr" dat:"id"`
	ItemSequence int    `xml:"ItemSequence,attr" dat:"item_sequence"`
	UniqueID1    int    `xml:"UniqueID1,attr" dat:"unique_id_1"`
	UniqueID2    int    `xml:"UniqueID2,attr" dat:"unique_id_2"`
	Comment      string `xml:"-" dat:"description"`
}
type CashItemPackage struct {
	Packages []Package `xml:"Package"`
}

func newBufioScanner(r io.Reader) *bufio.Scanner {
	bufr := bufio.NewReader(r)

	// detect bom
	bom, _, err := bufr.ReadRune()
	if err != nil {
		panic(err)
	}
	if bom != '\uFEFF' {
		bufr.UnreadRune()
	}
	return bufio.NewScanner(bufr)
}

func mustAtoi(a string) int {
	if a == "" || a == "test" {
		return 0
	}
	i, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	return i
}

func toDAT(i interface{}, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	bufw := bufio.NewWriter(f)

	// write head
	t := reflect.TypeOf(i) // ex: CashItemInfo or *CashItemInfo
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	field := t.Field(0) // ex: Infos []ItemInfo
	if field.Type.Kind() != reflect.Slice {
		return fmt.Errorf("%s.%s must be a slice", t.String(), field.Name)
	}
	t = field.Type.Elem() // ex: ItemInfo
	bufw.WriteString("// ")
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("dat")
		if tag == "-" {
			continue
		}
		bufw.WriteString(tag)
		if i < t.NumField()-1 {
			bufw.WriteString("  ")
		}
	}
	bufw.WriteByte('\n')

	// 1
	bufw.WriteString("1")
	bufw.WriteByte('\n')

	// body
	v := reflect.ValueOf(i) // ex: CashItemInfo or *CashItemInfo
	if v.Kind() == reflect.Ptr {
		v = reflect.ValueOf(i).Elem()
	}
	v = v.Field(0) // ex: Infos []ItemInfo
	for i := 0; i < v.Len(); i++ {
		v := v.Index(i) // ex: ItemInfo
		t := v.Type()   // t and tField are used to write comment field, can we use field.Type() instead ?
		for j := 0; j < v.NumField(); j++ {
			field := v.Field(j)
			tField := t.Field(j)
			if tField.Tag.Get("dat") == "-" {
				continue
			}
			if j == v.NumField()-1 && tField.Tag.Get("dat") == "description" {
				bufw.WriteString("//")
			}
			value := ""
			switch v := field.Interface().(type) {
			case int:
				value = strconv.Itoa(v)
			case string:
				value = v
			}
			bufw.WriteString(value)
			if j < v.NumField()-1 {
				bufw.WriteByte('\t')
			}
		}
		bufw.WriteByte('\n')
	}

	// end
	bufw.WriteString("end")
	bufw.WriteByte('\n')
	bufw.Flush()
	return nil
}

func toXML(v interface{}, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString(xml.Header); err != nil {
		return err
	}
	enc := xml.NewEncoder(f)
	enc.Indent("", "    ")
	return enc.Encode(v)
}

func code(section, index int) int {
	return section<<9 + index
}

func split(s string) []string {
	sep := ""
	switch {
	case strings.Contains(s, "^@"):
		sep = "^@"
	case strings.Contains(s, "@"):
		sep = "@"
	}
	return strings.Split(s, sep)
}

func replaceFile(fileName string) {
	oldFile := path.Join(conf.Version, fileName)
	newFile := path.Join(conf.Version, "out", fileName)
	s, err := ioutil.ReadFile(oldFile)
	if err != nil {
		log.Fatal(err)
	}
	// remove utf bom header
	if len(s) > 3 && s[0] == 0xEF && s[1] == 0xBB && s[2] == 0xBF {
		s = s[3:]
	}
	content := strings.ReplaceAll(string(s), "^@", "@")
	ioutil.WriteFile(newFile, []byte(content), 0666)
}

// IBSProduct.txt
// [0]: guid
// [1]: item name
// [2]: consumption/period/term/number
// [3]: value
// [4]: unit
// [5]: coin value
// [6]: id
// [7]: ---------------- unknown, always 142
// [8]: ---------------- unknown, always 145
// [9]: ---------------- unknown, always 1
// [10]: ---------------- unknown, always 144
// [11]: ---------------- unknown, always 673
// [12]: ---------------- unknown, always 518
// [13]: base code
// [14]: 7: number, 10/2: period
// [15]: ---------------- unknown, always 138
// [16]: ---------------- unknown, always 680
func convertItemInfo() (itemInfos map[int]ItemInfo) {
	var cii CashItemInfo
	itemInfos = make(map[int]ItemInfo)

	f, err := os.Open(path.Join(conf.Version, "IBSProduct.txt"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := newBufioScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		values := split(line)
		baseCode := mustAtoi(values[13])
		baseSection := baseCode >> 9
		baseIndex := baseCode % 512
		itemInfo := ItemInfo{
			GUID:    mustAtoi(values[0]),
			ID:      mustAtoi(values[6]),
			Cat:     baseSection,
			Index:   baseIndex,
			Comment: values[1],
		}
		kind := mustAtoi(values[14])
		switch kind {
		case 7: // quantity
			switch code(itemInfo.Cat, itemInfo.Index) {
			case code(13, 97), code(13, 98), code(14, 91), code(14, 169): // 魔剑士角色卡片 圣导师角色卡片 召唤术士角色卡 格斗家角色卡片
				itemInfo.Type = 3 // character card
			default:
				itemInfo.Type = 0 // quantity/number
			}
			itemInfo.Durability = mustAtoi(values[3])
		case 2, 10: // duration
			// type
			switch code(itemInfo.Cat, itemInfo.Index) {
			case code(13, 43), code(13, 44), code(13, 45): // 经验印章 神圣印章 贡献印章
				itemInfo.Type = 1
			case code(13, 62), code(13, 63): // 大师经验印章 大师神圣印章
				itemInfo.Type = 1
			case code(13, 93), code(13, 94): // 大师等级经验印章 大师等级神圣印章
				itemInfo.Type = 1
			case code(13, 103), code(13, 104), code(13, 105): // 组队经验值符咒 最大AG提升光环 最大SD提升光环
				itemInfo.Type = 1
			case code(13, 128), code(13, 129), code(13, 130), code(13, 132), code(13, 134): // 神鹰雕像 山羊雕像 兽人符文 黄金兽人符文 破旧铁蹄
				itemInfo.Type = 1
			case code(14, 72), code(14, 73), code(14, 74), code(14, 75), code(14, 76), code(14, 77): // 加速卷轴 防御卷轴 愤怒卷轴 魔力卷轴 体力卷轴 魔法卷轴
				itemInfo.Type = 1
			case code(14, 97), code(14, 98): // 幸运一击卷轴 卓越一击卷轴
				itemInfo.Type = 1
			case code(14, 140): // 恢复卷轴
				itemInfo.Type = 1
			default:
				itemInfo.Type = 2
				itemInfo.Durability = 255
			}
			// period
			period := mustAtoi(values[3])
			unit := values[4]
			if strings.Contains(unit, conf.Unit) {
				period /= 60
			}
			itemInfo.Period = period
		}
		cii.Infos = append(cii.Infos, itemInfo)
		itemInfos[itemInfo.GUID] = itemInfo // used by convertItemList
	}
	if err := toXML(&cii, path.Join(conf.Version, "out/IGC_CashItem_Info.xml")); err != nil {
		log.Fatal(err)
	}
	if err := toDAT(&cii, path.Join(conf.Version, "out/IGCCashItemInfo.dat")); err != nil {
		log.Fatal(err)
	}
	return
}

// IBSPackage.txt
// [0]: row category, refer to IBSCategory.txt
// [1]: row category index
// [2]: Index, equal to IGCCashItemList.dat Index, used to identity items between server and client
// [3]: Item Name
// [4]: ---------------- unknown, always 170
// [5]: coin value
// [6]: item description
// [7]: ---------------- unknown, always empty
// [8]: ---------------- unknown, always 182
// [9]: ---------------- unknown, always 185
// [10]: ---------------- unknown, maybe date begin
// [11]: ---------------- unknown, maybe date end
// [12]: ---------------- unknown, always 177
// [13]: ---------------- unknown, always 1
// [14]: cash or goblin point
// [15]: 游戏点 或者 哥布林点数
// [16]: ---------------- unknown, always 181
// [17]: ---------------- unknown, always 200
// [18]: ---------------- unknown, always 0
// [19]: unique ID1
// [20]: item base code
// [21]: ---------------- unknown
// [22]: ---------------- unknown
// [23]: unique ID2
// [24]: coin type, 0: cash, 1: goblin
// [25]: ---------------- optional(necessary for CashShop Editor), 508: cash, 0: goblin
// [26]: ---------------- optional(necessary for CashShop Editor), unknown
func convertItemList(itemInfos map[int]ItemInfo) {
	var cil CashItemList
	var cip CashItemPackage

	f, err := os.Open(path.Join(conf.Version, "IBSPackage.txt"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var bufw *bufio.Writer

	packageID := 1
	scanner := newBufioScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		values := split(line)
		item := Item{
			GUID:         len(cil.Items),
			ShopCategory: mustAtoi(values[0]),
			Index:        mustAtoi(values[2]),
			Comment:      values[3],
			CoinValue:    mustAtoi(values[5]),
			// SubIndex:     mustAtoi(values[25]),
			CanBuy:  1,
			CanGift: 1,
		}

		// coin type
		switch mustAtoi(values[24]) {
		case 0: // cash
			item.CoinType = 0
			if conf.Short {
				item.SubIndex = 0
			} else {
				item.SubIndex = 508
			}
		case 1: // goblin
			item.CoinType = 2
			if conf.Short {
				item.SubIndex = 1
			} else {
				item.SubIndex = 0
			}
		}

		// fix IBSPackage.txt
		if conf.Short {
			if bufw == nil {
				IBSPackage2, err := os.Create(path.Join(conf.Version, "out/IBSPackage.txt"))
				if err != nil {
					panic(err)
				}
				defer IBSPackage2.Close()
				bufw = bufio.NewWriter(IBSPackage2)
			}
			suffix1 := ""
			suffix2 := "669"
			switch mustAtoi(values[24]) {
			case 0:
				suffix1 = "508"
			case 1:
				suffix1 = "0"
			}
			line = strings.Join([]string{line, suffix1, suffix2}, "^@")
			line = strings.ReplaceAll(line, "^@", "@")
			bufw.WriteString(line)
			bufw.WriteByte('\n')
		}

		if values[23] == "" {
			// package UniqueID1
			uniqueID1 := strings.TrimSuffix(values[19], "|")
			id1s := strings.Split(uniqueID1, "|")
			// here, need to pre-handle some special id1 as goblin point and assign result to GPRewardValue ?
			for i, id1 := range id1s {
				p := Package{
					ID:           packageID,
					GUID:         i,
					ItemSequence: i,
					UniqueID1:    mustAtoi(id1),
				}
				itemInfo, ok := itemInfos[p.UniqueID1]
				if ok {
					p.UniqueID2 = itemInfo.ID
					p.Comment = itemInfo.Comment
				} else {
					log.Printf("can not find id(id2) by id(id1):%d", p.UniqueID1)
				}
				cip.Packages = append(cip.Packages, p)
			}
			item.PackageID = packageID
			packageID++
			cil.Items = append(cil.Items, item)
		} else {
			uniqueID1 := strings.TrimSuffix(values[19], "|")
			item.UniqueID1 = mustAtoi(uniqueID1)
			uniqueID2 := strings.TrimSuffix(values[23], "|")
			id2s := strings.Split(uniqueID2, "|")
			for _, id2 := range id2s {
				item.GUID = len(cil.Items)
				item.UniqueID2 = mustAtoi(id2)
				if len(id2s) > 1 {
					item.OptionSelect = item.UniqueID2
				}
				cil.Items = append(cil.Items, item)
			}
		}
	}
	if conf.Short && bufw != nil {
		bufw.Flush()
		// IBSCategory2, err := os.Create(path.Join(conf.Version, "out/IBSCategory.txt"))
		// if err != nil {
		// 	log.Fatal(err)
		// }
		replaceFile("IBSCategory.txt")
		replaceFile("IBSProduct.txt")
	}
	if err := toXML(&cil, path.Join(conf.Version, "out/IGC_CashItem_List.xml")); err != nil {
		log.Fatal(err)
	}
	if err := toDAT(&cil, path.Join(conf.Version, "out/IGCCashItemList.dat")); err != nil {
		log.Fatal(err)
	}
	if err := toXML(&cip, path.Join(conf.Version, "out/IGC_CashItem_Package.xml")); err != nil {
		log.Fatal(err)
	}
	if err := toDAT(&cip, path.Join(conf.Version, "out/IGCCashItemPackages.dat")); err != nil {
		log.Fatal(err)
	}
}

func main() {
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(f).Decode(&conf)

	// create out directory
	err = os.MkdirAll(path.Join(conf.Version, "out"), 0777)
	if err != nil {
		log.Fatal(err)
	}

	// IBSProduct.txt -> IGCCashItemInfo.dat
	itemInfos := convertItemInfo()

	// IBSPackage.txt -> IGCCashItemList.dat & IGCCashItemPackages.dat
	convertItemList(itemInfos)
}
