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
)

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
	pathCommon := conf.ServerEnv.PathCommon
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
	conf.ServerEnv.PathCommon = filepath.Dir(dir)
	defer func() {
		conf.ServerEnv.PathCommon = pathCommon
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
