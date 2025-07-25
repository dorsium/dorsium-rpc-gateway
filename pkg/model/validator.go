package model

import "time"

// ValidatorStatus represents the state of a validator.
type ValidatorStatus struct {
	Status string `json:"status"`
}

// ValidatorProfile holds validator metadata.
type ValidatorProfile struct {
	Address    string    `json:"address"`
	Bio        string    `json:"bio"`
	JoinDate   time.Time `json:"joinDate"`
	Reputation int       `json:"reputation"`
}

// ValidatorListItem represents a validator in list responses.
type ValidatorListItem struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

// ValidatorListResponse wraps paginated validator lists.
type ValidatorListResponse struct {
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
	Items []ValidatorListItem `json:"items"`
}

// Validator aggregates all validator information.
type Validator struct {
	Address    string          `json:"address"`
	Name       string          `json:"name"`
	Bio        string          `json:"bio"`
	JoinDate   time.Time       `json:"joinDate"`
	Reputation int             `json:"reputation"`
	Status     ValidatorStatus `json:"status"`
}
