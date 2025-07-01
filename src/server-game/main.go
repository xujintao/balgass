package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xujintao/balgass/src/c1c2"
	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game"
	"github.com/xujintao/balgass/src/server-game/handle"
)

func main() {
	// handle signal and error
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	errChan := make(chan error, 2)

	// start game
	slog.Info("start game")
	game.Game.Start()

	// start tcp server
	slog.Info("start tcp server")
	tcpServer := c1c2.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.GameServerInfo.Port),
		Handler: &handle.C1C2Handle,
		NeedXor: true,
	}
	go func() {
		errChan <- tcpServer.ListenAndServe()
	}()

	// start http server
	slog.Info("start http server")
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
		slog.Info("exit", "signal", s.String())
	case err := <-errChan:
		slog.Error("<-errChan", "err", err)
		os.Exit(1)
	}

	// close tcp server
	slog.Info("close tcp server")
	tcpServer.Close()

	// close http server
	slog.Info("close http server")
	if err := httpServer.Close(); err != nil {
		slog.Error("httpServer.Close()", "err", err)
	}

	// close game
	slog.Info("close game")
	game.Game.Close()

	time.Sleep(2 * time.Second)
}
