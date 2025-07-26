package tests

import (
	"testing"

	validatorrepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/validator"
)

func TestValidatorRepositoryPagination(t *testing.T) {
	repo := validatorrepo.New()

	vals, err := repo.List(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vals) != 1 {
		t.Fatalf("expected 1 validator, got %d", len(vals))
	}

	vals, err = repo.List(3, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vals) != 0 {
		t.Fatalf("expected 0 validators, got %d", len(vals))
	}
}
