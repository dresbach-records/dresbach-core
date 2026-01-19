package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)


// Product representa um plano de hospedagem que pode ser vendido.
type Product struct {
	ID           int                    `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Price        int64                  `json:"price"` // Usar centavos para evitar problemas com ponto flutuante
	BillingCycle BillingCycle             `json:"billing_cycle"`
	Features     map[string]interface{} `json:"features"` // Ex: {"disk_space_mb": 1000, "databases": 5}
	IsActive     bool                   `json:"is_active"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// CreateProduct insere um novo produto no banco de dados.
func CreateProduct(db *sql.DB, p *Product) (int64, error) {
	featuresJSON, err := json.Marshal(p.Features)
	if err != nil {
		return 0, fmt.Errorf("erro ao serializar features: %w", err)
	}

	query := `INSERT INTO products (name, description, price, billing_cycle, features, is_active)
			 VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, p.Name, p.Description, p.Price, p.BillingCycle, string(featuresJSON), p.IsActive)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetAllProducts busca todos os produtos do banco de dados.
func GetAllProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query(`SELECT id, name, description, price, billing_cycle, features, is_active, created_at, updated_at FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		var featuresJSON string

		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.BillingCycle, &featuresJSON, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(featuresJSON), &p.Features); err != nil {
			p.Features = make(map[string]interface{})
		}
		products = append(products, p)
	}
	return products, nil
}

// GetProductByID busca um único produto pelo seu ID.
func GetProductByID(db *sql.DB, id int) (*Product, error) {
	var p Product
	var featuresJSON string

	query := `SELECT id, name, description, price, billing_cycle, features, is_active, created_at, updated_at FROM products WHERE id = ?`
	row := db.QueryRow(query, id)

	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.BillingCycle, &featuresJSON, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Produto não encontrado
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(featuresJSON), &p.Features); err != nil {
		return nil, fmt.Errorf("erro ao deserializar features: %w", err)
	}

	return &p, nil
}

// UpdateProduct atualiza um produto existente no banco de dados.
func UpdateProduct(db *sql.DB, p *Product) error {
	featuresJSON, err := json.Marshal(p.Features)
	if err != nil {
		return fmt.Errorf("erro ao serializar features: %w", err)
	}

	query := `UPDATE products SET name = ?, description = ?, price = ?, billing_cycle = ?, features = ?, is_active = ? WHERE id = ?`

	_, err = db.Exec(query, p.Name, p.Description, p.Price, p.BillingCycle, string(featuresJSON), p.IsActive, p.ID)
	return err
}

// DeleteProduct remove um produto do banco de dados.
func DeleteProduct(db *sql.DB, id int) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}
