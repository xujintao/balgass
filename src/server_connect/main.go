package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_connect/conf"
	"github.com/xujintao/balgass/src/server_connect/handle"
	"github.com/xujintao/balgass/src/server_connect/service"
)

func main() {
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// start service
	log.Println("start service")
	service.Service.Start()

	// start tcp server
	tcpServer := c1c2.Server{
		Addr:    fmt.Sprintf(":%d", conf.Net.TCPPort),
		Handler: &handle.C1C2Handle,
		NeedXor: false,
	}
	log.Println("start tcp server")
	go func() {
		err := tcpServer.ListenAndServe()
		if err != nil {
			log.Printf("tcpServer.ListenAndServe failed [err]%v\n", err)
		}
	}()

	// start udp server
	log.Println("start udp server")
	udpServer := c1c2.ServerUDP{
		Addr:    fmt.Sprintf(":%d", conf.Net.UDPPort),
		Handler: &handle.C1C2Handle,
	}
	go func() {
		err := udpServer.Run()
		if err != nil {
			log.Printf("udpServer.Run failed [err]%v\n", err)
		}
	}()

	// wait
	s := <-exit
	log.Printf("exit [signal]%s\n", s.String())

	// close service
	log.Println("close service")
	service.Service.Close()

	// close tcp server
	log.Println("close tcp server")
	tcpServer.Close()

	// close udp server
	udpServer.Close()

	time.Sleep(2 * time.Second)
}
