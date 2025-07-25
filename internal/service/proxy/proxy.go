package proxy

// Repository abstracts the proxy repository layer.
type Repository interface {
	ForwardGet(path string, query string) ([]byte, error)
	SendTx(data []byte) ([]byte, error)
}

// Service defines proxy operations.
type Service interface {
	ProxyGet(path string, query string) ([]byte, error)
	SendTx(data []byte) ([]byte, error)
}

type service struct {
	repo Repository
}

// New creates a proxy service.
func New(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ProxyGet(path string, query string) ([]byte, error) {
	return s.repo.ForwardGet(path, query)
}

func (s *service) SendTx(data []byte) ([]byte, error) {
	return s.repo.SendTx(data)
}
