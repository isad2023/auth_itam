package models

import (
	"time"

	"github.com/google/uuid"
)

// Achievement представляет достижение пользователя
type Achievement struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string    `json:"title" example:"First Project"`
	Description *string   `json:"description,omitempty" example:"Successfully completed first project"`
	Points      float64   `json:"points" example:"100.0"`
	Approved    bool      `json:"approved" example:"true"`
	ImageURL    *string   `json:"image_url,omitempty" example:"/uploads/achievement.jpg"`
	CreatedBy   int64     `json:"created_by" example:"123"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}
