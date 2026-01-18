package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"hosting-backend/internal/models"

	"github.com/gorilla/mux"
)

// CreateClientHandler cria um novo cliente.
func CreateClientHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var client models.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			http.Error(w, "Pedido inválido", http.StatusBadRequest)
			return
		}

		// Inserir no banco de dados
		result, err := db.Exec(`INSERT INTO clients (user_id, company_name, contact_name, email, phone, address, city, state, zip, country)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			client.UserID, client.CompanyName, client.ContactName, client.Email, client.Phone, client.Address, client.City, client.State, client.Zip, client.Country)
		if err != nil {
			http.Error(w, "Erro ao criar cliente", http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()
		client.ID = int(id)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(client)
	}
}

// GetClientsHandler lista todos os clientes.
func GetClientsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, user_id, company_name, contact_name, email, phone, address, city, state, zip, country, created_at FROM clients")
		if err != nil {
			http.Error(w, "Erro ao listar clientes", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var clients []models.Client
		for rows.Next() {
			var client models.Client
			if err := rows.Scan(&client.ID, &client.UserID, &client.CompanyName, &client.ContactName, &client.Email, &client.Phone, &client.Address, &client.City, &client.State, &client.Zip, &client.Country, &client.CreatedAt); err != nil {
				http.Error(w, "Erro ao ler dados do cliente", http.StatusInternalServerError)
				return
			}
			clients = append(clients, client)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	}
}

// GetClientHandler obtém um cliente específico.
func GetClientHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		var client models.Client
		err = db.QueryRow("SELECT id, user_id, company_name, contact_name, email, phone, address, city, state, zip, country, created_at FROM clients WHERE id = ?", id).
			Scan(&client.ID, &client.UserID, &client.CompanyName, &client.ContactName, &client.Email, &client.Phone, &client.Address, &client.City, &client.State, &client.Zip, &client.Country, &client.CreatedAt)
		if err == sql.ErrNoRows {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Erro ao buscar cliente", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(client)
	}
}

// UpdateClientHandler atualiza um cliente.
func UpdateClientHandler(db *sql.DB) http.HandlerFunc {
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

		_, err = db.Exec(`UPDATE clients SET company_name = ?, contact_name = ?, email = ?, phone = ?, address = ?, city = ?, state = ?, zip = ?, country = ?
			WHERE id = ?`,
			client.CompanyName, client.ContactName, client.Email, client.Phone, client.Address, client.City, client.State, client.Zip, client.Country, id)
		if err != nil {
			http.Error(w, "Erro ao atualizar cliente", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// DeleteClientHandler deleta um cliente.
func DeleteClientHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM clients WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Erro ao deletar cliente", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
