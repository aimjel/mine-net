package new

import (
	"unsafe"
)

type parser struct {
	buf []byte
	at  int

	compoundsOpen int

	tags []tag
}

func (p *parser) getTag(id byte) tag {
	for _, t := range p.tags {
		if t.typ == id {
			return t
		}
	}

	return tag{}
}

func (p *parser) parse(t tag) {
	p.at = t.start
	p.tags = nil

	for {
		id := p.readByte()

		if id == tagEnd {
			p.compoundsOpen--

			if p.compoundsOpen > 0 {
				continue
			}

			return
		}

		name := p.readString()

		switch id {
		case tagByte, tagShort, tagInt, tagLong, tagFloat, tagDouble:
			p.at += tagPayload(id)
		case tagString:
			_ = p.readString()
		case tagList:
			p.parseList()

		case tagCompound:
			p.compoundsOpen++
		}

		if p.compoundsOpen == 1 {
			p.addTag(id, p.at, name)
		}
	}
}

func (p *parser) parseList() {
	id := p.readByte()
	length := p.readInt32()

	if id == 0 || length == 0 {
		return
	}

	switch id {
	case tagByte, tagShort, tagInt, tagLong, tagFloat, tagDouble:
		p.at += tagPayload(id) * length
	case tagCompound:
		p.compoundsOpen += length
	}
}

func (p *parser) addTag(typ byte, start int, name string) {
	if cap(p.tags) > len(p.tags) {
		p.tags = p.tags[:len(p.tags)+1]
	} else {
		p.tags = append(p.tags, tag{})
	}

	t := &p.tags[len(p.tags)-1]
	t.typ, t.start, t.name = typ, start, name
}

func tagPayload(id byte) int {
	switch id {

	case tagByte:
		return 1

	case tagShort:
		return 2

	case tagInt, tagFloat:
		return 4

	case tagLong, tagDouble:
		return 8
	}

	return 0
}

func (p *parser) readByte() byte {
	v := p.buf[p.at]
	p.at++
	return v
}

func (p *parser) readInt32() int {
	v := int(p.buf[p.at])<<24 | int(p.buf[p.at+1])<<16 | int(p.buf[p.at+2])<<8 | int(p.buf[p.at+3])
	p.at += 4
	return v
}

func (p *parser) readString() string {
	v := int(p.buf[p.at])<<8 | int(p.buf[p.at+1])
	p.at += 2

	str := p.buf[p.at : p.at+v]
	p.at += v
	return *(*string)(unsafe.Pointer(&str))
}

type tag struct {
	typ   uint8
	start int
	name  string //unsafe
}

func (t tag) isZero() bool {
	return t.typ == 0
}
