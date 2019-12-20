package service

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/xujintao/balgass/cmd/server_data/conf"

	"github.com/xujintao/balgass/cmd/server_data/db"
	"github.com/xujintao/balgass/cmd/server_data/model"
)

const (
	ItemDBEncLen      = 32
	CharBaseItemCount = 12
)

type charManager struct{}

func (m *charManager) CharListGet(index interface{}, req *model.CharListGetReq) (*model.CharListGetRes, error) {
	res := model.CharListGetRes{
		Username: req.Username,
		Number:   req.Number,
	}

	// user_chars
	stmt := db.Lookup("user_chars-find-count")
	count := 0
	if err := db.DBMuOnline.Get(&count, stmt, req.Username); err != nil {
		return nil, fmt.Errorf("%s, %v", stmt, err)
	}
	if count == 0 {
		stmt = db.Lookup("user_chars-insert")
		if _, err := db.DBMuOnline.Exec(stmt, req.Username); err != nil {
			return nil, fmt.Errorf("%s, %v", stmt, err)
		}
	}
	stmt = db.Lookup("user_chars-find-char")
	userChar := model.UserChar{}
	if err := db.DBMuOnline.Get(&userChar, stmt, req.Username); err != nil {
		return nil, fmt.Errorf("%s, %v", stmt, err)
	}
	log.Printf("[%s] - characters: [%s][%s][%s][%s][%s]", req.Username, userChar.Charname1, userChar.Charname2, userChar.Charname3, userChar.Charname4, userChar.Charname5)
	res.WarehouseExtension = int32(userChar.WarehouseExtension)
	res.MoveCount = int32(userChar.MoveCount)
	res.SecurityCode = int32(userChar.SecurityCode)
	if userChar.EnableSummoner {
		res.EnableCharCreate |= 1
	}
	if userChar.EnableRageFighter {
		res.EnableCharCreate |= 8
	}

	// chars
	stmt = db.Lookup("chars-find-chars")
	var chars []model.CharInfo
	if err := db.DBMuOnline.Select(&chars, stmt, req.Username); err != nil {
		return nil, fmt.Errorf("%s, %v", stmt, err)
	}
	if len(chars) != 0 {
		res.Chars = make([]*model.CharListGetResCharInfo, len(chars))
	}
	for _, char := range chars {
		c := model.CharListGetResCharInfo{
			Inventory: make([]byte, 4*CharBaseItemCount),
		}
		c.Level = int32(char.Level + char.LevelMaster)
		if c.Level >= int32(conf.Server.MagicGladiatorCreateMinLevel) {
			res.EnableCharCreate |= 4
		}
		if c.Level >= int32(conf.Server.DarkLordCreateMinLevel) {
			res.EnableCharCreate |= 2
		}
		if c.Level >= int32(conf.Server.GrowLancerCreateMinLevel) {
			res.EnableCharCreate |= 16
		}
		c.Class = int32(char.Class)
		c.Resets = int32(char.Resets)
		c.CtlCode = int32(char.CtlCode)
		for i := 0; i < CharBaseItemCount; i++ {
			if char.InventoryItem[i*ItemDBEncLen+0] == 0xFF && char.InventoryItem[i*ItemDBEncLen+7]&0x80 == 0x80 {
				c.Inventory[i*4+0] = 0xFF
				c.Inventory[i*4+1] = 0xFF
				c.Inventory[i*4+2] = 0xFF
				c.Inventory[i*4+3] = 0xFF
			} else {
				c.Inventory[i*4+0] = char.InventoryItem[i*ItemDBEncLen+0] // 0~7bit of item
				c.Inventory[i*4+1] = char.InventoryItem[i*ItemDBEncLen+1] // level of item
				c.Inventory[i*4+2] = char.InventoryItem[i*ItemDBEncLen+7] // 8th bit(extend) of item
				c.Inventory[i*4+3] = char.InventoryItem[i*ItemDBEncLen+9] // section of item
			}
		}
		c.GuildStatus = 0xFF
		c.ServerCode = 0
		c.BattleName = char.Name
		if req.UnityBattleFieldServer {
			stmt = db.Lookup("battlecore-find-bfname")
			if err := db.DBBattleCore.QueryRowx(char.Name).Scan(&c.BattleName, &c.ServerCode); err != nil && err != sql.ErrNoRows {
				return nil, fmt.Errorf("%s, %v", stmt, err)
			}
		} else {
			stmt = db.Lookup("guild_members-find-status")
			if err := db.DBMuOnline.Get(&c.GuildStatus, stmt, char.Name); err != nil && err != sql.ErrNoRows {
				return nil, fmt.Errorf("%s, %v", stmt, err)
			}
		}
		res.Chars = append(res.Chars, &c)
	}
	return &res, nil
}
