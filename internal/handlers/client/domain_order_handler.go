package client

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"hosting-backend/internal/middleware"
)

// OrderDomainRequest representa a requisição para registrar um domínio.
type OrderDomainRequest struct {
	Domain   string `json:"domain"`
	Document string `json:"document"` // CPF ou CNPJ
}

// OrderDomainResponse representa a resposta após um pedido de domínio bem-sucedido.
type OrderDomainResponse struct {
	Message   string `json:"message"`
	InvoiceID int    `json:"invoice_id"`
	Domain    string `json:"domain"`
}

const domainRegistrationPrice = 40.00

// OrderDomainHandler processa um novo pedido de registro de domínio.
func OrderDomainHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Decodifica o corpo da requisição
		var req OrderDomainRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Requisição inválida", http.StatusBadRequest)
			return
		}

		domain := strings.ToLower(strings.TrimSpace(req.Domain))

		// 2. Extrai o ID do cliente do token JWT
		claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.Claims)
		if !ok {
			http.Error(w, "Falha ao obter dados do cliente", http.StatusUnauthorized)
			return
		}

		// 3. Verifica novamente a disponibilidade do domínio
		if !isDomainAvailable(domain) {
			http.Error(w, fmt.Sprintf("O domínio %s não está mais disponível para registro.", domain), http.StatusConflict) // 409 Conflict
			return
		}

		// 4. Inicia uma transação no banco de dados
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Erro ao iniciar transação: %v", err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback() // Rollback em caso de erro

		// 5. Insere o pedido de domínio
		orderRes, err := tx.Exec("INSERT INTO domain_orders (client_id, domain_name, document, status) VALUES (?, ?, ?, ?)",
			claims.UserID, domain, req.Document, "pending_payment")
		if err != nil {
			log.Printf("Erro ao inserir pedido de domínio: %v", err)
			http.Error(w, "Erro ao processar seu pedido.", http.StatusInternalServerError)
			return
		}
		orderID, _ := orderRes.LastInsertId()

		// 6. Cria a fatura para o pedido
		now := time.Now()
		dueDate := now.Add(7 * 24 * time.Hour) // Vencimento em 7 dias
		invoiceRes, err := tx.Exec("INSERT INTO invoices (client_id, issue_date, due_date, total_amount, status) VALUES (?, ?, ?, ?, ?)",
			claims.UserID, now, dueDate, domainRegistrationPrice, "unpaid")
		if err != nil {
			log.Printf("Erro ao criar fatura: %v", err)
			http.Error(w, "Erro ao processar seu pedido.", http.StatusInternalServerError)
			return
		}
		invoiceID, _ := invoiceRes.LastInsertId()

		// 7. Adiciona o item à fatura (associando ao pedido de domínio)
		// Usaremos um service_id nulo para itens que não são serviços recorrentes, como um registro de domínio.
		description := fmt.Sprintf("Registro de domínio: %s - 1 ano", domain)
		_, err = tx.Exec("INSERT INTO invoice_items (invoice_id, service_id, description, amount) VALUES (?, NULL, ?, ?)",
			invoiceID, description, domainRegistrationPrice)
		if err != nil {
			log.Printf("Erro ao criar item de fatura: %v", err)
			http.Error(w, "Erro ao processar seu pedido.", http.StatusInternalServerError)
			return
		}

		// 8. Associa o pedido de domínio à fatura (para referência futura)
		_, err = tx.Exec("UPDATE domain_orders SET invoice_id = ? WHERE id = ?", invoiceID, orderID)
		if err != nil {
			log.Printf("Erro ao associar fatura ao pedido: %v", err)
			http.Error(w, "Erro ao processar seu pedido.", http.StatusInternalServerError)
			return
		}

		// Se tudo deu certo, comita a transação
		if err := tx.Commit(); err != nil {
			log.Printf("Erro ao comitar transação: %v", err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}

		// 9. Gera a notificação para o admin (neste ponto, apenas um log)
		log.Printf("### [ADMIN] NOVA PRIORIDADE: REGISTRO DE DOMÍNIO MANUAL ###")
		log.Printf("### Cliente: %d | Domínio: %s | Documento: %s | Fatura: %d ###", claims.UserID, domain, req.Document, invoiceID)

		// 10. Retorna sucesso para o cliente
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(OrderDomainResponse{
			Message:   fmt.Sprintf("Pedido para o domínio %s criado com sucesso! Pague a fatura para iniciar o processo de registro.", domain),
			InvoiceID: int(invoiceID),
			Domain:    domain,
		})
	}
}

// isDomainAvailable é uma função auxiliar para verificar a disponibilidade no RDAP.
func isDomainAvailable(domain string) bool {
	rdapURL := fmt.Sprintf("https://rdap.registro.br/domain/%s", domain)
	resp, err := http.Get(rdapURL)
	if err != nil {
		log.Printf("Erro ao consultar RDAP para %s: %v", domain, err)
		return false // Por segurança, se a API falhar, não permite o registro
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusNotFound
}
