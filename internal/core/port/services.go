package port

import (
	"context"

	"gitlab.socket9.com/sritrang/super-driver/backend/internal/core/domain"
)

// HealthService defines the contract for application health checks.
type HealthService interface {
	GetStatus(ctx context.Context) (domain.HealthStatus, error)
}
