package nbt

import (
	"bytes"
	_ "embed"
	"os"
	"reflect"
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
	Status        string
	ZPos          int32 `nbt:"zPos"`
	LastUpdate    int64
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
		BlockStates struct {
			Data    []int64 `nbt:"data"`
			Palette []struct {
				Name       string
				Properties map[string]string
			} `nbt:"palette"`
		} `nbt:"block_states"`
	} `nbt:"sections"`

	DataVersion int32
}

func TestUnmarshalChunk(t *testing.T) {
	var c chunk
	if err := Unmarshal(chunkData, &c); err != nil {
		t.Fatal(err)
	}

	var c2 chunk
	err := NewDecoder(bytes.NewReader(chunkData)).Decode(&c2)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(c2, c) {
		t.Fatal("corrupt decoding")
	}
}

func TestUnmarshalBigTest(t *testing.T) {
	var bg bigTest
	if err := Unmarshal(bigTestData, &bg); err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", bg)
}

func BenchmarkUnmarshalChunk(b *testing.B) {
	var c chunk
	for i := 0; i < b.N; i++ {
		_ = Unmarshal(chunkData, &c)
	}
	b.ReportAllocs()
}

func BenchmarkDecoder_DecodeChunk(b *testing.B) {
	var c chunk
	rd := bytes.NewReader(chunkData)
	dec := NewDecoder(rd)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = dec.Decode(&c)

		b.StopTimer()
		rd.Reset(chunkData)
		b.StartTimer()
	}
	b.ReportAllocs()
}

func BenchmarkUnmarshalBigTest(b *testing.B) {
	var bgTest bigTest
	for i := 0; i < b.N; i++ {
		_ = Unmarshal(bigTestData, &bgTest)
	}
	b.ReportAllocs()
}

func TestScanner(t *testing.T) {
	if err := checkValid(chunkData); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkScanner(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = checkValid(bigTestData)
	}
}

func TestEncoder(t *testing.T) {
	var ck chunk
	_ = Unmarshal(chunkData, &ck)

	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(ck); err != nil {
		t.Fatal(err)
	}

	if err := checkValid(buf.Bytes()); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile("test.nbt", buf.Bytes(), 0666)
}

func BenchmarkEncoder(b *testing.B) {
	var c chunk
	_ = Unmarshal(chunkData, &c)

	buf := bytes.NewBuffer(make([]byte, 0, len(chunkData)))
	enc := NewEncoder(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = enc.Encode(c)
		buf.Reset()
	}

	b.ReportAllocs()
}
