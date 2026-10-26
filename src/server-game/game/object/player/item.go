package player

import (
	"log/slog"
	"math"
	"math/rand"

	"github.com/xujintao/balgass/src/server-game/game/drop"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/shop"
)

func (p *Player) PickItem(msg *model.MsgPickItem) {
	reply := model.MsgPickItemReply{
		Result: -1,
	}
	itemDurChanged := false
	position := -1
	var it *item.Item
	defer func() {
		p.Push(&reply)
		if itemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   position,
				Durability: it.Durability,
				Flag:       0,
			}
			p.Push(&reply)
		}
	}()
	it2 := maps.MapManager.PeekItem(p.MapNumber, msg.Index)
	if it2 == nil {
		return
	}
	switch it2.Code {
	case item.Code(14, 15): // zen
		money := it2.Durability
		if p.Money+money > object.MaxZen {
			return
		}
		p.Money += money
		reply.Money = p.Money
		position = -2
	default:
		position = p.inventory.FindFreePositionForItem(it2)
		if position == -1 {
			return
		}
		it = p.inventory.Items[position]
		if it == nil {
			p.inventory.AddItem(position, it2)
		} else {
			it.Durability += it2.Durability
			itemDurChanged = true
		}
	}
	maps.MapManager.RemoveItem(p.MapNumber, msg.Index)
	reply.Result = position
	reply.Item = it2
}

func (p *Player) DropItem(msg *model.MsgDropItem) {
	reply := model.MsgDropItemReply{
		Result:   0,
		Position: msg.Position,
	}
	defer p.Push(&reply)
	// validate
	if msg.Position < 0 || msg.Position >= p.inventory.Size {
		return
	}
	it := p.inventory.Items[msg.Position]
	if it == nil {
		return
	}
	switch it.Code {
	case item.Code(14, 63): // Firecracker 爆竹
		cmdReply := model.MsgServerCMDReply{
			Type: 0,
			X:    p.X,
			Y:    p.Y,
		}
		p.PushViewport(&cmdReply)
	default:
		result := drop.DropManager.Drop(drop.Request{
			Trigger:   drop.TriggerPlayerDropItem,
			ItemCode:  it.Code,
			ItemLevel: it.Level,
		})
		if result.Handled {
			placed := 0
			for _, reward := range result.Rewards {
				if reward.Item == nil {
					continue
				}
				x, y := p.X, p.Y
				if reward.Placement == drop.PlaceNearSource {
					x, y = commonBagDropPosition(p.MapNumber, x, y)
				}
				if maps.MapManager.AddItem(p.MapNumber, x, y, reward.Item) {
					placed++
				}
			}
			if len(result.Rewards) > 0 && placed == 0 {
				return
			}
			if placed < len(result.Rewards) {
				slog.Warn("CommonBag rewards could not all be placed", "item", it.Code, "placed", placed, "rewards", len(result.Rewards))
			}
			if result.Effect != nil {
				x, y := p.X, p.Y
				if *result.Effect == 58 {
					x, y = p.Index>>8, p.Index&0xff
				}
				p.PushViewport(&model.MsgServerCMDReply{Type: *result.Effect, X: x, Y: y})
			}
		} else if !maps.MapManager.AddItem(p.MapNumber, msg.X, msg.Y, it) {
			return
		}
	}
	p.inventory.RemoveItem(msg.Position, it)
	if msg.Position < 12 || msg.Position == 126 {
		p.equipmentChanged()
	}
	reply.Result = 1
}

func commonBagDropPosition(mapNumber, x, y int) (int, int) {
	for i := 0; i < 10; i++ {
		dx, dy := rand.Intn(5), rand.Intn(5)
		if rand.Intn(2) == 0 {
			dx = -dx
		}
		if rand.Intn(2) == 0 {
			dy = -dy
		}
		dropX, dropY := x+dx, y+dy
		if maps.MapManager.GetMapAttr(mapNumber, dropX, dropY)&(4|8) == 0 {
			return dropX, dropY
		}
	}
	return x, y
}

