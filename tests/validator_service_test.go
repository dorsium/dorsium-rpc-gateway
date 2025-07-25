package tests

import (
	"fmt"
	"testing"

	validatorsvc "github.com/dorsium/dorsium-rpc-gateway/internal/service/validator"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

// fakeValidatorRepo implements validatorsvc.Repository for testing.
type fakeValidatorRepo struct {
	items []model.Validator
}

func (f *fakeValidatorRepo) Get(address string) (*model.Validator, error) {
	for _, v := range f.items {
		if v.Address == address {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (f *fakeValidatorRepo) List(page, limit int) ([]model.Validator, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = len(f.items)
	}
	start := (page - 1) * limit
	if start >= len(f.items) {
		return []model.Validator{}, nil
	}
	end := start + limit
	if end > len(f.items) {
		end = len(f.items)
	}
	return f.items[start:end], nil
}

func TestValidatorListLimitCap(t *testing.T) {
	repo := &fakeValidatorRepo{}
	for i := 0; i < validatorsvc.MaxListLimit+50; i++ {
		repo.items = append(repo.items, model.Validator{Address: fmt.Sprintf("addr%d", i), Name: fmt.Sprintf("Val %d", i)})
	}
	svc := validatorsvc.New(repo)
	res, err := svc.List(1, validatorsvc.MaxListLimit+50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Items) != validatorsvc.MaxListLimit {
		t.Fatalf("expected %d items, got %d", validatorsvc.MaxListLimit, len(res.Items))
	}
	if res.Limit != validatorsvc.MaxListLimit {
		t.Fatalf("expected limit %d, got %d", validatorsvc.MaxListLimit, res.Limit)
	}
}
