package player

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/class"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

func TestAttackProcRatesIncludeMasterSkills(t *testing.T) {
	p := Player{
		recoverMaxHP: 5,
		recoverMaxMP: 5,
	}
	p.Skills = skill.Skills{
		skill.SkillIndexMaceMastery:      {Level: 10},
		skill.SkillIndexRecoverManaFully: {Level: 10},
		skill.SkillIndexRecoverHPFully:   {Level: 10},
		skill.SkillIndexRecoverSDFully:   {Level: 10},
	}

	if p.GetMaceStunRate() <= 0 {
		t.Fatal("mace stun rate did not include master skill")
	}
	if p.GetFullHPRecoveryRate() <= 5 || p.GetFullMPRecoveryRate() <= 5 {
		t.Fatal("HP/MP recovery rates did not combine wing and master skill values")
	}
	if p.GetFullSDRecoveryRate() <= 0 {
		t.Fatal("SD recovery rate did not include master skill")
	}
}

func TestCanUseItemChecksRequirements(t *testing.T) {
	p := Player{}
	p.Class = int(class.Wizard)
	p.Level = 20
	p.energy = 30
	p.changeUp = 0
	it := &item.Item{ItemBase: &item.ItemBase{
		ReqLevel:  20,
		ReqEnergy: 30,
		ReqClass:  [8]int{class.Wizard: 1},
	}}
	if !p.CanUseItem(it) {
		t.Fatal("CanUseItem() = false, want true")
	}
	it.ReqEnergy = 31
	if p.CanUseItem(it) {
		t.Fatal("CanUseItem() = true with insufficient energy")
	}
	it.ReqEnergy = 30
	it.ReqClass[class.Wizard] = 2
	if p.CanUseItem(it) {
		t.Fatal("CanUseItem() = true without required change-up")
	}
}

func TestAutoSaveSchedule(t *testing.T) {
	now := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name             string
		player           Player
		wantAutoSaveTime time.Time
	}{
		{
			name: "not playing",
			player: Player{
				Object:       object.Object{Name: "hero", ConnectState: object.ConnectStateLogged},
				autoSaveTime: now.Add(-autoSaveInterval),
			},
			wantAutoSaveTime: now.Add(-autoSaveInterval),
		},
		{
			name: "empty character name",
			player: Player{
				Object:       object.Object{ConnectState: object.ConnectStatePlaying},
				autoSaveTime: now.Add(-autoSaveInterval),
			},
			wantAutoSaveTime: now.Add(-autoSaveInterval),
		},
		{
			name: "not due",
			player: Player{
				Object:       object.Object{Name: "hero", ConnectState: object.ConnectStatePlaying},
				autoSaveTime: now.Add(-autoSaveInterval + time.Second),
			},
			wantAutoSaveTime: now.Add(-autoSaveInterval + time.Second),
		},
		{
			name: "initializes missing save time",
			player: Player{
				Object: object.Object{Name: "hero", ConnectState: object.ConnectStatePlaying},
			},
			wantAutoSaveTime: now,
		},
		{
			name: "refreshes due save time",
			player: Player{
				Object:       object.Object{Name: "hero", ConnectState: object.ConnectStatePlaying},
				autoSaveTime: now.Add(-autoSaveInterval),
			},
			wantAutoSaveTime: now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.player.autoSave(now)
			if !tt.player.autoSaveTime.Equal(tt.wantAutoSaveTime) {
				t.Fatalf("autoSaveTime = %v, want %v", tt.player.autoSaveTime, tt.wantAutoSaveTime)
			}
		})
	}
}

func TestLoadWarehouseMarksWarehouseOpen(t *testing.T) {
	p := Player{}
	p.loadWarehouse(&model.Account{
		WarehouseMoney: 321,
	})

	if !p.warehouseOpen {
		t.Fatal("warehouseOpen = false, want true")
	}
	if p.warehouseMoney != 321 {
		t.Fatalf("warehouseMoney = %d, want 321", p.warehouseMoney)
	}
}

