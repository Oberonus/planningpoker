package repository

import (
	"encoding/json"
	"planningpoker/internal/domain"
	"sync"
)

type userDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type MemoryUserRepository struct {
	m     sync.RWMutex
	users map[string][]byte
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[string][]byte),
	}
}

func (r *MemoryUserRepository) Get(id string) (*domain.User, error) {
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

	user := domain.User{
		ID:    dto.ID,
		Name:  dto.Name,
	}

	return &user, nil
}

func (r *MemoryUserRepository) GetMany(ids []string) ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(ids))
	for _, id := range ids {
		u, err := r.Get(id)
		if err != nil {
			return nil, err
		}
		if u == nil {
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *MemoryUserRepository) Save(user *domain.User) error {
	dto := userDTO{
		ID:    user.ID,
		Name:  user.Name,
	}

	raw, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	r.m.Lock()
	defer r.m.Unlock()
	r.users[user.ID] = raw

	return nil
}
