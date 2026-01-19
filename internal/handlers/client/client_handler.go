package client

import (
	"encoding/json"
	"net/http"

	"hosting-backend/internal/services"
)

// GetServicesHandler lista os serviços do cliente logado.
func GetServicesHandler(clientService *services.ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			http.Error(w, "ID de usuário inválido", http.StatusUnauthorized)
			return
		}

		services, err := clientService.GetClientServices(userID)
		if err != nil {
			http.Error(w, "Erro ao listar serviços", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(services)
	}
}

// GetInvoicesHandler lista as faturas do cliente logado.
func GetInvoicesHandler(clientService *services.ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			http.Error(w, "ID de usuário inválido", http.StatusUnauthorized)
			return
		}

		invoices, err := clientService.GetClientInvoices(userID)
		if err != nil {
			http.Error(w, "Erro ao listar faturas", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invoices)
	}
}
