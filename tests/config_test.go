package tests

import (
	"testing"

	"github.com/dorsium/dorsium-rpc-gateway/internal/config"
)

func TestConfigNewFailsWithoutAdminToken(t *testing.T) {
	t.Setenv("ADMIN_TOKEN", "")
	if _, err := config.New(); err == nil {
		t.Fatal("expected error when ADMIN_TOKEN is missing")
	}
}
