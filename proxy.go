package minecraft

import (
	"errors"
	"fmt"
	"github.com/aimjel/minecraft/packet"
	"net"
)

type ProxyConfig struct {

	// Status displayed for the proxy server.
	// If value is nil, the proxy server will use the target servers status.
	Status *Status

	// OnReceive called when the client receives a packet from the server
	// Returning false drops the packet, not sending it to the client.
	OnReceive func(conn *Conn, pk packet.Packet) bool

	// OnSend called when the client sends a packet to the server
	// Returning false drops the packet, not sending it to the server.
	OnSend func(conn *Conn, pk packet.Packet) bool
}

func (cfg *ProxyConfig) Listen(addr, targetAddr string) error {
	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				fmt.Println(err)

				if errors.Is(err, net.ErrClosed) {
					return
				}
			}

			go cfg.handleConn(newConn(c.(*net.TCPConn)), targetAddr)
		}
	}()

	return nil
}

type proxyConn struct {
	// dialConn is connected to the target server.
	// Packets written by the target server
	// are read from this connection.
	dialConn *Conn

	// conn is the original server which originated
	//from the proxy server. Also known as the client.
	conn *Conn

	targetAddr string
}

func (cfg *ProxyConfig) handleConn(c *Conn, targetAddr string) {
	proxyC := &proxyConn{
		conn:       c,
		targetAddr: targetAddr,
	}

	cfg.start(proxyC)
}

func (cfg *ProxyConfig) start(pc *proxyConn) {
	var hs packet.Handshake
	if err := pc.conn.DecodePacket(&hs); err != nil {
		fmt.Println("decode hs", err)
		return
	}

	//addr, port, _ := net.SplitHostPort(pc.targetAddr)
	//portNum, _ := strconv.Atoi(port)

	//hs.ServerAddress = addr
	//hs.ServerPort = uint16(portNum)

	switch hs.NextState {

	//status
	case 0x01:
		if cfg.Status != nil {
			if err := cfg.Status.handleStatus(pc.conn); err != nil {
				err = fmt.Errorf("%v handling status state", err)
			}
			pc.conn.Close(nil)
			return
		}
		dialC, err := net.Dial("tcp4", pc.targetAddr)
		if err != nil {
			fmt.Println("dial", err)
			return
		}

		pc.dialConn = newConn(dialC.(*net.TCPConn))
		_ = pc.dialConn.SendPacket(&hs)
		cfg.proxy(pc)
		pc.conn.Close(nil)
		pc.dialConn.Close(nil)
		return

	case 0x02:
		dialC, err := net.Dial("tcp4", pc.targetAddr)
		if err != nil {
			fmt.Println(err)
			return
		}

		pc.dialConn = newConn(dialC.(*net.TCPConn))
		_ = pc.dialConn.SendPacket(&hs)

		if err = cfg.handleLogin(pc); err != nil {
			fmt.Printf("%v trying to login", err)
			return
		}

		cfg.proxy(pc)
	}
}

func (cfg *ProxyConfig) handleLogin(pc *proxyConn) error {
	var ls packet.LoginStart
	if err := pc.conn.DecodePacket(&ls); err != nil {
		return err
	}
	_ = pc.dialConn.SendPacket(&ls)

	for {
		pack, err := pc.dialConn.ReadPacket()
		if err != nil {
			return err
		}

		pk := pack.(packet.Unknown)
		switch pk.ID() {

		//encryption request
		case 0x01:
			pc.dialConn.Close(nil)
			pc.conn.Close(nil)
			return fmt.Errorf("online mode is not supported")

		case 0x02:
			var lgSuc packet.LoginSuccess
			if err = lgSuc.Decode(packet.NewReader(pk.Payload)); err != nil {
				return fmt.Errorf("%v decoding login success packet", err)
			}

			pc.dialConn.name = lgSuc.Name
			pc.dialConn.uuid = lgSuc.UUID
			pc.dialConn.properties = lgSuc.Properties

			pc.conn.name = lgSuc.Name
			pc.conn.uuid = lgSuc.UUID
			pc.conn.properties = lgSuc.Properties

			pc.conn.SendPacket(&lgSuc)
			return nil

		case 0x03:
			var com packet.SetCompression
			if err = com.Decode(packet.NewReader(pk.Payload)); err != nil {
				return fmt.Errorf("%v decoding compression threshold packet", err)
			}
			pc.conn.SendPacket(&com)

			pc.dialConn.enableCompression(com.Threshold)
			pc.conn.enableCompression(com.Threshold)
		}
	}
}

func (cfg *ProxyConfig) proxy(pc *proxyConn) {
	go func() {
		for {
			pk, err := pc.dialConn.ReadPacket()
			if err != nil {
				return
			}

			if cfg.OnReceive != nil {
				cfg.OnReceive(pc.dialConn, pk)
			}

			if err = pc.conn.SendPacket(pk); err != nil {
				return
			}

		}
	}()

	for {
		pk, err := pc.conn.ReadPacket()
		if err != nil {
			return
		}

		if cfg.OnSend != nil {
			cfg.OnSend(pc.conn, pk)
		}

		if err = pc.dialConn.SendPacket(pk); err != nil {
			return
		}
	}
}
