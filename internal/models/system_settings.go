package models

import (
	"database/sql"
)

// GetSystemSettingsAsMap busca todas as configurações do sistema e as retorna como um mapa.
func GetSystemSettingsAsMap(db *sql.DB) (map[string]string, error) {
	query := `SELECT setting_key, setting_value FROM system_settings`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settingsMap := make(map[string]string)
	for rows.Next() {
		var key string
		var value sql.NullString
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		settingsMap[key] = value.String
	}
	return settingsMap, nil
}

// UpdateSystemSettings atualiza múltiplas configurações do sistema em uma única transação.
// A entrada é um mapa de setting_key para setting_value.
func UpdateSystemSettings(db *sql.DB, settings map[string]string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Garante rollback em caso de pânico ou erro
	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE system_settings SET setting_value = ? WHERE setting_key = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for key, value := range settings {
		_, err := stmt.Exec(value, key)
		if err != nil {
			return err // tx.Rollback() será chamado pelo defer
		}
	}

	// Se tudo correu bem, confirma a transação
	return tx.Commit()
}
