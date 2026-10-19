package object

import (
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

func TestDropCommonBagConsumesItemAndPlaysEffect(t *testing.T) {
	obj, actor := newBasicAttackTestObject(32, ObjectTypePlayer)
	withTestObjectManager(t, obj)
	obj.X, obj.Y = findBasicAttackMapPosition(t, false)
	actor.inventory.Size = 237
	actor.inventory.Items = make([]*item.Item, actor.inventory.Size)
	actor.inventory.Flags = make([]bool, actor.inventory.Size)
	bag := item.NewItem(14, 11)
	bag.Level = 2 // Firecracker: one reward and effect 0.
	actor.inventory.AddItem(12, bag)
	t.Cleanup(func() {
		var indexes []int
		maps.MapManager.MapEachItem(obj.MapNumber, func(_ *item.Item, index, _, _ int, _ bool) {
			indexes = append(indexes, index)
		})
		for _, index := range indexes {
			maps.MapManager.RemoveItem(obj.MapNumber, index)
		}
	})

	obj.DropItem(&model.MsgDropItem{Position: 12, X: obj.X + 1, Y: obj.Y})
	if actor.inventory.Items[12] != nil {
		t.Fatal("CommonBag was not consumed")
	}
	var dropReply *model.MsgDropItemReply
	var effect *model.MsgServerCMDReply
	for _, msg := range actor.messages {
		switch value := msg.(type) {
		case *model.MsgDropItemReply:
			dropReply = value
		case *model.MsgServerCMDReply:
			effect = value
		}
	}
	if dropReply == nil || dropReply.Result != 1 || effect == nil || effect.Type != 0 || effect.X != obj.X || effect.Y != obj.Y {
		t.Fatalf("reply=%#v effect=%#v", dropReply, effect)
	}
	count := 0
	maps.MapManager.MapEachItem(obj.MapNumber, func(reward *item.Item, _ int, x, y int, _ bool) {
		if reward == bag || x != obj.X || y != obj.Y {
			t.Errorf("reward=%#v position=%d,%d", reward, x, y)
		}
		count++
	})
	if count != 1 {
		t.Fatalf("reward count = %d, want 1", count)
	}
}

func TestDropUnregisteredItemUsesOrdinaryDrop(t *testing.T) {
	obj, actor := newBasicAttackTestObject(32, ObjectTypePlayer)
	withTestObjectManager(t, obj)
	obj.X, obj.Y = findBasicAttackMapPosition(t, false)
	actor.inventory.Size = 237
	actor.inventory.Items = make([]*item.Item, actor.inventory.Size)
	actor.inventory.Flags = make([]bool, actor.inventory.Size)
	it := item.NewItem(14, 11)
	it.Level = 4 // This level has no BAG_COMMON binding.
	actor.inventory.AddItem(12, it)
	defer func() {
		var indexes []int
		maps.MapManager.MapEachItem(obj.MapNumber, func(_ *item.Item, index, _, _ int, _ bool) {
			indexes = append(indexes, index)
		})
		for _, index := range indexes {
			maps.MapManager.RemoveItem(obj.MapNumber, index)
		}
	}()

	obj.DropItem(&model.MsgDropItem{Position: 12, X: obj.X, Y: obj.Y})
	if actor.inventory.Items[12] != nil || maps.MapManager.PeekItem(obj.MapNumber, 0) != it {
		t.Fatal("unregistered item did not follow ordinary drop")
	}
	if reply, ok := actor.messages[len(actor.messages)-1].(*model.MsgDropItemReply); !ok || reply.Result != 1 {
		t.Fatalf("ordinary drop reply = %#v", actor.messages)
	}
}

func TestDropCommonBagKeepsItemWhenRewardCannotBePlaced(t *testing.T) {
	obj, actor := newBasicAttackTestObject(32, ObjectTypePlayer)
	withTestObjectManager(t, obj)
	obj.X, obj.Y = -1, -1
	actor.inventory.Size = 237
	actor.inventory.Items = make([]*item.Item, actor.inventory.Size)
	actor.inventory.Flags = make([]bool, actor.inventory.Size)
	bag := item.NewItem(14, 11)
	bag.Level = 2
	actor.inventory.AddItem(12, bag)

	obj.DropItem(&model.MsgDropItem{Position: 12})
	if actor.inventory.Items[12] != bag {
		t.Fatal("CommonBag was consumed without placing its reward")
	}
	if reply, ok := actor.messages[len(actor.messages)-1].(*model.MsgDropItemReply); !ok || reply.Result != 0 {
		t.Fatalf("failed CommonBag drop reply = %#v", actor.messages)
	}
	for _, msg := range actor.messages {
		if effect, ok := msg.(*model.MsgServerCMDReply); ok {
			t.Fatalf("effect played after failed drop: %#v", effect)
		}
	}
}

func TestDropItemEffectSpecialCase(t *testing.T) {
	obj, actor := newBasicAttackTestObject(32, ObjectTypePlayer)
	withTestObjectManager(t, obj)
	actor.inventory.Size = 237
	actor.inventory.Items = make([]*item.Item, actor.inventory.Size)
	actor.inventory.Flags = make([]bool, actor.inventory.Size)
	it := item.NewItem(14, 63)
	actor.inventory.AddItem(12, it)

	obj.DropItem(&model.MsgDropItem{Position: 12})
	if actor.inventory.Items[12] != nil {
		t.Fatal("special item was not consumed")
	}
	if len(actor.messages) != 2 {
		t.Fatalf("special item messages = %#v", actor.messages)
	}
	if effect, ok := actor.messages[0].(*model.MsgServerCMDReply); !ok || effect.Type != 0 || effect.X != obj.X || effect.Y != obj.Y {
		t.Fatalf("special item effect = %#v", actor.messages[0])
	}
	if reply, ok := actor.messages[1].(*model.MsgDropItemReply); !ok || reply.Result != 1 {
		t.Fatalf("special item reply = %#v", actor.messages[1])
	}
}
