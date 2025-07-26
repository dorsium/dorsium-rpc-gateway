package validator

import (
	"errors"
	"sync"
	"time"

	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

// Repository abstracts validator data storage.
type Repository interface {
	Get(address string) (*model.Validator, error)
	List(page, limit int) ([]model.Validator, error)
}

type repo struct {
	mu    sync.RWMutex
	store map[string]model.Validator
}

// ErrNotFound is returned when a validator does not exist.
var ErrNotFound = errors.New("validator not found")

// New returns an in-memory validator repository with mock data.
func New() Repository {
	r := &repo{store: make(map[string]model.Validator)}
	now := time.Now()
	r.store["0xvalidator1"] = model.Validator{
		Address:    "0xvalidator1",
		Name:       "Validator One",
		Bio:        "First validator",
		JoinDate:   now.AddDate(-1, 0, 0),
		Reputation: 80,
		Status:     model.ValidatorStatus{Status: "active"},
	}
	r.store["0xvalidator2"] = model.Validator{
		Address:    "0xvalidator2",
		Name:       "Validator Two",
		Bio:        "Second validator",
		JoinDate:   now.AddDate(-2, 0, 0),
		Reputation: 95,
		Status:     model.ValidatorStatus{Status: "jailed"},
	}
	return r
}

func (r *repo) Get(address string) (*model.Validator, error) {
	r.mu.RLock()
	v, ok := r.store[address]
	r.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	return &v, nil
}

func (r *repo) List(page, limit int) ([]model.Validator, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = len(r.store)
	}
	start := (page - 1) * limit

	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]model.Validator, 0, limit)
	i := 0
	for _, v := range r.store {
		if i >= start+limit {
			break
		}
		if i >= start {
			res = append(res, v)
		}
		i++
	}
	return res, nil
}
