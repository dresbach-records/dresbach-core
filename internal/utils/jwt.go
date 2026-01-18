package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JwtClaims define os dados (claims) que serão armazenados no token JWT.
type JwtClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GetJwtSecret busca o segredo JWT da variável de ambiente.
// É crucial que esta variável esteja definida no ambiente de produção.
func GetJwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Para ambientes de desenvolvimento, um segredo padrão pode ser usado.
		// NUNCA use isso em produção.
		return []byte("default_super_secret_key_for_dev_only")
	}
	return []byte(secret)
}

// GenerateToken cria um novo token JWT para um usuário.
func GenerateToken(userID int, role string) (string, error) {
	claims := JwtClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expira em 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "hosting-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJwtSecret())
}

// ValidateToken verifica a validade de um token string.
// Retorna os claims se o token for válido, caso contrário, retorna um erro.
func ValidateToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verifica se o método de assinatura é o esperado (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return GetJwtSecret(), nil
	})

	if err != nil {
		return nil, err // O erro pode ser de expiração, assinatura inválida, etc.
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("token inválido")
	}
}
