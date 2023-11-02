package model

import (
	"log"
	"os"
	"time"

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

type Account struct {
	ID         int          `json:"id,omitempty" gorm:"primarykey"`
	Name       string       `json:"name" validate:"required,max=10,min=1,ascii" gorm:"unique;not null"`
	Password   string       `json:"password,omitempty" validate:"required,max=10,min=1" gorm:"not null"`
	Characters []*Character `json:"-" validate:"-"`
	UserID     int          `json:"user_id" validate:"required" gorm:"not null"`
	CreatedAt  time.Time    `json:"-"`
	UpdatedAt  time.Time    `json:"-"`
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

func (db *db) GetAccountList(uid int) ([]*Account, error) {
	var accs []*Account
	var err error
	if uid == 0 {
		err = db.Order("id ASC").Find(&accs).Error
	} else {
		err = db.Order("id ASC").Where("user_id = ?", uid).Find(&accs).Error
	}
	if err != nil {
		return nil, err
	}
	for _, acc := range accs {
		acc.Password = ""
	}
	return accs, nil
}

func (db *db) DeleteAccount(id int) error {
	return db.Delete(&Account{ID: id}).Error
}

type Character struct {
	ID        int       `json:"id,omitempty" gorm:"primarykey"`
	AccountID int       `json:"-" validate:"-" gorm:"not null"`
	Position  int       `json:"position" validate:"-" gorm:"not null"`
	Name      string    `json:"name" validate:"required,max=10,min=1,ascii" gorm:"unique;not null"`
	Class     int       `json:"class" validate:"-" gorm:"not null"`
	Level     int       `json:"level" validate:"-" gorm:"not null"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (db *db) CreateCharacter(c *Character) error {
	result := db.FirstOrCreate(c, &Character{Name: c.Name})
	if result.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

func (db *db) GetCharacterList(aid int) ([]*Character, error) {
	var chars []*Character
	err := db.Order("position ASC").Where("account_id = ?", aid).Find(&chars).Error
	if err != nil {
		return nil, err
	}
	return chars, nil
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
