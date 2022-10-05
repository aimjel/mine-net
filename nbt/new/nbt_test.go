package new

import (
	_ "embed"
	"testing"
)

//go:embed testdata/bigTest.nbt
var bigTestBytes []byte

func TestUnmarshal(t *testing.T) {
	if err := Unmarshal(bigTestBytes, &bigTest); err != nil {
		t.Fatal(err)
	}
}

var bigTest struct {
	Level struct {
		LongTest           int64   `nbt:"longTest"`
		ShortTest          int16   `nbt:"shortTest"`
		StringTest         string  `nbt:"stringTest"`
		FloatTest          float32 `nbt:"floatTest"`
		IntTest            int32   `nbt:"intTest"`
		NestedCompoundTest struct {
			Ham struct {
				Name  string  `nbt:"name"`
				Value float32 `nbt:"value"`
			} `nbt:"ham"`
			Egg struct {
				Name  string  `nbt:"name"`
				Value float32 `nbt:"value"`
			} `nbt:"egg"`
		} `nbt:"nested compound test"`

		ListTestLong []int64 `nbt:"listTest (long)"`

		ListTestCompound []struct {
			Name      string `nbt:"name"`
			CreatedOn int64  `nbt:"created-on"`
		} `nbt:"listTest (compound)"`

		ByteTest int8 `nbt:"byteTest"`

		DoubleTest float64 `nbt:"doubleTest"`
	}
}
