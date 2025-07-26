package tests

import (
	"fmt"
	"testing"

	noderepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/node"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

func TestNodeRepositoryPagination(t *testing.T) {
	repo := noderepo.New()
	// add extra nodes
	for i := 3; i <= 12; i++ {
		repo.Update(&model.Node{ID: fmt.Sprintf("node%d", i), Label: fmt.Sprintf("Node %d", i)})
	}

	nodes, err := repo.List(2, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nodes) != 5 {
		t.Fatalf("expected 5 nodes, got %d", len(nodes))
	}
}
