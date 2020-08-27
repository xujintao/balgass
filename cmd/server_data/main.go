package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xujintao/balgass/cmd/server_data/conf"
	"github.com/xujintao/balgass/cmd/server_data/handle"
	"github.com/xujintao/balgass/network"
)

func main() {
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	type config struct {
		name string
		port int
		use  bool
		only bool
		h    network.Handler
	}
	configs := [...]config{
		{"join server", conf.Server.JoinServerPort, conf.Server.UseJoinServer, conf.Server.DataServerOnlyMode, handle.HandleJoin},
		{"data server", conf.Server.DataServerPort, conf.Server.UseDataServer, false, handle.HandleData},
		{"exdata server", conf.Server.ExDataServerPort, conf.Server.UseExDataServer, conf.Server.DataServerOnlyMode, handle.HandleExData},
	}
	var servers []*network.Server
	for _, c := range configs {
		if c.use && !c.only {
			log.Printf("start %s", c.name)
			server := &network.Server{
				Addr:    fmt.Sprintf(":%d", c.port),
				Handler: c.h,
			}
			servers = append(servers, server)
			go func() {
				err := server.ListenAndServe()
				log.Fatal(err)
			}()
		}
	}

	// start udp server
	// log.Printf("start udp server")
	// serverUDP := network.ServerUDP{
	// 	Addr:    fmt.Sprintf(":%d", conf.Net.UDPPort),
	// 	Handler: ch,
	// }
	// go func() {
	// 	err := serverUDP.Run()
	// 	log.Fatal(err)
	// }()

	<-exit
	for _, s := range servers {
		s.Close()
	}
	handle.HandleData.Exit()
	time.Sleep(5 * time.Second)
}
