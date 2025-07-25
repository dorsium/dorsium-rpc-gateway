package model

// WalletInfo represents wallet details.
type WalletInfo struct {
	Address    string `json:"address"`
	Balance    string `json:"balance"`
	Staking    string `json:"staking"`
	Reputation int    `json:"reputation"`
}

// Transaction represents a blockchain transaction.
type Transaction struct {
	Hash   string `json:"hash"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}

// NFT represents a non-fungible token.
type NFT struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
