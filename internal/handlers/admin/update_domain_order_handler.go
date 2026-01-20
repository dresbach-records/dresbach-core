package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"hosting-backend/internal/orchestrator"

	"github.com/gorilla/mux"
)

// UpdateDomainOrderStatusRequest é a estrutura para o corpo da requisição de atualização.
type UpdateDomainOrderStatusRequest struct {
	Status string `json:"status"`
}

// UpdateDomainOrderHandler atualiza o status de um pedido de domínio.
func UpdateDomainOrderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID, ok := vars["id"]
		if !ok {
			http.Error(w, "ID do pedido não fornecido", http.StatusBadRequest)
			return
		}

		var req UpdateDomainOrderStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		// Validação do status
		allowedStatus := map[string]bool{
			"pending_payment":      true, // Cliente ainda não pagou
			"pending_registration": true, // Pago, aguardando nossa ação
			"completed":            true, // Processo finalizado com sucesso
			"failed":               true, // Falha em alguma etapa
			"cancelled":            true, // Cancelado pelo cliente ou admin
		}
		if !allowedStatus[req.Status] {
			http.Error(w, "Status inválido fornecido", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			log.Printf("Erro ao iniciar transação: %v", err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}

		// Atualiza o status no banco de dados
		res, err := tx.ExecContext(ctx, "UPDATE domain_orders SET status = $1 WHERE id = $2", req.Status, orderID)
		if err != nil {
			tx.Rollback()
			log.Printf("Erro ao atualizar status do pedido de domínio #%s: %v", orderID, err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			log.Printf("Erro ao verificar linhas afetadas para o pedido #%s: %v", orderID, err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			tx.Rollback()
			http.Error(w, "Pedido de domínio não encontrado", http.StatusNotFound)
			return
		}

		// Commit da transação antes de iniciar processos assíncronos
		if err := tx.Commit(); err != nil {
			log.Printf("Erro ao commitar transação para o pedido #%s: %v", orderID, err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}

		// Se o status for 'pending_registration', inicie o processo de provisionamento
		if req.Status == "pending_registration" {
			log.Printf("Disparando provisionamento para o pedido de domínio #%s", orderID)
			// A função de provisionamento deve ser chamada em uma goroutine
			// para não bloquear a resposta HTTP.
			go orchestrator.ProcessDomainProvisioning(db, orderID) // Passando o ID do pedido
		}

		log.Printf("[Admin] Status do pedido de domínio #%s atualizado para '%s'", orderID, req.Status)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Status do pedido de domínio atualizado com sucesso."})
	}
}
