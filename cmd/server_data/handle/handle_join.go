package handle

import (
	"github.com/xujintao/balgass/network"
)

type handleJoin struct {
	handleBase
}

var cmdsJoin = map[int]func(index string, req *network.Request){
	// 0x05:   checkVersion,
	// 0xf406: getServerList,
	// 0xf403: getServerInfo,
}
