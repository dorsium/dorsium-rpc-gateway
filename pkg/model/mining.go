package model

// Proof represents a mining proof submission.
type Proof struct {
	MinerID string `json:"minerId" validate:"required"`
	Nonce   int64  `json:"nonce" validate:"required"`
}

// MiningStatus provides mining configuration details.
type MiningStatus struct {
	Mode       string `json:"mode"`
	Difficulty int    `json:"difficulty"`
	Challenge  string `json:"challenge"`
}
