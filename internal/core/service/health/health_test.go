package health_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.company.com/projectname/internal/core/service/health"
)

func TestGetStatus_ReturnsOKWithVersion(t *testing.T) {
	svc := health.New(health.Dependencies{Version: "1.2.3"})

	status, err := svc.GetStatus(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "ok", status.Status)
	assert.Equal(t, "1.2.3", status.Version)
	assert.False(t, status.Timestamp.IsZero())
}
