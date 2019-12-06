package handle

import (
	"github.com/xujintao/balgass/cmd/server_data/model"
	"github.com/xujintao/balgass/network"
)

type handleJoin struct {
	handleBase
}

var cmdsJoin = map[int]func(server *model.Server, req *network.Request){
	// 0x05:   checkVersion,
	// 0xf406: getServerList,
	// 0xf403: getServerInfo,
}
