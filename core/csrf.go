package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// generateCSRFToken генерирует CSRF-токен с использованием SHA-256.
// Возвращает токен в виде строки base64 или ошибку, если что-то пошло не так.
func generateCSRFToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", errors.New("failed to generate random bytes")
	}

	hash := sha256.Sum256(randomBytes)
	token := base64.StdEncoding.EncodeToString(hash[:])
	return token, nil
}
