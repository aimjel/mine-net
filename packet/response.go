package packet

type Response struct {
	JSON string
}

func (r Response) ID() int32 {
	return 0x00
}

func (r *Response) Decode(rd *Reader) error {
	return rd.String(&r.JSON)
}

func (r Response) Encode(w Writer) error {
	return w.String(r.JSON)
}
