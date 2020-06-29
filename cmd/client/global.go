package main

import (
	"encoding/binary"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/xujintao/balgass/win"
)

var v0131A240 uint32
var v0131A250 uint32
var v0131A26C bool
var v0131A26D bool
var v0131A270 int

var v08C88AB8 [100]uint8
var v08C88C69 bool
var v08C88C6A bool
var v08C88C74 uint32
var v086A3B94 func()
var v086A3BEC uint32

var v0638C6C4 uint32 = 0x4EB0E616

var v0114EE50 = [...]uint64{
	0x3FE3333340000000,
	0x3FC3333340000000,
	0x3F70624DE0000000,
}

var v01319F28 uint32 = 0x3F2DE762

// f004E46B3handleState2
var v09FB8736 uint32 = 0x2B2F4A3D
var v0A4E24C6 uint32 = 0x4849459E
var v09E035DB uint32 = 0x5433BC5D
var v0A43DD91 uint32 = 0x653CAB4B
var v0B06F40C uint32 = 0x08244423 // v0B06F40C = 0x310F1716 = v09E035DB ^ v0A43DD91
var v09FDFB22 uint32 = 0xFC5D207C
var v0AFD3C52 uint32 = 0xD577
var v09FB6D69 uint32 = 0xDEA7
var v0A441BD1 uint32 = 0x38
var v0AD2DD3A uint32 = 0xE4E2E401
var v0AF890C3 uint32 = 9
var v0B287022label1 uint32 = 0x0AD98B31
var v09E8DF92 uint32 = 8
var v0AAB7324label2 uint32 = 0x0AD2CAED
var v0A933705 uint32 = 7
var v0A935B85label3 uint32 = 0x0A43B47B
var v0ABFAB88blocks = [...]block{
	{0x000393EA, 0x0192},
	{0x00039581, 0x0100},
	{0x0003968F, 0x05},
	// ...
}
var v0A4E85ABimageBase uintptr = 0x00400000
var v0AA2E05D uint32 = 0x5AA7D865
var v09EB9D65 uint32 = 1

var v0A9360CB uint32 = 0x4B0FAFAC
var v0AFD365B uint32 = 0x5EC66E14 // v0AFD365B = v0A9360CB

type t7 struct {
	data [0x2A4]byte
}

var v08C86C50 = [10]t7{}

var v0114EC40 float64

var v012E234C = "data\\music\\mutheme.mp3"
var v012E239C = "data\\music\\main_theme.mp3"

// f004E17B9handleState4
var v0A443F74 uint32 = 0xE9497BA1
var v0AD3B896 uint32 = 0x13EFEDCE // v0AD3B896 = v0A443F74
var v0A890E43 uint32 = 0xC491
var v0A56EB4C uint32 = 0xE521
var v0A84A81D uint32 = 0x3E92EB75 // 0xEB
var v09FE1805 uint32 = 0x125      // v0A84A81D 与 v09FE1805比较

var v0ABE324Flabel1 uint32 = 0x09EBC65D
var v0A952A97 uint32 = 8
var v0A92FC07label2 uint32 = 0x09FDF5DF
var v0A920E49 uint32 = 7
var v0A32B0C2label3 uint32 = 0x0AFDEEEE
var v0A6039A7 uint32 = 6

var v0AFD3C89 uint32 = 0x9B20B37E
var v0AA30425 uint32 = 0x105DA588

var v0A88819B uint32 = 0xD266CC5A // key
var v0AF96824backupCode = [...]uint32{
	0x37664E4F,
	0x36FE23B9,
}

var v0A9F69B2blocks = [...]block{
	{0x000393EA, 0x192}, // f004393EA
	{0x00039581, 0x100},
	{0xFFFFFFFF, 0xFFFFFFFF},
}
var v09FE37F3imageBase uintptr = 0x00400000

var v0ABD75F3 uint32 = 0xE55020DE
var v0A3268A5 uint32 = 0xC01C390C

var v0A0C2FD2 uint32 = 0x05BE6BE3
var v0AF8FD45 uint32 = 0x34B17CF5
var v0AA09787 uint32 = 0x0824448B // v0AA09787 = v0A0C2FD2 ^ v0AF8FD45

