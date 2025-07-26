package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	proxyrepo "github.com/dorsium/dorsium-rpc-gateway/internal/repository/proxy"
)

func TestProxyRepository_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	repo := proxyrepo.New(srv.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	if _, err := repo.ForwardGet(ctx, "/", ""); err == nil {
		t.Fatalf("expected context cancellation error")
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()
	if _, err := repo.SendTx(ctx2, []byte(`{}`)); err == nil {
		t.Fatalf("expected context cancellation error")
	}
}
