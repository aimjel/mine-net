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

	written int
}

func NewEncoder() *Encoder {
	return &Encoder{
		buf:       bytes.NewBuffer(make([]byte, 0, 4096)),
		threshold: -1,
	}
}

func (enc *Encoder) EnableEncryption(block cipher.Block, iv []byte) {
	enc.cipher = crypto.NewCFB8(block, iv, false)
}

func (enc *Encoder) EnableCompression(threshold int) {
	enc.compressor = zlib.NewWriter(nil)
	enc.threshold = threshold
}

func (enc *Encoder) Encode(pk packet.Packet) error {
	tmp := bytes.NewBuffer(enc.buf.Bytes()[enc.written:enc.buf.Cap()][:0])

	pw := packet.NewWriter(tmp)

	if err := pw.VarInt(pk.ID()); err != nil {
		return err
	}

	if err := pk.Encode(pw); err != nil {
		return err
	}

	dataLength := -1
	if enc.threshold != -1 {
		dataLength = 0

		if tmp.Len() >= enc.threshold {
			return enc.compress(tmp)
		}
	}

	enc.writeHeader(tmp, tmp.Len(), dataLength)
	return nil
}

// compresses the bytes of the buffer object passed
func (enc *Encoder) compress(payload *bytes.Buffer) error {
	//default objects for compressing
	uncompressedLength := payload.Len()
	localBuf := bytes.NewBuffer(payload.Bytes()[:0])

	//go uses a default window size of 32768 for zlib
	//if the data length is over the window size we should
	//pull a new buffer so data isn't overwritten
	if payload.Len() >= 32768 {
		// Formula: Compressed Size = Raw Data Size + Raw Data Size / 1000 + 12 bytes
		estimatedSize := payload.Len() + payload.Len()/1000 + 12
		localBuf = bytes.NewBuffer(buffers.Get(estimatedSize + enc.buf.Len()).Bytes()[:0])
		defer buffers.Put(localBuf)

		//the payload is going be written to a new buffer not equal to the one in the encoder struct
		//we need to copy the enc.buf.Bytes() into the new buffer and free the enc.buffer
		localBuf.Write(enc.buf.Bytes()[:enc.written])
		enc.buf = localBuf
		buffers.Put(enc.buf)
	}
	enc.compressor.Reset(localBuf)

	localBuf.Grow(2) //guarantee space for the zlib headers
	copy(localBuf.Bytes()[2:localBuf.Cap()], localBuf.Bytes()[:localBuf.Cap()])

	_, err := enc.compressor.Write(payload.Bytes()[2 : payload.Len()+2])
	if err != nil {
		return err
	}

	if err = enc.compressor.Flush(); err != nil {
		return err
	}

	pkLen := localBuf.Len() + varIntSize(uncompressedLength)
	dataLen := uncompressedLength
	enc.writeHeader(localBuf, pkLen, dataLen)
	return nil
}

func (enc *Encoder) Flush() []byte {
	n := enc.written
	enc.written = 0
	enc.buf.Reset()
	data := enc.buf.Bytes()[:n]
	if enc.cipher != nil {
		enc.cipher.XORKeyStream(data, data)
	}

	return data
}

func (enc *Encoder) writeHeader(b *bytes.Buffer, pkLen, dataLen int) {
	headerLength := varIntSize(pkLen)

	if dataLen != -1 {
		headerLength += varIntSize(dataLen)
	}

	b.Grow(headerLength)
	sl := b.Bytes()[:b.Len()+headerLength]
	copy(sl[headerLength:], sl[:b.Len()])

	b.Reset()

	writeVarInt(b, pkLen)
	if dataLen != -1 {
		writeVarInt(b, dataLen)
	}

	enc.written += len(sl)

	if b.Cap() >= enc.buf.Cap() {
		buffers.Put(enc.buf)
		enc.buf = b
	}
}

func writeVarInt(b *bytes.Buffer, n int) {
	ux := uint32(n)

	for ux >= 0x80 {
		b.WriteByte(byte(ux&0x7F) | 0x80)
		ux >>= 7
	}

	b.WriteByte(byte(ux))
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
