package node

import (
	"errors"
	"sync"
	"time"

	"github.com/dorsium/dorsium-rpc-gateway/pkg/model"
)

// Repository abstracts node data storage.
type Repository interface {
	Get(id string) (*model.Node, error)
	Update(n *model.Node) error
	List(page, limit int) ([]model.Node, error)
}

type repo struct {
	mu    sync.RWMutex
	store map[string]model.Node
}

// ErrNotFound is returned when a node does not exist.
var ErrNotFound = errors.New("node not found")

// New returns an in-memory node repository with mock data.
func New() Repository {
	r := &repo{store: make(map[string]model.Node)}
	now := time.Now()
	r.store["node1"] = model.Node{
		ID:       "node1",
		Label:    "Gateway Node 1",
		Identity: "node-one",
		Location: "Earth",
		Status: model.NodeStatus{
			Health:    "healthy",
			LastPing:  now,
			SyncState: "synced",
		},
		Metrics: model.NodeMetrics{
			Uptime:       99.9,
			RequestCount: 1000,
			AvgResponse:  0.2,
		},
	}
	r.store["node2"] = model.Node{
		ID:       "node2",
		Label:    "Gateway Node 2",
		Identity: "node-two",
		Location: "Mars",
		Status: model.NodeStatus{
			Health:    "healthy",
			LastPing:  now,
			SyncState: "syncing",
		},
		Metrics: model.NodeMetrics{
			Uptime:       95.1,
			RequestCount: 800,
			AvgResponse:  0.25,
		},
	}
	return r
}

func (r *repo) Get(id string) (*model.Node, error) {
	r.mu.RLock()
	n, ok := r.store[id]
	r.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	return &n, nil
}

func (r *repo) Update(n *model.Node) error {
	r.mu.Lock()
	r.store[n.ID] = *n
	r.mu.Unlock()
	return nil
}

func (r *repo) List(page, limit int) ([]model.Node, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = len(r.store)
	}

	start := (page - 1) * limit

	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]model.Node, 0, limit)
	i := 0
	for _, n := range r.store {
		if i >= start+limit {
			break
		}
		if i >= start {
			res = append(res, n)
		}
		i++
	}
	return res, nil
}
