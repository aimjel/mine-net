package types

// Property represents the properties object used in the Mojang API
type Property struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature"`
}
