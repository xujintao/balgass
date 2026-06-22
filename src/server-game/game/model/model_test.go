package model

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestMasterSkillRepliesMarshalFloatBits(t *testing.T) {
	reply := MsgLearnMasterSkillReply{
		MsgMasterSkill: MsgMasterSkill{
			MasterSkillCurValue:  1.5,
			MasterSkillNextValue: 2.25,
		},
	}
	data, err := reply.Marshal()
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if got := binary.LittleEndian.Uint32(data[16:20]); got != math.Float32bits(1.5) {
		t.Fatalf("cur value bits = %#x, want %#x", got, math.Float32bits(1.5))
	}
	if got := binary.LittleEndian.Uint32(data[20:24]); got != math.Float32bits(2.25) {
		t.Fatalf("next value bits = %#x, want %#x", got, math.Float32bits(2.25))
	}

	list := MsgMasterSkillListReply{
		Skills: []*MsgMasterSkill{{
			MasterSkillCurValue:  1.5,
			MasterSkillNextValue: 2.25,
		}},
	}
	data, err = list.Marshal()
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if got := binary.LittleEndian.Uint32(data[11:15]); got != math.Float32bits(1.5) {
		t.Fatalf("list cur value bits = %#x, want %#x", got, math.Float32bits(1.5))
	}
	if got := binary.LittleEndian.Uint32(data[15:19]); got != math.Float32bits(2.25) {
		t.Fatalf("list next value bits = %#x, want %#x", got, math.Float32bits(2.25))
	}
}
