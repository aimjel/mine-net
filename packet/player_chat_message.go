package packet

import "github.com/aimjel/minecraft/chat"

// tubby ahh packet!

type PlayerChatMessage struct {
  //Header
	Sender [16]byte
  Index int32
  MessageSignature []byte
  //Body
  Message string
  Timestamp int64
  Salt int64
  //Previous Messages
  PreviousMessages []PreviousMessage
  //Other
  UnsignedContent string
  FilterType int32
  FilterTypeBits []int64
  //Network Target
  ChatType int32
  NetworkName string
  NetworkTargetName string
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
	  println("sig is present")
    w.Bool(true)
    w.FixedByteArray(m.MessageSignature)
  } else {
	  println("sig is not present")
    w.Bool(false)
  }

  w.String(m.Message)
  w.Int64(m.Timestamp)
  w.Int64(m.Salt)

	println(len(m.PreviousMessages), "prev msgs")
  w.VarInt(int32(len(m.PreviousMessages)))
  for _, p := range m.PreviousMessages {
    w.VarInt(p.MessageID + 1)
    if p.MessageID + 1 == 0 {
      w.FixedByteArray(p.Signature)
    }
  }

  if m.UnsignedContent != "" {
	  println("unsign content is present")
    w.Bool(true)
    msg := chat.NewMessage(m.UnsignedContent)
    w.String(msg.String())
  } else {
    w.Bool(false)
	  println("unsign content isnt present")
  }

	    w.VarInt(m.FilterType)
    if m.FilterType == 2 {
      w.VarInt(int32(len(m.FilterTypeBits)))
      for _, b := range m.FilterTypeBits {
        w.Int64(b)
      }
    }

  w.VarInt(m.ChatType)
  n := chat.NewMessage(m.NetworkName)
  w.String(n.String())
  if m.NetworkTargetName != "" {
    w.Bool(true)
    t := chat.NewMessage(m.NetworkTargetName)
    w.String(t.String())
  } else {
    w.Bool(false)
  }
	return nil
}
