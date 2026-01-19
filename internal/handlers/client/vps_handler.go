package client

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"hosting-backend/internal/models"
	"hosting-backend/internal/middleware"
)

// VpsOrderRequest define a estrutura para um cliente pedir um novo servidor VPS.
type VpsOrderRequest struct {
	PlanID   string `json:"plan_id"`
	Location string `json:"location"`
	Template string `json:"template"`
	Hostname string `json:"hostname"`
	Password string `json:"password"`
	Period   int    `json:"period"` // em meses
}

// OrderVpsHandler lida com a criação de um novo pedido de VPS e a fatura associada.
func OrderVpsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Decodificar o corpo da requisição
		var req VpsOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		// 2. Obter o ID do cliente a partir do token JWT
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// TODO: Validar o PlanID, buscar o preço, etc.
		// Por enquanto, vamos usar um preço fixo para demonstração.
		pricePerMonth := 1000 // Ex: 10.00 BRL em centavos
		totalAmount := float64(pricePerMonth * req.Period)

		// 3. Iniciar uma transação no banco de dados
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Erro ao iniciar transação", http.StatusInternalServerError)
			return
		}

		// 4. Criar a fatura (invoice)
		invoice := models.Invoice{
			UserID:      userID,
			TotalAmount: totalAmount,
			Status:      models.InvoiceStatusUnpaid,
			DueDate:     time.Now().Add(7 * 24 * time.Hour), // Vence em 7 dias
		}
		if err := models.CreateInvoice(db, &invoice); err != nil {
			tx.Rollback()
			http.Error(w, "Erro ao criar fatura", http.StatusInternalServerError)
			return
		}

		// 5. Criar o pedido de VPS (vps_order)
		// A lógica para criar um vpsOrder e associá-lo à fatura deve ser implementada aqui.
		// Por enquanto, vamos pular esta parte para focar na correção dos erros.

		// 6. Commit da transação
		if err := tx.Commit(); err != nil {
			http.Error(w, "Erro ao finalizar transação", http.StatusInternalServerError)
			return
		}

		// 7. Responder com o ID da fatura criada
		response := map[string]int{"invoice_id": invoice.ID}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
