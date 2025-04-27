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
	MapManager.init()
}

type sender interface {
	SendWeather(int, int)
}

var MapManager mapManager

type zenDropSystem struct {
	Enable                 bool
	MultiplyByMonsterLevel bool
	MultiplyChanceRate     int
}

type mapManager struct {
	maps [MAX_MAP_NUMBER]*_map
	zen  zenDropSystem
}

func (m *mapManager) init() {
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

	for _, v := range mapList.DefaultMaps.Map {
		number := v.Number
		file := path.Join(conf.PathCommon, "MapTerrains", v.File)
		mm := _map{}
		mm.init(number, file)
		m.maps[v.Number] = &mm
	}

	// MapAttribute was generated 2024-01-10 18:37:20 by https://xml-to-go.github.io/ in Ukraine.
	type MapAttribute struct {
		XMLName xml.Name `xml:"MapAttribute"`
		Text    string   `xml:",chardata"`
		Config  struct {
			Text string `xml:",chardata"`
			Map  []struct {
				Text              string  `xml:",chardata"`
				Number            int     `xml:"Number,attr"`
				PvPConfig         bool    `xml:"PvPConfig,attr"`
				ItemDropRateBonus float64 `xml:"ItemDropRateBonus,attr"`
				ExpBonus          float64 `xml:"ExpBonus,attr"`
				MasterExpBonus    float64 `xml:"MasterExpBonus,attr"`
				VipLevel          bool    `xml:"VipLevel,attr"`
				PkLevelIncrease   bool    `xml:"PkLevelIncrease,attr"`
				RegenOnSamePlace  bool    `xml:"RegenOnSamePlace,attr"`
				BlockEntry        bool    `xml:"BlockEntry,attr"`
			} `xml:"Map"`
		} `xml:"Config"`
	}
	var mapAttribute MapAttribute
	conf.XML(conf.PathCommon, "IGC_MapAttribute.xml", &mapAttribute)
	for _, v := range mapAttribute.Config.Map {
		m.maps[v.Number].expBonus = v.ExpBonus
		m.maps[v.Number].masterExpBonus = v.MasterExpBonus
	}

	// ZenDropSystem was generated 2024-01-16 16:31:26 by https://xml-to-go.github.io/ in Ukraine.
	type ZenDropSystem struct {
		XMLName                xml.Name `xml:"ZenDropSystem"`
		Text                   string   `xml:",chardata"`
		Enable                 bool     `xml:"Enable,attr"`
		MultiplyByMonsterLevel bool     `xml:"MultiplyByMonsterLevel,attr"`
		MultiplyChanceRate     int      `xml:"MultiplyChanceRate,attr"`
		Map                    []struct {
			Text          string `xml:",chardata"`
			Number        int    `xml:"Number,attr"`
			MinMoneyCount int    `xml:"MinMoneyCount,attr"`
			MaxMoneyCount int    `xml:"MaxMoneyCount,attr"`
		} `xml:"Map"`
	}
	var zenDropSystem ZenDropSystem
	conf.XML(conf.PathCommon, "IGC_ZenDrop.xml", &zenDropSystem)
	m.zen.Enable = zenDropSystem.Enable
	m.zen.MultiplyByMonsterLevel = zenDropSystem.MultiplyByMonsterLevel
	m.zen.MultiplyChanceRate = zenDropSystem.MultiplyChanceRate
	for _, v := range zenDropSystem.Map {
		m.maps[v.Number].minMoney = v.MinMoneyCount
		m.maps[v.Number].maxMoney = v.MaxMoneyCount
	}
}

func (m *mapManager) GetMapPots(number int) []*Pot {
	return m.maps[number].getPots()
}

// bit0(1): 安全区标志
// bit1(2): 被对象站立
// bit2(4): 障碍物标志
// bit3(8): ?
func (m *mapManager) GetMapAttr(number, x, y int) int {
	return m.maps[number].getAttr(x, y)
}

