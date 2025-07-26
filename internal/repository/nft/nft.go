package nft

import (
	"errors"
	"sync"

	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

// Repository abstracts NFT data source.
type Repository interface {
	Save(model.NFTMetadata) error
	GetByID(string) (*model.NFTMetadata, error)
}

type repo struct {
	mu    sync.RWMutex
	store map[string]model.NFTMetadata
}

// ErrNotFound is returned when an NFT is missing.
var ErrNotFound = errors.New("nft not found")

// New returns an in-memory NFT repository implementation.
func New() Repository {
	return &repo{store: make(map[string]model.NFTMetadata)}
}

func (r *repo) Save(n model.NFTMetadata) error {
	r.mu.Lock()
	r.store[n.ID] = n
	r.mu.Unlock()
	return nil
}

func (r *repo) GetByID(id string) (*model.NFTMetadata, error) {
	r.mu.RLock()
	n, ok := r.store[id]
	r.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	return &n, nil
}
