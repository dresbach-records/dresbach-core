package admin

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/balance"
	balancetransaction "github.com/stripe/stripe-go/v78/balancetransaction"
)

// --- Estruturas de Resposta ---

type BalanceResponse struct {
	AvailableAmount int64  `json:"available_amount"`
	PendingAmount   int64  `json:"pending_amount"`
	Currency        string `json:"currency"`
}

type Transaction struct {
	ID          string    `json:"id"`
	Amount      int64     `json:"amount"`      // Valor bruto
	Fee         int64     `json:"fee"`          // Taxa do Stripe
	Net         int64     `json:"net"`          // Valor líquido
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

// --- Handlers ---

func GetBalanceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := balance.Get(nil)
		if err != nil {
			log.Printf("Erro ao buscar saldo do Stripe: %v", err)
			http.Error(w, "Falha ao buscar saldo do Stripe", http.StatusInternalServerError)
			return
		}

		var availableAmount, pendingAmount int64
		var currency string
		if len(b.Available) > 0 {
			availableAmount = b.Available[0].Value
			currency = string(b.Available[0].Currency)
		}
		if len(b.Pending) > 0 {
			pendingAmount = b.Pending[0].Value
		}

		response := BalanceResponse{
			AvailableAmount: availableAmount,
			PendingAmount:   pendingAmount,
			Currency:        currency,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetTransactionsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := &stripe.BalanceTransactionListParams{}
		params.SetLimit("100") // Busca as últimas 100 transações
		i := balancetransaction.List(params)

		var transactions []Transaction
		for i.Next() {
			bt := i.BalanceTransaction()
			tx := Transaction{
				ID:          bt.ID,
				Amount:      bt.Amount,
				Fee:         bt.Fee,
				Net:         bt.Net,
				Currency:    string(bt.Currency),
				Description: bt.Description,
				CreatedAt:   time.Unix(bt.Created, 0),
			}
			transactions = append(transactions, tx)
		}

		if err := i.Err(); err != nil {
			log.Printf("Erro ao listar transações do Stripe: %v", err)
			http.Error(w, "Falha ao buscar transações do Stripe", http.StatusInternalServerError)
			return
		}

		response := TransactionsResponse{Transactions: transactions}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
