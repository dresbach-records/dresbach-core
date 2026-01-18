package models

import (
	"database/sql"
	"time"
)

// User representa um cliente final do serviço de hospedagem.
type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // O hash da senha nunca deve ser exposto na API
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateUser insere um novo cliente no banco de dados.
func CreateUser(db *sql.DB, user *User) (int64, error) {
	query := `INSERT INTO users (first_name, last_name, email, password_hash, is_active)
			 VALUES (?, ?, ?, ?, ?)`

	result, err := db.Exec(query, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.IsActive)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetUserByEmail busca um usuário pelo seu endereço de email.
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	var user User
	query := `SELECT id, first_name, last_name, email, password_hash, is_active, created_at, updated_at
			 FROM users WHERE email = ?`

	row := db.QueryRow(query, email)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Usuário não encontrado, não é um erro
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID busca um usuário pelo seu ID.
func GetUserByID(db *sql.DB, id int) (*User, error) {
	var user User
	query := `SELECT id, first_name, last_name, email, is_active, created_at, updated_at
			 FROM users WHERE id = ?`

	row := db.QueryRow(query, id)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Usuário não encontrado
		}
		return nil, err
	}
	// Note que não estamos buscando o PasswordHash aqui por segurança.
	return &user, nil
}
