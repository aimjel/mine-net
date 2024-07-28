package protocol

import "testing"

func TestPutBuffer(t *testing.T) {
	buf := GetBuffer(1024 * 30)
	PutBuffer(buf)
}

//TODO: re write tests

//var hs = &packet.Handshake{
//	ProtocolVersion: 758,
//	ServerAddress:   "localhost",
//	ServerPort:      25565,
//	NextState:       1,
//}
//
//var jg = &packet.JoinGame{}
//
//var enc = func() *Encoder {
//	e := NewEncoder()
//	e.EnableCompression(20)
//	return e
//}()
//
//func TestEncoder_Encode(t *testing.T) {
//	for i := 0; i < 2; i++ {
//		if err := enc.EncodePacket(hs); err != nil {
//			t.Fatal(err)
//		}
//	}
//	t.Log(enc.Flush())
//}
//
//func BenchmarkEncoder_EncodePacket(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		enc.EncodePacket(jg)
//		b.StopTimer()
//		enc.Flush()
//		b.StartTimer()
//	}
//}
//
//func TestDecoder_DecodePacket(t *testing.T) {
//	t.Run("TestEncoder_Encode", TestEncoder_Encode)
//
//	enc := NewEncoder()
//	enc.EnableCompression(15)
//
//	for i := 0; i < 1; i++ {
//		if err := enc.EncodePacket(hs); err != nil {
//			t.Fatal(err)
//		}
//	}
//
//	dec := NewDecoder(bytes.NewBuffer(enc.Flush()))
//	dec.EnableDecompression()
//
//	for i := 0; i < 1; i++ {
//		pk, err := dec.DecodePacket()
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		t.Logf("decoded packet %v\n", pk)
//	}
//}
