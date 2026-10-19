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
		position = p.GetInventory().FindFreePositionForItem(it2)
		if position == -1 {
			return
		}
		it = p.GetInventory().Items[position]
		if it == nil {
			p.GetInventory().AddItem(position, it2)
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
	if msg.Position < 0 || msg.Position >= p.GetInventory().Size {
		return
	}
	it := p.GetInventory().Items[msg.Position]
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
			MapNumber: p.MapNumber,
			X:         p.X,
			Y:         p.Y,
			Player: drop.PlayerContext{
				Class:       p.Class,
				ChangeUp:    p.GetChangeUp(),
				Level:       p.Level,
				MasterLevel: p.GetMasterLevel(),
				MapNumber:   p.MapNumber,
			},
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
	p.GetInventory().RemoveItem(msg.Position, it)
	if msg.Position < 12 || msg.Position == 126 {
		p.EquipmentChanged()
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
	position = p.GetInventory().FindFreePositionForItem(sit)
	if position == -1 {
		return
	}
	if p.Money < sit.BuyMoney {
		return
	}
	p.Money -= sit.BuyMoney
	p.PushMoney()
	it = p.GetInventory().Items[position]
	if it == nil {
		p.GetInventory().AddItem(position, sit)
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
		msg.Position >= p.GetInventory().Size ||
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
	it := p.GetInventory().Items[msg.Position]
	if it == nil {
		return
	}
	if p.Money+it.SellMoney > object.MaxZen {
		return
	}
	p.Money += it.SellMoney
	p.GetInventory().RemoveItem(msg.Position, it)
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
		if msg.SrcPosition >= p.GetInventory().Size {
			return
		}
		sitem = p.GetInventory().Items[msg.SrcPosition]
	case 2: // warehouse
		if msg.SrcPosition >= p.GetWarehouse().Size {
			return
		}
		sitem = p.GetWarehouse().Items[msg.SrcPosition]

	}
	if sitem == nil {
		return
	}
	// get destination item
	switch msg.DstFlag {
	case 0: // inventory
		if msg.DstPosition >= p.GetInventory().Size {
			return
		}
		titem = p.GetInventory().Items[msg.DstPosition]
	case 2: // warehouse
		if msg.DstPosition >= p.GetWarehouse().Size {
			return
		}
		titem = p.GetWarehouse().Items[msg.DstPosition]
	}

	switch msg.SrcFlag {
	case 0:
		switch msg.DstFlag {
		case 0:
			switch {
			case titem == nil: // move
				ok := p.GetInventory().CheckFlagsForItem(msg.DstPosition, sitem)
				if !ok {
					return
				}
				p.GetInventory().RemoveItem(msg.SrcPosition, sitem)
				p.GetInventory().AddItem(msg.DstPosition, sitem)
				if msg.SrcPosition < 12 || msg.SrcPosition == 236 ||
					msg.DstPosition < 12 || msg.DstPosition == 236 {
					p.EquipmentChanged()
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
					p.GetInventory().RemoveItem(msg.SrcPosition, sitem)
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
		case 2:
			if titem == nil {
				ok := p.GetWarehouse().CheckFlagsForItem(msg.DstPosition, sitem)
				if !ok {
					return
				}
				p.GetInventory().RemoveItem(msg.SrcPosition, sitem)
				p.GetWarehouse().AddItem(msg.DstPosition, sitem)
				if msg.SrcPosition < 12 || msg.SrcPosition == 236 {
					p.EquipmentChanged()
				}
				reply.Result = msg.DstFlag
			}
		default:
			return
		}
	case 2:
		switch msg.DstFlag {
		case 0:
			if titem == nil {
				ok := p.GetInventory().CheckFlagsForItem(msg.DstPosition, sitem)
				if !ok {
					return
				}
				p.GetWarehouse().RemoveItem(msg.SrcPosition, sitem)
				p.GetInventory().AddItem(msg.DstPosition, sitem)
				if msg.DstPosition < 12 || msg.DstPosition == 236 {
					p.EquipmentChanged()
				}
				reply.Result = msg.DstFlag
			}
		case 2:
			if titem == nil {
				ok := p.GetWarehouse().CheckFlagsForItem(msg.DstPosition, sitem)
				if !ok {
					return
				}
				p.GetWarehouse().RemoveItem(msg.SrcPosition, sitem)
				p.GetWarehouse().AddItem(msg.DstPosition, sitem)
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

func (p *Player) pushItemDurability(position, dur int) {
	p.Push(&model.MsgItemDurabilityReply{Position: position, Durability: dur, Flag: 1})
}

func (p *Player) pushDeleteItem(position int) {
	p.Push(&model.MsgDeleteInventoryItemReply{Position: position, Flag: 1})
}

func (p *Player) decreaseItemDurability(position int) {
	it := p.GetInventory().Items[position]
	it.Durability--
	if it.Durability > 0 {
		p.pushItemDurability(position, it.Durability)
	} else {
		p.GetInventory().RemoveItem(position, it)
		p.pushDeleteItem(position)
	}
}

func (p *Player) UseItem(msg *model.MsgUseItem) {
	// validate the position
	if msg.SrcPosition < 0 || msg.SrcPosition >= p.GetInventory().Size ||
		msg.DstPosition < 0 || msg.DstPosition >= p.GetInventory().Size ||
		msg.SrcPosition == msg.DstPosition {
		return
	}
	it := p.GetInventory().Items[msg.SrcPosition]
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
		if !p.CanUseItem(it) {
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
			p.GetInventory().RemoveItem(msg.SrcPosition, it)
			p.Push(&model.MsgDeleteInventoryItemReply{
				Position: msg.SrcPosition,
				Flag:     1,
			})
		}
	}
}

func (p *Player) ValidateItemPosition(position int) bool {
	ok := position >= 0 && position < p.GetInventory().Size
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
		for _, it := range p.GetInventory().Items {
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
	it := p.GetInventoryItem(msg.Position)
	if !it.Repair {
		return
	}
	p.repairItem(it, msg.Fast)
}
