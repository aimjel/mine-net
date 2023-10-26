package packet

import "github.com/aimjel/minecraft/chat"

// tubby ahh packet!

type PlayerChatMessage struct {
	//Header
	Sender           [16]byte
	Index            int32
	MessageSignature []byte
	//Body
	Message   string
	Timestamp int64
	Salt      int64
	//Previous Messages
	PreviousMessages []PreviousMessage
	//Other
	UnsignedContent *chat.Message
	FilterType      int32
	FilterTypeBits  []int64
	//Network Target
	ChatType          int32
	NetworkName       chat.Message
	NetworkTargetName *chat.Message
}

type PreviousMessage struct {
	MessageID int32
	Signature []byte
}

func (m PlayerChatMessage) ID() int32 {
	return 0x35
}

func (m *PlayerChatMessage) Decode(r *Reader) error {
	// no way!!
	return NotImplemented
}

func (m PlayerChatMessage) Encode(w Writer) error {
	w.UUID(m.Sender)
	w.VarInt(m.Index)
	if m.MessageSignature != nil {
		w.Bool(true)
		w.FixedByteArray(m.MessageSignature)
	} else {
		w.Bool(false)
	}

	w.String(m.Message)
	w.Int64(m.Timestamp)
	w.Int64(m.Salt)

	w.VarInt(int32(len(m.PreviousMessages)))
	for _, p := range m.PreviousMessages {
		w.VarInt(p.MessageID + 1)
		if p.MessageID+1 == 0 {
			w.FixedByteArray(p.Signature)
		}
	}

	if m.UnsignedContent != nil {
		w.Bool(true)
		w.String(m.UnsignedContent.String())
	} else {
		w.Bool(false)
	}

	w.VarInt(m.FilterType)
	if m.FilterType == 2 {
		w.VarInt(int32(len(m.FilterTypeBits)))
		for _, b := range m.FilterTypeBits {
			w.Int64(b)
		}
	}

	w.VarInt(m.ChatType)
	w.String(m.NetworkName.String())
	if m.NetworkTargetName != nil {
		w.Bool(true)
		w.String(m.NetworkTargetName.String())
	} else {
		w.Bool(false)
	}
	return nil
}
