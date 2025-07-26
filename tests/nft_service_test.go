package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	nftsvc "github.com/dorsium/dorsium-rpc-gateway/internal/service/nft"
	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

type fakeRepo struct {
	item    *model.NFTMetadata
	saveErr error
	getErr  error
	saved   *model.NFTMetadata
}

func (f *fakeRepo) Save(n model.NFTMetadata) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.saved = &n
	f.item = &n
	return nil
}

func (f *fakeRepo) GetByID(id string) (*model.NFTMetadata, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.item != nil && f.item.ID == id {
		return f.item, nil
	}
	return nil, errors.New("not found")
}

type fakeMintHandler struct {
	err  error
	last model.NFTMetadata
}

func (f *fakeMintHandler) Mint(n model.NFTMetadata) error {
	f.last = n
	return f.err
}

func TestNFTMintSuccess(t *testing.T) {
	repo := &fakeRepo{}
	mh := &fakeMintHandler{}
	svc := nftsvc.New(repo, mh)

	req := model.MintRequest{Name: "Test", ImageURL: "http://example.com/img.png"}
	nft, err := svc.MintNFT(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nft == nil || nft.ID == "" {
		t.Fatalf("expected minted NFT with ID")
	}
	if repo.saved == nil || repo.saved.ID != nft.ID {
		t.Fatalf("expected NFT saved with ID %s", nft.ID)
	}
	if mh.last.ID != nft.ID {
		t.Fatalf("mint handler not invoked")
	}
}

func TestNFTMintMintError(t *testing.T) {
	repo := &fakeRepo{}
	mh := &fakeMintHandler{err: errors.New("mint fail")}
	svc := nftsvc.New(repo, mh)
	req := model.MintRequest{Name: "Test", ImageURL: "http://example.com/img.png"}
	if _, err := svc.MintNFT(req); err == nil {
		t.Fatalf("expected error from mint handler")
	}
}

func TestNFTMintSaveError(t *testing.T) {
	repo := &fakeRepo{saveErr: errors.New("save fail")}
	mh := &fakeMintHandler{}
	svc := nftsvc.New(repo, mh)
	req := model.MintRequest{Name: "Test", ImageURL: "http://example.com/img.png"}
	if _, err := svc.MintNFT(req); err == nil {
		t.Fatalf("expected save error")
	}
}

func TestNFTGetImageSuccess(t *testing.T) {
	imgData := []byte("img")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(imgData)
	}))
	defer srv.Close()

	repo := &fakeRepo{item: &model.NFTMetadata{ID: "1", ImageURL: srv.URL}}
	svc := nftsvc.New(repo, &fakeMintHandler{})

	data, ct, err := svc.GetImage("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != string(imgData) {
		t.Fatalf("expected %s got %s", imgData, data)
	}
	if ct != "image/png" {
		t.Fatalf("expected content-type image/png got %s", ct)
	}
}

func TestNFTGetImageRepoError(t *testing.T) {
	repo := &fakeRepo{getErr: errors.New("missing")}
	svc := nftsvc.New(repo, &fakeMintHandler{})
	if _, _, err := svc.GetImage("x"); err == nil {
		t.Fatalf("expected repo error")
	}
}

func TestNFTGetImageRequestError(t *testing.T) {
	repo := &fakeRepo{item: &model.NFTMetadata{ID: "1", ImageURL: "http://%41"}}
	svc := nftsvc.New(repo, &fakeMintHandler{})
	if _, _, err := svc.GetImage("1"); err == nil {
		t.Fatalf("expected request error")
	}
}
