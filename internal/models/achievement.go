package models

import (
	"time"

	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type Achievement struct {
	ID          uuid.UUID
	Title       string
	Description string
	Points      float64
	Approved    bool
	CreatedBy   int64
	CreatedAt   time.Time
}
