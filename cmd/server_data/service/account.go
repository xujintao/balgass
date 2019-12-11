package service

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"

	"github.com/xujintao/balgass/cmd/server_data/db"
	"github.com/xujintao/balgass/network"

	"github.com/xujintao/balgass/cmd/server_data/model"
)

type account struct {
	username    string
	machineID   int
	number      uint32
	addr        string
	serverIndex string
}

type accountManager struct {
	mu       sync.RWMutex
	accounts map[string]*account
}

func (m *accountManager) AccountKickLocked(acc *account, force bool) error {
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
	delete(m.accounts, acc.serverIndex)
	return nil
}

func (m *accountManager) AccountAuth(index string, req *model.AccountAuthReq) (*model.AccountAuthRes, error) {
	res := &model.AccountAuthRes{
		Username: req.Username,
		Number:   req.Number,
	}
	username := req.Username
	account := &model.Account{}
	stmt := db.Lookup("account-find-passwd")
	if err := db.DBMuOnline.Get(account, stmt, username); err != nil {
		if err == sql.ErrNoRows {
			res.Result = 2
			return res, nil
		}
		return nil, fmt.Errorf("%s %v", stmt, err)
	}
	if account.Passwd != req.Password {
		res.Result = 0
		return res, nil
	}
	stmt = db.Lookup("ban-find-count")
	count := 0
	if err := db.DBMuOnline.Get(&count, stmt, req.MachineID); err != nil {
		return nil, fmt.Errorf("%s %v", stmt, err)
	}
	if count > 0 {
		res.Result = 5
		return res, nil
	}

	switch account.BlockCode {
	case 0x41:
		res.Result = 14
	case 0x46:
		res.Result = 15
	case 0x49:
		res.Result = 17
	default:
		res.Result = 5
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if acc, ok := m.accounts[username]; ok {
		if err := m.AccountKickLocked(acc, true); err != nil {
			return nil, fmt.Errorf("AccountKickLocked failed, %v", err)
		}
	}

	return nil, nil
}
