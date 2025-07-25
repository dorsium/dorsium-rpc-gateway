package model

// DappConfig provides configuration for DAPP frontends.
type DappConfig struct {
	RPCURL string `json:"rpcUrl"`
	Token  string `json:"token"`
}

// WalletVerifyRequest defines payload for wallet signature verification.
type WalletVerifyRequest struct {
	Wallet    string `json:"wallet" validate:"required"`
	Payload   string `json:"payload" validate:"required"`
	Signature string `json:"signature" validate:"required"`
}

// DappPermissions describes access details for a user.
type DappPermissions struct {
	Role  string   `json:"role"`
	NFTs  []string `json:"nfts"`
	Flags []string `json:"flags"`
}
