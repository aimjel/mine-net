package minecraft

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aimjel/minecraft/packet"
	"github.com/aimjel/minecraft/player"
	"github.com/aimjel/minecraft/protocol"
	"io"
	"net"
	"net/http"
	"strings"
)

type ListenConfig struct {

	// Status handles the information showed to the client on the server list
	// which includes description, favicon, online/max players and protocol version and name
	Status *Status

	// OnlineMode enables server side encryption.
	// cracked accounts will not be able to connect when online mode is true.
	OnlineMode bool

	// CompressionThreshold compresses packets when they exceed n bytes.
	//-1 disables compression
	// 0 compresses everything
	CompressionThreshold int
	//todo add more config fields
}

func (lc *ListenConfig) Listen(address string) (*Listener, error) {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}

	ln, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		return nil, err
	}

	var key *rsa.PrivateKey
	if lc.OnlineMode {
		key, _ = rsa.GenerateKey(rand.Reader, 1024)
	}

	l := &Listener{
		tcpLn:                ln,
		key:                  key,
		compressionThreshold: lc.CompressionThreshold,
		status:               lc.Status,

		await: make(chan *Conn, 4),
	}

	//starts listening for incoming connections
	go l.listen()

	return l, nil
}

type Listener struct {
	tcpLn *net.TCPListener

	key *rsa.PrivateKey

	status *Status

	compressionThreshold int

	err error

	await chan *Conn
}

func (l *Listener) listen() {
	for {
		c, err := l.tcpLn.AcceptTCP()
		if err != nil {
			l.err = err
			close(l.await)
			return
		}

		go l.handle(c)
	}
}

// handle new connections
func (l *Listener) handle(conn *net.TCPConn) {
	c := newConn(conn)
	pk, err := c.ReadPacket()
	if err != nil {
		c.Close(err)
		return
	}

	if hs, ok := pk.(*packet.Handshake); ok {

		switch hs.NextState {

		case 0x01: //status
			if err = l.handleStatus(c); err != nil && l.status != nil {
				c.Close(fmt.Errorf("%v while handling status", err))
			}

			c.Close(nil)

		case 0x02:
			if err = l.handleLogin(c); err != nil {
				c.Close(fmt.Errorf("%v while handling login", err))
			}

			l.await <- c
		}
	} else {
		c.Close(fmt.Errorf("unknown initial packet"))
	}
}

func (l *Listener) handleStatus(c *Conn) error {
	c.pool = protocol.NewPool([]packet.Packet{
		0x00: &packet.Request{},
		0x01: &packet.Ping{},
	})
	for {
		p, err := c.ReadPacket()
		if err != nil {
			return err
		}

		switch pk := p.(type) {

		case *packet.Request:
			if err = c.WritePacket(&packet.Response{JSON: l.status.json()}, true); err != nil {
				return fmt.Errorf("%v writing response packet", err)
			}

		case *packet.Ping:
			return c.WritePacket(&packet.Pong{Payload: pk.Payload}, true)
		}
	}
}

func (l *Listener) handleLogin(c *Conn) error {
	c.pool = protocol.NewPool([]packet.Packet{
		0x00: &packet.LoginStart{},
		0x01: &packet.EncryptionResponse{},
	})

	token := make([]byte, 8)
	for {
		p, err := c.ReadPacket()
		if err != nil {
			return err
		}

		switch pk := p.(type) {

		case *packet.LoginStart:
			var uuid [16]byte
			_, _ = rand.Read(uuid[:])
			c.Info = &player.Info{UUID: uuid, Name: pk.Name}

			if l.key == nil {
				return nil
			}

			key, err := x509.MarshalPKIXPublicKey(&l.key.PublicKey)
			if err != nil {
				return err
			}

			_, _ = rand.Read(token)

			if err = c.WritePacket(&packet.EncryptionRequest{PublicKey: key, VerifyToken: token}, true); err != nil {
				return err
			}

		case *packet.EncryptionResponse:
			var (
				sharedSecret, verifyToken []byte
			)
			if sharedSecret, err = l.key.Decrypt(nil, pk.SharedSecret, nil); err != nil {
				return err
			}

			if verifyToken, err = l.key.Decrypt(nil, pk.VerifyToken, nil); err != nil {
				return err
			}

			if bytes.Equal(verifyToken, token) == false {
				//todo send disconnect packet
				return fmt.Errorf("failed to verify token")
			}

			block, err := aes.NewCipher(sharedSecret)
			if err != nil {
				return err
			}
			c.enableEncryption(block, sharedSecret)

			loginHash, err := l.generateHash(sharedSecret)
			if err != nil {
				return err
			}

			r, err := http.DefaultClient.Get("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=" + c.Info.Name + "&serverId=" + loginHash)
			if err != nil {
				return fmt.Errorf("%v getting player data", err)
			}

			var data struct {
				Id         string `json:"id"`
				Name       string `json:"name"`
				Properties []struct {
					Name      string `json:"name"`
					Value     string `json:"value"`
					Signature string `json:"signature"`
				} `json:"properties"`
			}

			if err = json.NewDecoder(r.Body).Decode(&data); err != nil && err != io.EOF {
				return err
			}
			_ = r.Body.Close()

			uuid, err := hex.DecodeString(data.Id)
			if err != nil {
				return err
			}

			c.Info = &player.Info{Name: data.Name, Properties: []struct {
				Name      string
				Value     string
				Signature string
			}(data.Properties)}

			if n := copy(c.Info.UUID[:], uuid); n != 16 {
				return fmt.Errorf("expected 16 bytes from uuid got %v", n)
			}

		}
	}
}

func (l *Listener) Accept() (*Conn, error) {
	c, ok := <-l.await

	if ok == false {
		if l.err != nil {
			return nil, l.err
		}

		return nil, net.ErrClosed
	}

	return c, nil
}

// generateHash generates the login hash sent in the HTTP Get to retrieve uuid, name, textures
func (l *Listener) generateHash(sharedSecret []byte) (string, error) {
	h := sha1.New()
	h.Write(sharedSecret)

	key, err := x509.MarshalPKIXPublicKey(&l.key.PublicKey)
	if err != nil {
		return "", err
	}

	h.Write(key)
	loginHash := h.Sum(nil)

	neg := loginHash[0] >= 128

	if neg {
		for k, v := range loginHash {
			loginHash[k] = ^v
		}

		loginHash[19] |= 0x01
	}

	hs := strings.TrimLeft(hex.EncodeToString(loginHash), "0")
	if neg {
		hs = "-" + hs
	}

	return hs, nil
}
