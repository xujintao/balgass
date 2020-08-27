package handle

import (
	"github.com/xujintao/balgass/network"
)

type handleExData struct {
	handleBase
}

var cmdsExData = map[int]func(index interface{}, req *network.Request){
	// 0x05:   checkVersion,
	// 0xf406: getServerList,
	// 0xf403: getServerInfo,
}
