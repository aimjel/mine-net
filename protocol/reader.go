package protocol

import (
	"fmt"
	"github.com/aimjel/minecraft/protocol/crypto"
	"io"
)

// Reader is a buffer reader with decryption built-in
type Reader struct {
	rd io.Reader

	dec *crypto.CFB8

	buf []byte

	r, w int
}

func NewReader(r io.Reader, size int) *Reader {
	if size == 0 {
		size = 4096
	}

	return &Reader{
		rd:  r,
		buf: make([]byte, size),
	}
}

// len returns the number of unread bytes
func (rd *Reader) len() int {
	return rd.w - rd.r
}

func (rd *Reader) available() int {
	return len(rd.buf) - rd.w
}

// grow makes sure there's enough space for n bytes
func (rd *Reader) grow(n int) {
	if n < len(rd.buf)-rd.len() {
		rd.w, rd.r = copy(rd.buf, rd.buf[rd.r:rd.w]), 0
	} else {
		b := make([]byte, rd.len()+n)
		rd.w, rd.r = copy(b, rd.buf[rd.r:rd.w]), 0
		rd.buf = b
	}
}

// read data into b and performs decryption if needed
func (rd *Reader) fill() error {

	if rd.r > 0 {
		//slides the unread bytes to the front
		rd.w = copy(rd.buf, rd.buf[rd.r:rd.w])
		rd.r = 0
	}

	n, err := rd.rd.Read(rd.buf[rd.w:])
	if err != nil {
		return err
	}

	rd.w += n

	if rd.dec != nil {
		rd.dec.XORKeyStream(rd.buf[rd.w-n:rd.w], rd.buf[rd.w-n:rd.w])
	}

	return nil
}

func (rd *Reader) readByte() (byte, error) {
	if rd.r == rd.w {
		if err := rd.fill(); err != nil {
			return 0, err
		}
	}

	b := rd.buf[rd.r]
	rd.r++

	return b, nil
}

func (rd *Reader) readVarInt() (int, error) {
	var ux uint32

	for i := 0; i < 35; i += 7 {
		b, err := rd.readByte()
		if err != nil {
			return 0, err
		}

		ux |= uint32(b&0x7f) << i

		if b&0x80 == 0 {
			return int(int32(ux)), nil
		}
	}

	return 0, fmt.Errorf("var-int overflows 32-bit integer")
}

func (rd *Reader) next(n int) []byte {
	p := rd.buf[rd.r : rd.r+n]
	rd.r += n
	return p
}
