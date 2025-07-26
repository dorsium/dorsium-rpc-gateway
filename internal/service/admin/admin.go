package admin

import (
	"sync"

	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

// NodeRepository defines listing capabilities for nodes.
type NodeRepository interface {
	List() ([]model.Node, error)
}

// ValidatorRepository defines listing capabilities for validators.
type ValidatorRepository interface {
	List() ([]model.Validator, error)
}

// Service exposes admin operations.
type Service interface {
	Broadcast(msg string) error
	GetLogs() []string
}

type service struct {
	nodes      NodeRepository
	validators ValidatorRepository
	mu         sync.Mutex
	logs       []string
}

// New creates an admin service.
func New(nRepo NodeRepository, vRepo ValidatorRepository) Service {
	return &service{nodes: nRepo, validators: vRepo, logs: make([]string, 0)}
}

func (s *service) Broadcast(msg string) error {
	nodes, err := s.nodes.List()
	if err != nil {
		return err
	}
	validators, err := s.validators.List()
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, n := range nodes {
		s.logs = append(s.logs, "node "+n.ID+": "+msg)
	}
	for _, v := range validators {
		s.logs = append(s.logs, "validator "+v.Address+": "+msg)
	}
	return nil
}

func (s *service) GetLogs() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.logs) <= 100 {
		res := make([]string, len(s.logs))
		copy(res, s.logs)
		return res
	}
	res := make([]string, 100)
	copy(res, s.logs[len(s.logs)-100:])
	return res
}
