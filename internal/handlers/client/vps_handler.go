package client

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"hosting-backend/internal/models"
	"hosting-backend/internal/auth"

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
		claims, ok := auth.GetClaims(r.Context())
		if !ok {
			http.Error(w, "Token de autenticação inválido", http.StatusUnauthorized)
			return
		}
		clientID := claims.UserID

		// TODO: Validar o PlanID, buscar o preço, etc.
		// Por enquanto, vamos usar um preço fixo para demonstração.
		pricePerMonth := 1000 // Ex: 10.00 BRL em centavos
		totalAmount := pricePerMonth * req.Period

		// 3. Iniciar uma transação no banco de dados
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Erro ao iniciar transação", http.StatusInternalServerError)
			return
		}

		// 4. Criar a fatura (invoice)
		invoice := models.Invoice{
			ClientID:  clientID,
			Amount:    totalAmount,
			Status:    "unpaid",
			DueDate:   time.Now().Add(7 * 24 * time.Hour), // Vence em 7 dias
			ServiceType: "vps",
		}
		invoiceID, err := invoice.Create(tx)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Erro ao criar fatura", http.StatusInternalServerError)
			return
		}

		// 5. Criar o pedido de VPS (vps_order)
		vpsOrder := models.VpsOrder{
			ClientID:  clientID,
			InvoiceID: invoiceID,
			PlanID:    req.PlanID,
			Location:  req.Location,
			Template:  req.Template,
			Hostname:  req.Hostname,
			Password:  req.Password, // ATENÇÃO: Em produção, isso deve ser hasheado!
			Status:    "pending",
		}
		if _, err := vpsOrder.Create(tx); err != nil {
			tx.Rollback()
			http.Error(w, "Erro ao criar pedido de VPS", http.StatusInternalServerError)
			return
		}

		// 6. Commit da transação
		if err := tx.Commit(); err != nil {
			http.Error(w, "Erro ao finalizar transação", http.StatusInternalServerError)
			return
		}

		// 7. Responder com o ID da fatura criada
		response := map[string]int64{"invoice_id": invoiceID}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