func (p *Player) BuyItem(msg *model.MsgBuyItem) {
	reply := model.MsgBuyItemReply{
		Result: -1,
	}
	itemDurChanged := false
	position := -1
	var it *item.Item
	defer func() {
		p.Push(&reply)
		if itemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   position,
				Durability: it.Durability,
				Flag:       0,
			}
			p.Push(&reply)
		}
	}()
	// validate
	if msg.Position < 0 ||
		msg.Position >= shop.MaxShopItemCount ||
		p.TargetNumber < 0 {
		return
	}
	tobj := object.ObjectManager.GetObject(p.TargetNumber)
	if tobj == nil {
		return
	}
	if math.Abs(float64(p.X-tobj.X)) > 5 ||
		math.Abs(float64(p.Y-tobj.Y)) > 5 {
		return
	}
	if tobj.NpcType != object.NpcTypeShop {
		return
	}
	sit := shop.ShopManager.GetShopItem(tobj.Class, tobj.MapNumber, msg.Position)
	if sit == nil {
		return
	}
	position = p.inventory.FindFreePositionForItem(sit)
	if position == -1 {
		return
	}
	if p.Money < sit.BuyMoney {
		return
	}
	p.Money -= sit.BuyMoney
	p.PushMoney()
	it = p.inventory.Items[position]
	if it == nil {
		p.inventory.AddItem(position, sit)
	} else {
		it.Durability += sit.Durability
		itemDurChanged = true
	}
	reply.Result = position
	reply.Item = sit
}

func (p *Player) SellItem(msg *model.MsgSellItem) {
	reply := model.MsgSellItemReply{
		Result: 0,
	}
	defer p.Push(&reply)
	// validate
	if msg.Position < 0 ||
		msg.Position >= p.inventory.Size ||
		p.TargetNumber < 0 {
		return
	}
	tobj := object.ObjectManager.GetObject(p.TargetNumber)
	if tobj == nil {
		return
	}
	if math.Abs(float64(p.X-tobj.X)) > 5 ||
		math.Abs(float64(p.Y-tobj.Y)) > 5 {
		return
	}
	if tobj.NpcType != object.NpcTypeShop {
		return
	}
	it := p.inventory.Items[msg.Position]
	if it == nil {
		return
	}
	if p.Money+it.SellMoney > object.MaxZen {
		return
	}
	p.Money += it.SellMoney
	p.inventory.RemoveItem(msg.Position, it)
	reply.Result = 1
	reply.Money = p.Money
}

