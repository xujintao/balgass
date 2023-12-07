package main

import (
	"encoding/binary"
	"log"
	"unsafe"

	"github.com/xujintao/balgass/win"
)

// s9 f0067084A
func f0075C3B2handlecmd(code uint8, buf []uint8, len int, enc bool) {
	// 0x09FB655C
	if h, ok := cmds[int(code)]; ok {
		h(code, buf, len, enc)
		return
	}
	// 0x0075CAFD 0x0075FCAC
	log.Printf("invalid cmd, code:%02dx, buf: %v", code, buf)
}

// cmd table
// key: 0x0075FF6A, s9 0x00673FE6
// value: 0x0075FCB2, s9 0x00673D2E
var cmds = map[int]func(code uint8, buf []uint8, len int, enc bool){
	0x00:   f0075CB89handle00,                // server_connect is prepared
	0x01:   f0075CB97handle01,                // chat
	0x0D:   f007087BFhandle0D,                // handle notice message
	0x1D:   handle1DBeAttacked,               // hook, hash[DF]=hash[1D]
	0x12:   handleViewportPlayerAdd,          // viewport player add
	0x13:   handleViewportMonsterNPCAdd,      // viewport monster/npc add
	0x14:   handleViewportObjRemove,          // viewport object add
	0x1F:   handleViewportMonsterCallAdd,     // viewport monstercall add
	0x20:   f0075CC8Bhandle20ViewportItemAdd, // viewport item add
	0x21:   handleViewportItemRemove,         // viewport item remove
	0x26:   f0075CE21handle26hpsd,            // hp and sd
	0x27:   f0075CE2Fhandle27mpag,            // mp and ag
	0x42:   f0075D021handlePartyInfo,         // party info, 客户端主动请求以及队伍成员信息变化推送
	0x44:   f0075D03DhandlePartyHPMP,         // party member HP/MP, 队伍成员HP/MP数据变化了服务器才会推送而不是定时发送
	0x65:   handleViewportGuildInfoAdd,       // viewport guild info add
	0x45:   handleViewportPlayerChange,       // viewport player change
	0x81:   f0075D93DhandleWarehouseMoney,    // warehouse money
	0xA9:   f0075E87EhandlePetItemInfo,       // pet item info
	0xBF:   f0075F391handleBF,                // event experience
	0xD2:   f0075F7D0handleD2,                // cash shop
	0xD7:   handleD7positionSet,              // hook, hash[D4]=hash[D7]
	0xD9:   handleD9normalAttack,             // hook, hash[11]=hash[D9]
	0xDA:   handleDApositionGet,              // hook, hash[15]=hash[DA]
	0xF0:   nil,                              // reserved
	0xF1:   f0075C3E8handleF1,                // server_game is prepared and response with server's version, and the login logic also use code F1
	0xF2:   nil,                              // reserved
	0xF3:   f0075C794handleF3,                // character
	0xF4:   f0075CB02handleF4,                // server list and server info
	0xF5:   nil,                              // reserved
	0xF6:   f006C18B6handleF6,                // quest list
	0xF805: handleViewportGensInfoAdd,        // viewport gens info add
	0xFA:   nil,                              // reserved
	0xFB:   nil,                              // reserved
	0xFC:   nil,                              // reserved
}

