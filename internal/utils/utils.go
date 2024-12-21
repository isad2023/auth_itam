package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
)

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
