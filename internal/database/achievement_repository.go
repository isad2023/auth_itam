package database

import (
	"context"
	"database/sql"
	"fmt"
	"itam_auth/internal/models"

	"github.com/google/uuid"
)

var (
	saveNewAchievement = `INSERT INTO achievements (id, title, description, points, approved, created_by, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	getAchievementByID = `SELECT * FROM achievements WHERE id = $1`
	getAllAchievements = `SELECT * FROM achievements`
	updateAchievement  = `UPDATE achievements SET title = $1, description = $2, points = $3, approved = $4, created_by = $5, created_at = $6 WHERE id = $7`
	deleteAchievement  = `DELETE FROM achievements WHERE id = $1`

	getAchievementsByUserID = `SELECT * FROM achievements WHERE created_by = $1`
)

func (s *Storage) SaveAchievement(ctx context.Context, achievement models.Achievement) (uuid.UUID, error) {
	_, err := s.db.ExecContext(ctx, saveNewAchievement, achievement.ID, achievement.Title, achievement.Description, achievement.Points, achievement.Approved, achievement.CreatedBy, achievement.CreatedAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to save achievement: %w", err)
	}
	return achievement.ID, nil
}

func (s *Storage) GetAchievementByID(ctx context.Context, id uuid.UUID) (models.Achievement, error) {
	row := s.db.QueryRowContext(ctx, getAchievementByID, id)

	var achievement models.Achievement
	err := row.Scan(&achievement.ID, &achievement.Title, &achievement.Description, &achievement.Points, &achievement.Approved, &achievement.CreatedBy, &achievement.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return achievement, fmt.Errorf("achievement not found: %w", err)
		}
		return achievement, fmt.Errorf("failed to get achievement: %w", err)
	}
	return achievement, nil
}

func (s *Storage) GetAllAchievements(ctx context.Context) ([]models.Achievement, error) {
	rows, err := s.db.QueryContext(ctx, getAllAchievements)
	if err != nil {
		return nil, fmt.Errorf("failed to get all achievements: %w", err)
	}
	defer rows.Close()

	var achievements []models.Achievement
	for rows.Next() {
		var achievement models.Achievement
		err := rows.Scan(&achievement.ID, &achievement.Title, &achievement.Description, &achievement.Points, &achievement.Approved, &achievement.CreatedBy, &achievement.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan achievement: %w", err)
		}
		achievements = append(achievements, achievement)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return achievements, nil
}

func (s *Storage) UpdateAchievement(ctx context.Context, achievement models.Achievement) error {
	_, err := s.db.ExecContext(ctx, updateAchievement, achievement.Title, achievement.Description, achievement.Points, achievement.Approved, achievement.CreatedBy, achievement.CreatedAt, achievement.ID)
	if err != nil {
		return fmt.Errorf("failed to update achievement: %w", err)
	}
	return nil
}

func (s *Storage) DeleteAchievement(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, deleteAchievement, id)
	if err != nil {
		return fmt.Errorf("failed to delete achievement: %w", err)
	}
	return nil
}

func (s *Storage) GetAchievementsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Achievement, error) {
	rows, err := s.db.QueryContext(ctx, getAchievementsByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get achievements by user ID: %w", err)
	}
	defer rows.Close()

	var achievements []models.Achievement
	for rows.Next() {
		var achievement models.Achievement
		err := rows.Scan(&achievement.ID, &achievement.Title, &achievement.Description, &achievement.Points, &achievement.Approved, &achievement.CreatedBy, &achievement.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan achievement: %w", err)
		}
		achievements = append(achievements, achievement)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return achievements, nil
}
