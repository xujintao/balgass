package drop

import (
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/item"
)

func testRequest(monsterClass int) Request {
	return Request{Trigger: TriggerMonsterDeath, MonsterClass: monsterClass, MapNumber: 0, Player: PlayerContext{Class: int(class.Wizard), Level: 100, MapNumber: 0}}
}

func testItemBag(itemRate, setItemRate, money, count int) *itemBag {
	return &itemBag{itemRate: itemRate, setItemRate: setItemRate, moneyDrop: money, bagUseRate: 10000, allows: []itemBagAllow{{
		classes: [class.MaxClass]int{class.Wizard: 1}, minLevel: 1, maxLevel: 800, mapNumber: -1,
		drops: []itemBagDrop{{rate: 10000, count: count, items: []itemBagItem{{section: 0, index: 17, minLevel: 0, maxLevel: 0, skill: 1, luck: 1, option: 1}}}},
	}}}
}
func testItemBagEntry(bag *itemBag) *itemBagEntry {
	return &itemBagEntry{bag: bag}
}

func TestParseItemBagBindings(t *testing.T) {
	script := `
AddItemBag(BAG_COMMON, MakeItemID(14,11), 8, 'Item_(14,11,8)_Kundun_Box+1')
AddItemBag(BAG_MONSTER, 0, 44, 'Monster_(44)_Dragon_Red')
AddItemBag(BAG_EVENT, 0, 0, 'Buff_AkeronTower_Drop')
AddItemBag(BAG_EVENT, 26, 0, 'Monster_(275)_Kundun')
AddItemBag(BAG_EVENT, 46, 0, 'Monster_(673)_Lord_Silvester')
AddItemBag(BAG_EVENT, 150, 0, 'Event_Monster_(561)_Medusa')
AddItemBag(BAG_EVENT, 151, 0, 'Event_Monster_(560)_Sapi_Queen')
`
	bindings, err := parseItemBagBindings(script)
	if err != nil {
		t.Fatal(err)
	}
	for key, name := range map[itemBagKey]string{
		{kind: itemBagMonster, id: 44}:                         "Monster_(44)_Dragon_Red",
		{kind: itemBagCommon, id: item.Code(14, 11), level: 8}: "Item_(14,11,8)_Kundun_Box+1",
		{kind: itemBagEvent, id: 1}:                            "Buff_AkeronTower_Drop",
		{kind: itemBagEvent, id: 27}:                           "Monster_(275)_Kundun",
		{kind: itemBagEvent, id: 47}:                           "Monster_(673)_Lord_Silvester",
		{kind: itemBagEvent, id: 151}:                          "Event_Monster_(561)_Medusa",
		{kind: itemBagEvent, id: 152}:                          "Event_Monster_(560)_Sapi_Queen",
	} {
		if got, ok := bindings[key]; !ok || got != name {
			t.Fatalf("binding %v = %#v, want %q", key, got, name)
		}
	}
	if len(bindings) != 7 {
		t.Fatalf("bindings = %d, want 7", len(bindings))
	}
	if _, err := parseItemBagBindings("AddItemBag(BAG_MONSTER, 0, 44, 'a')\nAddItemBag(BAG_MONSTER, 0, 44, 'b')"); err == nil {
		t.Fatal("duplicate registration was accepted")
	}
	if _, err := parseItemBagBindings("AddItemBag(BAG_EVENT, 0, 0, 'a')\nAddItemBag(BAG_EVENT, 0, 0, 'b')"); err == nil {
		t.Fatal("duplicate EventBag registration was accepted")
	}
	if _, err := parseItemBagBindings("AddItemBag(BAG_COMMON, ItemCode(14,11), 8, 'x')"); err == nil {
		t.Fatal("non-static CommonBag expression was accepted")
	}
}

func TestItemBagStartupLoadsOnlyMonsterBags(t *testing.T) {
	monsterCount := 0
	for key, entry := range DropManager.itemBags {
		if key.kind != itemBagMonster {
			continue
		}
		if entry.bag == nil {
			t.Fatalf("MonsterBag %d was not loaded", key.id)
		}
		monsterCount++
	}
	if monsterCount != 20 {
		t.Fatalf("MonsterBag registration count = %d, want 20", monsterCount)
	}
	for _, key := range []itemBagKey{{kind: itemBagEvent, id: 1}, {kind: itemBagEvent, id: 27}, {kind: itemBagEvent, id: 47}, {kind: itemBagEvent, id: 151}, {kind: itemBagEvent, id: 152}} {
		if entry := DropManager.itemBags[key]; entry == nil || entry.bag != nil {
			t.Fatalf("EventBag %d should be registered but not activated by drop", key.id)
		}
	}
	for _, key := range []itemBagKey{{kind: itemBagCommon, id: item.Code(14, 11)}, {kind: itemBagEvent, id: 7}} {
		entry := DropManager.itemBags[key]
		if entry == nil || entry.bag != nil {
			t.Fatalf("inactive ItemBag %v/%d/%d was not registered lazily", key.kind, key.id, key.level)
		}
	}
}

