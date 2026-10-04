package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"planningpoker/internal/domain/events"
	"planningpoker/internal/domain/games"
)

// PostgresGameRepository stores each game aggregate as a JSONB document.
type PostgresGameRepository struct {
	db  *sql.DB
	bus events.EventBus
}

var _ games.GameRepository = (*PostgresGameRepository)(nil)

// NewPostgresGameRepository creates a repository using an initialized database.
func NewPostgresGameRepository(db *sql.DB, bus events.EventBus) *PostgresGameRepository {
	return &PostgresGameRepository{
		db:  db,
		bus: bus,
	}
}

// Get retrieves a game, returning nil when it does not exist.
func (r *PostgresGameRepository) Get(id string) (*games.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var raw []byte

	err := r.db.QueryRowContext(ctx, "SELECT document FROM games WHERE id = $1", id).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return decodeGame(raw)
}

// Save persists a game and publishes its events only after the write succeeds.
// Updates to an existing game should use ModifyExclusively to avoid lost updates.
func (r *PostgresGameRepository) Save(game *games.Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	raw, err := json.Marshal(newGameDTO(game))
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `INSERT INTO games (id, document) VALUES ($1, $2::jsonb)
		ON CONFLICT (id) DO UPDATE SET document = EXCLUDED.document, updated_at = now()`, game.ID(), string(raw))

	if err != nil {
		return err
	}

	publishEvents(r.bus, game.GetEvents())
	game.ClearEvents()

	return nil
}

// ModifyExclusively locks this game's row through read, modification, and commit.
// Independent connections and processes cannot overwrite concurrent votes.
func (r *PostgresGameRepository) ModifyExclusively(id string, cb func(*games.Game) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback() }()
	var raw []byte

	err = tx.QueryRowContext(ctx, "SELECT document FROM games WHERE id = $1 FOR UPDATE", id).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("game not found")
	}

	if err != nil {
		return fmt.Errorf("lock game: %w", err)
	}

	game, err := decodeGame(raw)
	if err != nil {
		return err
	}

	if err := cb(game); err != nil {
		return err
	}

	raw, err = json.Marshal(newGameDTO(game))
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "UPDATE games SET document = $2::jsonb, updated_at = now() WHERE id = $1", id, string(raw)); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	publishEvents(r.bus, game.GetEvents())
	game.ClearEvents()

	return nil
}

// GetActiveGamesByPlayerID finds started games containing a player, matching the memory repository.
func (r *PostgresGameRepository) GetActiveGamesByPlayerID(playerID string) ([]games.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `SELECT document FROM games
		WHERE document->>'state' = 'started' AND (document->'players') ? $1 ORDER BY id`, playerID)

	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	list := make([]games.Game, 0)

	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}

		game, err := decodeGame(raw)

		if err != nil {
			return nil, err
		}

		list = append(list, *game)
	}

	return list, rows.Err()
}

func decodeGame(raw []byte) (*games.Game, error) {
	var dto gameDTO
	if err := json.Unmarshal(raw, &dto); err != nil {
		return nil, err
	}

	return dto.toDomain()
}
