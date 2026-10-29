package player

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object"
)

func TestUsePetFollowsRequestedAction(t *testing.T) {
	p := newItemTestPlayer()
	pet := item.NewItem(13, 3)
	pet.Durability = 100
	other := item.NewItem(13, 2)
	other.Durability = 100
	late := item.NewItem(13, 2)
	late.Durability = 100
	p.inventory.AddItem(12, pet)
	p.inventory.AddItem(13, other)
	if !p.inventory.AddItem(204, late) {
		t.Fatal("could not place pet outside the main inventory")
	}

	use := func(position, value, wantResult int) {
		t.Helper()
		p.UsePet(&model.MsgUsePet{Position: position, Value: value})
		messages := itemTestMessages(p)
		reply, ok := messages[len(messages)-1].(*model.MsgUsePetReply)
		if !ok || reply.Result != wantResult {
			t.Fatalf("UsePet(%d, %d) reply = %#v, want %d", position, value, messages[len(messages)-1], wantResult)
		}
	}

	use(12, 0, 0)
	use(204, 254, 0)
	use(12, 254, -2)
	if p.pet != pet || pet.HarmonyOption != 1 || p.petIncreaseDamage != 15 {
		t.Fatal("pet was not activated")
	}
	use(13, 254, 0)
	use(13, 255, 0)
	if p.pet != pet || pet.HarmonyOption != 1 || other.HarmonyOption != 0 {
		t.Fatal("request for another pet changed the active pet")
	}
	pet.Durability = 0
	p.calc()
	itemTestMessages(p)
	if p.petIncreaseDamage != 0 || p.GetActivePetCode() != -1 {
		t.Fatal("broken pet still gives an effect")
	}
	use(12, 255, -1)
	if p.pet != nil || pet.HarmonyOption != 0 {
		t.Fatal("broken pet could not be deactivated")
	}
}

func TestFindAndUsePetKeepsOneValidMount(t *testing.T) {
	p := newItemTestPlayer()
	broken := item.NewItem(13, 2)
	broken.HarmonyOption = 1
	first := item.NewItem(13, 3)
	first.Durability, first.HarmonyOption = 100, 1
	second := item.NewItem(13, 2)
	second.Durability, second.HarmonyOption = 100, 1
	p.inventory.AddItem(12, broken)
	p.inventory.AddItem(13, first)
	p.inventory.AddItem(14, second)

	p.findAndUsePet()
	if p.pet != first || broken.HarmonyOption != 0 || first.HarmonyOption != 1 || second.HarmonyOption != 0 {
		t.Fatal("login did not keep exactly one valid mount")
	}
}

func TestLogoutRestoresPetOnNextCharacterLoad(t *testing.T) {
	p := newItemTestPlayer()
	pet := item.NewItem(13, 3)
	pet.Durability, pet.HarmonyOption = 100, 1
	p.inventory.AddItem(12, pet)
	p.pet = pet

	saved, err := json.Marshal(p.inventory)
	if err != nil {
		t.Fatal(err)
	}
	p.Logout(&model.MsgLogout{Flag: 1})
	if p.pet != nil || pet.HarmonyOption != 1 {
		t.Fatal("logout did not clear the active pet pointer while preserving its saved state")
	}
	if err := p.inventory.UnmarshalJSON(saved); err != nil {
		t.Fatal(err)
	}
	p.findAndUsePet()
	if p.pet != p.inventory.Item(12) || p.pet.HarmonyOption != 1 {
		t.Fatal("saved pet was not restored after logout")
	}
}

func TestGateMoveRestoresPetAppearance(t *testing.T) {
	for _, durability := range []int{100, 0, -1} {
		p := newItemTestPlayer()
		if durability >= 0 {
			pet := item.NewItem(13, 3)
			pet.Durability, pet.HarmonyOption = durability, 1
			p.inventory.AddItem(12, pet)
			p.pet = pet
		}
		p.Teleport(&model.MsgTeleport{GateNumber: 22})
		want := -1
		if durability > 0 {
			want = item.Code(13, 3)
		}
		teleported := false
		found := false
		for _, msg := range itemTestMessages(p) {
			switch reply := msg.(type) {
			case *model.MsgTeleportReply:
				teleported = true
			case *model.MsgCreateViewportPlayerReply:
				if !teleported || len(reply.Players) != 1 || reply.Players[0].MuunRideItem != want {
					t.Fatalf("durability %d: self viewport after teleport = %+v", durability, reply)
				}
				found = true
			}
		}
		if !found {
			t.Fatalf("durability %d: no self viewport after gate move", durability)
		}
	}
}

