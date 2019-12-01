package protocol

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
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

// Handler callback user to handle request
type Handler interface {
	Handle(*Request, *Response) bool
	HandleUDP(*Request, *Response) bool
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

// A ConnState represents the state of a client connection to a server.
type ConnState int

const (
	// StateNew represents a new connection that is expected to
	// send a request immediately. Connections begin at this
	// state and then transition to either StateActive or
	// StateClosed.
	StateNew ConnState = iota

	// StateActive represents a connection that has read 1 or more
	// bytes of a request.
	StateActive

	// StateIdle represents a connection that has finished
	// handling a request and is in the keep-alive state, waiting
	// for a new request. Connections transition from StateIdle
	// to either StateActive or StateClosed.
	StateIdle

	// StateClosed represents a closed connection.
	// This is a terminal state. Hijacked connections do not
	// transition to StateClosed.
	StateClosed
)

// A conn represents the server side of an connection.
type conn struct {
	// server is the server on which the connection arrived.
	// Immutable; never nil.
	server *Server

	// cancelCtx cancels the connection-level context.
	cancelCtx context.CancelFunc

	// rwc is the underlying network connection.
	// This is never wrapped by other types and is the value given out
	// to CloseNotifier callers. It is usually of type *net.TCPConn or
	// *tls.Conn.
	rwc net.Conn

	// remoteAddr is rwc.RemoteAddr().String(). It is not populated synchronously
	// inside the Listener's Accept goroutine, as some implementations block.
	// It is populated immediately inside the (*conn).serve goroutine.
	// This is the value of a Handler's (*Request).RemoteAddr.
	remoteAddr string

	bufr *bufio.Reader
	bufw *bufio.Writer

	curState struct{ atomic uint64 } // packed (unixtime<<8|uint8(ConnState))
}

func (c *conn) setState(nc net.Conn, state ConnState) {
	srv := c.server
	switch state {
	case StateNew:
		srv.trackConn(c, true)
	case StateClosed:
		srv.trackConn(c, false)
	}
	if state > 0xff || state < 0 {
		panic("internal error")
	}
	packedState := uint64(time.Now().Unix()<<8) | uint64(state)
	atomic.StoreUint64(&c.curState.atomic, packedState)
	if hook := srv.ConnState; hook != nil {
		hook(nc, state)
	}
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

func (c *conn) readRequest(ctx context.Context) (req *Request, err error) {
	//  read deadline
	if d := c.server.ReadTimeout; d != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(d))
	}

	// write deadline
	if d := c.server.WriteTimeout; d != 0 {
		defer func() {
			c.rwc.SetWriteDeadline(time.Now().Add(d))
		}()
	}

	// read
	frame, err := readFrame(c.bufr)
	if err != nil {
		return nil, fmt.Errorf("readFrame failed, %v", err)
	}

	// parse
	req, err = parseFrame(frame, c.server.NeedXor)

	return
}

func (c *conn) serve(ctx context.Context) {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("panic serving %v: %v, %s", c.remoteAddr, err, buf)
		}
		c.close()
		c.setState(c.rwc, StateClosed)
	}()

	ctx, cancelCtx := context.WithCancel(ctx)
	c.cancelCtx = cancelCtx
	defer cancelCtx()

	c.bufr = newBufioReader(c.rwc)
	c.bufw = newBufioWriteSize(c.rwc, 4<<10)

	// callback OnConn
	if hook := c.server.OnConn; hook != nil {
		res := new(Response)
		hook(nil, res)
		if err := c.write(res); err != nil {
			log.Printf("%s, Write failed, %v", c.remoteAddr, err)
		}
	}

	for {
		req, err := c.readRequest(ctx)
		if err != nil {
			log.Printf("%s, readRequest failed, %v", c.remoteAddr, err)
			return
		}

		// callback Handle
		res := new(Response)
		if c.server.Handler.Handle(req, res) {
			if err := c.write(res); err != nil {
				log.Printf("%s, Write failed, %v", c.remoteAddr, err)
			}
		}
	}
}

// write
func (c *conn) write(r *Response) error {
	frame, err := newFrame(r)
	if err != nil {
		return err
	}
	if _, err := c.bufw.Write(frame); err != nil {
		return err
	}
	return c.bufw.Flush()
}

// Server implements the specified protocol
type Server struct {
	Addr    string
	Handler Handler
	NeedXor bool

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	WriteTimeout time.Duration

	// ConnState specifies an optional callback function that is
	// called when a client connection changes state. See the
	// ConnState type and associated constants for details.
	ConnState func(net.Conn, ConnState)

	OnConn func(req *Request, res *Response) bool

	mu          sync.Mutex
	doneChan    chan struct{}
	listeners   map[net.Listener]struct{}
	activeConns map[*conn]struct{}
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// onceCloseListener wraps a net.Listener, protecting it from
// multiple Close calls.
type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(func() { oc.closeErr = oc.Listener.Close() })
	return oc.closeErr
}

func (s *Server) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getDoneChanLocked()
}

func (s *Server) getDoneChanLocked() chan struct{} {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

// ErrServerClosed is returned by the Server's Serve, ServeTLS, ListenAndServe,
// and ListenAndServeTLS methods after a call to Shutdown or Close.
var ErrServerClosed = errors.New("http: Server closed")

// Close immediately closes all active net.Listeners and any
// connections in state StateNew, StateActive, or StateIdle.
func (srv *Server) Close() {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	// close chan
	ch := srv.getDoneChanLocked()
	select {
	case <-ch:
	default:
		close(ch)
	}

	// close listen
	var err error
	for ln := range srv.listeners {
		if e := ln.Close(); e != nil && err == nil {
			err = e
		}
		delete(srv.listeners, ln)
	}

	// close all the conn
	for c := range srv.activeConns {
		c.rwc.Close()
	}
}

func (s *Server) trackConn(c *conn, add bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeConns == nil {
		s.activeConns = make(map[*conn]struct{})
	}
	if add {
		s.activeConns[c] = struct{}{}
	} else {
		delete(s.activeConns, c)
	}
}

// trackListener adds or removes a net.Listener to the set of tracked
// listeners.
//
// We store a pointer to interface in the map set, in case the
// net.Listener is not comparable. This is safe because we only call
// trackListener via Serve and can track+defer untrack the same
// pointer to local variable there. We never need to compare a
// Listener from another caller.
//
// It reports whether the server is still up (not Shutdown or Closed).
func (s *Server) trackListener(ln net.Listener, add bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listeners == nil {
		s.listeners = make(map[net.Listener]struct{})
	}
	if add {
		s.listeners[ln] = struct{}{}
	} else {
		delete(s.listeners, ln)
	}
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each. The service goroutines read requests and
// then call srv.Handler to reply to them.
func (srv *Server) Serve(l net.Listener) error {
	l = &onceCloseListener{Listener: l}
	defer l.Close()

	srv.trackListener(l, true)
	defer srv.trackListener(l, false)

	var tempDelay time.Duration // how long to sleep on accept failure
	baseCtx := context.Background()
	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
		rw, err := l.Accept()
		if err != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
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
		c.setState(c.rwc, StateNew) // before Serve can return
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
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

// ListenAndServe listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
func ListenAndServe(addr string, handler Handler, needxor bool) error {
	server := &Server{Addr: addr, Handler: handler, NeedXor: needxor}
	return server.ListenAndServe()
}
