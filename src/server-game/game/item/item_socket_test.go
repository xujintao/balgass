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

func TestSocketConfigRejectsInvalidRates(t *testing.T) {
	for _, tc := range []struct {
		name  string
		rate  string
		count int
	}{
		{"duplicate section", `CatID="0" Count1="10000"`, 13},
		{"missing section", `CatID="0" Count1="10000"`, 11},
		{"invalid weight", `CatID="0" Count1="-1"`, 12},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var xml strings.Builder
			xml.WriteString(`<SocketSystem><SocketRateSettings>`)
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
				if recover() == nil {
					t.Fatal("invalid socket rates were accepted")
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
