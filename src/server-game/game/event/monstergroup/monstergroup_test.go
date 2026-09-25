package monstergroup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/object"
)

const groupXML = `<MonsterGroupRegenSystem SpawnNotice="0">
<GroupSettings><Group Index="7" MapNumber="0" BossMonsterIndex="44" Duration="60"/></GroupSettings>
<SpotSettings><Group Index="7"><Area StartX="135" StartY="61" EndX="146" EndY="70"/></Group></SpotSettings>
<MonsterSettings><Group Index="7"><Monster Index="44" Count="2" OverrideDefaultSettings="0"/></Group></MonsterSettings>
<SpawnSettings><Group Index="7"><Spawn StartHour="12" StarMinute="0" AppearanceRate="100"/></Group></SpawnSettings>
</MonsterGroupRegenSystem>`

func loadGroup(t *testing.T, body string) (*monsterGroupManager, error) {
	t.Helper()
	base := t.TempDir()
	if err := os.WriteFile(filepath.Join(base, "IGC_MonsterGroupRegen.xml"), []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	m := &monsterGroupManager{}
	err := m.load(base)
	return m, err
}
func TestGroupConfigValidation(t *testing.T) {
	for _, tt := range []struct {
		name, body string
		bad        bool
	}{
		{"valid", groupXML, false},
		{"reference", strings.Replace(groupXML, `Index="7"`, `Index="8"`, 1), true},
		{"rate", strings.Replace(groupXML, `AppearanceRate="100"`, `AppearanceRate="101"`, 1), true},
		{"map", strings.Replace(groupXML, `MapNumber="0"`, `MapNumber="999"`, 1), true},
		{"missing boss", strings.Replace(groupXML, `BossMonsterIndex="44"`, `BossMonsterIndex="55"`, 1), true},
		{"count", strings.Replace(groupXML, `Count="2"`, `Count="256"`, 1), true},
		{"duplicate time", strings.Replace(groupXML, `<Spawn StartHour="12" StarMinute="0" AppearanceRate="100"/>`, `<Spawn StartHour="12" StarMinute="0" AppearanceRate="100"/><Spawn StartHour="12" StarMinute="0" AppearanceRate="50"/>`, 1), true},
		{"missing override fields", strings.Replace(groupXML, `OverrideDefaultSettings="0"`, `OverrideDefaultSettings="1"`, 1), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			m, err := loadGroup(t, tt.body)
			if (err != nil) != tt.bad {
				t.Fatalf("load: %v", err)
			}
			if !tt.bad && (len(m.Groups) != 1 || m.Groups[0].Index != 7) {
				t.Fatal("lost sparse group index")
			}
		})
	}
}
func TestGroupRefreshOverridesAndRollback(t *testing.T) {
	body := strings.Replace(groupXML, `OverrideDefaultSettings="0"`, `OverrideDefaultSettings="1" EventID="26" HP="12345" DamageMin="100" DamageMax="200" Defense="50" AttackRate="60" DefenseRate="70"`, 1)
	m, err := loadGroup(t, body)
	if err != nil {
		t.Fatal(err)
	}
	g := m.Groups[0]
	now := time.Now()
	t.Cleanup(func() {
		for _, obj := range g.monsters {
			object.ObjectManager.DeleteEventMonster(obj, time.Now())
		}
	})
	if err := g.Start(now); err != nil {
		t.Fatal(err)
	}
	if len(g.monsters) != 2 || !g.Running() {
		t.Fatal("group did not spawn")
	}
	old := append([]*object.Object(nil), g.monsters...)
	for _, obj := range old {
		if !obj.NoRegen || obj.HP != 12345 || obj.MaxHP != 12345 || obj.AttackMin != 100 || obj.AttackMax != 200 || obj.EventBagID == nil || *obj.EventBagID != 26 {
			t.Fatal("missing override")
		}
	}
	g.Tick(now.Add(2 * time.Hour))
	if !g.Running() {
		t.Fatal("reserved Duration unexpectedly expired group")
	}
	if err := g.Start(now.Add(3 * time.Hour)); err != nil {
		t.Fatal(err)
	}
	for _, obj := range old {
		if object.ObjectManager.GetObject(obj.Index) == obj {
			t.Fatal("old group still registered")
		}
	}
	g.members = append(g.members, member{class: 99999, count: 1})
	before := map[*object.Object]bool{}
	for i := 0; i < conf.Server.GameServerInfo.MaxMonsterCount; i++ {
		if obj := object.ObjectManager.GetObject(i); obj != nil {
			before[obj] = true
		}
	}
	if err := g.Start(now.Add(4 * time.Hour)); err == nil {
		t.Fatal("expected partial failure")
	}
	for i := 0; i < conf.Server.GameServerInfo.MaxMonsterCount; i++ {
		if obj := object.ObjectManager.GetObject(i); obj != nil && !before[obj] {
			t.Fatal("partial group leaked")
		}
	}
}
