package model

import "time"

type User struct {
	ID        int    `db:"memb_guid"`
	Username  string `db:"memb___id"`
	Password  string `db:"memb__pwd"`
	BlockCode int    `db:"bloc_code"`
}

type UserLoginHistory struct {
	Username   string    `db:"AccountID"`
	ServerName string    `db:"ServerName"`
	IP         string    `db:"IP"`
	Date       time.Time `db:"Date"`
	State      string    `db:"State"`
	MID        string    `db:"HWID"`
}

type UserState struct {
	Username       string    `db:"memb___id"`
	State          int       `db:"ConnectStat"`
	ServerName     string    `db:"ServerName"`
	IP             string    `db:"IP"`
	ConnectTime    time.Time `db:"ConnectTM"`
	DisConnectTime time.Time `db:"DisConnectTM"`
}

type UserChar struct {
	Charname1          string `db:"GameID1"`
	Charname2          string `db:"GameID2"`
	Charname3          string `db:"GameID3"`
	Charname4          string `db:"GameID4"`
	Charname5          string `db:"GameID5"`
	MoveCount          int    `db:"MoveCnt"`
	EnableSummoner     bool   `db:"Summoner"`
	EnableRageFighter  bool   `db:"RageFighter"`
	WarehouseExtension int    `db:"WarehouseExpansion"`
	SecurityCode       int    `db:"SecCode"`
}

type VIP struct {
	Username string    `db:"AccountID"`
	Date     time.Time `db:"Date"`
	Type     int       `db:"Type"`
}

type CharInfo struct {
	Name          string `db:"Name"`
	Class         int    `db:"Class"`
	Level         int    `db:"cLevel"`
	LevelMaster   int    `db:"mLevel"`
	CtlCode       int    `db:"CtlCode"`
	Resets        int    `db:"RESETS"`
	InventoryItem []byte `db:"Inventory"`
}
