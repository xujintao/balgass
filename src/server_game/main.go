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
	server := c1c2.Server{
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
