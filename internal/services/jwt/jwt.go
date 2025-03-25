package jwt

import (
	"fmt"
	"itam_auth/internal/models"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewToken(user models.User, duration time.Duration, hmacSecret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(hmacSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, hmacSecret string) (models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(hmacSecret), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return models.User{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return models.User{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return models.User{}, fmt.Errorf("invalid claims")
	}

	var authUser models.User
	authUser.ID = uuid.MustParse(claims.UID)
	authUser.Email = claims.Email
	return authUser, nil
}

func NewRefreshToken(user models.User, hmacSecret string) (string, error) {
	claims := jwt.MapClaims{
		"uid": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(hmacSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}
