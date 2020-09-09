package main

import (
		"bytes"
	//	"fmt"
		"time"
		"net"
		"errors"
	   )


type Request struct {
	data bytes.Buffer
}

func (rq *Request)readRequest(c conn) (int, error) {
	rq.data.Reset()
	c.rwc.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))

	// check circle buffer
	if ok := c.rbuf.ReadBlock(&rq.data); ok {
		return rq.data.Len(), nil
	}

	b := make([]byte, 100)
	n, err := c.rwc.Read(b)
	if err != nil {
		//client time out
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			return 0, nil
		}
		//fmt.Printf("conn read %d bytes,  error: %s, serve client QUIT.....\n", n, err)
		return 0, errors.New("Client connection failed.")
	}
	// append to circle buffer
	c.rbuf.Write(b[:n])

	// check circle buffer
	if ok := c.rbuf.ReadBlock(&rq.data); ok {
		return rq.data.Len(), nil
	}

	return 0, nil
}
