package models

import (
	"time"

	"github.com/google/uuid"
)

type Achievement struct {
	ID          uuid.UUID
	Title       string
	Description *string
	Points      float64
	Approved    bool
	CreatedBy   int64
	CreatedAt   time.Time
}
