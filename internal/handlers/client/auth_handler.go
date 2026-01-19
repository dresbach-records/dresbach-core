package client

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"hosting-backend/internal/models"
	"hosting-backend/internal/services/asaas"
	"hosting-backend/internal/utils"
)

// RegisterPayload define a estrutura de dados para o cadastro de um novo cliente.
type RegisterPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CpfCnpj   string `json:"cpf_cnpj"`
}

// LoginPayload define a estrutura de dados para o login de um cliente.
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler lida com o cadastro de novos clientes.
// Rota: POST /api/register
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload RegisterPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		// Validar se o email já existe
		existingUser, err := models.GetUserByEmail(db, payload.Email)
		if err != nil {
			http.Error(w, "Erro ao verificar o email", http.StatusInternalServerError)
			return
		}
		if existingUser != nil {
			http.Error(w, "Email já cadastrado", http.StatusConflict)
			return
		}

		// Gerar hash da senha
		hashedPassword, err := utils.HashPassword(payload.Password)
		if err != nil {
			http.Error(w, "Erro ao processar a senha", http.StatusInternalServerError)
			return
		}

		// Criar o novo usuário
		newUser := models.User{
			FirstName:    payload.FirstName,
			LastName:     payload.LastName,
			Email:        payload.Email,
			PasswordHash: hashedPassword,
			IsActive:     true, // Usuários são ativados por padrão
		}

		userID, err := models.CreateUser(db, &newUser)
		if err != nil {
			http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
			return
		}

		// Criar o novo cliente
		newClient := models.Client{
			UserID:      int(userID),
			ContactName: sql.NullString{String: payload.FirstName + " " + payload.LastName, Valid: true},
			Email:       payload.Email,
			CpfCnpj:     sql.NullString{String: payload.CpfCnpj, Valid: true},
		}

		// Criar cliente no Asaas
		asaasClient := asaas.NewAsaasClient()
		customerRequest := asaas.CustomerRequest{
			Name:    newClient.ContactName.String,
			CpfCnpj: newClient.CpfCnpj.String,
			Email:   newClient.Email,
		}

		asaasCustomer, err := asaasClient.CreateCustomer(customerRequest)
		if err != nil {
			log.Printf("Erro ao criar cliente no Asaas: %v", err)
			// Mesmo que a criação no Asaas falhe, o usuário é criado no nosso sistema.
			// O erro será logado, e poderemos tratar a sincronização depois.
		} else {
			newClient.AsaasCustomerID = sql.NullString{String: asaasCustomer.ID, Valid: true}
		}

		// Salvar o cliente no banco de dados
		_, err = db.Exec("INSERT INTO clients (user_id, contact_name, email, cpf_cnpj, asaas_customer_id) VALUES (?, ?, ?, ?, ?)",
			newClient.UserID, newClient.ContactName, newClient.Email, newClient.CpfCnpj, newClient.AsaasCustomerID)
		if err != nil {
			http.Error(w, "Erro ao salvar cliente", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser) // Não envia o hash da senha
	}
}

// LoginHandler autentica um cliente e retorna um token JWT.
// @Summary Autentica um cliente
// @Description Autentica um cliente com email e senha e retorna um token JWT se as credenciais estiverem corretas.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   credentials     body    LoginPayload     true        "Credenciais de Login"
// @Success 200 {object} map[string]string "Token JWT"
// @Failure 400 {string} string "Corpo da requisição inválido"
// @Failure 401 {string} string "Credenciais inválidas"
// @Failure 500 {string} string "Erro ao gerar token de autenticação"
// @Router /api/login [post]
func LoginHandler(db *sql.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload LoginPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		// Buscar usuário pelo email
		user, err := models.GetUserByEmail(db, payload.Email)
		if err != nil || user == nil {
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		}

		// Verificar a senha
		if !utils.CheckPasswordHash(payload.Password, user.PasswordHash) {
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		}

		// Gerar o token JWT para o cliente
		token, err := utils.GenerateToken(user.ID, "client")
		if err != nil {
			http.Error(w, "Erro ao gerar token de autenticação", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
