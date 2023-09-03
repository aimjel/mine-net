package minecraft

import (
	"github.com/aimjel/minecraft/packet"
	"net"
	"strconv"
)

type Dialer struct {
	Username string
}

func (d *Dialer) Dial(address string) (*Conn, error) {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}

	tcp, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		return nil, err
	}

	mccon := newConn(tcp)
	if err = doLogin(mccon, address, d.Username); err != nil {
		return nil, err
	}

	return mccon, nil
}

func doLogin(c *Conn, address, username string) error {
	host, p, _ := net.SplitHostPort(address)
	port, _ := strconv.Atoi(p)
	if err := c.SendPacket(&packet.Handshake{
		ProtocolVersion: 763,
		ServerAddress:   host,
		ServerPort:      uint16(port),
		NextState:       2,
	}); err != nil {
		return err
	}

	if err := c.SendPacket(&packet.LoginStart{Name: username}); err != nil {
		return err
	}

	c.pool = clientLoginPool{}
	return nil
}
