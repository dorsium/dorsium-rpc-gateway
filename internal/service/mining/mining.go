package mining

import (
	"errors"

	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

// ProofVerifier checks validity of mining proofs.
type ProofVerifier interface {
	Verify(model.Proof) bool
}

// Service defines mining-related business logic.
type Service interface {
	SubmitProof(model.Proof) error
	GetStatus() model.MiningStatus
}

type service struct {
	verifier ProofVerifier
	status   model.MiningStatus
}

// New creates a mining service.
func New(v ProofVerifier, status model.MiningStatus) Service {
	return &service{verifier: v, status: status}
}

func (s *service) SubmitProof(p model.Proof) error {
	if !s.verifier.Verify(p) {
		return ErrInvalidProof
	}
	return nil
}

func (s *service) GetStatus() model.MiningStatus {
	return s.status
}

// ErrInvalidProof is returned when a proof fails verification.
var ErrInvalidProof = errors.New("invalid proof")

// dummyVerifier is a placeholder ProofVerifier that always succeeds.
type dummyVerifier struct{}

// NewDummyVerifier returns a ProofVerifier for development.
func NewDummyVerifier() ProofVerifier { return dummyVerifier{} }

func (dummyVerifier) Verify(model.Proof) bool { return true }