func (p *Player) MoveItem(msg *model.MsgMoveItem) {
	reply := model.MsgMoveItemReply{
		Result: -1,
	}
	sitemDurChanged := false
	titemDurChanged := false
	var sitem *item.Item
	var titem *item.Item
	defer func() {
		p.Push(&reply)
		if sitemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   msg.SrcPosition,
				Durability: sitem.Durability,
				Flag:       0,
			}
			p.Push(&reply)
		}
		if titemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   msg.DstPosition,
				Durability: titem.Durability,
				Flag:       0,
			}
			p.Push(&reply)
		}
	}()
	// get source item
	switch msg.SrcFlag {
	case 0: // inventory
		sitem = p.inventory.Item(msg.SrcPosition)
	case 1: // Chaos Box
		if !p.chaosBoxOpen {
			return
		}
		sitem = p.chaosBox.Item(msg.SrcPosition)
	case 2: // warehouse
		sitem = p.warehouse.Item(msg.SrcPosition)
	}
	if sitem == nil {
		return
	}
	// get destination item
	switch msg.DstFlag {
	case 0: // inventory
		titem = p.inventory.Item(msg.DstPosition)
	case 1: // Chaos Box
		if !p.chaosBoxOpen {
			return
		}
		titem = p.chaosBox.Item(msg.DstPosition)
	case 2: // warehouse
		titem = p.warehouse.Item(msg.DstPosition)
	}
	// move item
	switch msg.SrcFlag {
	case 0: // inventory
		switch msg.DstFlag {
		case 0: // inventory to inventory
			switch {
			case titem == nil: // move
				ok := p.inventory.CanPlace(msg.DstPosition, sitem)
				if !ok {
					return
				}
				p.inventory.RemoveItem(msg.SrcPosition, sitem)
				p.inventory.AddItem(msg.DstPosition, sitem)
				if msg.SrcPosition < 12 || msg.SrcPosition == 236 ||
					msg.DstPosition < 12 || msg.DstPosition == 236 {
					p.equipmentChanged()
				}
				reply.Result = msg.DstFlag
			case titem.Overlap != 0 && // overlap
				titem.Code == sitem.Code &&
				titem.Level == sitem.Level &&
				titem.Durability < titem.Overlap:
				delta := titem.Overlap - titem.Durability
				if delta > sitem.Durability {
					delta = sitem.Durability
				}
				sitem.Durability -= delta
				sitemDurChanged = true
				if sitem.Durability <= 0 {
					reply.Result = msg.DstFlag
					sitem.Durability = 0
					sitemDurChanged = false
					p.inventory.RemoveItem(msg.SrcPosition, sitem)
					reply := model.MsgDeleteInventoryItemReply{
						Position: msg.SrcPosition,
						Flag:     1,
					}
					p.Push(&reply)
				} else {
					reply.Result = -1
				}
				titem.Durability += delta
				titemDurChanged = true
			default:
				return
			}
		case 1: // inventory to Chaos Box
			if msg.SrcPosition < 12 || !p.chaosBox.AddItem(msg.DstPosition, sitem) {
				return
			}
			p.inventory.RemoveItem(msg.SrcPosition, sitem)
			reply.Result = msg.DstFlag
		case 2: // inventory to warehouse
			if titem == nil {
				ok := p.warehouse.CanPlace(msg.DstPosition, sitem)
				if !ok {
					return
				}
				p.inventory.RemoveItem(msg.SrcPosition, sitem)
				p.warehouse.AddItem(msg.DstPosition, sitem)
				if msg.SrcPosition < 12 || msg.SrcPosition == 236 {
					p.equipmentChanged()
				}
				reply.Result = msg.DstFlag
			}
		default:
			return
		}
	case 1: // Chaos Box
		switch msg.DstFlag {
		case 0: // Chaos Box to inventory
			if msg.DstPosition < 12 {
				return
			}
			if !p.inventory.CanPlace(msg.DstPosition, sitem) {
				return
			}
			p.inventory.AddItem(msg.DstPosition, sitem)
			p.chaosBox.RemoveItem(msg.SrcPosition)
			reply.Result = msg.DstFlag
		case 1: // Chaos Box to Chaos Box
			if !p.chaosBox.MoveItem(msg.SrcPosition, msg.DstPosition) {
				return
			}
			reply.Result = msg.DstFlag
		default:
			return
		}
	case 2: // warehouse
		switch msg.DstFlag {
		case 0: // warehouse to inventory
			if titem == nil {
				ok := p.inventory.CanPlace(msg.DstPosition, sitem)
				if !ok {
					return
				}
				p.warehouse.RemoveItem(msg.SrcPosition, sitem)
				p.inventory.AddItem(msg.DstPosition, sitem)
				if msg.DstPosition < 12 || msg.DstPosition == 236 {
					p.equipmentChanged()
				}
				reply.Result = msg.DstFlag
			}
		case 2: // warehouse to warehouse
			if titem == nil {
				ok := p.warehouse.CanPlace(msg.DstPosition, sitem)
				if !ok {
					return
				}
				p.warehouse.RemoveItem(msg.SrcPosition, sitem)
				p.warehouse.AddItem(msg.DstPosition, sitem)
				reply.Result = msg.DstFlag
			}
		default:
			return
		}
	default:
		return
	}
	reply.Position = msg.DstPosition
	reply.Item = sitem
}

