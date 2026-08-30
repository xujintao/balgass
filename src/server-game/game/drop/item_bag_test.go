package drop

import (
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/class"
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

func TestParseItemBagBindings(t *testing.T) {
	script := `
AddItemBag(BAG_COMMON, MakeItemID(14,11), 8, 'Item_(14,11,8)_Kundun_Box+1')
AddItemBag(BAG_MONSTER, 0, 44, 'Monster_(44)_Dragon_Red')
AddItemBag(BAG_EVENT, 26, 0, 'Monster_(275)_Kundun')
`
	bindings, err := parseItemBagBindings(script)
	if err != nil {
		t.Fatal(err)
	}
	for key, name := range map[itemBagKey]string{
		{kind: itemBagMonster, id: 44}: "Monster_(44)_Dragon_Red",
		{kind: itemBagEvent, id: 26}:   "Monster_(275)_Kundun",
	} {
		if got, ok := bindings[key]; !ok || got != name {
			t.Fatalf("binding %v = %#v, want %q", key, got, name)
		}
	}
	if len(bindings) != 2 {
		t.Fatalf("active bindings = %d, want 2", len(bindings))
	}
	if _, err := parseItemBagBindings("AddItemBag(BAG_MONSTER, 0, 44, 'a')\nAddItemBag(BAG_MONSTER, 0, 44, 'b')"); err == nil {
		t.Fatal("duplicate registration was accepted")
	}
	if bindings, err := parseItemBagBindings("AddItemBag(BAG_COMMON, ItemCode(14,11), 8, 'x')"); err != nil || len(bindings) != 0 {
		t.Fatalf("inactive CommonBag should be ignored, bindings=%#v err=%v", bindings, err)
	}
}

func TestItemBagStartupLoadsMonsterAndSpecialEvents(t *testing.T) {
	monsterCount := 0
	for key := range DropManager.itemBags {
		if key.kind != itemBagMonster {
			continue
		}
		monsterCount++
	}
	if monsterCount != 20 {
		t.Fatalf("MonsterBag registration count = %d, want 20", monsterCount)
	}
	for _, key := range []itemBagKey{{kind: itemBagEvent, id: 26}, {kind: itemBagEvent, id: 46}} {
		if DropManager.itemBags[key] == nil {
			t.Fatalf("special EventBag %d was not loaded", key.id)
		}
	}
	if DropManager.itemBags[itemBagKey{kind: itemBagEvent, id: 6}] != nil {
		t.Fatal("inactive Common/Event XML was loaded")
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
	m := dropManager{itemBags: map[itemBagKey]*itemBag{
		{kind: itemBagEvent, id: 26}: testItemBag(0, 0, 275, 1),
		{kind: itemBagEvent, id: 46}: testItemBag(0, 0, 673, 1),
	}}
	for _, monsterClass := range []int{275, 673} {
		result := m.Drop(testRequest(monsterClass))
		if !result.Handled || len(result.Rewards) != 1 || result.Rewards[0].Zen != monsterClass {
			t.Fatalf("boss %d did not consume death with its EventBag: %#v", monsterClass, result)
		}
	}

	fallback := DropManager
	fallback.itemBags = map[itemBagKey]*itemBag{{kind: itemBagMonster, id: 295}: {bagUseRate: 0}}
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
