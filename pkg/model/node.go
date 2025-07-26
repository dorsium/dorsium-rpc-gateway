package model

import "time"

// NodeStatus represents health information for a node.
type NodeStatus struct {
	Health    string    `json:"health"`
	LastPing  time.Time `json:"lastPing"`
	SyncState string    `json:"syncState"`
}

// NodeProfile contains identity and location metadata.
type NodeProfile struct {
	ID       string `json:"id"`
	Identity string `json:"identity"`
	Location string `json:"location"`
}

// NodeMetrics holds usage statistics for a node.
type NodeMetrics struct {
	Uptime       float64 `json:"uptime"`
	RequestCount int     `json:"requestCount"`
	AvgResponse  float64 `json:"avgResponse"`
}

// NodeListItem represents a node in list responses.
type NodeListItem struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// NodeListResponse wraps paginated list responses.
type NodeListResponse struct {
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Items []NodeListItem `json:"items"`
}

// NodePing is received when nodes send heartbeat information.
type NodePing struct {
	ID        string `json:"id" validate:"required"`
	Health    string `json:"health" validate:"required"`
	SyncState string `json:"syncState" validate:"required"`
}

// Node aggregates all node information.
type Node struct {
	ID       string      `json:"id"`
	Label    string      `json:"label"`
	Identity string      `json:"identity"`
	Location string      `json:"location"`
	Status   NodeStatus  `json:"status"`
	Metrics  NodeMetrics `json:"metrics"`
}