func TestCharacterTableInitLoadsXML(t *testing.T) {
	pathCommon := conf.Server.PathCommon
	dir := filepath.Join(t.TempDir(), "class")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "class.xml"), []byte(`<ClassTable>
		<Class Annotation="darkload" Class="4" Level="1" MapNumber="0" Strength="26" Dexterity="20" Vitality="20" Energy="15" Leadership="24" HP="25" MP="90" LevelHP="1.5" LevelMP="1" VitalityHP="2" EnergyMP="1.5">
			<Skills>
				<Skill Index="60"/>
			</Skills>
			<Inventories>
				<Inventory Position="0" Section="0" Index="1" Durability="22"/>
			</Inventories>
		</Class>
	</ClassTable>`), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	conf.Server.PathCommon = filepath.Dir(dir)
	defer func() {
		conf.Server.PathCommon = pathCommon
	}()

	var table characterTable
	table.init()

	if len(table) != 1 {
		t.Fatalf("len(table) = %d, want 1", len(table))
	}
	character, ok := table[4]
	if !ok {
		t.Fatal("table[4] is missing")
	}
	if character.Level != 1 || character.MapNumber != 0 {
		t.Fatalf("level/map = %d/%d, want 1/0", character.Level, character.MapNumber)
	}
	if character.Strength != 26 || character.Dexterity != 20 || character.Vitality != 20 || character.Energy != 15 || character.Leadership != 24 {
		t.Fatalf("stats = %d/%d/%d/%d/%d, want 26/20/20/15/24", character.Strength, character.Dexterity, character.Vitality, character.Energy, character.Leadership)
	}
	if character.HP != 25 || character.MP != 90 {
		t.Fatalf("hp/mp = %d/%d, want 25/90", character.HP, character.MP)
	}
	if character.LevelHP != 1.5 || character.LevelMP != 1 || character.VitalityHP != 2 || character.EnergyMP != 1.5 {
		t.Fatalf("growth = %g/%g/%g/%g, want 1.5/1/2/1.5", character.LevelHP, character.LevelMP, character.VitalityHP, character.EnergyMP)
	}
	skill := character.Skills[60]
	if skill == nil || skill.SkillBase == nil {
		t.Fatal("skill 60 is not bound to SkillBase")
	}
	inventoryItem := character.Inventory.Items[0]
	if inventoryItem == nil || inventoryItem.ItemBase == nil {
		t.Fatal("inventory item 0 is not bound to ItemBase")
	}
	if inventoryItem.Section != 0 || inventoryItem.Index != 1 || inventoryItem.Durability != 22 {
		t.Fatalf("inventory item = %d/%d/%d, want 0/1/22", inventoryItem.Section, inventoryItem.Index, inventoryItem.Durability)
	}
	if inventoryItem.MaxDurability == 0 {
		t.Fatal("inventory item MaxDurability was not calculated")
	}
	if !character.Inventory.Flags[0] {
		t.Fatal("inventory flag 0 = false, want true")
	}
}

func TestSocketEquippedEffectsAndEvents(t *testing.T) {
	p := Player{msgChan: make(chan any, 20)}
	p.Object.Objecter = &p
	p.inventory.Items = make([]*item.Item, 237)
	p.MaxHP, p.MaxMP, p.MaxAG = 1000, 500, 100
	p.leftAttackMin, p.leftAttackMax = 100, 200
	p.DefenseRate = 100
	p.Level = 100
	p.monsterDieGetHP = 0.125 // existing non-Socket recovery
	weapon := item.NewItem(0, 26)
	weapon.Durability = 100
	weapon.SocketSlot1, weapon.SocketSlot2, weapon.SocketSlot3 = 100, 110, 116
	weapon.SocketSlot4, weapon.SocketSlot5 = 121, 129
	weapon.SocketBonus = 1
	p.inventory.Items[0] = weapon
	armor := item.NewItem(0, 26)
	armor.Durability = 100
	armor.SocketSlot1, armor.SocketSlot2 = 136, 105
	p.inventory.Items[2] = armor
	p.calcSocketItem()
	if p.leftAttackMin != 116 || p.DefenseRate != 110 || p.doubleDamageRate != 3 {
		t.Fatalf("Socket combat stats = attack %d, block %d, double %d", p.leftAttackMin, p.DefenseRate, p.doubleDamageRate)
	}
	if p.recoverHP != 8 || p.agReduce != 40 {
		t.Fatalf("Socket event stats = recover HP %d, AG reduction %d", p.recoverHP, p.agReduce)
	}
	if hp := p.GetMonsterDieGetHP(); hp != 0.25 {
		t.Fatalf("kill recovery fraction = %g", hp)
	}
	_, ag := p.GetSkillMPAG(&skill.Skill{SkillBase: &skill.SkillBase{BPUsage: 100}})
	if ag != 60 {
		t.Fatalf("skill AG cost = %d, want 60", ag)
	}
	p.HP, p.autoRecoverHPTick = 100, 6
	p.recoverHPSD()
	if p.HP != 108 {
		t.Fatalf("Socket HP recovery = %d, want 108", p.HP)
	}
	p.inventory.Items[0] = nil
	p.inventory.Items[2] = nil
	p.monsterDieGetHP, p.recoverHP, p.agReduce, p.doubleDamageRate = 0.125, 0, 0, 0
	p.calcSocketItem()
	if p.doubleDamageRate != 0 || p.GetMonsterDieGetHP() != 0.125 {
		t.Fatal("unequipped Socket item still contributes")
	}
}

