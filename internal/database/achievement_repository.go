package database

import (
	"context"
	"database/sql"
	"fmt"
	"itam_auth/internal/models"
	"log"

	"github.com/google/uuid"
)

const (
	saveAchievementQuery = `INSERT INTO achievements 
		(id, title, description, points, approved, image_url, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	getAchievementByIDQuery = `SELECT id, title, description, points, approved, image_url, created_by, created_at 
		FROM achievements WHERE id = $1`
	getAllAchievementsQuery = `SELECT id, title, description, points, approved, image_url, created_by, created_at 
		FROM achievements LIMIT $1 OFFSET $2`
	updateAchievementQuery = `UPDATE achievements 
		SET title = $1, description = $2, points = $3, approved = $4, image_url = $5, created_by = $6, created_at = $7 
		WHERE id = $8`
	deleteAchievementQuery       = `DELETE FROM achievements WHERE id = $1`
	getAchievementsByUserIDQuery = `
		SELECT achievements.id, achievements.title, achievements.description, achievements.points, 
			achievements.approved, achievements.image_url, achievements.created_by, achievements.created_at 
		FROM achievements
		JOIN user_achievements ON achievements.id = user_achievements.achievement_id
		WHERE user_achievements.user_id = $1
		LIMIT $2 OFFSET $3`
)

func validateAchievement(achievement models.Achievement) error {
	if achievement.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if achievement.Points < 0 || achievement.Points > 1000 {
		return fmt.Errorf("points must be between 0 and 1000, got %f", achievement.Points)
	}
	return nil
}

func scanAchievement(row interface{ Scan(...any) error }) (models.Achievement, error) {
	var achievement models.Achievement
	var description sql.NullString
	var imageURL sql.NullString

	err := row.Scan(
		&achievement.ID,
		&achievement.Title,
		&description,
		&achievement.Points,
		&achievement.Approved,
		&imageURL,
		&achievement.CreatedBy,
		&achievement.CreatedAt,
	)
	if err != nil {
		return models.Achievement{}, fmt.Errorf("failed to scan achievement: %w", err)
	}

	if description.Valid {
		achievement.Description = &description.String
	} else {
		achievement.Description = nil
	}

	if imageURL.Valid {
		achievement.ImageURL = &imageURL.String
	} else {
		achievement.ImageURL = nil
	}

	return achievement, nil
}

func (s *Storage) SaveAchievement(ctx context.Context, achievement models.Achievement) (uuid.UUID, error) {
	if err := validateAchievement(achievement); err != nil {
		log.Printf("Validation failed for achievement with ID %s: %v", achievement.ID, err)
		return uuid.Nil, fmt.Errorf("invalid achievement data: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction for saving achievement with ID %s: %v", achievement.ID, err)
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Failed to rollback transaction for achievement with ID %s: %v", achievement.ID, err)
		}
	}()

	_, err = tx.ExecContext(ctx, saveAchievementQuery,
		achievement.ID,
		achievement.Title,
		achievement.Description,
		achievement.Points,
		achievement.Approved,
		achievement.ImageURL,
		achievement.CreatedBy,
		achievement.CreatedAt,
	)
	if err != nil {
		log.Printf("Failed to save achievement with ID %s: %v", achievement.ID, err)
		return uuid.Nil, fmt.Errorf("failed to save achievement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction for achievement with ID %s: %v", achievement.ID, err)
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return achievement.ID, nil
}

func (s *Storage) GetAchievementByID(ctx context.Context, id uuid.UUID) (models.Achievement, error) {
	row := s.db.QueryRowContext(ctx, getAchievementByIDQuery, id)

	achievement, err := scanAchievement(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Achievement{}, err
		}
		log.Printf("Failed to get achievement with ID %s: %v", id, err)
		return models.Achievement{}, fmt.Errorf("failed to get achievement: %w", err)
	}

	return achievement, nil
}

func (s *Storage) GetAllAchievements(ctx context.Context, limit, offset int) ([]models.Achievement, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	rows, err := s.db.QueryContext(ctx, getAllAchievementsQuery, limit, offset)
	if err != nil {
		log.Printf("Failed to get all achievements (limit=%d, offset=%d): %v", limit, offset, err)
		return nil, fmt.Errorf("failed to get all achievements: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows in GetAllAchievements: %v", closeErr)
		}
	}()

	var achievements []models.Achievement
	for rows.Next() {
		achievement, err := scanAchievement(rows)
		if err != nil {
			log.Printf("Failed to scan achievement in GetAllAchievements: %v", err)
			return nil, err
		}
		achievements = append(achievements, achievement)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration in GetAllAchievements: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return achievements, nil
}

func (s *Storage) UpdateAchievement(ctx context.Context, achievement models.Achievement) error {

	if err := validateAchievement(achievement); err != nil {
		log.Printf("Validation failed for achievement with ID %s: %v", achievement.ID, err)
		return fmt.Errorf("invalid achievement data: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction for updating achievement with ID %s: %v", achievement.ID, err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Failed to rollback transaction for achievement with ID %s: %v", achievement.ID, err)
		}
	}()

	result, err := tx.ExecContext(ctx, updateAchievementQuery,
		achievement.Title,
		achievement.Description,
		achievement.Points,
		achievement.Approved,
		achievement.ImageURL,
		achievement.CreatedBy,
		achievement.CreatedAt,
		achievement.ID,
	)
	if err != nil {
		log.Printf("Failed to update achievement with ID %s: %v", achievement.ID, err)
		return fmt.Errorf("failed to update achievement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected for achievement with ID %s: %v", achievement.ID, err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no achievement found with ID: %s", achievement.ID)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction for achievement with ID %s: %v", achievement.ID, err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Storage) DeleteAchievement(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, deleteAchievementQuery, id)
	if err != nil {
		log.Printf("Failed to delete achievement with ID %s: %v", id, err)
		return fmt.Errorf("failed to delete achievement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected for achievement with ID %s: %v", id, err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no achievement found with ID: %s", id)
	}

	return nil
}

func (s *Storage) GetAchievementsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Achievement, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	rows, err := s.db.QueryContext(ctx, getAchievementsByUserIDQuery, userID, limit, offset)
	if err != nil {
		log.Printf("Failed to get achievements for user ID %s (limit=%d, offset=%d): %v", userID, limit, offset, err)
		return nil, fmt.Errorf("failed to get achievements by user ID: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows in GetAchievementsByUserID: %v", closeErr)
		}
	}()

	var achievements []models.Achievement
	for rows.Next() {
		achievement, err := scanAchievement(rows)
		if err != nil {
			log.Printf("Failed to scan achievement in GetAchievementsByUserID: %v", err)
			return nil, err
		}
		achievements = append(achievements, achievement)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration in GetAchievementsByUserID: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return achievements, nil
}
