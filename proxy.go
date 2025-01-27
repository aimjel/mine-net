package minenet

import (
	"errors"
	"fmt"
	"github.com/aimjel/minenet/packet"
	"github.com/aimjel/minenet/protocol/encoding"
	"net"
)

type ProxyConfig struct {
	// OnReceive called when a packet is received from the client or server.
	// Returning false will drop the packet.
	OnReceive func(conn *Conn, pk packet.Packet, fromServer bool, state int) bool

	//ErrCh receives errors from the client and server
	ErrCh chan ProxyError
}

type ProxyListener struct {
	cfg *ProxyConfig

	ln *net.TCPListener
}

type ProxyError struct {
	//State is the protocol state the error occurred in
	State int

	Addr net.Addr

	Err error
}

func (p ProxyError) Error() string {
	return fmt.Sprintf("[%v] %v: %v", p.State, p.Addr, p.Err)
}

func (cfg *ProxyConfig) Listen(addr, targetAddr string) (*ProxyListener, error) {
	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				if errors.Is(err, net.ErrClosed) {
					fmt.Println(err)
					return
				}
			}

			go cfg.handleConn(newConn(c.(*net.TCPConn)), targetAddr)
		}
	}()

	return &ProxyListener{cfg: cfg, ln: ln.(*net.TCPListener)}, nil
}

func (l *ProxyListener) Close() error {
	return l.ln.Close()
}

func (cfg *ProxyConfig) handleConn(client *Conn, targetAddr string) {
	var hs packet.Handshake
	if err := client.DecodePacket(&hs); err != nil {
		cfg.ErrCh <- ProxyError{State: 0, Addr: client.RemoteAddr(), Err: err}
		client.Close(nil)
		return
	}
	if hs.NextState != 1 && hs.NextState != 2 && hs.ProtocolVersion != 763 {
		client.Close(nil)
		return
	}

	dialC, err := net.Dial("tcp4", targetAddr)
	if err != nil {
		cfg.ErrCh <- ProxyError{State: 0, Addr: client.RemoteAddr(), Err: err}
		client.Close(nil)
		return
	}

	serverConn := newConn(dialC.(*net.TCPConn))
	if err = serverConn.SendPacket(&hs); err != nil {
		cfg.ErrCh <- ProxyError{State: 0, Addr: serverConn.RemoteAddr(), Err: err}
		client.Close(nil)
		return
	}

	switch hs.NextState {

	case 0x01:
		cfg.proxy(client, serverConn)
		client.Close(nil)
		serverConn.Close(nil)
		return

	case 0x02:
		if err = cfg.handleLogin(client, serverConn); err != nil {
			cfg.ErrCh <- err.(ProxyError)
			client.Close(nil)
			return
		}

		cfg.proxy(client, serverConn)
	}
}

func (cfg *ProxyConfig) handleLogin(client, serverConn *Conn) error {
	var ls packet.LoginStart
	if err := client.DecodePacket(&ls); err != nil {
		return ProxyError{State: 2, Addr: client.RemoteAddr(), Err: err}
	}
	if err := serverConn.SendPacket(&ls); err != nil {
		return ProxyError{State: 2, Addr: client.RemoteAddr(), Err: err}
	}

	for {
		pack, err := serverConn.ReadPacket()
		if err != nil {
			return ProxyError{State: 2, Addr: serverConn.RemoteAddr(), Err: err}
		}

		pk := pack.(packet.Unknown)
		switch pk.ID() {

		//encryption request
		case 0x01:
			serverConn.Close(nil)
			client.Close(nil)
			return ProxyError{State: 2, Addr: client.RemoteAddr(), Err: fmt.Errorf("online mode is not supported")}

		case 0x02:
			var lgSuc packet.LoginSuccess
			if err = lgSuc.Decode(encoding.NewReader(pk.Payload)); err != nil {
				return ProxyError{State: 2, Addr: client.RemoteAddr(), Err: err}
			}

			serverConn.name, client.name = lgSuc.Name, lgSuc.Name
			serverConn.uuid, client.uuid = lgSuc.UUID, lgSuc.UUID
			serverConn.properties, client.properties = lgSuc.Properties, client.properties

			if err = client.SendPacket(&lgSuc); err != nil {
				return ProxyError{State: 2, Addr: client.RemoteAddr(), Err: err}
			}
			return nil

		case 0x03:
			var com packet.SetCompression
			if err = com.Decode(encoding.NewReader(pk.Payload)); err != nil {
				return ProxyError{State: 2, Addr: client.RemoteAddr(), Err: err}
			}
			if err = client.SendPacket(&com); err != nil {
				return ProxyError{State: 2, Addr: client.RemoteAddr(), Err: err}
			}

			serverConn.enableCompression(com.Threshold)
			client.enableCompression(com.Threshold)
		}
	}
}

func (cfg *ProxyConfig) proxy(client, serverConn *Conn) {
	go func(srv, client *Conn) {

		//reads from the server and forwards to the client
		for {
			pk, err := srv.ReadPacket()
			if err != nil {
				return
			}

			if cfg.OnReceive != nil {
				if !cfg.OnReceive(srv, pk, true, 3) {
					continue
				}
			}

			if err = client.SendPacket(pk); err != nil {
				return
			}

		}
	}(serverConn, client)

	//reads from the client and forwards to the target server
	for {
		pk, err := client.ReadPacket()
		if err != nil {
			return
		}

		if cfg.OnReceive != nil {
			if !cfg.OnReceive(client, pk, false, 3) {
				continue
			}
		}

		if err = serverConn.SendPacket(pk); err != nil {
			return
		}
	}
}
