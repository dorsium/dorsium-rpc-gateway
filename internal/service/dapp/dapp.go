package dapp

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/dorsium/dorsium-rpc-gateway/internal/repository/nft"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

// Service exposes DAPP operations.
type Service interface {
	GetConfig() model.DappConfig
	VerifyNFT(id string) error
	VerifyWallet(req model.WalletVerifyRequest) (bool, error)
	GetPermissions(address string) (*model.DappPermissions, error)
}

type service struct {
	nftRepo nft.Repository
}

// New creates a DAPP service.
func New(nRepo nft.Repository) Service {
	return &service{nftRepo: nRepo}
}

func (s *service) GetConfig() model.DappConfig {
	return model.DappConfig{
		RPCURL: "https://rpc.example.com",
		Token:  "exampletoken",
	}
}

func (s *service) VerifyNFT(id string) error {
	if _, err := s.nftRepo.GetByID(id); err != nil {
		return err
	}
	return nil
}

func verifySignature(wallet, payload, signature string) bool {
	sum := sha256.Sum256([]byte(payload + wallet))
	return strings.EqualFold(hex.EncodeToString(sum[:]), signature)
}

func (s *service) VerifyWallet(req model.WalletVerifyRequest) (bool, error) {
	if !verifySignature(req.Wallet, req.Payload, req.Signature) {
		return false, errors.New("invalid signature")
	}
	return true, nil
}

func (s *service) GetPermissions(address string) (*model.DappPermissions, error) {
	perms := model.DappPermissions{
		Role:  "user",
		NFTs:  []string{},
		Flags: []string{},
	}
	return &perms, nil
}
