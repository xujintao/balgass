package object

import (
	"math"

	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/shop"
)

func (obj *Object) PickItem(msg *model.MsgPickItem) {
	reply := model.MsgPickItemReply{
		Result: -1,
	}
	itemDurChanged := false
	position := -1
	var it *item.Item
	defer func() {
		obj.Push(&reply)
		if itemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   position,
				Durability: it.Durability,
				Flag:       0,
			}
			obj.Push(&reply)
		}
	}()
	it2 := maps.MapManager.PeekItem(obj.MapNumber, msg.Index)
	if it2 == nil {
		return
	}
	switch it2.Code {
	case item.Code(14, 15): // zen
		money := it2.Durability
		objMoney := obj.GetMoney()
		if objMoney+money > MaxZen {
			return
		}
		objMoney += money
		obj.SetMoney(objMoney)
		reply.Money = objMoney
		position = -2
	default:
		position = obj.GetInventory().FindFreePositionForItem(it2)
		if position == -1 {
			return
		}
		it = obj.GetInventory().Items[position]
		if it == nil {
			obj.GetInventory().AddItem(position, it2)
		} else {
			it.Durability += it2.Durability
			itemDurChanged = true
		}
	}
	maps.MapManager.RemoveItem(obj.MapNumber, msg.Index)
	reply.Result = position
	reply.Item = it2
}

func (obj *Object) DropItem(msg *model.MsgDropItem) {
	reply := model.MsgDropItemReply{
		Result:   0,
		Position: msg.Position,
	}
	defer obj.Push(&reply)
	// validate
	if msg.Position >= obj.GetInventory().Size {
		return
	}
	it := obj.GetInventory().Items[msg.Position]
	if it == nil {
		return
	}
	switch it.Code {
	case item.Code(14, 63):
		cmdReply := model.MsgServerCMDReply{
			Type: 0,
			X:    obj.X,
			Y:    obj.Y,
		}
		obj.PushViewport(&cmdReply)
	default:
		ok := maps.MapManager.AddItem(obj.MapNumber, msg.X, msg.Y, it)
		if !ok {
			return
		}
	}
	obj.GetInventory().RemoveItem(msg.Position, it)
	if msg.Position < 12 || msg.Position == 126 {
		obj.EquipmentChanged()
	}
	reply.Result = 1
}

func (obj *Object) BuyItem(msg *model.MsgBuyItem) {
	reply := model.MsgBuyItemReply{
		Result: -1,
	}
	itemDurChanged := false
	position := -1
	var it *item.Item
	defer func() {
		obj.Push(&reply)
		if itemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   position,
				Durability: it.Durability,
				Flag:       0,
			}
			obj.Push(&reply)
		}
	}()
	// validate
	if msg.Position < 0 ||
		msg.Position >= shop.MaxShopItemCount ||
		obj.TargetNumber < 0 {
		return
	}
	tobj := ObjectManager.GetObject(obj.TargetNumber)
	if tobj == nil {
		return
	}
	if math.Abs(float64(obj.X-tobj.X)) > 5 ||
		math.Abs(float64(obj.Y-tobj.Y)) > 5 {
		return
	}
	if tobj.NpcType != NpcTypeShop {
		return
	}
	sit := shop.ShopManager.GetShopItem(tobj.Class, tobj.MapNumber, msg.Position)
	if sit == nil {
		return
	}
	position = obj.GetInventory().FindFreePositionForItem(sit)
	if position == -1 {
		return
	}
	it = obj.GetInventory().Items[position]
	if it == nil {
		obj.GetInventory().AddItem(position, sit)
	} else {
		it.Durability += sit.Durability
		itemDurChanged = true
	}
	reply.Result = position
	reply.Item = sit
}

func (obj *Object) SellItem(msg *model.MsgSellItem) {
	reply := model.MsgSellItemReply{
		Result: 0,
	}
	defer obj.Push(&reply)
	// validate
	if msg.Position < 0 ||
		msg.Position >= obj.GetInventory().Size ||
		obj.TargetNumber < 0 {
		return
	}
	tobj := ObjectManager.GetObject(obj.TargetNumber)
	if tobj == nil {
		return
	}
	if math.Abs(float64(obj.X-tobj.X)) > 5 ||
		math.Abs(float64(obj.Y-tobj.Y)) > 5 {
		return
	}
	if tobj.NpcType != NpcTypeShop {
		return
	}
	it := obj.GetInventory().Items[msg.Position]
	if it == nil {
		return
	}
	objMoney := obj.GetMoney()
	if objMoney+it.Money > MaxZen {
		return
	}
	objMoney += it.Money
	obj.SetMoney(objMoney)
	obj.GetInventory().RemoveItem(msg.Position, it)
	reply.Result = 1
	reply.Money = objMoney
}

