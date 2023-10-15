package types

type CommandNode struct {
	Flags           uint8
	Children        []int32
	RedirectNode    int32
	Name            string
	ParserID        int32
	SuggestionsType string

	Properties CommandProperties
}

type CommandProperties struct {
	Flags      uint8
	Min, Max   uint64 //min and max stores the bits of the actual data type
	Identifier string
}
