package types

import "github.com/aimjel/minecraft/chat"

type PlayerInfo struct {
	UUID [16]byte

	Name string

	Properties []Property

	GameMode       int32
	Ping           int32
	DisplayName *chat.Message

	ChatSessionID [16]byte
	ExpiresAt     int64
	PublicKey     []byte
	KeySignature  []byte

	// Listed lists the player's info in the server's player list
	Listed bool
}
