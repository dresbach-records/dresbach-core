package client

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

// CreateCheckoutSessionHandler cria uma sessão de checkout do Stripe para uma fatura.

// CreateCheckoutSessionHandler cria uma sessão de checkout do Stripe para uma fatura.
func CreateCheckoutSessionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		invoiceID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID da fatura inválido", http.StatusBadRequest)
			return
		}

		// Busca os detalhes da fatura no banco de dados
		var totalAmount float64
		var clientID int
		err = db.QueryRow("SELECT client_id, total_amount FROM invoices WHERE id = ? AND status = 'unpaid'", invoiceID).Scan(&clientID, &totalAmount)
		if err == sql.ErrNoRows {
			http.Error(w, "Fatura não encontrada ou já paga", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Erro ao buscar fatura", http.StatusInternalServerError)
			return
		}

		// Parâmetros para a sessão de checkout
		params := &stripe.CheckoutSessionParams{
			PaymentMethodTypes: stripe.StringSlice([]string{
				"card",
			}),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
						Currency:    stripe.String(string(stripe.CurrencyBRL)),
						ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
							Name: stripe.String(fmt.Sprintf("Fatura #%d", invoiceID)),
						},
						UnitAmount: stripe.Int64(int64(totalAmount * 100)), // O valor deve ser em centavos
					},
					Quantity: stripe.Int64(1),
				},
			},
			Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
			SuccessURL: stripe.String(os.Getenv("STRIPE_SUCCESS_URL")),
			CancelURL:  stripe.String(os.Getenv("STRIPE_CANCEL_URL")),
			Metadata: map[string]string{
				"invoice_id": strconv.Itoa(invoiceID),
			},
		}

		s, err := session.New(params)
		if err != nil {
			http.Error(w, "Erro ao criar sessão de checkout", http.StatusInternalServerError)
			return
		}

		// Retorna a URL da sessão de checkout para o cliente
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"url": s.URL})
	}
}
