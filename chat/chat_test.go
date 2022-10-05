package chat

import (
	"encoding/json"
	"strings"
	"testing"
)

func BenchmarkNewMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := NewMessage("&1A&2 Mine&3craft &4&5&6Server written &7i&8n &9Go&0!")
		_, _ = json.Marshal(m)
	}

	b.ReportAllocs()
}

func BenchmarkMessagev2_MarshalJSON(b *testing.B) {
	//m := NewMessage2("&1A&2 Mine&3craft &4&5&6Server written &7i&8n &9Go&0!")

	for i := 0; i < b.N; i++ {
		m := NewMessage2("&1A&2 Mine&3craft &4&5&6Server written &7i&8n &9Go&0!")
		_, _ = m.MarshalJSON()
	}

	b.ReportAllocs()
}

func TestSplit(t *testing.T) {
	t.Logf("%#v\n", strings.Split("&a a minecraft server", "&"))
}
