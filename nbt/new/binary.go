package new

import (
	"io"
	"unsafe"
)

type reader struct {
	buf []byte
	at  int
}

func (r *reader) readByte() (byte, error) {
	if r.at+1 > len(r.buf) {
		return 0, nil
	}

	v := r.buf[r.at]
	r.at++
	return v, nil
}

func (r *reader) readShort() (int16, error) {
	if r.at+2 > len(r.buf) {
		return 0, nil
	}

	v := int16(r.buf[r.at])<<8 | int16(r.buf[r.at+1])
	r.at += 2
	return v, nil
}

func (r *reader) readInt() (int32, error) {
	if r.at+4 > len(r.buf) {
		return 0, nil
	}

	v := int32(r.buf[r.at])<<24 | int32(r.buf[r.at+1])<<16 | int32(r.buf[r.at+2])<<8 | int32(r.buf[r.at+3])
	r.at += 4
	return v, nil
}

func (r *reader) readString() (string, error) {
	if r.at+2 > len(r.buf) {
		return "", io.ErrUnexpectedEOF
	}

	v := int(r.buf[r.at])<<8 | int(r.buf[r.at+1])
	r.at += 2

	if r.at+v > len(r.buf) {
		return "", io.ErrUnexpectedEOF
	}

	str := r.buf[r.at : r.at+v]
	r.at += v
	return *(*string)(unsafe.Pointer(&str)), nil
}

func (r *reader) seek(n int) error {
	if r.at+n > len(r.buf) {
		return io.ErrUnexpectedEOF
	}

	r.at += n
	return nil
}
