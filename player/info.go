package player

type Info struct {
	UUID [16]byte

	Name string

	Properties []struct {
		Name      string
		Value     string
		Signature string
	}

	GameMode       int32
	Ping           int32
	HasDisplayName bool
	DisplayName    string

	ExpiresAt    int64
	PublicKey    []byte
	KeySignature []byte

	//Listed tells the PlayerInfoUpdate packet
	//that the player should be on the player-list
	Listed bool
}
