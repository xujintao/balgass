package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/xujintao/balgass/cmd/server_connect/conf"
	"github.com/xujintao/balgass/cmd/server_connect/handle"
	"github.com/xujintao/balgass/protocol"
)

func main() {
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	ch := handle.CMDHandle{}

	// start tcp server
	server := protocol.Server{
		Addr:      fmt.Sprintf(":%d", conf.Net.TCPPort),
		Handler:   ch,
		NeedXor:   false,
		ConnState: ch.TrackConnState,
		OnConn:    ch.OnConn,
	}
	log.Printf("start tcp server")
	go func() {
		err := server.ListenAndServe()
		log.Fatal(err)
	}()

	// start udp server
	log.Printf("start udp server")
	serverUDP := protocol.ServerUDP{
		Addr:    fmt.Sprintf(":%d", conf.Net.UDPPort),
		Handler: ch,
	}
	go func() {
		err := serverUDP.Run()
		log.Fatal(err)
	}()

	<-exit
	server.Close()
}