func f0075CB89handle00(code uint8, buf []uint8, len int, enc bool) {
	// f006FA9EB(buf)
	func(buf []uint8) {
		// SEH
		f00DE8A70chkstk() // 0x4B34
		if v012E2340state == 2 {
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
	}(buf)
}

// chat
func f0075CB97handle01(code uint8, buf []uint8, len int, enc bool) {
	// f006C8922(buf)
	func(buf []uint8) {
		// ebpC := buf
		ebp8id := int(binary.BigEndian.Uint16(buf[3:]))
		ebp4obj := v01308D04objectPool.f00A38D5BgetObject(f00592888getIndex(ebp8id))
		if ebp4obj == nil {
			return
		}
		if v0805BBACobjectself.m18guildTitle == 0x80 && f00DE94F0strcmp(ebp4obj.m38name[:], []uint8("길드 마스터")) == 0 { // 战盟使者
			// f0090E94C().f0090D6A6(21)
			// v086A3BDE = 1
			// v012E31B4 = 0
			// f005A4C09(0)
			// v012E31C4 = 8
			// v012E31B8 = 0
			if v0805BBACobjectself.m5C == -1 {
				v0805BBACobjectself.m5C = 1999
			}
			return
		}
		// f005C72CD(ebp4obj.m38name[:], buf[5:], 0, -1)
	}(buf)
}

func f007087BFhandle0D(code uint8, buf []uint8, len int, enc bool) {
	subcode := buf[3]
	switch subcode {
	case 0:
	case 1:
		if v012E2340state == 4 {
			ebp14 := f004A7D34getWindowManager()
			ebp14.f004A9F3B(buf[13:]) // "you are logged in as Free Player"
		}
	}
}

func f0059D4F6bit4changeup2(class uint8) bool {
	// return f0059D1AF(x)
	return func() bool {
		if class>>4&1 > 0 { // bit4
			return true
		}
		return false
	}()
}

func handle1DBeAttacked(code uint8, buf []uint8, len int, enc bool) {}

func handleViewportPlayerAdd(code uint8, buf []uint8, len int, enc bool) {}

func handleViewportMonsterNPCAdd(code uint8, buf []uint8, len int, enc bool) {}

func handleViewportObjRemove(code uint8, buf []uint8, len int, enc bool) {}

func handleViewportMonsterCallAdd(code uint8, buf []uint8, len int, enc bool) {}

func f0075CC8Bhandle20ViewportItemAdd(code uint8, buf []uint8, len int, enc bool) {
	// f006C0242(buf)
	func(buf []uint8) {
		var items []struct {
			number [2]uint8 // 0 ~ 299
			x      uint8
			y      uint8
			attr   [12]uint8
		}
		for _, item := range items {
			var ebp1Cord [2]float32
			ebp1Cord[0] = (float32(item.x) + 0.5) / 100 // (item.x + v0114A660) / v0114C250
			ebp1Cord[1] = (float32(item.y) + 0.5) / 100 // (item.y + v0114A660) / v0114C250
			ebp20number := int(binary.BigEndian.Uint16(item.number[:]))
			ebp10msb := ebp20number >> 15
			ebp20number &= 0x7FFF
			v086D09C8items[ebp20number].item2.f00A2742Einit()
			v086D09C8items[ebp20number].number = ebp20number
			// v086D09C8items[ebp20number].item.m334 = v012E9995
			f00677421mapItemAdd(&v086D09C8items[ebp20number], item.attr[:], ebp1Cord[:], ebp10msb)
		}
	}(buf)
}

func handleViewportItemRemove(code uint8, buf []uint8, len int, enc bool) {}

// hp and sd
func f0075CE21handle26hpsd(code uint8, buf []uint8, len int, enc bool) {
	// f006CF07C(buf)
	func(buf []uint8) {
		// 0x28局部变量
		ebp4 := buf
		ebp20subcode := ebp4[3]
		switch ebp20subcode {
		case 0xFD:
			// v08C88F98 = 0
			return
		case 0xFE: // max
			if f0059D4F6bit4changeup2(v0805BBACobjectself.m13class) {
				v08C88E58hpMax = binary.BigEndian.Uint16(ebp4[4:])
				v08C88E5CsdMax = binary.BigEndian.Uint16(ebp4[7:])
			} else {
				v086105ECpanel.m126hpMax = binary.BigEndian.Uint16(ebp4[4:]) // the difference between v086105ECpanel and v0805BBACobjectself
				v086105ECpanel.m12CsdMax = binary.BigEndian.Uint16(ebp4[7:])
			}
			return
		case 0xFF: // current
			v086105ECpanel.m122hp = binary.BigEndian.Uint16(ebp4[4:])
			v086105ECpanel.m12Asd = binary.BigEndian.Uint16(ebp4[7:])
			return
		}
	}(buf)
}

// mp and ag
func f0075CE2Fhandle27mpag(code uint8, buf []uint8, len int, enc bool) {
	// f006C0559(buf)
	func(buf []uint8) {
		ebp4 := buf
		ebp8subcode := ebp4[3]
		switch ebp8subcode {
		case 0xFE: // max
			if f0059D4F6bit4changeup2(v0805BBACobjectself.m13class) {
				v08C88E5AmpMax = binary.BigEndian.Uint16(ebp4[4:])
				v08C88E5EsdMax = binary.BigEndian.Uint16(ebp4[6:])
			} else {
				v086105ECpanel.m128mpMax = binary.BigEndian.Uint16(ebp4[4:])
				v086105ECpanel.m142agMax = binary.BigEndian.Uint16(ebp4[6:])
			}
			return
		case 0xFF: // current
			v086105ECpanel.m124mp = binary.BigEndian.Uint16(ebp4[4:])
			v086105ECpanel.m140ag = binary.BigEndian.Uint16(ebp4[6:])
			return
		}
	}(buf)
}

// party info
func f0075D021handlePartyInfo(code uint8, buf []uint8, len int, enc bool) {
	// f006D4442
	func(buf []uint8) {
		var party msgPartyInfo
		// ebp149CplayerInfo := f00A49798ui().m160playerInfo
		// ebp149CplayerInfo.f00A9396B()

		ebp14A4partyFrame := f00A49798ui().m184partyFrame
		ebp14A4partyFrame.f00A82218setPartyInfo(int(party.memberSize), party.members)

		// ebp14ACnaviMap := f00A49798ui().m1C8naviMap
		// ebp14ACnaviMap.f00A3A41F()

		// request party member position info C1:E7
	}(buf)
}

// party member HP
func f0075D03DhandlePartyHPMP(code uint8, buf []uint8, len int, enc bool) {
	// f006D5608(buf)
	func(buf []uint8) {
		var party struct {
			memberSize uint8
			members    []struct {
				HPPercent int8
				MPPercent int8
				Name      [11]uint8
			}
		}

		var ebp1Cname [11]uint8
		for i := 0; i < int(party.memberSize); i++ {
			ebp10member := party.members[i]
			f00DE8100memset(ebp1Cname[:], 0, 11)
			f00DE9370strncpy(ebp1Cname[:], ebp10member.Name[:], 10)
			ebp24partyFrame := f00A49798ui().m184partyFrame
			if ebp24partyFrame.f00A83122validateName(ebp1Cname[:]) == false {
				v01319E08log.f00B38AE4printf("Party Member is Missing!! / index : %d \r\n", i)
				continue
			}
			var ebp34MP int8
			var ebp30MP int8 // ebp30MP = ebp10member.MPPercent<0 ? 0 : ebp10member.MPPercent
			if ebp30MP <= 100 {
				var ebp38MP int8 // ebp38MP = ebp10member.MPPercent<0 ? 0 : ebp10member.MPPercent
				ebp34MP = ebp38MP
			} else {
				ebp34MP = 100
			}

			var ebp40HP int8
			var ebp3CHP int8 // ebp3CHP = ebp10member.HPPercent<0 ? 0 : ebp10member.HPPercent
			if ebp3CHP <= 100 {
				var ebp44HP int8 // ebp44HP = ebp10member.HPPercent<0 ? 0 : ebp10member.HPPercent
				ebp40HP = ebp44HP
			} else {
				ebp40HP = 100
			}

			ebp2CpartyFrame := f00A49798ui().m184partyFrame
			ebp2CpartyFrame.f00A82A72setPartyHPMPPercent(ebp1Cname[:], ebp40HP, 100, ebp34MP, 100)
		}
	}(buf)
}

func handleViewportGuildInfoAdd(code uint8, buf []uint8, len int, enc bool) {}
func handleViewportPlayerChange(code uint8, buf []uint8, len int, enc bool) {}

func f0075D93DhandleWarehouseMoney(code uint8, buf []uint8, len int, enc bool) {
	// f006C07B0(buf)
	func(buf []uint8) {
		if buf[3] == 0 {
			return
		}
		v086105E8player.m1C08moneyWarehouse = uint(binary.LittleEndian.Uint32(buf[4:]))
		v086105E8player.m1C04money = uint(binary.LittleEndian.Uint32(buf[8:]))
	}(buf)
}

func f0075E87EhandlePetItemInfo(code uint8, buf []uint8, len int, enc bool) {
	// 0x006EAB21: jmp 0x09DEB47B
	// hook to hidden f0AF12EDA directly, jump over complicated shell logic
	// that would send trap message c1 08 f3 31 ...

	// 0x09DEB47B
	var label1 uint32 = 0x009E5B5B
	// push label1
	// push 0x09FB8327
	// ret

	// 0x09FB8327
	var label2 uint32 = 0x00ECAEB6
	// push label2
	// push 0x0AC398A9
	// ret

	// 0x0AC398A9
	var label3 uint32 = 0x007AB22A
	// push label3
	// push 0x0A56EABC
	// ret
	println(label3, label2, label1)

	// f0AF12EDA 隐藏函数
}

func f0075F391handleBF(code uint8, buf []uint8, len int, enc bool) {
	// 0x0A047622
	if buf[0] != 0xC1 {
		// 0x0075F3B9 0x0AD75D2F
	}
	// 0x0A32E264
	subcode := buf[3]
	// 0x0075F3D2 0x0A43A189
	switch subcode {
	case 0x18:
		// 0x0075F65C
	case 0x52: // event experience multiple
		// 0x0075F686
		// f006F1C2D
		func(buf []uint8) {
			// 0x0AAB24A3 0x0A440F8F 0x0A0C7617 0x0AA08BE8 0x0AF979E6
			ebpCpcBang := int16(binary.LittleEndian.Uint16(buf[4:]))
			ebp8eventExpMultiple := int16(binary.LittleEndian.Uint16(buf[6:]))
			ebp8goldLineExpMultiple := int16(binary.LittleEndian.Uint16(buf[8:]))
			if ebpCpcBang < 0 {
				ebpCpcBang = 0
			}
			// 0x006F1C68 0x0AF7F5A0
			if ebp8eventExpMultiple < 0 {
				ebp8eventExpMultiple = 0
			}
			// 0x006F1C76 0x09FFE867
			if ebp8goldLineExpMultiple < 0 {
				ebp8goldLineExpMultiple = 0
			}
			// 0x006F1C84 0x09E6ABC6 0x006F1CA4 0x0A9F856B 0x0B287414 0x006F1CC2
			f00A49798ui().m8CmainFrame.f00AAB45EsetExpMultiple(ebpCpcBang, ebp8eventExpMultiple, ebp8goldLineExpMultiple)
			f00A49798ui().m8CmainFrame.f00AAB447draw(false)
		}(buf)
	}
}

// s9 f00673854
func f0075F7D0handleD2(code uint8, buf []uint8, len int, enc bool) {
	// 0x0A5FBCA6 0x0A5D4C5A 0x0A4381F1
	subcode := buf[3]
	subcode--

	// 0x00760216 table, s9 0x00674292
	switch subcode {
	case 0: // s9 0x00673894
		// 0x0075F810 0x0075F797 0x0A0C1E4B 0x0075F813
		// f006C28D9, s9 f0066B85C
		func(buf []uint8) {
			// 0x0A6014ED 0x09F8B6AA 0x0A9044AC 0x09FBB3E7 0x0ABE308B
			// totalCash := *(*float64)(unsafe.Pointer(&buf[5]))
			// 0x006C28EE f009F0787
			// 0x006C287C 0x0ABF9A98 0x006C28F5
			// v09D8D548.f009F1748(totalCash)

			// 0x0A55EA23 0x0A1938CA
			// cashC := *(*float64)(unsafe.Pointer(&buf[13]))
			// 0x006C2905 f009F0787
			// 0x006C2898 0x0A9F59D7 0x006C290C
			// v09D8D548.f009F175F(cashC)

			// 0x0A5FED76
			// cashP := *(*float64)(unsafe.Pointer(&buf[21]))
			// 0x006C291C f009F0787
			// 0x006C28A3 0x0AF7CA9A 0x006C2923
			// v09D8D548.f009F1776(cashP)

			// 0x006C28BA 0x0AF76971
			return
		}(buf)
	case 1:
		// f0075F81E
	}
}
func handleD7positionSet(code uint8, buf []uint8, len int, enc bool)  {}
func handleD9normalAttack(code uint8, buf []uint8, len int, enc bool) {}
func handleDApositionGet(code uint8, buf []uint8, len int, enc bool)  {}

// version match
func f006C75A7handleF100(buf []uint8) {
	// jmp 0x0A4E7A49 hook to jmp 0x09F8A350
	// ebp4 := buf
	connRes := struct {
		resultOK bool
		number   [2]uint8
		version  [8]uint8
	}{}
	if v08C88F74 != 0 {
		// 0x0A43EEE8
		v08C88E0CconnID = binary.BigEndian.Uint16(connRes.number[:])
		// 0x006C75DA 0x0A38D7F5
		// v0131A290.f004FF0C2()
	} else {
		// 0x006C75E4 0x0AF1317E
		ebp8 := f004A7D34getWindowManager()
		ebp10resultOK := connRes.resultOK
		if !ebp10resultOK {
			// 0x0AA31A3D 0x006C7632 0x0AFD7872 0x006C763C 0x0A5597CD
			v01319E08log.f00B38AE4printf("Connectting error")
			// 0x006C764A
			v01319E08log.f00B38D49(1)
			// 0x0AD2E325 0x006C7656
			ebp8.f004A9EEB(112, 0)
			// 0x006C7668
		} else {
			// 0x006C75FD 0x0A9001D6 0x006C7609
			ebp8.f004A9123(&ebp8.s3)
			// 0x0A83F1CD 0x0A3287DB
			v08C88E0CconnID = binary.BigEndian.Uint16(connRes.number[:])
			// 0x09FC8D15
			v08C88E08 = 2
			// 0x006C765B
		}
		// 0x09E6A47E
		var ebpC uint8
		// 0x006C7668 0x0AA3097E
		for ebpC < 5 { // hook to 8
			// 0x012DDC6C 0x0AD8344E 0x0A334079 0x0ABD6962
			if v012E4018versionDLL[ebpC] != buf[7+ebpC]+ebpC+1 {
				// 0x0A92FF69 0x006C7698
				ebp8.f004A9146LRU(&ebp8.s3)
				// 0x0AD7F7D5 0x006C76A4
				ebp8.f004A9EEB(111, 0)
				// 0x0B10C6BC 0x006C76B3 0x006C7673 0x0A442DBB
				v01319E08log.f00B38AE4printf("Version dismatch - Join server.\r\n")
			}
			// 0x006C76BA 0x006C766E 0x0A445B1B 0x006C7661 0x09EB5D35
			ebpC++
		}
	}
	// 0x006C76BC 0x0AAB7BDF 0x006C76C1
	// v013523B0.f0051FB61()
	return
}

func f0075C3E8handleF1(code uint8, buf []uint8, len int, enc bool) {
	// 0x0AD35DC8 0x0AD3DDF9 0x0A9FA1B8
	subcode := buf[3]
	switch subcode {
	case 0: // connect result and version match
		// 0x0075C46E 0x0075C3F0 0x0A4DE2B7 0x0075C471
		f006C75A7handleF100(buf)
	case 1: // login result
		// 0x0075C47C
	case 2:
		// 0x0075C6FE
	case 3:
		// 0x0075C765
	case 4:
		// 0x0075C770
	case 5:
		// 0x0075C77B
	case 6:
		// 0x0075C786
	case 0x12:
		// 0x0075C718
	}
}

var v09FD6693 = [...]uint{
	0x7C847811,
	0x9A04CD68,
	0x21295446,
	0xECDDBC71,
	0xBA156216,
}
var v0A05044A = [...]uint{
	0x1BC15B10,
	0x75C966E2,
	0xB99388B8,
	0xFCEFE807,
	0x79C783E6,
}

// cipher key
var v0A846872 = [4]uint{
	0x6D6DA521,
	0x50FF5BCA,
	0x51F68227,
	0x45C5A1EC,
}

// block 0: (v0A8FCD55, 4)
var v0A8FCD55 = [4]uint8{0x9E, 0x01, 0x1D, 0xC6}

// block 1~n, cipher text
var v0A7445E4 = [485]block{
	{0x4D1250B3, 0x91A3D473}, // block 1: v00B98ED0, 10
	{0x02EAB896, 0x966E5E5B}, // block 2: v00B98EE5, 15
	// ...
	{0xEDF20833, 0x61320F98}, // 0x0A7455FC, block 484, ?
	{0x09609CAB, 0xDCB5C892}, // 0x0A745504, block 485: 0xFFFFFFFF, 0xFFFFFFFF
}

// xor key
var v0A935B1D uint = 0xAE5D1C94
var v0AD3C20F uint = 0x306A652D

// xor key
var v0B10C78C = [4]uint{}
var v09EAC72E = [4]uint{}

var v0AF7C3C4label1 = 0x0ABF8601
var v0AFD83CAlabel1 = 8
var v0AD7BED9label2 = 0x0A55FB91
var v0ABD891Flabel2 = 7
var v09FE4BFFlabel3 = 0x09E733E8
var v0AF836B5label3 = 6

// 一个字符串占了半页
var v09D96B40 = [10][2048]uint8{} // "Check Integrity... : data/local/Gameguard.csr"
var v09D9AB40 = 6

// character list
func f006C780DhandleF300(buf []uint8) {
	// 0x0A3A79D0 0x006C7813
	f006C64A1init()
	// 0x0AD82BA1 0x006C7826 0x09E731C7
	// ebp8 := buf
	f004DAAB3setMaxClass(buf[4])
	if v012E2340state == 4 {
		v012E3EC8mapNumber = 106
	}
	ebp4offset := 8
	for ebpC := 0; ebpC < int(buf[6]); ebpC++ {
		// 0x006C7853 0x0A8FA634
		// 0x0A55E0C6 0x0AFDCF11 0x006C7874 0x09FC90D1
		// ebp20char := buf[ebp4offset:]
		ebp20char := struct {
			index       uint8
			name        [10]uint8
			level       uint16
			ctlcode     uint8
			class       uint8
			inventory   [17]uint8
			guildstatus uint8
			pklevel     uint8
		}{}
		ebp1Cclass := f0059CFA6class(ebp20char.class)
		var ebp14 float32
		var ebp10 float32
		var ebp18dir float32
		if ebp20char.index >= 0 && ebp20char.index <= 4 {
			// 0x006C789B 0x0AF972CC 0x0A5584D2
			ebp14 = 0.0
			ebp10 = 0.0
			ebp18dir = 0.0
		}
		// 0x006C78AA 0x0A33925D 0x09FB4DD9 0x09FE45BB
		// 0x0A050C84 0x0ABDA84D 0x006C78CB 0x0A9325F7
		ebp18dir = 0.0
		ebp10 = 0.0
		ebp14 = 0.0
		ebp24obj := f0059CA40newObject(int(ebp20char.index), ebp1Cclass, 0, ebp14, ebp10, ebp18dir)
		ebp24obj.m164level = ebp20char.level
		ebp24obj.m15ctlCode = ebp20char.ctlcode
		// 0x006C7901 0x0AF766EB 0x0AD336E0
		f00DE7C90memcpy(ebp24obj.m38name[:], ebp20char.name[:], 10)
		ebp24obj.m42 = 0
		// 0x0AF89D0F 0x0A4E5B21 0x0A8FD421 0x006C7922 0x09F8EC05
		f00594B19objSet(int(ebp20char.index), ebp20char.inventory[:], 0, 0)
		ebp24obj.m18guildTitle = ebp20char.guildstatus
		ebp24obj.m20pkLevel = ebp20char.pklevel
		ebp4offset += 36
		// 0x006C784C 0x09E2C1AA
	}
	// 0x006C7950 0x0A84B435
	v08C88E08 = 51
	// 0x006C795A 0x0A43F003 0x0B076F66 0x0A895BFA 0x006C7979
	// f00A49798ui().m104.f00A636DF(buf[7])
	v0131A250 = win.GetTickCount()
}

// character create
func f006C798BhandleF301(buf []uint8) {
	//
	f004A7D34getWindowManager().f004A9EEB(1, 2)
}

// character info
func f00701DF4handleF303(buf []uint8, enc bool) { // (v0018C80C, true)
	// 0x0AC3581C hook to hide 0x0ABF8601, jump complicated shell logic: check integrity on load character info
	var label1 = 0x0A9368DD
	// push label1
	// push 0x0AD33DDE
	// ret

	// 0x0AD33DDE
	var label2 = 0x09FDA253
	// push label2
	// push 0x0A5D4A14
	// ret

	// 0x0A5D4A14
	var label3 = 0x00BF60A2
	// push label3
	// push 0x0AD3B967
	// ret

	// 0x0AD3B967
	// push edi esi ebx edx ecx fd
	// 0x0AC34158
	// push 0x0AF80858
	// push 0x0A604CD4
	// ret

	// 0x0A604CD4 f0A604CD4, 完整性校验
	func() {
		// 0x84局部变量
		v0018C068 := struct {
			m00  [5]uint // 0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476, 0xC3D2E1F0
			m14  uint
			m18  uint
			size uint16
			m1E  [64]uint8
			m60  uint
			m64  uint
		}{}
		var ebp8 block

		// push esi edi eax
		// push 0x0AF76043
		// push 0x0AD91F40
		// ret

		// 0x0AD91F40
		eax := &v0018C068 // eax为什么为0x0018C068?
		// push esi
		if eax == nil {
			// 0x0A392626
			// eax = 1
			// pop esi
			// 0x0AF76043
		} else {
			// 0x09E8C84C 0x09EB665D 0x0AD3DBA8
			edx := 0
			v0018C068.m14 = 0
			v0018C068.m18 = 0
			v0018C068.size = 0
			ecx := 0
			// push edi
			// 0x09E8C4D7
			for edx < 5 {
				// v0018C068.data[edx] = v09FD6693[ecx] ^ v0A05044A[ecx]
				// 0x0A4DD423 0x0A9010BC 0x09F8A6F5 0x09E9201B 0x0AD8957E 0x09FC3702 0x0AD56A99
				edx++
				ecx++
			}
			// 0x0A55D14E 0x09FDE315 0x09FD31D0 0x0AD71C58
			v0018C068.m60 = 0
			v0018C068.m64 = 0
			// pop edi
			// eax = 0
			// pop esi
			// 0x0AF76043
		}

		f0A4424B3calc := func() int {
			// 计算v0018C068.m00[5]
			// 0x15C局部变量
			var ebp1C = [4]uint{}   // 0x5A827999, 0x6ED9EBA1, 0x8F1BBCDC, 0xCA62C1D6
			var ebp15C = [80]uint{} // 0x9E011DC6,
			// push ebx esi edi
			ecx := 0
			eax := 0
			edi := 4
			// 0x09EB0A9C
			for {
				ecx++
				ebp1C[ecx] = v0B10C78C[eax] ^ v09EAC72E[eax] // *(ebp + ecx*4 - 20) = v0B10C78C[eax] ^ v09EAC72E[eax]
				eax++
				if ecx >= edi {
					break // 0x0A43FBE5
				}
			}
			// 0x0A43FBE5
			esi := 0
			data := v0018C068.m1E[:]
			// 0x09F886D1
			for {
				ebp15C[esi] = uint(binary.BigEndian.Uint32(data)) // *(ebp + esi*4 - 0x15C) = uint(binary.BigEndian.Uint32(data))
				esi++
				data = data[edi:]
				if esi > 16 {
					break // 0x09FE811D
				}
			}
			// 0x09FE811D
			esi = 64
			edx := 0
			// 0x0A84AD4E
			for {
				ecx = int(ebp15C[edx+13] ^ ebp15C[edx+8] ^ ebp15C[edx+0] ^ ebp15C[edx+2])
				ebp15C[edx+17] = uint(ecx>>31 | ecx*2)
				edx++
				esi--
				if esi <= 0 {
					break // 0x0A557E99
				}
			}
			// 0x0A557E99
			i := 0
			// 0x0AB4CE8F
			for {
				ecx = int(v0018C068.m00[0])
				edx = int(v0018C068.m00[1])
				esi = int(v0018C068.m00[2])
				edi = int(v0018C068.m00[3])
				ebp1C[1] = uint(edi)
				ebp1C[2] = v0018C068.m00[4]
				ebp1C[3] = uint(ecx)
				ecx = ^edx & edi
				edi = esi & edi
				ecx |= edi
				// ...
				i++
				if i >= 20 {
					break // 0x09EBA091
				}
			}
			// 0x09EBA091
			// for {
			// 	// ...
			// }
			// 0x0A000237
			// 0x0A33B671
			// for{
			// 	// ...
			// }
			// 0x09E037DC
			// 0x0AFDB396
			// for{
			// 	// ...
			// }
			// 0x09E314D1
			v0018C068.size = 0
			// pop edi esi ebx
			return 80
		}
		f0A04D0E9copy := func(buf []uint8, size int) {
			// 0x0AFD91A1 ... 0x0A38C729 0x0AF98DB7
			if size == 0 {
				// 0x0A84B8A2
				// eax = 0
				// return
			}
			// 0x0A0501C2 0x09E70719 0x0AA064A9 0x09FDB3C6 0x09FFEA27
			// push esi edi
			if buf == nil {
				// 0x0AF127E7
				// eax = 1
				// return 0x0B285753
			} else {
				// 0x0AA2DD14 0x0AD7F88C 0x0A888838 0x0A32F731
				if v0018C068.m60 != 0 {
					// 0x0A4447E1
					v0018C068.m64 = 3
					// return 0x0B285753
				}
				// 0x0AF999CA 0x0AAB2477 0x0AF12655 0x0A4E99B7 0x0A05B63E
				if v0018C068.m64 != 0 {
					// return 0x0B285753
				}
				// 0x0A9D5444
				for {
					size--
					// 0x0AC36ABB 0x0ABE21DB 0x0AA30D88 0x0A9F8A08 0x0A041A45 0x0A8F9951 0x0A0584FC
					if v0018C068.m64 != 0 {
						// break 0x0A333884
					}
					// 0x0A742D43 0x0AF80694 0x0ABD8DA9
					v0018C068.m1E[v0018C068.size] = buf[0]
					v0018C068.size++
					v0018C068.m14 += 8
					// 0x0A39139D
					if v0018C068.size == 64 {
						// 0x0ABDA41F
						// push esi
						// push 0x0A8912FD
						// push 0x0A4424B3
						// ret
						// 0x0A4424B3 f0A4424B3()
						f0A4424B3calc()
						// 0x0A8912FD 0x09E2FA6E
					}
					// 0x09E2FA6E
					buf = buf[1:]
					if size == 0 {
						break // 0x0A333884
					}
				} // for loop end 0x0A9D5444
				// 0x0A333884
				// eax = 0
				// return 0x0B285753
			}
			// 0x0B285753
			// pop edi
			// pop esi
			// return
		}

		f09FC0148dec := func(key []uint) {
			// 0x18局部变量
			ebp4 := ebp8.addr
			edx := v0A935B1D ^ v0AD3C20F
			eax := edx << 5
			// push esi
			esi := ebp8.size
			// push edi
			// ebp10 := key[3]
			// ebpC := key[2]
			// ebp18 := key[1]
			// ebp14 := key[0]
			cnt := 32

			// 0x0A436A1C
			for {
				// esi -= (ebp4<<4 + ebpC) ^ (ebp4>>5 + ebp10) ^ (ebp4 + eax)
				// ecx := (esi<<4 + ebp14) ^ (esi>>5 + ebp18) ^ (esi + eax)
				// ebp4 -= ecx
				eax -= edx
				cnt--
				if cnt == 0 {
					break // 0x09E2C971
				}
			}
			// 0x09E2C971
			ebp8.addr = ebp4
			ebp8.size = esi

			// pop edi
			// pop esi
			return
		}

		// 计算大块校验和
		// 0x0AF76043, copy section 0 -------------------------------------------
		// push 4
		// push v0A8FCD55[:]
		// push &v0018C068
		// push 0x0A952311
		// push f0A04D0E9
		// ret
		// f0A04D0E9(&v0018C068, v0A8FCD55[:], 4)
		f0A04D0E9copy(v0A8FCD55[:], 4)

		// 0x0A952311, decode and copy section 1 -------------------------------
		// var ebp8 = [2]uint{v0A7445E4[0], v0A7445E4[1]} // 密文，需要异或解密
		ebp8.addr = v0A7445E4[0].addr
		ebp8.size = v0A7445E4[0].size
		// push v0A846872[:]
		// push ebp8[:]
		esi := v0A7445E4[:]
		// push 0x0A32877C
		// push 0x09FC0148
		// ret
		// 0x09FC0148 f09FC0148(&ebp8, v0A846872[:]) // 解密
		f09FC0148dec(v0A846872[:])
		// 0x0A32877C
		label1 = v0AF7C3C4label1 // *(ebp + v0AFD83CAlabel1*4 + 8) = v0AF7C3C4label1
		// 0x0A915FFE
		for {
			if ebp8.addr == uintptr(^uint(0)) {
				break // 0x0A887B7F
			} else {
				// 0x0A338AC2
				if ebp8.addr == uintptr(^uint(1)) {
					// 0x0A0C5AFC
				}
			}
			// 0x09EB2544
			// push ebp8.size                push 10
			// push 0x00400000+ebp8.data     push 0x00B98ED0
			// push &v0018C068
			// push 0x09EB153E
			// push f0A04D0E9
			// ret
			f0A04D0E9copy((*[100]uint8)((unsafe.Pointer(ebp8.addr)))[:], int(ebp8.size))

			// 0x09EB153E, decode and copy section 2 --------------------------------
			esi = esi[2:]
			ebp8.addr = uintptr(esi[0].addr)
			ebp8.size = esi[0].size
			// 0x0A7465E3 0x0AAB49E0 0x09E90C91
			// push v0A846872[:]
			// push &ebp8
			// push 0x0A84AAEA
			// push 0x09FC0148
			// ret
			// 0x09FC0148 f09FC0148(&ebp8, v0A846872[:])
			f09FC0148dec(v0A846872[:])
			// 0x0A84AAEA
			// 0x0A915FFE
		}

		// 0x0A887B7F 开始校验
		// *(ebp+v0ABD891Flabel2*4+8) = v0AD7BED9label2
		label2 = v0AD7BED9label2
		var ebp1C [20]uint8
		// push &ebp1C
		// push &v0018C068
		// push 0x0A32D2F8
		// push 0x0AD3DA6E
		// ret
		// 0x0AD3DA6E f0AD3DA6E(&v0018C068, ebp1C[:]) 计算最后一帧校验和
		func() {
			// push esi
			// esi := &v0018C068
			// if esi == nil || ebp1C == nil {
			// 	// 0x0A9F8D92
			// 	// eax = 1
			// 	// pop esi
			// 	// 0x0A32D2F8
			// }
			// 0x09FDDBAB
			if v0018C068.m64 != 0 {
				// 0x0A4390F3
				// pop esi
				return // 0x0A32D2F8
			}
			// 0x09FFE6EA
			if v0018C068.m60 == 0 {
				// 0x0AFDA74C
				// push edi esi
				// push 0x0A56CE1C
				// push 0x0A842E15
				// ret
				// 0x0A842E15
				// push esi
				// esi = &v0018C068
				v0018C068.m1E[v0018C068.size] = 0x80 // 0x35
				if v0018C068.size <= 0x37 {
					// 0x09E8F0A1
					v0018C068.size++ // 0x36
					if v0018C068.size < 0x38 {
						// 0x0AFD8921
						for {
							v0018C068.m1E[v0018C068.size] = 0
							v0018C068.size++ // 0x37, 0x38
							if v0018C068.size >= 0x38 {
								break // 0x09EB8CE2
							}
						}
					}
					// 0x09EB8CE2 0x0AA0711A 0x0A0C7CA1 0x0AFDE4B8 0x09FD360D 0x0A84CCAB 0x0A90432C
					// 0x0A56D736 0x09FDB4F4 0x0A841B27 0x0AD831E8 0x0AF77ADB
					// little endian
					v0018C068.m1E[56] = uint8(v0018C068.m18 >> 24)
					v0018C068.m1E[57] = uint8(v0018C068.m18 >> 16)
					v0018C068.m1E[58] = uint8(v0018C068.m18 >> 8)
					v0018C068.m1E[59] = uint8(v0018C068.m18)
					// big endian
					v0018C068.m1E[60] = uint8(v0018C068.m14 >> 24)
					v0018C068.m1E[61] = uint8(v0018C068.m14 >> 16)
					v0018C068.m1E[62] = uint8(v0018C068.m14 >> 8)
					v0018C068.m1E[63] = uint8(v0018C068.m14)
					// push 0x0A5FDFCF
					// push 0x0A4424B3
					// ret
					f0A4424B3calc()
					// 0x0A5FDFCF
					// pop ecx esi
				} else {
					// 0x0AC3860F
					v0018C068.size++
					// 0x0A88E950 0x09FD9899 0x0AF9A5FA 0x0B0739E4
					if v0018C068.size < 0x40 {
						// 0x0AF897D5
						for {
							v0018C068.m1E[v0018C068.size] = 0
							v0018C068.size++
							// 0x0AA0D4A4 0x0ABD9808 0x0B10F6EB 0x09E6A6D1 0x0ABB65BD 0x0A4E2DEE 0x0A847593
							if v0018C068.size > 0x40 {
								break // 0x09F8AC4E
							}
						}
					}
					// 0x09F8AC4E 0x09DEB5D0
					// push esi
					// push 0x0AA0514D
					// push 0x0A4424B3
					// ret
					f0A4424B3calc()
					// 0x0AA0514D
					// pop ecx
					// 0x0A0C7A91
					if v0018C068.size < 0x38 {
						// 0x09F89A0A
					} else {
						// 0x09FD73A0
					}
				}
				// 0x0A56CE1C
				// pop ecx
				// eax := 0
				// ecx := 16
				// memset(v0018C068.m14, 0, 16*4)
				v0018C068.m14 = 0
				v0018C068.m18 = 0
				v0018C068.m60 = 1
				// pop edi
			}
			// 0x0A84F492
			eax := 0
			// 0x09FDBFF8, change byte endians
			for {
				edx := eax
				edx &= 3
				ecx := 3
				ecx -= edx
				edx = eax
				edx >>= 2
				edx = int(v0018C068.m00[edx])
				ecx <<= 3
				edx >>= uint(ecx)
				ebp1C[eax] = uint8(edx)
				eax++
				if eax >= 20 {
					break // 0x09FDCEBE
				}
			}
			// 0x09FDCEBE
			eax = 0
			// pop esi
			// 0x0B2846F4
			return // 0x0A32D2F8
		}()

		// 校验
		// 0x0A32D2F8 0x0AD7F532 0x0A56C5AD 0x0AD56B3D
		eax2 := 5 // eax2 := v0A5FA71Clabel // 5
		// pop ecx
		// pop ecx
		eax2--
		// pop edi
		ecx := eax2 // 4
		// pop esi
		// 0x0A84B91E 0x0A5FA4B4 0x09E259D4 0x0A903F4E 0x0A0017AA 0x0A5D6D14 0x09F87D9F 0x0AD3CD2E
		if eax2 > 0 {
			// 0x0AA2FE2E
			for {
				ecx-- // 3, 2, 1, 0, -1
				// v0AFE0A39[ecx] = ebp1C[ecx+1] // v0AFE0A39[3]=ebp1C[4], v0AFE0A39[2]=ebp1C[3], v0AFE0A39[1]=ebp1C[2], v0AFE0A39[0]=ebp1C[1], v0AFE0A39[-1]=ebp1C[0]
				// copy(v0AFE0A35[:], ebp1C[:])
				if ecx < 0 {
					break // 0x0B074B24
				}
			}
		}
		// 0x0B074B24
		if eax2 >= 0 { // eax2 = 4
			// 0x0A43DFB1
			for {
				// ecx := *(ebp + eax2*4 - 0x1C)
				ecx := ebp1C[eax2]
				// ecx -= v0A391F0E[eax2] // v0A391F0E
				if ecx != 0 {
					// 0x0A9FE958, 不一致
				}
				// 0x0AF13005
				eax2--
				if eax2 < 0 {
					// 0x09FCB8D5, 校验通过
					break // 0x0A9F773B
				}
				// 0x0A43DFB1
			}
		}
		// 0x0A9F773B
		// v09FE5B85 = v0A047F75
		label3 = v09FE4BFFlabel3 // *(ebp+v0AF836B5*4+8) = v09FE4BFF
		// return 0x0AF80858
	}()

	// 0x0AF80858
	// pop fd ecx edx ebx esi edi
	// 0x09FB8F6A 0x0AAB2DEA
	// label3 ret
	// label2 0x012DDD4C ret
	// label1 0x0ABF8601 f0ABF8601 隐藏函数 处理F303, s9 f0A28D090

	// 0x0ABF8601 0x00701E11
	f00DE8A70chkstk() // 0xDB80局部变量
	// 0x0A049CFF
	if enc == false {
		// 0x0A970D01 0x00701E26 构造个报文
		var ebp15FC pb
		ebp15FC.f00439178init()
		// ...
	}
	// 0x0070326F 0x0A441D0A
	// v08C88BC6 = 0
	ebp17Cmsg := struct {
		x            uint8  // 4
		y            uint8  // 5
		mapNumber    uint8  // 6
		dir          uint8  // 7
		exp          uint64 // 8
		expNext      uint64 // 10
		points       uint16 // 18
		strength     uint16 // 1A
		agility      uint16 // 1C
		vitality     uint16 // 1E
		energy       uint16 // 20
		hp           uint16 // 22
		hpMax        uint16 // 24
		mp           uint16 // 26
		mpMax        uint16 // 28
		sd           uint16 // 2A
		sdMax        uint16 // 2C
		ag           uint16 // 2E
		agMax        uint16 // 30
		money        int    // 34
		pkLevel      uint8  // 38
		ctlCode      uint8  // 39
		pointsAdd    uint16 // 3A
		pointsAddMax uint16 // 3C
		leadship     uint16 // 3E
		pointsDec    uint16 // 40
		pointsDecMax uint16 // 42
		inventoryExt uint16 // 44
		ruud         int    // 48
	}{}
	// 0x0A56F62D 0x007033DF 0x0AF113E5
	if f0059D4F6bit4changeup2(v086105ECpanel.m10Aclass) {
		// 0x09EB66C4
		v08C88E40exp = ebp17Cmsg.exp
		v08C88E48expNext = ebp17Cmsg.expNext
	} else {
		// 0x007033FF 0x0A848C6C
		v086105ECpanel.m110exp = uint32(ebp17Cmsg.exp)
		v086105ECpanel.m114expNext = uint32(ebp17Cmsg.exp)
	}
	// 0x0A88A77F
	v086105ECpanel.m174points = ebp17Cmsg.points
	v086105ECpanel.m118strength = ebp17Cmsg.strength
	v086105ECpanel.m11Aagility = ebp17Cmsg.agility
	v086105ECpanel.m11Cvitality = ebp17Cmsg.vitality
	v086105ECpanel.m11Eenergy = ebp17Cmsg.energy
	v086105ECpanel.m120leadship = ebp17Cmsg.leadship
	v086105ECpanel.m122hp = ebp17Cmsg.hp
	v086105ECpanel.m126hpMax = ebp17Cmsg.hpMax
	v086105ECpanel.m124mp = ebp17Cmsg.mp
	v086105ECpanel.m128mpMax = ebp17Cmsg.mpMax
	v086105ECpanel.m140ag = ebp17Cmsg.ag
	v086105ECpanel.m142agMax = ebp17Cmsg.agMax
	v086105ECpanel.m12Asd = ebp17Cmsg.sd
	v086105ECpanel.m12CsdMax = ebp17Cmsg.sdMax
	v086105ECpanel.m14CpointsAdd = ebp17Cmsg.pointsAdd
	v086105ECpanel.m14EpointsAddMax = ebp17Cmsg.pointsAddMax
	v086105ECpanel.m150pointsDec = ebp17Cmsg.pointsDec
	v086105ECpanel.m152pointsDecMax = ebp17Cmsg.pointsDecMax
	v086105E8player.m1C04money = uint(ebp17Cmsg.money)
	/*
		// 0x00703776 0x09E9141D
		// ExtensionBagFrame
		if f00A49798ui().m98ExtensionBag == nil {
			// 0x007037CF
		}
		// 0x0070379C 0x0AF7AE01 0x09FDD366 0x0ABD7268 0x007037CA
		f00A49798ui().m98ExtensionBag.f00A5971D(ebp17Cmsg.inventoryExt)
	*/

	/*
		// 0x0A60214B 0x007037E4 0x0070376B 0x0AAB2811 0x007037EC 0x0AF8AEA7
		v012E3EC8mapNumber = ebp17Cmsg.mapNumber
		f006B1B3A(v012E3EC8mapNumber)
		if f005105DC(92).m28 {
			// 0x0A849EFF
		}
		// 0x00703846 0x09E2DC21
		if v012E3EC8mapNumber == 34 { // 狼魂要塞
			// 0x0ABB58C9
		}
		// 0x007048E7 0x0A9F5AEA 0x007048ED 0x0AB4DAF4 0x0070491D 0x0A9FB42D
		f0051FE0F(v012E3EC8mapNumber)
	*/
	v08C88CACid = f00DE8AADrand() % 400
	ebp178obj := v01308D04objectPool.f00A38D5BgetObject(v08C88CACid)
	if ebp178obj == nil {
		// 0x00704CC1
	}
	// 0x0ABF8771 0x00704988 0x0A12E76C
	// f00592B90objSet(ebp178obj, 0x4C4, ebp17Cmsg.x, ebp17Cmsg.y, (ebp17Cmsg.dir-3.14/4)*180/3.14) // (dir-PI/4)*180/PI
	ebp178obj.m5Eid = v08C88E0CconnID
	ebp2A84 := &ebp178obj.m410
	ebp178obj.m13class = v086105ECpanel.m10Aclass
	ebp178obj.m14 = 0
	ebp178obj.m20pkLevel = ebp17Cmsg.pkLevel
	ebp178obj.m15ctlCode = ebp17Cmsg.ctlCode
	ebp2A84.m28 = 1
	// 0x00704A05 0x09F8CAF9
	// f00594247objSet(ebp178obj)
	// if v086A3BE3 {
	// 	// 0x00704A2A
	// }
	// 0x09FE7B46 0x0AD83F3F
	ebp178obj.m32 = v0805BBACobjectself.m32
	// 0x00704A31 0x0A5FEB28
	v0805BBACobjectself = ebp178obj
	// 0x0A4DD230 0x0AF947B4 0x00704A54 0x09FBD0F4
	f00DE7C90memcpy(ebp178obj.m64[:], v086105ECpanel.m11[:], 255)
	// 0x09FC3478 0x00704A6E 0x0A8FE907
	f00DE7C90memcpy(ebp178obj.m38name[:], v086105ECpanel.m00name[:], 10)

	/*
		// 0x00704A8C 0x0A041DE7
		for ebp2A88 := 0; ebp2A88 < 12; ebp2A88++ {
			// 0x0A4E92F8
			// v086105E8player.m1348inventorys[ebp2A88].m00 = 0
			// 0x00704A7F 0x0A4413A3
		}

		// 0x00704B20 0x00704AA7 0x0AF80484 0x00704B27 0x00704AAE 0x0AD82FE1 0x00704B2E
		f0090E94C().f0090EA51().f008F3ED2()
		// 0x00704ABA 0x0A5FC968 0x00704B3A 0x0ABE33DC 0x00704B51
		f0090E94C().f0090EA51().m138.f006FA001()
		// 0x00704ADD 0x0A05F7FA 0x00704B5D 0x0AB4FF17 0x00704B74
		f0090E94C().f0090EA51().m158.f006F9D22()
		// 0x00704B00 0x09E8F95B 0x00704B80 0x0A8FA447 0x00704B97
		f0090E94C().f0090EA51().m178.f006F9D77()

		// 0x09F85FE8 0x0A3278B9
		ebp178obj.m42 = 0
		// 0x0B28507C 0x0A55AFDB 0x0A4E20E9 0x0A5593C7
		// 0x0A556D91 0x09EB2F1E 0x00704BE7 0x0A008896
		f007EA8A8(0x7D8F, &ebp2A84.m114, &ebp2A84.m120, &ebp2A84.mB4, 0, ebp2A84, -1, 0, 0, 0, 0.0, -1)
		v08C88E08 = 61
		v086A3BDD = 0
		// 0x00704C04 0x09E2E4E6
		f005A4DD4(1, 0)
		v086A3B94 = 30
		v086A3B9C = 0
		v012E3140 = 6
		// 0x00704C2C 0x0A56B7EA
		f00511409(v0805BBACobjectself)
		// 0x0B10C423 0x0AF81D1D 0x0AB4D14C 0x09EB23D9
		// 0x0A442AB5 0x0AFE08CE 0x00704C73 0x0AFD473A
		f007EA8A8(0x7D8F, &ebp2A84.m114, &ebp2A84.m120, &ebp2A84.mB4, 0, ebp2A84, -1, 0, 0, 0, 0.0, -1)

		// 0x0A008E9E 0x0A5D41A1 0x09E2BEB8 0x00704CBC
		ebp2A84.mB0 = 0.0
		v086105E8player.m00 = 0
		v086105E8player.f005A3727(int(ebp178obj.m13class) & 7)

		f0090E94C().f0090E02B()
		f00A49798ui().f00A4E2C8()
		f0090E94C().f0090EBCB().f00936D4B(0)

		// 0x0A4EB782
		v012E3198 = -1
		v012E319C = -1
		v012E31A0 = -1
		v012E31A4 = -1
		v012E31A8 = -1
		v086D09AC = 0
		v0805BBACobjectself.m58 = 0
		// 0x00704D27 0x09F857A2 0x00704D33 0x0AD3551E
		f0055E188(v0805BBACobjectself)
		if f004F4048isBloodCastle(v012E3EC8mapNumber) == false {
			// 0x0AF7A6EA 0x00704D47 0x00704CF3 0x09FCAF0F
			f007DB145(139, 1)
		}
		// 0x00704D4E 0x09E2B233 0x00704D54
		if f004F4014isChaosCastle(v012E3EC8mapNumber) == false {
			// 0x0AD9A707 0x00704D68 0x0B072708 0x00704D76
			// 0x00704CFD 0x0B10F5B5
			f007DB145(153, 1)
			f007DB145(158, 1)
		}
		// 0x00704D7D 0x0AF921BC
		if f009A1053isImperial1() == false &&
			f009A546BisImperial2() == false &&
			f009A9C9EisImperial3() == false &&
			f009AF183isImperial456() == false {
			// 0x0A934ACC 0x00704DB4 0x09F87A94 0x00704DC2
			// 0x0A002840 0x00704DD0 0x0A83CF0D 0x00704DDE
			// 0x00704DAD 0x09EBC615
			f007DB145(904, 1)
			f007DB145(905, 1)
			f007DB145(906, 1)
			f007DB145(907, 1)
		}
		// 0x00704DE5 0x0A4E5A7C 0x09FD963F 0x00704E08
		f00A49798ui().m184partyFrame.f00A821E4()
		// 0x0A5FCD46 0x0AA2D86A 0x00704E30
		f00A49798ui().m164BuffFram.f00A99B92()
		// 0x0A39288A 0x00704E3F 0x0AF92E48
		if f0059D4F6bit4changeup2(v0805BBACobjectself.m13class) {
			// 0x0A9FD488 0x00704E57 0x0AF99B8E 0x00704E69
			// 0x00704DF0 0x0A902355 0x00704E70 0x00704DF7
			// 0x0AFE0144 0x00704E77
			f0090E94C().f0090EA51().f008F3F40(f0059CFDB(v0805BBACobjectself.m13class))
		}
		// 0x00704E7C 0x0AAB7E3A 0x00704E9F
		f00A49798ui().m174Alarm.f00A98C4A()
		// 0x09EAE426 0x09E288C5 0x00704EC7
		f00A49798ui().m174Alarm.f00A98CCA()
		// 0x09FE5C9B 0x0A05AD3B 0x00704EEF
		f00A49798ui().m17CMatchingGuild.f00A7BF00()
		// 0x0AF95361 0x0A92FC61 0x00704F17
		f00A49798ui().m180MatchingParty.f00A7EC7B()
		// 0x00704EAE 0x0A442249 0x00704F23
		f009EEA9B().f09EEAF8()
		// 0x00704EB3 0x0ABDFEE1 0x00704F2A 0x00704EB8 0x0A0C232D 0x00704F31
		f009EEA9B().f009EF5A2(0)
		// 0x09F8B6DF 0x09E8C79A 0x00704F5B
		f00A49798ui().m1C8Navimap.f00A3A4DE(1)
		// 0x0AD98847 0x00704F83
		f00A49798ui().m1F8GremoryCase.f00A52EB6()
	*/

	/*
		// 0x00704F0A 0x0A43D20B 0x00704F8B
		ebp64.f00B605D0()
		// 0x09FDEE6F 0x00704FC7
		ebp4 := 4
		ebp64.f00B60770(&ebp174)
		// 0x0AFDFF4B 0x00704FD8 0x0A329469
		var ebp16C [255]uint8
		f00DE817Asprintf(ebp16C[:], "data\\local\\Gameguard.csr")
		// 0x0A32B2FF 0x00704FEC
		var ebp402Ctmpstr stdstring
		ebp402Ctmpstr.f00406FC0stdstring(ebp16C[:])
		// 0x0AD7EC66 0x00704FFF 0x0A7AAB86 0x0AD33C03
		ebp4 = 6
		rg := false
		// rg := !ebp64.f00B607F0(&ebp402Ctmpstr) // 耗时检测 s9 f00A83580
		func() bool {
			// 0x0AA30CC3 0x0A4EB1BE
			// SEH设置
			// 0x0ABB6271 0x09FB828C
			// 0x768局部变量
			ebp760 := ""
			ebp75C := ""
			// ebp46C := v012F7B90 ^ ebp
			// 0x0AD83FBC
			// ebp758 := ecx
			if ebp402Ctmpstr.m18 < 16 { // 31
				// 0x00B6083B 0x0A33B326
				ebp75C := "" // &ebp402Ctmpstr.m04
			} else {
				// 0x0AFDEEC5
				ebp75C := ebp402Ctmpstr.m04 // "data/local/Gameguard.csr"
			}
			// 0x00B60847 0x0A33B482 0x00B60853 0x0A8F96C5 0x00B60862
			// f00B4AA83
			// func(fmt, text string) []uint8 {
			// 	buf := v09D96B40[v09D9AB40%8][:]
			// 	f00DF0805(buf, fmt, text)
			// 	v09D9AB40++
			// 	return buf
			// }("Check Integrity... : %s", ebp75C)
			// ebp4C4.f00406FC0stdstring(f00B4AA83("Check Integrity... : %s", ebp75C))
			// 0x09F8EEC5 0x0A9320B0 0x00B60882
			ebp4 := 0
			// push 0x8000000A
			// push 1
			// push &ebp4C4
			// ebp64.f00B61710(&ebp4C4, 1, 0x8000000A)
			// 0x0AD9A333 0x00B60898
			ebp4 = -1
			// ebp4C4.f00407AC0(1, 0)
			// 0x0A901716
			if ebp402Ctmpstr.m18 < 16 {
				// 0x00B608B4 0x0A44422D
				ebp760 := "" // &ebp402Ctmpstr.m04
			} else {
				// 0x0A4ED735
				ebp760 := ebp402Ctmpstr.m04
			}
			// 0x00B608C0 0x0A333B7C 0x00B608CC
			ebp14 := f00DE909Efopen(ebp760, "rb") // "data/local/Gameguard.csr", "rb"
			// 0x0A556EA3
			if ebp14 == nil {
				// 0x0AD3E7EC
			}
			// 0x00B60961 0x0A05BAA4 0x00B60969
			f00DEFA34fseek(ebp14, 0, SEEK_END)
			// 0x0A7ABF10 0x00B60975
			ebp34size := f00DEFCD4ftell(ebp14)
			// 0x0B0728F8 0x00B60988
			f00DEFA34fseek(ebp14, 0, SEEK_SET)
			// 0x0A05D721 0x0B600994 0x0A04BE04 0x0AF8F9C2 0x0A561B68 0x00B609AF 0x0A32ADF2
			ebp4E4 := f00DE64BCnew(ebp34size)
			ebp10 := ebp4E4
			ebp4E8 := f00DE64BCnew(ebp34size)
			ebp40 := ebp4E8
			// 0x00B609D4
			f00DE8FBDfread(ebp10, 1, ebp34size, ebp14)
			// 0x0A557BE6 0x00B609E0
			f00DE8C84close(ebp14)
			// 0x0A83C47B 0x00B609F0 0x00B6097F 0x0A8FF4EB 0x00B609F7
			f00B62640once().f00B625E0dec(ebp10, ebp34size) // v09D9AB48.f00B625E0(p, size)
			// 0x0A930885 0x00B60A08 0x00B6099E 0x09E7023F 0x00B60A0F
			f00B62640once().f00B624A0xor(ebp40, ebp10, ebp34size) // v09D9AB48.f00B624A0xor(pDst, pSrc, size)
			// 0x00B60A14 0x00B609A3 0x0AF120FB
			ebp30.f00B4BF7C() // ebp30 = &v01182470
			// 0x0A5D4AE5 0x0A849272 0x00B60A2E
			ebp4 = 2
			ebp30.f00B4BD0C(ebp40, ebp34size)
			// 0x00B609B9 0x012DDCD7 0x00B60A36 0x0A9FC2F8
			type crcHeader struct {
				m00magic int // 0x4752, "RG" means resource guard
				m04num   int // 0xBB, 187个待校验文件
			}
			ebp3C := new(crcHeader)
			ebp3C := ebp30.f00B4BE05() // ebp3C := ebp30.m04
			if ebp3C.m00magic != 0x4752 {
				// 0x0A8FF7C0
			}
			// 0x00B60AAA 0x0A43DF37 0x00B60AB4
			ebp35 := true
			ebp464.f00B411A4()
			// 0x0AF84996 0x00B60AC0 0x0AD7BD59
			ebp4 = 4
			type crcFile struct {
				val  int
				name string
			}
			ebp18 := make([]crcFile)
			ebp18 := crcFiles
			ebp18 := ebp30.f00B4BE05()
			ebp468cnt := 0
			// 0x00B60AE6 0x0A193CC9
			for {
				if ebp468cnt >= ebp3C.m04num {
					break // 0x00B60D65
				}
				// 0x0A050E36 0x00B60B05
				ebp4A4.f00406FC0stdstring(ebp18.name)
				// 0x0AF96B18 0x00B60B1B 0x0A3315F2 0x00B60B28
				ebp4 = 5
				f00406450(&ebp488, ebp758.f00B60750(&ebp4A4))
				// 0x0A88FAE2
				ebp4 = 6
				if ebp470 >= 16 {
					// 0x09FB807F
				}
				// 0x00B60B4B 0x09EB2740 0x00B60B57 0x0AC9D71F 0x00B60B60 0x0A05BB0E
				ebp768 := &ebp484 // "main.exe"
				if -1 == f00DF553F(ebp768, 0) {
					// 0x00B60CB2
				}
				// 0x0AA2CE76 0x00B60B89 0x09FDF779
				ebp4A8 := 0
				if f00B41CB4(&ebp488, &ebp4A8) == 0 {
					// 0x0ABF83B4
				}
				// 0x00B60B9C 0x0A601D0C
				if ebp4A8 != ebp18.val {
					// 0x00B60C20
				}
				// 0x0A4E56BB
				ebp35 = false
				if ebp48C >= 16 {
					// 0x0B28495F
				}
				// 0x00B60BC4 0x09EB5251 0x00B60BD0 0x09E91315 0x00B60BDC 0x0AD7208D 0x00B60BEB
				ebp76C := &ebp4A0 // "main.exe"
				ebp524.f00406FC0stdstring(f00B4AA83("%s file is modified.", ebp76C))
				// 0x0A32DD52 0x00B60C08
				ebp4 = 7
				ebp758.f00B61710(&ebp524, 2, 0x8000000A)
				// 0x0A56E280 0x00B60C1B
				ebp4 = 6
				ebp524.f00407AC0(1, 0)
				// 0x0A9309E1 0x00B60C2D 0x0A324BE0
				if ebp464.f00B412CE(&ebp488) == 0 {
					// 0x0A38FD97
				}
				// 0x00B60CB0 0x00B60C39 0x0A9FC2E2 0x00B60D29 0x0A55B184 0x0AF87662
				ebp18 = ebp18[1:]
				ebp4 = 5
				// 0x09F8F68A 0x00B60D48
				f00407AC0(1, 0)
				// 0x0A05F891 0x00B60D5B
				ebp4 = 4
				ebp4A4.f00407AC0(1, 0)
				// 0x09FCA3F1 0x00B60AD7 0x0AF0EE3B
				ebp468cnt++
				// 0x00B60AE6
			}
			// 0x00B60D65 0x0B07BEB9 crc finish
			if ebp35 {
				// 0x00B60DDF
			}
			// 0x0A0593C9 0x00B60D78 0x0A604945
			f00406FC0stdstring("Stop checking integrity.")
			// 0x09FD43E0 0x00B60D95
			ebp4 = 10
			ebp758.f00B61710(&ebp578, 1, 0x8000000A)
			// 0x09E27085 0x00B60DA8
			ebp4 = 4
			ebp578.f00407AC0(1, 0)
			// 0x09FC245D 0x00B60DC0
			ebp579 = ebp35
			ebp4 = 2
			ebp464.f00B411B8()
			// 0x0A3328F4 0x00B60DCF
			ebp4 = -1
			ebp30.f00B4C04A()
			// 0x09FD3683
			return ebp579
		}()
		ebp400D := rg
		// ebp4 = 5
		// 0x0A888B4A 0x00705020
		// ebp402Ctmpstr.f00407AC0(1, 0)
		// 0x0A9503C9
		if ebp400D { // hook always false, diabale ResourceGuard
			// x64dbg查找引用无法跨模块，比如0x0A84D09A位置引用了0x00705035，但却搜不到0x0A84D09A，因为跨模块了
			// 0x00705035: f006BD5BBclose()
			// 0x0070503A: jmp 0x09EB027B, ResourceGuard Error, jmp 0x00705044
			// 0x00705044: f00B38AE4() // v01319E08log.f00B38AE4printf()
			// 0x00705067: f00B60690()
			// 0x00706761: return // 0x0075C9B2

			// 0x0A84D095
			v08C88FF0conn.f006BD5BBclose()
			// 0x09EB027B 0x00705044 0x09E8D2C0
			v01319E08log.f00B38AE4printf("> ResourceGuard Error!!\r\n")
			// f00B60690()
			return
		}
	*/

	/*
		// 0x00705077 0x0ABE466D 0x0AD31EA2 0x09FF730D 0x007050AE
		f00A49798ui().m138MapName.do16()
		f0080A5ED()
		// 0x0A9F9C65 0x007050BE 0x0070504E 0x0A7455DF 0x007050C6
		f005105DC(92).f00974BEA(v0805BBACobjectself)
		// 0x0A4EB760
		if f009C3904map82to90() == false {
			// 0x007050DC 0x0B109088 0x007050E6 0x0A9F5570
			var ebpDB74 []uint8
			ebpA6C4text := v08610600textManager.f00436DF1findstdstring(484) // "欢迎您加入"
			if ebpA6C4text.m18cap < 16 {
				// 0x0070510E 0x0A001515
				ebpDB74 = (*[1000]uint8)(unsafe.Pointer(&ebpA6C4text.m04data))[:] // 指针自己优化为值
			} else {
				// 0x09FC2201
				ebpDB74 = ebpA6C4text.m04data
			}
			// 0x0070511D 0x0A391427 0x00705123 0x0ABDA21A
			// 0x0070513C 0x0A887EE9 0x00705151
			var ebp2B8C [255]uint8
			f00DE817Asprintf(ebp2B8C[:], "%s%s", ebpDB74, f0064E23BmapName(v012E3EC8mapNumber))
			var ebp404C stdstring
			ebp404C.f00406FC0stdstring(ebp2B8C[:])
			// 0x0A7AAD68 0x00705165
			ebp4 = 7
			var ebp4068 stdstring
			ebp4068.f00406FC0stdstring(nil)
			// 0x007050F0 0x09FB502B
			ebp4 = 8
			// 0x0070516E 0x09E28D49 0x007051A3
			f00A49798ui().m124ChatWindow.f00A9FB38print(&ebp4068, &ebp404C, 3, 0)
			// 0x0A7AC852 0x007051B6
			ebp4 = 7
			ebp4068.f00407AC0(1, 0)
			// 0x0A05D516 0x007051C9
			ebp4 = 5
			ebp404C.f00407AC0(1, 0)
		}
		// 0x0ABF67D8
		if v012E3EC8mapNumber == 30 && f0090E94C().f0090EB53() {
			// 0x007051D3 0x0070515B 0x0ABF7064 0x007051EE
			f0090E94C().f0090EB53().f0092B229()
		}
		// 0x0A56AF78
		if f009C3904map82to90() == false {
			// 0x0A84667C 0x00705208
			f0090E94C().f0090DC7E(65)
		}
		// 0x0AD7B5D2 帝国要塞
		if !(v012E3EC8mapNumber >= 69 && v012E3EC8mapNumber <= 72) {
			// 0x00705226 0x0070521D 0x0A04B923 0x00705228
			// 0x007051D5 0x00705160 0x0A8406ED 0x0070522F
			f0090E94C().f0090DC7E(69)
		}
		// 0x00705234
		f007C8B6A().f007C8BEE()
		// 0x0B073EE7 0x0070524C
		v01308D04objectPool.f00A390A6(v0805BBACobjectself)
		// 0x0AC378EA 0x00705257 0x0A43DF16
		f004DAACA(v01319D6ChWnd)

		// 0x00705267 0x0AF9434E
		ebpA70C := v08610600textManager.f00436DF1findstdstring(4802) // "黄名、红名角色将受到更加严厉的PK惩罚。"
		// 0x0AD56F1E
		ebpDB78 := ebpA70C.m04data
		// 0x0070529E 0x0AFDD2D2 0x007052AA
		var ebp4084 stdstring
		ebp4084.f00406FC0stdstring(ebpDB78)
		// 0x0A84D1D5 0x007052BE
		ebp4 = 9
		var ebp40A0 stdstring
		ebp40A0.f00406FC0stdstring(nil)
		// 0x00705245 0x0A329236
		ebp4 = 10
		// 0x007052C7 0x09FD9308 0x0A5FD276 0x0AD2CB23 0x0A4EC44D 0x007052FC
		f00A49798ui().m124ChatWindow.f00A9FB38print(&ebp40A0, &ebp4084, 4, 0)
		// 0x0A43A54F
		ebp4 = 9
		// 0x0070530F
		ebp40A0.f00407AC0(1, 0)
		// 0x09FC1931
		ebp4 = 5
		ebp4084.f00407AC0(1, 0)

		// 0x09E22836 0x00705331 0x0B10E3C2
		v08610600textManager.f00436DF1findstdstring(4803) // "角色死亡后损失更多金钱及经验。"
		v08610600textManager.f00436DF1findstdstring(4804) // "角色死亡后有较大概率掉落装备栏及背包道具（包含+7以上的装备及各种贵重道具）。"
		v08610600textManager.f00436DF1findstdstring(4805) // "PVP线路存在强制PVP游戏内容，请玩家谨慎选择。"

		// 0x0B072574 0x09F83B10
		if f00A49798ui().m1F4ImageNotice.m74 == false {
			// 0x007055C0 0x0A88BFD3 0x09E27C7D 0X007055E3
			// 0x0A745C78
			if f00A49798ui().m1F4ImageNotice.f00A3A41F() == false {
				// 0x0AF7673E 0x007055F9
				var ebp400C pb
				ebp400C.f00439178init() // C1 04 F3 26
			}
		}

		// 0x00706681 0x0A05D752
		f00A49798ui().m124ChatWindow.m1D9 = 0
		// 0x007066AB 0x0AF8B0B3
		if f00A49798ui().m124ChatWindow.m1D8 == false {
			// 0x007066E5 0x0A55CFC7
			f00A49798ui().m124ChatWindow.m1D8 = true
			// 0x0070670F 0x0AF8951F 0x0AB50375 0x00706732
			f00A49798ui().f00A9F5F8()
		}
		// 0x0A5653DF 0x0B06F999 0x00706756
		ebp414C := true
		ebp4 = 4
		ebp174 := v0114E7AC[:]
		ebp4 = -1
		ebp64.f00B60690()
		// 0x0AF8ECA0 0x00706761 0x0ABB356A
		return ebp414C
	*/
}

// character add point
func f006C03F3handleF306(buf []uint8) {
	addPointRes := struct {
		result  int
		hpmpMax int
		sdMax   int
		agMax   int
	}{}
	if addPointRes.result&0x10 != 0 {
		// 0x09FBA724
		v086105ECpanel.m174points--
		statKind := addPointRes.result & 0x0F
		switch statKind {
		case 0:
			// 0x006C0461
			v086105ECpanel.m118strength++
		case 1:
			// 0x006C0483
			v086105ECpanel.m11Aagility++
		case 2:
			// 0x006C04A5
			v086105ECpanel.m11Cvitality++
			v086105ECpanel.m126hpMax = uint16(addPointRes.hpmpMax)
		case 3:
			// 0x006C04D7
			v086105ECpanel.m11Eenergy++
			v086105ECpanel.m128mpMax = uint16(addPointRes.hpmpMax)
		case 4:
			// 0x006C0509
			v086105ECpanel.m120leadship++
		}
		// 0x006C0526 0x0ABF96F6
		v086105ECpanel.m142agMax = uint16(addPointRes.agMax)
		v086105ECpanel.m12CsdMax = uint16(addPointRes.sdMax)
	}
	// 0x006C054C 0x0AF95874 0x006C0552
	v086105E8player.f005A36E1()
}

// s9 f00670C47
func f0075C794handleF3(code uint8, buf []uint8, len int, enc bool) {
	if buf[0] != 0xC1 {
		// 0x0075C7B3
	}
	// 0x09FB56EA 0x0075C7C3 0x0AD3CA7D
	subcode := buf[3]
	switch subcode {
	case 0:
		// 0x0075C97D
		f006C780DhandleF300(buf) // character list
	case 1:
		// 0x0075C98B 0x0075C90D 0x09FBC920 0x0075C98E
		f006C798BhandleF301(buf) // character create
	case 2:
		// 0x0075C999
	case 3: // character info
		// 0x0075C9A7, s9 0x00670E5A
		f00701DF4handleF303(buf, enc) // character info
	case 4:
		// 0x0075C9C4
	case 5:
		// 0x0075C9EF
	case 6:
		// 0x0075C9FD
		f006C03F3handleF306(buf) // character add point
	case 0x13:
		// 0x0075CA35
	case 0x14:
		// 0x0075CA43
	case 0x15:
		// 0x0075CAE9
	case 0x20:
		// 0x0075CA51
	case 0x22:
		// 0x0075CA6D
	case 0x23:
		// 0x0075CA91
	case 0x24:
		// 0x0075CA7B
	case 0x25:
		// 0x0075CA86
	case 0x26:
		// 0x0075CA5F
	case 0x30:
		// 0x0075CA9C
	case 0x32:
		// 0x0075CAA7
	case 0x40:
		// 0x0075CAB2
	case 0x50:
		// 0x0075CABD
	case 0x51:
		// 0x0075CAC8
	case 0x52:
		// 0x0075CAD3
	case 0x53:
		// 0x0075CADE
	}
}

func f0075CB02handleF4(code uint8, buf []uint8, len int, enc bool) {
	subcode := buf[3]
	switch subcode {
	case 3: // server info
	// close连接服务器
	// f006BF89ADial(ip, port) // 拨号游戏服务器
	// set state
	case 5:
	case 6: // server list
		// f006BFA3E(buf)
		func(buf []uint8) {
			// 0x0A1B0A0B
			// push 0x00591FC5 // lable1 <- 0x09FC87EB
			// push 0x0AFD9E19
			// ret

			// 0x0AFD9E19
			// push 0x006A7909 // label2
			// push 0x0A9361E5
			// ret

			// 0x0A9361E5
			// push 0x00964365 // label3
			// push 0x0A43B202
			// ret

			// 0x0A43B202
			// push ecx edx edi ebx esi fd
			// 0x0ABE3CDB
			// push 0x09FD8750
			// push 0x09FD8029
			// ret

			// 0x09FD8029 f09FD8029
			func() {
				// ...
			}()
			// ...
			// 0x09FC87EB 0x0AF82CE7
			// 0x1C局部变量
			// ebp10buf := buf
			ebp4 := 5 // sizeof(frameHead)
			// f00AF7DC3serverList().f00AF7E20()
			count := int(binary.BigEndian.Uint16(buf[ebp4:]))
			ebp4 += 2
			f00AF7DC3serverList().f00AF883FsetCount(count)
			ebp14 := 0
			for {
				// 0x006BFAAD
				if ebp14 >= f00AF7DC3serverList().f00AF8853getCount() {
					break // 0x006BFAED
				}
				ebp18server := buf[ebp4:]
				code := int(binary.LittleEndian.Uint16(ebp18server[:]))
				percent := int(ebp18server[2])
				f00AF7DC3serverList().f00AF81FF(code, percent) // so, what about the server.type field?
				ebp4 += 4
				ebp14++
			}
			// 0x006BFAED
			ebp8 := f004A7D34getWindowManager()
			ebp19 := ebp8.s2.window.m0Cshow
			if ebp19 == false {
				// 0x0A441368 呈现
				ebp8.f004A9123(ebp8.s4)
				// ebp8.s4.f00446625()
				ebp8.f004A9123(ebp8.s5)
			}
			// 0x006BFB39 0x09FC6DF2
			v01319E08log.f00B38AE4printf("Success Receive Server List.\r\n")
		}(buf)
	}
}

func handleFA(code uint8, buf []uint8, len int, enc bool) {}
func f006C18B6handleF6(code uint8, buf []uint8, len int, enc bool) {
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
			ebp148Ctrap.f0043922CwriteHead(0xC1, 0xF3)
			ebp148Ctrap.f004397B1writeUint8(0x31)
			// big endian
			ebp148Ctrap.f004397B1writeUint8(uint8(unk >> 24))
			ebp148Ctrap.f004397B1writeUint8(uint8(unk >> 16))
			ebp148Ctrap.f004397B1writeUint8(uint8(unk >> 8))
			ebp148Ctrap.f004397B1writeUint8(uint8(unk))
			ebp148Ctrap.f004393EAsend(true, false)
			ebp148Ctrap.f004391CF()

			for {
				ebp1490 := 0
				if ebp1490 >= 0x190 {
					break // 0x007C8001
				}
				if ebp1490 == v08C88CACid {
					break
				}
				ebp1494obj := f004373C5objectPool().f00A38D5BgetObject(ebp1490)
				if ebp1494obj.m5Eid == 0 {
					break
				}
				ebp1490++
			}
			// ebp1494obj.m404 = 0x2710
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

func handleViewportGensInfoAdd(code uint8, buf []uint8, len int, enc bool) {}

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
