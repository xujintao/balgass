package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xujintao/balgass/cmd/server_connect/conf"
	"github.com/xujintao/balgass/cmd/server_connect/handle"
	"github.com/xujintao/balgass/network"
)

func main() {
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	ch := handle.CMDHandle{}

	// start tcp server
	server := network.Server{
		Addr:    fmt.Sprintf(":%d", conf.Net.TCPPort),
		Handler: ch,
		NeedXor: false,
	}
	log.Printf("start tcp server")
	go func() {
		err := server.ListenAndServe()
		log.Fatal(err)
	}()

	// start udp server
	log.Printf("start udp server")
	serverUDP := network.ServerUDP{
		Addr:    fmt.Sprintf(":%d", conf.Net.UDPPort),
		Handler: ch,
	}
	go func() {
		err := serverUDP.Run()
		log.Fatal(err)
	}()

	<-exit
	server.Close()
	time.Sleep(2 * time.Second)
}
