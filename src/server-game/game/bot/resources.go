package bot

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/xujintao/balgass/src/server-game/conf"
)

type Position struct {
	X int
	Y int
}

type terrain struct {
	width  int
	height int
	attrs  []byte
}

func (t *terrain) valid(p Position) bool {
	return p.X >= 0 && p.X < t.width && p.Y >= 0 && p.Y < t.height
}

func (t *terrain) walkable(p Position) bool {
	if !t.valid(p) {
		return false
	}
	return t.attrs[p.X+p.Y*t.width]&(4|8) == 0
}

type spawnArea struct {
	class int
	min   Position
	max   Position
}

func (a spawnArea) contains(p Position) bool {
	return p.X >= a.min.X && p.X <= a.max.X &&
		p.Y >= a.min.Y && p.Y <= a.max.Y
}

func (a spawnArea) nearest(p Position) Position {
	return Position{
		X: clamp(p.X, a.min.X, a.max.X),
		Y: clamp(p.Y, a.min.Y, a.max.Y),
	}
}

func (a spawnArea) center() Position {
	return Position{
		X: (a.min.X + a.max.X) / 2,
		Y: (a.min.Y + a.max.Y) / 2,
	}
}

type resources struct {
	terrains          map[int]*terrain
	spawnAreas        map[int][]spawnArea
	attackableClasses map[int]struct{}
	npcClasses        map[int]struct{}
}

func (r *resources) terrain(number int) *terrain {
	return r.terrains[number]
}

func (r *resources) attackable(class int) bool {
	_, ok := r.attackableClasses[class]
	return ok
}

var (
	defaultResourcesOnce sync.Once
	defaultResources     *resources
	defaultResourcesErr  error
)

func getDefaultResources() (*resources, error) {
	defaultResourcesOnce.Do(func() {
		defaultResources, defaultResourcesErr = loadResources(conf.PathCommon)
	})
	return defaultResources, defaultResourcesErr
}

func loadResources(basePath string) (*resources, error) {
	r := &resources{
		terrains:          make(map[int]*terrain),
		spawnAreas:        make(map[int][]spawnArea),
		attackableClasses: make(map[int]struct{}),
		npcClasses:        make(map[int]struct{}),
	}
	if err := r.loadTerrains(basePath); err != nil {
		return nil, err
	}
	if err := r.loadMonsterSpawn(basePath); err != nil {
		return nil, err
	}
	if err := r.loadShopNPCs(basePath); err != nil {
		return nil, err
	}
	for class := range r.npcClasses {
		delete(r.attackableClasses, class)
	}
	return r, nil
}

func (r *resources) loadTerrains(basePath string) error {
	type mapList struct {
		DefaultMaps struct {
			Map []struct {
				Number *int   `xml:"Number,attr"`
				File   string `xml:"File,attr"`
			} `xml:"Map"`
		} `xml:"DefaultMaps"`
	}
	var cfg mapList
	if err := readXML(filepath.Join(basePath, "IGC_MapList.xml"), &cfg); err != nil {
		return err
	}
	if len(cfg.DefaultMaps.Map) == 0 {
		return fmt.Errorf("bot map list has no default maps")
	}
	for _, m := range cfg.DefaultMaps.Map {
		if m.Number == nil {
			return fmt.Errorf("bot map has no number")
		}
		number := *m.Number
		if number < 0 {
			return fmt.Errorf("bot map has invalid number %d", number)
		}
		if m.File == "" {
			return fmt.Errorf("bot map %d has no terrain file", number)
		}
		if _, ok := r.terrains[number]; ok {
			return fmt.Errorf("bot map %d is duplicated", number)
		}
		t, err := loadTerrain(filepath.Join(basePath, "MapTerrains", m.File))
		if err != nil {
			return fmt.Errorf("bot map %d: %w", number, err)
		}
		r.terrains[number] = t
	}
	return nil
}

func loadTerrain(file string) (*terrain, error) {
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read terrain %s: %w", file, err)
	}
	if len(buf) < 3 {
		return nil, fmt.Errorf("terrain %s header is too short", file)
	}
	width := int(buf[1]) + 1
	height := int(buf[2]) + 1
	size := width * height
	if len(buf[3:]) != size {
		return nil, fmt.Errorf("terrain %s size is %d, want %d", file, len(buf[3:]), size)
	}
	attrs := make([]byte, size)
	copy(attrs, buf[3:])
	return &terrain{width: width, height: height, attrs: attrs}, nil
}

