package database

import (
	"context"
	"database/sql"
	"fmt"
	"itam_auth/internal/models"
	"log"

	"github.com/google/uuid"
)

var (
	saveNewNotification = `INSERT INTO notifications 
	(id, user_id, content, is_read, created_at) 
	VALUES ($1, $2, $3, $4, $5)`
	getNotificationByUserID    = `SELECT * FROM notifications WHERE user_id = $1 LIMIT $2 OFFSET $3`
	updateExistingNotification = `UPDATE notifications SET user_id = $1, content = $2, is_read = $3, created_at = $4 WHERE id = $5`
	getNotificationByID        = `SELECT * FROM notifications WHERE id = $1`
	getAllNotifications        = `SELECT * FROM notifications LIMIT $1 OFFSET $2`
)

func validateNotification(notification models.Notification) error {
	if notification.ID == uuid.Nil {
		return fmt.Errorf("notification ID cannot be empty")
	}
	if notification.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}
	if notification.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	if notification.CreatedAt.IsZero() {
		return fmt.Errorf("created_at must be set")
	}
	return nil
}

func scanNotification(row interface{ Scan(...any) error }) (models.Notification, error) {
	var notification models.Notification
	err := row.Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Content,
		&notification.IsRead,
		&notification.CreatedAt,
	)
	if err != nil {
		return models.Notification{}, fmt.Errorf("failed to scan notification: %w", err)
	}
	return notification, nil
}

func (s *Storage) SaveNotification(ctx context.Context, notification models.Notification) (uuid.UUID, error) {
	if err := validateNotification(notification); err != nil {
		log.Printf("Validation failed for notification with ID %s: %v", notification.ID, err)
		return uuid.Nil, fmt.Errorf("invalid notification data: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction for saving notification with ID %s: %v", notification.ID, err)
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Failed to rollback transaction for notification with ID %s: %v", notification.ID, err)
		}
	}()

	_, err = tx.ExecContext(ctx, saveNewNotification,
		notification.ID,
		notification.UserID,
		notification.Content,
		notification.IsRead,
		notification.CreatedAt,
	)
	if err != nil {
		log.Printf("Failed to save notification with ID %s: %v", notification.ID, err)
		return uuid.Nil, fmt.Errorf("failed to save notification: %w", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction for notification with ID %s: %v", notification.ID, err)
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return notification.ID, nil
}

func (s *Storage) GetNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be empty")
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := s.db.QueryContext(ctx, getNotificationByUserID, userID, limit, offset)
	if err != nil {
		log.Printf("Failed to get notifications for user ID %s (limit=%d, offset=%d): %v", userID, limit, offset, err)
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows in GetNotifications: %v", closeErr)
		}
	}()

	var notifications []models.Notification
	for rows.Next() {
		notification, err := scanNotification(rows)
		if err != nil {
			log.Printf("Failed to scan notification in GetNotifications: %v", err)
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration in GetNotifications: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return notifications, nil
}

func (s *Storage) GetAllNotifications(ctx context.Context, limit, offset int) ([]models.Notification, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := s.db.QueryContext(ctx, getAllNotifications, limit, offset)
	if err != nil {
		log.Printf("Failed to get all notifications (limit=%d, offset=%d): %v", limit, offset, err)
		return nil, fmt.Errorf("failed to get all notifications: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows in GetAllNotifications: %v", closeErr)
		}
	}()

	var notifications []models.Notification
	for rows.Next() {
		notification, err := scanNotification(rows)
		if err != nil {
			log.Printf("Failed to scan notification in GetAllNotifications: %v", err)
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration in GetAllNotifications: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return notifications, nil
}

func (s *Storage) UpdateNotification(ctx context.Context, notification models.Notification) error {
	if err := validateNotification(notification); err != nil {
		log.Printf("Validation failed for notification with ID %s: %v", notification.ID, err)
		return fmt.Errorf("invalid notification data: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction for updating notification with ID %s: %v", notification.ID, err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Failed to rollback transaction for notification with ID %s: %v", notification.ID, err)
		}
	}()

	result, err := tx.ExecContext(ctx, updateExistingNotification,
		notification.UserID,
		notification.Content,
		notification.IsRead,
		notification.CreatedAt,
		notification.ID,
	)
	if err != nil {
		log.Printf("Failed to update notification with ID %s: %v", notification.ID, err)
		return fmt.Errorf("failed to update notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected for notification with ID %s: %v", notification.ID, err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no notification found with ID: %s", notification.ID)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction for notification with ID %s: %v", notification.ID, err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Storage) GetNotificationByID(ctx context.Context, notificationID uuid.UUID) (models.Notification, error) {
	if notificationID == uuid.Nil {
		return models.Notification{}, fmt.Errorf("notification ID cannot be empty")
	}

	row := s.db.QueryRowContext(ctx, getNotificationByID, notificationID)

	notification, err := scanNotification(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Notification{}, err
		}
		log.Printf("Failed to get notification with ID %s: %v", notificationID, err)
		return models.Notification{}, fmt.Errorf("failed to get notification: %w", err)
	}

	return notification, nil
}
