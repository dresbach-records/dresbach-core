package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// ServerType define os tipos de painel de controle de servidor suportados.
type ServerType string

const (
	ServerTypeCPanel ServerType = "cpanel"
	ServerTypePlesk  ServerType = "plesk"
	// Outros tipos podem ser adicionados aqui
)

// ServerStatus define o status operacional de um servidor.
type ServerStatus string

const (
	ServerStatusActive      ServerStatus = "active"
	ServerStatusInactive    ServerStatus = "inactive"
	ServerStatusMaintenance ServerStatus = "maintenance"
)

// Server representa um servidor de hospedagem (e.g., cPanel/WHM).
type Server struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Hostname      string          `json:"hostname"`
	IPAddress     string          `json:"ip_address"`
	ServerType    ServerType      `json:"server_type"`
	APIKeyID      int             `json:"api_key_id"` // FK para api_credentials
	Status        ServerStatus    `json:"status"`
	IsDefault     bool            `json:"is_default"`
	Metrics       json.RawMessage `json:"metrics"` // Armazena dados como uso de disco, etc.
	LastSyncAt    sql.NullTime    `json:"last_sync_at"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// GetAllServers busca todos os servidores no banco de dados.
func GetAllServers(db *sql.DB) ([]Server, error) {
	query := `SELECT id, name, hostname, ip_address, server_type, api_key_id, status, is_default, metrics, last_sync_at, created_at, updated_at FROM servers`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []Server
	for rows.Next() {
		var s Server
		if err := rows.Scan(&s.ID, &s.Name, &s.Hostname, &s.IPAddress, &s.ServerType, &s.APIKeyID, &s.Status, &s.IsDefault, &s.Metrics, &s.LastSyncAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}
	return servers, nil
}

// CreateServer insere um novo servidor no banco de dados.
func CreateServer(db *sql.DB, s *Server) (int64, error) {
	query := `INSERT INTO servers (name, hostname, ip_address, server_type, api_key_id, status, is_default, metrics) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(query, s.Name, s.Hostname, s.IPAddress, s.ServerType, s.APIKeyID, s.Status, s.IsDefault, s.Metrics)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateServer atualiza um servidor existente.
func UpdateServer(db *sql.DB, s *Server) error {
	query := `UPDATE servers SET name = ?, hostname = ?, ip_address = ?, server_type = ?, api_key_id = ?, status = ?, is_default = ?, metrics = ? WHERE id = ?`
	_, err := db.Exec(query, s.Name, s.Hostname, s.IPAddress, s.ServerType, s.APIKeyID, s.Status, s.IsDefault, s.Metrics, s.ID)
	return err
}

// DeleteServer remove um servidor do banco de dados.
func DeleteServer(db *sql.DB, id int) error {
	query := `DELETE FROM servers WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}
