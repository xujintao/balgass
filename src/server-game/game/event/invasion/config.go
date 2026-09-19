package invasion

import (
	"encoding/xml"
	"fmt"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/object/monster"
)

const RedDragon = 0

type dragonMap struct {
	number int
	spawns []monster.EventSpawn
}

type dragonConfig struct {
	count int
	maps  []dragonMap
}

type configuration struct {
	dragon dragonConfig
}

var config configuration

func init() { config.init() }

func (c *configuration) init() {
	if err := c.load(conf.PathCommon); err != nil {
		panic(fmt.Errorf("invasion: %w", err))
	}
}

func (c *configuration) load(basePath string) error {
	// Pointer attributes distinguish missing required values from valid zeroes.
	type dragonXML struct {
		XMLName  xml.Name `xml:"DragonEvent"`
		Monsters []struct {
			Index    *int `xml:"Index,attr"`
			Distance *int `xml:"Distance,attr"`
			Count    *int `xml:"Count,attr"`
			Maps     []struct {
				Number *int `xml:"Number,attr"`
				Spawns []struct {
					StartX *int `xml:"StartX,attr"`
					StartY *int `xml:"StartY,attr"`
					EndX   *int `xml:"EndX,attr"`
					EndY   *int `xml:"EndY,attr"`
				} `xml:"Spawn"`
			} `xml:"Map"`
		} `xml:"Monster"`
	}
	next := configuration{}
	var dragons dragonXML
	conf.XML(basePath, "Events/IGC_DragonEvent.xml", &dragons)
	if len(dragons.Monsters) != 1 {
		return fmt.Errorf("DragonEvent requires exactly one Monster definition")
	}
	dragon := dragons.Monsters[0]
	if dragon.Index == nil || *dragon.Index != 44 || dragon.Distance == nil || dragon.Count == nil ||
		*dragon.Count <= 0 || *dragon.Count > conf.Server.GameServerInfo.MaxMonsterCount || len(dragon.Maps) == 0 {
		return fmt.Errorf("invalid red dragon class, distance, count or maps")
	}
	next.dragon.count = *dragon.Count
	mapIDs := make(map[int]bool)
	for _, dm := range dragon.Maps {
		if dm.Number == nil || mapIDs[*dm.Number] || len(dm.Spawns) == 0 {
			return fmt.Errorf("missing/duplicate dragon map or empty spawn list")
		}
		mapIDs[*dm.Number] = true
		m := dragonMap{number: *dm.Number}
		areas := make(map[monster.EventSpawn]bool)
		for _, area := range dm.Spawns {
			if area.StartX == nil || area.StartY == nil || area.EndX == nil || area.EndY == nil {
				return fmt.Errorf("dragon map %d: missing spawn coordinates", m.number)
			}
			s := monster.EventSpawn{Class: 44, MapNumber: m.number, StartX: *area.StartX, StartY: *area.StartY,
				EndX: *area.EndX, EndY: *area.EndY, Direction: -1, Distance: *dragon.Distance}
			if err := s.Validate(); err != nil {
				return err
			}
			if areas[s] {
				return fmt.Errorf("duplicate dragon spawn: %+v", s)
			}
			areas[s] = true
			m.spawns = append(m.spawns, s)
		}
		next.dragon.maps = append(next.dragon.maps, m)
	}
	*c = next
	return nil
}
