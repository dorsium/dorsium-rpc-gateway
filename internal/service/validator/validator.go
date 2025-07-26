package validator

import (
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

// MaxListLimit caps the maximum number of items returned by List.
const MaxListLimit = 100

// Repository describes persistence layer requirements.
type Repository interface {
	Get(address string) (*model.Validator, error)
	List(page, limit int) ([]model.Validator, error)
}

// Service exposes validator operations.
type Service interface {
	GetStatus(address string) (*model.ValidatorStatus, error)
	GetProfile(address string) (*model.ValidatorProfile, error)
	List(page, limit int) (*model.ValidatorListResponse, error)
}

type service struct {
	repo Repository
}

// New creates a validator service.
func New(repo Repository) Service { return &service{repo: repo} }

func (s *service) GetStatus(address string) (*model.ValidatorStatus, error) {
	v, err := s.repo.Get(address)
	if err != nil {
		return nil, err
	}
	return &v.Status, nil
}

func (s *service) GetProfile(address string) (*model.ValidatorProfile, error) {
	v, err := s.repo.Get(address)
	if err != nil {
		return nil, err
	}
	return &model.ValidatorProfile{
		Address:    v.Address,
		Bio:        v.Bio,
		JoinDate:   v.JoinDate,
		Reputation: v.Reputation,
	}, nil
}

func (s *service) List(page, limit int) (*model.ValidatorListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	} else if limit > MaxListLimit {
		limit = MaxListLimit
	}

	vals, err := s.repo.List(page, limit)
	if err != nil {
		return nil, err
	}

	items := make([]model.ValidatorListItem, 0, len(vals))
	for _, v := range vals {
		items = append(items, model.ValidatorListItem{Address: v.Address, Name: v.Name})
	}
	return &model.ValidatorListResponse{Page: page, Limit: limit, Items: items}, nil
}
