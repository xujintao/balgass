package model

import (
	"bytes"
	"encoding/binary"

	"github.com/xujintao/balgass/src/server_game/game/item"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/utils"
)

type CreateViewportMonster struct {
	Index                  int
	Class                  int
	X                      int
	Y                      int
	TX                     int
	TY                     int
	Dir                    int
	PentagramMainAttribute int
	Level                  int
	MaxHP                  int
	HP                     int
	BuffEffects            []int
}

// pack(1)
type MsgCreateViewportMonsterReply struct {
	Monsters []*CreateViewportMonster
}

func (msg *MsgCreateViewportMonsterReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(len(msg.Monsters)))
	for _, monster := range msg.Monsters {
		binary.Write(&bw, binary.BigEndian, uint16(monster.Index))
		binary.Write(&bw, binary.BigEndian, uint16(monster.Class))
		bw.WriteByte(byte(monster.X))
		bw.WriteByte(byte(monster.Y))
		bw.WriteByte(byte(monster.TX))
		bw.WriteByte(byte(monster.TY))
		bw.WriteByte(byte(monster.Dir << 4))
		bw.WriteByte(byte(monster.PentagramMainAttribute))
		binary.Write(&bw, binary.BigEndian, uint16(monster.Level))
		binary.Write(&bw, binary.BigEndian, uint32(monster.MaxHP))
		binary.Write(&bw, binary.BigEndian, uint32(monster.HP))
		bw.WriteByte(byte(len(monster.BuffEffects)))
		for _, buff := range monster.BuffEffects {
			bw.WriteByte(byte(buff))
		}
	}
	return bw.Bytes(), nil
}

type DestroyViewportObject struct {
	Index int
}

// pack(1)
type MsgDestroyViewportObjectReply struct {
	Objects []*DestroyViewportObject
}

func (msg *MsgDestroyViewportObjectReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(len(msg.Objects)))
	for _, obj := range msg.Objects {
		binary.Write(&bw, binary.BigEndian, uint16(obj.Index))
	}
	return bw.Bytes(), nil
}

// pack(1)
type MsgAttack struct {
	Target int
	Action int
	Dir    int
}

func (msg *MsgAttack) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// target
	var target uint16
	err := binary.Read(br, binary.BigEndian, &target)
	if err != nil {
		return err
	}
	msg.Target = int(target)

	// action
	action, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Action = int(action)

	// dir
	dir, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Dir = int(dir)

	return nil
}

// pack(1)
type MsgAttackReply struct {
	Index      int
	Damage     int
	DamageType int
	SDDamage   int
}

func (msg *MsgAttackReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	binary.Write(&bw, binary.BigEndian, uint16(msg.Index))
	binary.Write(&bw, binary.BigEndian, uint16(msg.Damage))
	binary.Write(&bw, binary.BigEndian, uint16(msg.DamageType))
	binary.Write(&bw, binary.BigEndian, uint16(msg.SDDamage))
	return bw.Bytes(), nil
}

// pack(1)
type MsgAttackDieReply struct {
	Target int
	Skill  int
	Killer int
}

func (msg *MsgAttackDieReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	binary.Write(&bw, binary.BigEndian, uint16(msg.Target))
	bw.WriteByte(byte(msg.Skill))
	binary.Write(&bw, binary.BigEndian, uint16(msg.Killer))
	return bw.Bytes(), nil
}

// pack(1)
type MsgAttackEffectReply struct {
	Target       int
	HP           int
	MaxHP        int
	Level        int
	IceEffect    int
	PoisonEffect int
}

func (msg *MsgAttackEffectReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	binary.Write(&bw, binary.LittleEndian, uint16(msg.Target))
	binary.Write(&bw, binary.LittleEndian, uint32(msg.HP))
	binary.Write(&bw, binary.LittleEndian, uint32(msg.MaxHP))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.Level))
	bw.WriteByte(byte(msg.IceEffect))
	bw.WriteByte(byte(msg.PoisonEffect))
	return bw.Bytes(), nil
}

// pack(1)
type MsgAttackHPReply struct {
	Target int
	MaxHP  int
	HP     int
}

func (msg *MsgAttackHPReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	binary.Write(&bw, binary.BigEndian, uint16(msg.Target))
	// binary.Write(&bw, binary.BigEndian, uint32(msg.MaxHP))
	bw.WriteByte(byte(msg.MaxHP >> 24))
	bw.WriteByte(byte(msg.MaxHP >> 8))
	bw.WriteByte(byte(msg.MaxHP >> 16))
	bw.WriteByte(byte(msg.MaxHP))
	// binary.Write(&bw, binary.BigEndian, uint32(msg.HP))
	bw.WriteByte(byte(msg.HP >> 24))
	bw.WriteByte(byte(msg.HP >> 8))
	bw.WriteByte(byte(msg.HP >> 16))
	bw.WriteByte(byte(msg.HP))
	return bw.Bytes(), nil
}

