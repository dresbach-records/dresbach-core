package models

import (
	"database/sql"
	"time"
)

// DomainType define os tipos de operação de domínio.
type DomainType string

const (
	DomainTypeRegister DomainType = "register"
	DomainTypeTransfer DomainType = "transfer"
	DomainTypeExisting DomainType = "existing"
)

// DomainStatus define os possíveis status de um domínio no sistema.
type DomainStatus string

const (
	StatusPendingPayment      DomainStatus = "pending_payment"
	StatusPendingProvisioning DomainStatus = "pending_provisioning"
	StatusActive              DomainStatus = "active"
	StatusFailed              DomainStatus = "failed"
	StatusCancelled           DomainStatus = "cancelled"
)

// Domain representa um registro na tabela `domains`.
type Domain struct {
	ID                int            `json:"id"`
	ClientID          int            `json:"client_id"`
	ServiceID         int            `json:"service_id"`
	DomainName        string         `json:"domain_name"`
	Type              DomainType     `json:"type"`
	Status            DomainStatus   `json:"status"`
	Provider          sql.NullString `json:"provider"`
	ProviderOrderID   sql.NullString `json:"provider_order_id"`
	ExpiresAt         sql.NullTime   `json:"expires_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// CreateDomain cria um novo registro de domínio e retorna o ID gerado.
func CreateDomain(tx *sql.Tx, clientID, serviceID int, domainName string, domainType DomainType) (int, error) {
	var domainID int
	query := `
        INSERT INTO domains (client_id, service_id, domain_name, type, status)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	err := tx.QueryRow(query, clientID, serviceID, domainName, domainType, StatusPendingProvisioning).Scan(&domainID)
	if err != nil {
		return 0, err
	}
	return domainID, nil
}

// UpdateDomainStatus atualiza o status de um domínio.
func UpdateDomainStatus(db *sql.DB, domainID int, status DomainStatus) error {
	_, err := db.Exec("UPDATE domains SET status = $1 WHERE id = $2", status, domainID)
	return err
}

// UpdateDomainProviderOrderID atualiza o ID do pedido no provedor.
func UpdateDomainProviderOrderID(db *sql.DB, domainID int, provider, providerOrderID string) error {
	query := "UPDATE domains SET provider = $1, provider_order_id = $2 WHERE id = $3"
	_, err := db.Exec(query, provider, providerOrderID, domainID)
	return err
}
