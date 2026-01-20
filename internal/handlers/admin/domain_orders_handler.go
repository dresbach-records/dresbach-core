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
				do.id,
				do.client_id,
				c.name AS client_name,
				do.domain_name,
				c.document,
				do.status,
				do.invoice_id,
				do.created_at
			FROM domain_orders do
			JOIN clients c ON do.client_id = c.id`

		var params []interface{}
		if statusFilter != "" {
			query += " WHERE do.status = $1"
			params = append(params, statusFilter)
		}

		query += " ORDER BY do.created_at DESC"

		rows, err := db.Query(query, params...)
		if err != nil {
			log.Printf("Error querying domain orders: %v", err)
			http.Error(w, "Erro ao buscar pedidos de domínio", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var orders []DomainOrderView
		for rows.Next() {
			var order DomainOrderView
			if err := rows.Scan(
				&order.ID,
				&order.ClientID,
				&order.ClientName,
				&order.DomainName,
				&order.Document,
				&order.Status,
				&order.InvoiceID,
				&order.CreatedAt,
			); err != nil {
				log.Printf("Error scanning domain order: %v", err)
				http.Error(w, "Erro ao processar pedidos de domínio", http.StatusInternalServerError)
				return
			}
			orders = append(orders, order)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Error iterating domain orders: %v", err)
			http.Error(w, "Erro ao processar lista de pedidos de domínio", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	}
}
