package service

// Service defines business logic interfaces.
type Service interface {
	Ping() string
}

type service struct{}

// New returns a Service implementation.
func New() Service {
	return &service{}
}

func (s *service) Ping() string {
	return "pong"
}
