package drop

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/item"
)

const (
	itemBagCommon itemBagKind = iota
	itemBagMonster
	itemBagEvent
)

type itemBagKind uint8

type itemBagKey struct {
	kind  itemBagKind
	id    int
	level int
}

type itemBag struct {
	name        string
	itemRate    int
	setItemRate int
	moneyDrop   int
	bagUseRate  int
	allows      []itemBagAllow
}

type itemBagAllow struct {
	classes   [class.MaxClass]int
	minLevel  int
	maxLevel  int
	mapNumber int
	drops     []itemBagDrop
}

type itemBagDrop struct {
	rate  int
	count int
	items []itemBagItem
}

type itemBagItem struct {
	section    int
	index      int
	minLevel   int
	maxLevel   int
	durability *int
	skill      int
	luck       int
	option     int
	excellent  int
}

func (m *dropManager) initItemBags() {
	scriptPath := filepath.Join(conf.PathCommon, "Scripts/ItemBags/ItemBagScript.lua")
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		panic(fmt.Errorf("drop: read ItemBagScript.lua: %w", err))
	}
	bindings, err := parseItemBagBindings(string(script))
	if err != nil {
		panic(fmt.Errorf("drop: parse ItemBagScript.lua: %w", err))
	}
	m.itemBags = make(map[itemBagKey]*itemBag, len(bindings))
	for key, name := range bindings {
		path := filepath.Join(conf.PathCommon, "ItemBags", name+".xml")
		if _, err := os.Stat(path); err != nil {
			panic(fmt.Errorf("drop: ItemBag %v/%d/%d file %q: %w", key.kind, key.id, key.level, name, err))
		}
		m.itemBags[key] = loadItemBag(name)
	}
}

func parseItemBagBindings(script string) (map[itemBagKey]string, error) {
	// The loader intentionally recognizes only static registration calls. It
	// never evaluates Lua, so script code cannot affect server startup.
	re := regexp.MustCompile(`(?m)^\s*AddItemBag\(\s*(BAG_COMMON|BAG_MONSTER|BAG_EVENT)\s*,\s*(.+?)\s*,\s*(\d+)\s*,\s*'([^']+)'\s*\)`)
	number := regexp.MustCompile(`^\d+$`)
	matches := re.FindAllStringSubmatch(script, -1)
	bindings := make(map[itemBagKey]string)
	for _, match := range matches {
		first := strings.TrimSpace(match[2])
		second, _ := strconv.Atoi(match[3])
		name := match[4]
		if filepath.Base(name) != name {
			return nil, fmt.Errorf("invalid ItemBag file name %q", name)
		}
		var key itemBagKey
		switch match[1] {
		case "BAG_COMMON":
			continue
		case "BAG_MONSTER":
			if first != "0" || !number.MatchString(first) || second < 0 {
				return nil, fmt.Errorf("invalid MonsterBag key %q, %d", first, second)
			}
			if class.MonsterTable[second] == nil {
				return nil, fmt.Errorf("invalid MonsterBag monster class %d", second)
			}
			key = itemBagKey{kind: itemBagMonster, id: second}
		case "BAG_EVENT":
			if !number.MatchString(first) || second != 0 {
				return nil, fmt.Errorf("invalid EventBag key %q, %d", first, second)
			}
			eventID, _ := strconv.Atoi(first)
			if eventID != 26 && eventID != 46 {
				continue
			}
			key = itemBagKey{kind: itemBagEvent, id: eventID}
		}
		if _, exists := bindings[key]; exists {
			return nil, fmt.Errorf("duplicate ItemBag registration %v/%d/%d", key.kind, key.id, key.level)
		}
		bindings[key] = name
	}
	return bindings, nil
}

func monsterEventItemBagKey(monsterClass int) (itemBagKey, bool) {
	switch monsterClass {
	case 275:
		return itemBagKey{kind: itemBagEvent, id: 26}, true
	case 673:
		return itemBagKey{kind: itemBagEvent, id: 46}, true
	default:
		return itemBagKey{}, false
	}
}

