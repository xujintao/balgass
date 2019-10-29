package wzudp

import (
	"log"
	"net"
)

// Handler callback handle logic
type Handler interface {
	Handle(req *Message) *Message
}

type conn struct {
	server *Server
	rw     *net.UDPConn
	addr   *net.UDPAddr
	bufr   []byte
}

func (c *conn) serve() {
	for {
		n, addr, err := c.rw.ReadFromUDP(c.bufr[:0])
		if err != nil {
			log.Fatal(err)
		}
		c.addr = addr
		req, _ := parseFrame(c.bufr[:n])
		res := c.server.Handler.Handle(req)
		if res == nil {
			continue
		}
		frame := newFrame(res)
		_, err = c.rw.WriteToUDP(frame, c.addr)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Server implements the specified protocol
type Server struct {
	Addr    string
	Handler Handler
}

// Run on the UDP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
func (srv *Server) Run() error {
	addr := srv.Addr
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	rw, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	c := &conn{
		server: srv,
		rw:     rw,
		bufr:   make([]byte, 4<<10),
	}

	c.serve()
	return nil
}

// Run on the UDP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
func Run(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.Run()
}