var v0AD6FFBE uint32 = 0x90B9
var v0ABB4649 uint32 = 0x975D
var v0AAB36E2 uint32 = 0x531BDC38 // v0AD6FFBE * v0AAB36E2 % v0ABB4649
var v09EAE2F8 uint32 = 0x26

var v0AF781CF uint32 = 0x35C75981
var v0A33BBC9blocks = [...]block{}
var v0A5FCB8Bcodes = [...]uint32{}
var v0A84B964 uint32 = 1

var v09FFF8E0imageBase uintptr = 0x00400000

var v0A9FB214 uint32 = 7
var v0A443580label1 uint32 = 0x09EB9FEC
var v0AFDFD32 uint32 = 6
var v09FB8B32label2 uint32 = 0x0AC9CB58
var v0A849070 uint32 = 5
var v0A32C145label3 uint32 = 0x0A338F26

// f004E0E03handleState5
// 确保第一次不覆盖
var v0A56C1BA uint32 = 0x865EB72E
var v0AF8AF60 uint32 = 0xDCD68D02
var v0AD73D05 uintptr = 0xD2435CB5
var v0ABD3FEA uintptr = 0xD36C2725
var v0B07C274 uintptr = 0x09FDEC49 // v0AD73D05 ^ v0ABD3FEA - v0B07C274 + 0x09FDEC49

// 后面在一个确定的时机进行覆盖
var v0A339687 uint32 = 0xB2F9
var v09F87E15 uint32 = 0x8E25
var v09DEA85C uint32 = 0x76237DA5 //
var v0AC3A339 uint32 = 0xE5

var v0A9FD787backupCode = [...]uint32{}
var v0A56AF60blocks = [...]block{
	{0x000393EA, 0x192},
	{0x00039581, 0x100},
	{0xFFFFFFFF, 0xFFFFFFFF},
}

var v09E72FAA uint32 = 0x03B1C032
var v0A933991 uint32 = 0x922C3447 // v0A933991 = v09E72FAA

var v0B07545C uint32 = 8
var v0ABB631Flabel1 uint32 = 0x0A84CD1A
var v0A38E67E uint32 = 7
var v0AC36B4Dlabel2 uint32 = 0x09FC4127
var v0AF87A7B uint32 = 6
var v0AA2D054label3 uint32 = 0x0AF7BD1B

// f00A0A5E1 shell
var v0AA08598 uint32 = 0xE2209428
var v0AD7A978 uint32 = 0xABF8FFC9
var v0A9D361D uint32 = 0xD2B99D3B
var v0A391C78 uint32 = 0x52B79D4F // 0x7FFE0014 = v0A391C78 - v0A9D361D
var v0AD82FDD uint32 = 0x9EA1
var v0A338B93 uint32 = 0xE003
var v0A562D19 uint32 = 0x5AFBC57D
var v09F876D2 uint32 = 0x4334

var v0AD93A11 uint32 = 0x1A65E9C0
var v0AD93A0D = [...]uint32{0x078202AD}

var v0AD574B8 uint32 = 8
var v0A0519E0label1 uint32 = 0x09E2729C

type stdstring struct {
	m04data []uint8
	m14len  int
	m18cap  int
}

func (t *stdstring) f00406A20init() {
	t.m04data = nil
	t.m14len = 0
	t.m18cap = 15
}

func (t *stdstring) f00406EB0(buf []uint8, len int) {
	// s长度小于16就放在栈上(m04~m13)，否则存在堆上
	t.m04data = buf
	t.m14len = len
	t.m18cap = 0x1F
}

func (t *stdstring) f0043D7E2stdstring(s string) {
	t.f00406EB0([]uint8(s), len(s))
}

func (t *stdstring) f00406FC0stdstring(buf []uint8) {
	t.f00406EB0(buf, len(buf))
}

func (t *stdstring) f004079A0stdstring(s *stdstring) {

}

func (t *stdstring) f004073E0cstr() []uint8 {
	return t.m04data
}

func (t *stdstring) f00407B10free() {
	f00DE7538free(nil)
}

// ----------------------xstring----------------------------------
var v012F4D80empty emptystring = 0x80000000
var v012F4D84cnt uint32

func f00BAC850xstring() *xstring {
	atomic.AddUint32(&v012F4D84cnt, 1)
	return (*xstring)(unsafe.Pointer(&v012F4D80empty))
}

