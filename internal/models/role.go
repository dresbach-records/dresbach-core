package models

import (
	"database/sql"
	"time"
)

// Role representa uma função de um usuário administrativo no sistema.
type Role struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Permission representa uma ação específica que pode ser permitida.
type Permission struct {
	ID          int    `json:"id"`
	Name        string `json:"name"` // Ex: "users.create"
	Description string `json:"description"`
}

// GetAllRoles busca todas as funções de admin.
func GetAllRoles(db *sql.DB) ([]Role, error) {
	rows, err := db.Query(`SELECT id, name, description, created_at, updated_at FROM roles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var r Role
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}

// AssignPermissionsToRole associa uma lista de IDs de permissão a uma função.
// Ele limpa as permissões antigas e define as novas em uma transação.
func AssignPermissionsToRole(db *sql.DB, roleID int, permissionIDs []int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback em caso de erro

	// 1. Limpa permissões existentes para esta função
	_, err = tx.Exec(`DELETE FROM role_permissions WHERE role_id = ?`, roleID)
	if err != nil {
		return err
	}

	// 2. Insere as novas permissões
	if len(permissionIDs) > 0 {
		stmt, err := tx.Prepare(`INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, permID := range permissionIDs {
			if _, err := stmt.Exec(roleID, permID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// GetPermissionsForRole busca todas as permissões associadas a uma função.
func GetPermissionsForRole(db *sql.DB, roleID int) ([]Permission, error) {
	query := `SELECT p.id, p.name, p.description FROM permissions p 
			   INNER JOIN role_permissions rp ON p.id = rp.permission_id
			   WHERE rp.role_id = ?`
	rows, err := db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}