func TestSocketEquipmentRecalculation(t *testing.T) {
	p := Player{msgChan: make(chan any, 30)}
	p.Object.Objecter = &p
	p.Class = int(class.Wizard)
	p.Level = 100
	p.strength, p.dexterity, p.vitality, p.energy = 20, 20, 20, 100
	p.inventory.Items = make([]*item.Item, 237)
	p.calc()
	baseHP := p.MaxHP
	it := item.NewItem(0, 26)
	it.Durability = 100
	it.SocketSlot1 = 122 // first-level maximum HP increase
	it.SocketSlot2 = 116 // first-level HP recovery on kill
	p.inventory.Items[0] = it
	p.calc()
	if p.MaxHP <= baseHP {
		t.Fatalf("Socket maximum HP = %d, base = %d", p.MaxHP, baseHP)
	}
	if p.GetMonsterDieGetHP() != 0.125 {
		t.Fatalf("Socket kill HP recovery = %g, want 0.125", p.GetMonsterDieGetHP())
	}
	withItem := p.MaxHP
	p.calc()
	if p.MaxHP != withItem {
		t.Fatalf("Socket maximum HP accumulated across recalculations: %d, want %d", p.MaxHP, withItem)
	}
	p.inventory.Items[0] = nil
	p.calc()
	if p.MaxHP != baseHP || p.GetMonsterDieGetHP() != 0 {
		t.Fatalf("Socket effects after unequip = HP %d, kill recovery %g", p.MaxHP, p.GetMonsterDieGetHP())
	}
}

func TestSocketAttackRecalculation(t *testing.T) {
	p := Player{msgChan: make(chan any, 30)}
	p.Object.Objecter = &p
	p.Class = int(class.Wizard)
	p.Level = 100
	p.strength, p.dexterity, p.vitality, p.energy = 20, 20, 20, 100
	p.inventory.Items = make([]*item.Item, 237)
	weapon := item.NewItem(0, 26)
	weapon.Durability = 100
	p.inventory.Items[0] = weapon
	p.calc()
	baseMin, baseMax := p.AttackMin, p.AttackMax
	weapon.SocketSlot1 = 100
	weapon.SocketBonus = 1
	p.calc()
	if p.AttackMin != baseMin+16 || p.AttackMax != baseMax+16 {
		t.Fatalf("Socket attack after hand merge = %d/%d, base %d/%d", p.AttackMin, p.AttackMax, baseMin, baseMax)
	}
	p.calc()
	if p.AttackMin != baseMin+16 || p.AttackMax != baseMax+16 {
		t.Fatal("Socket attack accumulated across recalculations")
	}
	weapon.SocketSlot1, weapon.SocketBonus = 0, 0
	p.calc()
	if p.AttackMin != baseMin || p.AttackMax != baseMax {
		t.Fatalf("Socket attack after removing seed = %d/%d, want %d/%d", p.AttackMin, p.AttackMax, baseMin, baseMax)
	}
}

func TestSocketManaAndAGRecovery(t *testing.T) {
	p := Player{msgChan: make(chan any, 10)}
	p.Object.Objecter = &p
	p.inventory.Items = make([]*item.Item, 237)
	p.MaxHP, p.MaxMP, p.MaxAG = 1000, 500, 100
	p.monsterDieGetMP = 0.125 // existing non-Socket recovery
	it := item.NewItem(0, 26)
	it.Durability = 100
	it.SocketSlot1, it.SocketSlot2, it.SocketSlot3 = 117, 124, 125
	it.SocketSlot4 = 126
	p.inventory.Items[0] = it
	p.calcSocketItem()
	if p.GetMonsterDieGetMP() != 0.25 || p.recoverMP != 7 || p.MaxAG != 150 || p.recoverAG != 5 {
		t.Fatalf("Socket mana/AG effects = kill %g, recover %d, max AG %d, AG recovery %d", p.GetMonsterDieGetMP(), p.recoverMP, p.MaxAG, p.recoverAG)
	}
	p.MP, p.AG = 0, 0
	p.autoRecoverMPTick = 2
	p.recoverMPAG()
	if p.MP != 25 || p.AG != 14 {
		t.Fatalf("Socket MP/AG recovery = %d/%d, want 25/14", p.MP, p.AG)
	}
	p.inventory.Items[0] = nil
	p.inventory.Items[12] = it
	p.monsterDieGetMP, p.recoverMP, p.recoverAG = 0.125, 0, 0
	p.calcSocketItem()
	if p.GetMonsterDieGetMP() != 0.125 || p.recoverMP != 0 || p.recoverAG != 0 {
		t.Fatal("Socket item in backpack contributes")
	}
}

