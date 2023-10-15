package packet

import (
	"bytes"
	"github.com/aimjel/minecraft/protocol/types"
	"reflect"
	"testing"
)

var pk = TestPacket{
	Boolean:       false,
	Byte:          -50,
	UnsignedByte:  255,
	Short:         -25565,
	UnsignedShort: 45402,
	Int:           245724852,
	Long:          9123172033854775830,
	Float:         9.91235342,
	Double:        6.91986579173569175,
	String:        "hello world!",
	VarInt:        756,
	VarLong:       00, //TODO not implemented
	UUID:          [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
	ByteArray:     []byte("hello world again!"),
}

var b bytes.Buffer

func TestPacketEncode_Decode(t *testing.T) {
	if err := pk.Encode(NewWriter(&b)); err != nil {
		t.Fatal(err)
	}

	var test_pk TestPacket
	if err := test_pk.Decode(NewReader(b.Bytes())); err != nil {
		t.Fatalf("%v decoding packet", err)
	}

	if !reflect.DeepEqual(test_pk, pk) {
		ar := reflect.ValueOf(pk)
		br := reflect.ValueOf(test_pk)

		for i := 0; i < ar.NumField(); i++ {
			f1 := ar.Field(i).Interface()
			f2 := br.Field(i).Interface()

			if !reflect.DeepEqual(f1, f2) {
				t.Logf("expected %v but got %v for type %T in field %v\n", f1, f2, f1, ar.Type().Field(i).Name)
			}
		}
		t.FailNow()
	}
}

func TestPlayerInfoUpdate_Encode(t *testing.T) {
	p := PlayerInfoUpdate{Actions: 63, Players: make([]types.PlayerInfo, 1)}
	if err := p.Encode(NewWriter(&b)); err != nil {
		t.Fatal(err)
	}
}
