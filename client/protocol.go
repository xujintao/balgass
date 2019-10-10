package main

import (
	"encoding/binary"
	"unsafe"
)

type command struct {
	code     int
	handle   func(buf []uint8, len int) bool
	comamnds []*command
}

var cmds = []*command{
	{0x00, handle00, nil},
	{0x01, handle01, nil},
	{0x02, handle02, nil},
	{0x0D, f007087BFhandle0D, nil},
	{0xF1, nil, []*command{
		{0x00, handleF100versionmatch, nil},
		{0x04, handleF104, nil},
	}},
	{0xFA, nil, nil},
	{0xF4, nil, []*command{
		{0x03, handleF403, nil},
		{0x05, handleF405, nil},
		{0x06, handleF406, nil},
	}},
}

func handlecmd(code int, subcode int, buf []uint8, len int) bool {
	for _, cmd := range cmds {
		if code == cmd.code {
			if cmd.handle != nil {
				return cmd.handle(buf, len)
			}
			for _, subcmd := range cmd.comamnds {
				if subcode == subcmd.code {
					return subcmd.handle(buf, len)
				}
			}
		}
	}
	return true
}

func handle00(buf []uint8, len int) bool { return true }
func handle01(buf []uint8, len int) bool { return true }
func handle02(buf []uint8, len int) bool { return true }
func f007087BFhandle0D(buf []uint8, len int) bool {
	switch buf[3] {
	case 0:
	case 1:
		if v012E2340 == 4 {
			ebp14 := f004A7D34()
			ebp14.f004A9F3B(buf[13:]) // "you are logged in as Free Player"
		}
	}
	return true
}
func handleF100versionmatch(buf []uint8, len int) bool { return true }
func handleF104(buf []uint8, len int) bool             { return true }
func handleFA(buf []uint8, len int) bool               { return true }
func handleF403(buf []uint8, len int) bool             { return true }
func handleF405(buf []uint8, len int) bool             { return true }
func handleF406(buf []uint8, len int) bool             { return true }

// cmd hook
var cmdhooks = []*command{
	{0xF1, nil, []*command{
		{0x00, joinResultSend, nil},
	}},
	{0xF3, nil, []*command{
		{0x03, charMapJoinResult, nil},
		{0x04, charRegen, nil},
		{0x05, levelUpSend, nil},
		{0x06, levelUpPointAdd, nil},
		{0x51, masterLevelUpSend, nil},
	}},
	{0xBF, nil, []*command{
		{0x51, ansMuBotRecvUse, nil},
	}},
	{0xF4, nil, []*command{
		{0x03, reconnect, nil},
	}},
	{0xCA, friendRoomCreate, nil},
	{0xD4, AHCheckGetTickHook, nil},
	{0xFA, nil, []*command{
		{0x0D, offTradeReq, nil},
		{0x10, customPost, nil},
		{0x11, customPost, nil},
		{0x12, setChatColors, nil},
		{0x90, IGCStatAdd, nil},
		{0xA0, enableSiegeSkills, nil},
		{0xA2, setAgilityFix, nil},
		{0xA3, setCashItemMoveEnable, nil},
		{0xA4, luckyItemFix, nil}, // 区分C1和C2
		{0xA5, disableReconnect, nil},
		{0xA6, handleExitProcess, nil},
		{0xA7, disableReconnectSystem, nil},
		{0xA8, alterPShopVault, nil},
		{0xB0, dropSellMod, nil},
		{0xF8, setBattleZoneData, nil},
	}},
	{0x11, attackResult, nil},
	{0x17, diePlayerSend, nil},
	{0x26, refillSend, nil},
	{0x27, manaSend, nil},
}

func handlecmdhook(code int, subcode int, buf []uint8, len int) bool {
	for _, cmd := range cmdhooks {
		if code == cmd.code {
			if cmd.handle != nil {
				return cmd.handle(buf, len)
			}
			for _, subcmd := range cmd.comamnds {
				if subcode == subcmd.code {
					return subcmd.handle(buf, len)
				}
			}
		}
	}
	return true
}

