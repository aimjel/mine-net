package minecraft

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aimjel/minecraft/protocol/types"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/aimjel/minecraft/packet"
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
	CompressionThreshold int32

	Messages *Messages

	//todo add more config fields
}

func (lc *ListenConfig) Listen(address string) (*Listener, error) {
	if lc.Messages == nil {
		lc.Messages = &DefaultMessages
	}
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
		key, err = rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			return nil, err
		}
	}

	l := &Listener{
		cfg:   *lc,
		tcpLn: ln,
		key:   key,

		await: make(chan *Conn, 4),
	}

	//starts listening for incoming connections
	go l.listen()

	return l, nil
}

type Listener struct {
	cfg ListenConfig

	tcpLn *net.TCPListener

	key *rsa.PrivateKey

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

	var pk packet.Handshake
	var err = c.DecodePacket(&pk)

	switch pk.NextState {

	case 0x01: //status
		if l.cfg.Status != nil {
			if err = l.cfg.Status.handleStatus(c); err != nil {
				c.Close(fmt.Errorf("%v while handling status", err))
			}
		}

	case 0x02:
		if pk.ProtocolVersion != int32(l.cfg.Status.s.Version.Protocol) {
			reason := l.cfg.Messages.ProtocolTooNew
			if pk.ProtocolVersion < int32(l.cfg.Status.s.Version.Protocol) {
				reason = l.cfg.Messages.ProtocolTooOld
			}

			err = c.SendPacket(&packet.DisconnectLogin{
				Reason: chat.NewMessage(reason),
			})
			break
		}

		if err = l.handleLogin(c); err == nil {
			if x := l.cfg.CompressionThreshold; x >= 0 {
				if err = c.SendPacket(&packet.SetCompression{Threshold: x}); err != nil {
					err = fmt.Errorf("%v sending set compression", err)
					break
				}

				c.enableCompression(x)
			}

			if err = c.SendPacket(&packet.LoginSuccess{
				Name:       c.name,
				UUID:       c.uuid,
				Properties: c.properties,
			}); err != nil {
				err = fmt.Errorf("%v sending login success", err)
				break
			}

			c.Pool = &basicPool{}
			l.await <- c
			return //return so it doesn't close the connection
		}
		err = fmt.Errorf("%v handling login state", err)
	}

	c.Close(err)
}

func (l *Listener) handleLogin(c *Conn) error {
	var ls packet.LoginStart
	if err := c.DecodePacket(&ls); err != nil {
		return err
	}

	if l.key == nil {
		var uuid [16]byte
		newUUIDv3(ls.Name, uuid[:])
		c.name, c.uuid = ls.Name, uuid
		return nil
	}

	key, err := x509.MarshalPKIXPublicKey(&l.key.PublicKey)
	if err != nil {
		return err
	}

	token := make([]byte, 8)
	_, _ = rand.Read(token)

	if err = c.SendPacket(&packet.EncryptionRequest{PublicKey: key, VerifyToken: token}); err != nil {
		return err
	}

	var encryptResp packet.EncryptionResponse
	if err = c.DecodePacket(&encryptResp); err != nil {
		return err
	}

	var sharedSecret, verifyToken []byte

	if sharedSecret, err = l.key.Decrypt(nil, encryptResp.SharedSecret, nil); err != nil {
		return err
	}

	if verifyToken, err = l.key.Decrypt(nil, encryptResp.VerifyToken, nil); err != nil {
		return err
	}

	if !bytes.Equal(verifyToken, token) {
		return fmt.Errorf("failed to verify token")
	}

	if err := c.enableEncryption(sharedSecret); err != nil {
		return err
	}

	loginHash, err := l.generateHash(sharedSecret)
	if err != nil {
		return err
	}

	r, err := http.DefaultClient.Get("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=" + ls.Name + "&serverId=" + loginHash)
	if err != nil {
		return fmt.Errorf("%v getting player data", err)
	}

	var data struct {
		Id         string           `json:"id"`
		Name       string           `json:"name"`
		Properties []types.Property `json:"properties"`
	}

	if err = json.NewDecoder(r.Body).Decode(&data); err != nil && err != io.EOF {
		return err
	}
	_ = r.Body.Close()

	uuid, err := hex.DecodeString(data.Id)
	if err != nil {
		return err
	}

	c.name, c.properties = data.Name, data.Properties

	if n := copy(c.uuid[:], uuid); n != 16 {
		c.SendPacket(&packet.DisconnectLogin{Reason: chat.NewMessage(l.cfg.Messages.OnlineMode)})
		return fmt.Errorf("offline player on online server")
	}
	return nil
}

func (l *Listener) Accept() (*Conn, error) {
	c, ok := <-l.await
	if !ok {
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
		twosComplement(loginHash)
	}

	hs := strings.TrimLeft(hex.EncodeToString(loginHash), "0")
	if neg {
		hs = "-" + hs
	}

	return hs, nil
}

func twosComplement(p []byte) {
	//invert all the bites
	for k, v := range p {
		p[k] = ^v
	}

	// Add 1
	carry := byte(1)
	for i := len(p) - 1; i >= 0; i-- {
		p[i] += carry
		carry = p[i] >> 8
		p[i] &= 0xFF
		if carry == 0 {
			break
		}
	}
}

func newUUIDv3(name string, out []byte) {
	h := md5.New()
	h.Write([]byte("OfflinePlayer:" + name))
	id := h.Sum(nil)

	id[6] = (id[6] & 0x0f) | uint8((3&0xf)<<4)
	id[8] = (id[8] & 0x3f) | 0x80 // RFC 4122 variant

	copy(out, id)
}
