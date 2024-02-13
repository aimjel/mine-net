package nbt

import (
	"errors"
	"fmt"
)

type ErrUnknownTag struct {
	tag byte
	at  int
}

func (e ErrUnknownTag) Error() string {
	return fmt.Sprintf("unknown tag id %v at index %v", e.tag, e.at)
}

func checkValid(data []byte) (err error) {
	if len(data) < 4 || data[0] != tagCompound {
		return fmt.Errorf("data must always start with a compound tag")
	}

	defer func() {
		if x := recover(); x != nil {
			switch e := x.(type) {

			case error:
				err = e

			case string:
				err = errors.New(e)
			}
		}
	}()
	//set reader position to 1, since we supplied the first byte
	s := &scanner{buf: data, at: 1}
	s.scanString()
	s.scan(data[0])

	if s.at != len(data) {
		panic("scanner did not consume all data")
	}
	return err
}

// scanner verifies nbt input syntax
type scanner struct {
	buf []byte
	at  int
}

func (s *scanner) scan(id byte) {
	switch id {
	default:
		panic(ErrUnknownTag{tag: id, at: s.at})

	case tagByte, tagShort, tagInt, tagLong, tagFloat, tagDouble:
		n := tagPayload(id)
		s.at += n

	case tagString:
		s.scanString()

	case tagCompound:
		s.scanCompound()

	case tagList:
		s.scanList()

	case tagByteArray, tagIntArray, tagLongArray:
		ln := s.read32()
		s.at += tagPayload(id) * ln
	}
}

func (s *scanner) scanList() {
	id := s.buf[s.at]
	s.at++

	ln := s.read32()

	if id == 0 || ln == 0 {
		return
	}

	for i := 0; i < ln; i++ {
		s.scan(id)
	}
}

func (s *scanner) scanCompound() {
	for {
		id := s.buf[s.at]
		s.at++

		if id == tagEnd {
			return
		}

		s.scanString()
		s.scan(id)
	}
}

func (s *scanner) scanString() {
	length := int(uint16(s.buf[s.at])<<8 | uint16(s.buf[s.at+1]))
	s.at += 2

	s.at += length
}

func (s *scanner) read32() int {
	data := s.buf[s.at : s.at+4]
	s.at += 4
	return int(uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3]))
}

func tagPayload(id byte) int {
	switch id {
	case tagByte, tagByteArray:
		return 1
	case tagShort:
		return 2
	case tagInt, tagFloat, tagIntArray:
		return 4
	case tagLong, tagDouble, tagLongArray:
		return 8
	default:
		return 0
	}
}
func tagName(id byte) string {
	switch id {
	case tagEnd:
		return "tag_end"
	case tagByte:
		return "tag_byte"
	case tagShort:
		return "tag_short"
	case tagInt:
		return "tag_int"
	case tagLong:
		return "tag_long"
	case tagFloat:
		return "tag_float"
	case tagDouble:
		return "tag_double"
	case tagByteArray:
		return "tag_byte_array"
	case tagString:
		return "tag_string"
	case tagList:
		return "tag_list"
	case tagCompound:
		return "tag_compound"
	case tagIntArray:
		return "tag_int_array"
	case tagLongArray:
		return "tag_long_array"
	default:
		return "unknown"
	}
}
