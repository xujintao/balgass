package protocol

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var (
	// ServerContextKey is a context key.
	ServerContextKey = &contextKey{"go-server"}

	// LocalAddrContextKey is a context key.
	LocalAddrContextKey = &contextKey{"local-addr"}
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "net/http context value " + k.name }

// Handler callback handle logic
type Handler interface {
	Handle(req *Message) *Message
}

// pool
var (
	bufioReaderPool   sync.Pool
	bufioWriter2kPool sync.Pool
	bufioWriter4kPool sync.Pool
)

func bufioWriterPool(size int) *sync.Pool {
	switch size {
	case 2 << 10:
		return &bufioWriter2kPool
	case 4 << 10:
		return &bufioWriter4kPool
	}
	return nil
}

func newBufioReader(r io.Reader) *bufio.Reader {
	if v := bufioReaderPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReader(r)
}

func putBufioReader(br *bufio.Reader) {
	br.Reset(nil)
	bufioReaderPool.Put(br)
}

func newBufioWriteSize(w io.Writer, size int) *bufio.Writer {
	pool := bufioWriterPool(size)
	if pool != nil {
		if v := pool.Get(); v != nil {
			bw := v.(*bufio.Writer)
			bw.Reset(w)
			return bw
		}
	}
	return bufio.NewWriterSize(w, size)
}

func putBufioWriter(bw *bufio.Writer) {
	bw.Reset(nil)
	if pool := bufioWriterPool(bw.Available()); pool != nil {
		pool.Put(bw)
	}
}

type conn struct {
	server     *Server
	cancelCtx  context.CancelFunc
	rwc        net.Conn
	remoteAddr string
	bufr       *bufio.Reader
	bufw       *bufio.Writer
}

func (c *conn) finalFlush() {
	if c.bufr != nil {
		putBufioReader(c.bufr)
		c.bufr = nil
	}
	if c.bufw != nil {
		c.bufw.Flush()
		putBufioWriter(c.bufw)
		c.bufw = nil
	}
}

func (c *conn) close() {
	c.finalFlush()
	c.rwc.Close()
}

func (c *conn) readRequest(ctx context.Context) (*Message, error) {
	// peek 3 bytes
	frameHead, err := c.bufr.Peek(3)
	if err != nil {
		return nil, err
	}

	flag := frameHead[0]
	size := 0
	switch flag {
	case 0xC1, 0xC3:
		size = int(frameHead[1])
	case 0xC2, 0xC4:
		size = int(binary.BigEndian.Uint16(frameHead[1:]))
	default:
		return nil, fmt.Errorf("invalid flag: %x", flag)
	}

	// peek size bytes
	if _, err := c.bufr.Peek(size); err != nil {
		return nil, err
	}

	// read
	frame := make([]byte, size)
	if _, err := c.bufr.Read(frame); err != nil {
		return nil, err
	}

	// parse
	return parseFrame(frame, c.server.isGame)
}

func (c *conn) serve(ctx context.Context) {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	defer c.close()

	ctx, cancelCtx := context.WithCancel(ctx)
	c.cancelCtx = cancelCtx
	defer cancelCtx()

	c.bufr = newBufioReader(c.rwc)
	c.bufw = newBufioWriteSize(c.rwc, 4<<10)
	for {
		req, err := c.readRequest(ctx)
		if err != nil {
			log.Println(err)
		}
		res := c.server.Handler.Handle(req) // callback
		resFrame := createFrame(res)
		if _, err = c.bufw.Write(resFrame); err != nil {
			log.Println(err)
		}
		if err := c.bufw.Flush(); err != nil {
			log.Println(err)
		}
	}
}

// Server implements the specified protocol
type Server struct {
	Addr    string
	Handler Handler
	isGame  bool
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each. The service goroutines read requests and
// then call srv.Handler to reply to them.
func (srv *Server) Serve(l net.Listener) error {
	var tempDelay time.Duration
	baseCtx := context.Background()
	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
		rw, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("server: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0
		c := &conn{
			server: srv,
			rwc:    rw,
		}
		go c.serve(ctx)
	}
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

// ListenAndServe listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
func ListenAndServe(addr string, handler Handler, isGame bool) error {
	server := &Server{Addr: addr, Handler: handler, isGame: isGame}
	return server.ListenAndServe()
}
