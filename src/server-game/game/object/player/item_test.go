package player

import (
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object"
)

func TestDropCommonBagConsumesItemAndPlaysEffect(t *testing.T) {
	p := newItemTestPlayer()
	obj := &p.Object
	obj.X, obj.Y = findItemTestMapPosition(t)
	bag := item.NewItem(14, 11)
	bag.Level = 2 // Firecracker: one reward and effect 0.
	p.inventory.AddItem(12, bag)
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
	messages := itemTestMessages(p)
	if p.inventory.Items[12] != nil {
		t.Fatal("CommonBag was not consumed")
	}
	var dropReply *model.MsgDropItemReply
	var effect *model.MsgServerCMDReply
	for _, msg := range messages {
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
	p := newItemTestPlayer()
	obj := &p.Object
	obj.X, obj.Y = findItemTestMapPosition(t)
	it := item.NewItem(14, 11)
	it.Level = 4 // This level has no BAG_COMMON binding.
	p.inventory.AddItem(12, it)
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
	messages := itemTestMessages(p)
	if p.inventory.Items[12] != nil || maps.MapManager.PeekItem(obj.MapNumber, 0) != it {
		t.Fatal("unregistered item did not follow ordinary drop")
	}
	if reply, ok := messages[len(messages)-1].(*model.MsgDropItemReply); !ok || reply.Result != 1 {
		t.Fatalf("ordinary drop reply = %#v", messages)
	}
}

func TestDropCommonBagKeepsItemWhenRewardCannotBePlaced(t *testing.T) {
	p := newItemTestPlayer()
	obj := &p.Object
	obj.X, obj.Y = -1, -1
	bag := item.NewItem(14, 11)
	bag.Level = 2
	p.inventory.AddItem(12, bag)

	obj.DropItem(&model.MsgDropItem{Position: 12})
	messages := itemTestMessages(p)
	if p.inventory.Items[12] != bag {
		t.Fatal("CommonBag was consumed without placing its reward")
	}
	if reply, ok := messages[len(messages)-1].(*model.MsgDropItemReply); !ok || reply.Result != 0 {
		t.Fatalf("failed CommonBag drop reply = %#v", messages)
	}
	for _, msg := range messages {
		if effect, ok := msg.(*model.MsgServerCMDReply); ok {
			t.Fatalf("effect played after failed drop: %#v", effect)
		}
	}
}

func TestDropItemEffectSpecialCase(t *testing.T) {
	p := newItemTestPlayer()
	obj := &p.Object
	it := item.NewItem(14, 63)
	p.inventory.AddItem(12, it)

	obj.DropItem(&model.MsgDropItem{Position: 12})
	messages := itemTestMessages(p)
	if p.inventory.Items[12] != nil {
		t.Fatal("special item was not consumed")
	}
	if len(messages) != 2 {
		t.Fatalf("special item messages = %#v", messages)
	}
	if effect, ok := messages[0].(*model.MsgServerCMDReply); !ok || effect.Type != 0 || effect.X != obj.X || effect.Y != obj.Y {
		t.Fatalf("special item effect = %#v", messages[0])
	}
	if reply, ok := messages[1].(*model.MsgDropItemReply); !ok || reply.Result != 1 {
		t.Fatalf("special item reply = %#v", messages[1])
	}
}

func newItemTestPlayer() *Player {
	p := &Player{msgChan: make(chan any, 100)}
	p.Init()
	p.Index = 32
	p.Level = 100
	p.Type = object.ObjectTypePlayer
	p.Objecter = p
	p.inventory.Size = 237
	p.inventory.Items = make([]*item.Item, p.inventory.Size)
	p.inventory.Flags = make([]bool, p.inventory.Size)
	return p
}

func itemTestMessages(p *Player) []any {
	messages := make([]any, 0, len(p.msgChan))
	for len(p.msgChan) > 0 {
		messages = append(messages, <-p.msgChan)
	}
	return messages
}

func findItemTestMapPosition(t *testing.T) (int, int) {
	t.Helper()
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			if maps.MapManager.GetMapAttr(0, x, y)&1 == 0 {
				return x, y
			}
		}
	}
	t.Fatal("map 0 has no available position")
	return 0, 0
}
