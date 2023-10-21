package model

import (
	"log"
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
	); err != nil {
		log.Panicf("gorm.AutoMigrate failed [err]%v", err)
	}
	db.DB = gdb.Debug()
}

type Account struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Account   string         `json:"account" gorm:"unique;not null"`
	Password  string         `json:"password,omitempty" gorm:"not null"`
	WebID     int            `json:"web_id" gorm:"unique;not null"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (db *db) CreateAccount(p *Account) (*Account, error) {
	err := db.Create(p).Error
	if err != nil {
		return nil, err
	}
	return db.GetAccountByAccount(p.Account)
}

func (db *db) GetAccountByAccount(account string) (*Account, error) {
	p := Account{}
	err := db.Where(&Account{
		Account: account,
	}).First(&p).Error
	return &p, err
}

func (db *db) DeleteAccount(account string) error {
	if account == "test_account" {
		return db.Unscoped().Where("account = ?", account).Delete(&Account{}).Error

	} else {
		return db.Where("account = ?", account).Delete(&Account{}).Error
	}
}
