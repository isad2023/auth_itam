package database

import (
	"context"
	"database/sql"
	"itam_auth/internal/models"

	"github.com/gofrs/uuid"
)

var (
	saveNewRequest      = `INSERT INTO requests (id, user_id, description, certificate, status, type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	getRequestsByUserID = `SELECT * FROM requests WHERE user_id = $1`
)

func SaveRequest(ctx context.Context, db *sql.DB, request models.Request) (uuid.UUID, error) {
	_, err := db.ExecContext(ctx, saveNewRequest, request.ID, request.UserID, request.Description, request.Certificate, request.Status, request.Type, request.CreatedAt)
	if err != nil {
		return uuid.Nil, err
	}
	return request.ID, nil
}

func GetRequests(ctx context.Context, db *sql.DB, userID uuid.UUID) ([]models.Request, error) {
	rows, err := db.QueryContext(ctx, getRequestsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.Request
	for rows.Next() {
		var request models.Request
		err := rows.Scan(&request.ID, &request.UserID, &request.Description, &request.Certificate, &request.Status, &request.Type, &request.CreatedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}
