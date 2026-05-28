package fixture

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/bot"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/object/player"
	"gorm.io/gorm"
)

func init() {
	Fixture.init()
}

const (
	ModeReset   string = "reset"
	ModeGrow    string = "grow"
	ModeDefault string = ModeReset
)

var Fixture fixture

type fixture struct {
	Enable     bool
	Mode       string
	Accounts   []*model.Account
	Characters []*struct {
		Account   string
		Character model.Character
	}
	Bots []*model.MsgAddBot
}

func (f *fixture) init() {
	// Fixture was generated 2026-05-26 17:44:23 by https://xml-to-go.github.io/ in Ukraine.
	type FixtureConfig struct {
		XMLName  xml.Name `xml:"Fixture"`
		Text     string   `xml:",chardata"`
		Enable   string   `xml:"Enable,attr"`
		Mode     string   `xml:"Mode,attr"`
		Accounts struct {
			Text    string `xml:",chardata"`
			Account []struct {
				Text      string `xml:",chardata"`
				Name      string `xml:"Name,attr"`
				Password  string `xml:"Password,attr"`
				UserEmail string `xml:"UserEmail,attr"`
			} `xml:"Account"`
		} `xml:"Accounts"`
		Characters struct {
			Text      string `xml:",chardata"`
			Character []struct {
				Text     string `xml:",chardata"`
				Account  string `xml:"Account,attr"`
				Name     string `xml:"Name,attr"`
				Class    string `xml:"Class,attr"`
				Position string `xml:"Position,attr"`
			} `xml:"Character"`
		} `xml:"Characters"`
		Bots struct {
			Text string `xml:",chardata"`
			Bot  []struct {
				Text     string `xml:",chardata"`
				Account  string `xml:"Account,attr"`
				Password string `xml:"Password,attr"`
				Name     string `xml:"Name,attr"`
			} `xml:"Bot"`
		} `xml:"Bots"`
	}
	var fc FixtureConfig
	conf.XML(conf.ServerEnv.PathCommon, "fixture/fixture.xml", &fc)

	// Convert FixtureConfig to fixture.
	// enable
	f.Enable = strings.EqualFold(strings.TrimSpace(fc.Enable), "true")
	if !f.Enable {
		return
	}

	// mode
	mode := strings.ToLower(strings.TrimSpace(fc.Mode))
	if mode == "" {
		mode = ModeDefault
	}
	switch mode {
	case ModeReset, ModeGrow:
		f.Mode = mode
	default:
		slog.Error("unknown fixture mode", "mode", fc.Mode)
		f.Enable = false
		return
	}

	// accounts
	accounts := make(map[string]struct{}, len(fc.Accounts.Account))
	f.Accounts = make([]*model.Account, 0, len(fc.Accounts.Account))
	for _, v := range fc.Accounts.Account {
		account := &model.Account{
			Name:      strings.TrimSpace(v.Name),
			Password:  strings.TrimSpace(v.Password),
			UserEmail: strings.TrimSpace(v.UserEmail),
		}
		if account.Name == "" || account.Password == "" || account.UserEmail == "" {
			slog.Error("fixture account requires name, password and user_email", "account", account.Name)
			continue
		}
		key := fixtureKey(account.Name)
		if _, ok := accounts[key]; ok {
			slog.Error("duplicated fixture account", "account", account.Name)
			continue
		}
		accounts[key] = struct{}{}
		f.Accounts = append(f.Accounts, account)
	}

	// characters
	characters := make(map[string]struct{}, len(fc.Characters.Character))
	accountCharacters := make(map[string]struct{}, len(fc.Characters.Character))
	nextPosition := make(map[string]int)
	f.Characters = make([]*struct {
		Account   string
		Character model.Character
	}, 0, len(fc.Characters.Character))
	for _, v := range fc.Characters.Character {
		account := strings.TrimSpace(v.Account)
		name := strings.TrimSpace(v.Name)
		classText := strings.TrimSpace(v.Class)
		if account == "" || name == "" || classText == "" {
			slog.Error("fixture character requires account, name and class", "account", account, "character", name)
			continue
		}
		accountKey := fixtureKey(account)
		if _, ok := accounts[accountKey]; !ok {
			slog.Error("fixture character references missing account", "account", account, "character", name)
			continue
		}
		nameKey := fixtureKey(name)
		if _, ok := characters[nameKey]; ok {
			slog.Error("duplicated fixture character", "character", name)
			continue
		}
		accountCharacterKey := accountKey + ":" + nameKey
		if _, ok := accountCharacters[accountCharacterKey]; ok {
			slog.Error("duplicated fixture account character", "account", account, "character", name)
			continue
		}

		class, err := strconv.Atoi(classText)
		if err != nil {
			slog.Error("invalid fixture character class", "account", account, "character", name, "class", v.Class, "err", err)
			continue
		}
		character, ok := player.CharacterTable[class]
		if !ok {
			slog.Error("fixture character class not found", "account", account, "character", name, "class", class)
			continue
		}
		character.ID = 0
		character.AccountID = 0
		character.Name = name
		character.Class = class
		if positionText := strings.TrimSpace(v.Position); positionText != "" {
			position, err := strconv.Atoi(positionText)
			if err != nil {
				slog.Error("invalid fixture character position", "account", account, "character", name, "position", v.Position, "err", err)
				continue
			}
			character.Position = position
		} else {
			character.Position = nextPosition[accountKey]
		}
		if character.Position >= nextPosition[accountKey] {
			nextPosition[accountKey] = character.Position + 1
		}
		characters[nameKey] = struct{}{}
		accountCharacters[accountCharacterKey] = struct{}{}
		f.Characters = append(f.Characters, &struct {
			Account   string
			Character model.Character
		}{
			Account:   account,
			Character: character,
		})
	}

	// bots
	bots := make(map[string]struct{}, len(fc.Bots.Bot))
	f.Bots = make([]*model.MsgAddBot, 0, len(fc.Bots.Bot))
	for _, v := range fc.Bots.Bot {
		bot := &model.MsgAddBot{
			Account:  strings.TrimSpace(v.Account),
			Password: strings.TrimSpace(v.Password),
			Name:     strings.TrimSpace(v.Name),
		}
		if bot.Account == "" || bot.Password == "" || bot.Name == "" {
			slog.Error("fixture bot requires account, password and name", "account", bot.Account, "character", bot.Name)
			continue
		}
		accountKey := fixtureKey(bot.Account)
		if _, ok := accounts[accountKey]; !ok {
			slog.Error("fixture bot references missing account", "account", bot.Account, "character", bot.Name)
			continue
		}
		if _, ok := accountCharacters[accountKey+":"+fixtureKey(bot.Name)]; !ok {
			slog.Error("fixture bot references missing character", "account", bot.Account, "character", bot.Name)
			continue
		}
		key := botKey(bot.Account, bot.Name)
		if _, ok := bots[key]; ok {
			slog.Error("duplicated fixture bot", "bot", key)
			continue
		}
		bots[key] = struct{}{}
		f.Bots = append(f.Bots, bot)
	}
}

