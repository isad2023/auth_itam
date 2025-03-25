package middleware

import (
	"itam_auth/internal/services/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(hmacSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		var tokenString string

		// Проверяем, начинается ли заголовок с "Bearer "
		if strings.HasPrefix(authHeader, "Bearer ") {
			// Разбиваем на части и берем токен
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
				return
			}
			tokenString = parts[1]
		} else {
			// Считаем, что передан просто токен
			tokenString = authHeader
		}

		user, err := jwt.ValidateToken(tokenString, hmacSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token", "details": err.Error()})
			return
		}

		c.Set("user", user)
		c.Set("user_id", user.ID.String())

		c.Next()
	}
}
