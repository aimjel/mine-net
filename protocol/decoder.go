package protocol

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
)

const MaxPacket = 1 << 21

type Decoder struct {
	r *Reader

	decompression bool

	tmp      *bytes.Buffer
	tmpInUse bool
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: NewReader(r)}
}

func (dec *Decoder) EnableDecryption(block cipher.Block, iv []byte) {
	dec.r.EnableDecryption(block, iv)
}

// EnableDecompression Enables zlib decompression
func (dec *Decoder) EnableDecompression() {
	dec.decompression = true
}

// DecodePacket reads from the underlying reader and returns a packet's payload.
// The slice returned is valid until the next call.
func (dec *Decoder) DecodePacket() ([]byte, error) {
	dec.tmpInUse = false

	pkLen, err := dec.r.ReadVarInt()
	if err != nil {
		return nil, fmt.Errorf("%w reading packet length", err)
	}

	if pkLen > MaxPacket || pkLen == 0 {
		return nil, fmt.Errorf("invalid packet length of %v", pkLen)
	}

	if dec.decompression {
		dataLength, err := dec.r.ReadVarInt()
		if err != nil {
			return nil, fmt.Errorf("%w reading data length", err)
		}

		if dataLength != 0 {
			return dec.decompress(dataLength)
		}

		pkLen--
	}

	payload, err := dec.r.Next(pkLen)
	if err != nil {
		if errors.Is(err, bufio.ErrBufferFull) {
			buf := dec.buffer(pkLen)
			return buf.Bytes()[:pkLen], dec.r.readFull(buf, pkLen)
		}
	}

	if dec.tmpInUse == false && dec.tmp != nil {
		buffers.Put(dec.tmp)
		dec.tmp = nil
	}

	return payload, err
}

func (dec *Decoder) decompress(len int) ([]byte, error) {
	zr := zlibReaders.Get().(io.ReadCloser)
	if err := zr.(zlib.Resetter).Reset(dec.r, nil); err != nil {
		return nil, fmt.Errorf("%w resetting decompresor", err)
	}

	buf := dec.buffer(len)
	buf.Grow(len)

	var written int
	for written != len {
		n, err := zr.Read(buf.Bytes()[written:len])
		written += n
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("%w decompressing payload", err)
		}
	}

	if written != len {
		return nil, fmt.Errorf("decompressed %v bytes but expected %v bytes", written, len)
	}

	return buf.Bytes()[:written], nil
}

func (dec *Decoder) buffer(size int) *bytes.Buffer {
	dec.tmpInUse = true
	if dec.tmp != nil {
		dec.tmp.Grow(size)
		dec.tmp.Reset()
		return dec.tmp
	}

	dec.tmp = buffers.Get(size)
	return dec.tmp
}
