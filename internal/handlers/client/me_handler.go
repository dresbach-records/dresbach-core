package client

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"hosting-backend/internal/models"
	"hosting-backend/internal/middleware"
)

// MeHandler busca e retorna as informações do usuário logado.
func MeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		user, err := models.GetUserByID(db, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Usuário não encontrado", http.StatusNotFound)
				return
			}
			http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
