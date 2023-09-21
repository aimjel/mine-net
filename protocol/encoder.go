package protocol

import (
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"github.com/aimjel/minecraft/packet"
	"github.com/aimjel/minecraft/protocol/crypto"
)

type Encoder struct {
	buf *bytes.Buffer

	cipher *crypto.CFB8

	compressor *zlib.Writer

	threshold int

	headerSize int
}

func NewEncoder() *Encoder {
	return &Encoder{
		buf:        bytes.NewBuffer(make([]byte, 0, 4096)),
		threshold:  -1,
		headerSize: 3, //max pk length in bytes
	}
}

func (enc *Encoder) EnableEncryption(block cipher.Block, iv []byte) {
	enc.cipher = crypto.NewCFB8(block, iv, false)
}

func (enc *Encoder) EnableCompression(threshold int) {
	enc.compressor = zlib.NewWriter(nil)
	enc.threshold = threshold
	enc.headerSize = 3 + 5 //max pk length and data length in bytes
}

func (enc *Encoder) EncodePacket(pk packet.Packet) error {
	start := enc.buf.Len() //records the start of the packet data

	pw := packet.NewWriter(enc.buf)

	if err := pw.VarInt(pk.ID()); err != nil {
		return err
	}

	if err := pk.Encode(pw); err != nil {
		return err
	}

	payload := enc.buf.Bytes()[start:enc.buf.Len()]
	pkLen := enc.buf.Len() - start
	dataLength := -1
	if enc.threshold >= 0 {
		if pkLen >= enc.threshold {
			enc.buf.Truncate(start)
			return enc.compress(payload)
		}

		dataLength = 0
		pkLen++
	}

	enc.buf.Grow(enc.headerSize) //ensures the max header can fit

	//makes room for the header
	copy(enc.buf.Bytes()[start+enc.headerSize:enc.buf.Cap()], payload)

	//changes the writers position to write at the space we created
	enc.buf.Truncate(start)

	enc.writeHeader(pkLen, dataLength)

	start += enc.headerSize //updates the position where the data starts

	enc.buf.Write(enc.buf.Bytes()[start : start+len(payload)])

	return nil
}

// compresses the bytes of the buffer object passed
func (enc *Encoder) compress(payload []byte) error {
	buf := buffers.Get(len(payload))
	defer buffers.Put(buf)

	enc.compressor.Reset(buf)

	if _, err := enc.compressor.Write(payload); err != nil {
		return err
	}

	if err := enc.compressor.Flush(); err != nil {
		return err
	}

	enc.writeHeader(buf.Len()+varIntSize(len(payload)), len(payload))

	enc.buf.Write(buf.Bytes())
	return nil
}

func (enc *Encoder) Flush() []byte {
	data := enc.buf.Bytes()
	enc.buf.Reset()
	if enc.cipher != nil {
		enc.cipher.XORKeyStream(data, data)
	}

	return data
}

func writeVarInt(b *bytes.Buffer, n int) {
	ux := uint32(n)

	for ux >= 0x80 {
		b.WriteByte(byte(ux&0x7F) | 0x80)
		ux >>= 7
	}

	b.WriteByte(byte(ux))
}

func (enc *Encoder) writeHeader(pkLen, dataLen int) {
	writeVarInt(enc.buf, pkLen)
	if dataLen != -1 {
		writeVarInt(enc.buf, dataLen)
	}
}

// varIntSize returns the number of bytes n takes up
func varIntSize(n int) (i int) {
	for n >= 0x80 {
		n >>= 7
		i++
	}
	i++
	return
}
