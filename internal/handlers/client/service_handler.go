package client

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"hosting-backend/internal/models"

	"github.com/go-chi/chi/v5"
)

// ListMyServicesHandler lista todos os serviços associados ao cliente autenticado.
// Rota: GET /api/my-services
func ListMyServicesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementar a obtenção do UserID a partir do contexto da requisição
		userID := 1 // Placeholder

		services, err := models.GetServicesByUserID(db, userID)
		if err != nil {
			http.Error(w, "Erro ao buscar os serviços", http.StatusInternalServerError)
			return
		}

		if services == nil {
			services = []models.Service{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(services)
	}
}

// GetServiceDetailsHandler busca os detalhes de um serviço específico do cliente.
// Rota: GET /api/my-services/{id}
func GetServiceDetailsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Obter o ID do usuário autenticado a partir do contexto.
		// TODO: Implementar a obtenção do UserID a partir do contexto da requisição
		userID := 1 // Placeholder

		// 2. Obter o ID do serviço a partir dos parâmetros da URL.
		serviceIDStr := chi.URLParam(r, "id")
		serviceID, err := strconv.Atoi(serviceIDStr)
		if err != nil {
			http.Error(w, "ID de serviço inválido", http.StatusBadRequest)
			return
		}

		// 3. Buscar o serviço no banco de dados, garantindo que ele pertença ao usuário.
		service, err := models.GetServiceByID(db, serviceID, userID)
		if err != nil {
			// Erro interno do servidor ao consultar o banco de dados.
			http.Error(w, "Erro ao buscar o serviço", http.StatusInternalServerError)
			return
		}

		// 4. Se o serviço for nil, significa que não foi encontrado ou não pertence ao usuário.
		if service == nil {
			http.Error(w, "Serviço não encontrado", http.StatusNotFound)
			return
		}

		// 5. Retornar os detalhes do serviço.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(service)
	}
}
