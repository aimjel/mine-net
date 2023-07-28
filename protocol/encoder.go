package protocol

import (
	"compress/zlib"
	"crypto/cipher"
	"encoding/binary"
	"github.com/aimjel/minecraft/protocol/crypto"
	"io"
)

type Encoder struct {
	wr io.Writer

	enc *crypto.CFB8

	zlib *zlib.Writer

	buf []byte

	start int // start index of the current packet
	at    int // current write index
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		wr:  w,
		buf: make([]byte, 4096),
	}
}

func (e *Encoder) EnableEncryption(b cipher.Block, sharedSecret []byte) {
	e.enc = crypto.NewCFB8(b, sharedSecret, false)
}

// Encode the unread bytes into the packet format
func (e *Encoder) Encode() {
	var l [3]byte

	n := binary.PutUvarint(l[:], uint64(e.at-e.start))

	//moves the written data n bytes to make space for the packet length
	copy(e.buf[e.start+n:e.at+n], e.buf[e.start:e.at])

	//copies the length into the space we made
	copy(e.buf[e.start:], l[:n])

	//add n bytes so we don`t overwrite data
	e.at += n
	e.start = e.at // move the start pointer so, we don`t re-encode the packet
	return
}

// Flush writes all the formatted packet data out to the underlying writer
func (e *Encoder) Flush() error {
	_, err := e.wr.Write(e.buf[:e.at])
	e.at, e.start = 0, 0
	return err
}

func (e *Encoder) Write(p []byte) (int, error) {
	if e.at+len(p) > len(e.buf) {
		e.grow(len(p))
	}

	n := copy(e.buf[e.at:], p)
	e.at += n
	return n, nil
}

func (e *Encoder) WriteByte(p byte) error {
	if e.at+1 > len(e.buf) {
		e.grow(1)
	}

	e.buf[e.at] = p
	e.at++
	return nil
}

func (e *Encoder) grow(n int) {
	c := make([]byte, len(e.buf)+(len(e.buf)/2)+n) // Incremental growth.
	e.at = copy(c, e.buf[:e.at])
	e.buf = c
}
