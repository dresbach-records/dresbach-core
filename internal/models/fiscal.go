package models

import (
	"database/sql"
	"time"
)

// SystemModule representa a tabela system_modules
type SystemModule struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Enabled   bool           `json:"enabled"`
	EnabledAt sql.NullTime   `json:"enabled_at"`
	EnabledBy sql.NullInt64  `json:"enabled_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// AuditLog representa a tabela audit_logs
type AuditLog struct {
	ID         int64          `json:"id"`
	UserID     int64          `json:"user_id"`
	Action     string         `json:"action"`
	TargetType sql.NullString `json:"target_type"`
	TargetID   sql.NullInt64  `json:"target_id"`
	Details    sql.NullString `json:"details"`
	IPAddress  sql.NullString `json:"ip_address"`
	CreatedAt  time.Time      `json:"created_at"`
}

// FiscalSettings representa a tabela fiscal_settings
type FiscalSettings struct {
	ID                   int64           `json:"id"`
	Provider             sql.NullString  `json:"provider"`
	CompanyName          sql.NullString  `json:"company_name"`
	CNPJ                 sql.NullString  `json:"cnpj"`
	MunicipalRegistration sql.NullString  `json:"municipal_registration"`
	City                 sql.NullString  `json:"city"`
	State                sql.NullString  `json:"state"`
	ISSRate              sql.NullFloat64 `json:"iss_rate"`
	Environment          sql.NullString  `json:"environment"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

// FiscalDocument representa a tabela fiscal_documents
type FiscalDocument struct {
	ID               int64          `json:"id"`
	InvoiceID        int64          `json:"invoice_id"`
	Provider         sql.NullString `json:"provider"`
	NFNumber         sql.NullString `json:"nf_number"`
	VerificationCode sql.NullString `json:"verification_code"`
	Status           string         `json:"status"`
	PDFURL           sql.NullString `json:"pdf_url"`
	XMLURL           sql.NullString `json:"xml_url"`
	ErrorMessage     sql.NullString `json:"error_message"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}
