package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/h3ll0kitt1/loyality-system/internal/config"
)

type Claims struct {
	Login string
	jwt.RegisteredClaims
}

func GenerateToken(login string, cfg *config.Config) (string, error) {

	claims := &Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.TokenExpire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.SecretKey))
	if err != nil {
		return "", fmt.Errorf("generate JWT token failed: %w", err)
	}

	return tokenString, nil
}

func CheckToken(tokenString string, SecretKey string) (string, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

	// token, err := jwt.ParseWithClaims(tokenString, claims,
	// 	func(t *jwt.Token) (interface{}, error) {
	// 		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
	// 			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	// 		}
	// 		return []byte(SECRET_KEY), nil
	// 	})

	if err != nil {
		return "", fmt.Errorf("check JWT tocken failed: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("check JWT tocken is not valid: %w", err)
	}

	return claims.Login, nil
}