// type xstring interface {
// 	f00A3AF8Acstr() string
// 	f00A3AF9Ecstr() string
// 	f00BAC0B0len()
// }

type emptystring uintptr

type xstring struct {
	m00size int
	m04cnt  int32
	m08data [256]uint8
}

func f00BAC8C0xstring(mm *mm, l1 int, unk int, data []uint8, l2 int) *xstring {
	var w *xstring
	if l1 == 0 {
		atomic.AddUint32(&v012F4D84cnt, 1)
		w = (*xstring)(unsafe.Pointer(v012F4D80empty))
	} else {
		// mm.do11malloc(uint(l1)+0xC, 0) // 0x13, 0
		w = new(xstring)
		w.m08data[8+l1] = 0
		w.m04cnt = 1
		w.m00size = l1
	}
	f00DE7C90memcpy(w.m08data[:], data, l2)
	return w
}

func (t *xstring) f00BADDD0xstring(s string) {
	l := len(s)
	t = f00BAC8C0xstring(v09D9BD74mm, l, 0, []uint8(s), l)
}

func (t *xstring) f00BACD30xstring(s string) {

}

func (t *xstring) f00BAE080xstring(s *xstring, offset int, l int) *xstring {
	if offset >= t.m00size || offset >= l {
		// s = &v012F4D80empty
		// return s
	}
	s = f00BAC8C0xstring(v09D9BD74mm, l, 0, t.m08data[offset:], l)
	return s
}

func (t *xstring) f00BADF50assign(s *xstring) *xstring {
	atomic.AddInt32(&s.m04cnt, 1)
	return s
}

func (t *xstring) f00BACE40assign(s *xstring) {

}

func (t *xstring) f00A3AF6C() string {
	return string(t.m08data[:])
}

func (t *xstring) f00A3AF8Acstr() string {
	return t.f00A3AF6C()
}

func (t *xstring) f00A3AF9Ecstr() string {
	return t.f00A3AF6C()
}

func f00BF0C80(data []uint8, l int) int {
	cnt := 0
	if l == -1 {
		return cnt
	}
	if l == 0 {
		return cnt
	}
	for {
		if data[0] == 0 { // f00BF0630(data)
			return cnt
		}
		data = data[1:]
		cnt++
	}
}

func (t *xstring) f00BAC0B0len() int {
	if t.m00size < 0 {
		return ^(1 << 31) & t.m00size
	}
	size := f00BF0C80(t.m08data[:], t.m00size)
	if t.m00size == size {
		t.m00size |= 1 << 31
	}
	return size
}

func (t *xstring) f00BAE060cat(sub *xstring) *xstring {
	return nil
}

func (t *xstring) f00A4FF11cat(s string) {

}

func (t *xstring) f00A3AF16destruct() {

}

func (t *xstring) destruct() {
	atomic.AddInt32(&t.m04cnt, -1)
	if t.m04cnt == 0 {
		// delete
	}
}

// load serverlist.bmd
type server struct { // 0x54
	// head
	m00 [32]uint8 // "电信1区"
	m20 uint8     // 0
	m21 uint8     // 1,2
	m22 bool      // 0
	m23 [20]uint8 // ?
	// body
	m38 stdstring
}

func (s *server) f00AF88D8server()           {}
func (s *server) f00AF8903server(s1 *server) {}

func f00AF7DC3getServerListManager() *serverListManager {
	v09D96728.Do(func() {
		v09D965B0serverListManager.f00AF7CC6()
	})
	return &v09D965B0serverListManager
}

var v09D965B0serverListManager serverListManager
var v09D96728 sync.Once

type serverListManager struct {
	m2Ccount int

	// 结构
	m13C *treeNode
	size int
}

// 构造函数
func (t *serverListManager) f00AF7CC6() {}

func (t *serverListManager) f00AF7EB6dec(buf []uint8, size uint) {

}

