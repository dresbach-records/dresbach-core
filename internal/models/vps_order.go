package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// VpsOrder representa um pedido de servidor VPS no banco de dados.
type VpsOrder struct {
	ID            int            `json:"id"`
	ClientID      int            `json:"client_id"`
	InvoiceID     int64          `json:"invoice_id"`
	VpsInstanceID sql.NullString `json:"vps_instance_id"`
	PlanID        string         `json:"plan_id"`
	Location      string         `json:"location"`
	Template      string         `json:"template"`
	Hostname      string         `json:"hostname"`
	Password      string         `json:"-"` // Apenas para entrada, não para armazenamento
	PasswordHash  string         `json:"-"`
	Status        string         `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// Create insere um novo pedido de VPS no banco de dados dentro de uma transação.
func (vo *VpsOrder) Create(tx *sql.Tx) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(vo.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	vo.PasswordHash = string(hashedPassword)

	stmt, err := tx.Prepare(`INSERT INTO vps_orders (client_id, invoice_id, plan_id, location, template, hostname, password_hash, status) 
						   VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(vo.ClientID, vo.InvoiceID, vo.PlanID, vo.Location, vo.Template, vo.Hostname, vo.PasswordHash, vo.Status)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// UpdateVpsInstance atualiza o status e o ID da instância de um pedido de VPS.
func UpdateVpsInstance(tx *sql.Tx, orderID int, vpsInstanceID, status string) error {
	stmt, err := tx.Prepare("UPDATE vps_orders SET vps_instance_id = ?, status = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(vpsInstanceID, status, orderID)
	return err
}
