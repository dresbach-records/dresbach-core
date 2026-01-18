package webhooks

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"hosting-backend/internal/orchestrator"
    "hosting-backend/internal/models"
	"github.com/stripe/stripe-go/v74/webhook"
)

// StripeWebhookHandler processa os eventos recebidos do Stripe.
func StripeWebhookHandler(db *sql.DB) http.HandlerFunc {
	// A chave secreta do endpoint do webhook, obtida do seu dashboard do Stripe.
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	return func(w http.ResponseWriter, r *http.Request) {
		const MaxBodyBytes = int64(65536)
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Erro ao ler o corpo da requisição", http.StatusServiceUnavailable)
			return
		}

		// Validar a assinatura do webhook
		event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), endpointSecret)
		if err != nil {
			http.Error(w, "Falha na verificação da assinatura do webhook", http.StatusBadRequest)
			return
		}

		// Processar apenas os eventos de interesse
		if event.Type == "checkout.session.completed" {
			var session struct {
				Metadata map[string]string `json:"metadata"`
			}
			if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
				http.Error(w, "Erro ao decodificar a sessão do Stripe", http.StatusBadRequest)
				return
			}

			// Extrair o ID do domínio dos metadados
			domainIDStr, ok := session.Metadata["domain_id"]
			if !ok {
                // Logar o evento, mas não tratar como um erro fatal
				return
			}
            domainID, _ := strconv.Atoi(domainIDStr)

            // LOG DE AUDITORIA: Pagamento recebido
            models.LogDomainEvent(db, domainID, "payment.succeeded", "Webhook do Stripe recebido.", event.Data.Raw)

            // ATUALIZA STATUS DO DOMÍNIO
            models.UpdateDomainStatus(db, domainID, models.StatusPendingProvisioning)

			// *** INICIAR O ORQUESTRADOR ***
			go orchestrator.ProcessDomainProvisioning(db, domainID)
		}

		w.WriteHeader(http.StatusOK)
	}
}
