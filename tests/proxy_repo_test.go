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

	repo := proxyrepo.New(srv.URL, 1024)

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

func TestProxyRepository_LargeResponse(t *testing.T) {
	payload := make([]byte, 2048)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(payload)
	}))
	defer srv.Close()

	repo := proxyrepo.New(srv.URL, 1024)
	if _, err := repo.ForwardGet(context.Background(), "/", ""); err == nil {
		t.Fatalf("expected size limit error")
	}
}
