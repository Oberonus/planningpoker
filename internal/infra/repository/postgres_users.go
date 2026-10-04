package repository

import (
	"context"
	"database/sql"
	"errors"

	"planningpoker/internal/domain/events"
	"planningpoker/internal/domain/state"
	"planningpoker/internal/domain/users"
)

// PostgresUserRepository persists identities across application restarts.
type PostgresUserRepository struct {
	db  *sql.DB
	bus events.EventBus
}

var _ users.Repository = (*PostgresUserRepository)(nil)
var _ state.UsersRepository = (*PostgresUserRepository)(nil)

// NewPostgresUserRepository creates a repository using an initialized database.
func NewPostgresUserRepository(db *sql.DB, bus events.EventBus) *PostgresUserRepository {
	return &PostgresUserRepository{
		db:  db,
		bus: bus,
	}
}

// Get retrieves a user, returning nil when it does not exist.
func (r *PostgresUserRepository) Get(id string) (*users.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var name string
	err := r.db.QueryRowContext(ctx, "SELECT name FROM users WHERE id = $1", id).Scan(&name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return users.NewRaw(id, name), nil
}

// GetMany retrieves existing users in the supplied order, skipping missing users.
func (r *PostgresUserRepository) GetMany(ids []string) ([]users.User, error) {
	if len(ids) == 0 {
		return []users.User{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)

	defer cancel()

	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM users WHERE id = ANY($1::text[])", ids)

	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	found := make(map[string]users.User)

	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}

		found[id] = *users.NewRaw(id, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	list := make([]users.User, 0, len(ids))

	for _, id := range ids {
		if user, ok := found[id]; ok {
			list = append(list, user)
		}
	}

	return list, nil
}

// Save persists a user before publishing events to subscribers.
func (r *PostgresUserRepository) Save(user users.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `INSERT INTO users (id, name) VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name`, user.ID(), user.Name())

	if err != nil {
		return err
	}

	publishEvents(r.bus, user.GetEvents())

	return nil
}