// pack(1)
type MsgAction struct {
	Dir    int
	Action int
}

func (msg *MsgAction) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// dir
	dir, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Dir = int(dir)

	// action
	action, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Action = int(action)

	return nil
}

// pack(1)
type MsgActionReply struct {
	Index  int
	Dir    int
	Action int
	Target int
}

func (msg *MsgActionReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	binary.Write(&bw, binary.BigEndian, uint16(msg.Index))
	bw.WriteByte(byte(msg.Dir))
	bw.WriteByte(byte(msg.Action))
	binary.Write(&bw, binary.BigEndian, uint16(msg.Target))
	return bw.Bytes(), nil
}

// pack(1)
type MsgHPReply struct {
	HP int
	SD int
}

func (msg *MsgHPReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(0xFE)
	binary.Write(&bw, binary.BigEndian, uint16(msg.HP))
	bw.WriteByte(0)
	binary.Write(&bw, binary.BigEndian, uint16(msg.SD))
	return bw.Bytes(), nil
}

// pack(1)
type MsgMPReply struct {
	MP int
	AG int
}

func (msg *MsgMPReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(0xFF)
	binary.Write(&bw, binary.BigEndian, uint16(msg.MP))
	binary.Write(&bw, binary.BigEndian, uint16(msg.AG))
	return bw.Bytes(), nil
}

type MsgMuunSystem struct {
}

func (msg *MsgMuunSystem) Unmarshal(buf []byte) error {
	return nil
}

// pack(1)
type MsgMove struct {
	Dirs []int
	Path maps.Path
}

func (msg *MsgMove) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// x
	x, err := br.ReadByte()
	if err != nil {
		return err
	}

	// y
	y, err := br.ReadByte()
	if err != nil {
		return err
	}

	// bufDir
	var bufDir [8]byte
	_, err = br.Read(bufDir[:])
	if err != nil {
		return err
	}

	size := bufDir[0] & 0x0F
	dirs := make([]int, size)
	path := make(maps.Path, size)
	for i := range path {
		path[i] = &maps.Pot{}
	}
	for i := 0; i < int(size); i++ {
		if i == 0 {
			dir := int(bufDir[(i+1)/2] >> 4 & 0x0F)
			dirs[i] = dir
			dirPot := maps.Dirs[dir]
			path[i].X = int(x) + dirPot.X
			path[i].Y = int(y) + dirPot.Y
			continue
		}
		dir := 0
		if i%2 == 1 {
			dir = int(bufDir[(i+1)/2] >> 4 & 0x0F)
		} else {
			dir = int(bufDir[(i+1)/2] & 0x0F)
		}
		dirs[i] = dir
		dirPot := maps.Dirs[dir]
		path[i].X = path[i-1].X + dirPot.X
		path[i].Y = path[i-1].Y + dirPot.Y
	}
	msg.Dirs = dirs
	msg.Path = path
	return nil
}

// pack(1)
type MsgMoveReply struct {
	Number int
	X      int
	Y      int
	Dir    int
}

func (msg *MsgMoveReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	binary.Write(&bw, binary.BigEndian, uint16(msg.Number))
	bw.WriteByte(byte(msg.X))
	bw.WriteByte(byte(msg.Y))
	bw.WriteByte(byte(msg.Dir))
	return bw.Bytes(), nil
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

type MsgLogout struct {
	Flag int
}

func (msg *MsgLogout) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)
	// flag
	flag, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Flag = int(flag)
	return nil
}

type MsgLogoutReply struct {
	Flag int
}

func (msg *MsgLogoutReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Flag))
	return bw.Bytes(), nil
}

type MsgHack struct {
	Flag1 int
	Flag2 int
}

func (msg *MsgHack) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// flag1
	flag, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Flag1 = int(flag)

	// flag2
	flag2, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Flag2 = int(flag2)

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
	EnableCharacterClass int
	MoveCnt              int
	Count                int
	WarehouseExpansion   int
	CharacterList        []*MsgCharacter
}

