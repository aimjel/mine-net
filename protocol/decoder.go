package protocol

import (
	"crypto/cipher"
	"fmt"
	"github.com/aimjel/minecraft/protocol/crypto"
	"io"
)

const maxPacketLength = (1 << 21) - 1 //2mb

type Decoder struct {
	rd *Reader

	//todo zlib decompressor

	compressionThreshold int
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{rd: NewReader(r, 0)}
}

func (dec *Decoder) EnableDecryption(b cipher.Block, sharedSecret []byte) {
	dec.rd.dec = crypto.NewCFB8(b, sharedSecret, true)
}

//todo work on new decode packet system

func (dec *Decoder) DecodePacket() ([]byte, error) {
	length, err := dec.rd.readVarInt()
	if err != nil {
		return nil, err
	}

	if length > maxPacketLength || length == 0 {
		return nil, fmt.Errorf("invalid packet length")
	}

	//make sure there's space for the packet
	dec.rd.grow(length)

	for length > dec.rd.len() {
		if err = dec.rd.fill(); err != nil {
			return nil, err
		}
	}

	//todo decompression

	return dec.rd.next(length), nil
}
