package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Repository forwards RPC requests to a Dorsium node.
type Repository interface {
	ForwardGet(ctx context.Context, path string, query string) ([]byte, error)
	SendTx(ctx context.Context, data []byte) ([]byte, error)
}

type repo struct {
	baseURL     string
	client      *http.Client
	maxRespSize int64
}

// New creates a proxy repository.
func New(baseURL string, maxSize int64) Repository {
	return &repo{
		baseURL:     baseURL,
		client:      &http.Client{Timeout: 10 * time.Second},
		maxRespSize: maxSize,
	}
}

func (r *repo) tryRequest(ctx context.Context, method, url string, body []byte) ([]byte, error) {
	var resp *http.Response
	var err error
	for i := 0; i < 3; i++ {
		var b io.Reader
		if body != nil {
			b = bytes.NewReader(body)
		}
		reqCtx, cancel := context.WithTimeout(ctx, r.client.Timeout)
		req, reqErr := http.NewRequestWithContext(reqCtx, method, url, b)
		if reqErr != nil {
			cancel()
			return nil, reqErr
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err = r.client.Do(req)
		cancel()
		if err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	reader := io.LimitReader(resp.Body, r.maxRespSize+1)
	out, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if int64(len(out)) > r.maxRespSize {
		return nil, fmt.Errorf("response too large")
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("rpc error: %s", resp.Status)
	}
	return out, nil
}

func (r *repo) ForwardGet(ctx context.Context, path string, query string) ([]byte, error) {
	url := r.baseURL + path
	if query != "" {
		url += "?" + query
	}
	return r.tryRequest(ctx, http.MethodGet, url, nil)
}

func (r *repo) SendTx(ctx context.Context, data []byte) ([]byte, error) {
	url := r.baseURL + "/tx/send"
	return r.tryRequest(ctx, http.MethodPost, url, data)
}
