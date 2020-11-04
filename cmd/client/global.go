package main

import (
	"encoding/binary"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/xujintao/balgass/win"
)

var v0131A233 uint8
var v0131A240 uint32
var v0131A250 uint32
var v0131A26C bool
var v0131A26D bool
var v0131A270 int

var v08C88AB8 [100]uint8
var v08C88C69 bool
var v08C88C6A bool
var v08C88C74 uint32

// var v086A3B94 func()
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

// stdstring
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

func (t *stdstring) f00407420len() int {
	return t.m14len
}

func (t *stdstring) f004078B0data() []uint8 {
	return t.m04data
}

func (t *stdstring) f00407B10free() {
	f00DE7538free(nil)
}

func (t *stdstring) f004AA08Elen() int {
	return t.m14len
}

func (t *stdstring) f004CE1BD(data []uint8, offset int, l int) int {
	if l == 0 && offset <= t.m14len {
		return offset
	}
	if offset >= t.m14len {
		return -1
	}
	ebp4 := t.m14len - offset
	if ebp4 < l {
		return -1
	}
	ebp4 = ebp4 - l + 1
	ebpCdata := t.f004078B0data()[offset:]
	for {
		// f004CE2B3(ebpCdata, ebp4, data)
		ebp8data := f00DEC9B0memchr(ebpCdata, data[0], ebp4)
		if ebp8data == nil {
			return -1
		}
		// if f0042DA80(ebp8data, data, l) == 0 {
		// 	return int(uintptr(unsafe.Pointer(&ebp8data[0])) - uintptr(unsafe.Pointer(&t.f004078B0data()[0])))
		// }
		ebp4 -= int(uintptr(unsafe.Pointer(&ebp8data[0])) - uintptr(unsafe.Pointer(&ebpCdata[0])) + 1)
		ebpCdata = ebp8data[1:]
	}
	return 0
}

func (t *stdstring) f004CE187(s *stdstring, offset int) int {
	return t.f004CE1BD(s.f004078B0data(), offset, s.f00407420len())
}

