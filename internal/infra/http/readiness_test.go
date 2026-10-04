package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	api "planningpoker/internal/infra/http"
)

func TestReadinessReportsDatabaseFailure(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{
			name:   "reachable",
			status: http.StatusOK,
		},
		{
			name:   "unavailable",
			err:    errors.New("private connection details"),
			status: http.StatusServiceUnavailable,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := &api.API{}
			h.SetReadinessCheck(func(ctx context.Context) error {
				_, bounded := ctx.Deadline()
				assert.True(t, bounded)

				return tc.err
			})

			r := gin.New()
			r.GET("/alive", h.Alive)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/alive", nil))
			assert.Equal(t, tc.status, w.Code)
			assert.NotContains(t, w.Body.String(), "private connection details")
		})
	}
}
