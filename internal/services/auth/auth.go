package auth

import (
	"context"
	"database/sql"
	"fmt"
	"itam_auth/internal/database"
	"itam_auth/internal/models"
	"itam_auth/internal/services/jwt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultUserSpecification = "Frontend"          // Спецификация пользователя по умолчанию
	tokenDuration            = 30 * 24 * time.Hour // Длительность действия токена (30 дней)
	minPasswordLength        = 8                   // Минимальная длина пароля
	bcryptCost               = bcrypt.DefaultCost  // Стоимость хеширования пароля
)

func validateUserData(name, email, password string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("invalid email format")
	}

	if len(password) < minPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", minPasswordLength)
	}

	return nil
}

func RegisterUser(ctx context.Context, storage *database.Storage, name, email, password string) (models.User, error) {
	if err := validateUserData(name, email, password); err != nil {
		log.Printf("Validation failed for user registration (email=%s): %v", email, err)
		return models.User{}, fmt.Errorf("invalid user data: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password for user (email=%s): %v", email, err)
		return models.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	userID := uuid.New()
	user := models.User{
		ID:            userID,
		Name:          name,
		Email:         email,
		PasswordHash:  string(hashedPassword),
		Specification: defaultUserSpecification,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = storage.SaveUser(ctx, user)
	if err != nil {
		log.Printf("Failed to save user (email=%s, id=%s): %v", email, userID, err)
		return models.User{}, fmt.Errorf("failed to save user: %w", err)
	}

	log.Printf("User registered successfully (email=%s, id=%s)", email, userID)
	return user, nil
}

func AuthenticateUser(ctx context.Context, storage *database.Storage, email, password string) (string, error) {
	if strings.TrimSpace(email) == "" {
		return "", fmt.Errorf("email cannot be empty")
	}
	if strings.TrimSpace(password) == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	user, err := storage.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("Failed to get user by email (email=%s): %v", email, err)
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Printf("Invalid password for user (email=%s, id=%s)", email, user.ID)
		return "", fmt.Errorf("invalid password: %w", err)
	}

	tokenString, err := jwt.NewToken(user, tokenDuration)
	if err != nil {
		log.Printf("Failed to generate JWT token for user (email=%s, id=%s): %v", email, user.ID, err)
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	log.Printf("User authenticated successfully (email=%s, id=%s)", email, user.ID)
	return tokenString, nil
}
