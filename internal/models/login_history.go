package models

import (
	"database/sql"
	"net/http"
	"time"
)

// LoginHistoryEntry representa um único registro no histórico de login para o cliente.
type LoginHistoryEntry struct {
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	WasSuccessful bool      `json:"was_successful"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetLoginHistoryForClient busca o histórico de login recente para um cliente específico.
func GetLoginHistoryForClient(db *sql.DB, clientID int, limit int) ([]LoginHistoryEntry, error) {
	query := `
		SELECT ip_address, user_agent, was_successful, created_at 
		FROM login_history 
		WHERE client_id = ? 
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := db.Query(query, clientID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []LoginHistoryEntry
	for rows.Next() {
		var entry LoginHistoryEntry
		if err := rows.Scan(&entry.IPAddress, &entry.UserAgent, &entry.WasSuccessful, &entry.CreatedAt); err != nil {
			return nil, err // ou logar o erro e continuar
		}
		history = append(history, entry)
	}

	return history, nil
}

// CheckAndBlockIP verifica se um IP deve ser bloqueado.
func CheckAndBlockIP(db *sql.DB, ipAddress string, maxAttempts int, blockDuration time.Duration) (bool, error) {
	// ... (código existente)
	query := `
		SELECT COUNT(*) FROM login_history 
		WHERE ip_address = ? AND was_successful = FALSE AND created_at >= ?
	`
	var failedAttempts int
	timeThreshold := time.Now().Add(-blockDuration)
	err := db.QueryRow(query, ipAddress, timeThreshold).Scan(&failedAttempts)
	if err != nil {
		return false, err
	}
	return failedAttempts >= maxAttempts, nil
}

// LogLoginAttempt registra uma tentativa de login.
func LogLoginAttempt(db *sql.DB, email, ipAddress, userAgent string, wasSuccessful bool) error {
	// ... (código existente)
	var clientID sql.NullInt64
	q := "SELECT id FROM clients WHERE email = ?"
	err := db.QueryRow(q, email).Scan(&clientID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO login_history (client_id, ip_address, user_agent, was_successful) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(clientID, ipAddress, userAgent, wasSuccessful)
	return err
}

// GetIP retorna o endereço IP real de uma requisição.
func GetIP(r *http.Request) string {
	// ... (código existente)
    ip := r.Header.Get("X-Forwarded-For")
    if ip != "" {
        return ip
    }
    return r.RemoteAddr
}
