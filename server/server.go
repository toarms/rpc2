package main

import (
		"fmt"
		"sync"
		"time"
		"net"
		"log"
		"io"
		"bytes"
		"errors"
		"private.github.com/toarms/rpc1/yuan"
		"private.github.com/toarms/rpc1/blockbuf"
	   )

// HandleFunc
type HandleFunc interface {
	handle(w io.Writer, rq Request)
}

//conn
type conn struct {
	Addr string
	rwc net.Conn
	srv *Server
	rbuf yuan.YBuf
	wbuf yuan.YBuf
}

// Server
type Server struct {
	doneChan chan struct{}
	mu sync.Mutex
	ln net.Listener
	activeConn map[*conn]struct{}
	h map[int]HandleFunc
	nclients int
}

// doneChan
func (srv *Server)getDoneChan() <-chan struct{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.getDoneChanLocked()
}
func (srv *Server)getDoneChanLocked() chan struct{} {
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	return srv.doneChan
}
func (srv *Server)closeDoneChan() {
	ch := srv.getDoneChanLocked()
	select {
		case <-ch:
			// already closed, don't close again
		default:
			// safe to close here. we are the only closer, guarded
			// by srv.mu
			close(ch)
	}
}

func (srv *Server)add() {
	fmt.Println("this is add")
}

func (srv *Server)AddHandleFunc(servid int, h HandleFunc) {

	if srv.h == nil {
		srv.h = make(map[int]HandleFunc)
	}
	srv.h[servid] = h
}

func (srv *Server)Start() error {
	fmt.Println("Start...")

	var err error
	srv.ln, err = net.Listen("tcp4", "127.0.0.1:7890")
	if err != nil {
		log.Println("start failed: ", err)
		return nil
	}

	for {

		client, err := srv.ln.Accept()
		if err != nil {
			select {
				case <-srv.getDoneChan():
					fmt.Println("get done chan, server exit")
					return errors.New("Server closed")
				default:
					fmt.Println("default")
			}
			return errors.New("Accept")
		}
		c := srv.newConn(client)
		srv.trackConn(c, true)
		go c.serve()
	}

	fmt.Println("Start quit...")

	return nil
}

func (srv *Server)Close() {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	fmt.Println("to Close...")
	srv.closeDoneChan()
	srv.ln.Close()
	for c := range srv.activeConn {
		srv.trackConn(c, false)
		c.rwc.Close()
	}
}

func (srv *Server)newConn(rwc net.Conn) *conn {
	yr := yuan.YBuf{}
	yr.Init(1024*1024)
	yw := yuan.YBuf{}
	yw.Init(1024*1024)
	c := &conn {
			Addr:rwc.RemoteAddr().String(),
			rwc: rwc,
			srv: srv,
			rbuf: yr,
			wbuf: yw,
	   }

   return c
}

func (srv *Server)trackConn(c *conn, add bool) {
	if srv.activeConn == nil {
		srv.activeConn = make(map[*conn]struct{})
	}

	if add {
		srv.activeConn[c] = struct{}{}
		srv.nclients++
		fmt.Printf("new client: %s [%d]\n", c.Addr, srv.nclients)
	} else {
		fmt.Printf("server close client %s\n", c.Addr)
		delete(srv.activeConn, c)
		srv.nclients--
	}
}


// conn
func (c *conn)serve() {
	fmt.Println("new connect handle.")
	var rq Request
	var bb bytes.Buffer
	bf := blockbuf.New(c.rwc, c.rwc)
	for {
		bb.Reset()
		bf.ReadBlock(&bb)
		bf.WriteBlock([]byte("hello"))

		c.srv.h[0x01].handle(c.rwc, rq)
	}
}

func (c *conn)ReadBlock(b *bytes.Buffer) bool {
	return false
}
func (c *conn)WriteBlock(b []byte) {
}


// handle
type echo struct {}
func (e echo)handle(w io.Writer, rq Request) {
		rq.data.WriteTo(w)
}

func main() {
	var s Server
	var e echo
	s.AddHandleFunc(0x01, e)
	go s.Start()

	time.Sleep(3*time.Second)

	s.Close()

	for {
	}
}
