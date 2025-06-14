package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server-connect/conf"
	"github.com/xujintao/balgass/src/server-connect/handle"
	"github.com/xujintao/balgass/src/server-connect/service"
)

func main() {
	// handle signal and error
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	errChan := make(chan error, 2)

	// start service
	log.Println("start service")
	service.Service.Start()

	// start tcp server
	log.Println("start tcp server")
	tcpServer := c1c2.Server{
		Addr:    fmt.Sprintf(":%d", conf.Net.TCPPort),
		Handler: &handle.C1C2Handle,
		NeedXor: false,
	}
	go func() {
		errChan <- tcpServer.ListenAndServe()
	}()

	// start udp server
	log.Println("start udp server")
	udpServer := c1c2.ServerUDP{
		Addr:    fmt.Sprintf(":%d", conf.Net.UDPPort),
		Handler: &handle.C1C2Handle,
	}
	go func() {
		errChan <- udpServer.Run()
	}()

	// wait signal and error
	select {
	case s := <-exit:
		log.Printf("exit [signal]%s\n", s.String())
	case err := <-errChan:
		log.Fatalf("server failed: [err]%v\n", err)
	}

	// close tcp server
	log.Println("close tcp server")
	tcpServer.Close()

	// close udp server
	log.Println("close udp server")
	udpServer.Close()

	// close service
	log.Println("close service")
	service.Service.Close()

	time.Sleep(2 * time.Second)
}
