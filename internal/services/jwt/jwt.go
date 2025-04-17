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
	UID           string   `json:"uid"`
	Email         string   `json:"email"`
	AdminServices []string `json:"admin_services"`
	jwt.RegisteredClaims
}

func NewToken(user models.User, duration time.Duration, hmacSecret string, userRoles []models.UserRole,
	roles []models.Role, rolePermissions []models.RolePermission, permissions []models.Permission) (string, error) {
	claims := Claims{
		UID:           user.ID.String(),
		Email:         user.Email,
		AdminServices: user.GetAdminServices(userRoles, roles, rolePermissions, permissions),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

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

func NewRefreshToken(user models.User, hmacSecret string, userRoles []models.UserRole,
	roles []models.Role, rolePermissions []models.RolePermission, permissions []models.Permission) (string, error) {
	claims := jwt.MapClaims{
		"uid":            user.ID.String(),
		"admin_services": user.GetAdminServices(userRoles, roles, rolePermissions, permissions),
		"exp":            time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(hmacSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}
