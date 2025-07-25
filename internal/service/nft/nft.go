package nft

import (
	"io"
	"net/http"

	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
	"github.com/google/uuid"
)

// MintHandler executes NFT minting.
type MintHandler interface {
	Mint(model.NFTMetadata) error
}

// Repository defines persistence layer methods.
type Repository interface {
	Save(model.NFTMetadata) error
	GetByID(string) (*model.NFTMetadata, error)
}

// Service exposes NFT operations.
type Service interface {
	GetMetadata(id string) (*model.NFTMetadata, error)
	MintNFT(model.MintRequest) (*model.NFTMetadata, error)
	GetImage(id string) ([]byte, string, error)
}

type service struct {
	repo        Repository
	mintHandler MintHandler
}

// New creates an NFT service.
func New(repo Repository, mh MintHandler) Service {
	return &service{repo: repo, mintHandler: mh}
}

func (s *service) GetMetadata(id string) (*model.NFTMetadata, error) {
	return s.repo.GetByID(id)
}

func (s *service) MintNFT(req model.MintRequest) (*model.NFTMetadata, error) {
	nft := model.NFTMetadata{
		ID:         uuid.New().String(),
		Name:       req.Name,
		ImageURL:   req.ImageURL,
		Attributes: req.Attributes,
	}
	if err := s.mintHandler.Mint(nft); err != nil {
		return nil, err
	}
	if err := s.repo.Save(nft); err != nil {
		return nil, err
	}
	return &nft, nil
}

func (s *service) GetImage(id string) ([]byte, string, error) {
	nft, err := s.repo.GetByID(id)
	if err != nil {
		return nil, "", err
	}
	resp, err := http.Get(nft.ImageURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = http.DetectContentType(data)
	}
	return data, ct, nil
}

// dummyMintHandler is a placeholder MintHandler.
type dummyMintHandler struct{}

// NewDummyMintHandler creates a no-op MintHandler.
func NewDummyMintHandler() MintHandler { return dummyMintHandler{} }

func (dummyMintHandler) Mint(model.NFTMetadata) error { return nil }
