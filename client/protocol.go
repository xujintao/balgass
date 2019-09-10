package main

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

// t08C8D014
var v08C8D014encrypt packetEncrypt

type packetEncrypt struct{}

func (t *packetEncrypt) f00B98FC0(conn *conn, bufdst []uint8, bufsrc []uint8, len int) int { return 0 }

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
				ebp182C = v08C8D014encrypt.f00B98FC0(conn, ebp181Cbuf[2:], ebp18buf[2:], ebp10len-2)

			// 0AF964D8
			default: // 0xC4
				// ebp1834 := ebp18buf
				ebp10len = int(ebp18buf[1]<<8 + ebp18buf[2])
				ebp182C = v08C8D014encrypt.f00B98FC0(conn, ebp181Cbuf[2:], ebp18buf[3:], ebp10len-3)
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
