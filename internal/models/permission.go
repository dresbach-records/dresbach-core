package models

import "database/sql"

// GetAllPermissions busca todas as permissões definidas no sistema.
func GetAllPermissions(db *sql.DB) ([]Permission, error) {
	rows, err := db.Query(`SELECT id, name, description FROM permissions ORDER BY name`)
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
