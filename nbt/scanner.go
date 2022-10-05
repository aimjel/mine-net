package nbt

import (
	"fmt"
	"io"
)

func checkValid(data []byte) error {
	if len(data) < 4 || data[0] != tagCompound {
		return fmt.Errorf("data must always start with a compound tag")
	}

	return (&scanner{buf: data}).scan(data[0])
}

//scanner verifies nbt input syntax
type scanner struct {
	buf []byte
	at  int
}

func (s *scanner) scan(id byte) error {
	switch id {

	default:
		return fmt.Errorf("nbt: unknown tag id at %v", s.at)

	case tagByte:
		if s.at+1 > len(s.buf) {
			return io.ErrUnexpectedEOF
		}
		s.at++

	case tagShort:
		if s.at+2 > len(s.buf) {
			return io.ErrUnexpectedEOF
		}
		s.at += 2

	case tagInt, tagFloat:
		if s.at+4 > len(s.buf) {
			return io.ErrUnexpectedEOF
		}
		s.at += 4

	case tagLong, tagDouble:
		if s.at+8 > len(s.buf) {
			return io.ErrUnexpectedEOF
		}
		s.at += 8

	case tagByteArray:
		length, err := s.read32()
		if err != nil {
			return err
		}

		if s.at+length > len(s.buf) {
			return io.ErrUnexpectedEOF
		}
		s.at += length

	case tagString:
		if s.at+2 > len(s.buf) {
			return io.ErrUnexpectedEOF
		}

		length := int(uint16(s.buf[s.at])<<8 | uint16(s.buf[s.at+1]))
		s.at += 2

		if s.at+length > len(s.buf) {
			return io.ErrUnexpectedEOF
		}

		s.at += length

	case tagList:
		return s.list()

	case tagCompound:
		return s.compound()

	case tagIntArray:
		length, err := s.read32()
		if err != nil {
			return err
		}

		length *= 4

		if s.at+length > len(s.buf) {
			return io.ErrUnexpectedEOF
		}

		s.at += length

	case tagLongArray:
		length, err := s.read32()
		if err != nil {
			return err
		}

		length *= 8

		if s.at+length > len(s.buf) {
			return io.ErrUnexpectedEOF
		}

		s.at += length
	}

	return nil
}

func (s *scanner) list() error {
	if s.at+1 > len(s.buf) {
		return io.ErrUnexpectedEOF
	}

	id := s.buf[s.at]
	s.at++

	length, err := s.read32()
	if err != nil {
		return err
	}

	if id == tagEnd || length == 0 {
		return nil
	}

	for i := 0; i < length; i++ {
		if err = s.scan(id); err != nil {
			return err
		}
	}

	return nil
}

func (s *scanner) compound() error {
	for {
		if s.at+1 > len(s.buf) {
			if s.buf[s.at-1] == tagEnd {
				return nil
			}

			return io.ErrUnexpectedEOF
		}

		id := s.buf[s.at]
		s.at++

		if id == tagEnd {
			return nil
		}

		if s.at+2 > len(s.buf) { // check if we have enough data to read a short
			return io.ErrUnexpectedEOF
		}

		length := int(uint16(s.buf[s.at])<<8 | uint16(s.buf[s.at+1]))
		s.at += 2

		if s.at+length > len(s.buf) { // checks the length of the string is valid
			return io.ErrUnexpectedEOF
		}

		s.at += length

		if err := s.scan(id); err != nil {
			return err
		}
	}
}

func (s *scanner) read32() (int, error) {
	if s.at+4 > len(s.buf) {
		return 0, io.ErrUnexpectedEOF
	}

	data := s.buf[s.at : s.at+4]
	s.at += 4
	return int(uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])), nil
}

const (
	tagEnd uint8 = iota
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
