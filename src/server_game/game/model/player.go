package model

import (
	"bytes"
	"encoding/binary"

	"github.com/xujintao/balgass/src/server_game/game/item"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/utils"
)

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
type MsgConnectReply struct {
	// 1: success
	// others: failed
	Result  int
	ID      int
	Version string
}

func (msg *MsgConnectReply) Marshal() ([]byte, error) {
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

// invalid api [body]f101cdfd98c8faabfccfabfccdfd98c8faabfccfabfccfabfccfabfccfabfccf000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007dfaa614302e312e350000004d374234564d3443356938424334396240000000
// cdfd98c8faabfccfabfc
// cdfd98c8faabfccfabfccfabfccfabfccfabfccf
// 00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
// 7dfaa614
// 302e312e35000000
// 4d374234564d34433569384243343962
// 40000000
type MsgLogin struct {
	Account   string
	Password  string
	HWID      string
	TickCount int
	Version   string
	Serial    string
}

func (msg *MsgLogin) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// account
	var account [10]byte
	_, err := br.Read(account[:])
	if err != nil {
		return err
	}
	msg.Account = string(bytes.TrimRight(utils.Xor(account[:]), "\x00"))

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

type MsgLoginReply struct {
	// 0: password doesn't match 密码错误
	// 1: success
	// 2: account doesn't exist 账号错误
	// 3: already online 该账号正在使用中
	// 4: machine id limit 本服务器可容纳的人数已满
	// 5: machine id banned 客服提示：该账号目前被禁止使用
	// 6: version unmatched 您的游戏版本不对，请到官方网站下载最新的版本
	// 8: login counts exceeded 3 失败3次连接中断
	// 9: 没有付款信息
	// 10: 本账号的使用期限已到期
	// 11: 本账号的储值点数不足
	// 12: 这个IP的使用期限已到期
	// 13: 这个IP的储值点数不足
	// 17: 是15岁以上能够使用的服务器
	// 64: non-vip 未购买收费服务器入场券而无法进入
	Result int
}

func (msg *MsgLoginReply) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(byte(msg.Result))
	return buf.Bytes(), nil
}

