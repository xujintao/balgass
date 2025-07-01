package model

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/skill"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	DB.init()
}

var DB db

type db struct {
	*gorm.DB
}

func (db *db) init() {
	// connect postgress
	slog.Info("gorm connect postgress", "host", conf.Server.GameServerInfo.DBHost, "name", conf.Server.GameServerInfo.DBName)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Server.GameServerInfo.DBUser,
		conf.Server.GameServerInfo.DBPassword,
		conf.Server.GameServerInfo.DBHost,
		conf.Server.GameServerInfo.DBPort,
		conf.Server.GameServerInfo.DBName,
	)
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("gorm.Open", "err", err)
		os.Exit(1)
	}

	// auto migrate
	slog.Info("auto migrate")
	if err := gdb.AutoMigrate(
		&Account{},
		&Character{},
	); err != nil {
		slog.Error("gdb.AutoMigrate", "err", err)
		os.Exit(1)
	}
	if conf.ServerEnv.Debug {
		db.DB = gdb.Debug()
	} else {
		db.DB = gdb
	}
}

type Account struct {
	ID                 int            `json:"id,omitempty" gorm:"primarykey"`
	Name               string         `json:"name" validate:"required,max=10,min=1,ascii" gorm:"unique"`
	Password           string         `json:"password,omitempty" validate:"required,max=10,min=1"`
	Characters         []*Character   `json:"characters" validate:"-"`
	UserEmail          string         `json:"user_email" validate:"required"`
	Warehouse          item.Warehouse `json:"warehouse,omitempty" validate:"-" gorm:"type:jsonb"`
	WarehouseExpansion int            `json:"warehouse_expansion,omitempty" validate:"-"`
	WarehouseMoney     int            `json:"warehouse_money,omitempty" validate:"-"`
	CreatedAt          time.Time      `json:"-"`
	UpdatedAt          time.Time      `json:"-"`
}

// https://stackoverflow.com/questions/66135691/automatically-delete-not-null-constraints-on-unused-columns-after-migration
func (db *db) CreateAccount(p *Account) error {
	result := db.FirstOrCreate(p, &Account{Name: p.Name})
	if result.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	p.Password = ""
	return nil
}

func (db *db) UpdateAccountWarehouse(a *Account) error {
	return db.Model(a).
		Select("Warehouse", "WarehouseMoney").
		Updates(a).Error
}

func (db *db) GetAccountByName(name string) (*Account, error) {
	p := Account{}
	err := db.Where(&Account{
		Name: name,
	}).First(&p).Error
	return &p, err
}

func (db *db) GetAccountByID(id int) (*Account, error) {
	p := Account{ID: id}
	err := db.First(&p).Error
	return &p, err
}

func (db *db) GetAccountList(email string) ([]*Account, error) {
	var accs []*Account
	var err error
	query := db.Preload("Characters").Order("id ASC").Omit("Password")
	if email == "" {
		err = query.
			Find(&accs).
			Error
	} else {
		err = query.Where("user_email = ?", email).
			Find(&accs).
			Error
	}
	if err != nil {
		return nil, err
	}
	return accs, nil
}

func (db *db) DeleteAccount(id int) error {
	return db.Delete(&Account{ID: id}).Error
}

type Character struct {
	ID                 int            `json:"id" gorm:"primarykey"`
	AccountID          int            `json:"-" validate:"-" gorm:"not null"`
	Position           int            `json:"position" validate:"-" gorm:"not null"`
	Name               string         `json:"name" validate:"required,max=10,min=1,ascii" gorm:"unique"`
	Class              int            `json:"class"`
	ChangeUp           int            `json:"change_up"`
	Level              int            `json:"level"`
	LevelPoint         int            `json:"level_point,omitempty"`
	Experience         int            `json:"experience,omitempty"`
	Strength           int            `json:"strength,omitempty"`
	Dexterity          int            `json:"dexterity,omitempty"`
	Vitality           int            `json:"vitality,omitempty"`
	Energy             int            `json:"energy,omitempty"`
	Leadership         int            `json:"leadership,omitempty"`
	MasterLevel        int            `json:"master_level"`
	MasterPoint        int            `json:"master_point,omitempty"`
	MasterExperience   int            `json:"master_experience,omitempty"`
	HP                 int            `json:"hp,omitempty"`
	MP                 int            `json:"mp,omitempty"`
	LevelHP            float32        `json:"level_hp,omitempty" gorm:"-"`
	LevelMP            float32        `json:"level_mp,omitempty" gorm:"-"`
	VitalityHP         float32        `json:"vitality_hp,omitempty" gorm:"-"`
	EnergyMP           float32        `json:"energy_mp,omitempty" gorm:"-"`
	Skills             skill.Skills   `json:"skills,omitempty" validate:"-" gorm:"type:jsonb;default:'[]'"`
	Inventory          item.Inventory `json:"inventory,omitempty" validate:"-" gorm:"type:jsonb;default:'[]'"`
	InventoryExpansion int            `json:"inventory_expansion,omitempty"`
	MuKey              MsgMuKey       `json:"mu_key,omitempty" validate:"-" gorm:"type:jsonb;default:'{}'"`
	MuBot              MsgMuBot       `json:"mu_bot,omitempty" validate:"-" gorm:"type:jsonb;default:'{}'"`
	Money              int            `json:"money,omitempty"`
	MapNumber          int            `json:"map_number"`
	X                  int            `json:"x,omitempty"`
	Y                  int            `json:"y,omitempty"`
	Dir                int            `json:"dir,omitempty"`
	CreatedAt          time.Time      `json:"-"`
	UpdatedAt          time.Time      `json:"-"`
}

func (db *db) CreateCharacter(c *Character) error {
	result := db.FirstOrCreate(c, &Character{Name: c.Name})
	if result.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

func (db *db) UpdateCharacter(c *Character) error {
	return db.Model(c).
		Select("*").Omit(
		"ID",
		"AccountID",
		"Position",
		"Name",
		"Class",
		"MuKey",
		"MuBot",
		"CreatedAt").
		Updates(c).Error
}

func (db *db) UpdateCharacterMuKey(id int, key *MsgDefineMuKey) error {
	return db.Model(&Character{ID: id}).Update("MuKey", key).Error
}

func (db *db) UpdateCharacterMuBot(id int, bot *MsgDefineMuBot) error {
	return db.Model(&Character{ID: id}).Update("MuBot", bot).Error
}

func (db *db) GetCharacterList(aid int) ([]*Character, error) {
	var chars []*Character
	err := db.Order("position ASC").Where("account_id = ?", aid).Find(&chars).Error
	return chars, err
}

func (db *db) GetCharacterByName(aid int, name string) (*Character, error) {
	c := Character{}
	err := db.First(&c, &Character{
		AccountID: aid,
		Name:      name,
	}).Error
	return &c, err
}

func (db *db) DeleteCharacterByName(aid int, name string) error {
	return db.Delete(&Character{}, &Character{
		AccountID: aid,
		Name:      name,
	}).Error
}
