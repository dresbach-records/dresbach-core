package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"hosting-backend/internal/models"
	"hosting-backend/internal/whm" // Caminho corrigido
)

// CreateServicePayload é a estrutura de dados para criar uma nova conta de hospedagem.
type CreateServicePayload struct {
	Domain         string `json:"domain"`
	CpanelUser     string `json:"cpanel_user"`
	CpanelPassword string `json:"cpanel_password"`
	PlanName       string `json:"plan_name"` // Nome do plano no WHM
	UserID         int    `json:"user_id"`
	ProductID      int    `json:"product_id"`
}

// CreateServiceHandler cria uma nova conta de hospedagem no WHM e a registra como um serviço.
// Rota: POST /admin/services
func CreateServiceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload CreateServicePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		if payload.Domain == "" || payload.CpanelUser == "" || payload.CpanelPassword == "" || payload.PlanName == "" || payload.UserID == 0 || payload.ProductID == 0 {
			http.Error(w, "Campos inválidos. Domain, CpanelUser, CpanelPassword, PlanName, UserID e ProductID são obrigatórios.", http.StatusBadRequest)
			return
		}

		// Buscar o email do cliente usando o UserID
		user, err := models.GetUserByID(db, payload.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Usuário não encontrado", http.StatusNotFound)
				return
			}
			http.Error(w, "Erro ao buscar dados do usuário", http.StatusInternalServerError)
			return
		}

		whmClient, err := whm.NewClient()
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao inicializar o cliente WHM: %v", err), http.StatusInternalServerError)
			return
		}

		// Chamada corrigida para CreateAccount, incluindo o email do usuário
		err = whmClient.CreateAccount(payload.Domain, payload.CpanelUser, payload.CpanelPassword, user.Email, payload.PlanName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Falha ao criar conta no WHM: %v", err), http.StatusServiceUnavailable)
			return
		}

		service := models.Service{
			UserID:     payload.UserID,
			ProductID:  payload.ProductID,
			Domain:     payload.Domain,
			CpanelUser: payload.CpanelUser,
			Status:     models.ServiceStatusActive,
		}

		id, err := models.CreateService(db, &service)
		if err != nil {
			// Este é um estado crítico. A conta foi criada no servidor, mas não no nosso sistema.
			// A mensagem de erro deve ser clara sobre a necessidade de uma ação manual.
			http.Error(w, "Conta criada no WHM, mas falhou ao salvar localmente. Verifique os logs e registre o serviço manualmente.", http.StatusInternalServerError)
			return
		}
		service.ID = int(id)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(service)
	}
}
