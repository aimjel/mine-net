package packet

import "github.com/aimjel/minenet/protocol/encoding"

type BundleDelimiter struct{}

func (b BundleDelimiter) ID() int32 {
	return 0x00
}

func (b BundleDelimiter) Decode(r *encoding.Reader) error {
	return nil
}

func (b BundleDelimiter) Encode(w *encoding.Writer) error {
	return nil
}
