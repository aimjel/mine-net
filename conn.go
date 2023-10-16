package minecraft

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"github.com/aimjel/minecraft/packet"
	"github.com/aimjel/minecraft/protocol"
	"github.com/aimjel/minecraft/protocol/types"
	"net"
	"sync"
)

type Conn struct {
	tcpCn *net.TCPConn

	dec *protocol.Decoder

	enc *protocol.Encoder
	//buf which encoder writes to
	buf *bytes.Buffer

	pool Pool

	//encMu protects the Encoder from data races if two goroutines try to write a packet
	encMu sync.Mutex

	//name is the clients in-game name
	name string

	uuid [16]byte

	properties []types.Property
}

func (c *Conn) Name() string {
	return c.name
}

func (c *Conn) UUID() [16]byte {
	return c.uuid
}

func (c *Conn) Properties() []types.Property {
	return c.properties
}

func newConn(c *net.TCPConn) *Conn {
	b := bytes.NewBuffer(make([]byte, 0, 4096))
	return &Conn{
		tcpCn: c,

		dec: protocol.NewDecoder(c),

		enc: protocol.NewEncoder(b),
		buf: b,
	}
}

func (c *Conn) ReadPacket() (packet.Packet, error) {
	data, err := c.dec.DecodePacket()
	if err != nil {
		return nil, err
	}

	pw := packet.NewReader(data)
	var id int32
	if err = pw.VarInt(&id); err != nil {
		return nil, fmt.Errorf("%v decoding packet id", err)
	}

	pk := c.pool.Get(id)
	if pk == nil {
		return packet.Unknown{Id: id, Payload: data}, nil
	}

	if err = pk.Decode(pw); err != nil {
		return nil, fmt.Errorf("%v decoding packet contents for %#v", err, pk)
	}

	return pk, nil
}

func (c *Conn) DecodePacket(pk packet.Packet) error {
	payload, err := c.dec.DecodePacket()
	if err != nil {
		return err
	}

	rd := packet.NewReader(payload)

	var id int32
	if err = rd.VarInt(&id); err != nil {
		return fmt.Errorf("%v decoding packet id", err)
	}

	if id != pk.ID() {
		return fmt.Errorf("unexpected packet ID %x received. Expected packet ID to be %x", id, pk.ID())
	}

	if err = pk.Decode(rd); err != nil {
		return fmt.Errorf("%v decoding packet contents for %+v", err, pk)
	}

	return nil
}

// SendPacket writes and immediately sends the packet.
// Use for critical information. overusing can cause
// more latency and bandwidth to be used.
func (c *Conn) SendPacket(pk packet.Packet) error {
	c.encMu.Lock()
	defer c.encMu.Unlock()

	if err := c.writePacket(pk); err != nil {
		return fmt.Errorf("%w writing %+v to buffer", err, pk)
	}

	data := c.enc.Flush()
	if _, err := c.tcpCn.Write(data); err != nil {
		return fmt.Errorf("%w sending packet %v", err, pk)
	}

	return nil
}

// WritePacket writes the packet to a buffer.
// Use for situations where packets don't need to be sent IMMEDIATELY.
// Chat messages etc.
// Can also be used to improve bandwidth and client side latency by sending all the data at once.
// Just make sure it's done in a timely way
func (c *Conn) WritePacket(pk packet.Packet) error {
	c.encMu.Lock()
	defer c.encMu.Unlock()
	return c.writePacket(pk)
}

func (c *Conn) writePacket(pk packet.Packet) error {
	start := c.buf.Len() //records the start of the packet data
	pw := packet.NewWriter(c.buf)

	//ignore errors since writing to a bytes.Buffer object
	//always returns nil
	_ = pw.VarInt(pk.ID())
	_ = pk.Encode(pw)

	end := c.buf.Len()

	c.buf.Truncate(start)

	return c.enc.EncodePacket(c.buf.Bytes()[start:end])
}

func (c *Conn) FlushPackets() error {
	c.encMu.Lock()
	defer c.encMu.Unlock()

	if _, err := c.tcpCn.Write(c.enc.Flush()); err != nil {
		return fmt.Errorf("%w writing packets", err)
	}

	return nil
}

func (c *Conn) enableEncryption(sharedSecret []byte) error {
	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return err
	}

	c.dec.EnableDecryption(block, sharedSecret)
	c.enc.EnableEncryption(block, sharedSecret)
	return nil
}

func (c *Conn) enableCompression(threshold int32) {
	c.dec.EnableDecompression()
	c.enc.EnableCompression(int(threshold))
}

func (c *Conn) Close(err error) {
	if err != nil {
		fmt.Printf("%v: Connection Closed: %v\n", c.tcpCn.RemoteAddr(), err)
	}

	c.tcpCn.Close()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.tcpCn.RemoteAddr()
}
