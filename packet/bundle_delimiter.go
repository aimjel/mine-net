package packet

type BundleDelimiter struct{}

func (b BundleDelimiter) ID() int32 {
	return 0x00
}

func (b BundleDelimiter) Decode(r *Reader) error {
	return nil
}

func (b BundleDelimiter) Encode(w *Writer) error {
	return nil
}
