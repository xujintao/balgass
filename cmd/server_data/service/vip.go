package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/xujintao/balgass/cmd/server_data/db"
	"github.com/xujintao/balgass/cmd/server_data/model"
)

type vipManager struct{}

func (v *vipManager) VIPAdd(index interface{}, req *model.VipAddReq) error {
	stmt := db.Lookup("vip-find-account-update")
	var count1, count2 int
	if err := db.DBMuOnline.QueryRowx(stmt, req.Username, time.Now()).Scan(&count1, &count2); err != nil {
		return fmt.Errorf("%s, %v", stmt, err)
	}
	vip := model.VIP{
		Username: req.Username,
		Date:     time.Now().Add(time.Duration(req.Days) * time.Hour * 24),
		Type:     int(req.Type),
	}
	switch {
	case count1 == 0:
		// insert
		stmt = db.Lookup("vip-insert")
		if _, err := db.DBMuOnline.NamedExec(stmt, &vip); err != nil {
			return fmt.Errorf("%s, %v", stmt, err)
		}
	case count1 == 1 && count2 == 0:
		// update
		stmt = db.Lookup("vip-update-renew")
		if _, err := db.DBMuOnline.NamedExec(stmt, &vip); err != nil {
			return fmt.Errorf("%s, %v", stmt, err)
		}
	case count1 == 1 && count2 == 0:
		// update
		stmt = db.Lookup("vip-update-upgrade")
		if _, err := db.DBMuOnline.NamedExec(stmt, &vip); err != nil {
			return fmt.Errorf("%s, %v", stmt, err)
		}
	}

	return nil
}

func (v *vipManager) VIPCheck(index interface{}, req *model.VipCheckReq) (*model.VipCheckRes, error) {
	res := &model.VipCheckRes{
		Username: req.Username,
		PayCode:  5,
	}

	stmt := db.Lookup("vip-find-account")
	var vip model.VIP
	if err := db.DBMuOnline.Get(&vip, stmt, req.Username, time.Now()); err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return nil, fmt.Errorf("%s, %v", stmt, err)
	}
	return res, nil
}
