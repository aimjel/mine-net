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
	"io"
	"log"
	"minecraft/packet"
	"minecraft/player"
	"net"
	"net/http"
	"strings"
)

type Listener struct {
	tcpLn *net.TCPListener

	status *Status

	key *rsa.PrivateKey

	await chan *Conn
}

// Listen acts like net.ListenTCP for minecraft servers.
func Listen(address string) (*Listener, error) {
	ln, err := net.Listen("tcp4", address)
	if err != nil {
		return nil, err
	}

	key, _ := rsa.GenerateKey(rand.Reader, 1024)

	listener := &Listener{
		tcpLn:  ln.(*net.TCPListener),
		status: NewStatus(756, 1, nil, "&aA Minecraft Server written in Go!"),
		key:    key,
		await:  make(chan *Conn, 255),
	}

	go listener.listen()

	return listener, nil
}

// listen accepts incoming connections
func (l *Listener) listen() {
	for {
		c, err := l.tcpLn.AcceptTCP()

		if err != nil {
			log.Printf("%v: %v\n", l.tcpLn.Addr(), err)
			close(l.await)
			break
		}

		go l.newConn(c)
	}
}

func (l *Listener) newConn(tcp *net.TCPConn) {
	c := newConn(tcp)

	var hs packet.Handshake
	if err := c.DecodePacket(&hs); err != nil {
		//todo log err close connection?
		return
	}

	switch hs.NextState {

	case 0x01: //status
		if err := l.handleStatus(c); err != nil {
			log.Printf("%v: %v during status state", c.RemoteAddr(), err)
		}

		c.Close()

	case 0x02: //login
		if err := l.handleLogin(c); err != nil {
			log.Printf("%v: %v during login state", c.RemoteAddr(), err)
			c.Close()
			return
		}

		if err := c.WritePacket(&packet.LoginSuccess{UUID: c.Info.UUID, Username: c.Info.Name}); err != nil {
			log.Printf("%v: %v sending login sucess\n", c.RemoteAddr(), err)
			c.Close()
			return
		}

		c.pool = serverBoundPlayPool
		l.await <- c
		return
	}
}

func (l *Listener) handleStatus(c *Conn) error {
	var r packet.Request
	if err := c.DecodePacket(&r); err != nil {
		return err
	}

	data, err := json.Marshal(l.status)
	if err != nil {
		return fmt.Errorf("%v marshaling status", err)
	}

	if err = c.WritePacket(&packet.Response{JSON: string(data)}); err != nil {
		return fmt.Errorf("%v writing response packet", err)
	}

	var ping packet.Ping
	if err = c.DecodePacket(&ping); err != nil {
		return err
	}

	return c.WritePacket(&packet.Pong{Payload: ping.Payload})
}

func (l *Listener) handleLogin(c *Conn) error {
	var ls packet.LoginStart
	if err := c.DecodePacket(&ls); err != nil {
		return err
	}

	if l.key == nil {
		var uuid [16]byte
		_, _ = rand.Read(uuid[:])
		c.Info = &player.Info{UUID: uuid, Name: ls.Name}
		return nil
	}

	key, err := x509.MarshalPKIXPublicKey(&l.key.PublicKey)
	if err != nil {
		return err
	}

	token := make([]byte, 8, 8)
	_, _ = rand.Read(token)

	if err = c.WritePacket(&packet.EncryptionRequest{PublicKey: key, VerifyToken: token}); err != nil {
		return err
	}

	var encryptResp packet.EncryptionResponse
	if err = c.DecodePacket(&encryptResp); err != nil {
		return err
	}

	var (
		sharedSecret, verifyToken []byte
	)
	if sharedSecret, err = l.key.Decrypt(nil, encryptResp.SharedSecret, nil); err != nil {
		return err
	}

	if verifyToken, err = l.key.Decrypt(nil, encryptResp.VerifyToken, nil); err != nil {
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

	r, err := http.DefaultClient.Get("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=" + ls.Name + "&serverId=" + loginHash)
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
	return nil
}

func (l *Listener) Accept() (*Conn, error) {
	conn, ok := <-l.await

	if !ok {
		return nil, net.ErrClosed
	}

	return conn, nil
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
