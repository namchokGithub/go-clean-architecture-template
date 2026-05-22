package service

import (
	"gitlab.company.com/projectname/configs"
	"gitlab.company.com/projectname/internal/core/port"
	"gitlab.company.com/projectname/internal/core/service/health"
)

// Dependencies holds all service-layer dependencies.
type Dependencies struct {
	Config configs.Config
}

// Service is the aggregate that holds all feature services.
type Service struct {
	Health port.HealthService
}

// New creates the Service aggregate. Register feature services here.
func New(d Dependencies) *Service {
	return &Service{
		Health: health.New(health.Dependencies{
			Version: d.Config.App.Version,
		}),
	}
}
