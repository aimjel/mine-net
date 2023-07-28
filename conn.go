package minecraft

import (
	"crypto/cipher"
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

	pool protocol.Pool

	Info *player.Info

	//encMu protects the enc from data races if two goroutines try to write a packet
	encMu sync.Mutex
}

func newConn(c *net.TCPConn) *Conn {
	return &Conn{
		tcpCn: c,

		dec: protocol.NewDecoder(c),

		enc: protocol.NewEncoder(c),

		pool: protocol.NewPool([]packet.Packet{
			0: &packet.Handshake{},
		}),
	}
}

func (c *Conn) ReadPacket() (packet.Packet, error) {
	data, err := c.dec.DecodePacket()
	if err != nil {
		return nil, fmt.Errorf("%v decoding packet", err)
	}

	pw := packet.NewReader(data)

	var id int32
	if err = pw.VarInt(&id); err != nil {
		return nil, fmt.Errorf("%v decoding packet id", err)
	}

	pk := c.pool.Get(id)
	if pk == nil {
		return nil, fmt.Errorf("packet %x is unknown", id)
	}

	if err = pk.Decode(pw); err != nil {
		return nil, fmt.Errorf("%v decoding packet contents", err)
	}

	return pk, nil
}

func (c *Conn) WritePacket(pk packet.Packet, forceSend bool) error {
	c.encMu.Lock()
	defer c.encMu.Unlock()

	wr := packet.NewWriter(c.enc)

	if err := wr.VarInt(pk.ID()); err != nil {
		return fmt.Errorf("%v encoding packet id", err)
	}

	if err := pk.Encode(wr); err != nil {
		return fmt.Errorf("%v encoding packet", err)
	}

	c.enc.Encode()

	if forceSend {
		return c.enc.Flush()
	}

	return nil
}

func (c *Conn) FlushPackets() error {
	c.encMu.Lock()
	defer c.encMu.Unlock()

	return c.enc.Flush()
}

func (c *Conn) enableEncryption(block cipher.Block, sharedSecret []byte) {
	c.dec.EnableDecryption(block, sharedSecret)
	c.enc.EnableEncryption(block, sharedSecret)
}

func (c *Conn) Close(err error) {
	if err != nil {
		fmt.Printf("%v: Connection Closed: %v\n", c.tcpCn.RemoteAddr(), err)
	}

	c.tcpCn.Close()
}
