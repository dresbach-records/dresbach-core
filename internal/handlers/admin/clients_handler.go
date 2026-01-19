package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"hosting-backend/internal/models"
	"hosting-backend/internal/services"
)

// ProvisionRequest é a estrutura para a solicitação de criação e provisionamento de cliente.
type ProvisionRequest struct {
	Client   models.Client `json:"client"`
	Domain   string        `json:"domain"`
	Username string        `json:"username"`
	Password string        `json:"password"`
	Plan     string        `json:"plan"`
}

// CreateClientAndProvisionHandler lida com a criação de um cliente e o provisionamento de sua conta.
func CreateClientAndProvisionHandler(adminService *services.AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ProvisionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Pedido inválido: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validação básica
		if req.Client.Email == "" || req.Domain == "" || req.Username == "" || req.Password == "" || req.Plan == "" {
			http.Error(w, "Campos obrigatórios ausentes no pedido", http.StatusBadRequest)
			return
		}

		clientID, err := adminService.CreateClientAndProvisionAccount(&req.Client, req.Domain, req.Username, req.Plan, req.Password)
		if err != nil {
			// O erro do serviço já é bem descritivo.
			http.Error(w, "Erro no processo de criação e provisionamento: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Retorna uma resposta de sucesso
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "Cliente criado e conta provisionada com sucesso.",
			"clientID": clientID,
		})
	}
}

// CreateClientHandler cria um novo cliente.
func CreateClientHandler(adminService *services.AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var client models.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			http.Error(w, "Pedido inválido", http.StatusBadRequest)
			return
		}

		id, err := adminService.CreateClient(&client)
		if err != nil {
			http.Error(w, "Erro ao criar cliente", http.StatusInternalServerError)
			return
		}

		client.ID = int(id)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(client)
	}
}

// GetClientsHandler lista todos os clientes.
func GetClientsHandler(adminService *services.AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clients, err := adminService.GetAllClients()
		if err != nil {
			http.Error(w, "Erro ao listar clientes", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	}
}

// GetClientHandler obtém um cliente específico.
func GetClientHandler(adminService *services.AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		client, err := adminService.GetClientByID(id)
		if err != nil {
			http.Error(w, "Erro ao buscar cliente", http.StatusInternalServerError)
			return
		}
		if client == nil {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(client)
	}
}

// UpdateClientHandler atualiza um cliente.
func UpdateClientHandler(adminService *services.AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		var client models.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			http.Error(w, "Pedido inválido", http.StatusBadRequest)
			return
		}

		client.ID = id
		if err := adminService.UpdateClient(&client); err != nil {
			http.Error(w, "Erro ao atualizar cliente", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// DeleteClientHandler deleta um cliente.
func DeleteClientHandler(adminService *services.AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		if err := adminService.DeleteClient(id); err != nil {
			http.Error(w, "Erro ao deletar cliente", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