func (m *mapManager) CheckMapAttrStand(number, x, y int) bool {
	return m.maps[number].checkAttrStand(x, y)
}

func (m *mapManager) SetMapAttrStand(number, x, y int) {
	m.maps[number].setAttrStand(x, y)
}

func (m *mapManager) ClearMapAttrStand(number, x, y int) {
	m.maps[number].clearAttrStand(x, y)
}

func (m *mapManager) GetMapRandomPos(number, x1, y1, x2, y2 int) (int, int) {
	return m.maps[number].getRandomPos(x1, y1, x2, y2)
}

func (m *mapManager) CheckMapNoWall(number, x1, y1, x2, y2 int) bool {
	return m.maps[number].checkNoWall(x1, y1, x2, y2)
}

func (m *mapManager) FindMapPath(number, x1, y1, x2, y2 int) (Path, bool) {
	return m.maps[number].findPath(x1, y1, x2, y2)
}

func (m *mapManager) GetMapWeather(number int) int {
	return m.maps[number].getWeather()
}

func (m *mapManager) ProcessWeather(sender sender) {
	for _, v := range m.maps {
		v.processWeather(sender)
	}
}

func (m *mapManager) GetExpBonus(number int) float64 {
	return m.maps[number].expBonus
}

func (m *mapManager) GetMasterExpBonus(number int) float64 {
	return m.maps[number].masterExpBonus
}

func (m *mapManager) GetZen(number, mLevel int) int {
	if !m.zen.Enable {
		return 0
	}
	money := m.maps[number].GetZen()
	if m.zen.MultiplyByMonsterLevel {
		if rand.Intn(10000) < m.zen.MultiplyChanceRate {
			return money * mLevel
		}
	}
	return money * 2
}

type _map struct {
	number         int
	file           string
	width          int
	height         int
	buf            []byte
	pots           []*Pot
	inventory      []*mapItem
	cnt            int
	weather        int
	expBonus       float64
	masterExpBonus float64
	minMoney       int
	maxMoney       int
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
	for i := range m.buf {
		p := Pot{}
		p.X = i % 256
		p.Y = i >> 8
		if m.buf[i]&4 != 0 {
			m.pots = append(m.pots, &p)
		}
	}
	m.inventory = make([]*mapItem, conf.Server.GameServerInfo.MaxObjectItemCount)
	m.cnt = 1
}

func (m *_map) valid(x, y int) bool {
	return x >= 0 && x < m.width && y >= 0 && y < m.height
}

// x + y<<8
func (m *_map) pos2index(x, y int) int {
	return x + y*m.width
}

func (m *_map) getPots() []*Pot {
	return m.pots
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

func (m *_map) getRandomPos(x1, y1, x2, y2 int) (int, int) {
	w := x2 - x1
	if w <= 0 {
		w = 1
	}
	h := y2 - y1
	if h <= 0 {
		h = 1
	}
	if w == 1 && h == 1 {
		return x1, y1
	}
	for i := 0; i < 100; i++ {
		x := x1 + rand.Intn(w)
		y := y1 + rand.Intn(h)
		attr := m.getAttr(x, y)
		if attr&2 == 0 && attr&4 == 0 && attr&8 == 0 {
			return x, y
		}
	}
	return x1, y1
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

func (m *_map) findPath(x1, y1, x2, y2 int) (Path, bool) {
	path := _path{
		validator: func(x, y int) bool {
			if !m.valid(x, y) {
				return false
			}
			pos := x + y*m.width
			return !(m.buf[pos] > 1)
		},
		hits: make(map[Pot]struct{}),
	}
	// return path.findPath(x1, y1, x2, y2)
	return path.findPathBFS(x1, y1, x2, y2)
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

func (m *_map) GetZen() int {
	sub := m.maxMoney - m.minMoney
	if sub < 0 {
		return 0
	}
	return m.minMoney + rand.Intn(sub+1)
}
