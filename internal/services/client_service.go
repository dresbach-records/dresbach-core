package services

import (
	"database/sql"

	"hosting-backend/internal/models"
)

// ClientService encapsula a lógica de negócios para clientes.
type ClientService struct {
	db *sql.DB
}

// NewClientService cria uma nova instância de ClientService.
func NewClientService(db *sql.DB) *ClientService {
	return &ClientService{db: db}
}

// GetClientServices retorna os serviços de um cliente.
func (s *ClientService) GetClientServices(userID int) ([]models.Service, error) {
	return models.GetServicesByUserID(s.db, userID)
}

// GetClientInvoices retorna as faturas de um cliente.
func (s *ClientService) GetClientInvoices(userID int) ([]models.Invoice, error) {
	return models.GetInvoicesByUserID(s.db, userID)
}
