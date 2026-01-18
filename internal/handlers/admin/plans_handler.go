package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"hosting-backend/internal/models"
	"github.com/go-chi/chi/v5"
)

// GetPlansHandler busca todos os planos de serviço.
// Rota: GET /admin/plans
func GetPlansHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		plans, err := models.GetAllPlans(db)
		if err != nil {
			http.Error(w, "Erro ao buscar planos", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plans)
	}
}

// CreatePlanHandler cria um novo plano de serviço.
// Rota: POST /admin/plans
func CreatePlanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var plan models.Plan
		if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		id, err := models.CreatePlan(db, &plan)
		if err != nil {
			http.Error(w, "Erro ao criar o plano", http.StatusInternalServerError)
			return
		}
		plan.ID = int(id)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(plan)
	}
}

// UpdatePlanHandler atualiza um plano existente.
// Rota: PUT /admin/plans/{id}
func UpdatePlanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		planID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "ID de plano inválido", http.StatusBadRequest)
			return
		}

		var plan models.Plan
		if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}
		plan.ID = planID

		if err := models.UpdatePlan(db, &plan); err != nil {
			http.Error(w, "Erro ao atualizar o plano", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// UpdatePlanStatusHandler atualiza o status de um plano (ativo, oculto, arquivado).
// Rota: PATCH /admin/plans/{id}/status
func UpdatePlanStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		planID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "ID de plano inválido", http.StatusBadRequest)
			return
		}

		var payload struct {
			Status models.PlanStatus `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		if err := models.UpdatePlanStatus(db, planID, payload.Status); err != nil {
			http.Error(w, "Erro ao atualizar o status do plano", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
