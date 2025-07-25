package repository

// Repository abstracts access to external resources.
type Repository interface {
	// TODO: add methods
}

// Dummy repository implementation for wiring.
type repo struct{}

// New returns a Repository implementation.
func New() Repository {
	return &repo{}
}
