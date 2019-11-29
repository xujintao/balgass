package main

import (
	"fmt"

	"github.com/xujintao/balgass/cmd/server_connect/conf"
	"github.com/xujintao/balgass/cmd/server_connect/handle"
	"github.com/xujintao/balgass/protocol"
	"github.com/xujintao/balgass/wzudp"
)

func main() {
	// start udp server
	udpAddr := fmt.Sprintf(":%d", conf.Net.UDPPort)
	go wzudp.Run(udpAddr, handle.UDPCMDHandle{})

	// start tcp server
	addr := fmt.Sprintf(":%d", conf.Net.TCPPort)
	protocol.ListenAndServe(addr, handle.CMDHandle{}, false)
}
