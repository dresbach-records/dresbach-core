package client

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"hosting-backend/internal/auth"
	"hosting-backend/internal/models"
)

// ListMyInvoicesHandler retorna a lista de faturas do cliente autenticado.
// @Summary Lista as faturas do cliente
// @Description Retorna um array com todas as faturas associadas ao cliente que faz a requisição.
// @Tags Client
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} models.Invoice "Lista de faturas"
// @Failure 401 {string} string "Erro ao identificar o usuário"
// @Failure 500 {string} string "Erro interno ao buscar faturas"
// @Router /api/my-invoices [get]
func ListMyInvoicesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extrai o ID do usuário do token JWT, que foi validado pelo middleware.
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "Erro ao identificar o usuário", http.StatusUnauthorized)
			return
		}

		// Busca as faturas do usuário no banco de dados.
		invoices, err := models.GetInvoicesByUserID(db, userID)
		if err != nil {
			// Log do erro para depuração no servidor.
			// log.Printf("Erro ao buscar faturas para o usuário %d: %v", userID, err)
			http.Error(w, "Erro interno ao buscar faturas", http.StatusInternalServerError)
			return
		}

		// Define o cabeçalho como JSON e envia a resposta.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invoices)
	}
}
