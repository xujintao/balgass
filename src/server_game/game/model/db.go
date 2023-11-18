package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xujintao/balgass/src/server_game/game/item"
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
	dsn := "postgres://root:Kdk7yTkCsvfvvEWg3d3H@localhost/game?sslmode=disable"
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("gorm.Open failed [err]%v", err)
	}
	if err := gdb.AutoMigrate(
		&Account{},
		&Character{},
	); err != nil {
		log.Panicf("gorm.AutoMigrate failed [err]%v", err)
	}
	if os.Getenv("DEBUG") == "1" {
		db.DB = gdb.Debug()
	} else {
		db.DB = gdb
	}
}

type PositionedItems struct {
	size  int
	Items []*item.Item
}

func (pi *PositionedItems) MarshalJSON() ([]byte, error) {
	var items []*item.Item
	for i, v := range pi.Items {
		if v == nil {
			continue
		}
		v.Position = i
		items = append(items, v)
	}
	data, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (pi *PositionedItems) UnmarshalJSON(buf []byte) error {
	var items []*item.Item
	err := json.Unmarshal(buf, &items)
	if err != nil {
		return err
	}
	pi.Items = make([]*item.Item, pi.size)
	for _, v := range items {
		v.Code = item.Code(v.Section, v.Index)
		itemBase, err := item.ItemTable.GetItemBase(v.Section, v.Index)
		if err != nil {
			return err
		}
		v.ItemBase = itemBase
		pi.Items[v.Position] = v
	}
	return nil
}

func (pi *PositionedItems) Scan(value any) error {
	buf, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to Scan Inventory value:", value))
	}
	return pi.UnmarshalJSON(buf)
}

func (pi PositionedItems) Value() (driver.Value, error) {
	return pi.MarshalJSON()
}

type Warehouse struct {
	PositionedItems
}

func (w *Warehouse) Scan(value any) error {
	w.size = 240
	return w.PositionedItems.Scan(value)
}

type Account struct {
	ID                 int          `json:"id,omitempty" gorm:"primarykey"`
	Name               string       `json:"name" validate:"required,max=10,min=1,ascii" gorm:"unique"`
	Password           string       `json:"password,omitempty" validate:"required,max=10,min=1"`
	Characters         []*Character `json:"-" validate:"-"`
	UserID             int          `json:"user_id" validate:"required"`
	Warehouse          Warehouse    `json:"warehouse,omitempty" validate:"-" gorm:"type:jsonb"`
	WarehouseExpansion int          `json:"warehouse_expansion,omitempty" validate:"-"`
	Money              int          `json:"money,omitempty" validate:"-"`
	CreatedAt          time.Time    `json:"-"`
	UpdatedAt          time.Time    `json:"-"`
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

func (db *db) GetAccountList(uid int) ([]*Account, error) {
	var accs []*Account
	var err error
	if uid == 0 {
		err = db.Order("id ASC").
			Omit("Password").
			Find(&accs).
			Error
	} else {
		err = db.Order("id ASC").
			Where("user_id = ?", uid).
			Omit("Password").
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

type Inventory [237]*item.Item

func (i *Inventory) Scan(value any) error {
	buf, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to Scan Inventory value:", value))
	}
	var inventory []*item.Item
	err := json.Unmarshal(buf, &inventory)
	if err != nil {
		return err
	}
	for _, v := range inventory {
		v.Code = item.Code(v.Section, v.Index)
		itemBase, err := item.ItemTable.GetItemBase(v.Section, v.Index)
		if err != nil {
			return err
		}
		v.ItemBase = itemBase
		i[v.Position] = v
	}
	return nil
}

func (i Inventory) Value() (driver.Value, error) {
	var inventory []*item.Item
	for i, v := range i {
		if v == nil {
			continue
		}
		v.Position = i
		inventory = append(inventory, v)
	}
	data, err := json.Marshal(inventory)
	if err != nil {
		return nil, err
	}
	return data, err
}

type Character struct {
	ID                 int       `json:"id" gorm:"primarykey"`
	AccountID          int       `json:"-" validate:"-" gorm:"not null"`
	Position           int       `json:"position" validate:"-" gorm:"not null"`
	Name               string    `json:"name" validate:"required,max=10,min=1,ascii" gorm:"unique;not null"`
	Class              int       `json:"class" validate:"-" gorm:"not null"`
	ChangeUp           int       `json:"change_up" validate:"-" gorm:"not null"`
	Level              int       `json:"level" validate:"-" gorm:"not null"`
	LevelUpPoint       int       `json:"level_up_point,omitempty" validate:"-" gorm:"not null"`
	MapNumber          int       `json:"map_number,omitempty" validate:"-" gorm:"not null"`
	X                  int       `json:"x,omitempty" validate:"-" gorm:"not null"`
	Y                  int       `json:"y,omitempty" validate:"-" gorm:"not null"`
	Dir                int       `json:"dir,omitempty" validate:"-" gorm:"not null"`
	Strength           int       `json:"strength,omitempty" validate:"-" gorm:"not null"`
	Dexterity          int       `json:"dexterity,omitempty" validate:"-" gorm:"not null"`
	Vitality           int       `json:"vitality,omitempty" validate:"-" gorm:"not null"`
	Energy             int       `json:"energy,omitempty" validate:"-" gorm:"not null"`
	Leadership         int       `json:"leadership,omitempty" validate:"-" gorm:"not null"`
	Inventory          Inventory `json:"inventory" validate:"-" gorm:"type:jsonb;not null"`
	InventoryExpansion int       `json:"inventory_expansion,omitempty" validate:"-" gorm:"not null"`
	Money              int       `json:"money,omitempty" validate:"-" gorm:"not null"`
	Experience         int       `json:"experience,omitempty" validate:"-" gorm:"not null"`
	CreatedAt          time.Time `json:"-"`
	UpdatedAt          time.Time `json:"-"`
}

func (db *db) CreateCharacter(c *Character) error {
	result := db.FirstOrCreate(c, &Character{Name: c.Name})
	if result.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

func (db *db) UpdateCharacter(name string, c *Character) error {
	return db.Model(c).
		Where("name = ?", name).
		Select("*").Omit(
		"ID",
		"AccountID",
		"Position",
		"Name",
		"Class",
		"CreatedAt").
		Updates(c).Error
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
