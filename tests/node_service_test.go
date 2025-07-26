package tests

import (
	"fmt"
	"testing"

	noderepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/node"
	nodesvc "github.com/dorsium/dorsium-rpc-gateway/internal/service/node"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

func TestNodeListLimitCap(t *testing.T) {
	repo := noderepo.New()
	for i := 3; i <= nodesvc.MaxListLimit+50; i++ {
		repo.Update(&model.Node{ID: fmt.Sprintf("node%d", i), Label: fmt.Sprintf("Node %d", i)})
	}
	svc := nodesvc.New(repo)
	res, err := svc.List(1, nodesvc.MaxListLimit+50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Items) != nodesvc.MaxListLimit {
		t.Fatalf("expected %d items, got %d", nodesvc.MaxListLimit, len(res.Items))
	}
	if res.Limit != nodesvc.MaxListLimit {
		t.Fatalf("expected limit %d, got %d", nodesvc.MaxListLimit, res.Limit)
	}
}
