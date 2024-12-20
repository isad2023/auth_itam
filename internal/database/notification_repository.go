package database

import (
	"context"
	"database/sql"
	"itam_auth/internal/models"

	"github.com/gofrs/uuid"
)

var (
	saveNewNotification        = `INSERT INTO notifications (id, user_id, content, is_read, created_at) VALUES ($1, $2, $3, $4, $5)`
	getNotificationByUserID    = `SELECT * FROM notifications WHERE user_id = $1`
	updateExistingNotification = `UPDATE notifications SET user_id = $1, content = $2, is_read = $3, created_at = $4 WHERE id = $5`
	getNotificationByID        = `SELECT * FROM notifications WHERE id = $1`
	getAllNotifications        = `SELECT * FROM notifications`
)

func (s *Storage) SaveNotification(ctx context.Context, notification models.Notification) (uuid.UUID, error) {
	id := notification.ID.UUID
	_, err := s.db.ExecContext(ctx, saveNewNotification, id, notification.UserID, notification.Content, notification.IsRead, notification.CreatedAt)
	if err != nil {
		return uuid.Nil, err
	}
	return notification.ID.UUID, nil
}

func (s *Storage) GetNotifications(ctx context.Context, userID uuid.UUID) ([]models.Notification, error) {
	rows, err := s.db.QueryContext(ctx, getNotificationByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(&notification.ID, &notification.UserID, &notification.Content, &notification.IsRead, &notification.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (s *Storage) GetAllNotifications(ctx context.Context) ([]models.Notification, error) {
	rows, err := s.db.QueryContext(ctx, getAllNotifications)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(&notification.ID, &notification.UserID, &notification.Content, &notification.IsRead, &notification.CreatedAt)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (s *Storage) UpdateNotification(ctx context.Context, notification models.Notification) error {
	_, err := s.db.ExecContext(ctx, updateExistingNotification, notification.UserID, notification.Content, notification.IsRead, notification.CreatedAt, notification.ID)
	return err
}

func (s *Storage) GetNotificationByID(ctx context.Context, notificationID uuid.UUID) (models.Notification, error) {
	row := s.db.QueryRowContext(ctx, getNotificationByID, notificationID)

	var notification models.Notification
	err := row.Scan(&notification.ID, &notification.UserID, &notification.Content, &notification.IsRead, &notification.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return notification, sql.ErrNoRows
		}
		return notification, err
	}
	return notification, nil
}
