package model

// NFTMetadata describes token metadata.
type NFTMetadata struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	ImageURL   string            `json:"imageUrl"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// MintRequest represents payload for minting an NFT.
type MintRequest struct {
	Name       string            `json:"name" validate:"required"`
	ImageURL   string            `json:"imageUrl" validate:"required,url"`
	Attributes map[string]string `json:"attributes"`
}
