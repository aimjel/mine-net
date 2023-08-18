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
}
