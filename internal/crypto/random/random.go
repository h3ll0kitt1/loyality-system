package random

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateSalt() (string, error) {
	salt, err := generateRandomString(16)
	if err != nil {
		return ``, fmt.Errorf("generate salt failed: %w", err)
	}

	return salt, nil
}

func GenerateSecretKey() (string, error) {
	secret, err := generateRandomString(128)
	if err != nil {
		return ``, fmt.Errorf("generate secret key failed: %w", err)
	}

	return secret, nil
}

func generateRandomString(size int) (string, error) {

	b := make([]byte, size)

	_, err := rand.Read(b)
	if err != nil {
		return ``, fmt.Errorf("generate random string failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
