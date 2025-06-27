package models

import (
	"time"

	"github.com/google/uuid"
)

type Achievement struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Points      float64    `json:"points"`
	Approved    bool       `json:"approved"`
	CreatedBy   int64      `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
}
