package database

import (
	"context"
	"itam_auth/internal/models"
	"time"

	"github.com/google/uuid"
)

var (
	saveNewRequest      = `INSERT INTO requests (id, user_id, description, certificate, status, type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	getRequestsByUserID = `SELECT * FROM requests WHERE user_id = $1`
	updateRequest       = `UPDATE requests SET status = $1, updated_at = $2 WHERE id = $3`
)

func (s *Storage) SaveRequest(ctx context.Context, request models.Request) (uuid.UUID, error) {
	_, err := s.db.ExecContext(ctx, saveNewRequest, request.ID, request.UserID, request.Description, request.Certificate, request.Status, request.Type, request.CreatedAt)
	if err != nil {
		return uuid.Nil, err
	}
	return request.ID, nil
}

func (s *Storage) GetRequests(ctx context.Context, userID uuid.UUID) ([]models.Request, error) {
	rows, err := s.db.QueryContext(ctx, getRequestsByUserID, userID)
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

func (s *Storage) UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error {
	_, err := s.db.ExecContext(ctx, updateRequest, status, time.Now(), requestID)
	return err
}
