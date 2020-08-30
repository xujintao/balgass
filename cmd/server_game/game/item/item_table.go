package item

import (
	"fmt"
	"path"

	"github.com/xujintao/balgass/cmd/server_game/conf"
)

func init() {
	type itemListConfig struct {
		Sections []*struct {
			Index string      `xml:"Index,attr"`
			Name  string      `xml:"Name,attr"`
			Items []*ItemBase `xml:"Item"`
		} `xml:"Section"`
	}
	var itemList itemListConfig
	conf.XML(path.Join(conf.PathCommon, "Items/IGC_ItemList.xml"), &itemList)

	// [][]array -> []map
	ItemTable = make(itemTable, len(itemList.Sections))
	for i, section := range itemList.Sections {
		ItemTable[i] = make(map[int]*ItemBase)
		for _, v := range section.Items {
			ItemTable[i][v.Index] = v
		}
	}
}

func Code(section, index int) int {
	return section*512 + index
}

var ItemTable itemTable

type itemTable []map[int]*ItemBase

func (table itemTable) GetItemBase(i, j int) (*ItemBase, error) {
	if i >= len(table) {
		return nil, fmt.Errorf("item section over bound")
	}
	items := table[i]
	item, ok := items[j]
	if ok == false {
		return nil, fmt.Errorf("item index over bound")
	}
	return item, nil
}

func (table itemTable) GetItemBaseMust(i, j int) *ItemBase {
	return table[i][j]
}
