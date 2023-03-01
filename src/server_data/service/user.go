package service

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/golang/protobuf/proto"
	"github.com/xujintao/balgass/cmd/server_data/db"
	"github.com/xujintao/balgass/cmd/server_data/model"
	"github.com/xujintao/balgass/network"
)

type userServerMoveData struct {
	serverMoveReq *model.UserServerMoveReq
	serverMoveRes *model.UserServerMoveRes
	movetime      int64
}

type user struct {
	username       string
	machineID      string
	number         uint32
	addr           string
	serverIndex    interface{}
	offTrade       bool
	serverMoveData *userServerMoveData
}

func (a *user) getIP() string {
	addr := a.addr
	ip := addr[:strings.Index(addr, ":")]
	return ip
}

type userManager struct {
	mu    sync.RWMutex
	users map[string]*user
}

func (m *userManager) UserDelete(index interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for username, v := range m.users {
		if v.serverIndex == index {
			delete(m.users, username)
		}
	}
}

func (m *userManager) userKickRes(u *user, force bool) error {
	msgRes := &model.UserKickRes{
		Username: u.username,
		Number:   u.number,
		Force:    force,
	}
	buf, err := proto.Marshal(msgRes)
	if err != nil {
		return fmt.Errorf("Marshal failed, %v", err)
	}

	res := &network.Response{}
	res.WriteHead(0xc1, 0x07).Write(buf)
	if err := ServerManager.Send(u.serverIndex, res); err != nil {
		return fmt.Errorf("Send failed, %v", err)
	}

	return nil
}

func (m *userManager) userJoinOtherRes(index interface{}, u *user) error {
	sid := index.(string)
	msgRes := &model.UserJoinOtherRes{
		Username: u.username,
	}
	buf, err := proto.Marshal(msgRes)
	if err != nil {
		return fmt.Errorf("Marshal failed, %v", err)
	}

	res := &network.Response{}
	res.WriteHead(0xc1, 0x08).Write(buf)
	if err := ServerManager.Send(sid, res); err != nil {
		return fmt.Errorf("Send failed, %v", err)
	}

	return nil
}

func (m *userManager) userStateUpdate(tx *sqlx.Tx, username string) error {
	if tx == nil {
		tx = db.DBMuOnline.MustBegin()
		defer tx.Commit() // return err, but ok
	}
	state := model.UserState{
		Username:       username,
		State:          0,
		DisConnectTime: time.Now(),
	}
	stmt := db.Lookup("user_states-find-user")
	num := 0
	if err := tx.Get(&num, stmt, username); err != nil {
		tx.Rollback()
		return fmt.Errorf("%s %v", stmt, err)
	}
	if num == 1 {
		// update state
		stmt = db.Lookup("user_states-update-disconn")
		if _, err := tx.NamedExec(stmt, &state); err != nil {
			tx.Rollback()
			return fmt.Errorf("%s, %v", stmt, err)
		}
	}
	return nil
}