func (f *fixture) Start() error {
	if !f.Enable {
		slog.Info("fixture is disabled")
		return nil
	}
	slog.Info("starting fixture", "mode", f.Mode)

	// accounts
	reset := f.Mode == ModeReset
	store := databaseFixtureStore{}
	accountIDs := make(map[string]int, len(f.Accounts))
	for _, account := range f.Accounts {
		if account == nil {
			continue
		}
		saved, err := store.UpsertFixtureAccount(*account, reset)
		if err != nil {
			slog.Error("fixture account seed failed", "err", err, "account", account.Name)
			continue
		}
		accountIDs[fixtureKey(saved.Name)] = saved.ID
	}

	// characters
	for _, c := range f.Characters {
		if c == nil {
			continue
		}
		accountID, ok := accountIDs[fixtureKey(c.Account)]
		if !ok {
			slog.Error("fixture character seed skipped", "account", c.Account, "character", c.Character.Name)
			continue
		}
		character := c.Character
		character.AccountID = accountID
		if err := store.UpsertFixtureCharacter(character, reset); err != nil {
			slog.Error("fixture character seed failed", "err", err, "account", c.Account, "character", character.Name)
		}
	}

	// bots
	for _, b := range f.Bots {
		if b == nil {
			continue
		}
		if _, err := bot.BotManager.AddBot(b); err != nil {
			slog.Error("test bot fixture add bot failed", "err", err, "account", b.Account, "character", b.Name)
		}
	}
	return nil
}

func fixtureKey(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func botKey(account, name string) string {
	return fixtureKey(account) + ":" + fixtureKey(name)
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