func joinResultSend(buf []uint8, len int) bool         { return true }
func charMapJoinResult(buf []uint8, len int) bool      { return true }
func charRegen(buf []uint8, len int) bool              { return true }
func levelUpSend(buf []uint8, len int) bool            { return true }
func levelUpPointAdd(buf []uint8, len int) bool        { return true }
func masterLevelUpSend(buf []uint8, len int) bool      { return true }
func ansMuBotRecvUse(buf []uint8, len int) bool        { return true }
func reconnect(buf []uint8, len int) bool              { return true }
func friendRoomCreate(buf []uint8, len int) bool       { return true }
func AHCheckGetTickHook(buf []uint8, len int) bool     { return true }
func offTradeReq(buf []uint8, len int) bool            { return true }
func customPost(buf []uint8, len int) bool             { return true }
func setChatColors(buf []uint8, len int) bool          { return true }
func IGCStatAdd(buf []uint8, len int) bool             { return true }
func enableSiegeSkills(buf []uint8, len int) bool      { return true }
func setAgilityFix(buf []uint8, len int) bool          { return true }
func setCashItemMoveEnable(buf []uint8, len int) bool  { return true }
func luckyItemFix(buf []uint8, len int) bool           { return true }
func disableReconnect(buf []uint8, len int) bool       { return true }
func handleExitProcess(buf []uint8, len int) bool      { return true }
func disableReconnectSystem(buf []uint8, len int) bool { return true }
func alterPShopVault(buf []uint8, len int) bool        { return true }
func dropSellMod(buf []uint8, len int) bool            { return true }
func setBattleZoneData(buf []uint8, len int) bool      { return true }
func attackResult(buf []uint8, len int) bool           { return true }
func diePlayerSend(buf []uint8, len int) bool          { return true }
func refillSend(buf []uint8, len int) bool             { return true }
func manaSend(buf []uint8, len int) bool               { return true }

type t11 struct {
	f00 *t11
	f04 *t11
	f08 *t11
	f10 uint32
	f0C *conn
	f15 uint8
}

var v08C8D014encrypt packetEncrypt // {f00:0x0ECA3300, f18:0x0ECA82E0}
var v08C88FB4decrypt packetEncrypt

type packetEncrypt struct {
	f00 *t12
	f18 *t11
}
type t12 struct {
	encrypt *packetEncrypt
	f04     *t11
}

// f00B98ED0encrypt 加密报文 被花了
func (t *packetEncrypt) f00B98ED0encrypt(conn *conn, bufdst []uint8, bufsrc []uint8, lensrc int) int {
	// // 0x0AC34017
	// // 8字节局部变量
	// var ebp8 t12
	// // 0x0A327C79 0x09FE6629 0x0A561513 0x0AF9A9FC 0x0A059A62 0x0A88BAD0 0x09FD2349 0x0AD981DC 0x0A561A1B
	// // push ebx
	// // push ebp
	// // push esi
	// // push edi
	// // push ecx
	// // pop edi // edi = ecx
	// // push &conn
	// // push &ebp8
	// // push edi
	// // pop ecx // ecx = edi
	// // 0x00B98EE5
	// t.f00B98D90(&ebp8, &conn) // ebp8:{0x08C8D014, 0}
	// // 0x0AFD452E
	// esi := ebp8.encrypt
	// ebp := t.f18
	// eax := t.f00
	// if esi != nil {
	// 	// 0x0A8FA15D
	// 	if esi == eax {
	// 		// 0x00B98F00
	// 	}
	// }
	// // 0x00B98EFB
	// // f00DE84C8()
	// // 0x0A9D324B
	// ebx := ebp8.f04
	// if ebx != ebp {
	// 	// 0x00B98F15 0x00B98ED5 0x0A43C4FF
	// 	if esi == nil {
	// 		// 0x00B98F19
	// 		// f00DE84C8()
	// 	} else {
	// 		// 0x00B98F8A 0x00B98F0C 0x0AD81E47
	// 		esi := esi.f00
	// 		// 0x00B98F1E
	// 	}
	// 	// 0x0A83E939 0x0A748151 0x0A4E9023 0x09E30AF0 0x0A7ABDD1
	// 	if ebx == esi.f18 {
	// 		// 0x00B98F23
	// 		// f00DE84C8()
	// 		// 0x0A5FAB24
	// 	} else {
	// 		// 0x00B98F28 0x0A5FAB24
	// 	}
	// 	// 0x0A5FAB24
	// 	edx := lensrc
	// 	esi := ebx.f10
	// 	eax := bufsrc
	// 	ebp := bufdst
	// 	// push edx
	// 	// push eax
	// 	// push ebp
	// 	// ebx = &esi.f90
	// 	// ecx = ebx
	// 	// 0x00B98F42
	// 	// ebp8 := xx.f00BA1B70(ebp, eax, edx) // (bufdst, bufsrc, lensrc)
	// 	// 0x09EAF588 0x0AF7C87D 0x0A38CE99 0x0A05D16B 0x0AF7C2BF 0x0A4E15DD 0x09FE8547
	// 	if ebp == nil {
	// 		// 0x00B9BFA7
	// 	} else {
	// 		// 0x09FBA820
	// 	}
	// }
	// // 0x0A603AE1
	// // pop edi
	// // pop esi
	return 0
}

