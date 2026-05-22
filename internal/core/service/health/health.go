package health

import (
	"context"
	"time"

	"gitlab.company.com/projectname/internal/core/domain"
)

// Dependencies holds construction-time config for the health service.
type Dependencies struct {
	Version string
}

// service implements port.HealthService.
type service struct {
	version string
}

// New creates a health service.
func New(d Dependencies) *service {
	return &service{version: d.Version}
}

func (s *service) GetStatus(_ context.Context) (domain.HealthStatus, error) {
	return domain.HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
		Version:   s.version,
	}, nil
}
