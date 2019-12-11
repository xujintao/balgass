package db

import (
	"fmt"
	"log"
	"net/url"

	_ "github.com/denisenkom/go-mssqldb" // sqlserver drive
	"github.com/jmoiron/sqlx"
	"github.com/xujintao/balgass/cmd/server_data/conf"
)

var xdb *sqlx.DB

var (
	DBMuOnline   *sqlx.DB
	DBEvent      *sqlx.DB
	DBRank       *sqlx.DB
	DBBattleCore *sqlx.DB
)

func init() {
	// "sqlserver://username:password@host:port?param1=value&param2=value"
	// "sqlserver://username:password@host/instance?param1=value&param2=value"

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(conf.SQL.User, conf.SQL.Pass),
		Host:   fmt.Sprintf("%s:%d", conf.SQL.SQLServerName, 1433),
		// Path:  instance, // if connecting to an instance instead of a port
		// RawQuery: query.Encode(),
	}
	type config struct {
		name string
		db   *sqlx.DB
	}
	configs := [...]config{
		{"MuOnline", DBMuOnline},
		{"Events", DBEvent},
		{"Ranking", DBRank},
		{"BattleCore", DBBattleCore},
	}
	for _, c := range configs {
		query := url.Values{}
		query.Set("database", c.name)
		u.RawQuery = query.Encode()
		c.db = newDB("sqlserver", u.String())
	}

	log.Println("db connected")
}

func newDB(driverName, dsn string) (db *sqlx.DB) {
	db, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(10)
	return
}