func TestSocketBonusRejectsOtherSection(t *testing.T) {
	p := Player{}
	p.inventory.Items = make([]*item.Item, 237)
	p.Level = 100
	p.leftAttackMin = 100
	shield := item.NewItem(6, 17)
	shield.Durability = 50
	shield.SocketSlot1 = 100
	shield.SocketBonus = 1 // bonus 0 applies only to sections 0-4
	p.inventory.Items[1] = shield
	p.calcSocketItem()
	if p.leftAttackMin != 105 {
		t.Fatalf("attack with inapplicable Socket bonus = %d, want 105", p.leftAttackMin)
	}
}

func TestSocketAttackUsesEachHand(t *testing.T) {
	p := Player{}
	p.inventory.Items = make([]*item.Item, 237)
	p.Level = 100
	p.leftAttackMin, p.leftAttackMax = 100, 200
	p.rightAttackMin, p.rightAttackMax = 40, 80
	p.magicAttackMin, p.magicAttackMax = 10, 20
	p.curseAttackMin, p.curseAttackMax = 7, 9
	weapon := item.NewItem(0, 26)
	weapon.Durability = 100
	weapon.SocketSlot1, weapon.SocketSlot2, weapon.SocketSlot3 = 100, 102, 103
	weapon.SocketBonus = 1
	p.inventory.Items[0] = weapon
	p.calcSocketItem()
	if p.leftAttackMin != 136 || p.leftAttackMax != 246 || p.rightAttackMin != 76 || p.rightAttackMax != 126 {
		t.Fatalf("Socket hand attacks = %d/%d, %d/%d", p.leftAttackMin, p.leftAttackMax, p.rightAttackMin, p.rightAttackMax)
	}
	if p.magicAttackMin != 46 || p.magicAttackMax != 66 || p.curseAttackMin != 7 || p.curseAttackMax != 9 {
		t.Fatalf("Socket magic/curse attacks = %d/%d, %d/%d", p.magicAttackMin, p.magicAttackMax, p.curseAttackMin, p.curseAttackMax)
	}
	p.applySocketEffect(item.SocketEffect{Index: item.SocketEffectAttack, ValueType: 2, Value: 10}, 0, 0, 0)
	if p.leftAttackMin != 149 || p.leftAttackMax != 270 || p.rightAttackMin != 83 || p.rightAttackMax != 138 {
		t.Fatalf("Socket percent attacks = %d/%d, %d/%d", p.leftAttackMin, p.leftAttackMax, p.rightAttackMin, p.rightAttackMax)
	}
}

func TestSetAttackBonusPrecedesSocketPercent(t *testing.T) {
	p := Player{}
	p.leftAttackMin, p.leftAttackMax = 100, 200
	p.rightAttackMin, p.rightAttackMax = 40, 80
	p.applySetEffect(item.SetEffectIncAttackMin, 20, false)
	p.applySetEffect(item.SetEffectIncAttackMax, 30, false)
	p.applySocketEffect(item.SocketEffect{Index: item.SocketEffectAttack, ValueType: 2, Value: 10}, 0, 0, 0)
	if p.leftAttackMin != 132 || p.leftAttackMax != 253 || p.rightAttackMin != 66 || p.rightAttackMax != 121 {
		t.Fatalf("set and socket attack = %d/%d, %d/%d", p.leftAttackMin, p.leftAttackMax, p.rightAttackMin, p.rightAttackMax)
	}
}

func TestSocketDefenseAndResources(t *testing.T) {
	p := Player{}
	p.inventory.Items = make([]*item.Item, 237)
	p.Defense, p.MagicDefense, p.MaxSD = 100, 40, 200
	p.MaxHP, p.MaxMP, p.MaxAG = 1000, 500, 100
	shield := item.NewItem(6, 17)
	shield.Durability = 100
	shield.SocketSlot1, shield.SocketSlot2 = 111, 112
	shield.SocketBonus = 5
	p.inventory.Items[1] = shield
	weapon := item.NewItem(0, 26)
	weapon.Durability = 100
	weapon.SocketSlot1, weapon.SocketSlot2, weapon.SocketSlot3 = 122, 122, 123
	weapon.SocketSlot4 = 125
	p.inventory.Items[0] = weapon
	p.calcSocketItem()
	if p.Defense != 164 || p.MagicDefense != 94 || p.MaxSD != 200 {
		t.Fatalf("Socket defenses = %d/%d, SD %d", p.Defense, p.MagicDefense, p.MaxSD)
	}
	if p.MaxHP != 1080 || p.MaxMP != 520 || p.MaxAG != 150 {
		t.Fatalf("Socket resources = %d/%d/%d", p.MaxHP, p.MaxMP, p.MaxAG)
	}
	p.inventory.Items[1] = nil
	p.inventory.Items[2] = shield
	p.Defense, p.MagicDefense = 100, 40
	p.shieldDefense = 0
	p.calcSocketItem()
	if p.Defense != 154 || p.MagicDefense != 94 {
		t.Fatalf("Socket defenses without shield = %d/%d", p.Defense, p.MagicDefense)
	}
}
