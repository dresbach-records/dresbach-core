package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

// Claims representa as reivindicações personalizadas para o token JWT.
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// contextKey é um tipo usado para definir chaves de contexto de forma segura.
type contextKey string

const claimsContextKey = contextKey("claims")

// GetClaims extrai as reivindicações JWT do contexto da solicitação.
func GetClaims(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	return claims, ok
}