func loadItemBag(name string) *itemBag {
	type itemXML struct {
		Cat                int  `xml:"Cat,attr"`
		Index              int  `xml:"Index,attr"`
		ItemMinLevel       int  `xml:"ItemMinLevel,attr"`
		ItemMaxLevel       int  `xml:"ItemMaxLevel,attr"`
		Durability         *int `xml:"Durability,attr"`
		Skill              int  `xml:"Skill,attr"`
		Luck               int  `xml:"Luck,attr"`
		Option             int  `xml:"Option,attr"`
		Exc                int  `xml:"Exc,attr"`
		SetItem            int  `xml:"SetItem,attr"`
		SocketCount        int  `xml:"SocketCount,attr"`
		ElementalItem      int  `xml:"ElementalItem,attr"`
		MuunEvolutionCat   int  `xml:"MuunEvolutionItemCat,attr"`
		MuunEvolutionIndex int  `xml:"MuunEvolutionItemIndex,attr"`
		Duration           int  `xml:"Duration,attr"`
	}
	type config struct {
		BagConfig struct {
			Name         string `xml:"Name,attr"`
			ItemRate     int    `xml:"ItemRate,attr"`
			SetItemRate  int    `xml:"SetItemRate,attr"`
			MoneyDrop    int    `xml:"MoneyDrop,attr"`
			BagUseEffect int    `xml:"BagUseEffect,attr"`
			BagUseRate   int    `xml:"BagUseRate,attr"`
		} `xml:"BagConfig"`
		SummonBook struct {
			Enable int `xml:"Enable,attr"`
		} `xml:"SummonBook"`
		AddCoin struct {
			Enable int `xml:"Enable,attr"`
		} `xml:"AddCoin"`
		Allows []struct {
			Wizard         int `xml:"Wizard,attr"`
			Knight         int `xml:"Knight,attr"`
			FairyElf       int `xml:"FairyElf,attr"`
			MagicGladiator int `xml:"MagicGladiator,attr"`
			DarkLord       int `xml:"DarkLord,attr"`
			Summoner       int `xml:"Summoner,attr"`
			RageFighter    int `xml:"RageFighter,attr"`
			PlayerMin      int `xml:"PlayerMinLevel,attr"`
			PlayerMax      int `xml:"PlayerMaxLevel,attr"`
			MapNumber      int `xml:"MapNumber,attr"`
			Drops          []struct {
				Rate  int       `xml:"Rate,attr"`
				Count int       `xml:"Count,attr"`
				Items []itemXML `xml:"Item"`
			} `xml:"Drop"`
		} `xml:"DropAllow"`
	}

	var cfg config
	conf.XML(conf.PathCommon, "ItemBags/"+name+".xml", &cfg)
	if cfg.BagConfig.Name == "" || cfg.BagConfig.ItemRate < 0 || cfg.BagConfig.ItemRate > 10000 ||
		cfg.BagConfig.SetItemRate < 0 || cfg.BagConfig.SetItemRate > 10000 ||
		cfg.BagConfig.MoneyDrop < 0 || cfg.BagConfig.BagUseRate < 0 || cfg.BagConfig.BagUseRate > 10000 ||
		(cfg.BagConfig.ItemRate < 10000 && cfg.BagConfig.MoneyDrop == 0) {
		panic(fmt.Errorf("drop: invalid BagConfig in %s", name))
	}
	if cfg.BagConfig.BagUseEffect != -1 || cfg.SummonBook.Enable != 0 || cfg.AddCoin.Enable != 0 {
		panic(fmt.Errorf("drop: unsupported BagUseEffect/SummonBook/AddCoin enabled in %s", name))
	}
	bag := &itemBag{name: cfg.BagConfig.Name, itemRate: cfg.BagConfig.ItemRate, setItemRate: cfg.BagConfig.SetItemRate, moneyDrop: cfg.BagConfig.MoneyDrop, bagUseRate: cfg.BagConfig.BagUseRate}
	validateItem := func(allowIndex, dropIndex, itemIndex int, configured itemXML) {
		if configured.Cat < 0 || configured.Index < 0 || configured.ItemMinLevel < 0 || configured.ItemMaxLevel < configured.ItemMinLevel || configured.ItemMaxLevel > 15 {
			panic(fmt.Errorf("drop: invalid item %d/%d/%d in %s", allowIndex, dropIndex, itemIndex, name))
		}
		if _, err := item.ItemTable.GetItemBase(configured.Cat, configured.Index); err != nil {
			panic(fmt.Errorf("drop: invalid item %d/%d/%d in %s: %w", allowIndex, dropIndex, itemIndex, name, err))
		}
		if configured.Durability != nil && (*configured.Durability < 0 || *configured.Durability > 255) ||
			(configured.Skill != -1 && configured.Skill != 0 && configured.Skill != 1) ||
			(configured.Luck != -1 && configured.Luck != 0 && configured.Luck != 1) ||
			(configured.Option < -1 || configured.Option > 7) ||
			(configured.Exc < -1 || configured.Exc > 63) || configured.SetItem != 0 ||
			configured.SocketCount != 0 || configured.ElementalItem != 0 || configured.MuunEvolutionCat != 0 || configured.MuunEvolutionIndex != 0 || configured.Duration != 0 {
			panic(fmt.Errorf("drop: unsupported item fields at %d/%d/%d in %s", allowIndex, dropIndex, itemIndex, name))
		}
	}
	for allowIndex, allow := range cfg.Allows {
		if allow.PlayerMin < 0 || allow.PlayerMax < allow.PlayerMin || allow.MapNumber < -1 {
			panic(fmt.Errorf("drop: invalid DropAllow %d in %s", allowIndex, name))
		}
		classes := [class.MaxClass]int{}
		classes[class.Wizard] = allow.Wizard
		classes[class.Knight] = allow.Knight
		classes[class.Elf] = allow.FairyElf
		classes[class.Magumsa] = allow.MagicGladiator
		classes[class.DarkLord] = allow.DarkLord
		classes[class.Summoner] = allow.Summoner
		classes[class.RageFighter] = allow.RageFighter
		for _, changeUp := range classes {
			if changeUp < 0 || changeUp > 3 {
				panic(fmt.Errorf("drop: invalid class filter in %s", name))
			}
		}
		runtimeAllow := itemBagAllow{classes: classes, minLevel: allow.PlayerMin, maxLevel: allow.PlayerMax, mapNumber: allow.MapNumber}
		totalRate := 0
		for dropIndex, configuredDrop := range allow.Drops {
			if configuredDrop.Rate < 0 || configuredDrop.Rate > 10000 || configuredDrop.Count <= 0 || len(configuredDrop.Items) == 0 {
				panic(fmt.Errorf("drop: invalid DropAllow %d Drop %d in %s", allowIndex, dropIndex, name))
			}
			runtimeDrop := itemBagDrop{rate: configuredDrop.Rate, count: configuredDrop.Count}
			for itemIndex, configuredItem := range configuredDrop.Items {
				validateItem(allowIndex, dropIndex, itemIndex, configuredItem)
				runtimeDrop.items = append(runtimeDrop.items, itemBagItem{
					section: configuredItem.Cat, index: configuredItem.Index, minLevel: configuredItem.ItemMinLevel, maxLevel: configuredItem.ItemMaxLevel,
					durability: configuredItem.Durability, skill: configuredItem.Skill, luck: configuredItem.Luck, option: configuredItem.Option, excellent: configuredItem.Exc,
				})
			}
			totalRate += configuredDrop.Rate
			runtimeAllow.drops = append(runtimeAllow.drops, runtimeDrop)
		}
		if totalRate <= 0 {
			panic(fmt.Errorf("drop: DropAllow %d in %s has no weighted Drop", allowIndex, name))
		}
		bag.allows = append(bag.allows, runtimeAllow)
	}
	if len(bag.allows) == 0 {
		panic(fmt.Errorf("drop: ItemBag %s has no DropAllow", name))
	}
	return bag
}

