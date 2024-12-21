package auth

import (
	"context"
	"itam_auth/internal/database"
	"itam_auth/internal/models"
	"itam_auth/internal/services/jwt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx context.Context, storage *database.Storage, name, email, password string) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	userID := uuid.New()
	user := models.User{
		ID:            userID,
		Name:          name,
		Email:         email,
		PasswordHash:  string(hashedPassword),
		Specification: "Frontend",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = storage.SaveUser(ctx, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func AuthenticateUser(ctx context.Context, storage *database.Storage, email, password string) (string, error) {
	user, err := storage.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}

	tokenString, err := jwt.NewToken(user, 30*24*time.Hour)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
