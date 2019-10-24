package main

import (
	"fmt"

	"github.com/xujintao/balgass/server/protocol"
)

func main() {
	addr := fmt.Sprintf(":%d", Config.Net.TCPPort)
	protocol.ListenAndServe(addr, cmdHandle{}, false)
}
