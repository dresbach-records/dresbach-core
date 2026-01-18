package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"hosting-backend/internal/models"

	"github.com/gorilla/mux"
)

// ActionResponse é uma estrutura padrão para respostas de ação.
type ActionResponse struct {
	Message string `json:"message"`
}

// SuspendServiceHandler suspende um serviço.
// Rota: PUT /admin/services/{id}/suspend
func SuspendServiceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID de serviço inválido", http.StatusBadRequest)
			return
		}

		err = models.UpdateServiceStatus(db, serviceID, models.ServiceStatusSuspended)
		if err != nil {
			// TODO: Adicionar um log mais detalhado aqui
			http.Error(w, "Erro ao suspender o serviço", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ActionResponse{Message: "Serviço suspenso com sucesso"})
	}
}

// ReactivateServiceHandler reativa um serviço.
// Rota: PUT /admin/services/{id}/reactivate
func ReactivateServiceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID de serviço inválido", http.StatusBadRequest)
			return
		}

		err = models.UpdateServiceStatus(db, serviceID, models.ServiceStatusActive)
		if err != nil {
			// TODO: Adicionar um log mais detalhado aqui
			http.Error(w, "Erro ao reativar o serviço", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ActionResponse{Message: "Serviço reativado com sucesso"})
	}
}
