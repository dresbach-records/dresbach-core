package middleware

import (
	"context"
	"net/http"
	"strings"

	"hosting-backend/internal/utils"
)

// ContextKey para o ID do usuário
const UserIDKey ContextKey = "user_id"

// AuthMiddleware verifica o token JWT e coloca o ID do usuário no contexto.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Cabeçalho de autorização não fornecido", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Unauthorized: Formato do cabeçalho de autorização inválido", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: Token inválido ou expirado", http.StatusUnauthorized)
			return
		}

		// Coloca o ID do usuário no contexto para uso posterior (ex: RBAC)
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
