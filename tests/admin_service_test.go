package tests

import (
	"fmt"
	"testing"

	adminsvc "github.com/dorsium/dorsium-rpc-gateway/internal/service/admin"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

type fakeNodeRepo struct{}

func (f *fakeNodeRepo) List(page, limit int) ([]model.Node, error) {
	return []model.Node{{ID: "n1"}}, nil
}

type fakeAdminValidatorRepo struct{}

func (f *fakeAdminValidatorRepo) List(page, limit int) ([]model.Validator, error) {
	return []model.Validator{{Address: "v1"}}, nil
}

func TestAdminBroadcastPrunesLogs(t *testing.T) {
	nRepo := &fakeNodeRepo{}
	vRepo := &fakeAdminValidatorRepo{}
	svc := adminsvc.New(nRepo, vRepo)

	iterations := adminsvc.MaxLogLength/2 + 10
	for i := 0; i < iterations; i++ {
		if err := svc.Broadcast(fmt.Sprintf("msg%d", i)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	logs := svc.GetLogs()
	if len(logs) != adminsvc.MaxLogLength {
		t.Fatalf("expected %d logs, got %d", adminsvc.MaxLogLength, len(logs))
	}
	firstExpected := fmt.Sprintf("node n1: msg%d", iterations-adminsvc.MaxLogLength/2)
	if logs[0] != firstExpected {
		t.Fatalf("expected first log %q got %q", firstExpected, logs[0])
	}
	lastExpected := fmt.Sprintf("validator v1: msg%d", iterations-1)
	if logs[len(logs)-1] != lastExpected {
		t.Fatalf("expected last log %q got %q", lastExpected, logs[len(logs)-1])
	}
}
