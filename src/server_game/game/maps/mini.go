package maps

import (
	"encoding/xml"

	"github.com/xujintao/balgass/src/server_game/conf"
)

func init() {
	MiniManager.init()
}

var MiniManager miniManager

type mini struct {
	ID          int    `xml:"Id,attr"`
	MapNumber   int    `xml:"MapNumber,attr"`
	NPCIndex    int    `xml:"NPCIndex,attr"`
	X           int    `xml:"X,attr"`
	Y           int    `xml:"Y,attr"`
	DisplayType int    `xml:"DisplayType,attr"`
	SyncType    int    `xml:"SyncType,attr"`
	Name        string `xml:"Name,attr"`
	Annotation  string `xml:"annotation,attr"`
}

type miniManager struct {
	npcTable      map[int][]*mini
	entranceTable map[int][]*mini
}

func (m *miniManager) init() {
	// MiniMapObjects was generated 2023-11-23 22:33:51 by https://xml-to-go.github.io/ in Ukraine.
	type MiniMapObjects struct {
		XMLName xml.Name `xml:"MiniMapObjects"`
		Text    string   `xml:",chardata"`
		TypeOne struct {
			Text string  `xml:",chardata"`
			Tag  []*mini `xml:"Tag"`
		} `xml:"TypeOne"`
		TypeTwo struct {
			Text string  `xml:",chardata"`
			Tag  []*mini `xml:"Tag"`
		} `xml:"TypeTwo"`
	}
	var miniMapObjects MiniMapObjects
	conf.XML(conf.PathCommon, "IGC_MiniMap.xml", &miniMapObjects)
	m.npcTable = make(map[int][]*mini)
	m.entranceTable = make(map[int][]*mini)
	for _, tag := range miniMapObjects.TypeOne.Tag {
		m.npcTable[tag.MapNumber] = append(m.npcTable[tag.MapNumber], tag)
	}
	for _, tag := range miniMapObjects.TypeTwo.Tag {
		m.entranceTable[tag.MapNumber] = append(m.entranceTable[tag.MapNumber], tag)
	}
}

func (m *miniManager) ForEachMapNpc(mapNumber int, f func(int, int, int, int, string)) {
	npcs, ok := m.npcTable[mapNumber]
	if !ok {
		return
	}
	for _, npc := range npcs {
		f(npc.ID, npc.DisplayType, npc.X, npc.Y, npc.Name)
	}
}

func (m *miniManager) ForEachMapEntrance(mapNumber int, f func(int, int, int, int, string)) {
	entrances, ok := m.entranceTable[mapNumber]
	if !ok {
		return
	}
	for _, entrance := range entrances {
		f(entrance.ID, entrance.DisplayType, entrance.X, entrance.Y, entrance.Name)
	}
}
