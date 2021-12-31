package repository

import (
	"encoding/json"
	"planningpoker/internal/domain/users"
	"sync"
)

type userDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// MemoryUserRepository is a simple in-memory linear users repository.
type MemoryUserRepository struct {
	m     sync.RWMutex
	users map[string][]byte
}

// NewMemoryUserRepository creates an in-memory users repository instance.
func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[string][]byte),
	}
}

// Get retrieves the user by ID.
func (r *MemoryUserRepository) Get(id string) (*users.User, error) {
	r.m.RLock()
	raw, ok := r.users[id]
	r.m.RUnlock()

	if !ok {
		return nil, nil
	}

	dto := userDTO{}
	err := json.Unmarshal(raw, &dto)
	if err != nil {
		return nil, err
	}

	return users.NewRaw(dto.ID, dto.Name), nil
}

// GetMany retrieves many users.
func (r *MemoryUserRepository) GetMany(ids []string) ([]users.User, error) {
	list := make([]users.User, 0, len(ids))
	for _, id := range ids {
		u, err := r.Get(id)
		if err != nil {
			return nil, err
		}
		if u == nil {
			continue
		}
		list = append(list, *u)
	}

	return list, nil
}

// Save persists the user.
func (r *MemoryUserRepository) Save(user users.User) error {
	dto := userDTO{
		ID:   user.ID(),
		Name: user.Name(),
	}

	raw, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	r.m.Lock()
	defer r.m.Unlock()
	r.users[user.ID()] = raw

	return nil
}
