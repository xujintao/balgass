package item

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xujintao/balgass/src/server-game/conf"
)

func TestSocketSlotRates(t *testing.T) {
	if SocketManager.rates[0][0] != 6000 || SocketManager.twoHand[0] != 6000 {
		t.Fatal("socket slot rates were not loaded")
	}
	m := socketManager{rates: [12][5]int{{0, 0, 10000}}, twoHand: [5]int{0, 0, 0, 0, 10000}}
	it := NewItem(0, 26)
	if got := m.RandomSlotCount(it); got != 3 {
		t.Fatalf("one-hand slot count = %d, want 3", got)
	}
	it.ItemBase = &ItemBase{Type: TypeSocket, TwoHand: true}
	if got := m.RandomSlotCount(it); got != 5 {
		t.Fatalf("two-hand slot count = %d, want 5", got)
	}
}

func TestSphereDropSelection(t *testing.T) {
	if !SocketManager.sphereEnabled || SocketManager.sphereRates[0] != 50 || SocketManager.sphereLevels[0] != 102 {
		t.Fatal("Sphere drop settings were not loaded")
	}
	m := socketManager{sphereEnabled: true, sphereRates: [5]int{10000, 0, 10000}, sphereLevels: [5]int{100, 0, 150}}
	if _, ok := m.RandomSphereIndex(99); ok {
		t.Fatal("Sphere dropped below the minimum monster level")
	}
	if index, ok := m.RandomSphereIndex(100); !ok || index != 70 {
		t.Fatalf("level 100 Sphere = %d, %t", index, ok)
	}
	if index, ok := m.RandomSphereIndex(150); !ok || index != 72 {
		t.Fatalf("higher-tier Sphere = %d, %t", index, ok)
	}
	m.sphereEnabled = false
	if _, ok := m.RandomSphereIndex(150); ok {
		t.Fatal("Sphere dropped while disabled")
	}
	m.sphereEnabled = true
	m.sphereRates = [5]int{}
	if _, ok := m.RandomSphereIndex(150); ok {
		t.Fatal("Sphere dropped at zero rate")
	}
}

func TestTetraDropSelection(t *testing.T) {
	if !SocketManager.tetraEnabled || SocketManager.tetraRates != [3]int{} {
		t.Fatal("Tetra settings were not loaded")
	}
	m := socketManager{
		tetraEnabled:   true,
		tetraRates:     [3]int{10000, 0, 10000},
		tetraMinLevels: [3]int{50, 100, 150},
		tetraMaxLevels: [3]int{300, 300, 300},
	}
	if _, ok := m.RandomTetraIndex(49); ok {
		t.Fatal("Tetra dropped below level range")
	}
	if index, ok := m.RandomTetraIndex(50); !ok || index != 247 {
		t.Fatalf("Bronze = %d, %t", index, ok)
	}
	if index, ok := m.RandomTetraIndex(150); !ok || index != 245 {
		t.Fatalf("Gold = %d, %t", index, ok)
	}
	m.tetraEnabled = false
	if _, ok := m.RandomTetraIndex(150); ok {
		t.Fatal("Tetra dropped when disabled")
	}
	m.tetraEnabled, m.tetraRates = true, [3]int{}
	if _, ok := m.RandomTetraIndex(150); ok {
		t.Fatal("Tetra dropped at zero rate")
	}
}

