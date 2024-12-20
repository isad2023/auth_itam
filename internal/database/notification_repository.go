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
)

func SaveNotification(ctx context.Context, db *sql.DB, notification models.Notification) (uuid.UUID, error) {
	_, err := db.ExecContext(ctx, saveNewNotification, notification.ID, notification.UserID, notification.Content, notification.IsRead, notification.CreatedAt)
	if err != nil {
		return uuid.Nil, err
	}
	return notification.ID.UUID, nil
}

func GetNotifications(ctx context.Context, db *sql.DB, userID uuid.UUID) ([]models.Notification, error) {
	rows, err := db.QueryContext(ctx, getNotificationByUserID, userID)
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

func UpdateNotification(ctx context.Context, db *sql.DB, notification models.Notification) error {
	_, err := db.ExecContext(ctx, updateExistingNotification, notification.UserID, notification.Content, notification.IsRead, notification.CreatedAt, notification.ID)
	return err
}

func GetNotificationByID(ctx context.Context, db *sql.DB, notificationID uuid.UUID) (models.Notification, error) {
	row := db.QueryRowContext(ctx, getNotificationByID, notificationID)

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
