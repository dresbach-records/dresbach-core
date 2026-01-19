package services

import (
	"database/sql"
	"fmt"
	"hosting-backend/internal/models"
	"hosting-backend/internal/provisioning"
	"log"
)

// AdminService encapsula a lógica de negócios para administração.
type AdminService struct {
	db             *sql.DB
	whmProvisioner *provisioning.WhmProvisioner
}

// NewAdminService cria uma nova instância de AdminService.
func NewAdminService(db *sql.DB, whmProvisioner *provisioning.WhmProvisioner) *AdminService {
	return &AdminService{
		db:             db,
		whmProvisioner: whmProvisioner,
	}
}

// CreateClientAndProvisionAccount cria um novo cliente no banco de dados e provisiona a conta no WHM.
func (s *AdminService) CreateClientAndProvisionAccount(client *models.Client, domain, username, plan, password string) (int64, error) {
	// Passo 1: Inserir o cliente no banco de dados.
	log.Printf("Iniciando a criação do cliente '%s' no banco de dados.", client.Email)
	clientID, err := models.CreateClient(s.db, client)
	if err != nil {
		return 0, fmt.Errorf("falha ao criar cliente no banco de dados: %w", err)
	}
	log.Printf("Cliente '%s' criado com sucesso no banco de dados com ID: %d.", client.Email, clientID)

	// Passo 2: Provisionar a conta no cPanel/WHM.
	log.Printf("Iniciando o provisionamento da conta no WHM para o domínio '%s'.", domain)
	_, err = s.whmProvisioner.CreateAccount(username, domain, plan, password, client.Email)
	if err != nil {
		// Opcional: Implementar lógica de rollback.
		// Por exemplo, deletar o cliente do banco de dados se o provisionamento falhar.
		log.Printf("ERRO CRÍTICO: Falha ao provisionar a conta no WHM para o cliente ID %d: %v", clientID, err)
		// Neste ponto, você pode querer retornar um erro específico para que o chamador possa tratar.
		// Por enquanto, vamos apenas retornar o erro do provisionador.
		return clientID, fmt.Errorf("o cliente foi criado no banco de dados (ID: %d), mas falhou ao provisionar no WHM: %w", clientID, err)
	}

	log.Printf("Conta para o domínio '%s' provisionada com sucesso no WHM.", domain)
	return clientID, nil
}

// CreateClient cria um novo cliente (sem provisionamento).
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
