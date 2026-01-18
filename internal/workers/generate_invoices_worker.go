package workers

import (
	"database/sql"
	"log"
	"time"

	"hosting-backend/internal/models"
)

const (
	// Define o intervalo para o worker verificar por serviços a faturar.
	invoicingCheckInterval = 24 * time.Hour
	// Define quantos dias o cliente tem para pagar a fatura após a emissão.
	invoiceGraceDays = 7
)

// GenerateInvoicesWorker é um processo de fundo que gera faturas recorrentes para serviços ativos.
func GenerateInvoicesWorker(db *sql.DB) {
	log.Println("[Invoicing Worker] Iniciado. Verificando serviços a faturar a cada", invoicingCheckInterval)
	ticker := time.NewTicker(invoicingCheckInterval)
	defer ticker.Stop()

	// Loop infinito que executa a lógica do worker a cada tick do temporizador.
	for {
		select {
		case <-ticker.C:
			generateDueInvoices(db)
		}
	}
}

// generateDueInvoices busca por serviços que precisam de uma nova fatura e as cria.
func generateDueInvoices(db *sql.DB) {
	log.Println("[Invoicing Worker] Procurando por serviços que precisam de faturamento...")

	// 1. Buscar todos os serviços ativos cuja data de faturamento (next_due_date) já passou.
	servicesToInvoice, err := models.GetServicesDueForInvoicing(db)
	if err != nil {
		log.Printf("[Invoicing Worker] ERRO ao buscar serviços para faturar: %v", err)
		return
	}

	if len(servicesToInvoice) == 0 {
		log.Println("[Invoicing Worker] Nenhum serviço encontrado para faturamento no momento.")
		return
	}

	log.Printf("[Invoicing Worker] Encontrado(s) %d serviço(s) para gerar fatura.", len(servicesToInvoice))

	// 2. Itera sobre cada serviço e cria a fatura correspondente.
	for _, service := range servicesToInvoice {
		log.Printf("[Invoicing Worker] Gerando fatura para o serviço #%d (Domínio: %s)...", service.ID, service.Domain)

		// 3. Criar o objeto da nova fatura.
		now := time.Now()
		invoice := models.Invoice{
			UserID:      service.UserID,
			ServiceID:   service.ID,
			IssueDate:   now,
			DueDate:     now.AddDate(0, 0, invoiceGraceDays),
			TotalAmount: service.Price,
			Status:      models.InvoiceStatusUnpaid,
		}

		// 4. Salvar a nova fatura no banco de dados.
		if err := models.CreateInvoice(db, &invoice); err != nil {
			log.Printf("[Invoicing Worker] ERRO ao criar a fatura para o serviço #%d: %v", service.ID, err)
			continue // Pula para o próximo serviço
		}

		// 5. Calcular a próxima data de vencimento com base no ciclo de faturamento.
		var newNextDueDate time.Time
		if service.BillingCycle == models.BillingCycleMonthly {
			newNextDueDate = service.NextDueDate.AddDate(0, 1, 0) // Adiciona 1 mês
		} else if service.BillingCycle == models.BillingCycleAnnually {
			newNextDueDate = service.NextDueDate.AddDate(1, 0, 0) // Adiciona 1 ano
		} else {
			log.Printf("[Invoicing Worker] AVISO: Ciclo de faturamento desconhecido ('%s') para o serviço #%d. A data da fatura não será atualizada.", service.BillingCycle, service.ID)
			continue
		}

		// 6. Atualizar a 'next_due_date' do serviço para evitar faturas duplicadas.
		if err := models.UpdateNextDueDate(db, service.ID, newNextDueDate); err != nil {
			log.Printf("[Invoicing Worker] ERRO CRÍTICO: A fatura para o serviço #%d foi criada, mas falhou ao atualizar a next_due_date: %v", service.ID, err)
			// Este é um erro crítico porque pode levar a faturas duplicadas. Requer atenção.
			continue
		}

		log.Printf("[Invoicing Worker] Sucesso: Fatura gerada e next_due_date atualizada para %s para o serviço #%d.", newNextDueDate.Format("2006-01-02"), service.ID)
	}

	log.Println("[Invoicing Worker] Processo de geração de faturas concluído.")
}
