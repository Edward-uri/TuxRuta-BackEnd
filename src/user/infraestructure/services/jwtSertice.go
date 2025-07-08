package services

import (
	"errors"
	"main/src/user/application/services"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ✅ Implementación concreta en infraestructura
type JWTServiceImpl struct {
	secretKey []byte
}

// ✅ Constructor que retorna la interfaz del dominio
func NewJWTService() services.JWTService { // ✅ Cambiar a services.JWTService
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}
	return &JWTServiceImpl{
		secretKey: []byte(secret),
	}
}

// ✅ Claims para JWT
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Rol    string `json:"rol"`
	jwt.RegisteredClaims
}

// ✅ Implementar GenerateToken
func (s *JWTServiceImpl) GenerateToken(userID int, email, rol string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Rol:    rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "tuxruta-api",
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ✅ Implementar ValidateToken
func (s *JWTServiceImpl) ValidateToken(tokenString string) (int, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return 0, errors.New("token expired")
	}

	return claims.UserID, nil
}
