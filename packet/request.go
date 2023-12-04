package packet

type Request struct{}

func (r Request) ID() int32 {
	return 0x00
}

func (r Request) Encode(w *Writer) error {
	return nil
}

func (r Request) Decode(rd *Reader) error {
	return nil
}
