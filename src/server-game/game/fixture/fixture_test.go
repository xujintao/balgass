package fixture

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/model"
)

type testGame struct {
	commands chan testCommand
}

type testCommand struct {
	name string
	msg  any
}

func newTestGame() *testGame {
	return &testGame{
		commands: make(chan testCommand, 10),
	}
}

func (g *testGame) Command(name string, msg any) (any, error) {
	g.commands <- testCommand{name: name, msg: msg}
	return nil, nil
}

func TestStartFixtureMissingFile(t *testing.T) {
	g := newTestGame()
	store := newMemoryFixtureStore()
	file := filepath.Join(t.TempDir(), testBotFixtureFile)

	if err := startFixture(file, g, store, testCharacterTemplate); err != nil {
		t.Fatalf("startFixture() error = %v", err)
	}
	if store.accountUpserts != 0 || store.characterUpserts != 0 {
		t.Fatalf("unexpected upserts: accounts=%d characters=%d", store.accountUpserts, store.characterUpserts)
	}
	assertNoCommand(t, g)
}

func TestStartFixtureDisabled(t *testing.T) {
	g := newTestGame()
	store := newMemoryFixtureStore()
	file := writeFixture(t, `<Fixture Enable="false" Mode="reset">
		<Accounts>
			<Account Name="account1" Password="password" UserEmail="bot@test.local"/>
		</Accounts>
		<Characters>
			<Character Account="account1" Name="bot1" Class="0"/>
		</Characters>
		<Bots>
			<Bot Account="account1" Password="password" Name="bot1"/>
		</Bots>
	</Fixture>`)

	if err := startFixture(file, g, store, testCharacterTemplate); err != nil {
		t.Fatalf("startFixture() error = %v", err)
	}
	if store.accountUpserts != 0 || store.characterUpserts != 0 {
		t.Fatalf("unexpected upserts: accounts=%d characters=%d", store.accountUpserts, store.characterUpserts)
	}
	assertNoCommand(t, g)
}

func TestStartFixtureResetOverwritesExistingData(t *testing.T) {
	g := newTestGame()
	store := newMemoryFixtureStore()
	store.accounts["account1"] = &model.Account{
		ID:        1,
		Name:      "account1",
		Password:  "old-password",
		UserEmail: "old@test.local",
	}
	store.characters["bot1"] = &model.Character{
		ID:        2,
		AccountID: 1,
		Name:      "bot1",
		Class:     0,
		Level:     99,
	}
	file := writeFixture(t, `<Fixture Enable="true" Mode="reset">
		<Accounts>
			<Account Name="account1" Password="password" UserEmail="bot@test.local"/>
		</Accounts>
		<Characters>
			<Character Account="account1" Name="bot1" Class="0" Level="5"/>
		</Characters>
		<Bots>
			<Bot Account="account1" Password="password" Name="bot1"/>
		</Bots>
	</Fixture>`)

	if err := startFixture(file, g, store, testCharacterTemplate); err != nil {
		t.Fatalf("startFixture() error = %v", err)
	}
	if got := store.accounts["account1"].Password; got != "password" {
		t.Fatalf("account password = %q, want password", got)
	}
	if got := store.accounts["account1"].UserEmail; got != "bot@test.local" {
		t.Fatalf("account user email = %q, want bot@test.local", got)
	}
	if got := store.characters["bot1"].Level; got != 5 {
		t.Fatalf("character level = %d, want 5", got)
	}
	command := waitCommand(t, g)
	if command.name != "AddBot" {
		t.Fatalf("command name = %q, want AddBot", command.name)
	}
}

func TestStartFixtureGrowKeepsExistingData(t *testing.T) {
	g := newTestGame()
	store := newMemoryFixtureStore()
	store.accounts["account1"] = &model.Account{
		ID:        1,
		Name:      "account1",
		Password:  "old-password",
		UserEmail: "old@test.local",
	}
	store.characters["bot1"] = &model.Character{
		ID:        2,
		AccountID: 1,
		Name:      "bot1",
		Class:     0,
		Level:     99,
	}
	file := writeFixture(t, `<Fixture Enable="true" Mode="grow">
		<Accounts>
			<Account Name="account1" Password="password" UserEmail="bot@test.local"/>
		</Accounts>
		<Characters>
			<Character Account="account1" Name="bot1" Class="0" Level="5"/>
		</Characters>
		<Bots>
			<Bot Account="account1" Password="password" Name="bot1"/>
		</Bots>
	</Fixture>`)

	if err := startFixture(file, g, store, testCharacterTemplate); err != nil {
		t.Fatalf("startFixture() error = %v", err)
	}
	if got := store.accounts["account1"].Password; got != "old-password" {
		t.Fatalf("account password = %q, want old-password", got)
	}
	if got := store.characters["bot1"].Level; got != 99 {
		t.Fatalf("character level = %d, want 99", got)
	}
	command := waitCommand(t, g)
	if command.name != "AddBot" {
		t.Fatalf("command name = %q, want AddBot", command.name)
	}
}