func (p *Player) MixItems(msg *model.MsgMixItems) {
	reply := &model.MsgMixItemsReply{Result: 7}
	defer p.Push(reply)
	if !p.chaosBoxOpen {
		return
	}
	fee, ok := item.SocketManager.SeedFee(msg.Type)
	if !ok {
		return
	}
	var output *item.Item
	var targetPosition = -1
	var rate = 100
	var remove = false
	items := make([]int, 0, 12)
	for position, it := range p.chaosBox.Items() {
		if it != nil {
			items = append(items, position)
		}
	}
	switch msg.Type {
	case item.SocketMixSeedExtract: // Extract a Seed from one excellent and one set item.
		if len(items) != 5 {
			return
		}
		excellent, setItem, harmony, chaos, creation := -1, -1, -1, -1, -1
		for _, position := range items {
			it := p.chaosBox.Item(position)
			switch {
			case it.IsExcellent() && it.Level >= 4 && excellent < 0:
				excellent = position
			case it.Set != 0 && it.Level >= 4 && setItem < 0:
				setItem = position
			case it.Code == item.Code(14, 42) && harmony < 0:
				harmony = position
			case it.Code == item.Code(12, 15) && chaos < 0:
				chaos = position
			case it.Code == item.Code(14, 22) && creation < 0:
				creation = position
			default:
				return
			}
		}
		if excellent < 0 || setItem < 0 || harmony < 0 || chaos < 0 || creation < 0 {
			return
		}
		code, level, found := item.SocketManager.RandomSeed()
		if !found {
			return
		}
		output = item.NewItem(code/512, code%512)
		output.Level, output.Durability = level, 1
		output.Calc()
		rate = 80 + (p.chaosBox.Item(excellent).BuyMoney+p.chaosBox.Item(setItem).BuyMoney)/2000000
		if rate > 90 {
			rate = 90
		}
	case item.SocketMixSeedSphereComposite: // Seed + Sphere, or third-rank Seed Sphere + Sphere.
		seed, sphere, chaos, creation, chaosBundle, creationBundle := -1, -1, -1, -1, -1, -1
		tetra := -1
		for _, position := range items {
			it := p.chaosBox.Item(position)
			switch {
			case it.KindB == item.KindBSocketSeed && seed < 0:
				seed = position
			case it.KindB == item.KindBSocketSphere && sphere < 0:
				sphere = position
			case it.Section == 12 && it.Index >= 112 && it.Index <= 117 && seed < 0:
				seed = position
			case it.Code == item.Code(12, 15) && chaos < 0:
				chaos = position
			case it.Code == item.Code(14, 22) && creation < 0:
				creation = position
			case it.Code == item.Code(12, 141) && it.Level == 0 && chaosBundle < 0:
				chaosBundle = position
			case it.Code == item.Code(12, 137) && it.Level == 0 && creationBundle < 0:
				creationBundle = position
			case it.Section == 14 && it.Index >= 245 && it.Index <= 247 && tetra < 0:
				tetra = position
			default:
				return
			}
		}
		if seed < 0 || sphere < 0 {
			return
		}
		high := p.chaosBox.Item(seed).KindB == item.KindBSocketSeedSphere
		if high {
			if chaosBundle < 0 || creationBundle < 0 || chaos >= 0 || creation >= 0 {
				return
			}
		} else if chaos < 0 || creation < 0 || chaosBundle >= 0 || creationBundle >= 0 {
			return
		}
		if len(items) != 4 && len(items) != 5 {
			return
		}
		if len(items) == 5 && tetra < 0 {
			return
		}
		seedItem := p.chaosBox.Item(seed)
		var index int
		var found bool
		if high {
			index, _, found = item.SocketManager.SeedSphereIndex(seedItem.Code, seedItem.Level)
		} else {
			index, found = item.SocketManager.SeedIndex(seedItem.Code, seedItem.Level)
		}
		if !found {
			return
		}
		output, found = item.SocketManager.SeedSphereItem(index, p.chaosBox.Item(sphere).Index-69)
		if !found {
			return
		}
		money := p.chaosBox.Item(sphere).BuyMoney
		if high {
			rate = 20 + money/200000
			if tetra >= 0 {
				rate = map[int]int{247: 50, 246: 65, 245: 80}[p.chaosBox.Item(tetra).Index] + money/200000
			}
		} else {
			money += seedItem.BuyMoney
			rate = 80 + money/200000
		}
	case item.SocketMixInsert: // Insert a Seed Sphere into an empty socket.
		if len(items) != 4 {
			return
		}
		seedSphere, chaos, creation := -1, -1, -1
		for _, position := range items {
			it := p.chaosBox.Item(position)
			switch {
			case it.Type == item.TypeSocket && targetPosition < 0:
				targetPosition = position
			case it.KindB == item.KindBSocketSeedSphere && seedSphere < 0:
				seedSphere = position
			case it.Code == item.Code(12, 15) && chaos < 0:
				chaos = position
			case it.Code == item.Code(14, 22) && creation < 0:
				creation = position
			default:
				return
			}
		}
		if targetPosition < 0 || seedSphere < 0 || chaos < 0 || creation < 0 {
			return
		}
		target := p.chaosBox.Item(targetPosition)
		slot, valid := target.SocketSlotAt(msg.Position)
		if !valid || slot != 1 {
			return
		}
		seedItem := p.chaosBox.Item(seedSphere)
		index, sphereLevel, found := item.SocketManager.SeedSphereIndex(seedItem.Code, seedItem.Level)
		if !found {
			return
		}
		element, _ := item.SocketManager.SeedElement(index)
		if target.Section <= 5 {
			if element != 1 && element != 3 && element != 5 {
				return
			}
		} else if target.Section <= 11 {
			if element != 2 && element != 4 && element != 6 {
				return
			}
		} else {
			return
		}
		for i := 0; i < 5; i++ {
			existing, _ := target.SocketSlotAt(i)
			if existing >= 100 && existing%100 == index {
				return
			}
		}
		output = target.Copy()
		output.SetSocketSlotAt(msg.Position, sphereLevel*100+index)
		item.SocketManager.RollBonus(output)
		output.Calc()
		if target.MaxDurability > 0 {
			output.Durability = output.MaxDurability * target.Durability / target.MaxDurability
		}
	case item.SocketMixRemove: // Remove one fitted Seed Sphere.
		if len(items) != 12 {
			return
		}
		chaos, harmony, guardian := 0, 0, 0
		for _, position := range items {
			it := p.chaosBox.Item(position)
			switch {
			case it.Type == item.TypeSocket && targetPosition < 0:
				targetPosition = position
			case it.Code == item.Code(12, 15):
				chaos++
			case it.Code == item.Code(14, 42):
				harmony++
			case it.Code == item.Code(14, 31):
				guardian++
			default:
				return
			}
		}
		if targetPosition < 0 || chaos != 5 || harmony != 5 || guardian != 1 {
			return
		}
		target := p.chaosBox.Item(targetPosition)
		slot, valid := target.SocketSlotAt(msg.Position)
		if !valid || slot < 100 {
			return
		}
		index := slot % 100
		if code, level, found := item.SocketManager.SeedItem(index); found {
			output = item.NewItem(code/512, code%512)
			output.Level, output.Durability = level, 1
			output.Calc()
		} else {
			return
		}
		rate, remove = 80, true
	}
	if p.Money < fee {
		reply.Result = 2
		return
	}
	if (msg.Type == item.SocketMixSeedExtract || msg.Type == item.SocketMixSeedSphereComposite || remove) && !p.chaosBox.OutputFits(output, targetPosition) {
		reply.Result = 6
		return
	}
	p.Money -= fee
	p.PushMoney()
	roll := rand.Intn(100)
	failed := roll > rate
	if remove {
		failed = roll >= rate
	}
	if failed {
		if msg.Type == item.SocketMixSeedExtract {
			for _, position := range items {
				it := p.chaosBox.Item(position)
				if it.IsExcellent() || it.Set != 0 {
					it.Level -= rand.Intn(3)
					if it.Level < 0 {
						it.Level = 0
					}
					it.Calc()
				} else {
					p.consumeChaosBoxItem(position)
				}
			}
		} else {
			for _, position := range items {
				if position != targetPosition {
					p.consumeChaosBoxItem(position)
				}
			}
		}
		p.pushChaosBoxItems()
		reply.Result = 0
		return
	}
	if msg.Type == item.SocketMixInsert || remove {
		for _, position := range items {
			if position != targetPosition {
				p.consumeChaosBoxItem(position)
			}
		}
		if remove {
			target := p.chaosBox.Item(targetPosition)
			oldMax, durability := target.MaxDurability, target.Durability
			target.SetSocketSlotAt(msg.Position, 1)
			item.SocketManager.RollBonus(target)
			target.Calc()
			if oldMax > 0 {
				target.Durability = target.MaxDurability * durability / oldMax
			}
			if !p.createMixInventoryItem(output) {
				return
			}
		} else {
			p.chaosBox.ReplaceItem(targetPosition, output)
		}
		if targetPosition != 0 {
			p.chaosBox.MoveItem(targetPosition, 0)
			targetPosition = 0
		}
		reply.Item = p.chaosBox.Item(targetPosition)
	} else {
		for _, position := range items {
			p.consumeChaosBoxItem(position)
		}
		if !p.createMixInventoryItem(output) {
			return
		}
		reply.Item = output
	}
	p.pushChaosBoxItems()
	reply.Result = 1
}

