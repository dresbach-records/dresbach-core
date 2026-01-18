package models

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// AdminUser representa um usuário administrativo.
type AdminUser struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Nunca expor a senha
	Email     string    `json:"email"`
	RoleID    int       `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SetPassword criptografa e define a senha para o usuário.
func (u *AdminUser) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifica se a senha fornecida corresponde ao hash.
func (u *AdminUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// GetUserPermissions busca todas as permissões de um usuário através de sua função.
func GetUserPermissions(db *sql.DB, userID int) (map[string]bool, error) {
	query := `SELECT p.name FROM permissions p 
			   INNER JOIN role_permissions rp ON p.id = rp.permission_id
			   INNER JOIN admin_users u ON rp.role_id = u.role_id
			   WHERE u.id = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		permissions[name] = true
	}
	return permissions, nil
}
