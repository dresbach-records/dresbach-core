package models

import (
	"database/sql"
)

// AllowedIP representa um endereço IP permitido para um cliente.
type AllowedIP struct {
	ID          int    `json:"id"`
	ClientID    int    `json:"-"` // Omitido do JSON de resposta
	IPAddress   string `json:"ip_address"`
	Description string `json:"description"`
}

// AddAllowedIP adiciona um novo IP à lista branca de um cliente.
func AddAllowedIP(db *sql.DB, clientID int, ipAddress, description string) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO allowed_ips (client_id, ip_address, description) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(clientID, ipAddress, description)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetAllowedIPsForClient retorna a lista de IPs permitidos para um cliente.
func GetAllowedIPsForClient(db *sql.DB, clientID int) ([]AllowedIP, error) {
	rows, err := db.Query("SELECT id, ip_address, description FROM allowed_ips WHERE client_id = ?", clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ips []AllowedIP
	for rows.Next() {
		var ip AllowedIP
		if err := rows.Scan(&ip.ID, &ip.IPAddress, &ip.Description); err != nil {
			return nil, err
		}
		ips = append(ips, ip)
	}
	return ips, nil
}

// DeleteAllowedIP remove um IP da lista branca de um cliente.
func DeleteAllowedIP(db *sql.DB, clientID, ipID int) (bool, error) {
	stmt, err := db.Prepare("DELETE FROM allowed_ips WHERE id = ? AND client_id = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(ipID, clientID)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	return rowsAffected > 0, err
}

// IsIPAllowedForClient verifica se um determinado IP é permitido para um cliente.
// Esta função será usada durante o login.
func IsIPAllowedForClient(db *sql.DB, clientID int, ipAddress string) (bool, error) {
    var count int
    query := "SELECT COUNT(*) FROM allowed_ips WHERE client_id = ? AND ip_address = ?"
    err := db.QueryRow(query, clientID, ipAddress).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}
