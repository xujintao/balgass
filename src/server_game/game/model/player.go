package model

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/xujintao/balgass/src/server_game/game/item"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/skill"
	"github.com/xujintao/balgass/src/utils"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type MsgChat struct {
	Name string
	Msg  string
}

func (msg *MsgChat) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// Name
	var Name [10]byte
	_, err := br.Read(Name[:])
	if err != nil {
		return err
	}
	utf8, err := simplifiedchinese.GBK.NewDecoder().Bytes(Name[:])
	if err != nil {
		return err
	}
	msg.Name = string(bytes.TrimRight(utf8[:], "\x00"))

	// Msg
	var Msg [90]byte
	_, err = br.Read(Msg[:])
	if err != nil {
		return err
	}
	utf8, err = simplifiedchinese.GBK.NewDecoder().Bytes(Msg[:])
	if err != nil {
		return err
	}
	msg.Msg = string(bytes.TrimRight(utf8[:], "\x00"))

	return nil
}

type MsgChatReply struct {
	MsgChat
}

func (msg *MsgChatReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	gbk, err := simplifiedchinese.GBK.NewEncoder().String(msg.Name)
	if err != nil {
		return nil, err
	}
	var Name [10]byte
	copy(Name[:], gbk)
	bw.Write(Name[:])
	gbk, err = simplifiedchinese.GBK.NewEncoder().String(msg.Msg)
	if err != nil {
		return nil, err
	}
	var Msg [90]byte
	copy(Msg[:], gbk)
	bw.Write(Msg[:])
	return bw.Bytes(), nil
}

type MsgWhisper struct {
	MsgChat
}

type MsgWhisperReply struct {
	MsgChatReply
}

type MsgWhisperReplyFailed struct {
	Flag int
}

func (msg *MsgWhisperReplyFailed) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Flag))
	return bw.Bytes(), nil
}

type CreateViewportPlayer struct {
	Index                  int
	X                      int
	Y                      int
	Class                  int
	ChangeUp               int
	Inventory              [9]*item.Item
	Name                   string
	TX                     int
	TY                     int
	Dir                    int
	PKLevel                int
	PentagramMainAttribute int
	MuunItem               int
	MuunSubItem            int
	MuunRideItem           int
	Level                  int
	MaxHP                  int
	HP                     int
	ServerCode             int
	BuffEffects            []int
}

// pack(1)
type MsgCreateViewportPlayerReply struct {
	Players []*CreateViewportPlayer
}

func (msg *MsgCreateViewportPlayerReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(len(msg.Players)))
	for _, player := range msg.Players {
		binary.Write(&bw, binary.BigEndian, uint16(player.Index))
		bw.WriteByte(byte(player.X))
		bw.WriteByte(byte(player.Y))
		chars := MakeCharacterFrame(player.Class, player.ChangeUp, player.Inventory)
		bw.Write(chars[:])
		gbk, err := simplifiedchinese.GBK.NewEncoder().String(player.Name)
		if err != nil {
			return nil, err
		}
		var name [10]byte
		copy(name[:], gbk)
		bw.Write(name[:])
		bw.WriteByte(byte(player.TX))
		bw.WriteByte(byte(player.TY))
		bw.WriteByte(byte(player.Dir<<4 | player.PKLevel))
		bw.WriteByte(byte(player.PentagramMainAttribute))
		binary.Write(&bw, binary.BigEndian, uint16(player.MuunItem))
		binary.Write(&bw, binary.BigEndian, uint16(player.MuunSubItem))
		binary.Write(&bw, binary.BigEndian, uint16(player.MuunRideItem))
		binary.Write(&bw, binary.BigEndian, uint16(player.Level))
		// binary.Write(&bw, binary.BigEndian, uint32(player.MaxHP))
		bw.WriteByte(byte(player.MaxHP >> 24))
		bw.WriteByte(byte(player.MaxHP >> 8))
		bw.WriteByte(byte(player.MaxHP >> 16))
		bw.WriteByte(byte(player.MaxHP))
		// binary.Write(&bw, binary.BigEndian, uint32(player.HP))
		bw.WriteByte(byte(player.HP >> 24))
		bw.WriteByte(byte(player.HP >> 8))
		bw.WriteByte(byte(player.HP >> 16))
		bw.WriteByte(byte(player.HP))
		binary.Write(&bw, binary.LittleEndian, uint16(player.ServerCode))
		bw.WriteByte(byte(len(player.BuffEffects)))
		for _, buff := range player.BuffEffects {
			bw.WriteByte(byte(buff))
		}
	}
	return bw.Bytes(), nil
}

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
		// binary.Write(&bw, binary.BigEndian, uint32(monster.MaxHP))
		bw.WriteByte(byte(monster.MaxHP >> 24))
		bw.WriteByte(byte(monster.MaxHP >> 8))
		bw.WriteByte(byte(monster.MaxHP >> 16))
		bw.WriteByte(byte(monster.MaxHP))
		// binary.Write(&bw, binary.BigEndian, uint32(monster.HP))
		bw.WriteByte(byte(monster.HP >> 24))
		bw.WriteByte(byte(monster.HP >> 8))
		bw.WriteByte(byte(monster.HP >> 16))
		bw.WriteByte(byte(monster.HP))
		bw.WriteByte(byte(len(monster.BuffEffects)))
		for _, buff := range monster.BuffEffects {
			bw.WriteByte(byte(buff))
		}
	}
	return bw.Bytes(), nil
}

