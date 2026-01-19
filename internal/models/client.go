package models

import (
	"database/sql"
	"time"
)

// Client representa a estrutura de dados de um cliente.
type Client struct {
	ID              int            `json:"id"`
	UserID          int            `json:"user_id"`
	CompanyName     sql.NullString `json:"company_name,omitempty"`
	ContactName     sql.NullString `json:"contact_name,omitempty"`
	Email           string         `json:"email"`
	Phone           sql.NullString `json:"phone,omitempty"`
	Address         sql.NullString `json:"address,omitempty"`
	City            sql.NullString `json:"city,omitempty"`
	State           sql.NullString `json:"state,omitempty"`
	Zip             sql.NullString `json:"zip,omitempty"`
	Country         sql.NullString `json:"country,omitempty"`
	CpfCnpj         sql.NullString `json:"cpf_cnpj,omitempty"`
	AsaasCustomerID sql.NullString `json:"asaas_customer_id,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
}

// CreateClient insere um novo cliente no banco de dados.
func CreateClient(db *sql.DB, client *Client) (int64, error) {
	query := `INSERT INTO clients (user_id, company_name, contact_name, email, phone, address, city, state, zip, country) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, client.UserID, client.CompanyName, client.ContactName, client.Email, client.Phone, client.Address, client.City, client.State, client.Zip, client.Country)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetAllClients retorna todos os clientes do banco de dados.
func GetAllClients(db *sql.DB) ([]Client, error) {
	query := `SELECT id, user_id, company_name, contact_name, email, phone, address, city, state, zip, country, created_at FROM clients`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var c Client
		if err := rows.Scan(&c.ID, &c.UserID, &c.CompanyName, &c.ContactName, &c.Email, &c.Phone, &c.Address, &c.City, &c.State, &c.Zip, &c.Country, &c.CreatedAt); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

// GetClientByID busca um cliente pelo seu ID de cliente.
func GetClientByID(db *sql.DB, clientID int) (*Client, error) {
	query := `SELECT id, user_id, company_name, contact_name, email, phone, address, city, state, zip, country, cpf_cnpj, asaas_customer_id, created_at FROM clients WHERE id = ?`
	row := db.QueryRow(query, clientID)

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
			return nil, nil // Not Found
		}
		return nil, err // Outros erros
	}

	return client, nil
}

// GetClientByUserID busca um cliente pelo ID do usuário associado.
func GetClientByUserID(db *sql.DB, userID int) (*Client, error) {
	query := `SELECT id, user_id, company_name, contact_name, email, phone, address, city, state, zip, country, cpf_cnpj, asaas_customer_id, created_at FROM clients WHERE user_id = ?`
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
			return nil, nil // Not Found
		}
		return nil, err // Outros erros
	}

	return client, nil
}

// UpdateClient atualiza um cliente no banco de dados.
func UpdateClient(db *sql.DB, client *Client) error {
	query := `UPDATE clients SET company_name = ?, contact_name = ?, email = ?, phone = ?, address = ?, city = ?, state = ?, zip = ?, country = ? WHERE id = ?`
	_, err := db.Exec(query, client.CompanyName, client.ContactName, client.Email, client.Phone, client.Address, client.City, client.State, client.Zip, client.Country, client.ID)
	return err
}

// DeleteClient deleta um cliente do banco de dados.
func DeleteClient(db *sql.DB, id int) error {
	query := `DELETE FROM clients WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}
