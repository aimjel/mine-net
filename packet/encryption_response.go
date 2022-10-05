package packet

type EncryptionResponse struct {
	SharedSecret []byte
	VerifyToken  []byte
}

func (e EncryptionResponse) ID() int32 {
	return 0x01
}

func (e *EncryptionResponse) Decode(r *Reader) error {
	_ = r.ByteArray(&e.SharedSecret)
	return r.ByteArray(&e.VerifyToken)
}

func (e EncryptionResponse) Encode(w Writer) error {
	_ = w.ByteArray(e.SharedSecret)
	return w.ByteArray(e.VerifyToken)
}
