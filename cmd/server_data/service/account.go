package service

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/xujintao/balgass/cmd/server_data/db"
	"github.com/xujintao/balgass/cmd/server_data/model"
	"github.com/xujintao/balgass/network"
)

type account struct {
	username    string
	machineID   string
	number      uint32
	addr        string
	serverIndex string
	offTrade    bool
	tick        int64
}

type accountManager struct {
	mu       sync.RWMutex
	accounts map[string]*account
}

func (m *accountManager) accountKick(acc *account, force bool) error {
	msgRes := &model.AccountKickRes{
		Username: acc.username,
		Number:   acc.number,
		Force:    force,
	}
	buf, err := proto.Marshal(msgRes)
	if err != nil {
		return fmt.Errorf("Marshal failed, %v", err)
	}

	res := &network.Response{}
	res.WriteHead(0xc1, 0x07).Write(buf)
	if err := ServerManager.Send(acc.serverIndex, res); err != nil {
		return fmt.Errorf("Send failed, %v", err)
	}

	return nil
}

func (m *accountManager) accountJoinOther(index string, acc *account) error {
	msgRes := &model.AccountJoinOtherRes{
		Username: acc.username,
	}
	buf, err := proto.Marshal(msgRes)
	if err != nil {
		return fmt.Errorf("Marshal failed, %v", err)
	}

	res := &network.Response{}
	res.WriteHead(0xc1, 0x08).Write(buf)
	if err := ServerManager.Send(index, res); err != nil {
		return fmt.Errorf("Send failed, %v", err)
	}

	return nil
}

func (m *accountManager) accountLoginTrack(data *model.AccountLoginHistory) {
	stmt := db.Lookup("account_login_history-insert")
	db.DBMuOnline.Exec(stmt, data)
}

func (m *accountManager) AccountLogin(index string, req *model.AccountLoginReq) (*model.AccountLoginRes, error) {
	res := &model.AccountLoginRes{
		Username: req.Username,
		Number:   req.Number,
	}
	username := req.Username
	accDB := &model.Account{}

	// verify username
	stmt := db.Lookup("account-find-passwd")
	if err := db.DBMuOnline.Get(accDB, stmt, username); err != nil {
		if err == sql.ErrNoRows {
			res.Result = 2 // here, 0 is more safe
			return res, nil
		}
		return nil, fmt.Errorf("%s %v", stmt, err)
	}

	// verify password
	if accDB.Passwd != req.Password {
		res.Result = 0
		return res, nil
	}

	// verify ban machine
	stmt = db.Lookup("ban_machine-find-count")
	count := 0
	if err := db.DBMuOnline.Get(&count, stmt, req.MachineID); err != nil {
		return nil, fmt.Errorf("%s %v", stmt, err)
	}
	if count > 0 {
		res.Result = 5
		return res, nil
	}

	// verify block code
	if accDB.BlockCode > 0 {
		switch accDB.BlockCode {
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
	acc, ok := m.accounts[username]
	if ok {
		// verify off trade
		if acc.offTrade {
			if err := m.accountKick(acc, true); err != nil {
				return nil, fmt.Errorf("accountKick failed, %v", err)
			}
		}

		// verify move timeout
		if time.Now().Unix()-acc.tick > int64(60*time.Second) {
			if err := m.accountKick(acc, true); err != nil {
				return nil, fmt.Errorf("accountKick failed, %v", err)
			}
		} else {
			if err := m.accountJoinOther(index, acc); err != nil {
				return nil, fmt.Errorf("accountJoinOther failed, %v", err)
			}
			res.Result = 3
			return res, nil
		}
		delete(m.accounts, username)
	}

	// verify machine id use count limit
	count = 0
	for _, acc := range m.accounts {
		if acc.machineID == req.MachineID {
			count++
		}
	}
	if count >= ServerManager.ServerGetMIDUseCount(index) {
		if err := m.accountJoinOther(index, acc); err != nil {
			return nil, fmt.Errorf("accountJoinOther failed, %v", err)
		}
		res.Result = 4
		return res, nil
	}

	// verify machine id use count limit(group)

	// account add
	acc = &account{
		username:    username,
		machineID:   req.MachineID,
		number:      req.Number,
		addr:        req.Addr,
		serverIndex: index,
	}
	m.accounts[username] = acc
	tx := db.DBMuOnline.MustBegin()
	// history
	stmt = db.Lookup("account_Login_history-insert")
	loginHistory := model.AccountLoginHistory{
		UserName:   username,
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
	state := model.AccountState{
		UserName:    username,
		State:       1,
		ServerName:  ServerManager.ServerGetName(index),
		IP:          req.Addr[:strings.Index(req.Addr, ":")],
		ConnectTime: time.Now(),
	}
	stmt = db.Lookup("account_state-find-account")
	if err := tx.Get(nil, stmt, username); err != nil {
		if err != sql.ErrNoRows {
			tx.Rollback()
			return nil, fmt.Errorf("%s %v", stmt, err)
		}
		// insert state
		stmt = db.Lookup("account_state-insert")
		if _, err := tx.NamedExec(stmt, &state); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%s, %v", stmt, err)
		}
	} else {
		// update state
		stmt = db.Lookup("account_state-update")
		if _, err := tx.NamedExec(stmt, &state); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%s, %v", stmt, err)
		}
	}
	tx.Commit()

	return res, nil
}