func (p *Player) consumeChaosBoxItem(position int) {
	p.chaosBox.RemoveItem(position)
}

func (p *Player) pushChaosBoxItems() {
	p.Push(&model.MsgTypeItemListReply{Type: 3, MsgItemListReply: model.MsgItemListReply{Items: p.chaosBox.Items()}})
}

func (p *Player) createMixInventoryItem(it *item.Item) bool {
	for position := 12; position < p.inventory.Size; position++ {
		if !p.inventory.CanPlace(position, it) {
			continue
		}
		p.inventory.AddItem(position, it)
		p.Push(&model.MsgPickItemReply{Result: position, Item: it})
		return true
	}
	return false
}

func (p *Player) CloseChaosBox(*model.MsgEmpty) {
	p.closeChaosBox()
	p.TargetNumber = -1
	p.Push(&model.MsgCloseChaosBoxReply{})
}

func (p *Player) closeChaosBox() {
	if !p.chaosBoxOpen {
		return
	}
	p.chaosBox.Clear()
	p.chaosBoxOpen = false
}

func (p *Player) equipmentChanged() {
	// 1, change skill
	newItemSkills := make(map[int]struct{})
	primaryHandWeapon := p.inventory.Items[0]
	if primaryHandWeapon != nil && primaryHandWeapon.SkillIndex != 0 {
		newItemSkills[primaryHandWeapon.SkillIndex] = struct{}{}
	}
	secondaryHandWeapon := p.inventory.Items[1]
	if secondaryHandWeapon != nil && secondaryHandWeapon.SkillIndex != 0 {
		newItemSkills[secondaryHandWeapon.SkillIndex] = struct{}{}
	}
	oldItemSkills := make(map[int]struct{})
	for _, s := range p.Skills {
		if s.Index < 300 && s.SkillBase.ItemSkill {
			oldItemSkills[s.Index] = struct{}{}
		}
	}
	var needLearnSkills []int
	for newSkill := range newItemSkills {
		if _, ok := oldItemSkills[newSkill]; !ok {
			needLearnSkills = append(needLearnSkills, newSkill)
		}
	}
	var needForgetSkills []int
	for oldSkill := range oldItemSkills {
		if _, ok := newItemSkills[oldSkill]; !ok {
			needForgetSkills = append(needForgetSkills, oldSkill)
		}
	}
	for _, index := range needLearnSkills {
		if s, ok := p.LearnSkill(index); ok {
			p.Push(&model.MsgSkillOneReply{
				Flag:  -2,
				Skill: s,
			})
		}
	}
	for _, index := range needForgetSkills {
		if s, ok := p.ForgetSkill(index); ok {
			p.Push(&model.MsgSkillOneReply{
				Flag:  -1,
				Skill: s,
			})
		}
	}
	// 2, calculate player
	p.calc()
}

