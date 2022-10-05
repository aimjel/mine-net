package nbt

import (
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

type parser struct {
	buf []byte
	at  int
}

func (p *parser) parseTag(id byte) tag {
	switch id {

	case tagByte:
		data := p.buf[p.at]
		p.at++

		return tag{typ: tagByte, rv: reflect.ValueOf(int8(data))}

	case tagShort:
		data := p.buf[p.at : p.at+2]
		p.at += 2
		return tag{typ: tagShort, rv: reflect.ValueOf(int16(uint16(data[0])<<8 | uint16(data[1])))}

	case tagInt:
		return tag{typ: tagInt, rv: reflect.ValueOf(p.read32())}

	case tagLong:
		return tag{typ: tagLong, rv: reflect.ValueOf(p.read64())}

	case tagFloat:
		return tag{typ: tagFloat, rv: reflect.ValueOf(math.Float32frombits(uint32(p.read32())))}

	case tagDouble:
		return tag{typ: tagDouble, rv: reflect.ValueOf(math.Float64frombits(uint64(p.read64())))}

	case tagByteArray:
		length := p.read32()

		data := p.buf[p.at : p.at+int(length)]
		p.at += int(length)
		v := make([]byte, 0, length)
		copy(v, data)

		return tag{typ: tagByteArray, rv: reflect.ValueOf(v)}

	case tagString:
		l := p.buf[p.at : p.at+2]
		p.at += 2

		length := int(uint16(l[0])<<8 | uint16(l[1]))

		v := p.buf[p.at : p.at+length]
		p.at += length

		return tag{typ: tagString, rv: reflect.ValueOf(string(v))}

	case tagList:
		return p.parseList()

	case tagCompound:
		return p.parseCompound()

	case tagIntArray:
		return p.parseIntArray()

	case tagLongArray:
		return p.parseLongArray()
	}

	// should never happen unless the byte slice is changed unsafely
	panic("unknown tag " + strconv.FormatInt(int64(id), 10))
}

func (p *parser) parseList() tag {
	t := tag{typ: tagList}

	id := p.buf[p.at]
	p.at++

	length := p.read32()

	if id == 0 || length == 0 {
		return t
	}

	tags := make([]tag, length, length)
	for i := int32(0); i < length; i++ {
		tags[i] = p.parseTag(id)
	}

	t.l = tags

	return t
}

func (p *parser) parseCompound() tag {
	t := tag{typ: tagCompound, c: compoundTag{}}
	for {
		id := p.buf[p.at]
		p.at++

		if id == tagEnd {
			return t
		}

		l := p.buf[p.at : p.at+2]
		p.at += 2

		length := int(uint16(l[0])<<8 | uint16(l[1]))

		tagName := p.buf[p.at : p.at+length]
		p.at += length

		pTag := p.parseTag(id)

		if cap(t.c.namedTags) > len(t.c.namedTags) {
			t.c.namedTags = t.c.namedTags[:len(t.c.namedTags)+1]
		} else {
			t.c.namedTags = append(t.c.namedTags, namedTag{})
		}

		nt := &t.c.namedTags[len(t.c.namedTags)-1]
		nt.k, nt.v = *(*string)(unsafe.Pointer(&tagName)), pTag
	}
}

func (p *parser) parseIntArray() tag {
	t := tag{typ: tagIntArray}
	length := p.read32()

	v := make([]int32, length)
	for i := int32(0); i < length; i++ {
		v[i] = p.read32()
	}

	t.rv = reflect.ValueOf(v)
	return t
}

func (p *parser) parseLongArray() tag {
	t := tag{typ: tagLongArray}
	length := int(p.read32())

	v := make([]int64, length)

	for i := 0; i < length; i++ {
		v[i] = p.read64()
	}

	t.rv = reflect.ValueOf(v)
	return t
}

func (p *parser) read32() int32 {
	data := p.buf[p.at : p.at+4]
	p.at += 4
	return int32(uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3]))
}

func (p *parser) read64() int64 {
	data := p.buf[p.at : p.at+8]
	p.at += 8
	return int64(uint64(data[0])<<56 | uint64(data[1])<<48 | uint64(data[2])<<40 | uint64(data[3])<<32 |
		uint64(data[4])<<24 | uint64(data[5])<<16 | uint64(data[6])<<8 | uint64(data[7]))
}

type compoundTag struct {
	namedTags []namedTag
}

func (c *compoundTag) getTag(id byte, key string) (tag, bool) {
	for _, t := range c.namedTags {
		if t.k == key && t.v.typ == id {
			return t.v, true
		}
	}

	return tag{}, false
}

func (c *compoundTag) Map() map[string]interface{} {
	m := make(map[string]interface{})

	for _, t := range c.namedTags {
		if t.v.typ == tagCompound {
			m[clone(t.k)] = t.v.c.Map()
			continue
		}

		m[clone(t.k)] = t.v.rv
	}

	return m
}

func clone(s string) string {
	b := make([]byte, len(s))
	copy(b, s)
	return *(*string)(unsafe.Pointer(&b))
}

type namedTag struct {
	k string
	v tag
}

type tag struct {
	typ uint8

	rv reflect.Value

	c compoundTag

	l []tag
}
