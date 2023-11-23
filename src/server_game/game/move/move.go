package move

import (
	"encoding/xml"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/maps"
)

func init() {
	MapMoveManager.init()
	GateMoveManager.init()
}

var MapMoveManager mapMoveManager
var GateMoveManager gateMoveManager

type mapMove struct {
	Index      int    `xml:"Index,attr"`
	ServerName string `xml:"ServerName,attr"`
	ClientName string `xml:"ClientName,attr"`
	MinLevel   int    `xml:"MinLevel,attr"`
	MaxLevel   int    `xml:"MaxLevel,attr"`
	ReqMoney   int    `xml:"ReqMoney,attr"`
	GateNumber int    `xml:"GateNumber,attr"`
}

type mapMoveManager struct {
	mapMoveTable map[int]*mapMove
}

func (m *mapMoveManager) init() {
	// WarpSettings was generated 2023-11-22 18:37:12 by https://xml-to-go.github.io/ in Ukraine.
	type WarpSettings struct {
		XMLName xml.Name   `xml:"WarpSettings"`
		Text    string     `xml:",chardata"`
		Warp    []*mapMove `xml:"Warp"`
	}
	var warpSettings WarpSettings
	conf.XML(conf.PathCommon, "Warps/IGC_MoveReq.xml", &warpSettings)
	m.mapMoveTable = make(map[int]*mapMove)
	for _, warp := range warpSettings.Warp {
		m.mapMoveTable[warp.Index] = warp
	}
}

func (m *mapMoveManager) Move(index int, f func(int, int, int)) {
	move, ok := m.mapMoveTable[index]
	if !ok {
		return
	}
	f(move.GateNumber, move.MinLevel, move.ReqMoney)
}

type gateMove struct {
	Index     int    `xml:"Index,attr"`
	Flag      int    `xml:"Flag,attr"`
	MapNumber int    `xml:"MapNumber,attr"`
	StartX    int    `xml:"StartX,attr"`
	StartY    int    `xml:"StartY,attr"`
	EndX      int    `xml:"EndX,attr"`
	EndY      int    `xml:"EndY,attr"`
	Target    int    `xml:"Target,attr"`
	Direction int    `xml:"Direction,attr"`
	MinLevel  int    `xml:"MinLevel,attr"`
	Name      string `xml:"Name,attr"`
}

type gateMoveManager struct {
	gateMoveTable map[int]*gateMove
}

func (m *gateMoveManager) init() {
	// GateSettings was generated 2023-11-22 19:07:36 by https://xml-to-go.github.io/ in Ukraine.
	type GateSettings struct {
		XMLName xml.Name    `xml:"GateSettings"`
		Text    string      `xml:",chardata"`
		Gate    []*gateMove `xml:"Gate"`
	}
	var gateSettings GateSettings
	conf.XML(conf.PathCommon, "Warps/IGC_GateSettings.xml", &gateSettings)
	m.gateMoveTable = make(map[int]*gateMove)
	for _, gate := range gateSettings.Gate {
		m.gateMoveTable[gate.Index] = gate
	}
}

func (m *gateMoveManager) Move(index int, f func(int, int, int, int)) {
	move, ok := m.gateMoveTable[index]
	if !ok {
		return
	}
	if move.Target != 0 {
		move, ok = m.gateMoveTable[move.Target]
		if !ok {
			return
		}
	}
	mapNumber := move.MapNumber
	x1 := move.StartX
	y1 := move.StartY
	x2 := move.EndX
	y2 := move.EndY
	x, y := maps.MapManager.GetMapRandomPos(mapNumber, x1, y1, x2, y2)
	f(mapNumber, x, y, move.Direction)
}
