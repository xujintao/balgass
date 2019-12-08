package model

import (
	"github.com/xujintao/balgass/network"
)

type Server struct {
	Addr string
	Conn network.ConnWriter
	Type int32
	Port int32
	Name string
	Code int32
	Vip  int32
}
