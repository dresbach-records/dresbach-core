package webhooks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"hosting-backend/internal/email"
	"hosting-backend/internal/models"
	"hosting-backend/internal/provisioning"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

func StripeWebhookHandler(db *sql.DB) http.HandlerFunc {
    // ... (código existente)
}

func handleSuccessfulPayment(db *sql.DB, session *stripe.CheckoutSession) error {
	invoiceID, err := strconv.Atoi(session.ClientReferenceID)
	if err != nil {
		return fmt.Errorf("invalid client_reference_id (invoice_id): %s", session.ClientReferenceID)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on any error

	var serviceType string
	query := "UPDATE invoices SET status = ?, stripe_payment_intent_id = ? WHERE id = ? RETURNING service_type"
	err = tx.QueryRow(query, "paid", session.PaymentIntent.ID, invoiceID).Scan(&serviceType)
	if err != nil {
		return fmt.Errorf("failed to update invoice #%d: %w", invoiceID, err)
	}

	log.Printf("Fatura #%d paga. Tipo de serviço: %s", invoiceID, serviceType)

	// AQUI ESTÁ O ORQUESTRADOR
	switch serviceType {
	case "hosting", "cpanel":
		err = processCpanelOrder(tx, invoiceID)
	case "domain":
		err = processDomainOrder(tx, invoiceID)
	case "vps":
		err = processVpsOrder(tx, invoiceID)
	default:
		log.Printf("Tipo de serviço desconhecido ou não requer provisionamento: %s", serviceType)
	}

	if err != nil {
        // O erro já foi logado dentro da função de processo
		return err // Causa rollback
	}

	return tx.Commit()
}

// processCpanelOrder orquestra a criação de uma conta cPanel e registro de domínio.
func processCpanelOrder(tx *sql.Tx, invoiceID int) error {
    // 1. Buscar detalhes do pedido de cPanel (precisamos de um modelo para isso, ex: `CpanelOrder`)
    // Por enquanto, vamos simular os dados.
    log.Printf("Processando pedido de cPanel para a fatura #%d", invoiceID)
    
    // DADOS SIMULADOS (Em um caso real, viriam do banco de dados)
    order := struct{
        Domain string
        Plan string
        Username string
        Password string
        Email string
        NeedsDomainRegistration bool
    }{
        Domain: "meunovosite.com",
        Plan: "plan-essencial",
        Username: "meunovo",
        Password: "SenhaMuitoForte!123",
        Email: "cliente@exemplo.com",
        NeedsDomainRegistration: true,
    }

    // 2. Provisionar a conta no cPanel/WHM
    whm, err := provisioning.NewWhmProvisioner()
    if err != nil {
        return fmt.Errorf("falha ao inicializar o provisioner WHM: %w", err)
    }
    
    _, err = whm.CreateAccount(order.Username, order.Domain, order.Plan, order.Password, order.Email)
    if err != nil {
        return fmt.Errorf("falha ao criar conta cPanel para o domínio %s: %w", order.Domain, err)
    }

    // 3. Se for um domínio novo, registrá-lo e apontar DNS
    if order.NeedsDomainRegistration {
        hostinger, err := provisioning.NewHostingerProvisioner()
        if err != nil {
            return fmt.Errorf("falha ao inicializar o provisioner Hostinger: %w", err)
        }

        // TODO: Buscar ou criar um WhoisProfile
        whoisID := 12345 

        _, err = hostinger.RegisterDomain(order.Domain, 1, whoisID)
        if err != nil {
            return fmt.Errorf("falha ao registrar domínio %s: %w", order.Domain, err)
        }
        
        // TODO: Apontar o DNS do domínio para o IP do servidor WHM
        log.Printf("TODO: Apontar DNS para o domínio %s", order.Domain)
    }

    // 4. Atualizar o status do serviço no banco de dados para "active"
    log.Printf("SUCESSO: Conta cPanel para %s criada e domínio registrado.", order.Domain)
    return nil
}

func processDomainOrder(tx *sql.Tx, invoiceID int) error {
	// ... (código existente)
	return nil
}

func processVpsOrder(tx *sql.Tx, invoiceID int) error {
	// ... (código existente)
	return nil
}