type DestroyViewport struct {
	Index int
}

// pack(1)
type MsgDestroyViewportObjectReply struct {
	Objects []*DestroyViewport
}

func (msg *MsgDestroyViewportObjectReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(len(msg.Objects)))
	for _, obj := range msg.Objects {
		binary.Write(&bw, binary.BigEndian, uint16(obj.Index))
	}
	return bw.Bytes(), nil
}

type CreateViewportItem struct {
	Index int
	X     int
	Y     int
	Item  *item.Item
}

// pack(1)
type MsgCreateViewportItemReply struct {
	Items []*CreateViewportItem
}

func (msg *MsgCreateViewportItemReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(len(msg.Items)))
	for _, item := range msg.Items {
		binary.Write(&bw, binary.BigEndian, uint16(item.Index))
		bw.WriteByte(byte(item.X))
		bw.WriteByte(byte(item.Y))
		data, err := item.Item.Marshal()
		if err != nil {
			return nil, err
		}
		bw.Write(data)
	}
	return bw.Bytes(), nil
}

// pack(1)
type MsgDestroyViewportItemReply struct {
	Items []*DestroyViewport
}

func (msg *MsgDestroyViewportItemReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(len(msg.Items)))
	for _, item := range msg.Items {
		binary.Write(&bw, binary.BigEndian, uint16(item.Index))
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
type MsgAttackDamageReply struct {
	Target     int
	Damage     int
	DamageType int
	SDDamage   int
}

func (msg *MsgAttackDamageReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	binary.Write(&bw, binary.BigEndian, uint16(msg.Target))
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

type MsgTeleport struct {
	GateNumber int
	X          int
	Y          int
}

func (msg *MsgTeleport) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// padding 1 byte
	_, err := br.ReadByte()
	if err != nil {
		return err
	}

	// GateNumber
	var GateNumber uint16
	err = binary.Read(br, binary.LittleEndian, &GateNumber)
	if err != nil {
		return err
	}
	msg.GateNumber = int(GateNumber)

	return nil
}

type MsgTeleportReply struct {
	GateNumber int
	MapNumber  int
	X          int
	Y          int
	Dir        int
}

func (msg *MsgTeleportReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(0) // padding 1 byte
	flag := 0
	if msg.GateNumber > 0 {
		flag = 256
	}
	binary.Write(&bw, binary.LittleEndian, uint16(flag))
	bw.WriteByte(byte(msg.MapNumber))
	bw.WriteByte(byte(msg.X))
	bw.WriteByte(byte(msg.Y))
	bw.WriteByte(byte(msg.Dir))
	bw.WriteByte(0) // padding 1 byte
	bw.WriteByte(0) // padding 1 byte
	return bw.Bytes(), nil
}

// pack(1)
type MsgGetItem struct {
	Index int
}

func (msg *MsgGetItem) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	var index uint16
	err := binary.Read(br, binary.BigEndian, &index)
	if err != nil {
		return err
	}
	msg.Index = int(index)

	return nil
}

