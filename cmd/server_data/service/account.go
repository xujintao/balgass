package service

import "sync"

type account struct {
	username string
	machieID int
	number   int
	addr     string
}

type accountManager struct {
	mu      sync.RWMutex
	servers map[string]*account
}
