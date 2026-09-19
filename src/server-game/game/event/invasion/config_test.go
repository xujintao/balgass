package invasion

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xujintao/balgass/src/server-game/conf"
)

func TestConfigLoadsAndValidates(t *testing.T) {
	const dragon = `<DragonEvent><Monster Index="44" Distance="30" Count="2"><Map Number="0"><Spawn StartX="135" StartY="61" EndX="146" EndY="70"/></Map></Monster></DragonEvent>`
	for _, tt := range []struct {
		name, dragon string
		bad          bool
	}{
		{"valid", dragon, false},
		{"map", strings.Replace(dragon, `Number="0"`, `Number="999"`, 1), true},
		{"area", strings.Replace(dragon, `EndX="146"`, `EndX="135"`, 1), true},
		{"class", strings.Replace(dragon, `Index="44"`, `Index="55"`, 1), true},
		{"count", strings.Replace(dragon, `Count="2"`, `Count="0"`, 1), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			if err := os.Mkdir(filepath.Join(base, "Events"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(base, "Events/IGC_DragonEvent.xml"), []byte(tt.dragon), 0600); err != nil {
				t.Fatal(err)
			}
			var c configuration
			err := c.load(base)
			if (err != nil) != tt.bad {
				t.Fatalf("load error=%v", err)
			}
			if !tt.bad && c.dragon.count != 2 {
				t.Fatal("wrong configured dragon count")
			}
		})
	}
}

func TestStartupConfig(t *testing.T) {
	var c configuration
	if err := c.load(conf.PathCommon); err != nil {
		t.Fatal(err)
	}
	if len(c.dragon.maps) == 0 {
		t.Fatal("startup config empty")
	}
}
