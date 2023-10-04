package minecraft

import (
	"crypto/aes"
	"fmt"
	"github.com/aimjel/minecraft/packet"
	"github.com/aimjel/minecraft/player"
	"github.com/aimjel/minecraft/protocol"
	"net"
	"sync"
)

type Conn struct {
	tcpCn *net.TCPConn

	dec *protocol.Decoder

	enc *protocol.Encoder

	pool Pool

	Info *player.Info

	//encMu protects the Encoder from data races if two goroutines try to write a packet
	encMu sync.Mutex
}

func newConn(c *net.TCPConn) *Conn {
	return &Conn{
		tcpCn: c,

		dec: protocol.NewDecoder(c),

		enc: protocol.NewEncoder(),
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

	if err := c.enc.EncodePacket(pk); err != nil {
		return err
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
	return c.enc.EncodePacket(pk)
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
	if err := c.SendPacket(&packet.SetCompression{Threshold: threshold}); err != nil {
		panic(err)
	}

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
