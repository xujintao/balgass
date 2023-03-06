package c1c2

import (
	"log"
	"net"
)

// UDPHandler callback user to handle request
type UDPHandler interface {
	HandleUDP(*Request, *Response) bool
}

type connUDP struct {
	server *ServerUDP
	rw     *net.UDPConn
	addr   *net.UDPAddr
	bufr   []byte
}

func (c *connUDP) serve() {
	for {
		n, addr, err := c.rw.ReadFromUDP(c.bufr)
		if err != nil {
			log.Println(err)
			continue
		}
		c.addr = addr

		// parse
		req, err := parseFrame(c.bufr[:n], false)
		if err != nil {
			log.Println(err)
			continue
		}

		// callback
		res := new(Response)
		if c.server.Handler.HandleUDP(req, res) {
			if err := c.write(res); err != nil {
				log.Printf("Write failed, %v", err)
			}
		}
	}
}

// Write
func (c *connUDP) write(r *Response) error {
	frame, err := newFrame(r)
	if err != nil {
		return err
	}

	if _, err := c.rw.WriteToUDP(frame, c.addr); err != nil {
		return err
	}
	return nil
}

// ServerUDP implements the specified protocol
type ServerUDP struct {
	Addr    string
	Handler UDPHandler
}

// Run on the UDP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
func (srv *ServerUDP) Run() error {
	addr := srv.Addr
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	rw, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}
	const size = 4 << 10
	c := &connUDP{
		server: srv,
		rw:     rw,
		bufr:   make([]byte, size),
	}

	c.serve()
	return nil
}

func (srv *ServerUDP) Close() {

}

// Run on the UDP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
func Run(addr string, handler UDPHandler) error {
	server := &ServerUDP{Addr: addr, Handler: handler}
	return server.Run()
}