func (m *dropManager) dropItemBag(bag *itemBag, consume bool, request Request) (Result, bool) {
	if rand.Intn(10000) >= bag.bagUseRate {
		return Result{Handled: consume}, consume
	}
	result := Result{Handled: true}
	if rand.Intn(10000) >= bag.itemRate {
		if bag.moneyDrop > 0 {
			result.Rewards = append(result.Rewards, Reward{Zen: bag.moneyDrop, Nearby: true})
		}
		return result, true
	}
	matching := bag.matchingAllows(request)
	if len(matching) == 0 {
		return result, true
	}
	allow := matching[rand.Intn(len(matching))]
	totalRate := 0
	for _, configuredDrop := range allow.drops {
		totalRate += configuredDrop.rate
	}
	if totalRate <= 0 {
		return result, true
	}
	selected := selectItemBagDrop(allow.drops, rand.Intn(totalRate))
	if selected == nil {
		return result, true
	}
	for i := 0; i < selected.count; i++ {
		var reward *item.Item
		if bag.setItemRate > 0 && rand.Intn(10000) < bag.setItemRate {
			reward = m.makeRandomSetItem()
		} else {
			reward = m.makeBagItem(selected.items[rand.Intn(len(selected.items))])
		}
		if reward != nil {
			result.Rewards = append(result.Rewards, Reward{Item: reward, Nearby: true})
		}
	}
	return result, true
}

