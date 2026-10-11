package model

import (
	"encoding/binary"
	"testing"
)

func registration(code uint16, ip, name string, port uint16, visible byte) []byte {
	buf := []byte{0, 0, 0, 35, 1}
	binary.LittleEndian.PutUint16(buf[1:3], code)
	var address [16]byte
	copy(address[:], ip)
	buf = append(buf, address[:]...)
	buf = binary.LittleEndian.AppendUint16(buf, port)
	buf = append(buf, visible, byte(len(name)))
	return append(buf, name...)
}

func TestRegisterUnmarshal(t *testing.T) {
	var msg MsgRegister
	if err := msg.Unmarshal(registration(12, "127.0.0.1", "Market", 56912, 0)); err != nil {
		t.Fatal(err)
	}
	if msg.Code != 12 || msg.Config.Code != 12 || msg.Config.IP != "127.0.0.1" || msg.Config.Port != 56912 || msg.Config.Visible || msg.Config.Name != "Market" || msg.Percent != 35 || msg.Type_ != 1 {
		t.Fatalf("unexpected registration: %+v", msg)
	}
}

func TestRegisterUnmarshalRejectsInvalid(t *testing.T) {
	valid := registration(1, "127.0.0.1", "Regular", 56900, 1)
	for _, change := range []func([]byte) []byte{
		func(b []byte) []byte { return b[:24] },
		func(b []byte) []byte { return b[:len(b)-1] },
		func(b []byte) []byte { b[21], b[22] = 0, 0; return b },
		func(b []byte) []byte { b[23] = 2; return b },
		func(b []byte) []byte { b[5] = 0; return b },
	} {
		buf := append([]byte(nil), valid...)
		if err := new(MsgRegister).Unmarshal(change(buf)); err == nil {
			t.Fatal("accepted invalid registration")
		}
	}
}