func (t *serverListManager) f00AF7F07load() {
	// 0x66C局部变量
	// ebp670 := t
	ebp54 := f00DE909Efopen("Data/Local/ServerList.bmd", "rb")
	var ebp10 uint = 0x3B // 59
	var ebp4AC server
	ebp4AC.f00AF88D8server()

	// unmarshal
	var ebp50head [60]uint8  // head 59字节
	var ebp454body [10]uint8 // body
	for {
		if f00DE8FBDfread(ebp50head[:], ebp10, 1, ebp54) == 0 {
			break
		}
		t.f00AF7EB6dec(ebp50head[:], ebp10)
		size := uint(binary.LittleEndian.Uint16(ebp50head[58:]))

		f00DE8FBDfread(ebp454body[:], size, 1, ebp54)
		t.f00AF7EB6dec(ebp454body[:], size)

		x := binary.LittleEndian.Uint16(ebp50head[:])
		if x >= 500 {

		}
		f00DE9370strncpy(ebp4AC.m00[:], ebp50head[2:], 32)
		ebp4AC.m20 = ebp50head[34]
		ebp4AC.m21 = ebp50head[35]
		if ebp50head[36] != 0 {
			ebp4AC.m22 = true
		} else {
			ebp4AC.m22 = false
		}
		ebp14 := 0
		for {
			if ebp14 >= 20 {
				break
			}
			ebp4AC.m23[ebp14] = ebp50head[37+ebp14]
			ebp14++
		}
		ebp4AC.m38.f0043D7E2stdstring(string(ebp454body[:]))

		var tmp server
		tmp.f00AF8903server(&ebp4AC)
		// f00AFC1D3(&ebp660, binary.LittleEndian.Uint32(ebp50[:]))
		// ebp604.f00AFC22B(ebp678)
		// t.m13C.f00AF8D6E(&ebp66C, &ebp604)
		// ebp604.f00AF897B()
		// ebp660.f00AF8967()
	}
}

func (t *serverListManager) f00AF883FsetCount(count int) {
	t.m2Ccount = count
}

func (t *serverListManager) f00AF8853getCount() int {
	return t.m2Ccount
}

func (t *serverListManager) f00AF81FF(code, percent int) {
	// 0x58局部变量
	// ebp5C := t
	// t.m04.f00AF89F1(&ebp18)
	// ebp19 := false
	// ebp18.f00AF8CCF(t.m04.f00AF8A14(&ebp24))
	// ebp34 := f00DE852Fnew()
	// 	ebp60 := ebp34.f00AF6EB3()
	// 	ebp30 := ebp60
	// 	ebp10 := ebp30
	// 	t.f00AF8387(code/20, ebp10)
	// 	t.m04.f00AF8A42(ebp50, ebp44.f00AF8F97(ebp10, ebp10.m04))
	// 	t.f00AF846C(ebp10, code, percent)
	// 	t.m04.f00AF89F1(&ebp58)
}

// objectManager
var v08C88E58hpMax uint16
var v08C88E5AmpMax uint16
var v08C88E5CsdMax uint16
var v08C88E5EsdMax uint16
var v086105ECobject = &object{}
var v0805BBACself = &object{}

// object sizeof=0x6BC
type object struct {
	m8        uint32
	m13       uint8
	m24       bool
	m2B       uint8
	m38name   [32]uint8
	m5Eid     uint16
	m10C      uint16
	m122hp    uint16 // hp
	m124mp    uint16 // mp
	m126hpMax uint16 // max hp
	m128mpMax uint16 // max mp
	m12Asd    uint16 // sd
	m12CsdMax uint16 // max sd
	m140ag    uint16 // ag
	m142agMax uint16 // max ag
	m154      uint16
	m160      uint16
	m166      uint16
	m178      [100]uint8
	m410      struct {
		m04 bool
		m0E bool
		m28 uint8
	}
	m438 uint8
}

var v012E31B0 int = 0
var v012E3200 int = 5
var v01308D04objectManager = &objectManager{}

var v01308D08 battleCoreSpec
var v01308D10 sync.Once

func f0043A2DFgetBattleCoreSpec() *battleCoreSpec {
	v01308D10.Do(func() {
		return &v01308D08
	})
}

type battleCoreSpec struct {
	m00       uintptr // &v0114D394
	m04Enable bool
}

func (b *battleCoreSpec) f0043A33Cenable() bool {
	return b.m04Enable
}

type objectManager struct {
	m08objects []object
}