func (m *userManager) UserLogin(index interface{}, req *model.UserLoginReq) (*model.UserLoginRes, error) {
	res := &model.UserLoginRes{
		Username: req.Username,
		Number:   req.Number,
	}
	username := req.Username
	uDB := &model.User{}

	// verify username
	stmt := db.Lookup("users-find-user")
	if err := db.DBMuOnline.Get(uDB, stmt, username); err != nil {
		if err == sql.ErrNoRows {
			res.Result = 2 // here, 0 is more safe
			return res, nil
		}
		return nil, fmt.Errorf("%s %v", stmt, err)
	}

	// verify password
	if uDB.Password != req.Password {
		res.Result = 0
		return res, nil
	}

	// verify ban machine
	stmt = db.Lookup("ban_machines-find-count")
	count := 0
	if err := db.DBMuOnline.Get(&count, stmt, req.MachineID); err != nil {
		return nil, fmt.Errorf("%s %v", stmt, err)
	}
	if count > 0 {
		res.Result = 5
		return res, nil
	}

	// verify block code
	if uDB.BlockCode > 0 {
		switch uDB.BlockCode {
		case 0x41:
			res.Result = 14
		case 0x46:
			res.Result = 15
		case 0x49:
			res.Result = 17
		default:
			res.Result = 5
		}
		return res, nil // need return ?
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[username]
	if ok {
		// verify off trade
		if u.offTrade {
			if err := m.userKickRes(u, true); err != nil {
				return nil, fmt.Errorf("userKick failed, %v", err)
			}
		}

		// verify server move timeout
		if u.serverMoveData != nil {
			if time.Now().Unix()-u.serverMoveData.movetime < int64(60*time.Second) {
				// keep server move state until timeout
				if err := m.userJoinOtherRes(index, u); err != nil {
					return nil, fmt.Errorf("userJoinOtherRes failed, %v", err)
				}
				res.Result = 3
				return res, nil
			}
			if err := m.userKickRes(u, true); err != nil {
				return nil, fmt.Errorf("userKick failed, %v", err)
			}
		}

		delete(m.users, username)
	}

	// verify machine id use count limit
	count = 0
	for _, u := range m.users {
		if u.machineID == req.MachineID {
			count++
		}
	}
	if count >= ServerManager.ServerGetMIDUseCount(index) {
		if err := m.userJoinOtherRes(index, u); err != nil {
			return nil, fmt.Errorf("userJoinOtherRes failed, %v", err)
		}
		res.Result = 4
		return res, nil
	}

	// verify machine id use count limit(group)
	// ...

	tx := db.DBMuOnline.MustBegin()
	// history
	stmt = db.Lookup("user_Login_historys-insert")
	loginHistory := model.UserLoginHistory{
		Username:   username,
		ServerName: ServerManager.ServerGetName(index),
		IP:         req.Addr[:strings.Index(req.Addr, ":")],
		State:      "Connected",
		MID:        req.MachineID,
	}
	if _, err := tx.NamedExec(stmt, &loginHistory); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s, %v", stmt, err)
	}
	// state
	state := model.UserState{
		Username:    username,
		State:       1,
		ServerName:  ServerManager.ServerGetName(index),
		IP:          req.Addr[:strings.Index(req.Addr, ":")],
		ConnectTime: time.Now(),
	}
	stmt = db.Lookup("user_states-find-user")
	num := 0
	if err := tx.Get(&num, stmt, username); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s %v", stmt, err)
	}
	if num == 0 {
		// insert state
		stmt = db.Lookup("user_states-insert")
		if _, err := tx.NamedExec(stmt, &state); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%s, %v", stmt, err)
		}
	} else {
		// update state
		stmt = db.Lookup("user_states-update-connect")
		if _, err := tx.NamedExec(stmt, &state); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%s, %v", stmt, err)
		}
	}
	tx.Commit()

	// user add
	u = &user{
		username:    username,
		machineID:   req.MachineID,
		number:      req.Number,
		addr:        req.Addr,
		serverIndex: index,
	}
	m.users[username] = u
	return res, nil
}

func (m *userManager) UserLoginFailed(index interface{}, req *model.UserLoginFailedReq) error {
	username := req.Username
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[username]
	if !ok {
		return fmt.Errorf("user not exist")
	}

	tx := db.DBMuOnline.MustBegin()
	// history
	stmt := db.Lookup("user_Login_historys-insert")
	loginHistory := model.UserLoginHistory{
		Username:   username,
		ServerName: ServerManager.ServerGetName(index),
		IP:         u.getIP(),
		State:      "Login Failed",
		MID:        "N/A",
	}
	if _, err := tx.NamedExec(stmt, &loginHistory); err != nil {
		tx.Rollback()
		return fmt.Errorf("%s, %v", stmt, err)
	}
	// state
	if err := m.userStateUpdate(tx, username); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	delete(m.users, username)
	return nil
}

func (m *userManager) UserExit(index interface{}, req *model.UserExitReq) error {
	username := req.Username
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[username]
	if !ok {
		return fmt.Errorf("user not exist")
	}

	tx := db.DBMuOnline.MustBegin()
	// history
	stmt := db.Lookup("user_Login_historys-insert")
	loginHistory := model.UserLoginHistory{
		Username:   username,
		ServerName: ServerManager.ServerGetName(index),
		IP:         u.getIP(),
		State:      "Login Failed",
		MID:        "N/A",
	}
	if _, err := tx.NamedExec(stmt, &loginHistory); err != nil {
		tx.Rollback()
		return fmt.Errorf("%s, %v", stmt, err)
	}
	// state
	if err := m.userStateUpdate(tx, username); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	delete(m.users, username)
	return nil
}

