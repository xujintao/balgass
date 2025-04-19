package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game"
	"github.com/xujintao/balgass/src/server_game/handle"
)

func main() {
	// handle signal and error
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	errChan := make(chan error, 2)

	// start game
	log.Println("start game")
	game.Game.Start()

	// start tcp server
	log.Printf("start tcp server")
	tcpServer := c1c2.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.GameServerInfo.Port),
		Handler: &handle.C1C2Handle,
		NeedXor: true,
	}
	go func() {
		errChan <- tcpServer.ListenAndServe()
	}()

	// start http server
	log.Printf("start http server")
	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.GameServerInfo.HTTPPort),
		Handler: &handle.HTTPHandle,
	}
	go func() {
		errChan <- httpServer.ListenAndServe()
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

	// close http server
	log.Println("close http server")
	if err := httpServer.Close(); err != nil {
		log.Printf("httpServer.Close failed [err]%v\n", err)
	}

	// close game
	log.Println("close game")
	game.Game.Close()

	time.Sleep(2 * time.Second)
}