func (t *objectManager) f00A38D5BgetObject(index int) *object {
	if index < 0 || index >= 0x190 {
		return nil
	}
	return &t.m08objects[index]
}

type mainFrame struct {
	m30   uintptr
	m2844 []uint8
}

func (t *mainFrame) f00BB1E00(buf []uint8, format string, param ...interface{}) {
	if t.m30 == 0 {
		return
	}
	// var ebp4 struct{}
	// ebp4.f00C33740(t.m2844, 1, 0x1B, 0)
	// var ebp2C struct{}
	// f00DF75E5(&ebp2C, 0, 0)
	// f00DF75E5(&ebp28, 0x10000, 0x30000)
}

// game
func f00A49798game() *game {
	if v09D96438 == nil {
		// ebp14 := f00DE852Fnew(0x228) // 0x0DDBB9F0
		var ebp18 *game
		ebp14 := &game{}
		if ebp14 != nil {
			ebp18 = ebp14.f00A49808construct()
		} else {
			ebp18 = nil
		}
		v09D96438 = ebp18
	}
	return v09D96438
}

var v09D96438 *game

// size:0x228 虚表:v01174338
type game struct {
	m00vtabptr  []uintptr
	m08         windowManager1
	m18         struct{}
	m1C         struct{}
	m20         struct{}
	m24         struct{}
	m2C         struct{}
	m28         int
	m30         *xstring // Taiwan
	m34         int
	m38         struct{}
	m58         struct{}
	m68hwndNext win.HWND

	m6Cstate int
	m70      *windowManager01173C54
	m74      uintptr                // state2
	m78      uintptr                // state4
	m7C      *windowManager01173C54 // state5

	m80            *windowgame0116A864
	m84infoTooltip *windowgameInfoTooltip
	m88caution     *windowgameCaution   // Caution
	m8CmainFrame   *windowgameMainFrame // MainFrame
	m90            uintptr              // characterFrame
	m94            uintptr              // inventoryFrame
	m98            uintptr              // ExtensionBagFrame
	m9C            uintptr              // PrivateStoreFrame
	mA0            uintptr              // MixInventoryFrame
	mA4            uintptr              // MixSel
	mA8            uintptr              // MixGoblinSel
	mAC            uintptr              // MixStone
	mB0            uintptr              // MixInventoryFrame
	mB4            uintptr              // GensRanking
	mB8            uintptr              // GuildCreateFrame
	mBC            uintptr              // GuildInfoFrame
	mC0            uintptr              // GuildPosition
	mC4            uintptr              // infoPopup
	mC8            uintptr              // PetInfoPopup
	mCC            uintptr              // ?
	mD0            uintptr              // QuickCommand
	mD4            uintptr              // duel
	mD8            uintptr              // duelwatch
	mDC            uintptr              // trade
	mE0            uintptr              // purchaseinven
	mE4            uintptr              // goldarcher
	mE8            uintptr              // ?
	mEC            uintptr              // npcdialogue
	mF0            uintptr              // NpcShop
	mF4            uintptr              // PetInfo
	mF8            uintptr              // Pentagram
	mFC            uintptr              // ArkaResultInfo
	m100           uintptr              // StorageInven
	m104           uintptr              // StorageInven
	m108           uintptr              // NpcJobChangeQuest
	m10C           uintptr              // NpcQuestProgress
	m110           uintptr              // NpcGateKeeper
	m114           uintptr              // SystemMenu
	m118           uintptr              // Option
	m11C           uintptr              // SelectMenu
	m120           uintptr              // MyQuestInfo
	m124           uintptr              // ChatWindow
	m128           uintptr              // MoveCommand
	m12C           uintptr              // HelpWindow
	m130           uintptr              // EventMapHelper
	m134           uintptr              // Trainer 驯兽师
	m138           uintptr              // MapName
	m13C           uintptr              // Notice
	m140           uintptr              // HeroPosition
	m144           uintptr              // EnterBloodCastle
	m148           uintptr              // EnterEmpireGuardian
	m14C           uintptr              // CommandWindow
	m150           uintptr              // ArkaBattleRegister
	m154           uintptr              // ArkaBattleProgress
	m158           uintptr              // ArkaBattleNotice
	m15C           uintptr              // MonsterInfo
	m160           uintptr              // PlayerInfo
	m164           uintptr              // BuffFrame
	m168           uintptr              // EquipDurabilityInfo
	m16C           uintptr              // MacroMain
	m170           uintptr              // MacroSub
	m174           uintptr              // Alarm
	m178           uintptr              // MatchingSelect
	m17C           uintptr              // MatchingGuild
	m180           uintptr              // MatchingParty
	m184           uintptr              // PartyFrame
	m188           uintptr              // ItemInfo
	m18C           uintptr              // StrongestMatchMenu
	m190           uintptr              // StrongestMatchRank
	m194           uintptr              // StrongestMatchResult
	m198           uintptr              // SearchPrivateStore
	m19C           uintptr              // MiniGameRummy
	m1A0           uintptr              // EventInventoryFrame
	m1A4           uintptr              // EventInfo
	m1A8           uintptr              // CursedTempleMatchMenu
	m1AC           uintptr              // CursedTempleResult
	m1B0           uintptr              // CursedTempleScore
	m1B4           uintptr              // CursedTempleInfo
	m1B8           uintptr              // ListOfMatches
	m1BC           uintptr              // EventMapProgressInfo
	m1C0           uintptr              // EventMapRoundCount
	m1C4           uintptr              // ProgressFrame
	m1C8           uintptr              // Navimap
	m1CC           uintptr              // PetInventoryFrame
	m1D0           uintptr              // MuunExchange
	m1D4           uintptr              // OrdealSquareEnterMenu
	m1D8           uintptr              // OrdealSquareScore
	m1DC           uintptr              // OrdealSquareResult
	m1E0           uintptr              // EventMapTutorial
	m1E4           uintptr              // EventInfo
	m1E8           uintptr              // BattleField_BigDialog
	m1EC           uintptr              // RideButton
	m1F0           uintptr              // MonsterCompensation
	m1F4           uintptr              // ImageNotice
	m1F8           uintptr              // GremoryCase 活动奖励背包
	m1FC           uintptr              // RebuyList
	m200           uintptr              // PetFrame
	m204           uintptr              // do5 ebp7C

	m20Cwidth  int
	m210height int
}

