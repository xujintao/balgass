package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/handle"
)

func main() {
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// start tcp server
	log.Printf("start tcp server")
	server := c1c2.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.GameServerInfo.Port),
		Handler: &handle.APIHandleDefault,
		NeedXor: true,
	}
	go func() {
		err := server.ListenAndServe()
		log.Fatal(err)
	}()

	// wait
	<-exit
	log.Println("SIGINT or SIGTERM")

	// close tcp server
	log.Printf("close tcp server")
	server.Close()

	// close handle
	handle.APIHandleDefault.Close()
	time.Sleep(2 * time.Second)
}
