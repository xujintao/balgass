package maps

import (
	"encoding/xml"
	"log"
	"math/rand"
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

	rects := map[int]*Rect{
		Lorencia:    {130, 116, 151, 137},
		Dungeon:     {106, 236, 112, 243},
		Devias:      {197, 35, 218, 50},
		Noria:       {174, 101, 187, 125},
		LostTower:   {201, 70, 213, 81},
		Exile:       {89, 135, 90, 136},
		Arena:       {89, 135, 90, 136},
		Atlans:      {14, 11, 27, 23},
		Tarkan:      {187, 54, 203, 69},
		Aida:        {82, 8, 87, 14},
		Crywolf:     {133, 41, 140, 44},
		Elbeland:    {40, 214, 43, 224},
		LorenMarket: {126, 142, 129, 148},
	}
	MapManager = make(mapManager)
	for _, v := range mapList.DefaultMaps.Map {
		number := v.Number
		file := path.Join(conf.PathCommon, "MapTerrains", v.File)
		rect, ok := rects[number]
		if !ok {
			rect = rects[Lorencia]
		}
		m := _map{regenRect: rect}
		m.init(number, file)
		MapManager[v.Number] = &m
	}
}

type Rect struct {
	left   int
	top    int
	right  int
	bottom int
}

type sender interface {
	SendWeather(int, int)
}

var MapManager mapManager

type mapManager map[int]*_map

func (m mapManager) GetMapAttr(number, x, y int) int {
	return m[number].getAttr(x, y)
}

func (m mapManager) CheckMapAttrStand(number, x, y int) bool {
	return m[number].checkAttrStand(x, y)
}

func (m mapManager) SetMapAttrStand(number, x, y int) {
	m[number].setAttrStand(x, y)
}

func (m mapManager) ClearMapAttrStand(number, x, y int) {
	m[number].clearAttrStand(x, y)
}

func (m mapManager) GetMapRegenPos(number int) (int, int) {
	return m[number].getRegenPos()
}

func (m mapManager) CheckMapNoWall(number, x1, y1, x2, y2 int) bool {
	return m[number].checkNoWall(x1, y1, x2, y2)
}

func (m mapManager) FindMapPath(number, x1, y1, x2, y2 int) (Path, bool) {
	return m[number].findPath(x1, y1, x2, y2)
}

func (m mapManager) GetMapWeather(number int) int {
	return m[number].getWeather()
}

func (m mapManager) ProcessWeather(sender sender) {
	for _, v := range m {
		v.processWeather(sender)
	}
}

type _map struct {
	number    int
	file      string
	width     int
	height    int
	buf       []byte
	regenRect *Rect
	mapItems  []mapItem
	cnt       int
	weather   int
}

func (m *_map) init(number int, file string) {
	m.number = number
	m.file = file
	buf, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("read file failed [file]%s [err]%v", file, err)
	}
	m.width = int(buf[1]) + 1
	m.height = int(buf[2]) + 1
	m.buf = buf[3:]
	m.mapItems = make([]mapItem, conf.Server.GameServerInfo.MaxObjectItemCount)
	m.cnt = 1
}

func (m *_map) valid(x, y int) bool {
	return x >= 0 && x < m.width && y >= 0 && y < m.height
}

// x + y<<8
func (m *_map) pos2index(x, y int) int {
	return x + y*m.width
}

func (m *_map) getAttr(x, y int) int {
	if !m.valid(x, y) {
		return 4
	}
	return int(m.buf[m.pos2index(x, y)])
}

func (m *_map) checkAttrStand(x, y int) bool {
	if !m.valid(x, y) {
		return false
	}
	attr := m.buf[m.pos2index(x, y)]
	if attr&2 != 0 || attr&4 != 0 || attr&8 != 0 {
		return false
	}
	return true
}

func (m *_map) setAttrStand(x, y int) {
	if !m.valid(x, y) {
		return
	}
	m.buf[m.pos2index(x, y)] |= 2
}

func (m *_map) clearAttrStand(x, y int) {
	if !m.valid(x, y) {
		return
	}
	if m.buf[m.pos2index(x, y)]&2 != 0 {
		m.buf[m.pos2index(x, y)] &^= 2
	}
}

func (m *_map) getRegenPos() (int, int) {
	left := m.regenRect.left
	top := m.regenRect.top
	right := m.regenRect.right
	bottom := m.regenRect.bottom
	for cnt := 50; cnt >= 0; cnt-- {
		w := right - left
		h := bottom - top
		var x, y int
		if w <= 0 {
			x = left
		} else {
			x = left + rand.Intn(w)%w
		}
		if h <= 0 {
			y = top
		} else {
			y = top + rand.Intn(h)%h
		}
		attr := m.getAttr(x, y)
		if attr&4 == 0 && attr&8 == 0 {
			return x, y
		}
	}
	log.Printf("cannot find position [file]%s", m.file)
	return left, top
}

func (m *_map) checkNoWall(x1, y1, x2, y2 int) bool {
	if !m.valid(x1, y1) || !m.valid(x2, y2) {
		return false
	}
	w, h, x, y := x2-x1, y2-y1, 1, 256
	if w < 0 {
		w, x = -w, -1
	}
	if h < 0 {
		h, y = -h, -256
	}
	len1, len2, d1, d2 := w, h, x, y
	if w <= h {
		len1, len2, d1, d2 = h, w, y, x
	}
	factor := 0
	index := m.pos2index(x1, y1)
	for i := 0; i <= len1; i++ {
		if m.buf[index]&4 != 0 {
			return false
		}
		index += d2
		factor += len2
		if factor > len1>>1 {
			factor -= len1
			index += d1
		}
	}
	return true
}

func (m *_map) canMoveForward(pos int) bool {
	return !(m.buf[pos] > 1)
}

func (m *_map) findPath(x1, y1, x2, y2 int) (Path, bool) {
	path := _path{
		validator: m,
		width:     m.width,
		height:    m.height,
		hits:      make([]bool, m.width*m.height),
	}
	return path.findPath(x1, y1, x2, y2)
}

func (m *_map) getWeather() int {
	return m.weather
}

func (m *_map) processWeather(sender sender) {
	m.cnt--
	if m.cnt <= 0 {
		m.cnt = rand.Intn(10) + 10
		weather1 := rand.Intn(3)
		weather2 := rand.Intn(10)
		m.weather = weather1<<4 | weather2
		sender.SendWeather(m.number, m.weather)
	}
}
