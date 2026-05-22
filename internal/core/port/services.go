package port

import (
	"context"

	"gitlab.company.com/projectname/internal/core/domain"
)

// HealthService defines the contract for application health checks.
type HealthService interface {
	GetStatus(ctx context.Context) (domain.HealthStatus, error)
}
