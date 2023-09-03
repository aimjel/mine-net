package protocol

import (
	"bytes"
	"crypto/cipher"
	"fmt"
	"github.com/aimjel/minecraft/protocol/crypto"
	"io"
)

type Reader struct {
	rd io.Reader

	buf []byte

	r, w int

	cipher *crypto.CFB8
}

func NewReader(r io.Reader) *Reader {
	return &Reader{rd: r, buf: make([]byte, 4096)}
}

func (r *Reader) EnableDecryption(block cipher.Block, iv []byte) {
	r.cipher = crypto.NewCFB8(block, iv, true)
}

// reset the reader and writer index to 0 so there's more space to write into
func (r *Reader) reset() {
	//should only be called if the readers index is the same as the writers.
	//if there's nothing to read then we can reset the reader and writer pointers
	//to the beginning of the slice so there's more space to write into.
	r.r, r.w = 0, 0
}

// len returns the number of unread bytes in the buffer
func (r *Reader) len() int {
	return r.w - r.r
}

// fill the buffer up
func (r *Reader) fill() error {
	//shifts the unread bytes to the front
	if r.r != 0 && r.r < r.w {
		r.w = copy(r.buf, r.buf[r.r:r.w])
		r.r = 0
	}

	x := 1
	for x > 0 { //Read till at least one byte has been Read
		n, err := r.rd.Read(r.buf[r.w:])
		if err != nil {
			return err
		}

		//decrypts the data
		if r.cipher != nil {
			r.cipher.XORKeyStream(r.buf[r.w:r.w+n], r.buf[r.w:r.w+n])
		}

		r.w += n
		x -= n
	}

	return nil
}

func (r *Reader) writeTo(b *bytes.Buffer, n int) (error, bool) {
	if n > len(r.buf) {
		var written int
		for written < n {
			nn, err := r.Read(b.Bytes()[written:n])
			if err != nil {
				return err, true
			}

			written += nn
		}

		b.Write(b.Bytes()[0:n])
		return nil, true
	}

	for r.len() < n {
		if err := r.fill(); err != nil {
			return err, false
		}
	}

	buf := r.buf[r.r : r.r+n]
	r.r += n
	*b = *bytes.NewBuffer(buf)
	return nil, false
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.r == r.w {
		r.reset()
		if err := r.fill(); err != nil {
			return 0, err
		}
	}

	n := copy(p, r.buf[r.r:r.w])
	r.r += n

	return n, nil
}

func (r *Reader) ReadByte() (byte, error) {
	if r.r == r.w {
		r.reset()
		if err := r.fill(); err != nil {
			return 0, err
		}
	}

	b := r.buf[r.r]
	r.r++
	return b, nil
}

func (r *Reader) ReadVarInt() (int, error) {
	var ux uint32

	for i := 0; i < 35; i += 7 {
		b, err := r.ReadByte()
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
