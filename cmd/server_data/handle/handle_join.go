package handle

import (
	"github.com/xujintao/balgass/network"
)

type handleJoin struct {
	handleBase
}

var cmdsJoin = map[int]func(index interface{}, req *network.Request){
	0x00: serverLogin,
	0x01: userLogin,
	0x02: userLoginFailed,
	0x04: userBlock,
	0x05: userExit,
	0x06: vipCheck,
	// 0x30: loveHeartEventReq,
	// 0x31: loveHeartCreate,
	0x7A: userServerMove,
	0x7B: userServerMoveAuth,
	0x7C: serverUserCountSet,
	0x80: userOffTradeSet,
	0xD3: vipAdd,
	0xFE: userKill,
	// 0xFF: serverDisconnect, // service.ServerManager.ServerExit(index) instead
}
