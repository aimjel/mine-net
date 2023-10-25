package minecraft

import (
	"fmt"
	"github.com/aimjel/minecraft/packet"
	"net"
	"strconv"
)

type Dialer struct {
	Username string
}

func (d *Dialer) Dial(address string) (*Conn, error) {
	c, err := net.Dial("tcp4", address)
	if err != nil {
		return nil, err
	}

	mccon := newConn(c.(*net.TCPConn))
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

	for {
		pack, err := c.ReadPacket()
		if err != nil {
			return err
		}

		pk := pack.(packet.Unknown)
		switch pk.ID() {

		case 0x01:
			return fmt.Errorf("online mode is not supported")

		case 0x02:
			var lgSuccess packet.LoginSuccess
			if err = lgSuccess.Decode(packet.NewReader(pk.Payload)); err != nil {
				return err
			}

			c.name = lgSuccess.Name
			c.uuid = lgSuccess.UUID
			c.properties = lgSuccess.Properties
			return nil

		case 0x03:
			var com packet.SetCompression
			if err = com.Decode(packet.NewReader(pk.Payload)); err != nil {
				return err
			}

			c.enableCompression(com.Threshold)
		}
	}
}