// pack(1)
type MsgGetItemReply struct {
	Result int // -1=failed -2=money 0~237=postion
	Item   *item.Item
}

func (msg *MsgGetItemReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Result))
	itemFrame, err := msg.Item.Marshal()
	if err != nil {
		return nil, err
	}
	bw.Write(itemFrame)
	return bw.Bytes(), nil
}

// pack(1)
type MsgMoneyReply struct {
	Result int // -2
	Money  int
}

func (msg *MsgMoneyReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Result))
	var data [12]byte
	binary.BigEndian.PutUint32(data[:], uint32(msg.Money))
	binary.Write(&bw, binary.BigEndian, data[:])
	return bw.Bytes(), nil
}

// pack(1)
type MsgDropInventoryItem struct {
	X        int
	Y        int
	Position int
}

func (msg *MsgDropInventoryItem) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// x
	x, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.X = int(x)

	// y
	y, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Y = int(y)

	// position
	position, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Position = int(position)

	return nil
}

// pack(1)
type MsgDropInventoryItemReply struct {
	Result   int // 0=failed 1=success
	Position int
}

func (msg *MsgDropInventoryItemReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Result))
	bw.WriteByte(byte(msg.Position))
	return bw.Bytes(), nil
}

// pack(1)
type MsgMoveItem struct {
	SrcFlag     int
	SrcPosition int
	Item        *item.Item // ItemFrame   [12]byte
	DstFlag     int
	DstPosition int
}

func (msg *MsgMoveItem) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// srcFlag
	srcFlag, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.SrcFlag = int(srcFlag)

	// srcPosition
	srcPosition, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.SrcPosition = int(srcPosition)

	// itemFrame
	var data [12]byte
	n, err := br.Read(data[:])
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("item frame invalid [frame]%s",
			hex.EncodeToString(data[:n]))
	}
	// msg.ItemFrame = data

	// dstFlag
	dstFlag, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.DstFlag = int(dstFlag)

	// dstPosition
	dstPosition, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.DstPosition = int(dstPosition)

	return nil
}

// pack(1)
type MsgMoveItemReply struct {
	Result   int // -1=failed dstFlag=success
	Position int
	Item     *item.Item // ItemFrame [12]byte
}

func (msg *MsgMoveItemReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Result))
	bw.WriteByte(byte(msg.Position))
	// bw.Write(msg.ItemFrame[:])
	itemFrame, err := msg.Item.Marshal()
	if err != nil {
		return nil, err
	}
	bw.Write(itemFrame)
	return bw.Bytes(), nil
}

// pack(1)
type MsgUseItem struct {
	SrcPosition int
	DstPosition int
	UseType     int
}

func (msg *MsgUseItem) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// srcPosition
	srcPosition, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.SrcPosition = int(srcPosition)

	// dstPosition
	dstPosition, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.DstPosition = int(dstPosition)

	// useType
	useType, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.UseType = int(useType)

	return nil
}

// pack(1)
type MsgHPReply struct {
	Position int // -1=HP -2=maxHP
	HP       int
	Flag     int
	SD       int
}

func (msg *MsgHPReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Position))
	binary.Write(&bw, binary.BigEndian, uint16(msg.HP))
	bw.WriteByte(byte(msg.Flag))
	binary.Write(&bw, binary.BigEndian, uint16(msg.SD))
	return bw.Bytes(), nil
}

// pack(1)
type MsgMPReply struct {
	Position int // -1=MP -2=maxMP
	MP       int
	AG       int
}

