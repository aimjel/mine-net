package minecraft

type Messages struct {
	ProtocolTooNew string
	ProtocolTooOld string
	OnlineMode     string
}

var DefaultMessages = Messages{
	OnlineMode:     "This server is in online mode.",
	ProtocolTooNew: "Your protocol is too new!",
	ProtocolTooOld: "Your protocol is too old!",
}
