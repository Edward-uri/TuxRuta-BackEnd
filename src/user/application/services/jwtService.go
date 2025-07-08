package services

type JWTService interface {
	GenerateToken(userID int, email string, rol string) (string, error)
	ValidateToken(token string) (int, error)
}

type TokenClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Rol    string `json:"rol"`
	Exp    int64  `json:"exp"`
}
