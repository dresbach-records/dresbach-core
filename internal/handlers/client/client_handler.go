package client

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"hosting-backend/internal/models"
)

// GetServicesHandler lista os serviços do cliente logado.

// GetServicesHandler lista os serviços do cliente logado.
func GetServicesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, _ := r.Context().Value("email").(string)

		// Busca o ID do cliente com base no email
		var clientID int
		err := db.QueryRow("SELECT id FROM clients WHERE email = ?", email).Scan(&clientID)
		if err != nil {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		}

		rows, err := db.Query("SELECT id, name, description, billing_cycle, price, status, next_due_date FROM services WHERE client_id = ?", clientID)
		if err != nil {
			http.Error(w, "Erro ao listar serviços", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var services []models.Service
		for rows.Next() {
			var service models.Service
			if err := rows.Scan(&service.ID, &service.Name, &service.Description, &service.BillingCycle, &service.Price, &service.Status, &service.NextDueDate); err != nil {
				http.Error(w, "Erro ao ler dados do serviço", http.StatusInternalServerError)
				return
			}
			services = append(services, service)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(services)
	}
}

// GetInvoicesHandler lista as faturas do cliente logado.

// GetInvoicesHandler lista as faturas do cliente logado.
func GetInvoicesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, _ := r.Context().Value("email").(string)

		// Busca o ID do cliente com base no email
		var clientID int
		err := db.QueryRow("SELECT id FROM clients WHERE email = ?", email).Scan(&clientID)
		if err != nil {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		}

		rows, err := db.Query("SELECT id, issue_date, due_date, total_amount, status FROM invoices WHERE client_id = ?", clientID)
		if err != nil {
			http.Error(w, "Erro ao listar faturas", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var invoices []models.Invoice
		for rows.Next() {
			var invoice models.Invoice
			if err := rows.Scan(&invoice.ID, &invoice.IssueDate, &invoice.DueDate, &invoice.TotalAmount, &invoice.Status); err != nil {
				http.Error(w, "Erro ao ler dados da fatura", http.StatusInternalServerError)
				return
			}
			invoices = append(invoices, invoice)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invoices)
	}
}
