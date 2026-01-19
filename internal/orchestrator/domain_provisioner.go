package orchestrator

import (
	"database/sql"
	"fmt"
	"log"

	"hosting-backend/internal/models"
	"hosting-backend/internal/services/hostinger"
	"hosting-backend/internal/services/whm"
)

// ProcessDomainProvisioning é o maestro do fluxo de provisionamento.
func ProcessDomainProvisioning(db *sql.DB, domainID int) {
	hasStarted, err := models.HasEventOccurred(db, domainID, "provisioning.started")
	if err != nil || hasStarted {
		log.Printf("Provisionamento para o domínio %d já iniciado ou erro ao verificar: %v", domainID, err)
		return
	}
	models.LogDomainEvent(db, domainID, "provisioning.started", "Orquestrador iniciou o processo de provisionamento.", nil)

	// BUSCAR DADOS (exemplo simplificado)
	var domain models.Domain
	// Em uma implementação real, você buscaria os detalhes completos do domínio, cliente e serviço.

	// INICIAR CLIENTES DE SERVIÇO
	hostingerClient, err := hostinger.NewClient()
	if err != nil {
		finalizeAsFailed(db, domainID, "Erro ao iniciar cliente Hostinger", err)
		return
	}
	whmClient, err := whm.NewClient()
	if err != nil {
		finalizeAsFailed(db, domainID, "Erro ao iniciar cliente WHM", err)
		return
	}

	// ETAPA 1: AÇÃO DE DOMÍNIO (HOSTINGER)
	switch domain.Type {
	case models.DomainTypeRegister, models.DomainTypeTransfer:
		var orderID string
		var actionErr error
		if domain.Type == models.DomainTypeRegister {
			orderID, actionErr = hostingerClient.RegisterDomain(domain.DomainName, nil) // Detalhes do cliente omitidos
		} else {
			authCode := "" // Auth-code viria do pedido do cliente
			orderID, actionErr = hostingerClient.TransferDomain(domain.DomainName, authCode)
		}

		if actionErr != nil {
			finalizeAsFailed(db, domainID, fmt.Sprintf("Falha na API da Hostinger (%s)", domain.Type), actionErr)
			return
		}
		models.LogDomainEvent(db, domainID, fmt.Sprintf("hostinger.%s.initiated", domain.Type), "Pedido criado com sucesso na Hostinger", map[string]string{"order_id": orderID})
	case models.DomainTypeExisting:
		models.LogDomainEvent(db, domainID, "dns.setup.skipped", "Domínio existente, pulando etapa de registro/transferência.", nil)
	}

	// ETAPA 2: CRIAÇÃO DA CONTA (WHM)
	// O nome de usuário do cPanel precisa ser único, uma lógica para gerá-lo seria necessária
	cpanelUsername := "user" + fmt.Sprintf("%d", domainID)
	planName := "plano_default" // O nome do plano viria dos detalhes do serviço
	// Em um cenário real, a senha deveria ser gerada de forma segura
	tempPassword := "aVeryComplexP@ssw0rd!"

	_, err = whmClient.CreateAccount(domain.DomainName, cpanelUsername, tempPassword, planName)
	if err != nil {
		finalizeAsFailed(db, domainID, "Falha na criação da conta WHM", err)
		// *** PONTO CRÍTICO DE ROLLBACK LÓGICO ***
		// Aqui, você adicionaria uma lógica para notificar a equipe ou tentar compensar a ação.
		// Ex: Enviar um email para o suporte: "Domínio 'X' registrado mas falha no WHM."
		return
	}
	models.LogDomainEvent(db, domainID, "whm.account.created", "Conta cPanel criada com sucesso.", map[string]string{"username": cpanelUsername})

	// ETAPA 3: FINALIZAÇÃO
	models.UpdateDomainStatus(db, domainID, models.StatusActive)
	models.LogDomainEvent(db, domainID, "provisioning.completed", "Provisionamento concluído com sucesso.", nil)
	log.Printf("Provisionamento concluído para o domínio: %s (ID: %d)", domain.DomainName, domain.ID)
}

// finalizeAsFailed é uma função helper para padronizar o tratamento de falhas.
func finalizeAsFailed(db *sql.DB, domainID int, message string, err error) {
	log.Printf("ERRO no provisionamento do domínio %d: %s - %v", domainID, message, err)
	models.UpdateDomainStatus(db, domainID, models.StatusFailed)
	models.LogDomainEvent(db, domainID, "provisioning.failed", message, map[string]string{"error": err.Error()})
}
