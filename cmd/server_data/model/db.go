package model

import "time"

type Account struct {
	ID        int    `db:"memb_guid"`
	Passwd    string `db:"memb__pwd"`
	BlockCode int    `db:"bloc_code"`
}

type AccountLoginHistory struct {
	UserName   string    `db:"AccountID"`
	ServerName string    `db:"ServerName"`
	IP         string    `db:"IP"`
	Date       time.Time `db:"Date"`
	State      string    `db:"State"`
	MID        string    `db:"HWID"`
}

type AccountState struct {
	UserName       string    `db:"memb___id"`
	State          int       `db:"ConnectStat"`
	ServerName     string    `db:"ServerName"`
	IP             string    `db:"IP"`
	ConnectTime    time.Time `db:"ConnectTM"`
	DisConnectTime time.Time `db:"DisConnectTM"`
}