func (msg *MsgMPReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Position))
	binary.Write(&bw, binary.BigEndian, uint16(msg.MP))
	binary.Write(&bw, binary.BigEndian, uint16(msg.AG))
	return bw.Bytes(), nil
}

// pack(1)
type MsgDeleteInventoryItemReply struct {
	Position int
	Flag     int
}

func (msg *MsgDeleteInventoryItemReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Position))
	bw.WriteByte(byte(msg.Flag))
	return bw.Bytes(), nil
}

// pack(1)
type MsgItemDurabilityReply struct {
	Position   int
	Durability int
	Flag       int
}

func (msg *MsgItemDurabilityReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Position))
	bw.WriteByte(byte(msg.Durability))
	bw.WriteByte(byte(msg.Flag))
	return bw.Bytes(), nil
}

// pack(1)
type MsgTalk struct {
	Target int
}

func (msg *MsgTalk) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)
	// target
	var target uint16
	err := binary.Read(br, binary.BigEndian, &target)
	if err != nil {
		return err
	}
	msg.Target = int(target)
	return nil
}

// pack(1)
type MsgTalkReply struct {
	Result      int
	SuccessRate [7]int
}

func (msg *MsgTalkReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Result))
	for i := range msg.SuccessRate {
		bw.WriteByte(byte(msg.SuccessRate[i]))
	}
	return bw.Bytes(), nil
}

type MsgCloseTalkWindow struct{}

func (msg *MsgCloseTalkWindow) Unmarshal(buf []byte) error {
	return nil
}

// pack(1)
type MsgTypeItemListReply struct {
	Type int
	MsgItemListReply
}

func (msg *MsgTypeItemListReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Type))
	data, err := msg.MsgItemListReply.Marshal()
	if err != nil {
		return nil, err
	}
	bw.Write(data)
	return bw.Bytes(), nil
}

type MsgBuyItem struct {
	Position int
}

func (msg *MsgBuyItem) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)
	position, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Position = int(position)
	return nil
}

type MsgBuyItemReply struct {
	Result int // -1=failed position=success
	Item   *item.Item
}

func (msg *MsgBuyItemReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Result))
	itemFrame, err := msg.Item.Marshal()
	if err != nil {
		return nil, err
	}
	bw.Write(itemFrame)
	return bw.Bytes(), nil
}

type MsgSellItem struct {
	Position int
}

func (msg *MsgSellItem) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)
	position, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Position = int(position)
	return nil
}

type MsgSellItemReply struct {
	Result int // 0=failed 1=success
	Money  int
}

func (msg *MsgSellItemReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Result))
	binary.Write(&bw, binary.LittleEndian, uint32(msg.Money))
	return bw.Bytes(), nil
}

type MsgMuunSystem struct{}

func (msg *MsgMuunSystem) Unmarshal(buf []byte) error {
	return nil
}

type MsgStatSpec struct {
	ID  int
	Min int
	Max int
}

type MsgStatSpecReply struct {
	Options []*MsgStatSpec
}

func (msg *MsgStatSpecReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	for _, v := range msg.Options {
		binary.Write(&bw, binary.LittleEndian, uint16(v.ID))
		binary.Write(&bw, binary.LittleEndian, uint16(v.Min))
		binary.Write(&bw, binary.LittleEndian, uint16(v.Max))
	}
	if len(msg.Options)%2 != 0 {
		bw.Write([]byte{0, 0}) // padding
	}
	return bw.Bytes(), nil
}

type MsgWarehouseMoneyReply struct {
	Result         int
	WarehouseMoney int
	InventoryMoney int
}

func (msg *MsgWarehouseMoneyReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Result))
	binary.Write(&bw, binary.LittleEndian, uint32(msg.WarehouseMoney))
	binary.Write(&bw, binary.LittleEndian, uint32(msg.InventoryMoney))
	return bw.Bytes(), nil
}

type MsgCloseWarehouseWindow struct{}

func (msg *MsgCloseWarehouseWindow) Unmarshal(buf []byte) error {
	return nil
}