func (msg *MsgGetCharacterListReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer

	// EnableCharacter
	bw.WriteByte(byte(msg.EnableCharacterClass))

	// MoveCnt
	bw.WriteByte(byte(msg.MoveCnt))

	// Count
	bw.WriteByte(byte(len(msg.CharacterList)))

	// WarehouseExpansion
	bw.WriteByte(byte(msg.WarehouseExpansion))

	// CharacterList
	for _, c := range msg.CharacterList {
		// index
		bw.WriteByte(byte(c.Index))

		// name
		var name [10]byte
		copy(name[:], c.Name)
		bw.Write(name[:])
		bw.WriteByte(0) // padding 1 byte

		// level
		binary.Write(&bw, binary.LittleEndian, uint16(c.Level))

		// ctlcode
		bw.WriteByte(byte(c.CtlCode))

		var chars [18]byte

		// class
		class := byte(c.Class << 5)
		switch c.ChangeUp {
		case 1:
			class |= 0x10
		case 2:
			class |= 0x18
		}
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

		// slot0~slot6 index -> chars[1]~chars[5]
		chars[1] = byte(inventory[0].Index)
		chars[2] = byte(inventory[1].Index)
		chars[3] = byte(inventory[2].Index&0x0F<<4 | inventory[3].Index&0x0F)
		chars[4] = byte(inventory[4].Index&0x0F<<4 | inventory[5].Index&0x0F)
		chars[5] = byte(inventory[6].Index & 0x0F << 4)

		// slot0~slot6 index extention1 -> chars[9] bit3~bit7
		extend := inventory[2].Index&0x10<<3 |
			inventory[3].Index&0x10<<2 |
			inventory[4].Index&0x10<<1 |
			inventory[5].Index&0x10<<0 |
			inventory[6].Index&0x10>>1
		chars[9] = byte(extend)

		// slot0~slot6 index extention2 -> chars[12]~chars[15]
		chars[13] |= (byte(inventory[2].Index & 0x1E0 >> 5))
		chars[14] |= (byte(inventory[3].Index&0x1E0>>1 | inventory[4].Index&0x1E0>>5))
		chars[15] |= (byte(inventory[5].Index&0x1E0>>1 | inventory[6].Index&0x1E0>>5))

		// slot0~slot6 level -> chars[6]~chars[8]
		var level uint32
		var data [4]byte
		for i, v := range inventory[0:7] {
			if v.Index == 512 {
				continue
			}
			level |= uint32(v.Level) << i * 3
		}
		binary.BigEndian.PutUint32(data[:], level)
		copy(chars[6:9], data[1:])

		// slot7 -> chars[5] bit2~bit3 4=1D, 8=2D, 12=3D 0=empty
		// slot7 -> chars[9] bit0~bit2
		// slot7 -> chars[16] bit2~bit4
		switch inventory[7].Index {
		case 0, 1, 2:
			chars[5] |= 4                            // 1D
			chars[9] |= byte(inventory[7].Index + 1) // 1=精灵之翼 2=天使之翼 3=恶魔之翼
		case 41:
			chars[5] |= 4 // 1D
			chars[9] |= 4 // 4=灾难之翼
		case 266, 267:
			chars[5] |= 4                              // 1D
			chars[9] |= byte(inventory[7].Index - 261) // 5=征服者的翅膀 6=善恶的翅膀
		case 3, 4, 5, 6:
			chars[5] |= 8                            // 2D
			chars[9] |= byte(inventory[7].Index - 2) // 1=圣灵之翼 2=魔魂之翼 3=飞龙之翼 4=暗黑之翼
		case 42:
			chars[5] |= 8 // 2D
			chars[9] |= 6 // 6=绝望之翼
		case 49:
			chars[5] |= 8 // 2D
			chars[9] |= 7 // 7=武者披风
		case 36, 37, 38, 39, 40:
			chars[5] |= 12                            // 3D
			chars[9] |= byte(inventory[7].Index - 35) // 1=暴风之翼 2=时空之翼 3=幻影之翼 4=破灭之翼 5=帝王披风
		case 43:
			chars[5] |= 12 // 3D
			chars[9] |= 6  // 6=次元之翼
		case 50:
			chars[5] |= 12 // 3D
			chars[9] |= 7  // 7=斗皇披风
		case 262:
			chars[5] |= 8  // 2.5D
			chars[16] |= 4 // 死亡披风
		case 263:
			chars[5] |= 8  // 2.5D
			chars[16] |= 8 // 混沌之翼
		case 264:
			chars[5] |= 8   // 2.5D
			chars[16] |= 12 // 魔力之翼
		case 265:
			chars[5] |= 8   // 2.5D
			chars[16] |= 16 // 生命之翼
		}

		// slot8 -> chars[5] bit0~bit1 0=守护天使 1=小恶魔 3=empty
		// slot8 -> chars[16] bit5~bit7 bit0~bit1
		switch inventory[8].Index {
		case 0, 1, 2:
			chars[5] |= byte(inventory[8].Index)
		case 64:
			chars[16] |= 0x20 // 强化恶魔
		case 65:
			chars[16] |= 0x40 // 强化天使
		case 80:
			chars[16] |= 0xE0 // 熊猫
		case 123:
			chars[16] |= 0x60 // 幼龙骨架
		default:
			chars[5] |= 3
		}

		// done
		bw.Write(chars[:])
		bw.WriteByte(byte(c.GuildStatus))
		bw.WriteByte(byte(c.PKLevel))
		bw.WriteByte(0) // padding 1 byte
	}
	return bw.Bytes(), nil
}

