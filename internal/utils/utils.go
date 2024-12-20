package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"itam_auth/internal/config"
	"log"
	"net/url"
	"sort"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey []byte

func init() {
	appConfig, _ := config.LoadConfig()
	// Загружаем JWT секретный ключ из переменной окружения
	jwtKey = []byte(appConfig.JwtSecretKey)
	if len(jwtKey) == 0 {
		log.Fatal("JWT_SECRET_KEY is not set in environment variables.")
	}
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorSignatureInvalid)
	}

	return claims, nil
}

// ValidateTelegramAuth проверяет данные, пришедшие от Telegram
func ValidateTelegramAuth(data url.Values, botToken string) bool {
	// Извлекаем hash из данных
	receivedHash := data.Get("hash")
	data.Del("hash") // Убираем hash из проверяемых данных

	// Формируем строку из оставшихся параметров
	var dataCheckString string
	var keys []string
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		dataCheckString += key + "=" + data.Get(key) + "\n"
	}
	dataCheckString = dataCheckString[:len(dataCheckString)-1] // Убираем последний перевод строки

	// Генерируем ключ из токена бота
	secretKey := sha256.Sum256([]byte(botToken))
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	// Сравниваем рассчитанный hash с переданным
	return calculatedHash == receivedHash
}
