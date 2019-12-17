package model

import "time"

type Account struct {
	ID        int    `db:"memb_guid"`
	Username  string `db:"memb___id"`
	Password  string `db:"memb__pwd"`
	BlockCode int    `db:"bloc_code"`
}

type AccountLoginHistory struct {
	Username   string    `db:"AccountID"`
	ServerName string    `db:"ServerName"`
	IP         string    `db:"IP"`
	Date       time.Time `db:"Date"`
	State      string    `db:"State"`
	MID        string    `db:"HWID"`
}

type AccountState struct {
	Username       string    `db:"memb___id"`
	State          int       `db:"ConnectStat"`
	ServerName     string    `db:"ServerName"`
	IP             string    `db:"IP"`
	ConnectTime    time.Time `db:"ConnectTM"`
	DisConnectTime time.Time `db:"DisConnectTM"`
}

type VIP struct {
	Username string    `db:"AccountID"`
	Date     time.Time `db:"Date"`
	Type     int       `db:"Type"`
}