func TestRegenRestoresPetAppearance(t *testing.T) {
	for _, ride := range []bool{true, false} {
		p := newItemTestPlayer()
		want := -1
		if ride {
			pet := item.NewItem(13, 3)
			pet.Durability, pet.HarmonyOption = 100, 1
			p.inventory.AddItem(12, pet)
			p.pet = pet
			want = pet.Code
		}
		p.X, p.Y = 1, 1
		p.CreateFrustum()
		p.Regen()
		if maps.MapManager.GetMapAttr(p.MapNumber, p.TX, p.TY)&2 == 0 {
			t.Fatal("regen position is not occupied")
		}
		maps.MapManager.ClearMapAttrStand(p.MapNumber, p.TX, p.TY)
		for i := range p.FrustumX {
			if p.FrustumX[i] != object.FrustumX[i]+p.X || p.FrustumY[i] != object.FrustumY[i]+p.Y {
				t.Fatalf("regen frustum vertex %d does not match new position", i)
			}
		}
		messages := itemTestMessages(p)
		if len(messages) != 2 {
			t.Fatalf("ride %t: regen sent %d messages, want 2", ride, len(messages))
		}
		reload, ok := messages[0].(*model.MsgReloadCharacterReply)
		if !ok {
			t.Fatalf("ride %t: first message = %T", ride, messages[0])
		}
		viewport, ok := messages[1].(*model.MsgCreateViewportPlayerReply)
		if !ok || len(viewport.Players) != 1 {
			t.Fatalf("ride %t: second message = %#v", ride, messages[1])
		}
		self := viewport.Players[0]
		if self.X != reload.X || self.Y != reload.Y || self.MuunRideItem != want {
			t.Fatalf("ride %t: self viewport = %+v, reload = %+v", ride, self, reload)
		}
	}
}

