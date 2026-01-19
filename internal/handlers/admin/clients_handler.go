package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"hosting-backend/internal/models"
	"hosting-backend/internal/services"
	"github.com/gorilla/mux"
)

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
