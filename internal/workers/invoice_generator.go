package workers

import (
	"database/sql"
	"log"
	"time"

	"hosting-backend/internal/models"
)

// GenerateInvoicesWorker é o worker que gera as faturas.
func GenerateInvoicesWorker(db *sql.DB) {
	ticker := time.NewTicker(24 * time.Hour) // Executa a cada 24 horas
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Executando worker de geração de faturas...")
		generateInvoices(db)
	}
}

func generateInvoices(db *sql.DB) {
	// Busca serviços que precisam de renovação
	rows, err := db.Query(`SELECT id, client_id, name, billing_cycle, price
		FROM services WHERE status = 'active' AND next_due_date <= NOW()`)
	if err != nil {
		log.Printf("Erro ao buscar serviços para faturamento: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var service models.Service
		if err := rows.Scan(&service.ID, &service.ClientID, &service.Name, &service.BillingCycle, &service.Price); err != nil {
			log.Printf("Erro ao ler dados do serviço: %v", err)
			continue
		}

		// Inicia uma transação para garantir a consistência dos dados
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Erro ao iniciar transação: %v", err)
			continue
		}

		if err := createInvoiceAndItems(tx, service); err != nil {
			tx.Rollback()
			log.Printf("Erro ao criar fatura e itens: %v", err)
			continue
		}

		if err := updateNextDueDate(tx, service); err != nil {
			tx.Rollback()
			log.Printf("Erro ao atualizar data de vencimento: %v", err)
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Erro ao commitar transação: %v", err)
		}

		log.Printf("Fatura gerada com sucesso para o serviço #%d", service.ID)
	}
}

func createInvoiceAndItems(tx *sql.Tx, service models.Service) error {
	issueDate := time.Now()
	dueDate := issueDate.AddDate(0, 0, 15) // Vencimento em 15 dias

	// Cria a fatura
	result, err := tx.Exec(`INSERT INTO invoices (client_id, issue_date, due_date, total_amount, status)
		VALUES (?, ?, ?, ?, 'unpaid')`,
		service.ClientID, issueDate, dueDate, service.Price)
	if err != nil {
		return err
	}
	invoiceID, _ := result.LastInsertId()

	// Adiciona o item na fatura
	_, err = tx.Exec(`INSERT INTO invoice_items (invoice_id, service_id, description, amount)
		VALUES (?, ?, ?, ?)`,
		invoiceID, service.ID, "Renovação de "+service.Name, service.Price)

	return err
}

func updateNextDueDate(tx *sql.Tx, service models.Service) error {
	nextDueDate := calculateNextDueDate(service.BillingCycle)
	_, err := tx.Exec("UPDATE services SET next_due_date = ? WHERE id = ?", nextDueDate, service.ID)
	return err
}

func calculateNextDueDate(billingCycle string) time.Time {
	now := time.Now()
	switch billingCycle {
	case "monthly":
		return now.AddDate(0, 1, 0)
	case "quarterly":
		return now.AddDate(0, 3, 0)
	case "semi-annually":
		return now.AddDate(0, 6, 0)
	case "annually":
		return now.AddDate(1, 0, 0)
	default:
		return now.AddDate(0, 1, 0) // Padrão é mensal
	}
}