func (m *userManager) UserKill(index interface{}, req *model.UserKillReq) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[req.Username]
	if !ok {
		return fmt.Errorf("user kill failed, %s do not exist", req.Username)
	}

	if err := m.userKickRes(u, true); err != nil {
		return err
	}
	if err := m.userStateUpdate(nil, req.Username); err != nil {
		return err
	}
	delete(m.users, req.Username)
	return nil
}

func (m *userManager) UserBlock(index interface{}, req *model.UserBlockReq) error {
	u := model.User{
		Username:  req.Username,
		BlockCode: 1,
	}
	stmt := db.Lookup("users-update-block")
	if _, err := db.DBMuOnline.NamedExec(stmt, &u); err != nil {
		return fmt.Errorf("%s, %v", stmt, err)
	}
	// delete(m.users, req.Username)
	return nil
}

func (m *userManager) UserOffTradeSet(index interface{}, req *model.UserOffTradeSetReq) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if u, ok := m.users[req.Username]; ok {
		u.offTrade = req.OffTrade
		return nil
	}
	return fmt.Errorf("%s set offtrade failed", req.Username)
}

func (m *userManager) UserServerMove(index interface{}, req *model.UserServerMoveReq) (res *model.UserServerMoveRes, err error) {
	err = nil
	res = &model.UserServerMoveRes{
		Index:         req.Index,
		Username:      req.Username,
		Charname:      req.Charname,
		ServerCodeSrc: req.ServerCodeSrc,
		ServerCodeDst: req.ServerCodeDst,
		MapNumber:     req.MapNumber,
	}
	if !ServerManager.ServerCodeValidate(uint16(req.ServerCodeSrc)) ||
		!ServerManager.ServerCodeValidate(uint16(req.ServerCodeDst)) {
		res.Result = 2
		return
	}
	if ServerManager.ServerUserCountLimit(uint16(req.ServerCodeDst)) {
		res.Result = 3
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[req.Username]
	if !ok {
		res.Result = 1
		return
	}
	if u.serverMoveData != nil {
		res.Result = 4
		return
	}
	u.serverMoveData = &userServerMoveData{
		serverMoveReq: req,
		serverMoveRes: res,
		movetime:      time.Now().Unix(),
	}

	res.AuthCode1 = int32(time.Now().Unix() << 32)
	res.AuthCode2 = int32(time.Now().Unix() % 1000 << 32)
	res.AuthCode3 = int32(time.Now().Unix() % 777 << 32)
	res.AuthCode4 = int32(time.Now().Unix() % 8911 << 32)
	return
}

func (m *userManager) UserServerMoveAuth(index interface{}, req *model.UserServerMoveAuthReq) (res *model.UserServerMoveAuthRes, err error) {
	err = nil
	res = &model.UserServerMoveAuthRes{
		Index:    req.Index,
		Username: req.Username,
		Charname: req.Charname,
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	defer func() {
		if res.Result != 0 {
			delete(m.users, req.Username)
		}
	}()
	u, ok := m.users[req.Username]
	if !ok {
		res.Result = 1
		return
	}
	if u.serverMoveData == nil {
		res.Result = 4
		return
	}
	if req.ServerCode != u.serverMoveData.serverMoveReq.ServerCodeDst {
		res.Result = 3
		return
	}
	if req.AuthCode1 != u.serverMoveData.serverMoveRes.AuthCode1 ||
		req.AuthCode2 != u.serverMoveData.serverMoveRes.AuthCode2 ||
		req.AuthCode3 != u.serverMoveData.serverMoveRes.AuthCode3 ||
		req.AuthCode4 != u.serverMoveData.serverMoveRes.AuthCode4 {
		res.Result = 2
		return
	}
	res.ServerCode = u.serverMoveData.serverMoveReq.ServerCodeSrc
	res.MapNumber = u.serverMoveData.serverMoveReq.MapNumber
	res.SecurityLock = u.serverMoveData.serverMoveReq.SecurityLock
	res.SecurityCode = u.serverMoveData.serverMoveReq.SecurityCode
	res.X = u.serverMoveData.serverMoveReq.X
	res.Y = u.serverMoveData.serverMoveReq.Y
	uDB := model.User{}
	stmt := db.Lookup("users-find-user")
	err = db.DBMuOnline.Get(uDB, stmt, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Result = 1
			return
		}
		return nil, fmt.Errorf("%s %v", stmt, err)
	}
	res.Password = uDB.Username
	u.serverIndex = index
	u.serverMoveData = nil
	return
}
