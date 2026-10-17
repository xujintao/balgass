package drop

import "testing"

func TestInvasionMonsterBags(t *testing.T) {
	for class, name := range map[int]string{
		44: "Monster_(44)_Dragon_Red",
		55: "Monster_(55)_Death_King",
	} {
		entry := DropManager.itemBags[itemBagKey{kind: itemBagMonster, id: class}]
		if entry == nil || entry.name != name || entry.bag == nil {
			t.Fatalf("monster %d did not load its MonsterBag", class)
		}
		result := DropManager.Drop(testRequest(class))
		if !result.Handled || len(result.Rewards) == 0 {
			t.Fatalf("monster %d did not drop through its MonsterBag: %#v", class, result)
		}
	}
}

func TestWhiteWizardMonsterBags(t *testing.T) {
	for class, name := range map[int]string{
		135: "Monster_(135)_White_Wizard",
		136: "Monster_(136)_Destructive_Ogre_Soldier",
		137: "Monster_(137)_Destructive_Ogre_Archer",
	} {
		entry := DropManager.itemBags[itemBagKey{kind: itemBagMonster, id: class}]
		if entry == nil || entry.name != name || entry.bag == nil {
			t.Fatalf("monster %d did not load its MonsterBag", class)
		}
		result := DropManager.Drop(testRequest(class))
		if !result.Handled || len(result.Rewards) == 0 {
			t.Fatalf("monster %d did not drop through its MonsterBag: %#v", class, result)
		}
	}
	soldier := DropManager.itemBags[itemBagKey{kind: itemBagMonster, id: 136}]
	archer := DropManager.itemBags[itemBagKey{kind: itemBagMonster, id: 137}]
	if soldier.name == archer.name {
		t.Fatal("soldier and archer share an ItemBag file")
	}
}

func TestSantaMonsterBags(t *testing.T) {
	for class, name := range map[int]string{
		466: "Monster_(466)_Cursed_Goblin",
		476: "Monster_(476)_Cursed_Santa",
	} {
		entry := DropManager.itemBags[itemBagKey{kind: itemBagMonster, id: class}]
		if entry == nil || entry.name != name || entry.bag == nil {
			t.Fatalf("monster %d did not load its MonsterBag", class)
		}
		result := DropManager.Drop(testRequest(class))
		if !result.Handled || len(result.Rewards) == 0 {
			t.Fatalf("monster %d did not drop through its MonsterBag: %#v", class, result)
		}
	}
}

func TestMedusaGroupMonsterBags(t *testing.T) {
	for class, name := range map[int]string{
		560: "Monster_(560)_Sapi_Queen",
		561: "Monster_(561)_Medusa",
	} {
		entry := DropManager.itemBags[itemBagKey{kind: itemBagMonster, id: class}]
		if entry == nil || entry.name != name || entry.bag == nil {
			t.Fatalf("monster %d did not load its existing ItemBag", class)
		}
		result := DropManager.Drop(testRequest(class))
		if !result.Handled || len(result.Rewards) == 0 {
			t.Fatalf("monster %d did not drop through its MonsterBag: %#v", class, result)
		}
	}
}

func TestKundunAndSilvesterMonsterBags(t *testing.T) {
	for class, name := range map[int]string{
		275: "Monster_(275)_Kundun",
		673: "Monster_(673)_Lord_Silvester",
	} {
		entry := DropManager.itemBags[itemBagKey{kind: itemBagMonster, id: class}]
		if entry == nil || entry.name != name || entry.bag == nil {
			t.Fatalf("monster %d did not load its MonsterBag", class)
		}
		result := DropManager.Drop(testRequest(class))
		if !result.Handled || len(result.Rewards) == 0 {
			t.Fatalf("monster %d did not drop through its MonsterBag: %#v", class, result)
		}
	}
}
