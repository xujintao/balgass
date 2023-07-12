package maps

import (
	"encoding/xml"
	"log"
	"os"
	"path"

	"github.com/xujintao/balgass/src/server_game/conf"
)

func init() {
	// MapList was generated 2023-07-12 18:10:44 by https://xml-to-go.github.io/ in Ukraine.
	type MapList struct {
		XMLName     xml.Name `xml:"MapList"`
		Text        string   `xml:",chardata"`
		DefaultMaps struct {
			Text string `xml:",chardata"`
			Map  []struct {
				Text   string `xml:",chardata"`
				Number int    `xml:"Number,attr"`
				File   string `xml:"File,attr"`
			} `xml:"Map"`
		} `xml:"DefaultMaps"`
		CrywolfEventAttr struct {
			Text string `xml:",chardata"`
			Map  []struct {
				Text            string `xml:",chardata"`
				Number          string `xml:"Number,attr"`
				File            string `xml:"File,attr"`
				OccupationState string `xml:"OccupationState,attr"`
			} `xml:"Map"`
		} `xml:"CrywolfEventAttr"`
		KanturuEventAttr struct {
			Text string `xml:",chardata"`
			Map  []struct {
				Text      string `xml:",chardata"`
				Number    string `xml:"Number,attr"`
				File      string `xml:"File,attr"`
				MayaState string `xml:"MayaState,attr"`
			} `xml:"Map"`
		} `xml:"KanturuEventAttr"`
	}
	var mapList MapList
	conf.XML(conf.PathCommon, "IGC_MapList.xml", &mapList)

	MapTable = make(mapTable)
	for _, v := range mapList.DefaultMaps.Map {
		number := v.Number
		file := path.Join(conf.PathCommon, "MapTerrains", v.File)
		m := Map{}
		m.init(number, file)
		MapTable[v.Number] = &m
	}
}

var MapTable mapTable

type mapTable map[int]*Map

type Map struct {
	number   int
	width    int
	height   int
	buf      []byte
	mapItems []mapItem
}

const (
	MAX_WIDTH  = 255
	MAX_HEIGHT = 255
)

func (m *Map) init(number int, file string) {
	m.number = number
	buf, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("read file failed [file]%s [err]%v", file, err)
	}
	m.width = int(buf[1])
	m.height = int(buf[2])
	m.buf = buf[3:]
	m.mapItems = make([]mapItem, conf.Server.GameServerInfo.MaxObjectItemCount)
}

func (m *Map) GetAttr(x, y int) int {
	if !(x >= 0 && x <= MAX_WIDTH &&
		y >= 0 && y <= MAX_HEIGHT) {
		return 4
	}
	return int(m.buf[x+y<<8])
}
