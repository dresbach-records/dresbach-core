package client

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"hosting-backend/internal/middleware"
	"hosting-backend/internal/models"
)

// ListMyInvoicesHandler retorna a lista de faturas do cliente autenticado.
// @Summary Lista as faturas do cliente
// @Description Retorna um array com todas as faturas associadas ao cliente autenticado.
// @Tags Client
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Invoice "Lista de faturas"
// @Failure 401 {string} string "Erro ao identificar o usuário"
// @Failure 500 {string} string "Erro interno ao buscar faturas"
// @Router /api/my-invoices [get]
func ListMyInvoicesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extrair o userID do contexto (injetado pelo middleware JWT)
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Buscar faturas do usuário
		invoices, err := models.GetInvoicesByUserID(db, userID)
		if err != nil {
			http.Error(w, "Erro interno ao buscar faturas", http.StatusInternalServerError)
			return
		}

		// Responder JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invoices)
	}
}