func TestTetraConfigValidation(t *testing.T) {
	file := filepath.Join(conf.PathCommon, "SocketSystem", "IGC_SocketSystem.xml")
	valid, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, old, replacement string }{
		{"switch", `TetraDrop="1"`, `TetraDrop="2"`},
		{"rate", `Bronze Rate="0"`, `Bronze Rate="10001"`},
		{"range", `Bronze Rate="0" LevelMin="50" LevelMax="300"`, `Bronze Rate="0" LevelMin="301" LevelMax="300"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.MkdirAll(filepath.Join(dir, "SocketSystem"), 0755); err != nil {
				t.Fatal(err)
			}
			changed := strings.Replace(string(valid), tc.old, tc.replacement, 1)
			if changed == string(valid) {
				t.Fatal("replacement did not change XML")
			}
			if err := os.WriteFile(filepath.Join(dir, "SocketSystem", "IGC_SocketSystem.xml"), []byte(changed), 0644); err != nil {
				t.Fatal(err)
			}
			oldPath := conf.PathCommon
			conf.PathCommon = dir
			defer func() { conf.PathCommon = oldPath }()
			defer func() {
				if got := recover(); got == nil || !strings.Contains(fmt.Sprint(got), "TetraDrop") {
					t.Fatalf("panic = %v", got)
				}
			}()
			var manager socketManager
			manager.init()
		})
	}
}

func TestSphereDropConfigValidation(t *testing.T) {
	file := filepath.Join(conf.PathCommon, "SocketSystem", "IGC_SocketSystem.xml")
	valid, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, old, replacement string }{
		{"invalid switch", `SphereDrop="1"`, `SphereDrop="2"`},
		{"invalid rate", `Level1="50"`, `Level1="10001"`},
		{"invalid level", `Level1="102"`, `Level1="-1"`},
		{"missing rate", `Level1="50"`, ``},
	} {
		t.Run(tc.name, func(t *testing.T) {
			changed := strings.Replace(string(valid), tc.old, tc.replacement, 1)
			if changed == string(valid) {
				t.Fatal("test XML replacement failed")
			}
			dir := t.TempDir()
			if err := os.MkdirAll(filepath.Join(dir, "SocketSystem"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "SocketSystem", "IGC_SocketSystem.xml"), []byte(changed), 0644); err != nil {
				t.Fatal(err)
			}
			oldPath := conf.PathCommon
			conf.PathCommon = dir
			defer func() { conf.PathCommon = oldPath }()
			defer func() {
				if got := recover(); got == nil || !strings.Contains(fmt.Sprint(got), "SphereDrop") {
					t.Fatalf("panic = %v, want SphereDrop config error", got)
				}
			}()
			var manager socketManager
			manager.init()
		})
	}
}

func TestSocketConfigRejectsInvalidRates(t *testing.T) {
	for _, tc := range []struct {
		name, rate, errorText string
		count                 int
	}{
		{"duplicate section", `CatID="0" Count1="10000"`, "duplicate section 0", 13},
		{"missing section", `CatID="0" Count1="10000"`, "missing section 11", 11},
		{"invalid weight", `CatID="0" Count1="-1"`, "invalid section 0 slot weight", 12},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var xml strings.Builder
			xml.WriteString(`<SocketSystem SphereDrop="0" TetraDrop="0"><SphereDrop><Rate Level1="0" Level2="0" Level3="0" Level4="0" Level5="0"/><Level Level1="0" Level2="0" Level3="0" Level4="0" Level5="0"/></SphereDrop><TetraDrop><Bronze Rate="0" LevelMin="0" LevelMax="0"/><Silver Rate="0" LevelMin="0" LevelMax="0"/><Gold Rate="0" LevelMin="0" LevelMax="0"/></TetraDrop><Seed ExtractMoney="0"><Sphere CompositeMoney="0" SetMoney="0" RemoveMoney="0"/></Seed><SocketRateSettings>`)
			for i := 0; i < tc.count; i++ {
				if i == 0 || i == 12 {
					fmt.Fprintf(&xml, "<Socket %s/>", tc.rate)
				} else {
					fmt.Fprintf(&xml, `<Socket CatID="%d" Count1="10000"/>`, i)
				}
			}
			xml.WriteString(`<TwoHandWeapons Count1="10000"/></SocketRateSettings></SocketSystem>`)
			dir := t.TempDir()
			if err := os.MkdirAll(filepath.Join(dir, "SocketSystem"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "SocketSystem", "IGC_SocketSystem.xml"), []byte(xml.String()), 0644); err != nil {
				t.Fatal(err)
			}
			oldPath := conf.PathCommon
			conf.PathCommon = dir
			defer func() { conf.PathCommon = oldPath }()
			defer func() {
				if got := recover(); got == nil || !strings.Contains(fmt.Sprint(got), tc.errorText) {
					t.Fatalf("panic = %v, want %q", got, tc.errorText)
				}
			}()
			var manager socketManager
			manager.init()
		})
	}
}

func TestSocketSlotsJSONAndProtocol(t *testing.T) {
	it := NewItem(0, 26)
	if it.SocketBonus != 0 || it.SocketSlot1 != 0 || it.SocketSlot5 != 0 {
		t.Fatalf("new Socket item has open slots: %#v", it)
	}
	jsonData, err := json.Marshal(it)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"socket_bonus", "socket_slot1", "socket_slot2", "socket_slot3", "socket_slot4", "socket_slot5"} {
		if strings.Contains(string(jsonData), `"`+field+`"`) {
			t.Fatalf("zero Socket field %s was saved: %s", field, jsonData)
		}
	}
	it.MakeSocketSlots(2)
	data, err := it.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if got := data[6:12]; got[0] != 0xFF || got[1] != 0xFE || got[2] != 0xFE || got[3] != 0xFF || got[4] != 0xFF || got[5] != 0xFF {
		t.Fatalf("Socket protocol bytes = %x", got)
	}
	for _, tc := range []struct{ slot, encoded int }{{100, 0}, {201, 51}, {300, 100}} {
		it.SocketSlot2 = tc.slot
		data, err = it.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if data[8] != byte(tc.encoded) {
			t.Fatalf("Socket slot %d encoded as %d, want %d", tc.slot, data[8], tc.encoded)
		}
	}
	for _, slot := range []int{-1, 2, 99, 150, 600} {
		it.SocketSlot2 = slot
		if _, err := it.Marshal(); err == nil {
			t.Fatalf("invalid Socket slot %d was encoded", slot)
		}
	}
	it.SocketSlot2 = 1
	it.SocketBonus = 1
	data, err = it.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if data[6] != 0 {
		t.Fatalf("Socket bonus option 0 encoded as %x", data[6])
	}
	it.SocketBonus = 0
	for _, container := range []struct {
		name    string
		marshal func() ([]byte, error)
		restore func([]byte) (*Item, error)
	}{
		{"inventory", func() ([]byte, error) {
			inv := Inventory{PositionedItems: PositionedItems{Size: 237, Items: make([]*Item, 237)}}
			inv.Items[12] = it
			return json.Marshal(inv)
		}, func(data []byte) (*Item, error) {
			var inv Inventory
			err := json.Unmarshal(data, &inv)
			return inv.Items[12], err
		}},
		{"warehouse", func() ([]byte, error) {
			w := Warehouse{PositionedItems: PositionedItems{Size: 240, Items: make([]*Item, 240)}}
			w.Items[12] = it
			return json.Marshal(w)
		}, func(data []byte) (*Item, error) {
			var w Warehouse
			err := json.Unmarshal(data, &w)
			return w.Items[12], err
		}},
	} {
		t.Run(container.name, func(t *testing.T) {
			data, err := container.marshal()
			if err != nil {
				t.Fatal(err)
			}
			restored, err := container.restore(data)
			if err != nil {
				t.Fatal(err)
			}
			if restored.SocketBonus != 0 || restored.SocketSlot1 != 1 || restored.SocketSlot2 != 1 || restored.SocketSlot3 != 0 {
				t.Fatalf("restored Socket slots = %#v", restored)
			}
			it.SocketSlot2 = 201
			saved, err := container.marshal()
			it.SocketSlot2 = 1
			if err != nil {
				t.Fatal(err)
			}
			restored, err = container.restore(saved)
			if err != nil || restored.SocketSlot2 != 201 {
				t.Fatalf("restored Socket option = %#v, err = %v", restored, err)
			}
			it.SocketBonus = 1
			saved, err = container.marshal()
			it.SocketBonus = 0
			if err != nil {
				t.Fatal(err)
			}
			restored, err = container.restore(saved)
			if err != nil || restored.SocketBonus != 1 {
				t.Fatalf("restored Socket bonus = %#v, err = %v", restored, err)
			}
		})
	}
}

func TestSocketOptionsLoaded(t *testing.T) {
	if len(SocketManager.seedOptions) != 27 || len(SocketManager.bonusOptions) != 12 || len(SocketManager.specialOptions) != 2 {
		t.Fatalf("Socket option counts = %d/%d/%d", len(SocketManager.seedOptions), len(SocketManager.bonusOptions), len(SocketManager.specialOptions))
	}
	seed := SocketManager.seedOptions[36]
	if seed.element != 6 || seed.values[0] != 30 || seed.values[4] != 38 || seed.seedItemCode != Code(12, 65) {
		t.Fatalf("seed option 36 = %#v", seed)
	}
	bonus := SocketManager.bonusOptions[0]
	if bonus.minSection != 0 || bonus.maxSection != 4 || bonus.elements != [5]int{1, 5, 3, 0, 0} {
		t.Fatalf("bonus option 0 = %#v", bonus)
	}
	special := SocketManager.specialOptions[1]
	if special.rate != 100 || special.seedCounts != [6]int{1, 3, 1, 3, 1, 2} {
		t.Fatalf("special option 1 = %#v", special)
	}
}

func TestSocketEffectLookup(t *testing.T) {
	for _, slot := range []int{0, 1, 2, 199, 600, 650} {
		if _, ok := SocketManager.SeedEffect(slot); ok {
			t.Fatalf("invalid slot %d has a Seed effect", slot)
		}
	}
	effect, ok := SocketManager.SeedEffect(201)
	if !ok || effect.Index != SocketEffectAttackSpeed || effect.Element != 1 || effect.Value != 8 {
		t.Fatalf("slot 201 effect = %+v, %t", effect, ok)
	}
	it := NewItem(0, 26)
	it.SocketBonus = 7
	bonus, ok := SocketManager.BonusEffect(it.SocketBonus)
	if !ok || bonus.Index != SocketBonusWeaponAttack2 || bonus.Value != 22 {
		t.Fatalf("Socket bonus 6 = %+v, %t", bonus, ok)
	}
	specials := SocketManager.SpecialEffects([6]int{1, 3, 1, 3, 1, 2})
	if len(specials) != 2 || specials[0].Index != SocketEffectDoubleDamage || specials[1].Index != SocketEffectIgnoreDefense {
		t.Fatalf("Socket special effects = %+v", specials)
	}
	for _, invalid := range []int{0, 13} {
		if _, ok := SocketManager.BonusEffect(invalid); ok {
			t.Fatalf("invalid Socket bonus %d was accepted", invalid)
		}
	}
}

func TestSocketOptionConfigValidation(t *testing.T) {
	const seed = `<Option Index="0" ElementType="1" Level="0" BonusType="1" BonusValue1="1" BonusValue2="2" BonusValue3="3" BonusValue4="4" BonusValue5="5" PoolValue="10" Name="seed"/>`
	const bonus = `<Option Index="0" MinItemType="0" MaxItemType="4" BonusType="1" BonusValue="1" BonusRate="30" ElementType1="1" ElementType2="0" ElementType3="0" ElementType4="0" ElementType5="0" Name="bonus"/>`
	const special = `<Option Index="0" Type="0" BonusType="1" BonusValue="1" BonusRate="100" SeedCountFire="1" SeedCountWater="0" SeedCountIce="0" SeedCountWind="0" SeedCountLight="0" SeedCountEarth="0" Name="special"/>`
	const xml = `<SocketSystem><SocketItemOptionSettings>` + seed + `</SocketItemOptionSettings><SocketBonusSettings>` + bonus + `</SocketBonusSettings><SpecialOptionSettings>` + special + `</SpecialOptionSettings></SocketSystem>`
	for _, tc := range []struct {
		name, old, replacement, errorText string
	}{
		{"duplicate seed", `</SocketItemOptionSettings>`, seed + `</SocketItemOptionSettings>`, "duplicate seed option 0"},
		{"missing value", ` BonusValue5="5"`, ``, "missing BonusValue5"},
		{"invalid bonus element", `ElementType1="1"`, `ElementType1="7"`, "invalid element 7"},
		{"invalid special type", ` Type="0"`, ` Type="256"`, "invalid special option 0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := strings.Replace(xml, tc.old, tc.replacement, 1)
			if data == xml {
				t.Fatal("test XML replacement failed")
			}
			dir := t.TempDir()
			if err := os.MkdirAll(filepath.Join(dir, "SocketSystem"), 0755); err != nil {
				t.Fatal(err)
			}
			rateData, err := os.ReadFile(filepath.Join(conf.PathCommon, "SocketSystem", "IGC_SocketSystem.xml"))
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "SocketSystem", "IGC_SocketSystem.xml"), rateData, 0644); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "SocketSystem", "IGC_SocketOption.xml"), []byte(data), 0644); err != nil {
				t.Fatal(err)
			}
			oldPath := conf.PathCommon
			conf.PathCommon = dir
			defer func() { conf.PathCommon = oldPath }()
			defer func() {
				if got := recover(); got == nil || !strings.Contains(fmt.Sprint(got), tc.errorText) {
					t.Fatalf("panic = %v, want %q", got, tc.errorText)
				}
			}()
			var manager socketManager
			manager.init()
		})
	}
}
