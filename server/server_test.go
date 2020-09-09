package main

import (
		"fmt"
		"net"
		"time"
		"bytes"
		"strings"
		"testing"
	   )

var y Server
func TestAccept(t *testing.T) {
	var e echo
	y.AddHandleFunc(0x01, e)
	go y.Start()


	// Now, support 500 clients at the same time
	c := make([]net.Conn, 500)
	time.Sleep(time.Second)
	for i := 0; i < 500; i++ {
		var err error
		c[i], err = net.Dial("tcp4", "127.0.0.1:7890")
		if err != nil {
			t.Errorf("dial failed: %v\n", err)
			return
		}

		var bw bytes.Buffer

		for j := 0; j < 5; j++ {
			bw.Reset()
			bw.WriteByte(0xAA)
			bw.WriteByte(0x55)
			bw.WriteByte(0x00)
			bw.WriteByte(0x00)
			bw.WriteByte(0x00)
			bw.WriteByte(0x09)
			bw.Write([]byte("hello,jkl"))
			//c[i].Write([]byte("hello,"))
			bw.WriteTo(c[i])
			p := make([]byte,100)
			c[i].SetReadDeadline(time.Now().Add(10000*time.Millisecond))
			_, err := c[i].Read(p)
			if err == nil {
				fmt.Printf("client recv: %s [%d]\n", string(p), j)
				if strings.Compare(string(p[:9]), "hello,jkl") != 0 {
					t.Errorf("not equal")
				}
			}
		}
		c[i].Close()
	}

	fmt.Println("sleep 3 seconds, then close server")
	time.Sleep(3*time.Second)
	y.Close()

}
