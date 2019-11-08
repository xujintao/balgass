package main

type command struct {
	code     uint8
	handle   func(buf []uint8)
	comamnds []*command
}

// cmd
func f0075C3B2handlecmd(code uint8, buf []uint8, len int, enc bool) {
	for _, cmd := range cmds {
		if code == cmd.code {
			if cmd.handle != nil {
				cmd.handle(buf)
				return
			}
			flag := buf[0]
			var subcode uint8
			switch flag {
			case 0xC1:
				subcode = buf[3]
			case 0xC2:
				subcode = buf[4]
			}
			for _, subcmd := range cmd.comamnds {
				if subcode == subcmd.code {
					subcmd.handle(buf)
					return
				}
			}
		}
	}
}

// cmd table
var cmds = [...]*command{
	{0x00, f006FA9EBhandle00, nil},
	{0x01, handle01, nil},
	{0x02, handle02, nil},
	{0x0D, f007087BFhandle0D, nil},
	{0xF1, nil, []*command{
		{0x00, handleF100versionmatch, nil},
		{0x04, handleF104, nil},
	}},
	{0xFA, nil, nil},
	{0xF4, nil, []*command{
		{0x03, handleF403serverInfo, nil},
		{0x05, handleF405, nil},
		{0x06, handleF406, nil},
	}},
}

func f006FA9EBhandle00(buf []uint8) {
	// SEH
	f00DE8A70chkstk() // 0x4B34
	if v012E2340 == 2 {
		v01319E08log.f00B38AE4printf("Send Request Server List.\r\n")
		// 0x006FAA2B 压缩
		var reqServerList pb // [c1 04 f4 06]
		reqServerList.f00439178init()
		reqServerList.buf[0] = 0xC1
		reqServerList.buf[2] = 0xF4
		reqServerList.buf[3] = 0x06
		reqServerList.len = 4
		reqServerList.buf[1] = uint8(reqServerList.len)
		reqServerList.f004393EAsend(false, false)
		// ...
	}
	// ...
}
func handle01(buf []uint8) {}
func handle02(buf []uint8) {}
func f007087BFhandle0D(buf []uint8) {
	switch buf[3] {
	case 0:
	case 1:
		if v012E2340 == 4 {
			ebp14 := f004A7D34()
			ebp14.f004A9F3B(buf[13:]) // "you are logged in as Free Player"
		}
	}
}
func handleF100versionmatch(buf []uint8) {}
func handleF104(buf []uint8)             {}
func handleFA(buf []uint8)               {}
func handleF403serverInfo(buf []uint8) {
	// close连接服务器
	// f006BF89ADial(ip, port) // 拨号游戏服务器
	// set state
}
func handleF405(buf []uint8) {}
func handleF406(buf []uint8) {}

// hijack cmd
func handlecmdhook(code uint8, subcode uint8, buf []uint8, len int) {
	for _, cmd := range cmdhooks {
		if code == cmd.code {
			if cmd.handle != nil {
				cmd.handle(buf)
				return
			}
			for _, subcmd := range cmd.comamnds {
				if subcode == subcmd.code {
					subcmd.handle(buf)
					return
				}
			}
		}
	}
}

// hijack cmd table
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

// hijack cmd handle
func joinResultSend(buf []uint8)         {}
func charMapJoinResult(buf []uint8)      {}
func charRegen(buf []uint8)              {}
func levelUpSend(buf []uint8)            {}
func levelUpPointAdd(buf []uint8)        {}
func masterLevelUpSend(buf []uint8)      {}
func ansMuBotRecvUse(buf []uint8)        {}
func reconnect(buf []uint8)              {}
func friendRoomCreate(buf []uint8)       {}
func AHCheckGetTickHook(buf []uint8)     {}
func offTradeReq(buf []uint8)            {}
func customPost(buf []uint8)             {}
func setChatColors(buf []uint8)          {}
func IGCStatAdd(buf []uint8)             {}
func enableSiegeSkills(buf []uint8)      {}
func setAgilityFix(buf []uint8)          {}
func setCashItemMoveEnable(buf []uint8)  {}
func luckyItemFix(buf []uint8)           {}
func disableReconnect(buf []uint8)       {}
func handleExitProcess(buf []uint8)      {}
func disableReconnectSystem(buf []uint8) {}
func alterPShopVault(buf []uint8)        {}
func dropSellMod(buf []uint8)            {}
func setBattleZoneData(buf []uint8)      {}
func attackResult(buf []uint8)           {}
func diePlayerSend(buf []uint8)          {}
func refillSend(buf []uint8)             {}
func manaSend(buf []uint8)               {}

func f00735DC7(uk1 int, code int, buf []uint8, len int, enc int) {

}

func f0075C3B2(code int, buf []uint8, len int, x int) {
	// 0x09FE6362
	ebp5D8 := code
	if ebp5D8 > 0xFD {
		// 0x0075FCAC
		return // return 1
	}
	// 0x09FB655C
	// 寻找处理器
	// v0075FCB2cmd[4*v0075FF6Acmd[ebp5D8]]

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

// t6000
var v01310798 t6000

type t6000 struct {
	fs []func()
}

func (t *t6000) f00446D6DreqServerList() {
	// 带SEH处理
	f00DE8A70chkstk() // 0x2994
	// ...
	var ebp149C pb
	ebp149C.f004393EAsend(false, false) // 发送协议报文
	// ...
}

func (t *t6000) f004CCC07() {
	// ...
	_ = t.fs[11]
	_ = t.fs[12] // f00446D6D
}
