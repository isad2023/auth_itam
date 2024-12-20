package database

import (
	"context"
	"database/sql"
	"itam_auth/internal/models"

	"github.com/gofrs/uuid"
)

var (
	saveNewNotification     = `INSERT INTO notifications (id, user_id, content, is_read, created_at) VALUES ($1, $2, $3, $4, $5)`
	getNotificationByUserID = `SELECT * FROM notifications WHERE user_id = $1`
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
