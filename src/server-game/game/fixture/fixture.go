package fixture

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object/player"
	"gorm.io/gorm"
)

const (
	testBotFixtureFile = "fixture.xml"
	testBotModeReset   = "reset"
	testBotModeGrow    = "grow"
)

// Game is the command surface needed to start fixture bots.
type Game interface {
	Command(string, any) (any, error)
}

// Start loads configured test bot seed data and starts bots when enabled.
func Start(game Game) {
	file := path.Join(conf.ServerEnv.PathCommon, "fixture", testBotFixtureFile)
	err := startFixture(file, game, databaseFixtureStore{}, playerCharacterTemplate)
	if err != nil {
		slog.Error("test bot fixture failed", "err", err, "file", file)
	}
}

type testBotFixture struct {
	Enable     bool
	Mode       string
	Accounts   []model.Account
	Characters []fixtureCharacter
	Bots       []model.MsgAddBot
}

type rawTestBotFixture struct {
	XMLName    xml.Name              `xml:"Fixture"`
	Enable     bool                  `xml:"Enable,attr"`
	Mode       string                `xml:"Mode,attr"`
	Accounts   []rawFixtureAccount   `xml:"Accounts>Account"`
	Characters []rawFixtureCharacter `xml:"Characters>Character"`
	Bots       []rawFixtureBot       `xml:"Bots>Bot"`
}

type rawFixtureAccount struct {
	Name      string `xml:"Name,attr"`
	Password  string `xml:"Password,attr"`
	UserEmail string `xml:"UserEmail,attr"`
}

type rawFixtureCharacter struct {
	Account  string `xml:"Account,attr"`
	Name     string `xml:"Name,attr"`
	Class    *int   `xml:"Class,attr"`
	Position *int   `xml:"Position,attr"`
	Level    *int   `xml:"Level,attr"`
}

type rawFixtureBot struct {
	Account  string `xml:"Account,attr"`
	Password string `xml:"Password,attr"`
	Name     string `xml:"Name,attr"`
}

type fixtureCharacter struct {
	Account   string
	Character model.Character
}

type fixtureTemplateFunc func(class int) (model.Character, bool)

type fixtureStore interface {
	UpsertFixtureAccount(account model.Account, reset bool) (*model.Account, error)
	UpsertFixtureCharacter(character model.Character, reset bool) error
}

func startFixture(file string, game Game, store fixtureStore, template fixtureTemplateFunc) error {
	fixture, err := loadTestBotFixture(file, template)
	if err != nil {
		return err
	}
	if fixture == nil {
		return nil
	}
	if !fixture.Enable {
		return nil
	}

	reset := fixture.Mode == testBotModeReset
	accountIDs := make(map[string]int, len(fixture.Accounts))
	for _, account := range fixture.Accounts {
		saved, err := store.UpsertFixtureAccount(account, reset)
		if err != nil {
			slog.Error("test bot fixture account seed failed", "err", err, "account", account.Name)
			continue
		}
		accountIDs[fixtureKey(saved.Name)] = saved.ID
	}

	for _, character := range fixture.Characters {
		accountID, ok := accountIDs[fixtureKey(character.Account)]
		if !ok {
			slog.Error("test bot fixture character seed skipped", "account", character.Account, "character", character.Character.Name)
			continue
		}
		c := character.Character
		c.AccountID = accountID
		if err := store.UpsertFixtureCharacter(c, reset); err != nil {
			slog.Error("test bot fixture character seed failed", "err", err, "account", character.Account, "character", c.Name)
		}
	}

	for _, bot := range fixture.Bots {
		if _, err := game.Command("AddBot", &bot); err != nil {
			slog.Error("test bot fixture add bot failed", "err", err, "account", bot.Account, "character", bot.Name)
		}
	}
	return nil
}

func loadTestBotFixture(file string, template fixtureTemplateFunc) (*testBotFixture, error) {
	if _, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	raw := rawTestBotFixture{}
	conf.XML(path.Dir(file), path.Base(file), &raw)
	return parseRawTestBotFixture(raw, template)
}

