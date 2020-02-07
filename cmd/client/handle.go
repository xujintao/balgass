package main

import (
	"encoding/binary"
	"log"
)

type command struct {
	code     uint8
	handle   func(buf []uint8)
	comamnds []*command
}

// cmd
func f0075C3B2handlecmd(code uint8, buf []uint8, len int, enc bool) {
	if h, ok := cmds[int(code)]; ok {
		h(buf)
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
	codes := []byte{code, subcode}
	code16 := binary.BigEndian.Uint16(codes)
	if h, ok := cmds[int(code16)]; ok {
		h(buf)
		return
	}
	log.Printf("invalid cmd, code:%02dx, buf: %v", code, buf)
}

// cmd table
var cmds = map[int]func(buf []uint8){
	0x00:   f006FA9EBhandle00,
	0x01:   handle01,
	0x02:   handle02,
	0x0D:   f007087BFhandle0D,
	0xD7:   handleD7position,     // hook D4->D7
	0xD9:   handleD9normalAttack, // hook D9->11
	0xF100: handleF100versionmatch,
	0xF104: handleF104,
	0xF403: handleF403serverInfo,
	0xF405: handleF405,
	0xF406: handleF406serverList,
	0xF6:   f006C18B6handleF6,
}

func f006FA9EBhandle00(buf []uint8) {
	// SEH
	f00DE8A70chkstk() // 0x4B34
	if v012E2340 == 2 {
		// 请求服务器列表
		v01319E08log.f00B38AE4printf("Send Request Server List.\r\n")
		// 0x006FAA2B 压缩
		var reqServerList pb // [c1 04 f4 06]
		reqServerList.f00439178init()
		reqServerList.buf[0] = 0xC1
		reqServerList.buf[2] = 0xF4
		reqServerList.buf[3] = 0x06
		reqServerList.len = 4
		reqServerList.buf[1] = uint8(reqServerList.len)
		// reqServerList.f004393EAsend(false, false) // 直接api发送了
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
			ebp14 := f004A7D34getServiceManager()
			ebp14.f004A9F3B(buf[13:]) // "you are logged in as Free Player"
		}
	}
}
func handleD7position(buf []uint8)     {}
func handleD9normalAttack(buf []uint8) {}
func handleF100versionmatch(buf []uint8) {
	// f006C75A7 被花到 0x0A4E7A49 // version match, buf: 01 2E 87 "10525"
	func(buf []uint8) {
		ebp4 := buf
		ebp10 := buf[4]
		if ebp10 == 1 {
			// ebp8 := &v0130F728
			// ebp8.f004A9123(&ebp8.f4880)
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
func handleF104(buf []uint8) {}
func handleF403serverInfo(buf []uint8) {
	// close连接服务器
	// f006BF89ADial(ip, port) // 拨号游戏服务器
	// set state
}
func handleF405(buf []uint8)           {}
func handleF406serverList(buf []uint8) {}
func handleFA(buf []uint8)             {}
func f006C18B6handleF6(buf []uint8) {
	// c1 05 f6 1a 00
	// 0x0A8FB4BE hook to hide f0A9F6026, jump complicated shell logic, which would send trap message
	var label1 uint32 = 0x00FF20B1
	// push label1
	// push 0x09E29528
	// ret

	// 0x09E29528
	var label2 uint32 = 0x0042D75C
	// push 0x0042D75C
	// push 0x09E27458
	// ret

	// 0x09E27458
	var label3 uint32 = 0x0A88AE46
	// push 0x0A88AE46
	// push 0x012DDF16
	// ret

	// 0x012DDF16 0x0ABF9D31 0x0ABE1253
	// push edi esi ebx ecx fd edx

	// 0x09FDE30A
	// push 0x0A32AC32
	// push 0x0A4425A9
	// ret

	// f0A4425A9
	func() {
		label1 = 0x0A9F6026
		label2 = 0x0AD70BD7
		// f0A4EC883 send trap message
		func(unk uint32) {
			var ebp148Ctrap pb
			ebp148Ctrap.f00439178init()
			ebp148Ctrap.f0043922CwritePrefix(0xC1, 0xF3)
			ebp148Ctrap.f004397B1writeUint8(0x31)
			// big endian
			ebp148Ctrap.f004397B1writeUint8(uint8(unk >> 24))
			ebp148Ctrap.f004397B1writeUint8(uint8(unk >> 16))
			ebp148Ctrap.f004397B1writeUint8(uint8(unk >> 8))
			ebp148Ctrap.f004397B1writeUint8(uint8(unk))
			ebp148Ctrap.f004393EAsend(true, false)
			ebp148Ctrap.f004391CF()

			// for {
			// 	ebp1490 := 0
			// 	if ebp1490 >= 0x190 {
			// 		break // 0x007C8001
			// 	}
			// 	if ebp1490 == v08C88CAC {
			// 		break
			// 	}
			// 	ebp1494 := f004373C5(ebp1490).f00A38D5BgetObject()
			// 	if ebp1494.m5E {
			// 		break
			// 	}
			// 	ebp1490++
			// }
			// ebp1494.m404 = 0x2710
		}(0x00007533)
		// 0x0A9F747B
		label3 = 0x0A4DD48C
	}()

	// 0x0A32AC32
	// pop edx fd ecx ebx esi edi
	// label3
	// label2
	// label1

	// f0A9F6026 隐藏函数
}

// hijack cmd
func handlecmdhook(code uint8, subcode uint8, buf []uint8, len int) {
	if h, ok := cmdhooks[int(code)]; ok {
		h(buf)
		return
	}
	codes := []byte{code, subcode}
	code16 := binary.BigEndian.Uint16(codes)
	if h, ok := cmdhooks[int(code16)]; ok {
		h(buf)
		return
	}
	log.Printf("invalid cmd, code:%02dx, buf: %v", code, buf)
}

// hijack cmd table
var cmdhooks = map[int]func(buf []uint8){
	0xF100: joinResultSend,
	0xF303: charMapJoinResult,
	0xF304: charRegen,
	0xF305: levelUpSend,
	0xF306: levelUpPointAdd,
	0xF351: masterLevelUpSend,
	0xBF51: ansMuBotRecvUse,
	0xF403: reconnect,
	0xCA:   friendRoomCreate,
	0xD4:   AHCheckGetTickHook,
	0xFA0D: offTradeReq,
	0xFA10: customPost,
	0xFA11: customPost,
	0xFA12: setChatColors,
	0xFA90: IGCStatAdd,
	0xFAA0: enableSiegeSkills,
	0xFAA2: setAgilityFix,
	0xFAA3: setCashItemMoveEnable,
	0xFAA4: luckyItemFix, // 区分C1和C2
	0xFAA5: disableReconnect,
	0xFAA6: handleExitProcess,
	0xFAA7: disableReconnectSystem,
	0xFAA8: alterPShopVault,
	0xFAB0: dropSellMod,
	0xFAF8: setBattleZoneData,
	0x11:   attackResult,
	0x17:   diePlayerSend,
	0x26:   refillSend,
	0x27:   manaSend,
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
	f0075C3B2handlecmd(uint8(code), buf, len, func() bool { return x == 1 }())
}