// f00B98FC0decrypt 解密报文 被花了
func (t *packetEncrypt) f00B98FC0decrypt(conn *conn, bufdst []uint8, bufsrc []uint8, lensrc int) int {
	// 8字节局部变量
	var ebp8 t12

	// 0x09FE83A3
	if conn == nil {
		// 0x0A8FB4A8
	}
	// 0x00B98FD7 0x0A05EBFC 0x0A8881B6 0x0A5FC1CD 0x0AD2B88E
	// push esi
	// push edi
	// push &conn
	// push &ebp8
	// 0x00B98FE5
	t.f00B98D90(&ebp8, &conn)
	// 0x09FC76D8
	return 0
}

func (t *packetEncrypt) f00B98D90(unk1 *t12, connp **conn) {
	// // 0x10字节局部变量
	// ebp := connp
	// esi := t
	// ecx := esi.f18
	// eax := ecx.f04
	// edi := ecx
	// for eax.f15 == 0 {
	// 	if eax.f0C >= (*connp) {
	// 		edi = eax
	// 		eax = eax.f00
	// 	} else {
	// 		eax = eax.f08
	// 	}
	// }
	// eax = esi
	// ebx = eax.f18
	// v2 := edi
	// v1 := eax
	// f00DE84C8()
	// if edi != ebx {
	// 	if (*connp) >= edi.f0C {
	// 		ecx = v2
	// 	}
	// } else {
	// 	v4 = esi.f18
	// 	v3 = esi.f00
	// 	ecx = v3
	// }
	// edx = ecx.f00
	// ecx = ecx.f04

	// eax = unk1
	// eax.f00 = edx
	// eax.f04 = ecx
	*unk1 = *t.f00
}

func f00735DC7(uk1 int, code int, buf []uint8, len int, enc int) {

}

