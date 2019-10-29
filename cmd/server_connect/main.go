package main

import (
	"fmt"

	"github.com/xujintao/balgass/wzudp"

	"github.com/xujintao/balgass/protocol"
)

func main() {
	// start udp server
	udpAddr := fmt.Sprintf(":%d", Config.Net.UDPPort)
	go wzudp.Run(udpAddr, udpcmdHandle{})

	// start tcp server
	addr := fmt.Sprintf(":%d", Config.Net.TCPPort)
	protocol.ListenAndServe(addr, cmdHandle{}, false)
}
