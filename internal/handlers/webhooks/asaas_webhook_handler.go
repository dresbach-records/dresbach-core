package webhooks

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"log"
)

// WebhookEvent define a estrutura básica de um evento de webhook do Asaas.
type WebhookEvent struct {
	Event   string `json:"event"`
	Payment struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"payment"`
}

// AsaasWebhookHandler processa os webhooks do Asaas.
func AsaasWebhookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event WebhookEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			log.Printf("[Webhook Asaas] Erro ao decodificar o corpo da requisição: %v", err)
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		log.Printf("[Webhook Asaas] Evento recebido: %s, ID do Pagamento: %s, Status: %s", event.Event, event.Payment.ID, event.Payment.Status)

		// Atualizar o status da fatura no banco de dados
		var newStatus string
		switch event.Event {
		case "PAYMENT_CONFIRMED", "PAYMENT_RECEIVED":
			newStatus = "paid"
		case "PAYMENT_OVERDUE":
			newStatus = "overdue"
		case "PAYMENT_CANCELED":
			newStatus = "canceled"
		default:
			log.Printf("[Webhook Asaas] Evento não tratado: %s", event.Event)
			w.WriteHeader(http.StatusOK) // Responde OK para eventos não tratados
			return
		}

		_, err := db.Exec("UPDATE invoices SET status = ? WHERE asaas_payment_id = ?", newStatus, event.Payment.ID)
		if err != nil {
			log.Printf("[Webhook Asaas] Erro ao atualizar o status da fatura: %v", err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}

		log.Printf("[Webhook Asaas] Fatura com asaas_payment_id %s atualizada para o status '%s'", event.Payment.ID, newStatus)
		w.WriteHeader(http.StatusOK)
	}
}