func (p *Player) pushItemDurability(position, dur int) {
	p.Push(&model.MsgItemDurabilityReply{Position: position, Durability: dur, Flag: 1})
}

func (p *Player) pushDeleteItem(position int) {
	p.Push(&model.MsgDeleteInventoryItemReply{Position: position, Flag: 1})
}

func (p *Player) decreaseItemDurability(position int) {
	it := p.inventory.Items[position]
	it.Durability--
	if it.Durability > 0 {
		p.pushItemDurability(position, it.Durability)
	} else {
		p.inventory.RemoveItem(position, it)
		p.pushDeleteItem(position)
	}
}

func (p *Player) UseItem(msg *model.MsgUseItem) {
	// validate the position
	if msg.SrcPosition < 0 || msg.SrcPosition >= p.inventory.Size ||
		msg.DstPosition < 0 || msg.DstPosition >= p.inventory.Size ||
		msg.SrcPosition == msg.DstPosition {
		return
	}
	it := p.inventory.Items[msg.SrcPosition]
	if it == nil {
		return
	}
	switch {
	case it.Code == item.Code(12, 30): // Bundle Jewel of Bless
	case it.Code == item.Code(13, 15): // Fruits 果实
	case it.Code >= item.Code(13, 43) && it.Code <= item.Code(13, 45): // Seal 印章
	case it.Code == item.Code(13, 48): // Kalima Ticket 卡利玛自由入场券
	case it.Code >= item.Code(13, 54) && it.Code <= item.Code(13, 58): // Reset Fruit 洗点果实
	case it.Code == item.Code(13, 60): // Indulgence 免罪符
	case it.Code == item.Code(13, 66): // Invitation to Santa Village 圣诞之地入场券
	case it.Code == item.Code(13, 69): // Talisman of Resurrection 复活符咒
	case it.Code == item.Code(13, 70): // Talisman of Mobility 移动符咒
	case it.Code == item.Code(13, 82): // Talisman of Item Protection 装备保护符咒
	case it.Code >= item.Code(13, 152) && it.Code <= item.Code(13, 159): // Scroll of Oblivion 忘却卷轴
	case it.Code >= item.Code(14, 0) && it.Code <= item.Code(14, 3): // HP Potion
		addRate := 0
		switch it.Code {
		case item.Code(14, 0): // Apple 苹果
			addRate = 10
		case item.Code(14, 1): // Small Healing Potion 小瓶治疗药水
			addRate = 20
		case item.Code(14, 2): // Healing Potion 中瓶治疗药水
			addRate = 30
		case item.Code(14, 3): // Large Healing Potion 大瓶治疗药水
			addRate = 40
		}
		if it.Level >= 1 {
			addRate += 5
		}
		hp := 0
		hp += p.MaxHP * addRate / 100
		// defer recover hp
		p.SetDelayRecoverHP(hp, hp)
		// decrease durability
		p.decreaseItemDurability(msg.SrcPosition)
	case it.Code >= item.Code(14, 4) && it.Code <= item.Code(14, 6): // MP Potion
		addRate := 0
		switch it.Code {
		case item.Code(14, 4):
			addRate = 20
		case item.Code(14, 5):
			addRate = 30
		case item.Code(14, 6):
			addRate = 40
		}
		mp := p.MaxMP * addRate / 100
		// recover mp immediately
		if p.MP < p.MaxMP {
			p.MP += mp
			if p.MP > p.MaxMP {
				p.MP = p.MaxMP
			}
			p.PushMPAG(p.MP, p.AG)
		}
		// decrease durability
		p.decreaseItemDurability(msg.SrcPosition)
	case it.Code == item.Code(14, 7): // Siege Potion 攻城药水
	case it.Code == item.Code(14, 8): // Antidote 解毒剂
	case it.Code == item.Code(14, 9) || it.Code == item.Code(14, 20): // Ale 酒 / Remedy of Love 爱情的魔力
	case it.Code == item.Code(14, 10): // Town Portal Scroll 回城卷轴
	case it.Code == item.Code(14, 13): // Jewel of Bless
	case it.Code == item.Code(14, 14): // Jewel of Soul
	case it.Code == item.Code(14, 16): // Jewel of Life
	case it.Code >= item.Code(14, 35) && it.Code <= item.Code(14, 37): // SD Potion
		addRate := 0
		switch it.Code {
		case item.Code(14, 35):
			addRate = 25
		case item.Code(14, 36):
			addRate = 35
		case item.Code(14, 37):
			addRate = 45
		}
		sd := p.MaxSD * addRate / 100
		// defer recover sd
		p.SetDelayRecoverSD(sd, sd)
		// decrease durability
		p.decreaseItemDurability(msg.SrcPosition)
	case it.Code >= item.Code(14, 38) && it.Code <= item.Code(14, 40): // comples/compound Potion
		addHPRate, addSDRate := 0, 0
		switch it.Code {
		case item.Code(14, 38):
			addHPRate = 10
			addSDRate = 5
		case item.Code(14, 39):
			addHPRate = 25
			addSDRate = 10
		case item.Code(14, 40):
			addHPRate = 45
			addSDRate = 20
		}
		hp := p.MaxHP * addHPRate / 100
		sd := p.MaxSD * addSDRate / 100
		// defer recover hp sd
		p.SetDelayRecoverHP(hp, hp)
		p.SetDelayRecoverSD(sd, sd)
		// decrease durability
		p.decreaseItemDurability(msg.SrcPosition)
	// case it.Code == item.Code(14, 42) && it2.Type != item.TypeSocket: // 再生强化
	case it.Code >= item.Code(14, 43) && it.Code <= item.Code(14, 44): // 进化道具
	case it.Code >= item.Code(14, 46) && it.Code <= item.Code(14, 50): // Jack O'Lantern 南瓜灯饮料
	case it.Code == item.Code(14, 70): // Elite HP Potion 精华HP药水
	case it.Code == item.Code(14, 71): // Elite MP Potion 精华MP药水
	case it.Code >= item.Code(14, 78) && it.Code <= item.Code(14, 82): // kindBPremiumElixir 会员圣水
	case it.Code >= item.Code(14, 85) && it.Code <= item.Code(14, 87): // Cherry Blossom 樱花
	case it.Code == item.Code(14, 133): // Elite SD Potion 精华防护值药水
	case it.Code == item.Code(14, 160): // Jewel of Extension 延长宝石
	case it.Code == item.Code(14, 161): // Jewel of Elevation 提高宝石
	case it.Code == item.Code(14, 162): // Magic Backpack 魔法背书
	case it.Code == item.Code(14, 163): // Vault Expansion Certificate 仓库拓展证书
	case it.Code == item.Code(14, 209): // Tradeable Seal 交易印章
	case it.Code == item.Code(14, 224): // Bless of Light (Greater) 光的祝福
	case it.Code >= item.Code(14, 263) && it.Code <= item.Code(14, 264): // Bless of Light 光之祝福
	case it.KindA == item.KindASkill:
		if !p.canUseItem(it) {
			return
		}
		// (15, 18) // Scroll of Nova 星辰一怒术
		skillIndex := it.SkillIndex
		if it.Code == item.Code(12, 11) { // Orb of Summoning 召唤之石
			skillIndex += it.Level
		}
		if s, ok := p.LearnSkill(skillIndex); ok {
			p.Push(&model.MsgSkillOneReply{
				Flag:  -2,
				Skill: s,
			})
			p.inventory.RemoveItem(msg.SrcPosition, it)
			p.Push(&model.MsgDeleteInventoryItemReply{
				Position: msg.SrcPosition,
				Flag:     1,
			})
		}
	}
}