func (t *game) f00A49808construct() *game {
	// 0x34局部变量
	// ebp34 := t
	// 构造子窗口1
	// ebp18 := f00A3BA10newobject(0x10) // 0x0EB41EB0
	// if ebp18 == nil {
	// 	ebp38 = nil
	// } else {
	// 	ebp38 = ebp18.f00BBBA00construct()
	// }
	// equal
	ebp18 := new(t01187D5C) // 0x0EB41EB0
	ebp18.f00BBBA00construct()
	ebp38 := ebp18
	ebp14 := ebp38
	ebp10 := i011737C4(ebp14) // ebp10.f00A50545(ebp14)

	// // 构造子窗口2
	// ebp24 := f00A3BA10newobject(0x0C) // 0x0EB41ED0
	// if ebp24 == nil {
	// 	ebp3C = nil
	// } else {
	// 	ebp3C = ebp24.f00A50053()
	// }
	// equal
	ebp24 := new(t0117437C) // 0x0EB41ED0
	ebp24.f00A50053construct()
	ebp3C := ebp24
	ebp20 := ebp3C
	ebp1C := i011737C4(ebp20) // ebp1C.f00A50576(ebp20)

	// // 构造子窗口3
	// ebp30 := f00A3BA10newobject(0x0C) // 0x0EB41EE0
	// if ebp30 == nil {
	// 	ebp40 = nil
	// } else {
	// 	ebp40 = ebp30.f00A50036()
	// }
	// equal
	ebp30 := new(window01174368) // 0x0EB41EE0
	ebp30.f00A50036construct()
	ebp40 := ebp30
	ebp2C := ebp40
	ebp28 := i011737C4(ebp2C) // ebp28.f00A504E3(ebp2C)

	// register
	t.m08.f00BB0D20construct(ebp28, ebp1C, ebp10)
	// ebp28.f00A504F9()
	// ebp1C.f00A5058C()
	// ebp10.f00A5055B()
	// t.m18 = 0 // t.m18.f00A4FCF7(0)
	// t.m1C = 0 // t.m1C.f00A4FCF7(0)
	// t.m20 = 0 // t.m20.f00A4FC2C(0)
	// t.m24 = 0 // t.m24.f00A4FC78(0)
	// t.m28 = 0
	// t.m2C = 0 // t.m2C.f00A4FDC2()
	// t.m30 = 0 // t.m30.f00BAC850()
	// t.m34 = 0
	// t.m38 = 0 // t.m38.f00BAC850()
	// t.m58.f00AECCBC()
	return nil
}