type MsgEnableCharacterClassReply struct {
	Class int
}

func (msg *MsgEnableCharacterClassReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Class))
	return bw.Bytes(), nil
}

type MsgResetCharacterReply struct {
	Reset string
}

func (msg *MsgResetCharacterReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer

	// reset
	bw.WriteString(msg.Reset)

	return bw.Bytes(), nil
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
	// 0x00 - Dark Wizard
	// 0x10 - Dark Knight
	// 0x20 - Elf
	// 0x30 - Magic Gladiator
	// 0x40 - Dark Lord
	// 0x50 - Summoner
	// 0x60 - Rage Fighter
	// 0x70 - GrowLancer
	class, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Class = int(class >> 4)

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
	bw.WriteByte(byte(msg.Class << 5))

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
	var password [7]byte
	_, err = br.Read(password[:])
	if err != nil {
		return err
	}
	msg.Password = string(bytes.TrimRight(password[:], "\x00"))

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

type MsgLoadCharacter struct {
	Name     string
	Position int
}

func (msg *MsgLoadCharacter) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// name
	var name [10]byte
	_, err := br.Read(name[:])
	if err != nil {
		return err
	}
	msg.Name = string(bytes.TrimRight(name[:], "\x00"))

	// position
	position, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Position = int(position)

	return nil
}

type MsgLoadCharacterReply struct {
	X                  int
	Y                  int
	MapNumber          int
	Dir                int
	Experience         int
	NextExperience     int
	LevelUpPoint       int
	Strength           int
	Dexterity          int
	Vitality           int
	Energy             int
	HP                 int
	MaxHP              int
	MP                 int
	MaxMP              int
	SD                 int
	MaxSD              int
	AG                 int
	MaxAG              int
	Money              int
	PKLevel            int
	CtlCode            int
	AddPoint           int
	MaxAddPoint        int
	Leadership         int
	MinusPoint         int
	MaxMinusPoint      int
	InventoryExpansion int
}

func (msg *MsgLoadCharacterReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.X))
	bw.WriteByte(byte(msg.Y))
	bw.WriteByte(byte(msg.MapNumber))
	bw.WriteByte(byte(msg.Dir))
	binary.Write(&bw, binary.BigEndian, uint64(msg.Experience))
	binary.Write(&bw, binary.BigEndian, uint64(msg.NextExperience))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.LevelUpPoint))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.Strength))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.Dexterity))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.Vitality))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.Energy))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.HP))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MaxHP))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MP))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MaxMP))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.SD))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MaxSD))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.AG))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MaxAG))
	bw.Write([]byte{0, 0}) // padding
	binary.Write(&bw, binary.LittleEndian, uint32(msg.Money))
	bw.WriteByte(byte(msg.PKLevel))
	bw.WriteByte(byte(msg.CtlCode))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.AddPoint))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MaxAddPoint))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.Leadership))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MinusPoint))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MaxMinusPoint))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.InventoryExpansion))
	bw.Write([]byte{0, 0}) // padding
	return bw.Bytes(), nil
}

// pack(1)
type MsgReloadCharacterReply struct {
	X          int
	Y          int
	MapNumber  int
	Dir        int
	HP         int
	MP         int
	SD         int
	AG         int
	Experience int
	Money      int
}

func (msg *MsgReloadCharacterReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.X))
	bw.WriteByte(byte(msg.Y))
	bw.WriteByte(byte(msg.MapNumber))
	bw.WriteByte(byte(msg.Dir))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.HP))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MP))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.SD))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.AG))
	binary.Write(&bw, binary.BigEndian, uint64(msg.Experience))
	binary.Write(&bw, binary.LittleEndian, uint32(msg.Money))
	return bw.Bytes(), nil
}

type MsgCheckCharacter struct {
	Name string
}

func (msg *MsgCheckCharacter) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// name
	var name [10]byte
	_, err := br.Read(name[:])
	if err != nil {
		return err
	}
	msg.Name = string(bytes.TrimRight(name[:], "\x00"))

	return nil
}

type MsgCheckCharacterReply struct {
	Name string

	// 0=success 1=failed
	Result int
}

func (msg *MsgCheckCharacterReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer

	// name
	var name [10]byte
	copy(name[:], msg.Name)
	bw.Write(name[:])

	// result
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

type MsgSkillAttack struct {
	Target int
	Skill  int
}
