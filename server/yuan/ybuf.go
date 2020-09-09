// Yuan protocol formtat:
//	| Magic(2) | dataLen(4) | Data(variable) |

package yuan

import (
		"bytes"
		"errors"
	   )

const (
		Magic1		= 0xAA
		Magic2		= 0x55
		YHeaderLen	= 6
	  )

// YHeader
type YHeader struct {
	magic1	byte
	magic2	byte
	datalen	uint32
}
// YBlock
type YBlock struct {
	YHeader
	data	[]byte
}

// YBuf
type YBuf struct {
	CBC
}

// YHeader
func newYHeader(a,b,c,d,e,f byte) YHeader {
	var yh YHeader
	yh.magic1 = a
	yh.magic2 = b
	yh.datalen= uint32(c << 24) | uint32(d << 16) | uint32(e << 8) | uint32(f)
	return yh
}

func (y *YBuf)seekmagic() bool {
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
func (yb *YBuf)ScanBlock() bool {
	if ok := yb.seekmagic(); !ok {
		return false
	}

	if yb.Count() < YHeaderLen {
		return false
	}

	yh := newYHeader(yb.At(0),  yb.At(1), yb.At(2), yb.At(3), yb.At(4), yb.At(5))
	if uint32(yb.Count()) < yh.datalen + YHeaderLen {
		return false
	}

	return true
}

// firstBlockSize
func (y *YBuf)firstBlockSize() int {
	return int(y.At(2) << 24 | y.At(3) << 16 | y.At(4) << 8 | y.At(5) )
}
// firstBlockData
// FIXME: use "block copy" instead of "byte copy"
func (y *YBuf)firstBlockData() []byte {
	l := y.firstBlockSize()
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = y.At(i + YHeaderLen)
	}
	y.DropN(l + YHeaderLen)
    return b
}

// ReadBlock
func (y *YBuf)ReadBlock(b *bytes.Buffer) bool {
	if !y.ScanBlock() {
		return false
	}

	b.Write(y.firstBlockData())
	return true
}
// Write
func (y *YBuf)Write(p []byte) (n int, err error) {
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