func (t *game) f00A49E40init1(hWnd win.HWND, width, height int) {
	// ebp7C := t
	t.m20Cwidth = width
	t.m210height = height
	t.m68hwndNext = 0 // dll.user32.SetClipboardViewer(hWnd) // Adds the specified window to the chain of clipboard viewers.

	// 构造窗口
	// ebp30 := f00A3BA10newobject(0x20) // 0x0EB42090
	// if ebp30 == nil {
	// 	ebp80 = nil
	// } else {
	// 	ebp80 = ebp30.f00A4FAD6()
	// }
	// equal
	ebp30 := new(window01174340)
	ebp80 := ebp30
	ebp2C := ebp80
	ebp28 := ebp2C // ebp28.f00A50514()
	t.m08.do3(0x18, ebp28)
	// ...
}

func (t *game) f00A4E46DLoadResource(w iwindowgame0117373C, x int, name *xstring, y bool, align int) bool {
	// 0x1C局部变量
	// ebp20 := t
	if t.m7C == nil {
		ebp14 := new(windowManager01173C54)
		ebp14.f00A472A7construct()
		ebp24 := ebp14
		ebp10 := ebp24
		t.m7C = ebp10
	}
	// ebp25 := t.f00A4E520()
	ebp25 := func(wm *windowManager01173C54, w iwindowgame0117373C, x int, name *xstring, y bool, align int) bool {
		// 0x48局部变量
		// ebp3C := t
		ebp1D := false
		if w == nil || wm == nil {
			ebp1E := ebp1D
			// name.f00A3AF16destruct()
			return ebp1E
		}
		var ebp14dir *xstring
		ebp14dir.f00BADDD0xstring("./Data/Interface/GFx/")
		var ebp18name *xstring
		ebp10 := 5
		ebp1C := 0
		// ebp40 := t.f00A4F04E(&ebp18name, ebp14dir, name.f00A3AF9Ecstr()) // ./Data/Interface/GFx/xxx
		// 			   &t.m08  5  "...dir/Caution.ozg"              "Taiwan"  5    0     0  0
		// 			   &t.m08  1  "...dir/MainFrame.ozg"            "Taiwan"  5    0     0  0
		ebp45 := w.do2(&t.m08, x, ebp14dir.f00BAE060cat(ebp18name), t.m30, ebp10, ebp1C, 0, 0)
		if ebp45 == true {
			// w.do3()
			if align != 9 {
				f00A3A4F2(w, "SetAlign", "%d", align)
			}
			if y == true {
				w.f00A3A717(true)
			}
			wm.f00A47FAB(w)
		} else {
			// dll.user32.MessageBox(0, name.f00A3AF9Ecstr(), "Load Fail!!", 0)
			if w != nil {
				ebp34 := w
				ebp30 := ebp34
				ebp30.do1(true)
			}
		}
		ebp35 := ebp1D
		// ebp18name.f00A3AF16destruct()
		// ebp14dir.f00A3AF16destruct()
		// name.f00A3AF16destruct()
		return ebp35
	}(t.m7C, w, x, name, y, align)
	ebp15 := ebp25
	// name.f00A3AF16destruct()
	return ebp15
}

