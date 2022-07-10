// Package users contains domain level users logic.
package users

import (
	"errors"
	"strings"

	"planningpoker/internal/domain"
	"planningpoker/internal/domain/events"

	"github.com/google/uuid"
)

// User is a user aggregate.
type User struct {
	domain.BaseAggregate
	id   string
	name string
}

// NewUser creates a new user.
func NewUser(name string) (*User, error) {
	u := &User{
		id: strings.ReplaceAll(uuid.New().String(), "-", ""),
	}
	if err := u.NameAs(name); err != nil {
		return nil, err
	}

	return u, nil
}

// NewRaw instantiates a user aggregate from raw data.
// It should never be used in any logic except aggregate hydration from any serialized format (db, etc...)
func NewRaw(id string, name string) *User {
	return &User{
		id:   id,
		name: name,
	}
}

// ID returns the unique user identifier.
func (u User) ID() string {
	return u.id
}

// Name returns the user name.
func (u User) Name() string {
	return u.name
}

// NameAs changes the user name.
func (u *User) NameAs(name string) error {
	if name == "" {
		return errors.New("user name should be provided")
	}
	u.name = name
	u.AddEvent(events.NewDomainEventBuilder(events.EventTypeUserUpdated).ForAggregate(u.ID()).Build())

	return nil
}