func parseTestBotFixture(buf []byte, template fixtureTemplateFunc) (*testBotFixture, error) {
	raw := rawTestBotFixture{}
	if err := xml.Unmarshal(buf, &raw); err != nil {
		return nil, err
	}
	return parseRawTestBotFixture(raw, template)
}

func parseRawTestBotFixture(raw rawTestBotFixture, template fixtureTemplateFunc) (*testBotFixture, error) {
	mode := strings.ToLower(strings.TrimSpace(raw.Mode))
	if mode == "" {
		mode = testBotModeReset
	}
	if !raw.Enable {
		return &testBotFixture{
			Enable: raw.Enable,
			Mode:   mode,
		}, nil
	}
	switch mode {
	case testBotModeReset, testBotModeGrow:
	default:
		return nil, fmt.Errorf("unknown test bot fixture mode: %s", raw.Mode)
	}

	accounts, err := parseFixtureAccounts(raw.Accounts)
	if err != nil {
		return nil, err
	}
	characters, err := parseFixtureCharacters(raw.Characters, accounts, template)
	if err != nil {
		return nil, err
	}
	bots, err := parseFixtureBots(raw.Bots, accounts, characters)
	if err != nil {
		return nil, err
	}

	return &testBotFixture{
		Enable:     raw.Enable,
		Mode:       mode,
		Accounts:   accounts,
		Characters: characters,
		Bots:       bots,
	}, nil
}