func (r *resources) loadMonsterSpawn(basePath string) error {
	type spawn struct {
		Index  *int `xml:"Index,attr"`
		StartX *int `xml:"StartX,attr"`
		StartY *int `xml:"StartY,attr"`
		EndX   *int `xml:"EndX,attr"`
		EndY   *int `xml:"EndY,attr"`
	}
	type spot struct {
		Type  *int    `xml:"Type,attr"`
		Spawn []spawn `xml:"Spawn"`
	}
	type monsterMap struct {
		Number *int   `xml:"Number,attr"`
		Spot   []spot `xml:"Spot"`
	}
	type monsterSpawn struct {
		Map []monsterMap `xml:"Map"`
	}
	var cfg monsterSpawn
	if err := readXML(filepath.Join(basePath, "Monsters", "IGC_MonsterSpawn.xml"), &cfg); err != nil {
		return err
	}
	for _, m := range cfg.Map {
		if m.Number == nil {
			return fmt.Errorf("bot monster spawn map has no number")
		}
		number := *m.Number
		t := r.terrain(number)
		if t == nil {
			return fmt.Errorf("bot monster spawn references missing map %d", number)
		}
		for _, spot := range m.Spot {
			if spot.Type == nil {
				return fmt.Errorf("bot map %d has monster spot without type", number)
			}
			spotType := *spot.Type
			if spotType < 0 || spotType > 3 {
				return fmt.Errorf("bot map %d has invalid monster spot type %d", number, spotType)
			}
			for _, spawn := range spot.Spawn {
				if spawn.Index == nil || spawn.StartX == nil || spawn.StartY == nil {
					return fmt.Errorf("bot map %d spot type %d has spawn without index or start position", number, spotType)
				}
				class := *spawn.Index
				start := Position{X: *spawn.StartX, Y: *spawn.StartY}
				if !t.valid(start) {
					return fmt.Errorf("bot map %d class %d has invalid spawn start %+v", number, class, start)
				}
				if spotType == 0 {
					r.npcClasses[class] = struct{}{}
					continue
				}
				area := spawnArea{class: class, min: start, max: start}
				switch spotType {
				case 1, 3:
					if spawn.EndX == nil || spawn.EndY == nil {
						return fmt.Errorf("bot map %d class %d area spawn has no end position", number, class)
					}
					area.max = Position{X: *spawn.EndX, Y: *spawn.EndY}
					if !t.valid(area.max) {
						return fmt.Errorf("bot map %d class %d has invalid spawn end %+v", number, class, area.max)
					}
					area.min.X, area.max.X = ordered(area.min.X, area.max.X)
					area.min.Y, area.max.Y = ordered(area.min.Y, area.max.Y)
				case 2:
					area.min.X = clamp(area.min.X-3, 0, t.width-1)
					area.min.Y = clamp(area.min.Y-3, 0, t.height-1)
					area.max.X = clamp(area.max.X+3, 0, t.width-1)
					area.max.Y = clamp(area.max.Y+3, 0, t.height-1)
				}
				if area.min.X > area.max.X || area.min.Y > area.max.Y {
					return fmt.Errorf("bot map %d class %d has invalid spawn area %+v", number, class, area)
				}
				r.spawnAreas[number] = append(r.spawnAreas[number], area)
				r.attackableClasses[class] = struct{}{}
			}
		}
	}
	return nil
}

func (r *resources) loadShopNPCs(basePath string) error {
	type shop struct {
		NPCIndex  *int `xml:"NPCIndex,attr"`
		MapNumber *int `xml:"MapNumber,attr"`
		PosX      *int `xml:"PosX,attr"`
		PosY      *int `xml:"PosY,attr"`
	}
	type shopList struct {
		Shop []shop `xml:"Shop"`
	}
	var cfg shopList
	if err := readXML(filepath.Join(basePath, "IGC_ShopList.xml"), &cfg); err != nil {
		return err
	}
	for _, shop := range cfg.Shop {
		if shop.NPCIndex == nil || shop.MapNumber == nil || shop.PosX == nil || shop.PosY == nil {
			return fmt.Errorf("bot shop requires NPCIndex, MapNumber, PosX and PosY")
		}
		class := *shop.NPCIndex
		number := *shop.MapNumber
		t := r.terrain(number)
		if t == nil {
			return fmt.Errorf("bot shop class %d references missing map %d", class, number)
		}
		p := Position{X: *shop.PosX, Y: *shop.PosY}
		if !t.valid(p) {
			return fmt.Errorf("bot shop class %d has invalid position %+v", class, p)
		}
		r.npcClasses[class] = struct{}{}
	}
	return nil
}

func readXML(file string, v any) error {
	buf, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read XML %s: %w", file, err)
	}
	if err := xml.Unmarshal(buf, v); err != nil {
		return fmt.Errorf("unmarshal XML %s: %w", file, err)
	}
	return nil
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func ordered(a, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}
