package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"planningpoker/internal/domain/users"
	api "planningpoker/internal/infra/http"
)

// Inject storage failures at the port; a database outage must not expose driver
// diagnostics to clients or turn a missing identity into a panic.
func TestUserEndpointsHandleStorageFailures(t *testing.T) {
	for _, tc := range []struct {
		name    string
		method  string
		path    string
		err     error
		status  int
		message string
	}{
		{
			name:    "registration failure",
			method:  http.MethodPost,
			path:    "/api/v1/register",
			err:     errors.New("private database diagnostics"),
			status:  http.StatusInternalServerError,
			message: "internal server error",
		},
		{
			name:    "profile read failure",
			method:  http.MethodGet,
			path:    "/api/v1/me",
			err:     errors.New("private database diagnostics"),
			status:  http.StatusInternalServerError,
			message: "internal server error",
		},
		{
			name:    "rename failure",
			method:  http.MethodPut,
			path:    "/api/v1/me",
			err:     errors.New("private database diagnostics"),
			status:  http.StatusInternalServerError,
			message: "internal server error",
		},
		{
			name:    "identity disappears after authentication",
			method:  http.MethodGet,
			path:    "/api/v1/me",
			status:  http.StatusBadRequest,
			message: "user not found",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := api.NewAPI(failingUsersService{err: tc.err}, authenticatedUser{})
			require.NoError(t, err)

			router := gin.New()
			handler.SetupRoutes(router)

			request := httptest.NewRequest(tc.method, tc.path, strings.NewReader(`{"name":"Alice"}`))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer user")

			response := httptest.NewRecorder()

			assert.NotPanics(t, func() { router.ServeHTTP(response, request) })
			assert.Equal(t, tc.status, response.Code)
			assert.JSONEq(t, `{"error":"`+tc.message+`"}`, response.Body.String())
		})
	}
}

type failingUsersService struct {
	err error
}

func (s failingUsersService) Register(users.RegisterCommand) (*users.User, error) {
	return nil, s.err
}

func (s failingUsersService) Update(users.UpdateCommand) (*users.User, error) {
	return nil, s.err
}

func (s failingUsersService) Get(string) (*users.User, error) {
	return nil, s.err
}

type authenticatedUser struct{}

func (authenticatedUser) AuthenticateByToken(string) (string, error) {
	return "user", nil
}