func TestActivePetCannotLeaveInventory(t *testing.T) {
	p := newItemTestPlayer()
	pet := item.NewItem(13, 3)
	pet.Durability, pet.HarmonyOption = 100, 1
	p.inventory.AddItem(12, pet)
	p.pet = pet
	p.chaosBoxOpen = true
	p.warehouse.Size = 240
	p.warehouse.Items = make([]*item.Item, 240)
	p.warehouse.Flags = make([]bool, 240)
	for _, msg := range []model.MsgMoveItem{
		{SrcFlag: 0, SrcPosition: 12, DstFlag: 0, DstPosition: 13},
		{SrcFlag: 0, SrcPosition: 12, DstFlag: 1, DstPosition: 0},
		{SrcFlag: 0, SrcPosition: 12, DstFlag: 2, DstPosition: 0},
	} {
		p.MoveItem(&msg)
		messages := itemTestMessages(p)
		if reply, ok := messages[len(messages)-1].(*model.MsgMoveItemReply); !ok || reply.Result != -1 {
			t.Fatalf("active pet move reply = %#v", messages[len(messages)-1])
		}
		if p.inventory.Item(12) != pet || p.chaosBox.Item(0) != nil || p.warehouse.Item(0) != nil {
			t.Fatal("active pet was moved")
		}
	}

	p.DropItem(&model.MsgDropItem{Position: 12})
	messages := itemTestMessages(p)
	if reply, ok := messages[len(messages)-1].(*model.MsgDropItemReply); !ok || reply.Result != 0 || p.inventory.Item(12) != pet {
		t.Fatal("active pet was dropped")
	}

	shop, err := object.ObjectManager.AddMonster(func() *object.Object {
		shop := &object.Object{NpcType: object.NpcTypeShop, NoRegen: true}
		shop.Init()
		return shop
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { object.ObjectManager.DeleteEventMonster(shop, time.Now()) })
	p.TargetNumber = shop.Index
	p.Money = 100
	p.SellItem(&model.MsgSellItem{Position: 12})
	messages = itemTestMessages(p)
	if reply, ok := messages[len(messages)-1].(*model.MsgSellItemReply); !ok || reply.Result != 0 || p.inventory.Item(12) != pet || p.Money != 100 {
		t.Fatal("active pet was sold")
	}
}

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

func TestMoveItemRejectsInvalidPositions(t *testing.T) {
	p := newItemTestPlayer()
	it := item.NewItem(12, 15)
	p.inventory.AddItem(12, it)
	p.warehouse.Size = 240
	p.warehouse.Items = make([]*item.Item, 240)
	p.warehouse.Flags = make([]bool, 240)
	for _, msg := range []model.MsgMoveItem{
		{SrcFlag: 0, SrcPosition: -1, DstFlag: 0, DstPosition: 13},
		{SrcFlag: 0, SrcPosition: 12, DstFlag: 0, DstPosition: -1},
		{SrcFlag: 2, SrcPosition: -1, DstFlag: 0, DstPosition: 13},
		{SrcFlag: 0, SrcPosition: 12, DstFlag: 2, DstPosition: 240},
	} {
		p.MoveItem(&msg)
		messages := itemTestMessages(p)
		if len(messages) != 1 {
			t.Fatalf("invalid move replies = %d, want 1", len(messages))
		}
		if reply, ok := messages[0].(*model.MsgMoveItemReply); !ok || reply.Result != -1 {
			t.Fatalf("invalid move reply = %#v", messages[0])
		}
		if p.inventory.Item(12) != it {
			t.Fatal("invalid move changed the item")
		}
	}
}

func TestMoveItemBetweenInventoryAndWarehouse(t *testing.T) {
	p := newItemTestPlayer()
	p.warehouse.Size = 240
	p.warehouse.Items = make([]*item.Item, 240)
	p.warehouse.Flags = make([]bool, 240)
	inventoryItem := item.NewItem(12, 15)
	warehouseItem := item.NewItem(12, 15)
	p.inventory.AddItem(12, inventoryItem)
	p.warehouse.AddItem(0, warehouseItem)

	for _, msg := range []model.MsgMoveItem{
		{SrcFlag: 0, SrcPosition: 12, DstFlag: 2, DstPosition: 0},
		{SrcFlag: 2, SrcPosition: 0, DstFlag: 0, DstPosition: 12},
	} {
		p.MoveItem(&msg)
		messages := itemTestMessages(p)
		if len(messages) != 1 {
			t.Fatalf("blocked move replies = %d, want 1", len(messages))
		}
		if reply, ok := messages[0].(*model.MsgMoveItemReply); !ok || reply.Result != -1 {
			t.Fatalf("blocked move reply = %#v", messages[0])
		}
		if p.inventory.Item(12) != inventoryItem || p.warehouse.Item(0) != warehouseItem {
			t.Fatal("blocked move changed the source or destination")
		}
	}

	p.MoveItem(&model.MsgMoveItem{SrcFlag: 0, SrcPosition: 12, DstFlag: 2, DstPosition: 1})
	itemTestMessages(p)
	if p.inventory.Item(12) != nil || p.warehouse.Item(1) != inventoryItem {
		t.Fatal("inventory item was not moved to the warehouse")
	}
	p.MoveItem(&model.MsgMoveItem{SrcFlag: 2, SrcPosition: 0, DstFlag: 0, DstPosition: 12})
	itemTestMessages(p)
	if p.warehouse.Item(0) != nil || p.inventory.Item(12) != warehouseItem {
		t.Fatal("warehouse item was not moved to the inventory")
	}
}

func TestChaosBoxReleasesInventoryPosition(t *testing.T) {
	p := newItemTestPlayer()
	p.chaosBoxOpen = true
	it := item.NewItem(12, 15)
	p.inventory.AddItem(12, it)
	p.MoveItem(&model.MsgMoveItem{SrcFlag: 0, SrcPosition: 12, DstFlag: 1, DstPosition: 32})
	if messages := itemTestMessages(p); len(messages) != 1 {
		t.Fatalf("failed move replies = %d, want 1", len(messages))
	} else if reply, ok := messages[0].(*model.MsgMoveItemReply); !ok || reply.Result != -1 {
		t.Fatalf("failed move reply = %#v", messages[0])
	}
	if p.inventory.Items[12] != it || !p.inventory.Flags[12] || p.chaosBox.Item(0) != nil {
		t.Fatal("failed move changed the item")
	}
	p.MoveItem(&model.MsgMoveItem{SrcFlag: 0, SrcPosition: 12, DstFlag: 1, DstPosition: 0})
	if messages := itemTestMessages(p); len(messages) != 1 {
		t.Fatalf("successful move replies = %d, want 1", len(messages))
	} else if reply, ok := messages[0].(*model.MsgMoveItemReply); !ok || reply.Result != 1 || reply.Position != 0 || reply.Item != it {
		t.Fatalf("successful move reply = %#v", messages[0])
	}
	if p.inventory.Items[12] != nil || p.inventory.Flags[12] || p.chaosBox.Item(0) != it {
		t.Fatal("moving into Chaos Box did not release the inventory position")
	}
	other := item.NewItem(12, 15)
	p.inventory.AddItem(12, other)
	p.pushChaosBoxItems()
	p.CloseChaosBox(&model.MsgEmpty{})
	if p.inventory.Items[12] != other || p.chaosBox.Item(0) != nil || p.chaosBoxOpen {
		t.Fatal("closing Chaos Box did not discard its contents")
	}
	if messages := itemTestMessages(p); len(messages) < 2 {
		t.Fatal("closing Chaos Box did not reply")
	} else if _, ok := messages[len(messages)-1].(*model.MsgCloseChaosBoxReply); !ok {
		t.Fatalf("last Chaos Box message = %T", messages[len(messages)-1])
	} else if list, ok := messages[len(messages)-2].(*model.MsgTypeItemListReply); !ok || list.Type != 3 {
		t.Fatalf("Chaos Box list = %#v", messages[len(messages)-2])
	}
}

func TestChaosBoxMoveItemDirections(t *testing.T) {
	p := newItemTestPlayer()
	it := item.NewItem(12, 15)
	p.inventory.AddItem(12, it)

	p.MoveItem(&model.MsgMoveItem{SrcFlag: 0, SrcPosition: 12, DstFlag: 1, DstPosition: 0})
	if messages := itemTestMessages(p); len(messages) != 1 {
		t.Fatalf("closed box replies = %d, want 1", len(messages))
	} else if reply, ok := messages[0].(*model.MsgMoveItemReply); !ok || reply.Result != -1 {
		t.Fatalf("closed box reply = %#v", messages[0])
	}
	if p.inventory.Items[12] != it {
		t.Fatal("closed box changed the item")
	}

	p.chaosBoxOpen = true
	p.MoveItem(&model.MsgMoveItem{SrcFlag: 0, SrcPosition: 12, DstFlag: 1, DstPosition: 0})
	itemTestMessages(p)
	p.MoveItem(&model.MsgMoveItem{SrcFlag: 1, SrcPosition: 0, DstFlag: 1, DstPosition: 1})
	if messages := itemTestMessages(p); len(messages) != 1 {
		t.Fatalf("box move replies = %d, want 1", len(messages))
	} else if reply, ok := messages[0].(*model.MsgMoveItemReply); !ok || reply.Result != 1 || reply.Item != it {
		t.Fatalf("box move reply = %#v", messages[0])
	}
	if p.chaosBox.Item(0) != nil || p.chaosBox.Item(1) != it {
		t.Fatal("box move did not preserve the item")
	}

	p.MoveItem(&model.MsgMoveItem{SrcFlag: 1, SrcPosition: 1, DstFlag: 0, DstPosition: 12})
	if messages := itemTestMessages(p); len(messages) != 1 {
		t.Fatalf("inventory return replies = %d, want 1", len(messages))
	} else if reply, ok := messages[0].(*model.MsgMoveItemReply); !ok || reply.Result != 0 || reply.Item != it {
		t.Fatalf("inventory return reply = %#v", messages[0])
	}
	if p.inventory.Items[12] != it || p.chaosBox.Item(1) != nil {
		t.Fatal("box item did not return to inventory")
	}
}

func TestChaosBoxLargerOutputUsesCurrentSize(t *testing.T) {
	p := newItemTestPlayer()
	p.chaosBoxOpen = true
	small := item.NewItem(12, 15)
	p.inventory.AddItem(12, small)
	p.MoveItem(&model.MsgMoveItem{SrcFlag: 0, SrcPosition: 12, DstFlag: 1, DstPosition: 0})
	itemTestMessages(p)

	other := item.NewItem(12, 15)
	p.inventory.AddItem(12, other)
	large := &item.Item{ItemBase: &item.ItemBase{Width: 2, Height: 2}}
	if !p.chaosBox.ReplaceItem(0, large) {
		t.Fatal("could not replace the box item")
	}
	p.MoveItem(&model.MsgMoveItem{SrcFlag: 1, SrcPosition: 0, DstFlag: 0, DstPosition: 12})
	if messages := itemTestMessages(p); len(messages) != 1 {
		t.Fatalf("blocked move replies = %d, want 1", len(messages))
	} else if reply, ok := messages[0].(*model.MsgMoveItemReply); !ok || reply.Result != -1 {
		t.Fatalf("blocked move reply = %#v", messages[0])
	}
	if p.inventory.Item(12) != other || !p.inventory.Flags[12] || p.chaosBox.Item(0) != large {
		t.Fatal("blocked move changed the inventory or box")
	}

	p.MoveItem(&model.MsgMoveItem{SrcFlag: 1, SrcPosition: 0, DstFlag: 0, DstPosition: 13})
	if messages := itemTestMessages(p); len(messages) != 1 {
		t.Fatalf("successful move replies = %d, want 1", len(messages))
	} else if reply, ok := messages[0].(*model.MsgMoveItemReply); !ok || reply.Result != 0 {
		t.Fatalf("successful move reply = %#v", messages[0])
	}
	if p.inventory.Item(12) != other || p.inventory.Item(13) != large || p.chaosBox.Item(0) != nil ||
		!p.inventory.Flags[13] || !p.inventory.Flags[14] || !p.inventory.Flags[21] || !p.inventory.Flags[22] {
		t.Fatal("larger output did not occupy its current size in the inventory")
	}
}

func TestChaosBoxInsertsSeedSphere(t *testing.T) {
	p := newItemTestPlayer()
	p.chaosBoxOpen = true
	weapon := item.NewItem(0, 26)
	weapon.MakeSocketSlots(1)
	seedSphere := item.SocketManager.SeedSphereItem(0, 1)
	materials := []*item.Item{weapon, seedSphere, item.NewItem(12, 15), item.NewItem(14, 22)}
	origins := []int{12, 20, 21, 22}
	for i, it := range materials {
		p.inventory.AddItem(origins[i], it)
		p.MoveItem(&model.MsgMoveItem{SrcFlag: 0, SrcPosition: origins[i], DstFlag: 1, DstPosition: i})
	}
	p.MixItems(&model.MsgMixItems{Type: 0x2c, Position: 0, MixCount: 1})
	if p.chaosBox.Item(1) != seedSphere || p.chaosBox.Item(0).SocketSlot1 != 1 {
		t.Fatal("insufficient Zen changed the materials")
	}
	p.Money = 1_000_000
	p.MixItems(&model.MsgMixItems{Type: 0x2c, Position: 0, MixCount: 1})
	if p.chaosBox.Item(0) == nil || p.chaosBox.Item(0).SocketSlot1 != 100 || p.Money != 0 {
		t.Fatalf("insert result: item=%#v money=%d", p.chaosBox.Item(0), p.Money)
	}
	for i := 1; i < 4; i++ {
		if p.chaosBox.Item(i) != nil || p.inventory.Flags[origins[i]] {
			t.Fatalf("material at %d survived insertion", i)
		}
	}
}

func TestChaosBoxComposesSeedSphere(t *testing.T) {
	p := newItemTestPlayer()
	p.chaosBoxOpen = true
	p.Money = 1_000_000
	seed := item.NewItem(12, 60)
	seed.Level = 0
	seed.Calc()
	seed.BuyMoney = 4_000_000 // Make the C++ success formula deterministic.
	materials := []*item.Item{seed, item.NewItem(12, 70), item.NewItem(12, 15), item.NewItem(14, 22)}
	for i, it := range materials {
		p.inventory.AddItem(12+i, it)
		p.MoveItem(&model.MsgMoveItem{SrcFlag: 0, SrcPosition: 12 + i, DstFlag: 1, DstPosition: i})
	}
	p.MixItems(&model.MsgMixItems{Type: 0x2b, MixCount: 1})
	if result := p.chaosBox.Item(0); result == nil || result.Code != item.Code(12, 100) || result.Level != 0 || p.Money != 0 {
		t.Fatalf("Seed Sphere not created in Chaos Box: item=%#v money=%d", result, p.Money)
	}
}

func TestMixItemsRejectsUnknownType(t *testing.T) {
	p := newItemTestPlayer()
	p.chaosBoxOpen, p.Money = true, 1_000_000
	material := item.NewItem(12, 15)
	if !p.chaosBox.AddItem(0, material) {
		t.Fatal("could not place material")
	}
	p.MixItems(&model.MsgMixItems{Type: 0x2e})
	if p.Money != 1_000_000 || p.chaosBox.Item(0) != material {
		t.Fatal("unknown mix changed money or materials")
	}
	messages := itemTestMessages(p)
	if len(messages) != 1 {
		t.Fatalf("unknown mix replies = %d, want 1", len(messages))
	}
	if reply, ok := messages[0].(*model.MsgMixItemsReply); !ok || reply.Result != 7 {
		t.Fatalf("unknown mix reply = %#v", messages[0])
	}
}

func TestChaosBoxCompositeSucceedsWhenInventoryIsFull(t *testing.T) {
	p := newItemTestPlayer()
	p.chaosBoxOpen, p.Money = true, 1_000_000
	seed := item.NewItem(12, 60)
	seed.Calc()
	seed.BuyMoney = 4_000_000
	materials := []*item.Item{seed, item.NewItem(12, 70), item.NewItem(12, 15), item.NewItem(14, 22)}
	for position, it := range materials {
		if !p.chaosBox.AddItem(position, it) {
			t.Fatalf("could not place material at %d", position)
		}
	}
	for position := 12; position < p.inventory.Size; position++ {
		p.inventory.Flags[position] = true
	}
	p.MixItems(&model.MsgMixItems{Type: mixSocketSeedSphereComposite, MixCount: 1})
	messages := itemTestMessages(p)
	if len(messages) < 2 {
		t.Fatalf("full inventory mix messages = %d", len(messages))
	}
	if update, ok := messages[len(messages)-2].(*model.MsgTypeItemListReply); !ok || update.Items[0] != p.chaosBox.Item(0) {
		t.Fatalf("full inventory box update = %#v", messages[len(messages)-2])
	}
	if reply, ok := messages[len(messages)-1].(*model.MsgMixItemsReply); !ok || reply.Result != 1 || reply.Item != p.chaosBox.Item(0) {
		t.Fatalf("full inventory mix reply = %#v", messages[len(messages)-1])
	}
	if result := p.chaosBox.Item(0); result == nil || result.Code != item.Code(12, 100) || p.Money != 0 {
		t.Fatalf("full inventory mix result: item=%#v money=%d", result, p.Money)
	}
	for position := 1; position < len(materials); position++ {
		if p.chaosBox.Item(position) != nil {
			t.Fatalf("material at %d survived", position)
		}
	}
}

func TestChaosBoxComposesHigherSeedSphere(t *testing.T) {
	p := newItemTestPlayer()
	p.chaosBoxOpen, p.Money = true, 1_000_000
	third := item.SocketManager.SeedSphereItem(0, 3)
	sphere := item.NewItem(12, 73)
	sphere.Calc()
	sphere.BuyMoney = 4_000_000
	materials := []*item.Item{third, sphere, item.NewItem(12, 141), item.NewItem(12, 137), item.NewItem(14, 245)}
	for i, it := range materials {
		p.inventory.AddItem(12+i, it)
		p.MoveItem(&model.MsgMoveItem{SrcFlag: 0, SrcPosition: 12 + i, DstFlag: 1, DstPosition: i})
	}
	p.MixItems(&model.MsgMixItems{Type: 0x2b, MixCount: 1})
	if result := p.chaosBox.Item(0); result != nil && result.Code == item.Code(12, 118) && result.Level == 0 {
		return
	}
	t.Fatal("higher-rank Seed Sphere was not created in Chaos Box")
}

func TestChaosBoxRemovesSeedSphere(t *testing.T) {
	p := newItemTestPlayer()
	p.chaosBoxOpen, p.Money = true, 1_000_000
	weapon := item.NewItem(0, 26)
	weapon.MakeSocketSlots(1)
	weapon.SocketSlot1 = 100
	weapon.Calc()
	materials := []*item.Item{weapon}
	for i := 0; i < 5; i++ {
		materials = append(materials, item.NewItem(12, 15))
	}
	for i := 0; i < 5; i++ {
		materials = append(materials, item.NewItem(14, 42))
	}
	materials = append(materials, item.NewItem(14, 31))
	origins := []int{12, 13, 14, 15, 16, 17, 18, 19, 21, 22, 23, 24}
	positions := []int{0, 1, 2, 3, 4, 5, 6, 7, 9, 10, 11, 12}
	for i, it := range materials {
		p.inventory.AddItem(origins[i], it)
		p.MoveItem(&model.MsgMoveItem{SrcFlag: 0, SrcPosition: origins[i], DstFlag: 1, DstPosition: positions[i]})
	}
	for position := 12; position < p.inventory.Size; position++ {
		if !p.inventory.AddItem(position, item.NewItem(12, 15)) {
			t.Fatalf("could not fill inventory position %d", position)
		}
	}
	itemTestMessages(p)
	p.MixItems(&model.MsgMixItems{Type: 0x2d, Position: 0, MixCount: 1})
	messages := itemTestMessages(p)
	if reply, ok := messages[len(messages)-1].(*model.MsgMixItemsReply); !ok || reply.Result != 6 {
		t.Fatalf("full inventory removal reply = %#v", messages[len(messages)-1])
	}
	if p.Money != 1_000_000 || p.chaosBox.Item(0) != weapon || weapon.SocketSlot1 != 100 {
		t.Fatal("full inventory changed the money or equipment")
	}
	for i, position := range positions[1:] {
		if p.chaosBox.Item(position) != materials[i+1] {
			t.Fatalf("full inventory consumed material at %d", position)
		}
	}
	p.inventory.RemoveItem(12, p.inventory.Item(12))
	p.MixItems(&model.MsgMixItems{Type: 0x2d, Position: 0, MixCount: 1})
	messages = itemTestMessages(p)
	if p.Money != 0 || p.chaosBox.Item(0) == nil {
		t.Fatalf("remove mix did not finish: money=%d item=%#v", p.Money, p.chaosBox.Item(0))
	}
	for _, msg := range messages {
		if _, ok := msg.(*model.MsgPickItemReply); ok {
			t.Fatal("removal sent an inventory item reply")
		}
	}
	seedCount := 0
	for position, it := range p.chaosBox.Items() {
		if position == 0 || it == nil {
			continue
		}
		if it.Code != item.Code(12, 60) || it.Level != 0 {
			t.Fatalf("unexpected Chaos Box item at %d: %#v", position, it)
		}
		seedCount++
	}
	if p.chaosBox.Item(0).SocketSlot1 == 1 {
		if seedCount != 1 {
			t.Fatalf("successful removal returned %d Seeds to Chaos Box", seedCount)
		}
	} else if p.chaosBox.Item(0).SocketSlot1 != 100 || seedCount != 0 {
		t.Fatalf("failed removal: socket=%d Seeds=%d", p.chaosBox.Item(0).SocketSlot1, seedCount)
	}
}

func TestChaosBoxExtractsSeed(t *testing.T) {
	p := newItemTestPlayer()
	p.chaosBoxOpen, p.Money = true, 1_000_000
	excellent := item.NewItem(0, 26)
	excellent.Level, excellent.ExcellentAttackRate = 4, true
	excellent.Calc()
	setItem := item.NewItem(0, 17)
	setItem.Level, setItem.Set = 4, 1
	setItem.Calc()
	materials := []*item.Item{excellent, setItem, item.NewItem(14, 42), item.NewItem(12, 15), item.NewItem(14, 22)}
	origins := []int{12, 16, 13, 14, 15}
	positions := []int{0, 4, 1, 2, 3}
	for i, it := range materials {
		p.inventory.AddItem(origins[i], it)
		p.MoveItem(&model.MsgMoveItem{SrcFlag: 0, SrcPosition: origins[i], DstFlag: 1, DstPosition: positions[i]})
	}
	p.MixItems(&model.MsgMixItems{Type: 0x2a, MixCount: 1})
	if p.Money != 0 {
		t.Fatalf("extract fee was not charged: %d", p.Money)
	}
	for _, position := range positions[2:] {
		if p.chaosBox.Item(position) != nil {
			t.Fatalf("jewel at %d survived", position)
		}
	}
	if p.chaosBox.Item(0) == nil && p.chaosBox.Item(4) == nil {
		t.Fatal("successful extraction did not create a Seed in Chaos Box")
	} else if seed := p.chaosBox.Item(0); seed != nil && seed.KindB == item.KindBSocketSeed {
		if p.chaosBox.Item(4) != nil {
			t.Fatal("successful extraction retained materials")
		}
	} else if p.chaosBox.Item(0) == nil || p.chaosBox.Item(4) == nil || excellent.Level > 4 || setItem.Level > 4 {
		t.Fatal("failed extraction increased equipment level")
	}
}