func f004CE0BDstrstr(s1 *stdstring, s2 *stdstring) bool {
	ebp8 := s1.f004CE187(s2, 0)
	ebp4 := s1.f004AA08Elen()
	if ebp8 <= ebp4 {
		return true
	}
	return false
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
	m138code uint8
	m13C     *treeNode
	size     int
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

func (t *serverListManager) f00AF8862isGoldLine() bool {
	if t.m138code == 2 || t.m138code == 3 {
		return true
	}
	return false
}

var v012E31B0 int = 0
var v012E3200 int = 5

var v01308D08 battleCoreSpec
var v01308D10 sync.Once

func f0043A2DFgetBattleCoreSpec() *battleCoreSpec {
	v01308D10.Do(func() {

	})
	return &v01308D08
}

type battleCoreSpec struct {
	m00       uintptr // &v0114D394
	m04Enable bool
}

func (b *battleCoreSpec) f0043A33Cenable() bool {
	return b.m04Enable
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
func f00A49798ui() *game {
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

	m6Cstate          int
	m70windowManager  *windowManager01173C54
	m74               uintptr                // state2
	m78               uintptr                // state4
	m7C               *windowManager01173C54 // state5
	m80               *windowgame0116A864
	m84infoTooltip    *windowgameInfoTooltip
	m88caution        *windowgameCaution    // Caution
	m8CmainFrame      *windowgameMainFrame  // MainFrame
	m90               uintptr               // characterFrame
	m94               uintptr               // inventoryFrame
	m98ExtensionBag   uintptr               // ExtensionBagFrame
	m9C               uintptr               // PrivateStoreFrame
	mA0               uintptr               // MixInventoryFrame
	mA4               uintptr               // MixSel
	mA8               uintptr               // MixGoblinSel
	mAC               uintptr               // MixStone
	mB0               uintptr               // MixInventoryFrame
	mB4               uintptr               // GensRanking
	mB8               uintptr               // GuildCreateFrame
	mBC               uintptr               // GuildInfoFrame
	mC0               uintptr               // GuildPosition
	mC4               uintptr               // infoPopup
	mC8               uintptr               // PetInfoPopup
	mCC               uintptr               // ?
	mD0               uintptr               // QuickCommand
	mD4               uintptr               // duel
	mD8               uintptr               // duelwatch
	mDC               uintptr               // trade
	mE0               uintptr               // purchaseinven
	mE4               uintptr               // goldarcher
	mE8               uintptr               // ?
	mEC               uintptr               // npcdialogue
	mF0               uintptr               // NpcShop
	mF4               uintptr               // PetInfo
	mF8               uintptr               // Pentagram
	mFC               uintptr               // ArkaResultInfo
	m100              uintptr               // StorageInven
	m104              uintptr               // StorageInven
	m108              uintptr               // NpcJobChangeQuest
	m10C              uintptr               // NpcQuestProgress
	m110              uintptr               // NpcGateKeeper
	m114              uintptr               // SystemMenu
	m118              uintptr               // Option
	m11C              uintptr               // SelectMenu
	m120              uintptr               // MyQuestInfo
	m124ChatWindow    *windowgameChat       // ChatWindow
	m128              uintptr               // MoveCommand
	m12C              uintptr               // HelpWindow
	m130              uintptr               // EventMapHelper
	m134              uintptr               // Trainer 驯兽师
	m138MapName       uintptr               // MapName
	m13C              uintptr               // Notice
	m140              uintptr               // HeroPosition
	m144              uintptr               // EnterBloodCastle
	m148              uintptr               // EnterEmpireGuardian
	m14C              uintptr               // CommandWindow
	m150              uintptr               // ArkaBattleRegister
	m154              uintptr               // ArkaBattleProgress
	m158              uintptr               // ArkaBattleNotice
	m15C              uintptr               // MonsterInfo
	m160              uintptr               // PlayerInfo
	m164BuffFram      uintptr               // BuffFrame
	m168              uintptr               // EquipDurabilityInfo
	m16C              uintptr               // MacroMain
	m170              uintptr               // MacroSub
	m174Alarm         uintptr               // Alarm
	m178              uintptr               // MatchingSelect
	m17CMatchingGuild uintptr               // MatchingGuild
	m180MatchingParty uintptr               // MatchingParty
	m184partyFrame    *windowgamePartyFrame // PartyFrame
	m188              uintptr               // ItemInfo
	m18C              uintptr               // StrongestMatchMenu
	m190              uintptr               // StrongestMatchRank
	m194              uintptr               // StrongestMatchResult
	m198              uintptr               // SearchPrivateStore
	m19C              uintptr               // MiniGameRummy
	m1A0              uintptr               // EventInventoryFrame
	m1A4              uintptr               // EventInfo
	m1A8              uintptr               // CursedTempleMatchMenu
	m1AC              uintptr               // CursedTempleResult
	m1B0              uintptr               // CursedTempleScore
	m1B4              uintptr               // CursedTempleInfo
	m1B8              uintptr               // ListOfMatches
	m1BC              uintptr               // EventMapProgressInfo
	m1C0              uintptr               // EventMapRoundCount
	m1C4              uintptr               // ProgressFrame
	m1C8Navimap       uintptr               // Navimap
	m1CC              uintptr               // PetInventoryFrame
	m1D0              uintptr               // MuunExchange
	m1D4              uintptr               // OrdealSquareEnterMenu
	m1D8              uintptr               // OrdealSquareScore
	m1DC              uintptr               // OrdealSquareResult
	m1E0              uintptr               // EventMapTutorial
	m1E4              uintptr               // EventInfo
	m1E8              uintptr               // BattleField_BigDialog
	m1EC              uintptr               // RideButton
	m1F0              uintptr               // MonsterCompensation
	m1F4ImageNotice   uintptr               // ImageNotice
	m1F8GremoryCase   uintptr               // GremoryCase 活动奖励背包
	m1FC              uintptr               // RebuyList
	m200              uintptr               // PetFrame
	m204              uintptr               // do5 ebp7C

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

func (t *game) f00A4E46DloadResource(w iwindowgame0117373C, x int, name *xstring, y bool, align int) bool {
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
			wm.f00A47FABappend(w)
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

func (t *game) f00A4A521loadResource() {
	// 0x610局部变量
	// ebp494 := t

	// 0x00A4A545: t.m80
	ebp14 := new(windowgame0116A864) // f00DE852Fnew
	ebp14.f008E20BBconstruct()
	ebp498 := ebp14
	ebp10 := ebp498
	t.m80 = ebp10

	// 0x00A4A58D: infoTooltip
	ebp1C := new(windowgameInfoTooltip)
	ebp1C.f00A3BBB9construct()
	ebp49C := ebp1C
	ebp18 := ebp49C
	t.m84infoTooltip = ebp18

	// 0x00A4A5D8: t.m70windowManager
	t.m70windowManager = nil

	// 0x00A4A5E2: Caution
	ebp24 := new(windowgameCaution)
	ebp24.f00A99F11construct()
	ebp4A0 := ebp24
	ebp20 := ebp4A0
	t.m88caution = ebp20
	var tmpName *xstring
	tmpName.f00BADDD0xstring("Caution")
	t.f00A4E46DloadResource(t.m88caution, 5, tmpName, false, 9)
	// ...
	// 0x00A4C560: PartyFrame
	ebp2B0 := new(windowgamePartyFrame)
	ebp2B0.f00A81B51construct()
	ebp57C := ebp2B0
	ebp2AC := ebp57C
	t.m184partyFrame = ebp2AC
	tmpName.f00BADDD0xstring("PartyFrame")
	t.f00A4E46DloadResource(t.m184partyFrame, 1, tmpName, true, 11)
	// ...
	// 0x00A4D41E: MainFrame
	ebp3F4 := new(windowgameMainFrame)
	ebp3F4.f00AA9021construct()
	ebp5E8 := ebp3F4
	ebp3F0 := ebp5E8
	t.m8CmainFrame = ebp3F0
	tmpName.f00BADDD0xstring("MainFrame")
	t.f00A4E46DloadResource(t.m8CmainFrame, 1, tmpName, true, 2)
	// ...
	// 0x00A4A9BD: ChatWindow
	ebp80 := new(windowgameChat)
	ebp80.f00A9D6F2construct()
	ebp4C0 := ebp80
	ebp7C := ebp4C0
	t.m124ChatWindow = ebp7C
	tmpName.f00BADDD0xstring("ChatWindow")
	t.f00A4E46DloadResource(t.m124ChatWindow, 6, tmpName, true, 10)
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
		// t.m70windowManager = t.m74
		// ebp4hIMC := dll.imm32.ImmGetContext(v01319D6ChWnd)
		// dll.imm32.ImmSetConversionStatus(ebp4hIMC, 0, 0)
		// dll.imm32.ImmReleaseContext(v01319D6ChWnd, ebp4hIMC)
	case 4:
		// t.m70windowManager = t.m78
		// ebp8hIMC := dll.imm32.ImmGetContext(v01319D6ChWnd)
		// dll.imm32.ImmSetConversionStatus(ebp8hIMC, 1, 0)
		// dll.imm32.ImmReleaseContext(v01319D6ChWnd, ebp8hIMC)
	case 5:
		// t.m70windowManager = t.m7C
		// ebpChIMC := dll.imm32.ImmGetContext(v01319D6ChWnd)
		// dll.imm32.ImmSetConversionStatus(ebpChIMC, 1, 0) // IME_CMODE_NATIVE, IME_CMODE_ALPHANUMERIC
		// dll.imm32.ImmReleaseContext(v01319D6ChWnd, ebpChIMC)
	}
	if t.m70windowManager == nil {
		return true
	}
	// t.m70windowManager.f00A47353()
	// t.m70windowManager.f00A47C55()
	return true
}

func (t *game) f00A4DC94fresh(state int) {
	if t.m80 != nil {
		// t.m80.f008E27CFfresh()
	}
	if t.m70windowManager != nil {
		t.m70windowManager.f00A47461fresh()
	}
}

func (t *game) f00A4DF93handleKeyPress(wParam, lParam uintptr) {
	if t.m70windowManager == nil {
		return
	}
	t.m70windowManager.f00A47A25handleKeyPress(wParam)
}

func (t *game) f00A4DFB7handleClick(hWnd win.HWND, msg uint32, wParam, lParam uintptr, unk bool) {
	if t.m70windowManager == nil {
		return
	}
	t.m70windowManager.f00A47A93handleClick(hWnd, msg, wParam, lParam, unk)
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

// --------------------------------------------------------------
var v01319B88 = [100]int{}

// --------------------------------------------------------------
var v01319D20 iwindow
var v01319D24 iedit
var v01319D28 iedit

// iwindow interface
type iwindow interface {
	do2(win.HWND)
	do12initWindow(hWnd win.HWND, width, height int, charNum uint, passwd bool)
}

// iedit interface
type iedit interface {
	do2showWindow(int)
	do12initWindow(hWnd win.HWND, width, height int, charNum uint, passwd bool)
}

// 有虚表的类型尽量使用虚表命名，否则使用构造函数命名
type t0048D468 struct {
	m10 int
	m14 int
	m18 int
	m1C int
}

func (t *t0048D468) f0048D468construct() {
	// f004973DE()
	// f004972E4()
	t.m10 = 0
	t.m14 = 0
	t.m18 = 0
	t.m1C = 0
}

type t0114AFAC struct {
	m00vtabptr []uintptr
	m04        t0048D468
}

func (t *t0114AFAC) f0045D2C3construct() {
	// t.m00vtabptr = v0114AFAC[:]
	t.m04.f0048D468construct()
}

// sizeof=0x6C
type t0114A7F0 struct {
	t0114AFAC
	m30     int
	m34     int
	m38show int
	m3C     int // 5
	m40     int
	// ...
	m68 int
}

func (t *t0114A7F0) f004499B2construct() {
	t.t0114AFAC.f0045D2C3construct()
	// t.m00vtabptr = v0114A7F0[:]
	// f0044987E()
	// f00449A57()
	t.f00449A7A(0, 0)
	// f00449A97()
	// f00449AB4()
	// f00449ADA()
}

func (t *t0114A7F0) f00449A7A(p1, p2 int) {}

func (t *t0114A7F0) f004397CF(p1 int) {}

func (t *t0114A7F0) f0045D373() int {
	return t.m30
}

func (t *t0114A7F0) f0045D391(p1 int) int {
	return t.m3C & p1
}

// sizeof=0xE0
type t0114AE34edit struct {
	t0114A7F0
	m6CwndProc    uintptr
	m70           int
	m74hWndParent win.HWND
	m78hWnd       win.HWND
	m7C           int
	m80           int
	m84           int
	m8C           int
	m90           int
	mA0           int
	mACx          int
	mB0y          int
	mB4vscroll    int
	mDC           bool
}

func f0045D29F(p1, p2, p3, p4 uint8) int { return 0 }

func (t *t0114AE34edit) f004521F5construct() *t0114AE34edit {
	t.t0114A7F0.f004499B2construct()
	// t.m00vtabptr = v0114AE34[:]
	t.m6CwndProc = 0
	t.m74hWndParent = 0
	t.m78hWnd = 0 // 在地址空间的pe seciton里面
	t.m7C = 0
	t.m80 = 0
	t.m84 = 0
	t.m8C = f0045D29F(0xFF, 0xFF, 0xFF, 0xFF)
	// ...
	return t
}

// f00452D78
func (t *t0114AE34edit) do2showWindow(flag int) {
	if t.m78hWnd == 0 {
		return
	}
	t.m38show = flag
	if t.m38show == 3 {
		win.ShowWindow(t.m78hWnd, 0 /*win.SW_HIDE*/)
	} else {
		win.ShowWindow(t.m78hWnd, 5 /*win.SW_SHOW*/)
	}
}

// f004529B4
func (t *t0114AE34edit) do3initDC(width, height int) {

}

func f008AEFAD(p1 int) bool {
	return true //f008AF00Dget().f008AF156(p1)
}

// f004524CC
func f004524CCeditWndProc(hWnd win.HWND, message uint32, wParam, lParam uintptr) uintptr {
	ebp4window := (*t0114AE34edit)(unsafe.Pointer(win.GetWindowLongPtr(hWnd, win.GWL_USERDATA)))
	if ebp4window == nil {
		return 0
	}
	if f008AEFAD(0x25) || f008AEFAD(0x26) || f008AEFAD(0x27) || f008AEFAD(0x28) {
		ebp4window.m70 = 0
	}
	switch message {
	case 0x102: // win.WM_CHAR
		// 0x00452541:
		ebp4window.m70 = 0
		ebp10 := int(wParam)
		switch ebp10 {
		case 0x9:
		case 0xD:
		case 0x1B:
		}
	case 0x10F, 0x282: // WM_IME_COMPOSITION, WM_IME_NOTIFY
		// 0x0045277C:
	default:
		// 0x004527C4:
		// if f005A4BC5(0x74) {
		// 	f008D6008().f004BF36C(1)
		// }
	}
	return win.CallWindowProc(ebp4window.m6CwndProc, hWnd, message, wParam, lParam) // 会经过 dll.user32.xxx 再回调到 f004D5F98mainWndProcOrigin
}

// f00452C00
func (t *t0114AE34edit) do12initWindow(hWnd win.HWND, width, height int, charNum uint, passwd bool) {
	t.m74hWndParent = hWnd
	var ebp4style uint32
	if passwd {
		ebp4style |= 0x20 // ES_PASSWORD/0x20
		t.mA0 = 1
	}
	if t.mB4vscroll == 1 {
		ebp4style |= 0x200044 // WS_VSCROLL/0x200000 | ES_AUTOVSCROLL/0x40 | SS_BLACKRECT/0x4
	} else {
		ebp4style |= 0x80 // ES_AUTOHSCROLL/0x80
	}
	// 0x00452C56: x
	// v01308EC4width = t.m40
	// t.mACx = f00DE76C0roundf(v01308EC4width)
	// 0x00452C81: y
	// v01308EC8height = t.m44
	// t.mB0y = f00DE76C0roundf(v01308EC8height)

	// 0x00452CD9: CreateWindow
	t.m78hWnd = win.CreateWindowEx(
		0,                    // ExStyle
		nil,                  // ClassName, "edit"预定义窗口类, syscall.UTF16FromString(v0114A710)
		nil,                  // WindowName
		ebp4style|0x50000000, // Style WS_VISIBLE/0x10000000 | WS_CHILD/0x40000000
		int32(t.mACx),        // x
		int32(t.mB0y),        // y
		0,                    //f00DE76C0roundf(width/v01308EC4width),   // width
		0,                    //f00DE76C0roundf(height/v01308EC8height), // height
		t.m74hWndParent,      // parent window
		win.HMENU(1),         // menu
		v01319D70hInstance,   // Instance
		nil,
	)

	// 0x00452D0A:
	t.do3initDC(width, height)

	// 0x00452D1B:
	if t.m78hWnd != 0 {
		t.do16sendMessage(charNum)
		t.m6CwndProc = win.SetWindowLongPtr(t.m78hWnd, win.GWL_WNDPROC, syscall.NewCallback(f004524CCeditWndProc))
		win.SetWindowLongPtr(t.m78hWnd, win.GWL_USERDATA, uintptr(unsafe.Pointer(t)))
		win.ShowWindow(t.m78hWnd, 0 /*SW_HIDE*/) // WM_SHOWWINDOW/0x18, WM_WINDOWPOSCHANGING/0x46, WM_WINDOWPOSCHANGED/0x47
	}
	t.mDC = true
}

// f00452DBE
func (t *t0114AE34edit) do13setFocus(sel bool) {
	if t.m78hWnd == 0 {
		return
	}
	if v012E2210 == 1 &&
		// user32.GetFocus() == v01319D6ChWnd &&
		t.f0045D391(2) == 0 &&
		t.f0045D391(1) == 0 { // win10 return 1
		win.SetFocus(t.m78hWnd)
		// f0045203A()
	} else {
		// This function sends a WM_KILLFOCUS message to the window that loses the keyboard focus
		// and a WM_SETFOCUS message to the window that receives the keyboard focus.
		// It also activates either the window that receives the focus
		// or the parent of the window that receives the focus.
		win.SetFocus(t.m78hWnd)
		// win7
		// mainWndProc.WM_KILLFOCUS/8, hEdit, 0
		// mainWndProc.WM_IME_SETCONTEXT/0x281, 0, ISC_SHOWUIALL/0xC000000F
		// editWndProc.WM_IME_SETCONTEXT/0x281, 1, ISC_SHOWUIALL/0xC000000F
		// editWndProc.WM_SETFOCUS/7, hMainWnd, 0
		// editWndProc.WM_IME_NOTIFY/0x282, A, 0
		// editWndProc.WM_IME_NOTIFY/0x282, F, ?
		// editWndProc.WM_IME_NOTIFY/0x282, B, 0
		// mainWndProc.WM_COMMAND/0x111, 0x10000001, hEdit

		// win10
		// mainWndProc.WM_KILLFOCUS/8
		// 然后就没有消息了，不知道从哪里进入的回调f00D86830，然后该回调产生了非法内存访问
		// 而client_vc在win7和win10上是一致的
	}
	// v01308EDC = t.f0045D373()
	if sel {
		// If the start is 0 and the end is -1, all the text in the edit control is selected.
		win.PostMessage(t.m78hWnd, 0xB1 /*EM_SETSEL*/, 0, ^uintptr(0)) // 0, -1
	} else {
		// If the start is -1, any current selection is deselected.
		win.PostMessage(t.m78hWnd, 0xB1 /*EM_SETSEL*/, ^uintptr(0), ^uintptr(0)) // -1, -2
	}
}

// f00452993
func (t *t0114AE34edit) do16sendMessage(charNum uint) {
	win.SendMessage(t.m78hWnd, 0xC5 /*EM_LIMITTEXT*/, uintptr(charNum), 0)
}

type t0114E1C0 struct {
	m00vtabptr []uintptr
	m004user   t0114AE34edit
	m0E4passwd t0114AE34edit
	// m2D0
	// m2F0
	// m2F8
}

func (t *t0114E1C0) f004D9B19construct() *t0114E1C0 {
	// t.m00vtabptr = v0114E1C0[:]

	return t
}

// f00454632
func (t *t0114E1C0) do2(hWnd win.HWND) {}

// sizeof=0x300
type t0114E1ACwindow struct {
	t0114E1C0
}

func (t *t0114E1ACwindow) f004D9AFCconstruct() *t0114E1ACwindow {
	t.t0114E1C0.f004D9B19construct()
	// t.m00vtabptr = v0114E1AC[:]
	return t
}

// f0045515C
func (t *t0114E1ACwindow) do2(hWnd win.HWND) {
	// 0x00455163: init user edit
	t.m004user.do12initWindow(hWnd, 100, 14, 10, false)
	t.m004user.f00449A7A(0x126, 0x12C)
	// t.m004user.do17(0xFF, 0xFF, 0xE6, 0xD2)
	// t.m004user.do18(0x0, 0x2D, 0x2D, 0x2D)
	// t.m004user.f004397CF(5)
	t.m004user.do13setFocus(false)

	// 0x004551EA: init passwd edit
	t.m0E4passwd.do12initWindow(hWnd, 100, 14, 10, false) // false?
	t.m0E4passwd.f00449A7A(0x15E, 0x126)
	// t.m0E4passwd.do17(0xFF, 0xFF, 0xE6, 0xD2)
	// t.m0E4passwd.do18(0x0, 0x2D, 0x2D, 0x2D)
	// t.m0E4passwd.f004397CF(5)
	win.SetFocus(v01319D6ChWnd)
}

// f00BAB280
func (t *t0114E1ACwindow) do12initWindow(hWnd win.HWND, width, height int, charNum uint, passwd bool) {

}

// ----------------------------------------------
var v09D96AAConce sync.Once
var v09D96A84modeManager modeManager

func f00B0EF1EmodeManager() *modeManager {
	v09D96AAConce.Do(func() {
		// v09D96A84modeManager.f00B0EE90construct()
	})
	return &v09D96A84modeManager
}

type DEVMODE struct {
	dmBitsPerPel       uint32 // 0x68
	dmPelsWidth        uint32 // 0x6C
	dmPelsHeight       uint32 // 0x70
	dmDisplayFrequency uint32 // 0x78
}

type DEVMODEx struct {
	devMode DEVMODE
	index   int
	w       int
	h       int
	isWide  bool
}

// size:0xC4
type mode struct {
	m00     [3]*mode
	index   int
	devMode DEVMODE
	mC0     bool
	mC1     bool
}
type modeManager struct {
	m18modes []*mode // tree
	m1Csize  int
}

func (t *modeManager) f00B0F0E0filter(devMode *DEVMODE) bool {
	if devMode.dmDisplayFrequency < 48 {
		return false
	}
	if devMode.dmPelsWidth < 800 || devMode.dmPelsHeight < 600 || devMode.dmBitsPerPel != 32 {
		return false
	}
	for _, node := range t.m18modes {
		m := node.devMode
		if devMode.dmPelsWidth == m.dmPelsWidth && devMode.dmPelsHeight == m.dmPelsHeight {
			// && devMode.dmDisplayFrequency == m.dmDisplayFrequency { // hook discard the equality operator of dmDisplayFrequency, fixed resolution mismatching
			return false
		}
	}
	return true
}

func (t *modeManager) f00B0EF7BenumMode() {
	var ebpB8devModex DEVMODEx
	ebp8index := 0
	ebp4id := 0
	for {
		// if user32.EnumDisplaySettings(0, ebp4id, &ebpB8devModex.devMode) == false {
		// 	break
		// }
		ebpB8devModex.index = ebp8index
		ebpB8devModex.isWide = false
		if t.f00B0F0E0filter(&ebpB8devModex.devMode) {
			switch {
			case ebpB8devModex.devMode.dmPelsWidth/4*3 == ebpB8devModex.devMode.dmPelsHeight:
				ebpB8devModex.w = 4
				ebpB8devModex.h = 3
			case ebpB8devModex.devMode.dmPelsWidth/16*9 == ebpB8devModex.devMode.dmPelsHeight:
				ebpB8devModex.w = 16
				ebpB8devModex.h = 9
				ebpB8devModex.isWide = true
			case ebpB8devModex.devMode.dmPelsWidth/16*10 == ebpB8devModex.devMode.dmPelsHeight:
				ebpB8devModex.w = 16
				ebpB8devModex.h = 10
				ebpB8devModex.isWide = true
			default:
				goto label1
			}
			// ebp168.f00B10FE8(ebp218.f00B0F6BD(&ebp8index, &ebpB8devModex))
			// t.f00B0F385addMode(&ebp224, &ebp168)
			ebp8index++
		}
	label1:
		ebp4id++
	}
}

func (t *modeManager) f00B0F1F6getMode(index int) {

}

// -------------------------------------------
type itemBase struct {
	m22EkindA uint8
	m22FkindB uint8
	m230kind  uint8
}

// size:0x8C
type item struct {
	m00itemBase    *itemBase
	m04code        int
	m08option      int
	m1Edurability  uint8
	m20option      uint8
	m21option      uint8
	m30            [8]uint16
	m40            [8]uint8
	m58is380       bool
	m59sockets     [5]uint8
	m5F            [5]uint8
	m6Csocketbonus int // harmony/socket/pentagram/muun bonus
	m80excel       uint8
	m81            [6]uint8
	m87            uint8
	m88            uint8
}

func (t *item) f00A2742Einit() {}

func (t *item) f00A270D6construct() {
	t.f00A2742Einit()
}

func (t *item) f00A276A1getItemBase() *itemBase {
	if t.m00itemBase == nil {
		// t.f00A2762E()
		func() {
			// t.m00 = f00A29EC7itemBaseTable().f00A2AD34getItemBase(t.m04code)
		}()
	}
	return t.m00itemBase
}

func (t *item) f00524B13kindA(kindA int) bool {
	if t.m04code < 0 || t.m04code > 0x2000 {
		return false
	}
	return int(t.f00A276A1getItemBase().m22EkindA) == kindA
}

func (t *item) f00524B55kindB(kindB int) bool {
	if t.m04code < 0 || t.m04code > 0x2000 {
		return false
	}
	return int(t.f00A276A1getItemBase().m22FkindB) == kindB
}

func (t *item) f004FE875kind(kind int) bool {
	if t.m04code < 0 || t.m04code > 0x2000 {
		return false
	}
	return int(t.f00A276A1getItemBase().m230kind) == kind
}

func (t *item) f00A29A7AhasExcel() bool {
	if t.f004FE875kind(1) || // itemTypeRegular
		t.f004FE875kind(3) || // itemType380
		t.f004FE875kind(6) || // itemTypeArchangel
		t.f004FE875kind(7) { // itemTypeChaos
		return true
	}
	return false
}

// size:0x388
type item2 struct {
	item
	m8C struct {
		m04 uint8
		m14 int
		m18 int
		m1C int
		m40 int
		m44 int
	}
	m334 uint8
	m384 int
}

var v086D09C8items [1000]mapItem

// sizeof=0x38C
type mapItem struct {
	number int
	item2
}

func f005A0292code(attr []uint8) int {
	return (int(attr[3]) & 0x80 << 1) | int(attr[0]) | (int(attr[5]) & 0xF0 << 5)
}

func f00677421mapItemAdd(mapItem *mapItem, attr []uint8, ord []float32, msb int) {
	ebp8code := f005A0292code(attr)
	ebp4item := &mapItem.item2
	ebp4item.f00A2742Einit()
	ebp4item.m04code = ebp8code
	if ebp4item.f00A276A1getItemBase() == nil {
		return
	}
	if ebp8code == 0x1C0F { // (14, 15) zen
		ebp4item.m08option = int(attr[1])<<16 | int(attr[2]<<8) | int(attr[4])
		ebp4item.m1Edurability = 0
		ebp4item.m20option = 0
		if msb == 1 {
			// f007DAFE0(0x1F, 0, 0)
		}
	} else {
		ebp4item.m08option = int(attr[1])
		ebp4item.m1Edurability = attr[2]
		ebp4item.m20option = attr[3] & 0x7F
		ebp4item.m21option = attr[4]
		ebp4item.m58is380 = false
		if attr[5]&0x08>>3 == 1 {
			ebp4item.m58is380 = true
		}

		// pentagram
		if ebp4item.f00524B13kindA(8) == true { // itemKindAPentagram
			ebp4item.m6Csocketbonus = int(attr[6]) & 0xF
			// ebp4item.mA0 = ebp4item.m6Csocketbonus
			// ebp18 := &ebp4item.m8C
			// ebp18.m14 = ebp4item.m6Csocketbonus
			// if ebp18.m14 >= 0 && ebp18.m14 < 6 {
			// 	ebp1C := f00A18027().f00A19CBF(ebp8code, 0)
			// 	if ebp1C != nil {
			// 		ebp18.m18 = ebp1C.m10
			// 		ebp18.m1C = ebp1C.m14[ebp18.m14]
			// 	}
			// }
			if ebp4item.f00524B55kindB(44) == true { // itemKindBPentagramJewel 艾尔特
				for i := 0; i < 5; i++ {
					ebp4item.m59sockets[i] = attr[7+i]
				}
			}
		}

		// excellent option
		if ebp4item.f00A29A7AhasExcel() == true {
			if attr[3]&0x3F != 0 {
				ebp4item.m80excel = 1
			}
			for i := 0; i < 3; i++ {
				if attr[7+i] != 0xFF {
					ebp4item.m80excel = 1
					break
				}
			}
		}

		if msb != 0 {
			if ebp8code == 0x1C0D || // (14, 13) 祝福宝石
				ebp8code == 0x1C0E || // (14, 14) 灵魂宝石
				ebp8code == 0x1C10 || // (14, 16) 生命宝石
				ebp8code == 0x180F || // (12, 15) 玛雅宝石
				ebp8code == 0x1C16 || // (14, 22) 创造宝石
				ebp8code == 0x1C1F { // (14, 31) 守护宝石
				// f007DAFE0(0x45, &ebp4item.m8C, 0)
			} else if ebp8code == 0x1C29 { // (14, 41) 再生原石
				// f007DAFE0(0x2D3, &ebp4item.m8C, 0)
			} else if ebp8code == 0x1CFE || // (14, 254) 实习骑士团的礼物
				(ebp8code >= 0x1D00 && ebp8code <= 0x1D06) { // (14, 256) ~ (14, 262) 骑士团的礼物 恶魔箱子 卫兵箱子
				// f007DAFE0(0x5A, &ebp4item.m8C, 0)
			} else {
				// f007DAFE0(0x1E, &ebp4item.m8C, 0)
			}
		}
	}
	// 0x6776FA:
	ebpC := &ebp4item.m8C
	ebpC.m04 = 1
	ebpC.m40 = ebp8code + 0x4F8 // 1272
	ebpC.m44 = 1

	// if ebp8code == 0x1C0B { // (14, 11) 幸运宝箱
	// 	switch ebp4item.m08option >> 3 {
	// 	case 1:
	// 		ebpC.m40 = 0x258E
	// 	case 2:
	// 		ebpC.m40 = 0x258F
	// 	case 3:
	// 		ebpC.m40 = 0x2590
	// 	case 5:
	// 		ebpC.m40 = 0x2592
	// 	case 6:
	// 		ebpC.m40 = 0x2593
	// 	case 8, 9, 10, 11, 12:
	// 		ebpC.m40 = 0x2593
	// 	case 13, 14:
	// 		ebpC.m40 = 0x259B
	// 	}

	// } else if ebp8code >= 0x1C2E && ebp8code <= 0x1C30 {

	// }
	switch {
	case ebp8code == 0x1C0B: // (14, 11) 幸运宝箱
		switch ebp4item.m08option >> 3 {
		case 1:
			ebpC.m40 = 0x258E
		case 2:
			ebpC.m40 = 0x258F
		case 3:
			ebpC.m40 = 0x2590
		case 5:
			ebpC.m40 = 0x2592
		case 6:
			ebpC.m40 = 0x2593
		case 8, 9, 10, 11, 12:
			ebpC.m40 = 0x2593
		case 13, 14:
			ebpC.m40 = 0x259B
		}
	case ebp8code >= 0x1C2E && ebp8code <= 0x1C30: // (14, 46/47/48) 南瓜灯祝福/愤怒/呐喊
		ebpC.m40 = ebp8code + 0x4F8
	case ebp8code >= 0x1C20 && ebp8code <= 0x1C22: // (14, 32/33/34) 粉红色/红色/蓝色宝箱
		if ebp4item.m08option>>3 == 1 {
			ebpC.m40 = ebp8code - 0x1C20 + 0x259F
		}
	case ebp8code == 0x1C15:
		ebp58 := ebp4item.m08option >> 3
		if ebp58 > 0 && ebp58 <= 2 {
			ebpC.m40 = 0x2595
		}
	case ebp8code == 0x1A13: // (13, 19) 大天使之武器
	case ebp8code == 0x1A17: // (13, 23) 风之指环
	case ebp8code == 0x1A18: // (13, 24) 魔之指环
	case ebp8code == 0x1A14: // (13, 20) 魔法戒指
	case ebp8code == 0x1C09: // (14, 9) 酒
	case ebp8code == 0x1A0E: // (13, 14) 洛克之羽
	case ebp8code == 0x1A0B: // (13, 11) 复活石
	case ebp8code == 0x1C29: // (14, 41) 再生原石
	case ebp8code == 0x1C2A: // (14, 42) 再生宝石
	case ebp8code == 0x2122:
	case ebp8code == 0x1C2B: // (14, 43) 处级进化宝石
	case ebp8code == 0x2123:
	case ebp8code == 0x2124:
	}
	// f00674D7A(ebpC)
	// f007EA8A8(...)
	// f00653E1E()
	// f00675564(ebpC)
	// f0067328Cdraw(mapItem)
}

func f006B8509() *t0114EF00 {
	return v0131A34C
}

var v0131A34C *t0114EF00

type tPremium struct {
}

func (t *tPremium) f004EC50D() {

}

// func (t *tPremium) f004EBB74(unk *int) {
// 	var ebp8 struct{}
// 	t.f004EC50D(&ebp8, unk)
// 	if ebp8.f004EC749(t.f004EBC7C(&ebp10)) != 0 ||
// 		t.m08.f00436FF3(unk, f004EBE4E(ebp8.f004EC78F())) != 0 {
// 		ebp1C := 0
// 		ebp24 := ebp8
// 		ebp20 := ebp4
// 		t.f004EBEBA(&ebp2C, ebp24, ebp20, ebp18.f004EC79E(unk, &ebp1C))
// 	}
// 	return ebp8.f004EC724() + 4
// }

type t0114EF00 struct {
	m00vtabptr []uintptr
	m04        tPremium
}

func (t *t0114EF00) f004EBD58construct() *t0114EF00 {
	// t.m00vtabptr = v0114EF00[:]
	if v0131A34C == nil {
		v0131A34C = t
	}
	return nil
}

func (t *t0114EF00) f004EBAE3(unk int) int {
	// p := t.m04.f004EBB74(&unk)
	// if *p == 0 {
	// 	return 0
	// }
	// p = t.m04.f004EBB74(&unk)
	// return p.m08
	return 0
}

// sizeof=0x24
type t0114EEF8 struct {
	t0114EF00
}

func (t *t0114EEF8) f004EB62Bconstruct() *t0114EEF8 {
	t.t0114EF00.f004EBD58construct()
	// t.m00vtabptr = v0114EEF8[:]
	return t
}

// map manager
func f0050E6D7() *t01351D58 {
	v01351D74once.Do(func() {
		v01351D58.f0050E734init()
		// f00DE8BF6atexit(f01148288)
	})
	return &v01351D58
}

var v01351D58 t01351D58
var v01351D74once sync.Once

type t01351D58 struct {
}

func (t *t01351D58) f0050E734init() {
	// t.f0050F974()
	// t.f0050E819()
}

func (t *t01351D58) f0050E893isIllusionTemple(mapnum int) bool {
	if (mapnum >= 45 && mapnum <= 50) ||
		(mapnum > 97 && mapnum <= 99) {
		return true
	}
	return false
}

func (t *t01351D58) f0050F876(unk int) bool {
	return false
}

//
func f0090E94C() *t09D8CFD0 {
	v09D8D080once.Do(func() {
		v09D8CFD0.f0090B0FEinit()
		// f00DE8BF6atexit(f0114865A)
	})
	return &v09D8CFD0
}

var v09D8D080once sync.Once
var v09D8CFD0 t09D8CFD0

type t09D8CFD0 struct {
}

func (t *t09D8CFD0) f0090B0FEinit() {

}
