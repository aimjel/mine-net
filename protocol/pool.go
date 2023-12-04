package protocol

import (
	"bytes"
	"compress/zlib"
	"sync"
)

type zlibHeader struct{}

func (z zlibHeader) Read(p []byte) (int, error) {
	return copy(p, []byte{0x78, 0x9c}), nil
}

var zlibReaders = sync.Pool{
	New: func() any {
		z, _ := zlib.NewReader(zlibHeader{})
		return z
	},
}

const s = 1 << 14

var (
	sizes = [...]int{
		1 << 10, //1  KB
		1 << 14, //16 KB
		1 << 20, //1  MB
	}

	buffers = [...]sync.Pool{
		{New: func() any { return bytes.NewBuffer(make([]byte, 0, 1<<10)) }},

		{New: func() any { return bytes.NewBuffer(make([]byte, 0, 1<<14)) }},

		{New: func() any { return bytes.NewBuffer(make([]byte, 0, 1<<20)) }},
	}
)

func GetBuffer(size int) *bytes.Buffer {
	for i, v := range sizes {
		if size < v {
			return buffers[i].Get().(*bytes.Buffer)
		}
	}

	return nil
}

func PutBuffer(b *bytes.Buffer) {
	b.Reset()

	for i, v := range sizes {
		if b.Cap() < v {
			buffers[i].Put(b)
		}
	}
}
