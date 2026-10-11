package model

import (
	"encoding/binary"
	"testing"
)

func TestServerRegisterMarshal(t *testing.T) {
	msg := MsgServerRegister{Code: 12, Percent: 35, Type_: 1, IP: "127.0.0.1", Port: 56912, Visible: true, Name: "Market"}
	data, err := msg.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 25+len(msg.Name) || binary.LittleEndian.Uint16(data[1:3]) != 12 ||
		string(data[5:14]) != msg.IP || binary.LittleEndian.Uint16(data[21:23]) != 56912 ||
		data[23] != 1 || data[24] != byte(len(msg.Name)) || string(data[25:]) != msg.Name {
		t.Fatalf("unexpected registration bytes: %v", data)
	}
	msg.IP = "::1"
	if _, err := msg.Marshal(); err == nil {
		t.Fatal("accepted non-IPv4 address")
	}
}
