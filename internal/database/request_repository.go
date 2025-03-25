package database

import (
	"context"
	"database/sql"
	"fmt"
	"itam_auth/internal/models"
	"log"
	"time"

	"github.com/google/uuid"
)

var (
	saveNewRequest = `INSERT INTO requests 
	(id, user_id, description, certificate, status, type, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	getRequestsByUserID = `SELECT * FROM requests WHERE user_id = $1 LIMIT $2 OFFSET $3`
	updateRequest       = `UPDATE requests SET status = $1, updated_at = $2 WHERE id = $3`
	deleteRequest       = `DELETE FROM requests WHERE id = $1`
)

var ValidRequestStatuses = map[string]bool{
	"pending":  true,
	"approved": true,
	"rejected": true,
}

func validateRequest(request models.Request) error {
	if request.ID == uuid.Nil {
		return fmt.Errorf("request ID cannot be empty")
	}
	if request.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}
	if request.Description == "" {
		return fmt.Errorf("description cannot be empty")
	}
	if request.Status == "" {
		return fmt.Errorf("status cannot be empty")
	}
	if !ValidRequestStatuses[request.Status] {
		return fmt.Errorf("invalid status: %s, must be one of %v", request.Status, ValidRequestStatuses)
	}
	if request.Type == "" {
		return fmt.Errorf("type cannot be empty")
	}
	if request.CreatedAt.IsZero() {
		return fmt.Errorf("created_at must be set")
	}
	return nil
}

func scanRequest(row interface{ Scan(...any) error }) (models.Request, error) {
	var request models.Request
	err := row.Scan(
		&request.ID,
		&request.UserID,
		&request.Description,
		&request.Certificate,
		&request.Status,
		&request.Type,
		&request.CreatedAt,
	)
	if err != nil {
		return models.Request{}, fmt.Errorf("failed to scan request: %w", err)
	}
	return request, nil
}

func (s *Storage) SaveRequest(ctx context.Context, request models.Request) (uuid.UUID, error) {
	if err := validateRequest(request); err != nil {
		log.Printf("Validation failed for request with ID %s: %v", request.ID, err)
		return uuid.Nil, fmt.Errorf("invalid request data: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction for saving request with ID %s: %v", request.ID, err)
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Failed to rollback transaction for request with ID %s: %v", request.ID, err)
		}
	}()

	_, err = tx.ExecContext(ctx, saveNewRequest,
		request.ID,
		request.UserID,
		request.Description,
		request.Certificate,
		request.Status,
		request.Type,
		request.CreatedAt,
	)
	if err != nil {
		log.Printf("Failed to save request with ID %s: %v", request.ID, err)
		return uuid.Nil, fmt.Errorf("failed to save request: %w", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction for request with ID %s: %v", request.ID, err)
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Request saved successfully with ID %s", request.ID)
	return request.ID, nil
}

func (s *Storage) GetRequests(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Request, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := s.db.QueryContext(ctx, getRequestsByUserID, userID, limit, offset)
	if err != nil {
		log.Printf("Failed to get requests for user ID %s (limit=%d, offset=%d): %v", userID, limit, offset, err)
		return nil, fmt.Errorf("failed to get requests: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows in GetRequests: %v", closeErr)
		}
	}()

	var requests []models.Request
	for rows.Next() {
		request, err := scanRequest(rows)
		if err != nil {
			log.Printf("Failed to scan request in GetRequests: %v", err)
			return nil, err
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration in GetRequests: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return requests, nil
}

func (s *Storage) UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error {
	if requestID == uuid.Nil {
		return fmt.Errorf("request ID cannot be empty")
	}
	if !ValidRequestStatuses[status] {
		return fmt.Errorf("invalid status: %s, must be one of %v", status, ValidRequestStatuses)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction for updating request with ID %s: %v", requestID, err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Failed to rollback transaction for request with ID %s: %v", requestID, err)
		}
	}()

	updatedAt := time.Now()
	result, err := tx.ExecContext(ctx, updateRequest, status, updatedAt, requestID)
	if err != nil {
		log.Printf("Failed to update request status with ID %s: %v", requestID, err)
		return fmt.Errorf("failed to update request status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected for request with ID %s: %v", requestID, err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no request found with ID: %s", requestID)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction for request with ID %s: %v", requestID, err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Request status updated successfully for ID %s to %s", requestID, status)
	return nil
}

func (s *Storage) DeleteRequest(ctx context.Context, requestID uuid.UUID) error {
	if requestID == uuid.Nil {
		return fmt.Errorf("request ID cannot be empty")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction for deleting request with ID %s: %v", requestID, err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Failed to rollback transaction for deleting request with ID %s: %v", requestID, err)
		}
	}()

	result, err := tx.ExecContext(ctx, deleteRequest, requestID)
	if err != nil {
		log.Printf("Failed to delete request with ID %s: %v", requestID, err)
		return fmt.Errorf("failed to delete request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected for deleting request with ID %s: %v", requestID, err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no request found with ID: %s", requestID)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction for deleting request with ID %s: %v", requestID, err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Request deleted successfully for ID %s", requestID)
	return nil
}
