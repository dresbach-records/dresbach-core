package client

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"hosting-backend/internal/models"
	"hosting-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// AddAllowedIPRequest define a estrutura para adicionar um novo IP.
type AddAllowedIPRequest struct {
	IPAddress   string `json:"ip_address"`
	Description string `json:"description"`
}

// AddAllowedIPHandler adiciona um novo IP à lista branca do cliente.
func AddAllowedIPHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, _ := middleware.GetClaims(r.Context())

		var req AddAllowedIPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Pedido inválido", http.StatusBadRequest)
			return
		}

		// (Validação adicional do formato do IP pode ser adicionada aqui)

		_, err := models.AddAllowedIP(db, claims.UserID, req.IPAddress, req.Description)
		if err != nil {
			// (Verificar erro de violação de chave única para mensagem mais amigável)
			http.Error(w, "Erro ao adicionar IP", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// GetAllowedIPsHandler lista todos os IPs da lista branca do cliente.
func GetAllowedIPsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, _ := middleware.GetClaims(r.Context())

		ips, err := models.GetAllowedIPsForClient(db, claims.UserID)
		if err != nil {
			http.Error(w, "Erro ao buscar IPs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ips)
	}
}

// DeleteAllowedIPHandler remove um IP da lista branca do cliente.
func DeleteAllowedIPHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, _ := middleware.GetClaims(r.Context())
		ipIDStr := chi.URLParam(r, "ipID")
		ipID, err := strconv.Atoi(ipIDStr)
		if err != nil {
			http.Error(w, "ID de IP inválido", http.StatusBadRequest)
			return
		}

		success, err := models.DeleteAllowedIP(db, claims.UserID, ipID)
		if err != nil {
			http.Error(w, "Erro ao deletar IP", http.StatusInternalServerError)
			return
		}

		if !success {
			http.Error(w, "IP não encontrado ou não pertence a você", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
