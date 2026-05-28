package fixture

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xujintao/balgass/src/server-game/conf"
)

func TestFixtureInitDisabled(t *testing.T) {
	cleanup := writeFixtureConfig(t, `<Fixture Enable="false" Mode="reset">
		<Accounts>
			<Account Name="account1" Password="password" UserEmail="bot@test.local"/>
		</Accounts>
		<Characters>
			<Character Account="account1" Name="bot1" Class="0" Position="0"/>
		</Characters>
		<Bots>
			<Bot Account="account1" Password="password" Name="bot1"/>
		</Bots>
	</Fixture>`)
	defer cleanup()

	f := fixture{}
	f.init()

	if f.Enable {
		t.Fatal("fixture enable = true, want false")
	}
	if len(f.Accounts) != 0 || len(f.Characters) != 0 || len(f.Bots) != 0 {
		t.Fatalf("disabled fixture loaded data: accounts=%d characters=%d bots=%d", len(f.Accounts), len(f.Characters), len(f.Bots))
	}
}

func TestFixtureInitParsesValidFixture(t *testing.T) {
	cleanup := writeFixtureConfig(t, `<Fixture Enable="true">
		<Accounts>
			<Account Name=" account1 " Password=" password " UserEmail=" bot@test.local "/>
		</Accounts>
		<Characters>
			<Character Account=" account1 " Name=" bot1 " Class="0"/>
			<Character Account="account1" Name="bot2" Class="0" Position="3"/>
			<Character Account="account1" Name="bot3" Class="0"/>
		</Characters>
		<Bots>
			<Bot Account=" account1 " Password=" password " Name=" bot1 "/>
		</Bots>
	</Fixture>`)
	defer cleanup()

	f := fixture{}
	f.init()

	if !f.Enable {
		t.Fatal("fixture enable = false, want true")
	}
	if f.Mode != ModeDefault {
		t.Fatalf("mode = %q, want %q", f.Mode, ModeDefault)
	}
	if len(f.Accounts) != 1 {
		t.Fatalf("accounts len = %d, want 1", len(f.Accounts))
	}
	if got := f.Accounts[0].Name; got != "account1" {
		t.Fatalf("account name = %q, want account1", got)
	}
	if got := f.Accounts[0].Password; got != "password" {
		t.Fatalf("account password = %q, want password", got)
	}
	if got := f.Accounts[0].UserEmail; got != "bot@test.local" {
		t.Fatalf("account user email = %q, want bot@test.local", got)
	}
	if len(f.Characters) != 3 {
		t.Fatalf("characters len = %d, want 3", len(f.Characters))
	}
	if got := f.Characters[0].Character.Name; got != "bot1" {
		t.Fatalf("character 0 name = %q, want bot1", got)
	}
	if got := f.Characters[0].Character.Position; got != 0 {
		t.Fatalf("character 0 position = %d, want 0", got)
	}
	if got := f.Characters[1].Character.Position; got != 3 {
		t.Fatalf("character 1 position = %d, want 3", got)
	}
	if got := f.Characters[2].Character.Position; got != 4 {
		t.Fatalf("character 2 position = %d, want 4", got)
	}
	if len(f.Bots) != 1 {
		t.Fatalf("bots len = %d, want 1", len(f.Bots))
	}
	if got := f.Bots[0].Account; got != "account1" {
		t.Fatalf("bot account = %q, want account1", got)
	}
	if got := f.Bots[0].Password; got != "password" {
		t.Fatalf("bot password = %q, want password", got)
	}
	if got := f.Bots[0].Name; got != "bot1" {
		t.Fatalf("bot name = %q, want bot1", got)
	}
}

func TestFixtureInitSkipsInvalidEntries(t *testing.T) {
	cleanup := writeFixtureConfig(t, `<Fixture Enable="true" Mode="grow">
		<Accounts>
			<Account Name="account1" Password="password" UserEmail="bot@test.local"/>
			<Account Name="ACCOUNT1" Password="password" UserEmail="dup@test.local"/>
			<Account Name="" Password="password" UserEmail="bad@test.local"/>
		</Accounts>
		<Characters>
			<Character Account="account1" Name="bot1" Class="0" Position="0"/>
			<Character Account="missing" Name="missing-account" Class="0" Position="1"/>
			<Character Account="account1" Name="invalid-class" Class="bad" Position="2"/>
			<Character Account="account1" Name="bot1" Class="0" Position="3"/>
		</Characters>
		<Bots>
			<Bot Account="account1" Password="password" Name="bot1"/>
			<Bot Account="account1" Password="password" Name="missing-character"/>
			<Bot Account="ACCOUNT1" Password="password" Name="BOT1"/>
		</Bots>
	</Fixture>`)
	defer cleanup()

	f := fixture{}
	f.init()

	if !f.Enable {
		t.Fatal("fixture enable = false, want true")
	}
	if f.Mode != ModeGrow {
		t.Fatalf("mode = %q, want grow", f.Mode)
	}
	if len(f.Accounts) != 1 {
		t.Fatalf("accounts len = %d, want 1", len(f.Accounts))
	}
	if len(f.Characters) != 1 {
		t.Fatalf("characters len = %d, want 1", len(f.Characters))
	}
	if got := f.Characters[0].Character.Name; got != "bot1" {
		t.Fatalf("character name = %q, want bot1", got)
	}
	if len(f.Bots) != 1 {
		t.Fatalf("bots len = %d, want 1", len(f.Bots))
	}
	if got := f.Bots[0].Name; got != "bot1" {
		t.Fatalf("bot name = %q, want bot1", got)
	}
}

func TestFixtureKeysNormalizeWhitespaceAndCase(t *testing.T) {
	if got := fixtureKey(" Account1 "); got != "account1" {
		t.Fatalf("fixtureKey() = %q, want account1", got)
	}
	if got := botKey(" Account1 ", " Bot1 "); got != "account1:bot1" {
		t.Fatalf("botKey() = %q, want account1:bot1", got)
	}
}

func writeFixtureConfig(t *testing.T, content string) func() {
	t.Helper()
	pathCommon := conf.ServerEnv.PathCommon
	dir := filepath.Join(t.TempDir(), "fixture")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "fixture.xml"), []byte(content), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	conf.ServerEnv.PathCommon = filepath.Dir(dir)
	return func() {
		conf.ServerEnv.PathCommon = pathCommon
	}
}
