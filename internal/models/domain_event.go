package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// DomainEvent representa um registro na tabela `domain_events`.
// É usado para auditoria e rastreamento de estado.
type DomainEvent struct {
	ID        int64
	DomainID  int
	Type      string
	Message   sql.NullString
	RawData   sql.NullString // Armazenado como JSON string
	CreatedAt time.Time
}

// LogDomainEvent registra um evento no ciclo de vida de um domínio.
// `data` é um objeto opcional que será serializado para JSON.
func LogDomainEvent(db *sql.DB, domainID int, eventType, message string, data interface{}) error {
	var rawDataJSON []byte
	var err error

	if data != nil {
		rawDataJSON, err = json.Marshal(data)
		if err != nil {
			return err // Falha ao serializar os dados
		}
	}

	stmt, err := db.Prepare(`
        INSERT INTO domain_events (domain_id, type, message, raw_data)
        VALUES ($1, $2, $3, $4)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(domainID, eventType, message, string(rawDataJSON))
	return err
}

// HasEventOccurred verifica se um evento específico já ocorreu para um domínio.
// Isso é chave para garantir a idempotência.
func HasEventOccurred(db *sql.DB, domainID int, eventType string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM domain_events WHERE domain_id = $1 AND type = $2"
	err := db.QueryRow(query, domainID, eventType).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
