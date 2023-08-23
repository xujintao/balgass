package game

import (
	"context"
	"fmt"
	"log"
)

func (g *game) addBot(name string) error {
	b, err := newbot(name)
	if err != nil {
		return err
	}
	g.bots[b.name] = b
	return nil
}

func (g *game) deleteBot(name string) {
	if b, ok := g.bots[name]; ok {
		b.Close()
		delete(g.bots, name)
	}
}

func (g *game) deleteAllBots() {
	for name, b := range g.bots {
		b.Close()
		delete(g.bots, name)
	}
}

func newbot(name string) (*bot, error) {
	b := &bot{name: name}
	b.msgChan = make(chan any, 100)
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	go func() {
		defer close(b.msgChan)
		id, err := Game.PlayerConn(b)
		if err != nil {
			return
		}
		defer Game.PlayerCloseConn(id)
		for {
			select {
			case msg := <-b.msgChan:
				fmt.Println(msg)
			case <-ctx.Done():
				return
			}
		}
	}()
	return b, nil
}

type bot struct {
	name    string
	cancel  context.CancelFunc
	msgChan chan any
}

func (b *bot) Addr() string {
	return b.name
}

func (b *bot) Write(msg any) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("*bot.Write [panic]%v [bot]%s [msg]%v\n", r, b.name, msg)
		}
	}()
	b.msgChan <- msg
	return nil
}

func (b *bot) Close() error {
	b.cancel()
	return nil
}
