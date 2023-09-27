package protocol

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"fmt"
	"github.com/aimjel/minecraft/protocol/crypto"
	"io"
)

// Reader wraps bufio.Reader and adds decryption functionality.
type Reader struct {
	buf *bufio.Reader

	cipher *crypto.CFB8
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		buf: bufio.NewReaderSize(r, 4096),
	}
}

func (r *Reader) EnableDecryption(block cipher.Block, iv []byte) {
	r.cipher = crypto.NewCFB8(block, iv, true)
}

// Next returns the next n bytes in the buffer and advances the reader.
// If n is bigger than the buffer then it returns bufio.ErrBufferFull
// without data and won't advance the reader.
func (r *Reader) Next(n int) ([]byte, error) {
	b, err := r.buf.Peek(n)
	if err != nil {
		return nil, err
	}

	if r.cipher != nil && err == nil {
		r.cipher.XORKeyStream(b, b)
	}

	_, err = r.buf.Discard(len(b))
	return b, err
}

func (r *Reader) Read(p []byte) (int, error) {
	n, err := r.buf.Read(p)

	if r.cipher != nil && n != 0 {
		r.cipher.XORKeyStream(p, p)
	}

	return n, err
}

func (r *Reader) ReadByte() (byte, error) {
	b, err := r.buf.ReadByte()

	x := []byte{b}

	if r.cipher != nil && err == nil {
		r.cipher.XORKeyStream(x, x)
	}

	return x[0], err
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

// readFull reads len bytes into b
func (r *Reader) readFull(b *bytes.Buffer, len int) error {
	var written int
	b.Grow(len) //makes sure n bytes can fit
	for written != len {
		x, err := r.Read(b.Bytes()[written:len])
		if err != nil {
			return err
		}

		written += x
	}

	return nil
}
