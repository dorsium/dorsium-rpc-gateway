package proxy

import "context"

// Repository abstracts the proxy repository layer.
type Repository interface {
	ForwardGet(ctx context.Context, path string, query string) ([]byte, error)
	SendTx(ctx context.Context, data []byte) ([]byte, error)
}

// Service defines proxy operations.
type Service interface {
	ProxyGet(ctx context.Context, path string, query string) ([]byte, error)
	SendTx(ctx context.Context, data []byte) ([]byte, error)
}

type service struct {
	repo Repository
}

// New creates a proxy service.
func New(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ProxyGet(ctx context.Context, path string, query string) ([]byte, error) {
	return s.repo.ForwardGet(ctx, path, query)
}

func (s *service) SendTx(ctx context.Context, data []byte) ([]byte, error) {
	return s.repo.SendTx(ctx, data)
}
