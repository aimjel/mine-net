package nbt

import (
	"bytes"
	_ "embed"
	"os"
	"testing"
)

//go:embed testdata/bigTest.nbt
var bigTestData []byte

//go:embed testdata/chunk.nbt
var chunkData []byte

type bigTest struct {
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

type chunk struct {
	Level struct {
		Status        string
		ZPos          int32 `nbt:"zPos"`
		LastUpdate    int64
		Biomes        []int32
		InhabitedTime int64
		XPos          int32 `nbt:"xPos"`
		Heightmaps    struct {
			OceanFloor             []int64 `nbt:"OCEAN_FLOOR"`
			MotionBlockingNoLeaves []int64 `nbt:"MOTION_BLOCKING_NO_LEAVES"`
			MotionBlocking         []int64 `nbt:"MOTION_BLOCKING"`
			WorldSurface           []int64 `nbt:"WORLD_SURFACE"`
		}
		IsLightOn int8 `nbt:"isLightOn"`
		Sections  []struct {
			Y           int8
			BlockStates []int64
			Palette     []struct {
				Name       string
				Properties map[string]string
			}
		}
	}

	DataVersion int32
}

func TestUnmarshal_chunk(t *testing.T) {
	var c chunk
	if err := Unmarshal(chunkData, &c); err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", c)
}

func TestUnmarshal_bigtest(t *testing.T) {
	var bg bigTest
	if err := Unmarshal(bigTestData, &bg); err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", bg)
}

func BenchmarkUnmarshalChunk(b *testing.B) {
	var c chunk
	for i := 0; i < b.N; i++ {
		if err := Unmarshal(chunkData, &c); err != nil {
			b.Fatal(err)
		}
	}
	b.ReportAllocs()
}

func BenchmarkUnmarshalBigTest(b *testing.B) {
	var bgTest bigTest
	for i := 0; i < b.N; i++ {
		if err := Unmarshal(bigTestData, &bgTest); err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
}

func TestScanner(t *testing.T) {
	err := checkValid(bigTestData)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkScanner(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = checkValid(bigTestData)
	}
}

func TestEncoder(t *testing.T) {
	b := bytes.NewBuffer(nil)
	hm := struct {
		Heightmaps struct {
			MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
			WorldSurface   []int64 `nbt:"WORLD_SURFACE"`
		}
	}{}
	if err := NewEncoder(b).Encode(hm); err != nil {
		t.Fatal(err)
	}

	os.WriteFile("heightMap.nbt", b.Bytes(), 0666)
}

func BenchmarkEncoder(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 1024*1024))

	enc := NewEncoder(buf)
	for i := 0; i < b.N; i++ {

		_ = enc.Encode(struct {
			Heightmaps struct {
				MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
				WorldSurface   []int64 `nbt:"WORLD_SURFACE"`
			}
		}{})

		buf.Reset()
	}

	b.ReportAllocs()
}
