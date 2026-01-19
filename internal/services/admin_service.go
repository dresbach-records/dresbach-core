package services

import (
	"database/sql"

	"hosting-backend/internal/models"
)

// AdminService encapsula a lógica de negócios para administração.
type AdminService struct {
	db *sql.DB
}

// NewAdminService cria uma nova instância de AdminService.
func NewAdminService(db *sql.DB) *AdminService {
	return &AdminService{db: db}
}

// CreateClient cria um novo cliente.
func (s *AdminService) CreateClient(client *models.Client) (int64, error) {
	return models.CreateClient(s.db, client)
}

// GetAllClients retorna todos os clientes.
func (s *AdminService) GetAllClients() ([]models.Client, error) {
	return models.GetAllClients(s.db)
}

// GetClientByID retorna um cliente pelo ID.
func (s *AdminService) GetClientByID(id int) (*models.Client, error) {
	return models.GetClientByID(s.db, id)
}

// UpdateClient atualiza um cliente.
func (s *AdminService) UpdateClient(client *models.Client) error {
	return models.UpdateClient(s.db, client)
}

// DeleteClient deleta um cliente.
func (s *AdminService) DeleteClient(id int) error {
	return models.DeleteClient(s.db, id)
}
