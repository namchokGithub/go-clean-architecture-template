package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.company.com/projectname/internal/core/port"
	internalhealth "gitlab.company.com/projectname/internal/handler/health"
)

// Handler is the aggregate that holds all feature handlers.
type Handler struct {
	Health *internalhealth.Handler
}

// Dependencies holds port interfaces required by handlers.
type Dependencies struct {
	HealthService port.HealthService
}

// New creates the Handler aggregate.
func New(d Dependencies) *Handler {
	return &Handler{
		Health: internalhealth.New(internalhealth.Dependencies{
			Service: d.HealthService,
		}),
	}
}

// RegisterRoutes wires all handler routes onto the Echo instance.
func (h *Handler) RegisterRoutes(e *echo.Echo, authGroup *echo.Group) {
	h.Health.RegisterRoutes(e)
}
