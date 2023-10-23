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
	ID        int            `json:"id,omitempty" gorm:"primarykey"`
	Name      string         `json:"name" validate:"required,max=10,min=1,ascii" gorm:"unique;not null"`
	Password  string         `json:"password,omitempty" validate:"required,max=10,min=1" gorm:"not null"`
	Mail      string         `json:"mail" validate:"required,email" gorm:"unique;not null"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// https://stackoverflow.com/questions/66135691/automatically-delete-not-null-constraints-on-unused-columns-after-migration
func (db *db) CreateAccount(p *Account) (*Account, error) {
	err := db.Create(p).Error
	if err != nil {
		return nil, err
	}
	p.Password = ""
	return p, nil
}

func (db *db) GetAccountByName(name string) (*Account, error) {
	p := Account{}
	err := db.Where(&Account{
		Name: name,
	}).First(&p).Error
	return &p, err
}

func (db *db) GetAccountListByMail(mail string) ([]*Account, error) {
	return nil, nil
}

func (db *db) GetAccountList() ([]*Account, error) {
	return nil, nil
}

func (db *db) DeleteAccount(id int) error {
	return db.Delete(&Account{ID: id}).Error
}
