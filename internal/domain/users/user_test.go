package users_test

import (
	"planningpoker/internal/domain/users"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		name     string
		expError string
	}{
		"success": {
			name: "foo",
		},
		"failed on wrong name": {
			expError: "user name should be provided",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			u, err := users.NewUser(tt.name)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
				assert.Nil(t, u)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, u)
				assert.NotEmpty(t, u.ID())
				assert.Equal(t, tt.name, u.Name())
			}
		})
	}
}

func TestNewRawUser(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		id       string
		name     string
		expError string
	}{
		"success": {
			id:   "bar",
			name: "foo",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			u := users.NewRaw(tt.id, tt.name)
			require.NotNil(t, u)
			assert.NotEmpty(t, u.ID())
			assert.Equal(t, tt.name, u.Name())
			assert.Equal(t, tt.id, u.ID())
		})
	}
}
