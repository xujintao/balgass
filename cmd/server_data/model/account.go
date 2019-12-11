package model

type Account struct {
	ID        int    `db:"memb_guid"`
	Passwd    string `db:"memb__pwd"`
	BlockCode int    `db:"bloc_code"`
}
