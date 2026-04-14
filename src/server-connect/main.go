package main

import (
	"fmt"
	"log/slog"
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
	slog.Info("start service")
	service.Service.Start()

	// start tcp server
	slog.Info("start tcp(c1c2) server")
	tcpServer := c1c2.Server{
		Addr:    fmt.Sprintf(":%d", conf.Net.TCPPort),
		Handler: &handle.C1C2Handle,
		NeedXor: false,
	}
	go func() {
		errChan <- tcpServer.ListenAndServe()
	}()

	// start udp server
	slog.Info("start udp(c1c2) server")
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
		slog.Info("exit", "signal", s.String())
	case err := <-errChan:
		slog.Error("<-errChan", "err", err)
		os.Exit(1)
	}

	// close tcp server
	slog.Info("close tcp(c1c2) server")
	tcpServer.Close()

	// close udp server
	slog.Info("close udp(c1c2) server")
	udpServer.Close()

	// close service
	slog.Info("close service")
	service.Service.Close()

	time.Sleep(2 * time.Second)
}
