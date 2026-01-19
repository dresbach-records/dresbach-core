package middleware

import (
	"context"
	"database/sql"
	"net/http"

	"hosting-backend/internal/models"
)

// ContextKey é um tipo usado para chaves de contexto para evitar colisões.
type ContextKey string

// Chaves usadas para armazenar valores no contexto da requisição.
const (
	UserPermissionsKey ContextKey = "user_permissions"
	// UserIDKey          ContextKey = "user_id" // Movido para auth.go para evitar re-declaração
)

// RBACMiddleware carrega as permissões do usuário e as coloca no contexto.
func RBACMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obter o ID do usuário que foi colocado no contexto pelo AuthMiddleware.
			userID, ok := r.Context().Value(UserIDKey).(int)
			if !ok {
				// Se não houver ID de usuário, não há permissões para carregar.
				// O AuthMiddleware já deve ter bloqueado a requisição se ela for protegida.
				// Apenas passamos para o próximo handler.
				next.ServeHTTP(w, r)
				return
			}

			// Buscar as permissões do usuário logado.
			permissions, err := models.GetUserPermissions(db, userID)
			if err != nil {
				// Um erro aqui é um problema no servidor, pois o usuário está autenticado.
				http.Error(w, "Erro interno ao buscar permissões do usuário", http.StatusInternalServerError)
				return
			}

			// Adicionar as permissões ao contexto para o RequirePermission usar.
			ctx := context.WithValue(r.Context(), UserPermissionsKey, permissions)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequirePermission verifica se o usuário tem a permissão necessária.
func RequirePermission(permissionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			permissions, ok := r.Context().Value(UserPermissionsKey).(map[string]bool)
			if !ok {
				// Isso acontece se o RBACMiddleware não conseguiu carregar as permissões,
				// ou se o usuário não está logado.
				http.Error(w, "Forbidden: Permissões de acesso não disponíveis.", http.StatusForbidden)
				return
			}

			if _, hasPermission := permissions[permissionName]; !hasPermission {
				http.Error(w, "Forbidden: Você não tem permissão para realizar esta ação.", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