type MsgCloseWarehouseWindowReply struct{}

func (msg *MsgCloseWarehouseWindowReply) Marshal() ([]byte, error) {
	return nil, nil
}

// pack(1)
type MsgMapMove struct {
	// RandomNumber int
	MoveIndex int
}

func (msg *MsgMapMove) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// randomNumber
	var randomNumber uint32
	err := binary.Read(br, binary.LittleEndian, &randomNumber)
	if err != nil {
		return err
	}
	// msg.RandomNumber = int(randomNumber)

	// moveIndex
	var moveIndex uint16
	err = binary.Read(br, binary.LittleEndian, &moveIndex)
	if err != nil {
		return err
	}
	msg.MoveIndex = int(moveIndex)

	return nil
}

type MsgMapMoveReply struct {
	// 0=success
	Result int
}

func (msg *MsgMapMoveReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Result))
	return bw.Bytes(), nil
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
	if size > 14 {
		return fmt.Errorf("MsgMove size invalid [size]%d", size)
	}
	dirs := make([]int, size)
	path := make(maps.Path, size)
	for i := range path {
		path[i] = &maps.Pot{}
	}
	for i := 0; i < int(size); i++ {
		if i == 0 {
			dir := int(bufDir[(i+2)/2] >> 4 & 0x0F)
			dirs[i] = dir
			dirPot := maps.Dirs[dir]
			path[i].X = int(x) + dirPot.X
			path[i].Y = int(y) + dirPot.Y
			continue
		}
		dir := 0
		if i%2 == 0 {
			dir = int(bufDir[(i+2)/2] >> 4 & 0x0F)
		} else {
			dir = int(bufDir[(i+2)/2] & 0x0F)
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

type MsgEnableCharacterClassReply struct {
	Class int
}

func (msg *MsgEnableCharacterClassReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Class))
	return bw.Bytes(), nil
}

// pack(1)
type MsgMiniMapReply struct {
	ID          int
	IsNpc       int
	DisplayType int
	Type        int
	X           int
	Y           int
	Name        string
}

func (msg *MsgMiniMapReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.ID))
	bw.WriteByte(byte(msg.IsNpc))
	bw.WriteByte(byte(msg.DisplayType))
	bw.WriteByte(byte(msg.Type))
	bw.WriteByte(byte(msg.X))
	bw.WriteByte(byte(msg.Y))
	gbk, err := simplifiedchinese.GBK.NewEncoder().String(msg.Name)
	if err != nil {
		return nil, err
	}
	var name [31]byte
	copy(name[:], gbk)
	bw.Write(name[:])
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

type MsgAttackSpeedReply struct {
	AttackSpeed int
	MagicSpeed  int
}

func (msg *MsgAttackSpeedReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	binary.Write(&bw, binary.LittleEndian, uint32(msg.AttackSpeed))
	binary.Write(&bw, binary.LittleEndian, uint32(msg.MagicSpeed))
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
	Inventory   [9]*item.Item
	GuildStatus int
	PKLevel     int
}

