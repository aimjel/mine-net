package packet

type PlayerSessionServer struct {
	SessionID [16]byte
  PublicKey PublicKey
}

type PublicKey struct {
  ExpiresAt int64
  PublicKey []byte
  KeySignature []byte
}

func (m PlayerSessionServer) ID() int32 {
	return 0x06
}

func (m *PlayerSessionServer) Decode(r *Reader) error {
	r.UUID(&m.SessionID)
  r.Int64(&m.PublicKey.ExpiresAt)
  r.ByteArray(&m.PublicKey.PublicKey)
  return r.ByteArray(&m.PublicKey.KeySignature)
}

func (m PlayerSessionServer) Encode(w Writer) error {
	w.UUID(m.SessionID)
  w.Int64(m.PublicKey.ExpiresAt)
  w.ByteArray(m.PublicKey.PublicKey)
  return w.ByteArray(m.PublicKey.KeySignature)
}
