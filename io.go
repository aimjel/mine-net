package minecraft

import (
	"crypto/cipher"
	"fmt"
	"io"
)

const maxPacketLength = (1 << 21) - 1 //2mb

type decoder struct {
	rd io.Reader

	cipher cipher.Stream

	buf []byte

	r, w int
}

func (dec *decoder) decodePacket() ([]byte, error) {
	length, err := dec.ReadVarInt()
	if err != nil {
		return nil, err
	}

	if length > maxPacketLength || length == 0 {
		return nil, fmt.Errorf("invalid packet length")
	}

	if length >= len(dec.buf) {
		b := make([]byte, (len(dec.buf)+length)*2)
		dec.w = copy(b, dec.buf[dec.r:dec.w])
		dec.r = 0
		dec.buf = b
	}

	for length > dec.w-dec.r {
		if err = dec.read(dec.buf[dec.w:]); err != nil {
			return nil, err
		}
	}
	p := dec.buf[dec.r : dec.r+length]
	dec.r += length

	if dec.r == dec.w {
		dec.r, dec.w = 0, 0
	}
	return p, nil
}

func (dec *decoder) ReadVarInt() (int, error) {
	var ux uint32

	for i := 0; i < 35; i += 7 {
		b, err := dec.ReadByte()
		if err != nil {
			return 0, err
		}

		ux |= uint32(b&0x7f) << i

		if b&0x80 == 0 {
			return int(ux), nil
		}
	}

	return 0, fmt.Errorf("var-int overflows 32-bit integer")
}

func (dec *decoder) ReadByte() (byte, error) {
	if dec.r == dec.w {
		if err := dec.read(dec.buf[dec.w:]); err != nil {
			return 0, err
		}
	}

	b := dec.buf[dec.r]
	dec.r++

	return b, nil
}

func (dec *decoder) read(p []byte) error {
	n, err := dec.rd.Read(p)
	dec.w += n

	if dec.cipher != nil {
		dec.cipher.XORKeyStream(p[:n], p[:n])
	}

	return err
}