func parseFixtureAccounts(raw []rawFixtureAccount) ([]model.Account, error) {
	seen := make(map[string]struct{}, len(raw))
	accounts := make([]model.Account, 0, len(raw))
	for _, rawAccount := range raw {
		account := model.Account{
			Name:      strings.TrimSpace(rawAccount.Name),
			Password:  strings.TrimSpace(rawAccount.Password),
			UserEmail: strings.TrimSpace(rawAccount.UserEmail),
		}
		if account.Name == "" || account.Password == "" || account.UserEmail == "" {
			return nil, errors.New("test bot fixture account requires name, password and user_email")
		}
		key := fixtureKey(account.Name)
		if _, ok := seen[key]; ok {
			return nil, fmt.Errorf("duplicated test bot fixture account: %s", account.Name)
		}
		seen[key] = struct{}{}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func parseFixtureCharacters(raw []rawFixtureCharacter, accounts []model.Account, template fixtureTemplateFunc) ([]fixtureCharacter, error) {
	accountSet := make(map[string]struct{}, len(accounts))
	for _, account := range accounts {
		accountSet[fixtureKey(account.Name)] = struct{}{}
	}
	seenName := make(map[string]struct{}, len(raw))
	seenAccountName := make(map[string]struct{}, len(raw))
	nextPosition := make(map[string]int)
	characters := make([]fixtureCharacter, 0, len(raw))
	for _, rawCharacter := range raw {
		rawCharacter.Account = strings.TrimSpace(rawCharacter.Account)
		rawCharacter.Name = strings.TrimSpace(rawCharacter.Name)
		if rawCharacter.Account == "" || rawCharacter.Name == "" || rawCharacter.Class == nil {
			return nil, errors.New("test bot fixture character requires account, name and class")
		}
		accountKey := fixtureKey(rawCharacter.Account)
		if _, ok := accountSet[accountKey]; !ok {
			return nil, fmt.Errorf("test bot fixture character references missing account: %s", rawCharacter.Account)
		}
		nameKey := fixtureKey(rawCharacter.Name)
		if _, ok := seenName[nameKey]; ok {
			return nil, fmt.Errorf("duplicated test bot fixture character: %s", rawCharacter.Name)
		}
		seenName[nameKey] = struct{}{}
		accountNameKey := accountKey + ":" + nameKey
		if _, ok := seenAccountName[accountNameKey]; ok {
			return nil, fmt.Errorf("duplicated test bot fixture account character: %s:%s", rawCharacter.Account, rawCharacter.Name)
		}
		seenAccountName[accountNameKey] = struct{}{}

		character, ok := template(*rawCharacter.Class)
		if !ok {
			return nil, fmt.Errorf("test bot fixture character class not found: %d", *rawCharacter.Class)
		}
		character.ID = 0
		character.AccountID = 0
		character.Name = rawCharacter.Name
		character.Class = *rawCharacter.Class
		if rawCharacter.Level != nil {
			character.Level = *rawCharacter.Level
		}
		if rawCharacter.Position == nil {
			character.Position = nextPosition[accountKey]
		} else {
			character.Position = *rawCharacter.Position
		}
		if character.Position >= nextPosition[accountKey] {
			nextPosition[accountKey] = character.Position + 1
		}
		characters = append(characters, fixtureCharacter{
			Account:   rawCharacter.Account,
			Character: character,
		})
	}
	return characters, nil
}

func parseFixtureBots(raw []rawFixtureBot, accounts []model.Account, characters []fixtureCharacter) ([]model.MsgAddBot, error) {
	accountSet := make(map[string]struct{}, len(accounts))
	for _, account := range accounts {
		accountSet[fixtureKey(account.Name)] = struct{}{}
	}
	characterSet := make(map[string]struct{}, len(characters))
	for _, character := range characters {
		characterSet[fixtureKey(character.Account)+":"+fixtureKey(character.Character.Name)] = struct{}{}
	}

	seen := make(map[string]struct{}, len(raw))
	bots := make([]model.MsgAddBot, 0, len(raw))
	for _, rawBot := range raw {
		bot := model.MsgAddBot{
			Account:  strings.TrimSpace(rawBot.Account),
			Password: strings.TrimSpace(rawBot.Password),
			Name:     strings.TrimSpace(rawBot.Name),
		}
		if bot.Account == "" || bot.Password == "" || bot.Name == "" {
			return nil, errors.New("test bot fixture bot requires account, password and name")
		}
		accountKey := fixtureKey(bot.Account)
		if _, ok := accountSet[accountKey]; !ok {
			return nil, fmt.Errorf("test bot fixture bot references missing account: %s", bot.Account)
		}
		if _, ok := characterSet[accountKey+":"+fixtureKey(bot.Name)]; !ok {
			return nil, fmt.Errorf("test bot fixture bot references missing character: %s:%s", bot.Account, bot.Name)
		}
		key := botKey(bot.Account, bot.Name)
		if _, ok := seen[key]; ok {
			return nil, fmt.Errorf("duplicated test bot fixture bot: %s", key)
		}
		seen[key] = struct{}{}
		bots = append(bots, bot)
	}
	return bots, nil
}

func fixtureKey(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func botKey(account, name string) string {
	return fixtureKey(account) + ":" + fixtureKey(name)
}

func playerCharacterTemplate(class int) (model.Character, bool) {
	character, ok := player.CharacterTable[class]
	return character, ok
}

type databaseFixtureStore struct{}

func (databaseFixtureStore) UpsertFixtureAccount(account model.Account, reset bool) (*model.Account, error) {
	existing := model.Account{}
	err := model.DB.Where(&model.Account{Name: account.Name}).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			account.ID = 0
			if err := model.DB.Create(&account).Error; err != nil {
				return nil, err
			}
			return &account, nil
		}
		return nil, err
	}
	if reset {
		if err := model.DB.Model(&existing).
			Select("Password", "UserEmail", "Warehouse", "WarehouseExpansion", "WarehouseMoney").
			Updates(&account).Error; err != nil {
			return nil, err
		}
		existing.Password = account.Password
		existing.UserEmail = account.UserEmail
		existing.Warehouse = account.Warehouse
		existing.WarehouseExpansion = account.WarehouseExpansion
		existing.WarehouseMoney = account.WarehouseMoney
	}
	return &existing, nil
}

func (databaseFixtureStore) UpsertFixtureCharacter(character model.Character, reset bool) error {
	existing := model.Character{}
	err := model.DB.Where(&model.Character{Name: character.Name}).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			character.ID = 0
			return model.DB.Create(&character).Error
		}
		return err
	}
	if !reset {
		if existing.AccountID != character.AccountID {
			return fmt.Errorf("test bot fixture character %s belongs to account id %d, want %d",
				character.Name, existing.AccountID, character.AccountID)
		}
		return nil
	}
	return model.DB.Model(&existing).
		Select("*").
		Omit("ID", "Name", "CreatedAt").
		Updates(&character).Error
}
