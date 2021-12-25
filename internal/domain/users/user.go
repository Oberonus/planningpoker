package users

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type User struct {
	id   string
	name string
}

func NewUser(name string) (*User, error) {
	u := &User{
		id: strings.ReplaceAll(uuid.New().String(), "-", ""),
	}
	if err := u.NameAs(name); err != nil {
		return nil, err
	}

	return u, nil
}

func NewRaw(id string, name string) *User {
	return &User{
		id:   id,
		name: name,
	}
}

func (u User) ID() string {
	return u.id
}

func (u User) Name() string {
	return u.name
}

func (u *User) NameAs(name string) error {
	if name == "" {
		return errors.New("user name should be provided")
	}
	u.name = name

	return nil
}
