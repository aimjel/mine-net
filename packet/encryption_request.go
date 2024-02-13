package packet

import "github.com/aimjel/minecraft/protocol/encoding"

type EncryptionRequest struct {
	ServerID    string
	PublicKey   []byte
	VerifyToken []byte
}

func (r EncryptionRequest) ID() int32 {
	return 0x01
}

func (r *EncryptionRequest) Decode(rd *encoding.Reader) error {
	_ = rd.String(&r.ServerID)
	_ = rd.ByteArray(&r.PublicKey)
	return rd.ByteArray(&r.VerifyToken)
}

func (r EncryptionRequest) Encode(w *encoding.Writer) error {
	_ = w.String(r.ServerID)
	_ = w.ByteArray(r.PublicKey)
	return w.ByteArray(r.VerifyToken)
}
