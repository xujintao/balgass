package server

import (
	"testing"

	"github.com/xujintao/balgass/src/server-connect/service/model"
)

type testActioner struct{}

func (testActioner) ServerAction(string, any) {}

func TestRegistrationUpdatesAndUnregisters(t *testing.T) {
	var manager serverManager
	manager.init()
	manager.RegisterActioner(testActioner{})
	manager.Register(&model.MsgRegister{
		ServerState: model.ServerState{Code: 1, Percent: 10},
		Config:      model.ServerConfig{Code: 1, IP: "127.0.0.1", Port: 56900, Visible: false, Name: "Old"},
	})
	if got := manager.GetServerList(); len(got) != 0 {
		t.Fatalf("hidden server listed: %+v", got)
	}
	manager.Register(&model.MsgRegister{
		ServerState: model.ServerState{Code: 1, Percent: 20},
		Config:      model.ServerConfig{Code: 1, IP: "127.0.0.2", Port: 56901, Visible: true, Name: "New"},
	})
	if got := manager.GetServer(1); got == nil || got.IP != "127.0.0.2" || got.Port != 56901 || got.Name != "New" {
		t.Fatalf("registration did not refresh config: %+v", got)
	}
	if got := manager.GetServerList(); len(got) != 1 || got[0].Percent != 20 {
		t.Fatalf("registration did not refresh list: %+v", got)
	}
	manager.Unregister(&model.MsgUnregister{Code: 1})
	if manager.GetServer(1) != nil || len(manager.GetServerList()) != 0 {
		t.Fatal("unregistered server is still available")
	}
}
