package models

import (
	"database/sql"
	"time"
)

// ServiceStatus define o status de um serviço contratado.
type ServiceStatus string

const (
	ServiceStatusActive    ServiceStatus = "active"
	ServiceStatusSuspended ServiceStatus = "suspended"
	ServiceStatusTerminated ServiceStatus = "terminated"
	ServiceStatusPending   ServiceStatus = "pending"
)

// BillingCycle define os possíveis ciclos de faturamento.
type BillingCycle string

const (
	BillingCycleMonthly  BillingCycle = "monthly"
	BillingCycleAnnually BillingCycle = "annually"
)

// Service representa um serviço contratado por um cliente.
type Service struct {
	ID           int           `json:"id"`
	UserID       int           `json:"user_id"`
	ProductID    int           `json:"product_id"`
	Domain       string        `json:"domain"`
	CpanelUser   string        `json:"cpanel_user,omitempty"`
	Status       ServiceStatus `json:"status"`
	Price        float64       `json:"price"`
	BillingCycle BillingCycle  `json:"billing_cycle"`
	NextDueDate  time.Time     `json:"next_due_date"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// CreateService insere um novo serviço no banco de dados.
func CreateService(db *sql.DB, s *Service) (int64, error) {
	query := `INSERT INTO services (user_id, product_id, domain, cpanel_user, status, price, billing_cycle, next_due_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(query, s.UserID, s.ProductID, s.Domain, s.CpanelUser, s.Status, s.Price, s.BillingCycle, s.NextDueDate)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetServicesByUserID busca todos os serviços associados a um ID de cliente.
func GetServicesByUserID(db *sql.DB, userID int) ([]Service, error) {
	query := `SELECT id, user_id, product_id, domain, cpanel_user, status, price, billing_cycle, next_due_date, created_at, updated_at FROM services WHERE user_id = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var s Service
		if err := rows.Scan(&s.ID, &s.UserID, &s.ProductID, &s.Domain, &s.CpanelUser, &s.Status, &s.Price, &s.BillingCycle, &s.NextDueDate, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

// GetServiceByID busca um serviço específico pelo seu ID e pelo ID do cliente.
func GetServiceByID(db *sql.DB, serviceID int, userID int) (*Service, error) {
	query := `SELECT id, user_id, product_id, domain, cpanel_user, status, price, billing_cycle, next_due_date, created_at, updated_at FROM services WHERE id = ? AND user_id = ?`
	row := db.QueryRow(query, serviceID, userID)

	var s Service
	if err := row.Scan(&s.ID, &s.UserID, &s.ProductID, &s.Domain, &s.CpanelUser, &s.Status, &s.Price, &s.BillingCycle, &s.NextDueDate, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &s, nil
}

// UpdateServiceStatus atualiza o status de um serviço específico.
func UpdateServiceStatus(db *sql.DB, serviceID int, newStatus ServiceStatus) error {
	query := `UPDATE services SET status = ?, updated_at = NOW() WHERE id = ?`
	_, err := db.Exec(query, newStatus, serviceID)
	return err
}

// GetServicesDueForInvoicing busca serviços ativos cuja data de próxima fatura já passou.
func GetServicesDueForInvoicing(db *sql.DB) ([]Service, error) {
	query := `SELECT id, user_id, product_id, domain, cpanel_user, status, price, billing_cycle, next_due_date, created_at, updated_at FROM services WHERE status = ? AND next_due_date <= NOW()`
	rows, err := db.Query(query, ServiceStatusActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var s Service
		if err := rows.Scan(&s.ID, &s.UserID, &s.ProductID, &s.Domain, &s.CpanelUser, &s.Status, &s.Price, &s.BillingCycle, &s.NextDueDate, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

// UpdateNextDueDate atualiza a data da próxima fatura de um serviço.
func UpdateNextDueDate(db *sql.DB, serviceID int, newDueDate time.Time) error {
	query := `UPDATE services SET next_due_date = ?, updated_at = NOW() WHERE id = ?`
	_, err := db.Exec(query, newDueDate, serviceID)
	return err
}