func TestItemBagFiltersWeightCountAndRewards(t *testing.T) {
	bag := testItemBag(10000, 0, 0, 2)
	if got := len(bag.matchingAllows(testRequest(44))); got != 1 {
		t.Fatalf("matching allows = %d, want 1", got)
	}
	request := testRequest(44)
	request.Player.Class = int(class.Knight)
	if got := len(bag.matchingAllows(request)); got != 0 {
		t.Fatalf("non-matching class allows = %d, want 0", got)
	}
	request = testRequest(44)
	request.Player.Level = 0
	if got := len(bag.matchingAllows(request)); got != 0 {
		t.Fatalf("below-level allows = %d, want 0", got)
	}
	bag.allows[0].mapNumber = 1
	if got := len(bag.matchingAllows(testRequest(44))); got != 0 {
		t.Fatalf("non-matching map allows = %d, want 0", got)
	}

	drops := []itemBagDrop{{rate: 1000}, {rate: 9000, count: 2}}
	if got := selectItemBagDrop(drops, 1000); got != &drops[1] {
		t.Fatalf("weighted Drop = %#v, want second Drop", got)
	}
	if result, used := (&dropManager{}).dropItemBag(testItemBag(0, 0, 12345, 1), false, testRequest(44)); !used || len(result.Rewards) != 1 || result.Rewards[0].Zen != 12345 {
		t.Fatalf("ItemRate Zen result = %#v", result)
	}
	if result, used := (&dropManager{}).dropItemBag(testItemBag(10000, 0, 0, 2), false, testRequest(44)); !used || len(result.Rewards) != 2 {
		t.Fatalf("Count result = %#v", result)
	}
	setItem := (&dropManager{}).makeRandomSetItem()
	if setItem == nil || setItem.Set <= 0 || !setItem.Skill || setItem.Durability != setItem.MaxDurability {
		t.Fatalf("random set item = %#v", setItem)
	}
}

func TestMonsterAndEventBagDispatch(t *testing.T) {
	m := dropManager{itemBags: map[itemBagKey]*itemBagEntry{
		{kind: itemBagEvent, id: 1}:  testItemBagEntry(testItemBag(0, 0, 1, 1)),
		{kind: itemBagEvent, id: 27}: testItemBagEntry(testItemBag(0, 0, 27, 1)),
		{kind: itemBagEvent, id: 47}: testItemBagEntry(testItemBag(0, 0, 47, 1)),
	}}
	for _, id := range []int{1, 27, 47} {
		req := testRequest(44)
		req.EventID = id
		result := m.Drop(req)
		if !result.Handled || len(result.Rewards) != 1 || result.Rewards[0].Zen != id {
			t.Fatalf("EventBag %d was not dispatched: %#v", id, result)
		}
	}

	fallback := DropManager
	fallback.itemBags = map[itemBagKey]*itemBagEntry{{kind: itemBagMonster, id: 295}: testItemBagEntry(&itemBag{bagUseRate: 0})}
	result := fallback.Drop(testRequest(295))
	if !result.Handled {
		t.Fatal("unhandled Erohim MonsterBag did not enter the generic drop path")
	}
	for _, reward := range result.Rewards {
		if reward.Nearby {
			t.Fatal("generic fallback reward was marked as an ItemBag reward")
		}
	}
}

func TestExplicitMonsterGroupEventBag(t *testing.T) {
	id := 151
	m := dropManager{itemBags: map[itemBagKey]*itemBagEntry{
		{kind: itemBagMonster, id: 44}: testItemBagEntry(testItemBag(0, 0, 44, 1)),
		{kind: itemBagEvent, id: id}:   testItemBagEntry(testItemBag(0, 0, 151, 1)),
	}}
	req := testRequest(44)
	req.EventID = id
	result := m.Drop(req)
	if !result.Handled || len(result.Rewards) != 1 || result.Rewards[0].Zen != 151 {
		t.Fatal("explicit bag did not override monster bag", result)
	}
	req.EventID = 0
	result = m.Drop(req)
	if len(result.Rewards) != 1 || result.Rewards[0].Zen != 44 {
		t.Fatal("normal monster bag changed")
	}
	bindings, err := parseItemBagBindings("AddItemBag(BAG_EVENT, 150, 0, 'Event_Monster_(561)_Medusa')")
	if err != nil || len(bindings) != 1 {
		t.Fatal("configured event bag not parsed", err)
	}
}
func TestLoadEventBagFromRegisteredBinding(t *testing.T) {
	key := itemBagKey{kind: itemBagEvent, id: 27}
	entry := DropManager.itemBags[key]
	if entry == nil || entry.name == "" {
		t.Fatalf("EventBag registration = %#v, want registered entry", entry)
	}
	m := dropManager{itemBags: map[itemBagKey]*itemBagEntry{
		key: {name: entry.name},
	}}
	if err := m.LoadEventBag(key.id); err != nil {
		t.Fatal(err)
	}
	if m.itemBags[key].bag == nil {
		t.Fatal("registered EventBag was not loaded")
	}
}
func TestMedusaSetItemKeepsConfiguredBaseItem(t *testing.T) {
	key := itemBagKey{kind: itemBagEvent, id: 151}
	entry := DropManager.itemBags[key]
	if entry == nil || entry.name == "" {
		t.Fatal("Medusa EventBag registration missing")
	}
	m := dropManager{itemBags: map[itemBagKey]*itemBagEntry{key: {name: entry.name}}}
	if err := m.LoadEventBag(151); err != nil {
		t.Fatal(err)
	}
	if err := m.LoadEventBag(151); err != nil {
		t.Fatal("cached EventBag reload:", err)
	}
	for _, allow := range m.itemBags[key].bag.allows {
		for _, group := range allow.drops {
			for _, configured := range group.items {
				if !configured.setItem {
					continue
				}
				reward := m.makeBagItem(configured)
				if reward.Section != configured.section || reward.Index != configured.index || reward.Set <= 0 {
					t.Fatalf("SetItem=1 changed base item or omitted set: %#v", reward)
				}
				return
			}
		}
	}
	t.Fatal("Medusa EventBag has no SetItem=1 reward")
}
