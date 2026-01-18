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
	StatusPendingPayment     DomainStatus = "pending_payment"
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

// CreateDomain cria um novo registro de domínio, geralmente no início do processo de compra.
func CreateDomain(db *sql.DB, clientID, serviceID int, domainName string, domainType DomainType) (int64, error) {
	stmt, err := db.Prepare(`
        INSERT INTO domains (client_id, service_id, domain_name, type, status)
        VALUES (?, ?, ?, ?, ?)
    `)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(clientID, serviceID, domainName, domainType, StatusPendingPayment)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateDomainStatus atualiza o status de um domínio.
// Esta será uma função central no nosso orquestrador.
func UpdateDomainStatus(db *sql.DB, domainID int, status DomainStatus) error {
	_, err := db.Exec("UPDATE domains SET status = ? WHERE id = ?", status, domainID)
	return err
}
