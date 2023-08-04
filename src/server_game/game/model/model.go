package model

import (
	"bytes"
	"encoding/binary"

	"github.com/xujintao/balgass/src/server_game/game/maps"
)

type MsgTest struct{}

func (msg *MsgTest) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(0x01)
	return buf.Bytes(), nil
}

func (msg *MsgTest) Unmarshal([]byte) error {
	return nil
}

type MsgChat struct {
	Name string
	Msg  string
}

func (*MsgChat) Unmarshal([]byte) error {
	// var buf bytes.Buffer
	// buf.WriteByte(byte(msg.Result))
	// var ids [2]uint8
	// binary.BigEndian.PutUint16(ids[:], uint16(msg.ID))
	// buf.Write(ids[:])
	// buf.WriteString(msg.Version)
	// return buf.Bytes(), nil
	var msg MsgChat
	msg.Name = "api"
	msg.Msg = "hello world"
	return nil
}

type MsgWhisper struct {
	Name string
	Msg  string
}

type MsgLive struct {
	Time         int
	AttackSpeed  int
	Agility      int
	MagicSpeed   int
	Version      string
	ServerSeason int
}

// struct PMSG_RESULT
//
//	{
//		PBMSG_HEAD h;
//		unsigned char subcode;	// 3
//		unsigned char result;	// 4
//	};
type MsgConnectFailed struct {
	Result int
}

func (msg *MsgConnectFailed) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(byte(msg.Result))
	return buf.Bytes(), nil
}

// #pragma pack (1)
// struct PMSG_JOINRESULT
//
//	{
//		PBMSG_HEAD h;	// C1:F1
//		BYTE scode;	// 3
//		BYTE result;	// 4
//		BYTE NumberH;	// 5
//		BYTE NumberL;	// 6
//		BYTE CliVersion[8];	// 7
//	};
//
// #pragma pack ()
type MsgConnectSuccess struct {
	Result  int
	ID      int
	Version string
}

func (msg *MsgConnectSuccess) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(byte(msg.Result))
	var ids [2]uint8
	binary.BigEndian.PutUint16(ids[:], uint16(msg.ID))
	buf.Write(ids[:])
	buf.WriteString(msg.Version)
	return buf.Bytes(), nil
}

type MsgUseItem struct {
	InventoryPos       int
	InventoryPosTarget int
	ItemUseType        int
}

type MsgLearnMasterSkill struct {
	SkillIndex int
}

type MsgSkillList struct {
}

// struct PMSG_MOVE
//
//	{
//		PBMSG_HEAD h;	// C1:1D
//		BYTE X;	// 3
//		BYTE Y;	// 4
//		BYTE Path[8];	// 5
//	};
type MsgMove struct {
	Path maps.Path
}

func (msg *MsgMove) Unmarshal([]byte) error {
	var x, y int
	var bufDir [8]int

	size := bufDir[0] & 0x0F
	path := make(maps.Path, size)
	for i := 0; i < size; i++ {
		if i == 0 {
			dir := bufDir[i] >> 4 & 0x0F
			dirPot := maps.Dirs[dir]
			path[i].X = x + dirPot.X
			path[i].Y = y + dirPot.Y
			continue
		}
		dir := bufDir[i] >> 4 & 0x0F
		dirPot := maps.Dirs[dir]
		path[i].X = path[i-1].X + dirPot.X
		path[i].Y = path[i-1].Y + dirPot.Y
		dir = bufDir[i] & 0x0F
		dirPot = maps.Dirs[dir]
		path[i+1].X = x + dirPot.X
		path[i+1].Y = y + dirPot.Y
	}
	msg.Path = path
	return nil
}

type MsgAttack struct {
	Target int
}

type MsgSkillAttack struct {
	Target int
	Skill  int
}
