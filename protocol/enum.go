package protocol

type PlayerListUpdateAction byte

const (
	AddPlayer PlayerListUpdateAction = 1 << iota
	InitializeChat
	UpdateGameMode
	UpdateListed
	UpdateLatency
	UpdateDisplayName
)
