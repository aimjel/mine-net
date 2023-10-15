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
	c.pool = nopPool{}

	for {
		pk, err := c.ReadPacket()
		if err != nil {
			return err
		}

		//set compression packet
		if pk.ID() == 0x03 {
			var sc packet.SetCompression

			uk := pk.(packet.Unknown)
			rd := packet.NewReader(uk.Payload)
			var id int32
			_ = rd.VarInt(&id) //reads the id
			sc.Decode(rd)

			c.enableCompression(sc.Threshold)
			fmt.Println("set compression packet", sc.Threshold)
			break
		}
	}
	return nil
}
