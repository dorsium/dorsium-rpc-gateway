package wallet

import "github.com/dorsium/dorsium-rpc-gateway/pkg/model"

// Repository abstracts wallet data source.
type Repository interface {
	GetInfo(address string) (*model.WalletInfo, error)
	GetTransactions(address string, limit int) ([]model.Transaction, error)
	GetNFTs(address string) ([]model.NFT, error)
}

type repo struct{}

// New returns a dummy wallet repository implementation.
func New() Repository { return &repo{} }

func (r *repo) GetInfo(address string) (*model.WalletInfo, error) {
	return &model.WalletInfo{
		Address:    address,
		Balance:    "0",
		Staking:    "0",
		Reputation: 0,
	}, nil
}

func (r *repo) GetTransactions(address string, limit int) ([]model.Transaction, error) {
	return []model.Transaction{}, nil
}

func (r *repo) GetNFTs(address string) ([]model.NFT, error) {
	return []model.NFT{}, nil
}
