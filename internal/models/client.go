package models

import (
	"database/sql"
	"time"
)

// Client representa a estrutura de dados de um cliente.

type Client struct {
	ID              int            `json:"id"`
	UserID          int            `json:"user_id"`
	CompanyName     string         `json:"company_name,omitempty"`
	ContactName     string         `json:"contact_name,omitempty"`
	Email           string         `json:"email"`
	Phone           string         `json:"phone,omitempty"`
	Address         string         `json:"address,omitempty"`
	City            string         `json:"city,omitempty"`
	State           string         `json:"state,omitempty"`
	Zip             string         `json:"zip,omitempty"`
	Country         string         `json:"country,omitempty"`
	CpfCnpj         string         `json:"cpf_cnpj,omitempty"`
	AsaasCustomerID sql.NullString `json:"asaas_customer_id,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
}

// GetUserByID busca um cliente pelo seu ID.
func GetUserByID(db *sql.DB, userID int) (*Client, error) {
	query := `SELECT id, user_id, company_name, contact_name, email, phone, address, city, state, zip, country, cpf_cnpj, asaas_customer_id, created_at
			  FROM clients WHERE id = ?`
	row := db.QueryRow(query, userID)

	client := &Client{}
	err := row.Scan(
		&client.ID,
		&client.UserID,
		&client.CompanyName,
		&client.ContactName,
		&client.Email,
		&client.Phone,
		&client.Address,
		&client.City,
		&client.State,
		&client.Zip,
		&client.Country,
		&client.CpfCnpj,
		&client.AsaasCustomerID,
		&client.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // Retorna o erro para o handler tratar o caso Not Found
		}
		return nil, err // Outros erros
	}

	return client, nil
}