func (obj *Object) MoveItem(msg *model.MsgMoveItem) {
	reply := model.MsgMoveItemReply{
		Result: -1,
	}
	sitemDurChanged := false
	titemDurChanged := false
	var sitem *item.Item
	var titem *item.Item
	defer func() {
		obj.Push(&reply)
		if sitemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   msg.SrcPosition,
				Durability: sitem.Durability,
				Flag:       0,
			}
			obj.Push(&reply)
		}
		if titemDurChanged {
			reply := model.MsgItemDurabilityReply{
				Position:   msg.DstPosition,
				Durability: titem.Durability,
				Flag:       0,
			}
			obj.Push(&reply)
		}
	}()

	// get source item
	switch msg.SrcFlag {
	case 0: // inventory
		if msg.SrcPosition >= obj.GetInventory().Size {
			return
		}
		sitem = obj.GetInventory().Items[msg.SrcPosition]
	case 2: // warehouse
		if msg.SrcPosition >= obj.GetWarehouse().Size {
			return
		}
		sitem = obj.GetWarehouse().Items[msg.SrcPosition]

	}
	if sitem == nil {
		return
	}
	// get destination item
	switch msg.DstFlag {
	case 0: // inventory
		if msg.DstPosition >= obj.GetInventory().Size {
			return
		}
		titem = obj.GetInventory().Items[msg.DstPosition]
	case 2: // warehouse
		if msg.DstPosition >= obj.GetWarehouse().Size {
			return
		}
		titem = obj.GetWarehouse().Items[msg.DstPosition]
	}

	switch msg.SrcFlag {
	case 0:
		switch msg.DstFlag {
		case 0:
			switch {
			case titem == nil: // move
				ok := obj.GetInventory().CheckFlagsForItem(msg.DstPosition, sitem)
				if !ok {
					return
				}
				obj.GetInventory().RemoveItem(msg.SrcPosition, sitem)
				obj.GetInventory().AddItem(msg.DstPosition, sitem)
				if msg.SrcPosition < 12 || msg.SrcPosition == 236 ||
					msg.DstPosition < 12 || msg.DstPosition == 236 {
					obj.EquipmentChanged()
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
					obj.GetInventory().RemoveItem(msg.SrcPosition, sitem)
					reply := model.MsgDeleteInventoryItemReply{
						Position: msg.SrcPosition,
						Flag:     1,
					}
					obj.Push(&reply)
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
				ok := obj.GetWarehouse().CheckFlagsForItem(msg.DstPosition, sitem)
				if !ok {
					return
				}
				obj.GetInventory().RemoveItem(msg.SrcPosition, sitem)
				obj.GetWarehouse().AddItem(msg.DstPosition, sitem)
				if msg.SrcPosition < 12 || msg.SrcPosition == 236 {
					obj.EquipmentChanged()
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
				ok := obj.GetInventory().CheckFlagsForItem(msg.DstPosition, sitem)
				if !ok {
					return
				}
				obj.GetWarehouse().RemoveItem(msg.SrcPosition, sitem)
				obj.GetInventory().AddItem(msg.DstPosition, sitem)
				if msg.DstPosition < 12 || msg.DstPosition == 236 {
					obj.EquipmentChanged()
				}
				reply.Result = msg.DstFlag
			}
		case 2:
			if titem == nil {
				ok := obj.GetWarehouse().CheckFlagsForItem(msg.DstPosition, sitem)
				if !ok {
					return
				}
				obj.GetWarehouse().RemoveItem(msg.SrcPosition, sitem)
				obj.GetWarehouse().AddItem(msg.DstPosition, sitem)
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

func (obj *Object) pushItemDurability(position, dur int) {
	obj.Push(&model.MsgItemDurabilityReply{Position: position, Durability: dur, Flag: 1})
}

func (obj *Object) pushDeleteItem(position int) {
	obj.Push(&model.MsgDeleteInventoryItemReply{Position: position, Flag: 1})
}

func (obj *Object) decreaseItemDurability(position int) {
	it := obj.GetInventory().Items[position]
	it.Durability--
	if it.Durability > 0 {
		obj.pushItemDurability(position, it.Durability)
	} else {
		obj.GetInventory().RemoveItem(position, it)
		obj.pushDeleteItem(position)
	}
}

func (obj *Object) UseItem(msg *model.MsgUseItem) {
	// validate the position
	if msg.SrcPosition < 0 || msg.SrcPosition >= obj.GetInventory().Size ||
		msg.DstPosition < 0 || msg.DstPosition >= obj.GetInventory().Size ||
		msg.SrcPosition == msg.DstPosition {
		return
	}
	it := obj.GetInventory().Items[msg.SrcPosition]
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
		hp += obj.MaxHP * addRate / 100
		// defer recover hp
		obj.SetDelayRecoverHP(hp, hp)
		// decrease durability
		obj.decreaseItemDurability(msg.SrcPosition)
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
		mp := obj.MaxMP * addRate / 100
		// recover mp immediately
		if obj.MP < obj.MaxMP {
			obj.MP += mp
			if obj.MP > obj.MaxMP {
				obj.MP = obj.MaxMP
			}
			obj.PushMPAG(obj.MP, obj.AG)
		}
		// decrease durability
		obj.decreaseItemDurability(msg.SrcPosition)
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
		sd := obj.MaxSD * addRate / 100
		// defer recover sd
		obj.SetDelayRecoverSD(sd, sd)
		// decrease durability
		obj.decreaseItemDurability(msg.SrcPosition)
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
		hp := obj.MaxHP * addHPRate / 100
		sd := obj.MaxSD * addSDRate / 100
		// defer recover hp sd
		obj.SetDelayRecoverHP(hp, hp)
		obj.SetDelayRecoverSD(sd, sd)
		// decrease durability
		obj.decreaseItemDurability(msg.SrcPosition)
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
		// (15, 18) // Scroll of Nova 星辰一怒术
		skillIndex := it.SkillIndex
		if it.Code == item.Code(12, 11) { // Orb of Summoning 召唤之石
			skillIndex += it.Level
		}
		if s, ok := obj.LearnSkill(skillIndex); ok {
			obj.Push(&model.MsgSkillOneReply{
				Flag:  -2,
				Skill: s,
			})
			obj.GetInventory().RemoveItem(msg.SrcPosition, it)
			obj.Push(&model.MsgDeleteInventoryItemReply{
				Position: msg.SrcPosition,
				Flag:     1,
			})
		}
	}
}
