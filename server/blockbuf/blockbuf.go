package blockbuf

import (
		"io"
		"bytes"
		"errors"
		"net"
	   )

type Blockbuf struct {
	r io.Reader
	w io.Writer
	CBC
}

func New(w io.Writer, r io.Reader) (*Blockbuf){
	var bf Blockbuf

	bf.r = r
	bf.w = w
	bf.CBC.Init(1024*1024)

	return &bf
}

func (bf *Blockbuf)String() string {
	return "Blockbuf"
}

// ReadBlock
func (bf *Blockbuf)ReadBlock(bb *bytes.Buffer) error {
	//c.rwc.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))

	b := make([]byte, 100)
	n, err := bf.r.Read(b)
	if err != nil {
		//client time out
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			return nil
		}
		return errors.New("Client connection failed.")
	}
	// append to circle buffer
	bf.w.Write(b[:n])

	if !bf.ScanBlock() {
		return nil
	}

	bb.Reset()
	bb.Write(bf.firstBlockData())
	return nil
}

// Write
func (y *Blockbuf)WriteBlock(p []byte) (n int, err error) {
	left := y.Space()

	if left < len(p) {
		return 0, errors.New("no space")
	}

	var i int
	for i = 0; i < len(p); i++ {
		y.Push(p[i])
	}
	return i, nil
}

// Yuan protocol formtat:
//	| Magic(2) | dataLen(4) | Data(variable) |
const (
		Magic1		= 0xAA
		Magic2		= 0x55
		HeaderLen	= 6
	  )

// Header
type Header struct {
	magic1	byte
	magic2	byte
	datalen	uint32
}
// Block
type Block struct {
	Header
	data	[]byte
}

// Header
func newHeader(a,b,c,d,e,f byte) Header {
	var yh Header
	yh.magic1 = a
	yh.magic2 = b
	yh.datalen= uint32(c << 24) | uint32(d << 16) | uint32(e << 8) | uint32(f)
	return yh
}

func (y *Blockbuf)seekmagic() bool {
	var found int = -1
	var ret bool = false

	if y.Count() < 2 {
		return ret
	}
	for i := 0; i < y.Count() - 1; i++ {
		b := y.At(i)
		if b == Magic1 && y.At(i+1) == Magic2 {
			found = i
			break
		}
	}

	if found == -1 {
		if y.Last() == Magic1 {
			y.DropN(y.Count()-1)
		} else {
			y.Reset()
		}
	} else {
		if found > 0 {
			y.DropN(found)
		}
		ret = true
	}
	return ret
}

// is Complete YuanBlock
func (yb *Blockbuf)ScanBlock() bool {
	if ok := yb.seekmagic(); !ok {
		return false
	}

	if yb.Count() < HeaderLen {
		return false
	}

	yh := newHeader(yb.At(0),  yb.At(1), yb.At(2), yb.At(3), yb.At(4), yb.At(5))
	if uint32(yb.Count()) < yh.datalen + HeaderLen {
		return false
	}

	return true
}

// firstBlockSize
func (y *Blockbuf)firstBlockSize() int {
	return int(y.At(2) << 24 | y.At(3) << 16 | y.At(4) << 8 | y.At(5) )
}
// firstBlockData
// FIXME: use "block copy" instead of "byte copy"
func (y *Blockbuf)firstBlockData() []byte {
	l := y.firstBlockSize()
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = y.At(i + HeaderLen)
	}
	y.DropN(l + HeaderLen)
    return b
}