func f0075C3B2(code int, buf []uint8, len int, x int) {
	ebp5D8 := code
	if ebp5D8 > 0xFD {
	}

	// 寻找处理器
	// v0075FCB2cmd[4*v0075FF6Acmd[ebp5D8]]

	// 0x00 处理器
	// _006FA9EB
	func(buf []uint8) {
		// 业务逻辑很复杂
		// send c1 04 f4 06
	}(buf)

	// 0xF1 处理器
	ebp4 := buf
	ebp5DC := ebp4[3]
	if ebp5DC > 4 {
	}
	switch ebp5DC {
	case 4:
	case 0: // version match, buf: 01 2E 87 "10525"
		// _006C75A7 被花到 0x0A4E7A49
		func(buf []uint8) {
			ebp4 := buf
			ebp10 := buf[4]
			if ebp10 == 1 {
				ebp8 := &v0130F728
				ebp8.f004A9123(&ebp8.f4880)
			}
			v08C88E0C = int(ebp4[5]<<8 + ebp4[6])
			v08C88E08 = 2
			var ebpC uint8
			if ebpC >= 5 {
			}
			for {
				if v012E4018[ebpC]-ebpC-1 != ebp4[7+ebpC] { // 改为jmp
					break
				}
				ebpC++
				if ebpC >= 5 {
					break
				}
			}

			// _004A9146
			func() {}()

			// _004A9EEB
			func() {}()
			v01319E08log.f00B38AE4printf("Version dismatch - Join server.\r\n")
			// _0051FB61
			func() {}()
		}(buf)
	}

	// 0xFA 处理器(没有)

	// 0xF4 处理器
	if buf[0] == 0xC1 {
	}
	ebp1C := buf
	ebp14 := ebp1C[4]
	ebp5EC := ebp14
	switch ebp5EC {
	case 3:
	case 5:
	case 6:
		// _006BFA3E
		func(buf []uint8) {
			ebp4 := 6
			// ebp9 := buf[ebp4]
			ebp4++
			if ebp14 >= buf[5]<<8|1 {
			}
			// ebp18 := buf[ebp4:]
		}(buf)
	}
}

// hook到 igc.dll了
func f00760292parse(conn *conn, unk1, unk2 int) {
	// SEH
	f00DE8A70chkstk() // 0x5F5C

	for {
		// 0AD3CD7D
		ebp18buf := conn.f006BDC33() // 从连接中获取buf
		if ebp18buf == nil {
			break
		}

		var ebp14code int
		var ebp10len int
		var ebp181Cbuf [7024]uint8
		var ebp1820enc int

		switch ebp18buf[0] {
		case 0xC1:
			// ebp1824 := ebp18buf
			ebp14code = int(ebp18buf[2])
			ebp10len = int(ebp18buf[1])
		case 0xC2:
			// ebp1828 := ebp18buf
			ebp14code = int(ebp18buf[3])
			ebp10len = int(ebp18buf[1]<<8 + ebp18buf[2])

		// 0A443ED9 0AD2AFB1
		case 0xC3, 0xC4:
			var ebp182C int

			switch ebp18buf[0] {
			// 0A0595EE
			case 0xC3:
				// ebp1830 := ebp18buf
				ebp10len = int(ebp18buf[1])
				ebp182C = v08C8D014encrypt.f00B98FC0decrypt(conn, ebp181Cbuf[2:], ebp18buf[2:], ebp10len-2)

			// 0AF964D8
			default: // 0xC4
				// ebp1834 := ebp18buf
				ebp10len = int(ebp18buf[1]<<8 + ebp18buf[2])
				ebp182C = v08C8D014encrypt.f00B98FC0decrypt(conn, ebp181Cbuf[2:], ebp18buf[3:], ebp10len-3)
			}
			ebp1820enc = 1

			// 0A049505
			if ebp182C < 0 {
				// var ebp2CB4 pb
				// ebp2CB4.f00439178init()
				// ebp2D6E := (uint8)0xF1
				// ebp2D6D := (uint8)0xC1
				// (uint16)&ebp2CB4[0] = 0
				// ebp5F60 := ebp2D6D
				// if ebp5F60 == 0xC1
			}

			// 0A83F8C1
			if v08C88F61 == ebp181Cbuf[2] {
				v08C88F61++
			} else {
				ebp1820enc = 0
				v08C88F61 = ebp181Cbuf[2]
				v01319E08log.f00B38AE4printf("Decrypt error : g_byPacketSerialRecv(0x%02X), byDec(0x%02X)\r\n", v08C88F61, ebp181Cbuf[2])
				v01319E08log.f00B38AE4printf("Dump : \r\n")
				v01319E08log.f00B38B43(ebp18buf, ebp10len)
				v01319E08log.f00B38AE4printf("\r\n")
				v01319E08log.f00B38B43(ebp181Cbuf[2:], ebp182C)
				// ...
			}

			switch ebp18buf[0] {
			// 0A43FECF
			case 0xC3:
				// ebp2CB8 := ebp181Cbuf[1:]
				ebp181Cbuf[1] = 0xC1
				ebp181Cbuf[2] = uint8(ebp182C)
				ebp14code = int(ebp181Cbuf[3])
				ebp18buf = ebp181Cbuf[1:]

			// 0A4417A0
			default:
				// ebp2CBC := ebp181Cbuf[:]
				ebp181Cbuf[0] = 0xC2
				ebp181Cbuf[1] = uint8(ebp182C >> 8)
				ebp181Cbuf[2] = uint8(ebp182C & 0xFF)
				ebp14code = int(ebp181Cbuf[3])
				ebp18buf = ebp181Cbuf[:]
			}
			ebp10len = ebp182C
		}

		// 0ABE2DAF
		v086A3C28 += ebp10len // 累计接收字节数

		//
		if unk1 == 1 {
			f00735DC7(unk2, ebp14code, ebp18buf, ebp10len, ebp1820enc)
		} else {
			f0075C3B2(ebp14code, ebp18buf, ebp10len, ebp1820enc)
		}
	}
}

