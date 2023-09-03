package protocol

import (
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"fmt"
	"io"
)

const MaxPacket = 1 << 21

type Decoder struct {
	r *Reader

	decompressor io.ReadCloser

	threshold int

	putBack *bytes.Buffer
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: NewReader(r), threshold: -1}
}

func (dec *Decoder) EnableDecryption(block cipher.Block, iv []byte) {
	dec.r.EnableDecryption(block, iv)
}

// EnableDecompression Enables zlib decompression
func (dec *Decoder) EnableDecompression() {
	dec.decompressor, _ = zlib.NewReader(bytes.NewBuffer([]byte{0x78, 0x9c}))
	dec.threshold = 0 //doesn't actually do anything just allows the check to go through
}

// DecodePacket reads from the underlying reader and returns a packet's payload.
// The slice returned is valid until the next call.
func (dec *Decoder) DecodePacket() ([]byte, error) {
	pkLen, err := dec.r.ReadVarInt()
	if err != nil {
		return nil, fmt.Errorf("%v reading packet length", err)
	}

	if pkLen > MaxPacket || pkLen == 0 {
		return nil, fmt.Errorf("invalid packet length of %v", pkLen)
	}

	if dec.threshold != -1 {
		dataLength, err := dec.r.ReadVarInt()
		if err != nil {
			return nil, fmt.Errorf("%v reading data length", err)
		}

		if dataLength != 0 {
			return dec.decompress(dataLength)
		}

		pkLen--
	}

	buf := buffers.Get(pkLen)
	err, putBack := dec.r.writeTo(buf, pkLen)
	if putBack {
		dec.putBack = buf
	}

	return buf.Bytes(), nil
}

func (dec *Decoder) decompress(len int) ([]byte, error) {
	if err := dec.decompressor.(zlib.Resetter).Reset(dec.r, nil); err != nil {
		return nil, fmt.Errorf("%v resetting decompresor", err)
	}

	buf := buffers.Get(len)
	defer buffers.Put(buf)

	n, err := dec.decompressor.Read(buf.Bytes()[:len])
	if err != nil && n != len {
		return nil, fmt.Errorf("%v decompressing payload", err)
	}

	if n != len {
		return nil, fmt.Errorf("decompressed an incorrect amount of data")
	}

	return buf.Bytes()[:n], nil
}

func (dec *Decoder) PutBufferBack() {
	if dec.putBack != nil {
		buffers.Put(dec.putBack)
		dec.putBack = nil
	}
}
