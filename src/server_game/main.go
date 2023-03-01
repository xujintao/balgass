package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xujintao/balgass/cmd/server_game/conf"
	"github.com/xujintao/balgass/cmd/server_game/handle"
	"github.com/xujintao/balgass/network"
)

func main() {
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// start tcp server
	server := network.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: &handle.APIHandleDefault,
		NeedXor: true,
	}
	log.Printf("start tcp server")
	go func() {
		err := server.ListenAndServe()
		log.Fatal(err)
	}()

	<-exit
	server.Close()
	time.Sleep(2 * time.Second)
}