func (p *Player) canUseItem(it *item.Item) bool {
	if p.Level < it.ReqLevel ||
		p.strength+p.addStrength < it.ReqStrength ||
		p.dexterity+p.addDexterity < it.ReqDexterity ||
		p.vitality+p.addVitality < it.ReqVitality ||
		p.energy+p.addEnergy < it.ReqEnergy ||
		p.leadership+p.addLeadership < it.ReqCommand {
		return false
	}
	reqClass := it.ReqClass[p.Class]
	if reqClass == 0 || (reqClass > p.changeUp+1) {
		return false
	}
	return true
}

func (p *Player) ValidateItemPosition(position int) bool {
	ok := position >= 0 && position < p.inventory.Size
	if !ok {
		slog.Warn("ValidateItemPosition", "position", position, "object", p.Name)
	}
	return ok
}

func (p *Player) repairItem(it *item.Item, fast bool) bool {
	reply := model.MsgRepairItemReply{}
	defer p.Push(&reply)
	money := it.CalculateRepairMoney(fast)
	if money > p.Money {
		return false
	}
	p.Money -= money
	reply.Money = p.Money
	it.Durability = it.MaxDurability
	p.Push(&model.MsgItemDurabilityReply{
		Position:   it.Position,
		Durability: it.Durability,
		Flag:       0,
	})
	return true
}

