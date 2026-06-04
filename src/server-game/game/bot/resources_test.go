package bot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadResourcesClassifiesAttackableMonstersConservatively(t *testing.T) {
	dir := t.TempDir()
	writeTestResources(t, dir, []byte{0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0}, `
<MonsterSpawn>
  <Map Number="0">
    <Spot Type="0"><Spawn Index="100" StartX="1" StartY="1"/></Spot>
    <Spot Type="1">
      <Spawn Index="1" StartX="1" StartY="1" EndX="2" EndY="2"/>
      <Spawn Index="100" StartX="1" StartY="1" EndX="2" EndY="2"/>
    </Spot>
  </Map>
</MonsterSpawn>`, `<ShopList><Shop NPCIndex="2" MapNumber="0" PosX="1" PosY="1"/></ShopList>`)

	resources, err := loadResources(dir)
	if err != nil {
		t.Fatalf("loadResources() error = %v", err)
	}
	if !resources.attackable(1) {
		t.Fatal("class 1 attackable = false, want true")
	}
	if resources.attackable(100) {
		t.Fatal("class 100 attackable = true, want false for NPC conflict")
	}
	if resources.attackable(2) {
		t.Fatal("class 2 attackable = true, want false for shop NPC")
	}
}

func TestLoadResourcesRejectsInvalidTerrainSize(t *testing.T) {
	dir := t.TempDir()
	writeTestResources(t, dir, []byte{0, 1, 1, 0}, `<MonsterSpawn/>`, `<ShopList/>`)

	if _, err := loadResources(dir); err == nil {
		t.Fatal("loadResources() error = nil, want invalid terrain size")
	}
}

func TestLoadResourcesRejectsMissingSpawnPosition(t *testing.T) {
	dir := t.TempDir()
	writeTestResources(t, dir, []byte{0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0}, `
<MonsterSpawn>
  <Map Number="0">
    <Spot Type="1"><Spawn Index="1" StartX="1" EndX="2" EndY="2"/></Spot>
  </Map>
</MonsterSpawn>`, `<ShopList/>`)

	if _, err := loadResources(dir); err == nil {
		t.Fatal("loadResources() error = nil, want missing spawn position")
	}
}

func writeTestResources(t *testing.T, dir string, terrain []byte, monsterSpawn, shops string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, "MapTerrains"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "Monsters"), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string][]byte{
		"IGC_MapList.xml":               []byte(`<MapList><DefaultMaps><Map Number="0" File="test.att"/></DefaultMaps></MapList>`),
		"MapTerrains/test.att":          terrain,
		"Monsters/IGC_MonsterSpawn.xml": []byte(monsterSpawn),
		"IGC_ShopList.xml":              []byte(shops),
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
