package client

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"hosting-backend/internal/models"
	"hosting-backend/internal/middleware"
)

// GetLoginHistoryHandler retorna o histórico de login recente para o cliente autenticado.
func GetLoginHistoryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obter claims do contexto, que foi adicionado pelo middleware de autenticação
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Token inválido ou não encontrado", http.StatusUnauthorized)
			return
		}

		// Definir um limite padrão, mas permitir que seja sobrescrito por query param
		limitStr := r.URL.Query().Get("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 20 // Limite padrão
		}

		// Buscar o histórico usando o UserID do token
		history, err := models.GetLoginHistoryForClient(db, userID, limit)
		if err != nil {
			http.Error(w, "Erro ao buscar histórico de login", http.StatusInternalServerError)
			return
		}

		// Se o histórico for nulo (nenhum registro), retorne um array vazio
		if history == nil {
			history = []models.LoginHistoryEntry{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(history)
	}
}