// pb
type pb struct {
	// len   uint16
	// buf   [5210]uint8 // 0x145A
	buf   [5212]uint8 // 0x145C
	f145C t1
}

var v012E4034conn *conn = &v08C88FF0

func (t *pb) f00439178init() {
	// SEH
	t.f145C.f007BF63B()
	t.f145C.f007BF68F()
	t.f0043921B()
}
func (t *pb) f004391CF() {}
func (t *pb) f0043921B() {
	// t.len = 0
	*(*uint16)(unsafe.Pointer(&t.buf[0])) = 0
}
func (t *pb) f0043922CwritePrefix(flag uint8, code uint8) {
	t.f0043921B()

	// 写flag
	var buf [1]uint8
	buf[0] = flag
	t.f00439298writeBuf(buf[:], 1, false)

	// 写len，没有必要
	switch flag {
	case 0xC1:
		t.f00439298writeBuf(t.buf[:], 1, false)
	case 0xC2:
		t.f00439298writeBuf(t.buf[:], 2, false)
	}

	// 写code
	buf[0] = code
	t.f00439298writeBuf(buf[:], 1, false)
}

func (t *pb) f00439298writeBuf(buf []uint8, len int, needEnc bool) {
	tlen := (*uint16)(unsafe.Pointer(&t.buf[0]))
	if int(*tlen)+len > 0x145A {
		return
	}
	f00DE7C90memcpy(t.buf[*tlen+2:], buf, len)
	if needEnc {
		t.f0043930Cenc(int(*tlen), int(*tlen)+len, 1) // (3,4,1), (4,8,1)
	}
	*tlen += uint16(len)
}

func (t *pb) f0043930Cenc(begin, end, interval int) {
	var buf = [32]uint8{
		0xAB, 0x11, 0xCD, 0xFE,
		0x18, 0x23, 0xC5, 0xA3,
		0xCA, 0x33, 0xC1, 0xCC,
		0x66, 0x67, 0x21, 0xF3,
		0x32, 0x12, 0x15, 0x35,
		0x29, 0xFF, 0xFE, 0x1D,
		0x44, 0xEF, 0xCD, 0x41,
		0x26, 0x3C, 0x4E, 0x4D,
	}
	ebp24 := begin
	for {
		if ebp24 == end {
			break
		}
		t.buf[ebp24+2] ^= (t.buf[ebp24+1] ^ buf[ebp24]%32) // t.buf[5] ^= t.buf[4] ^ buf[3%32]
		ebp24 += interval
	}
}

