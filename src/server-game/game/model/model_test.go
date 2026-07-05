package model

import (
	"bytes"
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

func TestUseSkillDurationProtocol(t *testing.T) {
	var duration MsgUseSkillDuration
	if err := duration.Unmarshal([]byte{
		11, 0x01, 22, 0x02, 3, 0x04, 5, 0x06, 7, 8,
	}); err != nil {
		t.Fatalf("duration Unmarshal() error = %v", err)
	}
	if duration.X != 11 || duration.Y != 22 || duration.Dir != 3 ||
		duration.Skill != 0x0102 || duration.Target != 0x0406 ||
		duration.Dis != 5 || duration.TargetPos != 7 || duration.MagicKey != 8 {
		t.Fatalf("duration = %#v", duration)
	}

	reply, err := (&MsgUseSkillDurationReply{
		X: 11, Y: 22, Dir: 3, Skill: 0x0102, Index: 0x0406,
	}).Marshal()
	if err != nil {
		t.Fatalf("duration reply Marshal() error = %v", err)
	}
	if want := []byte{11, 22, 3, 0x01, 0x04, 0x02, 0x06}; !bytes.Equal(reply, want) {
		t.Fatalf("duration reply = %v, want %v", reply, want)
	}

	buf := []byte{0x01, 6, 0x02, 11, 12, 13}
	for i := 0; i < 6; i++ {
		buf = append(buf, byte(i), byte(i+1), byte(i+10))
	}
	var multi MsgUseSkillAttackMultiTarget
	if err := multi.Unmarshal(buf); err != nil {
		t.Fatalf("multi target Unmarshal() error = %v", err)
	}
	if multi.Skill != 0x0102 || multi.X != 11 || multi.Serial != 12 || multi.Y != 13 {
		t.Fatalf("multi target = %#v", multi)
	}
	if len(multi.Targets) != 5 {
		t.Fatalf("target count = %d, want 5", len(multi.Targets))
	}
	if got := multi.Targets[0]; got.Target != 10 || got.MagicKey != 1 {
		t.Fatalf("first target = %#v", got)
	}
	if err := multi.Unmarshal(buf[:len(buf)-1]); err == nil {
		t.Fatal("truncated multi target packet was accepted")
	}
}
