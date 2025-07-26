package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dorsium/dorsium-rpc-gateway/internal/config"
	server "github.com/dorsium/dorsium-rpc-gateway/internal/http"
	"github.com/dorsium/dorsium-rpc-gateway/internal/service"
)

func TestMetricsAggregatesUnknown(t *testing.T) {
	t.Setenv("ADMIN_TOKEN", "secret")
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("config.New failed: %v", err)
	}
	cfg.DisableMetrics = false
	srv := server.NewServer(cfg, service.New())
	srv.RegisterRoutes()

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/unknown%d", i), nil)
		if _, err := srv.App().Test(req); err != nil {
			t.Fatalf("request failed: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	if _, err := srv.App().Test(req); err != nil {
		t.Fatalf("ping failed: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	resp, err := srv.App().Test(req)
	if err != nil {
		t.Fatalf("metrics failed: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if len(srv.Calls()) != 3 {
		t.Fatalf("expected 3 endpoints recorded, got %d", len(srv.Calls()))
	}
	if !strings.Contains(string(body), `endpoint_calls_total{endpoint="unknown"} 10`) {
		t.Fatalf("unexpected metrics output: %s", string(body))
	}
}