type MsgSetAccount struct {
	*MsgLogin
	*Account
	Err error
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

type MsgGetCharacterList struct{}

func (msg *MsgGetCharacterList) Unmarshal(buf []byte) error {
	return nil
}

type MsgCharacter struct {
	Index    int
	Name     string
	Level    int
	CtlCode  int
	Class    int
	ChangeUp int
	// CharSet     [18]byte
	Inventory   []*item.Item
	GuildStatus int
	PKLevel     int
}

type MsgGetCharacterListReply struct {
	EnableCharacter    int
	MoveCnt            int
	Count              int
	WarehouseExpansion int
	CharacterList      []*MsgCharacter
}

func (msg *MsgGetCharacterListReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer

	// EnableCharacter
	bw.WriteByte(byte(msg.EnableCharacter))

	// MoveCnt
	bw.WriteByte(byte(msg.MoveCnt))

	// Count
	bw.WriteByte(byte(len(msg.CharacterList)))

	// WarehouseExpansion
	bw.WriteByte(byte(msg.WarehouseExpansion))

	// CharacterList
	for i, c := range msg.CharacterList {
		// index
		bw.WriteByte(byte(i))

		// name
		var name [10]byte
		copy(name[:], c.Name)
		bw.Write(name[:])

		// level
		binary.Write(&bw, binary.LittleEndian, uint16(c.Level))

		// ctlcode
		bw.WriteByte(byte(c.CtlCode))

		var chars [18]byte

		// class
		class := byte(c.Class) << 5
		switch c.ChangeUp {
		case 1:
			class |= 0x10
		case 2:
			class |= 0x18
		}
		// bw.WriteByte(class)
		chars[0] = class

		// inventory
		inventory := make([]*item.Item, len(c.Inventory))
		for i, v := range c.Inventory {
			if v == nil {
				inventory[i] = item.NewItem(0, 512)
			} else {
				inventory[i] = c.Inventory[i]
			}
		}
		chars[1] = byte(inventory[0].Index)
		chars[2] = byte(inventory[1].Index)
		chars[3] = byte(inventory[2].Index&0x0F<<4 | inventory[3].Index&0x0F)
		chars[4] = byte(inventory[4].Index&0x0F<<4 | inventory[5].Index&0x0F)
		// slot8: 0=守护天使 1=小恶魔 3=empty
		// slot7: 4=1D, 8=2D, 12=3D 0=empty
		chars[5] = (byte(inventory[6].Index&0x0F<<4 | 0x03))
		var level uint32
		var levels [4]byte
		for i, v := range c.Inventory {
			if v == nil {
				continue
			}
			level |= uint32(v.Level) << i * 3
		}

		binary.BigEndian.PutUint32(levels[:], level)
		copy(chars[6:9], levels[1:])
		extend := inventory[2].Index&0x10<<3 |
			inventory[3].Index&0x10<<2 |
			inventory[4].Index&0x10<<1 |
			inventory[5].Index&0x10<<0 |
			inventory[6].Index&0x10>>1
		// 1=精灵之翼 1D
		// 2=天使之翼 1D
		// 3=恶魔之翼 1D
		// 4=灾难之翼 1D

		// 1=圣灵之翼 2D
		// 2=魔魂之翼 2D
		// 3=飞龙之翼 2D
		// 4=暗黑之翼 2D
		// 6=绝望之翼 2D
		// 7=武者披风 2D

		// 1=暴风之翼 3D
		// 2=时空之翼 3D
		// 3=幻影之翼 3D
		// 4=破灭之翼 3D
		// 5=帝王披风 3D
		// 6=次元之翼 3D
		// 7=斗皇披风 3D
		wingKind := 0
		chars[9] = (byte(extend | wingKind))
		chars[13] |= (byte(inventory[2].Index & 0x1E0 >> 5))
		chars[14] |= (byte(inventory[3].Index&0x1E0>>1 | inventory[4].Index&0x1E0>>5))
		chars[15] |= (byte(inventory[5].Index&0x1E0>>1 | inventory[6].Index&0x1E0>>5))
		bw.Write(chars[:])

		bw.WriteByte(byte(c.GuildStatus))
		bw.WriteByte(byte(c.PKLevel))

		bw.WriteByte(0) // padding 1 byte
	}
	return bw.Bytes(), nil
}

type MsgPickCharacter struct {
	Name string
}

type MsgSetCharacter struct {
	Name string
}

type MsgCreateCharacter struct {
	Name  string
	Class int
}

func (msg *MsgCreateCharacter) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// name
	var name [10]byte
	_, err := br.Read(name[:])
	if err != nil {
		return err
	}
	msg.Name = string(bytes.TrimRight(name[:], "\x00"))

	// class
	class, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Class = int(class)

	return nil
}

type MsgCreateCharacterReply struct {
	Result    int
	Name      string
	Index     int
	Level     int
	Class     int
	Equipment [24]byte
}

func (msg *MsgCreateCharacterReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer

	// result
	bw.WriteByte(byte(msg.Result))

	// name
	var name [10]byte
	copy(name[:], msg.Name)
	bw.Write(name[:])

	// index
	bw.WriteByte(byte(msg.Index))

	// level
	binary.Write(&bw, binary.LittleEndian, uint16(msg.Level))

	// class
	bw.WriteByte(byte(msg.Class))

	// equipment
	bw.Write(msg.Equipment[:])

	bw.WriteByte(0) // padding 1 byte

	return bw.Bytes(), nil
}

type MsgDeleteCharacter struct {
	Name     string
	Password string
}

func (msg *MsgDeleteCharacter) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// name
	var name [10]byte
	_, err := br.Read(name[:])
	if err != nil {
		return err
	}
	msg.Name = string(bytes.TrimRight(name[:], "\x00"))

	// password
	var password [20]byte
	_, err = br.Read(password[:])
	if err != nil {
		return err
	}
	msg.Password = string(bytes.TrimRight(utils.Xor(password[:]), "\x00"))

	return nil
}

type MsgDeleteCharacterReply struct {
	// 0: failed
	// 1: success
	// 2: password doesn't match 密码错误
	Result int
}

func (msg *MsgDeleteCharacterReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Result))
	return bw.Bytes(), nil
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
