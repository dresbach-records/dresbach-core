package client

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"log"

	"github.com/gorilla/mux"
	"hosting-backend/internal/services/asaas"
)

// CheckoutPayload define a estrutura de dados para o checkout.
type CheckoutPayload struct {
	InvoiceID int `json:"invoice_id"`
}

// CheckoutResponse define a estrutura de dados para a resposta do checkout.
type CheckoutResponse struct {
	PaymentID  string `json:"paymentId"`
	Status     string `json:"status"`
	InvoiceURL string `json:"invoiceUrl"`
}

// CheckoutHandler cria uma cobrança no Asaas para uma fatura.
func CheckoutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload CheckoutPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		invoiceID := payload.InvoiceID

		// Busca os detalhes da fatura no banco de dados
		var totalAmount float64
		var clientID int
		err := db.QueryRow("SELECT client_id, total_amount FROM invoices WHERE id = ? AND status = 'unpaid'", invoiceID).Scan(&clientID, &totalAmount)
		if err == sql.ErrNoRows {
			http.Error(w, "Fatura não encontrada ou já paga", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Erro ao buscar fatura", http.StatusInternalServerError)
			return
		}

		// Obter o asaas_customer_id do cliente
		var asaasCustomerID sql.NullString
		err = db.QueryRow("SELECT asaas_customer_id FROM clients WHERE id = ?", clientID).Scan(&asaasCustomerID)
		if err != nil {
			http.Error(w, "Erro ao buscar cliente", http.StatusInternalServerError)
			return
		}

		if !asaasCustomerID.Valid {
			http.Error(w, "Cliente não possui ID do Asaas", http.StatusBadRequest)
			return
		}

		// Criar a cobrança no Asaas
		asaasClient := asaas.NewAsaasClient()
		paymentRequest := asaas.PaymentRequest{
			Customer:    asaasCustomerID.String,
			BillingType: "BOLETO", // ou o tipo de cobrança desejado
			DueDate:     time.Now().AddDate(0, 0, 7).Format("2006-01-02"), // Vencimento em 7 dias
			Value:       totalAmount,
			Description: "Fatura #" + strconv.Itoa(invoiceID),
		}

		asaasPayment, err := asaasClient.CreatePayment(paymentRequest)
		if err != nil {
			log.Printf("Erro ao criar cobrança no Asaas: %v", err)
			http.Error(w, "Erro ao criar cobrança", http.StatusInternalServerError)
			return
		}

		// Salvar o ID de pagamento do Asaas na fatura
		_, err = db.Exec("UPDATE invoices SET asaas_payment_id = ? WHERE id = ?", asaasPayment.ID, invoiceID)
		if err != nil {
			log.Printf("Erro ao salvar ID de pagamento do Asaas: %v", err)
			// Não bloqueia a resposta para o cliente, mas loga o erro.
		}

		// Retorna a resposta do checkout para o cliente
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CheckoutResponse{
			PaymentID:  asaasPayment.ID,
			Status:     asaasPayment.Status,
			InvoiceURL: asaasPayment.InvoiceURL,
		})
	}
}
