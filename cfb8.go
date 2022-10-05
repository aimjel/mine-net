package minecraft

import (
	"crypto/cipher"
)

type cfb8 struct {
	b       cipher.Block
	iv      []byte
	out     []byte
	outUsed int

	decrypt bool
}

func (c *cfb8) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}

	for k, v := range src {
		copy(c.out, c.iv) //out is temp storage for the un-encrypted IV
		c.b.Encrypt(c.iv, c.iv)
		v = v ^ c.iv[0]

		copy(c.iv, c.out[1:]) //takes all but the first byte back into IV

		if c.decrypt {
			c.iv[15] = src[k]
		} else {
			c.iv[15] = v
		}

		dst[k] = v
	}
}

func newCFB(block cipher.Block, iv []byte, decrypt bool) cipher.Stream {
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		// stack trace will indicate whether it was de or encryption
		panic("cipher.newCFB: IV length must equal block size")
	}

	x := &cfb8{
		b:       block,
		out:     make([]byte, blockSize),
		iv:      make([]byte, blockSize),
		outUsed: blockSize,
		decrypt: decrypt,
	}
	copy(x.iv, iv)

	return x
}
