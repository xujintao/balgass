package drop

import (
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/item"
)

func TestActivityMonsterBags(t *testing.T) {
	for _, id := range []int{153, 154, 155, 156} {
		if _, exists := DropManager.itemBags[itemBagKey{kind: itemBagEvent, id: id}]; exists {
			t.Fatalf("activity EventBag %d is still registered", id)
		}
	}
	for _, class := range []int{135, 136, 137, 466, 476} {
		if bag := DropManager.itemBag(itemBagKey{kind: itemBagMonster, id: class}); bag == nil {
			t.Fatalf("monster %d has no ItemBag", class)
		}
	}
	soldier := DropManager.itemBags[itemBagKey{kind: itemBagMonster, id: 136}]
	archer := DropManager.itemBags[itemBagKey{kind: itemBagMonster, id: 137}]
	if soldier.name == archer.name {
		t.Fatal("soldier and archer share an ItemBag file")
	}
	request := Request{Trigger: TriggerMonsterDeath, MonsterClass: 135, Player: PlayerContext{Level: 200}}
	result := DropManager.Drop(request)
	if !result.Handled || len(result.Rewards) != 1 || result.Rewards[0].Item.Code != item.Code(14, 13) ||
		result.Rewards[0].Item.Durability != 0 {
		t.Fatal("White Wizard boss reward is invalid")
	}
	for _, class := range []int{136, 137} {
		request.MonsterClass = class
		for i := 0; i < 100; i++ {
			result = DropManager.Drop(request)
			if !result.Handled || len(result.Rewards) != 1 ||
				(result.Rewards[0].Zen != 10000 && (result.Rewards[0].Item == nil ||
					result.Rewards[0].Item.Code != item.Code(13, 20) || result.Rewards[0].Item.Durability != 30)) {
				t.Fatalf("White Wizard guard %d reward is invalid", class)
			}
		}
	}
	if bag := DropManager.itemBag(itemBagKey{kind: itemBagMonster, id: 466}); bag == nil || bag.itemRate != 1000 {
		t.Fatalf("Santa goblin ItemBag probability = %v, want 1000", bag)
	}
	request.MonsterClass = 466
	for i := 0; i < 100; i++ {
		result = DropManager.Drop(request)
		if len(result.Rewards) != 1 ||
			(result.Rewards[0].Zen != 10000 && (result.Rewards[0].Item == nil ||
				result.Rewards[0].Item.Code != item.Code(13, 66) || result.Rewards[0].Item.Durability != 30)) {
			t.Fatal("Santa goblin reward is invalid")
		}
	}
	request.MonsterClass = 476
	result = DropManager.Drop(request)
	if len(result.Rewards) != 1 || result.Rewards[0].Item.Code != item.Code(13, 66) {
		t.Fatal("Santa boss reward is invalid")
	}
}
