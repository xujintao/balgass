package model

import "github.com/xujintao/balgass/network"

type User struct {
	Addr string
	Conn network.ConnWriter
	New  bool
}
