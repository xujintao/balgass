package model

import (
	"bytes"
	"encoding/binary"

	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/utils"
)

// invalid api [body]f101cdfd98c8faabfccfabfccdfd98c8faabfccfabfccfabfccfabfccfabfccf000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007dfaa614302e312e350000004d374234564d3443356938424334396240000000
// cdfd98c8faabfccfabfc
// cdfd98c8faabfccfabfccfabfccfabfccfabfccf
// 00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
// 7dfaa614
// 302e312e35000000
// 4d374234564d34433569384243343962
// 40000000
type MsgLogin struct {
	Username  string
	Password  string
	HWID      string
	TickCount int
	Version   string
	Serial    string
}

func (msg *MsgLogin) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// username
	var username [10]byte
	_, err := br.Read(username[:])
	if err != nil {
		return err
	}
	msg.Username = string(bytes.TrimRight(utils.Xor(username[:]), "\x00"))

	// password
	var password [20]byte
	_, err = br.Read(password[:])
	if err != nil {
		return err
	}
	msg.Password = string(bytes.TrimRight(utils.Xor(password[:]), "\x00"))

	// hwid
	var hwid [100]byte
	_, err = br.Read(hwid[:])
	if err != nil {
		return err
	}
	msg.HWID = string(bytes.TrimRight(hwid[:], "\x00"))

	// time
	var tickCount uint32
	err = binary.Read(br, binary.LittleEndian, &tickCount)
	if err != nil {
		return err
	}
	msg.TickCount = int(tickCount)

	// version
	var version [8]byte
	_, err = br.Read(version[:])
	if err != nil {
		return err
	}
	msg.Version = string(bytes.TrimRight(version[:], "\x00"))

	// serial
	var serial [16]byte
	_, err = br.Read(serial[:])
	if err != nil {
		return err
	}
	msg.Serial = string(serial[:])
	return nil
}

// invalid api [body]f330ffffffffffffffffffffffffffffffffffffffff1dffffff16ff00000000
type MsgDefineKey struct {
}

type MsgTest struct{}

func (msg *MsgTest) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(0x01)
	return buf.Bytes(), nil
}

func (msg *MsgTest) Unmarshal([]byte) error {
	return nil
}

type MsgPickCharacter struct {
	Name string
}

type MsgSetCharacter struct {
	Name string
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

	// result
	buf.WriteByte(byte(msg.Result))

	// id
	binary.Write(&buf, binary.BigEndian, uint16(msg.ID))

	// version
	var version [8]byte
	copy(version[:], msg.Version)
	buf.Write(version[:])

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
	Dirs []int
	Path maps.Path
}

func (msg *MsgMove) Unmarshal([]byte) error {
	var x, y int
	var bufDir [8]int

	size := bufDir[0] & 0x0F
	dirs := make([]int, size)
	path := make(maps.Path, size)
	for i := 0; i < size; i++ {
		if i == 0 {
			dir := bufDir[i] >> 4 & 0x0F
			dirs[i] = dir
			dirPot := maps.Dirs[dir]
			path[i].X = x + dirPot.X
			path[i].Y = y + dirPot.Y
			continue
		}
		dir := bufDir[i] >> 4 & 0x0F
		dirs[i] = dir
		dirPot := maps.Dirs[dir]
		path[i].X = path[i-1].X + dirPot.X
		path[i].Y = path[i-1].Y + dirPot.Y

		dir = bufDir[i] & 0x0F
		dirs[i+1] = dir
		dirPot = maps.Dirs[dir]
		path[i+1].X = path[i].X + dirPot.X
		path[i+1].Y = path[i].Y + dirPot.Y
	}
	msg.Dirs = dirs
	msg.Path = path
	return nil
}

// struct PMSG_RECVMOVE
//
//	{
//		PBMSG_HEAD h;
//		BYTE NumberH;	// 3
//		BYTE NumberL;	// 4
//		BYTE X;	// 5
//		BYTE Y;	// 6
//		BYTE Path;	// 7
//	};
type MsgMoveReply struct {
	Number int
	X      int
	Y      int
	Dir    int
}

func (msg *MsgMoveReply) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	// buf.WriteByte(byte(msg.Result))
	return buf.Bytes(), nil
}

type MsgAttack struct {
	Target int
}

type MsgSkillAttack struct {
	Target int
	Skill  int
}
