package admin

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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

		// Validação simples do status
		allowedStatus := map[string]bool{
			"pending_registration": true,
			"completed":            true,
			"failed":               true,
		}
		if !allowedStatus[req.Status] {
			http.Error(w, "Status inválido fornecido", http.StatusBadRequest)
			return
		}

		// Atualiza o status no banco de dados
		res, err := db.Exec("UPDATE domain_orders SET status = ? WHERE id = ?", req.Status, orderID)
		if err != nil {
			log.Printf("Erro ao atualizar status do pedido de domínio #%s: %v", orderID, err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Printf("Erro ao verificar linhas afetadas para o pedido #%s: %v", orderID, err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Pedido de domínio não encontrado", http.StatusNotFound)
			return
		}

		log.Printf("[Admin] Status do pedido de domínio #%s atualizado para '%s'", orderID, req.Status)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Status do pedido de domínio atualizado com sucesso."}) 
	}
}
