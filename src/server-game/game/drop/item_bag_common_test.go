package drop

import (
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/item"
)

func commonBagRequest(code, level int) Request {
	return Request{Trigger: TriggerPlayerDropItem, ItemCode: code, ItemLevel: level, MapNumber: 0,
		Player: PlayerContext{Class: int(class.Wizard), Level: 100, MapNumber: 0}}
}

func TestCommonBagDrop(t *testing.T) {
	manager := &dropManager{itemBags: map[itemBagKey]*itemBagEntry{}}
	code := item.Code(14, 11)
	effect := 2
	bag := testItemBag(10000, 0, 0, 2)
	bag.useEffect = &effect
	manager.itemBags[itemBagKey{kind: itemBagCommon, id: code, level: 8}] = testItemBagEntry(bag)

	if got := manager.Drop(commonBagRequest(code, 7)); got.Handled {
		t.Fatalf("wrong item level handled: %#v", got)
	}
	if got := manager.Drop(commonBagRequest(item.Code(14, 12), 8)); got.Handled {
		t.Fatalf("wrong item code handled: %#v", got)
	}
	result := manager.Drop(commonBagRequest(code, 8))
	if !result.Handled || len(result.Rewards) != 2 || result.Effect == nil || *result.Effect != 2 {
		t.Fatalf("CommonBag result = %#v", result)
	}
	for _, reward := range result.Rewards {
		if reward.Item == nil || reward.Placement != PlaceNearSource {
			t.Fatalf("CommonBag reward = %#v", reward)
		}
	}

	bag.allows[0].mapNumber = 1
	if got := manager.Drop(commonBagRequest(code, 8)); got.Handled {
		t.Fatalf("unmatched DropAllow handled: %#v", got)
	}
	bag.allows[0].mapNumber = -1
	bag.bagUseRate = 0
	if got := manager.Drop(commonBagRequest(code, 8)); got.Handled {
		t.Fatalf("failed BagUseRate handled: %#v", got)
	}
	bag.bagUseRate = 10000
	bag.itemRate = 0
	bag.moneyDrop = 5000
	result = manager.Drop(commonBagRequest(code, 8))
	if !result.Handled || len(result.Rewards) != 1 || !isZenReward(result.Rewards[0], 5000) || result.Rewards[0].Placement != PlaceAtPlayer {
		t.Fatalf("CommonBag Zen result = %#v", result)
	}
}