// s9 00439378
func (t *pb) f004393EAsend(needEnc, isC2 bool) {
	// t.f00439612() 写len
	func() {}()

	t.f0043968F() // hook到 igc.dll了，直接ret

	// f00439420
	// hook到 IGC.dll
	func(buf []uint8, len int, needEnc, isC2 bool) {
		f00DE8A70chkstk() // 0x3124

		if !needEnc {
			v012E4034conn.f004397E3write(buf, len)
			return
		}

		var ebp311C int
		var ebp3118 [1000]uint8 // 原始缓存副本
		var ebp1914 int
		var ebp1910 [1000]uint8 // C4编码缓存
		var ebp108 [1000]uint8  // C3编码缓存

		// 编码数据
		f00DE7C90memcpy(ebp3118[:], buf, len)
		ebp3118[len] = uint8(f00DE8AADrand() % 256) //
		if isC2 {
			ebp3118[0] = 0xC2
		}
		if ebp3118[0] != 0xC1 {
			ebp311C = 1
		}
		ebp311C += 2
		ebp3118[ebp311C-1] = v08C88F60 // 流水号？
		v08C88F60++
		ebp311C--

		ebp1914 = v08C88FB4decrypt.f00B98ED0encrypt(v012E4034conn, nil, ebp3118[ebp311C:], len-ebp311C) // 得到len为0x11

		// len<256且C1 编码方案
		if ebp1914 < 0x100 && !isC2 {
			ebp3120 := ebp1914 + 2 // 0x13，超过255怎么办？
			ebp108[0] = 0xC3
			ebp108[1] = uint8(ebp3120)
			v08C88FB4decrypt.f00B98ED0encrypt(v012E4034conn, ebp108[2:], ebp3118[ebp311C:], len-ebp311C) // 编码数据
			v012E4034conn.f004397E3write(ebp108[:], ebp3120)
			return
		}
		// len>=256或者C2 编码方案
		ebp3124 := ebp1914 + 3
		ebp1910[0] = 0xC4
		binary.BigEndian.PutUint16(ebp1910[1:], uint16(ebp3124))
		v08C88FB4decrypt.f00B98ED0encrypt(v012E4034conn, ebp1910[3:], ebp3118[ebp311C:], len-ebp311C)
		v012E4034conn.f004397E3write(ebp1910[:], ebp3124)
	}(t.buf[2:], int(*(*uint16)(unsafe.Pointer(&t.buf[0]))), needEnc, isC2)
}

func (t *pb) f0043968F() {
	var ebp4len uint16
	var ebp8buf []uint8
	var ebp9code uint8
	if t.buf[2] == 0xC1 {
		ebp4len = uint16(t.buf[3])
		ebp9code = t.buf[4]
		ebp8buf = t.buf[5:]
		ebp4len -= 3

	}
	if t.buf[2] == 0xC2 {
		ebp4len = uint16(t.buf[3]<<8 + t.buf[4])
		ebp9code = t.buf[5]
		ebp8buf = t.buf[6:]
		ebp4len -= 4
	}

	t.f145C.f007C2763(ebp9code, ebp8buf, ebp4len)
}

func (t *pb) f0043974FwriteZero(x int) {
	var buf [5216]uint8
	f00DE8A70chkstk()

	f00DE8100memset(buf[:], 0, x)
	t.f00439298writeBuf(buf[:], x, true)
}

func (t *pb) f004397B1writeUint8(data uint8) *pb {
	var buf [1]uint8
	buf[0] = data
	t.f00439298writeBuf(buf[:], 1, true)
	return t
}

func (t *pb) f0043EDF5writeUint32(data uint32) {
	var buf [4]uint8
	binary.LittleEndian.PutUint32(buf[:], data)
	t.f00439298writeBuf(buf[:], 4, true)
}
func (t *pb) f004C65EFwriteUint16(data uint16) *pb {
	var buf [2]uint8
	binary.LittleEndian.PutUint16(buf[:], data)
	t.f00439298writeBuf(buf[:], 2, true)
	return t
}
