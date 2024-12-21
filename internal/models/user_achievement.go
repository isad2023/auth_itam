package models

import (
	"time"

	"github.com/google/uuid"
)

type UserAchievements struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	AchievementID uuid.UUID
	AwardedAt     time.Time
}
