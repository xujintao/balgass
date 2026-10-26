package model

import (
	"bytes"
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/item"
)

func TestMixItemsPackets(t *testing.T) {
	var request MsgMixItems
	if err := request.Unmarshal([]byte{0x2c, 3, 1}); err != nil {
		t.Fatal(err)
	}
	if request.Type != 0x2c || request.Position != 3 || request.MixCount != 1 {
		t.Fatalf("request = %#v", request)
	}
	if err := request.Unmarshal([]byte{0x2c}); err == nil {
		t.Fatal("short request was accepted")
	}
	data, err := (&MsgMixItemsReply{Result: 7}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 14 || data[0] != 7 {
		t.Fatalf("reply = %v", data)
	}
	it := item.NewItem(12, 60)
	frame, err := it.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	data, err = (&MsgMixItemsReply{Result: 1, Item: it, Position: 4}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 14 || data[0] != 1 || data[13] != 4 || !bytes.Equal(data[1:13], frame) {
		t.Fatalf("successful reply = %v", data)
	}
}