func (t *game) f00A4A521init2() {
	// 0x610局部变量
	// ebp494 := t

	// t.m80
	ebp14 := new(windowgame0116A864)
	ebp14.f008E20BBconstruct()
	ebp498 := ebp14
	ebp10 := ebp498
	t.m80 = ebp10

	// infoTooltip
	ebp1C := new(windowgameInfoTooltip)
	ebp1C.f00A3BBB9construct()
	ebp49C := ebp1C
	ebp18 := ebp49C
	t.m84infoTooltip = ebp18

	// t.m70
	t.m70 = nil

	// Caution
	ebp24 := new(windowgameCaution)
	ebp24.f00A99F11construct()
	ebp4A0 := ebp24
	ebp20 := ebp4A0
	t.m88caution = ebp20
	var tmpName *xstring
	tmpName.f00BADDD0xstring("Caution")
	t.f00A4E46DLoadResource(t.m88caution, 5, tmpName, false, 9)

	// MainFrame
	ebp3F4 := new(windowgameMainFrame)
	ebp3F4.f00AA9021construct()
	ebp5E8 := ebp3F4
	ebp3F0 := ebp5E8
	t.m8CmainFrame = ebp3F0
	tmpName.f00BADDD0xstring("MainFrame")
	t.f00A4E46DLoadResource(t.m8CmainFrame, 1, tmpName, true, 2)
	// ...
}

func (t *game) f00A4E1BFchangeState(state int) bool {
	// ebp10 := t
	if t.m6Cstate == state {
		return false
	}
	t.m6Cstate = state
	ebp14 := t.m6Cstate
	switch ebp14 {
	case 2:
		// t.m70 = t.m74
		// ebp4hIMC := dll.imm32.ImmGetContext(v01319D6ChWnd)
		// dll.imm32.ImmSetConversionStatus(ebp4hIMC, 0, 0)
		// dll.imm32.ImmReleaseContext(v01319D6ChWnd, ebp4hIMC)
	case 4:
		// t.m70 = t.m78
		// ebp8hIMC := dll.imm32.ImmGetContext(v01319D6ChWnd)
		// dll.imm32.ImmSetConversionStatus(ebp8hIMC, 1, 0)
		// dll.imm32.ImmReleaseContext(v01319D6ChWnd, ebp8hIMC)
	case 5:
		// t.m70 = t.m7C
		// ebpChIMC := dll.imm32.ImmGetContext(v01319D6ChWnd)
		// dll.imm32.ImmSetConversionStatus(ebpChIMC, 1, 0) // IME_CMODE_NATIVE, IME_CMODE_ALPHANUMERIC
		// dll.imm32.ImmReleaseContext(v01319D6ChWnd, ebpChIMC)
	}
	if t.m70 == nil {
		return true
	}
	// t.m70.f00A47353()
	// t.m70.f00A47C55()
	return true
}

func (t *game) f00A4DC94fresh(state int) {
	if t.m80 != nil {
		// t.m80.f008E27CFfresh()
	}
	if t.m70 != nil {
		t.m70.f00A47461fresh()
	}
}

//
func f004D9460get() *t01319F18 {
	v01319F24once.Do(func() {
		v01319F18.f004D9693construct()
	})
	ebp18 := &v01319F18 // ebp18 := f00405BF0(8, v01319F18)
	var ebp1C *t01319F18
	if ebp18 != nil {
		ebp1C = ebp18.f004D9D74vtab()
	} else {
		ebp1C = nil
	}
	ebp14 := ebp1C
	ebp10 := ebp14
	ebp10.m04 = &v01319F18
	v01319F18.m08 = true // v01319F20 = true
	return ebp10
}

var v01319F18 t01319F18
var v01319F24once sync.Once

type t01319F18 struct {
	m00 uintptr // 虚表指针: &v0114E1D4
	m04 *t01319F18
	m08 bool
}

func (t *t01319F18) f004D9693construct() {
	t.m08 = false
}

func (t *t01319F18) f004D9D74vtab() *t01319F18 {
	// t.f004DA0F4()
	func() {
		// t.f004DA118()
		func() {
			// t.f004DA135()
			func() {
				// t.m00 = &v0114E20C 基对象虚表
			}()
			// t.m00 = &v0114E1F0 基对象虚表
		}()
		// t.m00 = &v0114E21C 基对象虚表
		t.m04 = nil
	}()
	// t.m00 = &v0114E1D4 虚表
	return nil
}

// f00BAB280
func (t *t01319F18) do2(param uintptr) bool {
	// if v09D9D5E8 == nil {
	// 	v09D9BD40.f00BFA8B0()
	// }
	// f00BAAC30(param)
	if func(param uintptr) uintptr {
		// f00BFA710("Global", 0)
		func() {
			// 初始化 v09D9BD74mm: 0x0EB40B84
		}()
		return 0x0EB40B84
	}(param) != 0 {
		return true
	}
	return false
}
