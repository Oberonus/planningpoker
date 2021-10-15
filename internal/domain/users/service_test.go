package users_test

import (
	"errors"
	"planningpoker/internal/domain/users"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		usersRepo users.Repository
		expError  string
	}{
		"success": {
			usersRepo: usersRepoStub{},
		},
		"failed on no repository": {
			expError: "users repository should be provided",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			srv, err := users.NewService(tt.usersRepo)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
				assert.Nil(t, srv)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, srv)
			}
		})
	}
}

func TestService_Register(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		usersRepo users.Repository
		name      string
		expError  string
	}{
		"success": {
			usersRepo: usersRepoStub{},
			name:      "foo",
		},
		"failed on wrong name": {
			usersRepo: usersRepoStub{},
			name:      "",
			expError:  "user creation: user name should be provided",
		},
		"failed repository": {
			usersRepo: usersRepoStub{saveErr: errors.New("save failed")},
			name:      "foo",
			expError:  "save failed",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			srv, err := users.NewService(tt.usersRepo)
			require.NoError(t, err)

			cmd, err := users.NewRegisterCommand(tt.name)
			require.NoError(t, err)

			u, err := srv.Register(*cmd)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
				assert.Nil(t, u)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, u)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	t.Parallel()
	uid := "1"

	testCases := map[string]struct {
		usersRepo users.Repository
		name      string
		expError  string
	}{
		"success": {
			usersRepo: usersRepoStub{
				getUser: users.NewRaw(uid, "foo"),
			},
			name: "foo",
		},
		"failed on wrong name": {
			usersRepo: usersRepoStub{
				getUser: users.NewRaw(uid, "foo"),
			},
			name:     "",
			expError: "user name should be provided",
		},
		"failed fetch from repository": {
			usersRepo: usersRepoStub{
				getErr: errors.New("get failed"),
			},
			name:     "foo",
			expError: "get failed",
		},
		"failed on user not found": {
			usersRepo: usersRepoStub{},
			name:      "foo",
			expError:  "user not found",
		},
		"failed save repository": {
			usersRepo: usersRepoStub{
				getUser: users.NewRaw(uid, "foo"),
				saveErr: errors.New("save failed"),
			},
			name:     "foo",
			expError: "save failed",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			srv, err := users.NewService(tt.usersRepo)
			require.NoError(t, err)

			cmd, err := users.NewUpdateCommand(uid, tt.name)
			require.NoError(t, err)

			u, err := srv.Update(*cmd)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
				assert.Nil(t, u)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, u)
			}
		})
	}
}

func TestService_AuthenticateByID(t *testing.T) {
	t.Parallel()
	uid := "1"

	testCases := map[string]struct {
		usersRepo users.Repository
		expError  string
	}{
		"success": {
			usersRepo: usersRepoStub{
				getUser: users.NewRaw(uid, "foo"),
			},
		},
		"failed fetch from repository": {
			usersRepo: usersRepoStub{
				getErr: errors.New("get failed"),
			},
			expError: "get failed",
		},
		"failed on user not found": {
			usersRepo: usersRepoStub{},
			expError:  "user not found",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			srv, err := users.NewService(tt.usersRepo)
			require.NoError(t, err)

			cmd, err := users.NewAuthByIDCommand(uid)
			require.NoError(t, err)

			u, err := srv.AuthenticateByID(*cmd)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
				assert.Nil(t, u)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, u)
			}
		})
	}
}

type usersRepoStub struct {
	getErr  error
	getUser *users.User
	saveErr error
}

func (u usersRepoStub) Get(id string) (*users.User, error) {
	if u.getUser != nil && u.getUser.ID() != id {
		return nil, errors.New("test failed - user ID does not match")
	}
	return u.getUser, u.getErr
}

func (u usersRepoStub) GetMany([]string) ([]users.User, error) {
	return nil, errors.New("not implemented")
}

func (u usersRepoStub) Save(users.User) error {
	return u.saveErr
}
