package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// User represents the user data to be returned by the /auth/me endpoint.
	// User representa os dados do usuário a serem retornados pelo endpoint /auth/me.
type User struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// MeHandler processa as solicitações para o endpoint /auth/me.
func MeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extrai o email e o role do contexto da requisição, que foram adicionados pelo middleware.
		// Extrai o email e o role do contexto da requisição, que foram adicionados pelo middleware.
		email, ok := r.Context().Value("email").(string)
		if !ok {
			http.Error(w, "Email não encontrado no contexto", http.StatusInternalServerError)
			return
		}

		role, ok := r.Context().Value("role").(string)
		if !ok {
			http.Error(w, "Role não encontrado no contexto", http.StatusInternalServerError)
			return
		}

		// Retorna os dados do usuário como JSON.
		// Retorna os dados do usuário como JSON.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{Email: email, Role: role})
	}
}
