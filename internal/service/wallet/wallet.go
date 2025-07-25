package wallet

import "github.com/dorsium/dorsium-rpc-gateway/pkg/model"

// Service defines wallet business logic.
type Service interface {
	GetInfo(address string) (*model.WalletInfo, error)
	GetTransactions(address string, limit int) ([]model.Transaction, error)
	GetNFTs(address string) ([]model.NFT, error)
}

type service struct {
	repo Repository
}

// Repository describes required methods from repository layer.
type Repository interface {
	GetInfo(address string) (*model.WalletInfo, error)
	GetTransactions(address string, limit int) ([]model.Transaction, error)
	GetNFTs(address string) ([]model.NFT, error)
}

// New creates a wallet service.
func New(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetInfo(address string) (*model.WalletInfo, error) {
	return s.repo.GetInfo(address)
}

func (s *service) GetTransactions(address string, limit int) ([]model.Transaction, error) {
	return s.repo.GetTransactions(address, limit)
}

func (s *service) GetNFTs(address string) ([]model.NFT, error) {
	return s.repo.GetNFTs(address)
}
