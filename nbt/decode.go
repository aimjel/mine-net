package nbt

import (
	"io"
	"unsafe"
)

// decoder is a buffered decoder
type decoder struct {
	reader io.Reader

	buf []byte

	r, w int
}

func newDecoder(rd io.Reader) *decoder {
	return &decoder{
		reader: rd,
		buf:    make([]byte, 1024),
	}
}

func newDecoderWithBytes(b []byte) *decoder {
	return &decoder{
		buf: b,
		w:   len(b),
	}
}

func (rd *decoder) available() int {
	return rd.w - rd.r
}

func (rd *decoder) fill(min int) error {
	if rd.available() < min {
		//moves the unread portion to the front of the buffer
		rd.r, rd.w = 0, copy(rd.buf, rd.buf[rd.r:rd.w])

		for i := 0; i < 5; i++ {
			n, err := rd.reader.Read(rd.buf[rd.w:])
			rd.w += n
			if err != nil {
				return err
			}

			min -= n
			if min <= 0 || rd.w == len(rd.buf) {
				return nil
			}
		}
		return io.ErrNoProgress
	}

	return nil
}

func (rd *decoder) readByte() (byte, error) {
	if err := rd.fill(1); err != nil {
		return 0, err
	}

	b := rd.buf[rd.r]
	rd.r++
	return b, nil
}

func (rd *decoder) readInt16() (int, error) {
	if err := rd.fill(2); err != nil {
		return 0, err
	}

	b := rd.buf[rd.r : rd.r+2]
	rd.r += 2
	return int(b[0])<<8 | int(b[1]), nil
}

func (rd *decoder) readInt32() (int, error) {
	if err := rd.fill(4); err != nil {
		return 0, err
	}

	b := rd.buf[rd.r : rd.r+4]
	rd.r += 4
	return int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3]), nil
}

func (rd *decoder) readInt64() (int, error) {
	if err := rd.fill(8); err != nil {
		return 0, err
	}

	b := rd.buf[rd.r : rd.r+8]
	rd.r += 8
	return int(b[0])<<56 | int(b[1])<<48 | int(b[2])<<40 | int(b[3])<<32 | int(b[4])<<24 | int(b[5])<<16 | int(b[6])<<8 | int(b[7]), nil
}

func (rd *decoder) readTag() (byte, error) {
	id, err := rd.readByte()
	if err != nil {
		return 0, err
	}

	if id == tagCompound {
		if _, err = rd.readUnsafeString(); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (rd *decoder) readUnsafeString() (string, error) {
	v, err := rd.readInt16()
	if err != nil {
		return "", err
	}

	if err = rd.fill(v); err != nil {
		return "", err
	}

	if v > len(rd.buf) {
		panic("string length bigger than buffer")
	}

	buf := rd.buf[rd.r : rd.r+v]
	rd.r += v

	return *(*string)(unsafe.Pointer(&buf)), nil
}

func (rd *decoder) skip(n int) error {
	avail := rd.w - rd.r
	if avail > n {
		rd.r += n
		return nil
	}

	n -= avail
	rd.r += avail

	rd.r, rd.w = 0, copy(rd.buf, rd.buf[rd.r:rd.w])

	for n != 0 {
		mx := n
		if mx > len(rd.buf) {
			mx = len(rd.buf)
		}
		x, err := rd.reader.Read(rd.buf[rd.w : rd.w+mx])
		rd.w += x
		rd.r += x
		if err != nil {
			return err
		}

		if rd.w == len(rd.buf) {
			rd.r, rd.w = 0, 0
		}

		n -= x
	}

	return nil
}

const (
	tagEnd = iota
	tagByte
	tagShort
	tagInt
	tagLong
	tagFloat
	tagDouble
	tagByteArray
	tagString
	tagList
	tagCompound
	tagIntArray
	tagLongArray
)