func MakeCharacterFrame(Class, ChangeUp int, Inventory [9]*item.Item) [18]byte {
	var chars [18]byte

	// class
	class := byte(Class << 5)
	switch ChangeUp {
	case 1:
		class |= 0x10
	case 2:
		class |= 0x18
	}
	chars[0] = class

	// inventory
	inventory := make([]*item.Item, len(Inventory))
	for i, v := range Inventory {
		if v == nil {
			inventory[i] = &item.Item{Section: 15, Index: 511}
		} else {
			inventory[i] = Inventory[i]
		}
	}

	// slot0~slot1 index -> chars[1]~chars[2]
	chars[1] = byte(inventory[0].Index)
	chars[2] = byte(inventory[1].Index)

	// slot0~slot1 section -> chars[12]bit4~-bit7~chars[13]bit4~bit7
	chars[12] = byte(inventory[0].Section << 5)
	chars[13] = byte(inventory[1].Section << 5)

	// slot2~slot6 index -> chars[3]~chars[5]
	chars[3] = byte(inventory[2].Index&0x0F<<4 | inventory[3].Index&0x0F)
	chars[4] = byte(inventory[4].Index&0x0F<<4 | inventory[5].Index&0x0F)
	chars[5] = byte(inventory[6].Index & 0x0F << 4)

	// slot2~slot6 index extention1 -> chars[9] bit3~bit7
	extend := inventory[2].Index&0x10<<3 |
		inventory[3].Index&0x10<<2 |
		inventory[4].Index&0x10<<1 |
		inventory[5].Index&0x10<<0 |
		inventory[6].Index&0x10>>1
	chars[9] = byte(extend)

	// slot2~slot6 index extention2 -> chars[13]~chars[15]
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

	return chars
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
		gbk, err := simplifiedchinese.GBK.NewEncoder().String(c.Name)
		if err != nil {
			return nil, err
		}
		var name [10]byte
		copy(name[:], gbk)
		bw.Write(name[:])
		bw.WriteByte(0) // padding 1 byte

		// level
		binary.Write(&bw, binary.LittleEndian, uint16(c.Level))

		// ctlcode
		bw.WriteByte(byte(c.CtlCode))

		// chars
		chars := MakeCharacterFrame(c.Class, c.ChangeUp, c.Inventory)
		bw.Write(chars[:])

		bw.WriteByte(byte(c.GuildStatus))
		bw.WriteByte(byte(c.PKLevel))
		bw.WriteByte(0) // padding 1 byte
	}
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
	utf8, err := simplifiedchinese.GBK.NewDecoder().Bytes(name[:])
	if err != nil {
		return err
	}
	msg.Name = string(bytes.TrimRight(utf8[:], "\x00"))

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
	gbk, err := simplifiedchinese.GBK.NewEncoder().String(msg.Name)
	if err != nil {
		return nil, err
	}
	var name [10]byte
	copy(name[:], gbk)
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
	utf8, err := simplifiedchinese.GBK.NewDecoder().Bytes(name[:])
	if err != nil {
		return err
	}
	msg.Name = string(bytes.TrimRight(utf8[:], "\x00"))

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
	utf8, err := simplifiedchinese.GBK.NewDecoder().Bytes(name[:])
	if err != nil {
		return err
	}
	msg.Name = string(bytes.TrimRight(utf8[:], "\x00"))

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
	LevelPoint         int
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
	binary.Write(&bw, binary.LittleEndian, uint16(msg.LevelPoint))
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

// pack(1)
type MsgItemListReply struct {
	Items []*item.Item
}

func (msg *MsgItemListReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	count := 0
	for i, item := range msg.Items {
		if item == nil {
			continue
		}
		count++
		bw.WriteByte(byte(i))
		data, err := item.Marshal()
		if err != nil {
			return nil, err
		}
		bw.Write(data)
	}
	var bw2 bytes.Buffer
	bw2.WriteByte(byte(count))
	bw2.Write(bw.Bytes())
	return bw2.Bytes(), nil
}

type MsgSkillListReply struct {
	Skills []*skill.Skill
}

func (msg *MsgSkillListReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(len(msg.Skills)))
	bw.WriteByte(0)
	for i, v := range msg.Skills {
		bw.WriteByte(byte(i))
		data, _ := v.Marshal()
		bw.Write(data)
	}
	return bw.Bytes(), nil
}

type MsgSkillOneReply struct {
	Flag  int // -2=add -1=delete
	Skill *skill.Skill
}

func (msg *MsgSkillOneReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.WriteByte(byte(msg.Flag))
	bw.WriteByte(0)
	bw.WriteByte(byte(0))
	data, _ := msg.Skill.Marshal()
	bw.Write(data)
	return bw.Bytes(), nil
}

type MsgMapDataLoadingOK struct{}

