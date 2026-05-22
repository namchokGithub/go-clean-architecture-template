package domain

import "time"

// HealthStatus is the response model for the health check endpoint.
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}
