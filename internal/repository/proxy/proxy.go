package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Repository forwards RPC requests to a Dorsium node.
type Repository interface {
	ForwardGet(path string, query string) ([]byte, error)
	SendTx(data []byte) ([]byte, error)
}

type repo struct {
	baseURL string
	client  *http.Client
}

// New creates a proxy repository.
func New(baseURL string) Repository {
	return &repo{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *repo) tryRequest(method, url string, body []byte) ([]byte, error) {
	var resp *http.Response
	var err error
	for i := 0; i < 3; i++ {
		var b io.Reader
		if body != nil {
			b = bytes.NewReader(body)
		}
		req, reqErr := http.NewRequest(method, url, b)
		if reqErr != nil {
			return nil, reqErr
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err = r.client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("rpc error: %s", resp.Status)
	}
	return out, nil
}

func (r *repo) ForwardGet(path string, query string) ([]byte, error) {
	url := r.baseURL + path
	if query != "" {
		url += "?" + query
	}
	return r.tryRequest(http.MethodGet, url, nil)
}

func (r *repo) SendTx(data []byte) ([]byte, error) {
	url := r.baseURL + "/tx/send"
	return r.tryRequest(http.MethodPost, url, data)
}
