package minecraft

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"io"
	"minecraft/packet"
	"minecraft/player"
	"net"
	"sync"
)

type Conn struct {
	tcp *net.TCPConn

	dec *decoder

	enc *encoder

	pool map[int32]func() packet.Packet

	Info *player.Info
}

func Dial(address string) (*Conn, error) {
	cn, err := net.Dial("tcp4", address)
	if err != nil {
		return nil, err
	}

	return newConn(cn.(*net.TCPConn)), nil
}

func newConn(c *net.TCPConn) *Conn {
	return &Conn{
		tcp: c,
		dec: &decoder{
			rd:  c,
			buf: make([]byte, 1024),
		},
		enc: &encoder{
			wr:  c,
			buf: make([]byte, 1024*1024),
			w:   3, //ensures the packet length fits
			r:   3,
		},
	}
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.tcp.RemoteAddr()
}

func (c *Conn) Close() {
	if err := c.tcp.Close(); err != nil {
		fmt.Printf("%v: %v closing connection\n", c.RemoteAddr(), err)
	}
}

func (c *Conn) enableEncryption(b cipher.Block, iv []byte) {
	c.dec.cipher = newCFB(b, iv, true)
	c.enc.cipher = newCFB(b, iv, false)
}

func (c *Conn) WritePacket(pk packet.Packet) error {
	c.enc.mu.Lock()
	defer c.enc.mu.Unlock()
	w := packet.NewWriter(c.enc)

	if err := w.VarInt(pk.ID()); err != nil {
		return err
	}

	if err := pk.Encode(w); err != nil {
		return err
	}

    fmt.Printf("%v: packet %x sent\n", c.RemoteAddr(), pk.ID())
	return c.enc.flush()
}

func (c *Conn) ReadPacket() (packet.Packet, error) {
	p, err := c.dec.decodePacket()
	if err != nil {
		return nil, err
	}

	var id int32
	r := packet.NewReader(p)
	if err = r.VarInt(&id); err != nil {
		return nil, fmt.Errorf("%v reading packet id", err)
	}

	fn, ok := c.pool[id]
	if ok == false {
		//if id == 0 {
		//	fmt.Println(string(p))
		//}
		//fmt.Printf("%v: unknown packet with id %x\n", c.RemoteAddr(), id)
		return nil, nil
	}

	pk := fn()

	if err = pk.Decode(r); err != nil {
		return nil, fmt.Errorf("%v decoding packet", err)
	}

	return pk, nil
}

func (c *Conn) DecodePacket(pk packet.Packet) error {
	p, err := c.dec.decodePacket()
	if err != nil {
		return err
	}

	var id int32
	r := packet.NewReader(p)
	if err = r.VarInt(&id); err != nil {
		return fmt.Errorf("%v reading packet id", err)
	}

	if id != pk.ID() {
		return fmt.Errorf("unexpected id")
	}

	return pk.Decode(r)
}

type encoder struct {
	wr io.Writer //connection

	cipher cipher.Stream

	mu sync.Mutex //protects buffer

	buf []byte

	r, w int
}

func (enc *encoder) flush() error {
	n := binary.PutUvarint(enc.buf, uint64(enc.w-enc.r))

	//move the packet length value to the right index
	copy(enc.buf[enc.r-n:], enc.buf[:n])

	if enc.cipher != nil {
		enc.cipher.XORKeyStream(enc.buf[enc.r-n:enc.w], enc.buf[enc.r-n:enc.w])
	}

	_, err := enc.wr.Write(enc.buf[enc.r-n : enc.w])
	enc.r, enc.w = 3, 3

	return err
}

func (enc *encoder) Write(p []byte) (int, error) {
	if enc.w+len(p) > len(enc.buf) {
		enc.grow(len(p))
	}

	n := copy(enc.buf[enc.w:], p)
	enc.w += n
	return n, nil
}

func (enc *encoder) WriteByte(c byte) error {
	if enc.w+1 > len(enc.buf) {
		enc.grow(1)
	}

	enc.buf[enc.w] = c
	enc.w++
	return nil
}

func (enc *encoder) grow(n int) {
	c := make([]byte, (len(enc.buf)+n)*2)
	enc.w = copy(c, enc.buf[:enc.w])
	enc.buf = c
}
