package admin

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// DomainOrderView representa a estrutura de um pedido de domínio para o admin.
type DomainOrderView struct {
	ID          int       `json:"id"`
	ClientID    int       `json:"client_id"`
	ClientName  string    `json:"client_name"`
	DomainName  string    `json:"domain_name"`
	Document    string    `json:"document"`
	Status      string    `json:"status"`
	InvoiceID   int       `json:"invoice_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// GetDomainOrdersHandler lista todos os pedidos de registro de domínio.
func GetDomainOrdersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// A query pode ser filtrada por status, ex: /admin/domain-orders?status=pending_registration
		statusFilter := r.URL.Query().Get("status")

		query := `
			SELECT 
				do.id, do.client_id, c.name, do.domain_name, do.document, do.status, do.invoice_id, do.created_at
			FROM domain_orders do
			JOIN clients c ON do.client_id = c.id
		`
		args := []interface{}{}

		if statusFilter != "" {
			query += " WHERE do.status = ?"
			args = append(args, statusFilter)
		}

		query += " ORDER BY do.created_at DESC"

		rows, err := db.Query(query, args...)
		if err != nil {
			log.Printf("Erro ao consultar pedidos de domínio: %v", err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		orders := []DomainOrderView{}
		for rows.Next() {
			var order DomainOrderView
			var invoiceID sql.NullInt64 // invoice_id pode ser nulo

			if err := rows.Scan(&order.ID, &order.ClientID, &order.ClientName, &order.DomainName, &order.Document, &order.Status, &invoiceID, &order.CreatedAt); err != nil {
				log.Printf("Erro ao escanear pedido de domínio: %v", err)
				continue
			}
			if invoiceID.Valid {
				order.InvoiceID = int(invoiceID.Int64)
			}
			orders = append(orders, order)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	}
}
