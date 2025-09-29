package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-this-in-production") // TODO: mover para config

type Claims struct {
	ID       int    `json:"id"`
	UserUUID string `json:"uuid"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT gera um token JWT para o usuário
func GenerateJWT(id int, uuid, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token válido por 24h

	claims := &Claims{
		ID:       id,
		UserUUID: uuid,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ratatoskr",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT valida e extrai os claims de um token JWT
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verificar se o método de assinatura é o esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de assinatura inválido")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

// SetJWTSecret permite configurar o secret do JWT (deve ser chamado no startup)
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}
