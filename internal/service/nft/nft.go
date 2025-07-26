package nft

import (
	"context"
	"io"
	"net/http"
	"time"

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
	GetMetadata(ctx context.Context, id string) (*model.NFTMetadata, error)
	MintNFT(ctx context.Context, req model.MintRequest) (*model.NFTMetadata, error)
	GetImage(ctx context.Context, id string) ([]byte, string, error)
}

type service struct {
	repo        Repository
	mintHandler MintHandler
}

// New creates an NFT service.
func New(repo Repository, mh MintHandler) Service {
	return &service{repo: repo, mintHandler: mh}
}

func (s *service) GetMetadata(ctx context.Context, id string) (*model.NFTMetadata, error) {
	return s.repo.GetByID(id)
}

func (s *service) MintNFT(ctx context.Context, req model.MintRequest) (*model.NFTMetadata, error) {
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

func (s *service) GetImage(ctx context.Context, id string) ([]byte, string, error) {
	nft, err := s.repo.GetByID(id)
	if err != nil {
		return nil, "", err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, nft.ImageURL, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := client.Do(req)
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
