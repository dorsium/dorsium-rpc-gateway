package node

import (
	"time"

	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

// MaxListLimit caps the maximum number of items returned by List.
const MaxListLimit = 100

// Repository defines persistence layer requirements.
type Repository interface {
	Get(id string) (*model.Node, error)
	Update(n *model.Node) error
	List(page, limit int) ([]model.Node, error)
}

// Service exposes node operations.
type Service interface {
	GetStatus(id string) (*model.NodeStatus, error)
	Ping(p model.NodePing) error
	GetProfile(id string) (*model.NodeProfile, error)
	GetMetrics(id string) (*model.NodeMetrics, error)
	List(page, limit int) (*model.NodeListResponse, error)
}

type service struct {
	repo Repository
}

// New creates a node service.
func New(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetStatus(id string) (*model.NodeStatus, error) {
	n, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return &n.Status, nil
}

func (s *service) Ping(p model.NodePing) error {
	n, err := s.repo.Get(p.ID)
	if err != nil {
		return err
	}
	n.Status.Health = p.Health
	n.Status.SyncState = p.SyncState
	n.Status.LastPing = time.Now()
	return s.repo.Update(n)
}

func (s *service) GetProfile(id string) (*model.NodeProfile, error) {
	n, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return &model.NodeProfile{ID: n.ID, Identity: n.Identity, Location: n.Location}, nil
}

func (s *service) GetMetrics(id string) (*model.NodeMetrics, error) {
	n, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return &n.Metrics, nil
}

func (s *service) List(page, limit int) (*model.NodeListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	} else if limit > MaxListLimit {
		limit = MaxListLimit
	}

	nodes, err := s.repo.List(page, limit)
	if err != nil {
		return nil, err
	}

	items := make([]model.NodeListItem, 0, len(nodes))
	for _, n := range nodes {
		items = append(items, model.NodeListItem{ID: n.ID, Label: n.Label})
	}
	return &model.NodeListResponse{Page: page, Limit: limit, Items: items}, nil
}