func (msg *MsgMapDataLoadingOK) Unmarshal(buf []byte) error {
	return nil
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
	utf8, err := simplifiedchinese.GBK.NewDecoder().Bytes(name[:])
	if err != nil {
		return err
	}
	msg.Name = string(bytes.TrimRight(utf8[:], "\x00"))

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
	gbk, err := simplifiedchinese.GBK.NewEncoder().String(msg.Name)
	if err != nil {
		return nil, err
	}
	var name [10]byte
	copy(name[:], gbk)
	bw.Write(name[:])

	// result
	bw.WriteByte(byte(msg.Result))

	return bw.Bytes(), nil
}

// pack(1)
type MsgMasterDataReply struct {
	MasterLevel          int
	MasterExperience     int
	MasterNextExperience int
	MasterPoint          int
	MaxHP                int
	MaxMP                int
	MaxSD                int
	MaxAG                int
}

func (msg *MsgMasterDataReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MasterLevel))
	binary.Write(&bw, binary.BigEndian, uint64(msg.MasterExperience))
	binary.Write(&bw, binary.BigEndian, uint64(msg.MasterNextExperience))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MasterPoint))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MaxHP))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MaxMP))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MaxSD))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MaxAG))
	return bw.Bytes(), nil
}

type MsgLearnMasterSkill struct {
	SkillIndex int
}

func (msg *MsgLearnMasterSkill) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)
	var SkillIndex uint32
	err := binary.Read(br, binary.LittleEndian, &SkillIndex)
	if err != nil {
		return err
	}
	msg.SkillIndex = int(SkillIndex)
	return nil
}

type MsgMasterSkill struct {
	MasterSkillUIIndex   int
	MasterSkillLevel     int
	MasterSkillCurValue  float32
	MasterSkillNextValue float32
}

type MsgLearnMasterSkillReply struct {
	Result           int
	MasterPoint      int
	MasterSkillIndex int
	MsgMasterSkill
}

func (msg *MsgLearnMasterSkillReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	binary.Write(&bw, binary.LittleEndian, uint16(msg.Result))
	binary.Write(&bw, binary.LittleEndian, uint16(msg.MasterPoint))
	binary.Write(&bw, binary.LittleEndian, uint32(msg.MasterSkillUIIndex))
	binary.Write(&bw, binary.LittleEndian, uint32(msg.MasterSkillIndex))
	binary.Write(&bw, binary.LittleEndian, uint32(msg.MasterSkillLevel))
	binary.Write(&bw, binary.LittleEndian, uint32(msg.MasterSkillCurValue))
	binary.Write(&bw, binary.LittleEndian, uint32(msg.MasterSkillNextValue))
	return bw.Bytes(), nil
}

type MsgMasterSkillListReply struct {
	Skills []*MsgMasterSkill
}

func (msg *MsgMasterSkillListReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	bw.Write([]byte{0, 0, 0}) // padding
	binary.Write(&bw, binary.LittleEndian, uint32(len(msg.Skills)))
	for _, v := range msg.Skills {
		bw.WriteByte(byte(v.MasterSkillUIIndex))
		bw.WriteByte(byte(v.MasterSkillLevel))
		bw.Write([]byte{0, 0}) // padding
		binary.Write(&bw, binary.LittleEndian, uint32(v.MasterSkillCurValue))
		binary.Write(&bw, binary.LittleEndian, uint32(v.MasterSkillNextValue))
		binary.Write(&bw, binary.LittleEndian, uint32(0)) // unknown field
	}
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

type MsgResetCharacterReply struct {
	Reset string
}

func (msg *MsgResetCharacterReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer

	// reset
	bw.WriteString(msg.Reset)

	return bw.Bytes(), nil
}

// pack(1)
type MsgResetGameReply struct {
}

func (msg *MsgResetGameReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	binary.Write(&bw, binary.LittleEndian, uint16(0))
	return bw.Bytes(), nil
}

type MsgLive struct {
	Time         int
	AttackSpeed  int
	Agility      int
	MagicSpeed   int
	Version      string
	ServerSeason int
}

type MsgSkillAttack struct {
	Target int
	Skill  int
}
