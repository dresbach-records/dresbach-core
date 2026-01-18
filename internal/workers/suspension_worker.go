package workers

import (
	"database/sql"
	"log"
	"time"

	"hosting-backend/internal/models"
)

const (
	// Define o intervalo para o worker verificar por contas vencidas.
	suspensionCheckInterval = 24 * time.Hour
)

// SuspensionWorker é um processo de fundo que verifica e suspende contas com pagamentos atrasados.
func SuspensionWorker(db *sql.DB) {
	log.Println("[Suspension Worker] Iniciado. Verificando faturas vencidas a cada", suspensionCheckInterval)
	ticker := time.NewTicker(suspensionCheckInterval)
	defer ticker.Stop()

	// Loop infinito que executa a lógica do worker a cada tick do temporizador.
	for {
		select {
		case <-ticker.C:
			checkForOverdueServices(db)
		}
	}
}

// checkForOverdueServices busca por faturas vencidas e suspende os serviços associados.
func checkForOverdueServices(db *sql.DB) {
	log.Println("[Suspension Worker] Procurando por faturas vencidas para suspender serviços...")

	// 1. Buscar todas as faturas com status 'unpaid' e data de vencimento passada.
	overdueInvoices, err := models.GetOverdueInvoices(db)
	if err != nil {
		log.Printf("[Suspension Worker] ERRO ao buscar faturas vencidas: %v", err)
		return // Retorna para não continuar se houver um erro no banco de dados.
	}

	if len(overdueInvoices) == 0 {
		log.Println("[Suspension Worker] Nenhuma fatura vencida encontrada.")
		return
	}

	log.Printf("[Suspension Worker] Encontrado(s) %d serviço(s) para suspender baseado em faturas vencidas.", len(overdueInvoices))

	// 2. Itera sobre cada fatura vencida e suspende o serviço correspondente.
	for _, invoice := range overdueInvoices {
		log.Printf("[Suspension Worker] Suspendendo serviço #%d devido à fatura #%d vencida...", invoice.ServiceID, invoice.ID)

		// TODO: Aqui entraria a lógica de integração com um painel de controle real.
		// Ex: err := cpanel.SuspendAccount(service.CpanelUser)
		// Por enquanto, vamos apenas atualizar nosso banco de dados.

		// 3. Atualiza o status do serviço no banco de dados para 'suspended'.
		err := models.UpdateServiceStatus(db, invoice.ServiceID, models.ServiceStatusSuspended)
		if err != nil {
			log.Printf("[Suspension Worker] ERRO ao tentar suspender o serviço #%d no banco de dados: %v", invoice.ServiceID, err)
			// Continua para o próximo serviço, mesmo que este tenha falhado.
			continue
		}

		log.Printf("[Suspension Worker] Sucesso: Serviço #%d foi suspenso.", invoice.ServiceID)
	}

	log.Println("[Suspension Worker] Processo de suspensão concluído.")
}
