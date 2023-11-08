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

	//Listed tells the PlayerInfoUpdate packet
	//that the player should be on the player-list
	Listed bool
}
