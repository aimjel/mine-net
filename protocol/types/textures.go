package types

type TexturesProperty struct {
  Timestamp int64 `json:"timestamp"`
  ProfileID string `json:"profileId"`
  ProfileName string `json:"profileName"`
  SignatureRequired bool `json:"signatureRequired"`
  Textures Textures `json:"textures"`
}

type Textures struct {
  Skin Texture `json:"SKIN"`
  Cape Texture `json:"CAPE"`
}

type Texture struct {
  URL string `json:"url"`
  Metadata TextureMetadata `json:"metadata"`
}

type TextureMetadata struct {
  Model string `json:"model"`
}
