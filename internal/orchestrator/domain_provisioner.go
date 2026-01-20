package orchestrator

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"hosting-backend/internal/models"
	"hosting-backend/internal/services/hostinger"
	"hosting-backend/internal/services/whm"
)

// DomainOrderData é uma estrutura para guardar os dados do pedido de domínio.
type DomainOrderData struct {
	ClientID   int
	DomainName string
	// O tipo (register/transfer) não está no pedido, então será fixo por enquanto.
	DomainType models.DomainType
	// O service ID não está disponível no pedido de domínio, precisará ser inferido ou adicionado posteriormente.
	ServiceID int
}

// ProcessDomainProvisioning orquestra o provisionamento de um domínio a partir de um pedido.
func ProcessDomainProvisioning(db *sql.DB, orderIDStr string) {
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		log.Printf("ERRO: orderID inválido: %s", orderIDStr)
		return
	}

	// 1. Iniciar transação e buscar dados do pedido
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Printf("Erro ao iniciar transação para o pedido #%d: %v", orderID, err)
		return
	}

	var orderData DomainOrderData
	// Simplificamos a query para refletir a estrutura real da tabela `domain_orders`
	query := `SELECT client_id, domain_name FROM domain_orders WHERE id = $1`
	err = tx.QueryRow(query, orderID).Scan(&orderData.ClientID, &orderData.DomainName)
	if err != nil {
		tx.Rollback()
		log.Printf("Erro ao buscar dados do pedido #%d: %v", orderID, err)
		return
	}

	// === Lógica de Negócio Provisória ===
	// Como o tipo não está no pedido, definimos como 'register'.
	orderData.DomainType = models.DomainTypeRegister
	// Como o service_id não está no pedido, usamos um valor provisório (ex: 1).
	// Em um sistema real, isso deveria ser resolvido (ex: um serviço padrão para domínios).
	orderData.ServiceID = 1 // TODO: Implementar uma forma de determinar o service_id correto

	// 2. Criar o registro do domínio
	domainID, err := models.CreateDomain(tx, orderData.ClientID, orderData.ServiceID, orderData.DomainName, orderData.DomainType)
	if err != nil {
		tx.Rollback()
		log.Printf("Erro ao criar registro de domínio para o pedido #%d: %v", orderID, err)
		return
	}

	// Commit da transação para salvar o novo domínio
	if err := tx.Commit(); err != nil {
		log.Printf("Erro ao commitar a criação do domínio para o pedido #%d: %v", orderID, err)
		finalizeAsFailed(db, domainID, "Erro crítico de banco de dados", err)
		return
	}

	log.Printf("Domínio #%d criado para o pedido #%d. Iniciando provisionamento.", domainID, orderID)
	// Inicia o processo de provisionamento real
	runProvisioningSteps(db, domainID, orderData)
}

// runProvisioningSteps contém a lógica de provisionamento.
func runProvisioningSteps(db *sql.DB, domainID int, orderData DomainOrderData) {
	hasStarted, err := models.HasEventOccurred(db, domainID, "provisioning.started")
	if err != nil || hasStarted {
		log.Printf("Provisionamento para o domínio %d já foi iniciado ou ocorreu um erro: %v", domainID, err)
		return
	}
	models.LogDomainEvent(db, domainID, "provisioning.started", "Orquestrador iniciou o processo de provisionamento.", nil)

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
	switch orderData.DomainType {
	case models.DomainTypeRegister:
		// Detalhes do cliente seriam buscados do banco de dados para um registro real
		providerOrderID, actionErr := hostingerClient.RegisterDomain(orderData.DomainName, nil) // Usando nil por enquanto
		if actionErr != nil {
			finalizeAsFailed(db, domainID, "Falha na API da Hostinger (Register)", actionErr)
			return
		}
		models.UpdateDomainProviderOrderID(db, domainID, "hostinger", providerOrderID)
		models.LogDomainEvent(db, domainID, "hostinger.register.initiated", "Pedido de registro criado na Hostinger", map[string]string{"order_id": providerOrderID})
	
	case models.DomainTypeTransfer:
        // A lógica de transferência permanece, caso seja usada no futuro
        authCode := "" // Deveria vir do pedido
        providerOrderID, actionErr := hostingerClient.TransferDomain(orderData.DomainName, authCode)
        if actionErr != nil {
            finalizeAsFailed(db, domainID, "Falha na API da Hostinger (Transfer)", actionErr)
            return
        }
        models.UpdateDomainProviderOrderID(db, domainID, "hostinger", providerOrderID)
        models.LogDomainEvent(db, domainID, "hostinger.transfer.initiated", "Pedido de transferência criado na Hostinger", map[string]string{"order_id": providerOrderID})
    }

	// ETAPA 2: CRIAÇÃO DA CONTA (WHM)
	cpanelUsername := "user" + strconv.Itoa(domainID)
	planName := "plano_default" // TODO: Deve vir do serviço associado
	tempPassword := "aVeryComplexP@ssw0rd!2024" // TODO: Gerar senha segura

	_, err = whmClient.CreateAccount(orderData.DomainName, cpanelUsername, tempPassword, planName)
	if err != nil {
		finalizeAsFailed(db, domainID, "Falha na criação da conta WHM", err)
		return
	}
	models.LogDomainEvent(db, domainID, "whm.account.created", "Conta cPanel criada com sucesso.", map[string]string{"username": cpanelUsername})

	// ETAPA 3: FINALIZAÇÃO
	models.UpdateDomainStatus(db, domainID, models.StatusActive)
	models.LogDomainEvent(db, domainID, "provisioning.completed", "Provisionamento concluído com sucesso.", nil)
	log.Printf("Provisionamento concluído para o domínio: %s (ID: %d)", orderData.DomainName, domainID)
}

// finalizeAsFailed padroniza o tratamento de falhas.
func finalizeAsFailed(db *sql.DB, domainID int, message string, err error) {
	log.Printf("ERRO no provisionamento do domínio %d: %s - %v", domainID, message, err)
	models.UpdateDomainStatus(db, domainID, models.StatusFailed)
	models.LogDomainEvent(db, domainID, "provisioning.failed", message, map[string]string{"error": err.Error()})
}