func TestParseTestBotFixtureValidation(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "unknown mode",
			in: `<Fixture Enable="true" Mode="seed">
				<Accounts>
					<Account Name="account1" Password="password" UserEmail="bot@test.local"/>
				</Accounts>
				<Characters>
					<Character Account="account1" Name="bot1" Class="0"/>
				</Characters>
				<Bots>
					<Bot Account="account1" Password="password" Name="bot1"/>
				</Bots>
			</Fixture>`,
			want: "unknown test bot fixture mode",
		},
		{
			name: "character missing account",
			in: `<Fixture Enable="true" Mode="reset">
				<Accounts>
					<Account Name="account1" Password="password" UserEmail="bot@test.local"/>
				</Accounts>
				<Characters>
					<Character Account="missing" Name="bot1" Class="0"/>
				</Characters>
				<Bots>
					<Bot Account="account1" Password="password" Name="bot1"/>
				</Bots>
			</Fixture>`,
			want: "character references missing account",
		},
		{
			name: "bot missing character",
			in: `<Fixture Enable="true" Mode="reset">
				<Accounts>
					<Account Name="account1" Password="password" UserEmail="bot@test.local"/>
				</Accounts>
				<Characters>
					<Character Account="account1" Name="bot1" Class="0"/>
				</Characters>
				<Bots>
					<Bot Account="account1" Password="password" Name="missing"/>
				</Bots>
			</Fixture>`,
			want: "bot references missing character",
		},
		{
			name: "duplicated bot",
			in: `<Fixture Enable="true" Mode="reset">
				<Accounts>
					<Account Name="account1" Password="password" UserEmail="bot@test.local"/>
				</Accounts>
				<Characters>
					<Character Account="account1" Name="bot1" Class="0"/>
				</Characters>
				<Bots>
					<Bot Account="account1" Password="password" Name="bot1"/>
					<Bot Account="ACCOUNT1" Password="password" Name="BOT1"/>
				</Bots>
			</Fixture>`,
			want: "duplicated test bot fixture bot",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseTestBotFixture([]byte(tt.in), testCharacterTemplate)
			if err == nil {
				t.Fatal("parseTestBotFixture() error = nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("parseTestBotFixture() error = %q, want containing %q", err, tt.want)
			}
		})
	}
}

func TestParseTestBotFixtureDefaultsModeAndPosition(t *testing.T) {
	fixture, err := parseTestBotFixture([]byte(`<Fixture Enable="true">
		<Accounts>
			<Account Name="account1" Password="password" UserEmail="bot@test.local"/>
		</Accounts>
		<Characters>
			<Character Account="account1" Name="bot1" Class="0"/>
			<Character Account="account1" Name="bot2" Class="0" Position="3"/>
			<Character Account="account1" Name="bot3" Class="0"/>
		</Characters>
		<Bots>
			<Bot Account="account1" Password="password" Name="bot1"/>
		</Bots>
	</Fixture>`), testCharacterTemplate)
	if err != nil {
		t.Fatalf("parseTestBotFixture() error = %v", err)
	}
	if fixture.Mode != testBotModeReset {
		t.Fatalf("mode = %q, want reset", fixture.Mode)
	}
	if got := fixture.Characters[0].Character.Position; got != 0 {
		t.Fatalf("bot1 position = %d, want 0", got)
	}
	if got := fixture.Characters[1].Character.Position; got != 3 {
		t.Fatalf("bot2 position = %d, want 3", got)
	}
	if got := fixture.Characters[2].Character.Position; got != 4 {
		t.Fatalf("bot3 position = %d, want 4", got)
	}
}

func writeFixture(t *testing.T, content string) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), testBotFixtureFile)
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	return file
}

func testCharacterTemplate(class int) (model.Character, bool) {
	if class != 0 {
		return model.Character{}, false
	}
	return model.Character{
		Class:     class,
		Level:     1,
		MapNumber: 0,
		X:         125,
		Y:         125,
		HP:        100,
		MP:        50,
	}, true
}

type memoryFixtureStore struct {
	accounts         map[string]*model.Account
	characters       map[string]*model.Character
	nextID           int
	accountUpserts   int
	characterUpserts int
}

func newMemoryFixtureStore() *memoryFixtureStore {
	return &memoryFixtureStore{
		accounts:   make(map[string]*model.Account),
		characters: make(map[string]*model.Character),
		nextID:     1,
	}
}

func (s *memoryFixtureStore) UpsertFixtureAccount(account model.Account, reset bool) (*model.Account, error) {
	s.accountUpserts++
	key := fixtureKey(account.Name)
	existing, ok := s.accounts[key]
	if !ok {
		account.ID = s.nextID
		s.nextID++
		accountCopy := account
		s.accounts[key] = &accountCopy
		return &accountCopy, nil
	}
	if reset {
		existing.Password = account.Password
		existing.UserEmail = account.UserEmail
		existing.Warehouse = account.Warehouse
		existing.WarehouseExpansion = account.WarehouseExpansion
		existing.WarehouseMoney = account.WarehouseMoney
	}
	return existing, nil
}

func (s *memoryFixtureStore) UpsertFixtureCharacter(character model.Character, reset bool) error {
	s.characterUpserts++
	key := fixtureKey(character.Name)
	existing, ok := s.characters[key]
	if !ok {
		character.ID = s.nextID
		s.nextID++
		characterCopy := character
		s.characters[key] = &characterCopy
		return nil
	}
	if !reset {
		if existing.AccountID != character.AccountID {
			return errors.New("character account mismatch")
		}
		return nil
	}
	character.ID = existing.ID
	character.Name = existing.Name
	characterCopy := character
	s.characters[key] = &characterCopy
	return nil
}

func waitCommand(t *testing.T, g *testGame) testCommand {
	t.Helper()
	select {
	case command := <-g.commands:
		return command
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for Command")
	}
	return testCommand{}
}

func assertNoCommand(t *testing.T, g *testGame) {
	t.Helper()
	select {
	case command := <-g.commands:
		t.Fatalf("unexpected command = %#v", command)
	case <-time.After(50 * time.Millisecond):
	}
}
