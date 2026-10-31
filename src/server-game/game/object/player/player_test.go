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
	if !p.canUseItem(it) {
		t.Fatal("canUseItem() = false, want true")
	}
	it.ReqEnergy = 31
	if p.canUseItem(it) {
		t.Fatal("canUseItem() = true with insufficient energy")
	}
	it.ReqEnergy = 30
	it.ReqClass[class.Wizard] = 2
	if p.canUseItem(it) {
		t.Fatal("canUseItem() = true without required change-up")
	}

	harmonyWeapon := &item.Item{
		ItemBase: &item.ItemBase{
			ReqStrength: 30,
			ReqClass:    [8]int{class.Wizard: 1},
		},
		Section: 0,
		Level:   0,
		Harmony: 3,
	}
	harmonyWeapon.Calc()
	p.changeUp = 0
	p.strength = 24
	if !p.canUseItem(harmonyWeapon) {
		t.Fatal("active Harmony strength reduction was not applied to item requirements")
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

func Test380EffectsRequireEnabledDurableWornItems(t *testing.T) {
	p := newItemTestPlayer()
	it := item.NewItem(0, 22)
	it.Durability = 1
	p.inventory.Items[0] = it
	p.attackRatePVP, p.defenseRatePVP = 100, 200
	p.MaxHP, p.MaxSD = 1000, 500
	p.calc380Item()
	if p.GetAttackBonusPVP() != 0 || p.attackRatePVP != 100 {
		t.Fatal("disabled 380 option affected player stats")
	}

	it.Option380 = true
	p.attackRatePVP, p.defenseRatePVP = 100, 200
	p.calc380Item()
	if p.GetAttackBonusPVP() != 200 || p.attackRatePVP != 110 || p.GetDefenseBonusPVP() != 0 {
		t.Fatalf("enabled 380 weapon effects: attack=%d rate=%d defense=%d", p.GetAttackBonusPVP(), p.attackRatePVP, p.GetDefenseBonusPVP())
	}
	if p.MaxHP != 1000 || p.MaxSD != 500 {
		t.Fatalf("380 weapon changed resources: HP=%d SD=%d", p.MaxHP, p.MaxSD)
	}

	it.Durability = 0
	p.attackRatePVP = 100
	p.calc380Item()
	if p.GetAttackBonusPVP() != 0 || p.attackRatePVP != 100 {
		t.Fatal("broken 380 equipment still affected stats")
	}
	p.inventory.Items[0] = nil
	p.inventory.Items[12] = it
	p.attackRatePVP = 100
	p.calc380Item()
	if p.GetAttackBonusPVP() != 0 || p.attackRatePVP != 100 {
		t.Fatal("380 item in the backpack affected stats")
	}
}

func Test380SDRecoveryHonorsServerAndItemOptions(t *testing.T) {
	info := &conf.CommonServer.GameServerInfo
	oldShieldSystem, oldEnable, oldSafeZone := info.ShieldSystemEnable, info.SDAutoRefillEnable, info.SDAutoRefillSafeZoneEnable
	info.ShieldSystemEnable, info.SDAutoRefillEnable, info.SDAutoRefillSafeZoneEnable = true, true, true
	t.Cleanup(func() {
		info.ShieldSystemEnable, info.SDAutoRefillEnable, info.SDAutoRefillSafeZoneEnable = oldShieldSystem, oldEnable, oldSafeZone
	})
	p := newItemTestPlayer()
	p.MapNumber = 0
	p.X, p.Y = findItemTestMapPosition(t)
	p.MaxSD, p.SD = 30000, 100
	p.autoRecoverSDTime = time.Now().Add(-20 * time.Second)
	p.recoverHPSD()
	if p.SD != 100 {
		t.Fatalf("SD recovered outside safe zone without 380 option: %d", p.SD)
	}
	if time.Since(p.autoRecoverSDTime) > time.Second {
		t.Fatal("SD recovery timer was not reset outside the safe zone")
	}
	info.ShieldSystemEnable = false
	p.autoRecoverSDTime = time.Now().Add(-20 * time.Second)
	previousTimer := p.autoRecoverSDTime
	p.recoverHPSD()
	if p.SD != 100 || p.autoRecoverSDTime != previousTimer {
		t.Fatal("disabled shield system changed SD or its recovery timer")
	}
	info.ShieldSystemEnable = true
	p.autoRecoverSD = true
	p.autoRecoverSDRate = 20
	p.autoRecoverSDTime = time.Now().Add(-10 * time.Second)
	p.recoverHPSD()
	if p.SD != 196 {
		t.Fatalf("380 SD recovery = %d, want 196", p.SD)
	}
	info.SDAutoRefillEnable = false
	p.SD = 100
	p.autoRecoverSDTime = time.Now().Add(-10 * time.Second)
	p.recoverHPSD()
	if p.SD != 100 {
		t.Fatalf("380 option bypassed server SD recovery switch: %d", p.SD)
	}
	info.SDAutoRefillEnable = true
	p.SD = p.MaxSD
	p.autoRecoverSDTime = time.Now().Add(-20 * time.Second)
	p.recoverHPSD()
	if time.Since(p.autoRecoverSDTime) > time.Second {
		t.Fatal("SD recovery timer was not reset at full SD")
	}
}

func TestSDAutoRecoveryUsesElapsedTimeTiers(t *testing.T) {
	info := &conf.CommonServer.GameServerInfo
	oldShieldSystem, oldEnable, oldSafeZone := info.ShieldSystemEnable, info.SDAutoRefillEnable, info.SDAutoRefillSafeZoneEnable
	info.ShieldSystemEnable, info.SDAutoRefillEnable, info.SDAutoRefillSafeZoneEnable = true, true, false
	t.Cleanup(func() {
		info.ShieldSystemEnable, info.SDAutoRefillEnable, info.SDAutoRefillSafeZoneEnable = oldShieldSystem, oldEnable, oldSafeZone
	})

	for _, tt := range []struct {
		elapsed time.Duration
		wantSD  int
	}{
		{elapsed: 9 * time.Second, wantSD: 100},
		{elapsed: 10 * time.Second, wantSD: 180},
		{elapsed: 15 * time.Second, wantSD: 200},
		{elapsed: 25 * time.Second, wantSD: 220},
	} {
		p := newItemTestPlayer()
		p.MaxSD, p.SD = 30000, 100
		p.autoRecoverSDTime = time.Now().Add(-tt.elapsed)
		p.recoverHPSD()
		if p.SD != tt.wantSD {
			t.Errorf("SD after %s = %d, want %d", tt.elapsed, p.SD, tt.wantSD)
		}
	}
}

func TestSDAutoRecoveryInitializesZeroTimer(t *testing.T) {
	info := &conf.CommonServer.GameServerInfo
	oldShieldSystem, oldEnable, oldSafeZone := info.ShieldSystemEnable, info.SDAutoRefillEnable, info.SDAutoRefillSafeZoneEnable
	info.ShieldSystemEnable, info.SDAutoRefillEnable, info.SDAutoRefillSafeZoneEnable = true, true, false
	t.Cleanup(func() {
		info.ShieldSystemEnable, info.SDAutoRefillEnable, info.SDAutoRefillSafeZoneEnable = oldShieldSystem, oldEnable, oldSafeZone
	})

	p := newItemTestPlayer()
	p.MaxSD, p.SD = 30000, 100
	p.recoverHPSD()
	if p.SD != 100 {
		t.Fatalf("first SD recovery call changed SD to %d", p.SD)
	}
	if p.autoRecoverSDTime.IsZero() || time.Since(p.autoRecoverSDTime) > time.Second {
		t.Fatal("first SD recovery call did not initialize its timer")
	}
}

func TestRecalculationPreservesTemporaryResources(t *testing.T) {
	p := Player{msgChan: make(chan any, 30)}
	p.Object.Objecter = &p
	p.Class = int(class.Wizard)
	p.Level = 100
	p.strength, p.dexterity, p.vitality, p.energy = 20, 20, 20, 100
	p.inventory.Items = make([]*item.Item, 237)
	p.calc()
	baseHP, baseMP := p.MaxHP, p.MaxMP
	p.AddHP = 100
	p.HP = baseHP + 100
	p.AddMP = 50
	p.MP = baseMP + 50
	p.calc()
	if p.MaxHP != baseHP || p.AddHP != 100 || p.HP != baseHP+100 ||
		p.MaxMP != baseMP || p.AddMP != 50 || p.MP != baseMP+50 {
		t.Fatalf("recalculated HP/base/add and MP/base/add = %d/%d/%d and %d/%d/%d", p.HP, p.MaxHP, p.AddHP, p.MP, p.MaxMP, p.AddMP)
	}
	p.AddHP = -20
	p.AddMP = -20
	p.calc()
	if p.HP != p.MaxHP+p.AddHP || p.MP != p.MaxMP+p.AddMP {
		t.Fatalf("resources after reduced caps = HP %d/%d, MP %d/%d", p.HP, p.MaxHP+p.AddHP, p.MP, p.MaxMP+p.AddMP)
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

func TestSocketWeaponAndStaffAttackBonuses(t *testing.T) {
	for _, tc := range []struct {
		name         string
		section      int
		bonus        int
		wantPhysical int
		wantMagic    int
	}{
		{name: "weapon", section: 0, bonus: 1, wantPhysical: 11},
		{name: "fourth-tier weapon", section: 0, bonus: 7, wantPhysical: 22},
		{name: "staff", section: 5, bonus: 3, wantMagic: 5},
		{name: "fourth-tier staff", section: 5, bonus: 9, wantMagic: 10},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := Player{}
			p.inventory.Items = make([]*item.Item, 237)
			p.leftAttackMin, p.leftAttackMax = 100, 200
			p.magicAttackMin, p.magicAttackMax = 10, 20
			p.inventory.Items[0] = &item.Item{
				ItemBase:    &item.ItemBase{Type: item.TypeSocket},
				Section:     tc.section,
				Durability:  1,
				SocketBonus: tc.bonus,
			}
			p.calcSocketItem()
			if p.leftAttackMin != 100+tc.wantPhysical || p.leftAttackMax != 200+tc.wantPhysical ||
				p.magicAttackMin != 10+tc.wantMagic || p.magicAttackMax != 20+tc.wantMagic {
				t.Fatalf("weapon/staff bonus attacks = physical %d/%d, magic %d/%d; want physical %d/%d and magic %d/%d",
					p.leftAttackMin, p.leftAttackMax, p.magicAttackMin, p.magicAttackMax,
					100+tc.wantPhysical, 200+tc.wantPhysical, 10+tc.wantMagic, 20+tc.wantMagic)
			}
		})
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
	if p.leftAttackMin != 136 || p.leftAttackMax != 246 || p.rightAttackMin != 40 || p.rightAttackMax != 80 {
		t.Fatalf("Socket hand attacks = %d/%d, %d/%d", p.leftAttackMin, p.leftAttackMax, p.rightAttackMin, p.rightAttackMax)
	}
	if p.magicAttackMin != 46 || p.magicAttackMax != 66 || p.curseAttackMin != 7 || p.curseAttackMax != 9 {
		t.Fatalf("Socket magic/curse attacks = %d/%d, %d/%d", p.magicAttackMin, p.magicAttackMax, p.curseAttackMin, p.curseAttackMax)
	}

	rightHandPlayer := Player{}
	rightHandPlayer.inventory.Items = make([]*item.Item, 237)
	rightHandPlayer.Level = 100
	rightHandPlayer.leftAttackMin, rightHandPlayer.leftAttackMax = 100, 200
	rightHandPlayer.rightAttackMin, rightHandPlayer.rightAttackMax = 40, 80
	rightHandPlayer.magicAttackMin, rightHandPlayer.magicAttackMax = 10, 20
	rightHandPlayer.inventory.Items[1] = weapon.Copy()
	rightHandPlayer.calcSocketItem()
	if rightHandPlayer.leftAttackMin != 100 || rightHandPlayer.leftAttackMax != 200 ||
		rightHandPlayer.rightAttackMin != 76 || rightHandPlayer.rightAttackMax != 126 {
		t.Fatalf("right-hand Socket attacks = left %d/%d, right %d/%d",
			rightHandPlayer.leftAttackMin, rightHandPlayer.leftAttackMax,
			rightHandPlayer.rightAttackMin, rightHandPlayer.rightAttackMax)
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
	if p.Defense != 163 || p.MagicDefense != 94 || p.MaxSD != 200 {
		t.Fatalf("Socket defenses = %d/%d, SD %d", p.Defense, p.MagicDefense, p.MaxSD)
	}
	if p.MaxHP != 1080 || p.MaxMP != 520 || p.MaxAG != 150 {
		t.Fatalf("Socket resources = %d/%d/%d", p.MaxHP, p.MaxMP, p.MaxAG)
	}
	p.inventory.Items[1] = nil
	p.inventory.Items[2] = shield
	p.Defense, p.MagicDefense = 100, 40
	p.calcSocketItem()
	if p.Defense != 154 || p.MagicDefense != 94 {
		t.Fatalf("Socket defenses without shield = %d/%d", p.Defense, p.MagicDefense)
	}
}

func TestSocketShieldDefenseIsAppliedByPlayer(t *testing.T) {
	for _, tc := range []struct {
		name   string
		shield *item.Item
		want   int
	}{
		{name: "valid shield", shield: func() *item.Item {
			it := item.NewItem(6, 17)
			it.Durability = 1
			return it
		}(), want: 125},
		{name: "no shield", want: 100},
		{name: "broken shield", shield: item.NewItem(6, 17), want: 100},
		{name: "other equipment", shield: item.NewItem(0, 26), want: 100},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := Player{}
			p.inventory.Items = make([]*item.Item, 237)
			p.Defense = 100
			socketItem := item.NewItem(0, 26)
			socketItem.Durability = 1
			socketItem.SocketSlot1, socketItem.SocketSlot2, socketItem.SocketSlot3 = 112, 212, 112
			p.inventory.Items[0] = socketItem
			p.inventory.Items[1] = tc.shield
			p.calcSocketItem()
			if p.Defense != tc.want {
				t.Fatalf("Defense after Socket shield effects = %d, want %d", p.Defense, tc.want)
			}
		})
	}
}

func TestSetAndExcellentEquipmentRecalculation(t *testing.T) {
	p := newItemTestPlayer()
	p.Class = int(class.Knight)
	p.strength, p.dexterity, p.vitality, p.energy = 1110, 500, 25, 303
	p.calc()
	baseMin, baseMax, baseHP := p.AttackMin, p.AttackMax, p.MaxHP
	for position := 2; position <= 6; position++ {
		it := item.NewItem(position+5, 6)
		it.Set = 6
		it.Calc()
		it.Durability = it.MaxDurability
		p.inventory.Items[position] = it
	}
	weapon := item.NewItem(0, 19)
	weapon.Lucky, weapon.ExcellentAttackRate = true, true
	weapon.ExcellentAttackLevel, weapon.ExcellentAttackPercent = true, true
	weapon.Calc()
	weapon.Durability = weapon.MaxDurability
	p.inventory.Items[0] = weapon
	p.calc()
	if p.addStrength != 40 || p.addSkillAttack != 40 || p.criticalAttackRate != 5 || p.excellentAttackRate != 10 {
		t.Fatalf("set and excellent effects = strength %d, skill attack %d, critical %d, excellent %d", p.addStrength, p.addSkillAttack, p.criticalAttackRate, p.excellentAttackRate)
	}
	attackMin, attackMax, maxHP := p.AttackMin, p.AttackMax, p.MaxHP
	p.calc()
	if p.AttackMin != attackMin || p.AttackMax != attackMax || p.MaxHP != maxHP || p.addStrength != 40 || p.addSkillAttack != 40 {
		t.Fatal("set or excellent effects accumulated across recalculations")
	}
	for position := range p.inventory.WearingItems() {
		p.inventory.Items[position] = nil
	}
	p.calc()
	if p.AttackMin != baseMin || p.AttackMax != baseMax || p.MaxHP != baseHP || p.addStrength != 0 || p.addSkillAttack != 0 || p.criticalAttackRate != 0 || p.excellentAttackRate != 0 {
		t.Fatal("set or excellent effects remained after unequipping")
	}
}

func TestHarmonyEquipmentRecalculation(t *testing.T) {
	p := newItemTestPlayer()
	p.Class = int(class.Knight)
	p.strength, p.dexterity, p.vitality, p.energy = 100, 100, 100, 100

	weapon := item.NewItem(0, 0)
	weapon.Calc()
	p.inventory.Items[0] = weapon

	armor := item.NewItem(6, 0)
	armor.Level = 3
	armor.Calc()
	p.inventory.Items[2] = armor
	p.calc()
	baseAttackMin, baseHP := p.AttackMin, p.MaxHP

	weapon.Harmony, weapon.Level = 1, 0
	armor.Harmony, armor.Level = 3, 3
	weapon.Calc()
	armor.Calc()
	p.calc()
	if p.AttackMin != baseAttackMin+2 || p.MaxHP != baseHP+7 {
		t.Fatalf("Harmony equipment attack/max HP = %d/%d, want %d/%d", p.AttackMin, p.MaxHP, baseAttackMin+2, baseHP+7)
	}
	p.calc()
	if p.AttackMin != baseAttackMin+2 || p.MaxHP != baseHP+7 {
		t.Fatal("Harmony effects accumulated across recalculation")
	}
	p.inventory.Items[0], p.inventory.Items[2] = nil, nil
	p.calc()
	basePlayer := newItemTestPlayer()
	basePlayer.Class = int(class.Knight)
	basePlayer.strength, basePlayer.dexterity, basePlayer.vitality, basePlayer.energy = 100, 100, 100, 100
	basePlayer.calc()
	if p.AttackMin != basePlayer.AttackMin || p.AttackMax != basePlayer.AttackMax ||
		p.MaxHP != basePlayer.MaxHP || p.harmonyDecSDRate != 0 || p.harmonyIgnoreSDRate != 0 || p.harmonyIncSDRate != 0 {
		t.Fatalf("Harmony effects remained after unequipping: attack=%d/%d want %d/%d, max HP=%d want %d, SD rates=%d/%d/%d", p.AttackMin, p.AttackMax, basePlayer.AttackMin, basePlayer.AttackMax, p.MaxHP, basePlayer.MaxHP, p.harmonyDecSDRate, p.harmonyIgnoreSDRate, p.harmonyIncSDRate)
	}
}

func TestHarmonyAttackPanelBonusUsesEquippedHand(t *testing.T) {
	p := newItemTestPlayer()
	p.leftAttackMin, p.leftAttackMax = 10, 20
	p.rightAttackMin, p.rightAttackMax = 30, 40

	leftWeapon := item.NewItem(0, 0)
	leftWeapon.Harmony = int(item.HarmonyEffectWeaponIncDamageMin)
	rightWeapon := item.NewItem(0, 0)
	rightWeapon.Harmony = int(item.HarmonyEffectWeaponIncDamageMax)
	p.inventory.Items[0] = leftWeapon
	p.inventory.Items[1] = rightWeapon

	p.calcHarmonyItem()
	if p.leftAttackMin != 12 || p.leftAttackMax != 20 ||
		p.rightAttackMin != 30 || p.rightAttackMax != 43 {
		t.Fatalf("Harmony hand attack panels = %d/%d, %d/%d; want 12/20, 30/43",
			p.leftAttackMin, p.leftAttackMax, p.rightAttackMin, p.rightAttackMax)
	}
}

func TestSocketContributionsUseCurrentAttributes(t *testing.T) {
	p := newItemTestPlayer()
	p.MaxHP, p.MaxMP, p.MaxAG = 1001, 501, 100
	p.DefenseRate, p.agReduce = 99, 80
	weapon := item.NewItem(0, 26)
	weapon.Durability = 100
	weapon.SocketSlot1, weapon.SocketSlot2 = 110, 110
	weapon.SocketSlot3 = 105
	weapon.SocketSlot4, weapon.SocketSlot5 = 122, 122
	p.inventory.Items[0] = weapon
	p.calcSocketItem()
	if p.DefenseRate != 118 || p.agReduce != 100 {
		t.Fatalf("Socket percent defense and AG reduction = %d/%d, want 118/100", p.DefenseRate, p.agReduce)
	}
	if p.MaxHP != 1081 || p.MaxMP != 501 || p.MaxAG != 100 {
		t.Fatalf("Socket resources = %d/%d/%d, want 1081/501/100", p.MaxHP, p.MaxMP, p.MaxAG)
	}
}

func TestExcellentContributionsUseCurrentAttributes(t *testing.T) {
	p := newItemTestPlayer()
	p.Level = 388
	p.leftAttackMin, p.leftAttackMax = 100, 200
	p.rightAttackMin, p.rightAttackMax = 40, 80
	p.magicAttackMin, p.magicAttackMax = 10, 20
	p.MaxHP, p.MaxMP, p.DefenseRate = 1001, 501, 99
	p.monsterDieGetHP = 0.125
	weapon := item.NewItem(0, 19)
	weapon.ExcellentAttackRate = true
	weapon.ExcellentAttackLevel, weapon.ExcellentAttackPercent = true, true
	weapon.ExcellentAttackSpeed, weapon.ExcellentAttackHP = true, true
	weapon.ExcellentDefenseHP, weapon.ExcellentDefenseMP, weapon.ExcellentDefenseRate = true, true, true
	p.inventory.Items[0], p.inventory.Items[1] = weapon, weapon
	p.inventory.Items[7] = &item.Item{Level: 9, Wing2ExcellentHP: true, Wing2ExcellentMP: true, Wing2ExcellentAG: true, Wing2ExcellentSpeed: true, Wing3ExcellentIgnore: true}
	p.calcExcellentItem()
	if p.leftAttackMin != 121 || p.leftAttackMax != 223 || p.rightAttackMin != 60 || p.rightAttackMax != 100 || p.magicAttackMin != 48 || p.magicAttackMax != 59 {
		t.Fatalf("excellent hand attacks = %d/%d, %d/%d, magic %d/%d", p.leftAttackMin, p.leftAttackMax, p.rightAttackMin, p.rightAttackMax, p.magicAttackMin, p.magicAttackMax)
	}
	if p.MaxHP != 1177 || p.MaxMP != 636 || p.MaxAG != 50 || p.DefenseRate != 118 {
		t.Fatalf("excellent resources and defense = %d/%d/%d, %d", p.MaxHP, p.MaxMP, p.MaxAG, p.DefenseRate)
	}
	if p.criticalAttackRate != 0 || p.excellentAttackRate != 20 || p.AttackSpeed != 19 || p.magicSpeed != 19 || p.ignoreDefenseRate != 5 || p.monsterDieGetHP != 0.375 {
		t.Fatalf("excellent rates, speeds and recovery = %d/%d, %d/%d, %d, %g", p.criticalAttackRate, p.excellentAttackRate, p.AttackSpeed, p.magicSpeed, p.ignoreDefenseRate, p.monsterDieGetHP)
	}
}
