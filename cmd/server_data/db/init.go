package db

import (
	"log"

	_ "github.com/denisenkom/go-mssqldb" // mssql drive
	"github.com/jmoiron/sqlx"
)

var xdb *sqlx.DB

func init() {
	var dsn string
	db, err := sqlx.Connect("mssql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(10)
	log.Println("db connected")
	xdb = db
}
