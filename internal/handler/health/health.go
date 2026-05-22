package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.socket9.com/sritrang/super-driver/backend/internal/core/port"
	"gitlab.socket9.com/sritrang/super-driver/backend/internal/handler/common"
)

// Handler handles HTTP requests for health checks.
type Handler struct {
	svc port.HealthService
}

// Dependencies holds the port interfaces needed by this handler.
type Dependencies struct {
	Service port.HealthService
}

// New creates a health Handler.
func New(d Dependencies) *Handler {
	return &Handler{svc: d.Service}
}

// RegisterRoutes registers public routes on the Echo root instance.
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/health", h.GetStatus)
}

// GetStatus handles GET /health.
func (h *Handler) GetStatus(c echo.Context) error {
	status, err := h.svc.GetStatus(c.Request().Context())
	if err != nil {
		return common.NewInternalServerError()
	}
	return c.JSON(http.StatusOK, common.NewSuccessResponse(status))
}
