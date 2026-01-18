package middleware

import (
	"net/http"
)

// RBACMiddleware é um middleware para controle de acesso baseado em roles.
func RBACMiddleware(next http.Handler, allowedRoles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("role").(string)
		if !ok {
			http.Error(w, "Role não encontrado no contexto", http.StatusInternalServerError)
			return
		}

		allowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			http.Error(w, "Acesso negado", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