func (p *Player) RepairItem(msg *model.MsgRepairItem) {
	if msg.Position == 0xFF {
		for _, it := range p.inventory.Items {
			if it == nil {
				continue
			}
			if !it.Repair {
				continue
			}
			if ok := p.repairItem(it, msg.Fast); !ok {
				return
			}
		}
		return
	}
	// validate
	if !p.ValidateItemPosition(msg.Position) {
		return
	}
	it := p.inventory.Item(msg.Position)
	if !it.Repair {
		return
	}
	p.repairItem(it, msg.Fast)
}

func (p *Player) UsePet(msg *model.MsgUsePet) {
	reply := model.MsgUsePetReply{
		Position: msg.Position,
		Result:   0,
	}
	defer p.Push(&reply)

	if msg.Position < 12 || msg.Position > 204 {
		return
	}
	it := p.inventory.Items[msg.Position]
	if it == nil || it.Durability <= 0 {
		return
	}
	switch it.Code {
	case item.Code(13, 2), // Horn of Uniria 兽角
		item.Code(13, 3),  // Horn of Dinorant 彩云兽
		item.Code(13, 4),  // Dark Horse 黑王马之角
		item.Code(13, 37): // Horn of Fenrir 炎狼兽之角
		if p.pet != nil {
			it.HarmonyOption = 0
			p.pet = nil
			reply.Result = -1
		} else {
			it.HarmonyOption = 1
			p.pet = it
			reply.Result = -2
		}
		p.equipmentChanged()
		p.pushChangedEquipment(it)
	}
}

func (p *Player) pushChangedEquipment(it *item.Item) {
	reply := model.MsgChangeEquipmentReply{
		Number: p.Index,
		Item:   it,
	}
	p.PushViewport(&reply)
}
