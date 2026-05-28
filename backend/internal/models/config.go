package models

import (
	"encoding/json"
	"time"
)

type UserConfig struct {
	ID        int             `json:"id"`
	UserID    int             `json:"user_id"`
	Config    json.RawMessage `json:"config"`
	State     json.RawMessage `json:"state"`
	UpdatedAt time.Time       `json:"updated_at"`
}
