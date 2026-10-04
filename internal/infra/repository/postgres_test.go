package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"planningpoker/internal/infra/repository"
)

func TestOpenPostgresRejectsInvalidConfigurationWithoutExposingCredentials(t *testing.T) {
	const password = "private-review-password"

	db, err := repository.OpenPostgres(context.Background(),
		"postgres://app:"+password+"@localhost:invalid/app")

	assert.Nil(t, db)
	assert.EqualError(t, err, "invalid DATABASE_URL")
}
