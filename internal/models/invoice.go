package models

import (
	"database/sql"
	"time"
)

// InvoiceStatus define os status possíveis para uma fatura.
type InvoiceStatus string

const (
	InvoiceStatusPaid   InvoiceStatus = "paid"
	InvoiceStatusUnpaid InvoiceStatus = "unpaid"
	InvoiceStatusVoid   InvoiceStatus = "void"
)

// Invoice representa a estrutura de dados de uma fatura.
type Invoice struct {
	ID          int           `json:"id"`
	UserID      int           `json:"user_id"`
	ServiceID   int           `json:"service_id"` // Essencial para o worker saber qual serviço suspender
	IssueDate   time.Time     `json:"issue_date"`
	DueDate     time.Time     `json:"due_date"`
	TotalAmount float64       `json:"total_amount"`
	Status      InvoiceStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
}

// GetOverdueInvoices busca faturas com status 'unpaid' e cuja data de vencimento já passou.
func GetOverdueInvoices(db *sql.DB) ([]Invoice, error) {
	query := `SELECT id, user_id, service_id, issue_date, due_date, total_amount, status, created_at FROM invoices WHERE status = ? AND due_date < NOW()`

	rows, err := db.Query(query, InvoiceStatusUnpaid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(&i.ID, &i.UserID, &i.ServiceID, &i.IssueDate, &i.DueDate, &i.TotalAmount, &i.Status, &i.CreatedAt); err != nil {
			return nil, err
		}
		invoices = append(invoices, i)
	}

	return invoices, nil
}

// CreateInvoice insere uma nova fatura no banco de dados.
func CreateInvoice(db *sql.DB, invoice *Invoice) error {
	query := `INSERT INTO invoices (user_id, service_id, issue_date, due_date, total_amount, status) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, invoice.UserID, invoice.ServiceID, invoice.IssueDate, invoice.DueDate, invoice.TotalAmount, invoice.Status)
	return err
}

// GetInvoicesByUserID busca todas as faturas de um usuário específico, ordenadas pela data de emissão.
func GetInvoicesByUserID(db *sql.DB, userID int) ([]Invoice, error) {
	query := `SELECT id, user_id, service_id, issue_date, due_date, total_amount, status, created_at FROM invoices WHERE user_id = ? ORDER BY issue_date DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(&i.ID, &i.UserID, &i.ServiceID, &i.IssueDate, &i.DueDate, &i.TotalAmount, &i.Status, &i.CreatedAt); err != nil {
			return nil, err
		}
		invoices = append(invoices, i)
	}
	return invoices, nil
}
