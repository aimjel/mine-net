package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type PlayerSessionServer struct {
	SessionID    [16]byte
	ExpiresAt    int64
	PublicKey    []byte
	KeySignature []byte
}

func (m PlayerSessionServer) ID() int32 {
	return 0x06
}

func (m *PlayerSessionServer) Decode(r *encoding.Reader) error {
	r.UUID(&m.SessionID)
	r.Int64(&m.ExpiresAt)
	r.ByteArray(&m.PublicKey)
	return r.ByteArray(&m.KeySignature)
}

func (m PlayerSessionServer) Encode(w *encoding.Writer) error {
	w.UUID(m.SessionID)
	w.Int64(m.ExpiresAt)
	w.ByteArray(m.PublicKey)
	return w.ByteArray(m.KeySignature)
}