func (bag *itemBag) matchingAllows(request Request) []*itemBagAllow {
	var matching []*itemBagAllow
	for i := range bag.allows {
		allow := &bag.allows[i]
		if request.Player.Class < 0 || request.Player.Class >= class.MaxClass || allow.classes[request.Player.Class] < request.Player.ChangeUp+1 ||
			request.Player.Level+request.Player.MasterLevel < allow.minLevel || request.Player.Level+request.Player.MasterLevel > allow.maxLevel ||
			(allow.mapNumber != -1 && allow.mapNumber != request.MapNumber) {
			continue
		}
		matching = append(matching, allow)
	}
	return matching
}

func selectItemBagDrop(drops []itemBagDrop, number int) *itemBagDrop {
	for i := range drops {
		if number < drops[i].rate {
			return &drops[i]
		}
		number -= drops[i].rate
	}
	return nil
}

func (m *dropManager) makeBagItem(configured itemBagItem) *item.Item {
	it := item.NewItem(configured.section, configured.index)
	it.Level = configured.minLevel + rand.Intn(configured.maxLevel-configured.minLevel+1)
	it.Skill = configured.skill == 1 || configured.skill == -1 && rand.Intn(2) == 0
	it.Lucky = configured.luck == 1 || configured.luck == -1 && rand.Intn(2) == 0
	option := configured.option
	if option == -1 {
		option = rand.Intn(8)
	}
	it.Addition = option * 4
	if configured.excellent == -1 {
		it.DecodeExcellent(item.ExcellentDropManager.DropExcellent(it.KindA, it.KindB))
	} else if configured.excellent > 0 {
		it.DecodeExcellent(configured.excellent)
	}
	it.Calc()
	if configured.durability != nil {
		it.Durability = *configured.durability
	} else {
		it.Durability = it.MaxDurability
	}
	return it
}

func (m *dropManager) makeRandomSetItem() *item.Item {
	section, index, setIndex, ok := item.SetManager.RandomItem()
	if !ok {
		return nil
	}
	it := item.NewItem(section, index)
	it.Set = setIndex
	it.Skill = true
	it.Lucky = rand.Intn(100) < 4
	optionClass := rand.Intn(3)
	optionRoll := rand.Intn(100)
	switch optionClass {
	case 0:
		if optionRoll < 4 {
			it.Addition = 12
		}
	case 1:
		if optionRoll < 8 {
			it.Addition = 8
		}
	case 2:
		if optionRoll < 12 {
			it.Addition = 4
		}
	}
	it.Calc()
	it.Durability = it.MaxDurability
	return it
}
