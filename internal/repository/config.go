package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/RaphaelAZ/go-wordle/internal/models"
)

type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

func (r *ConfigRepository) Get(userID int) (*models.UserConfig, error) {
	cfg := &models.UserConfig{}
	err := r.db.QueryRow(
		`SELECT id, user_id, config, state, updated_at FROM user_configs WHERE user_id = $1`,
		userID,
	).Scan(&cfg.ID, &cfg.UserID, &cfg.Config, &cfg.State, &cfg.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get config: %w", err)
	}
	return cfg, nil
}

func (r *ConfigRepository) Upsert(userID int, config, state []byte) (*models.UserConfig, error) {
	cfg := &models.UserConfig{}
	err := r.db.QueryRow(
		`INSERT INTO user_configs (user_id, config, state, updated_at)
		 VALUES ($1, $2, $3, NOW())
		 ON CONFLICT (user_id) DO UPDATE
		   SET config = EXCLUDED.config,
		       state  = EXCLUDED.state,
		       updated_at = NOW()
		 RETURNING id, user_id, config, state, updated_at`,
		userID, config, state,
	).Scan(&cfg.ID, &cfg.UserID, &cfg.Config, &cfg.State, &cfg.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("upsert config: %w", err)
	}
	return cfg, nil
}
