package models

import (
	"database/sql"
	"time"

	"hosting-backend/internal/utils"
)

// ApiProvider define os provedores de API suportados.
type ApiProvider string

const (
	Stripe    ApiProvider = "stripe"
	WHM       ApiProvider = "whm"
	Hostinger ApiProvider = "hostinger"
	SMTP      ApiProvider = "smtp"
)

// ApiCredential armazena credenciais de API de forma segura.
type ApiCredential struct {
	ID           int          `json:"-"`
	Provider     ApiProvider  `json:"provider"`
	EncryptedKey string       `json:"-"` // Nunca exposto!
	Status       string       `json:"status"`
	LastTestAt   sql.NullTime `json:"last_test_at"`
	LastError    sql.NullString `json:"last_error"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// GetApiCredentialsForAdmin lista as credenciais para o painel admin (sem a chave).
func GetApiCredentialsForAdmin(db *sql.DB) ([]ApiCredential, error) {
	query := `SELECT provider, status, last_test_at, last_error, updated_at FROM api_credentials`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creds []ApiCredential
	for rows.Next() {
		var c ApiCredential
		if err := rows.Scan(&c.Provider, &c.Status, &c.LastTestAt, &c.LastError, &c.UpdatedAt); err != nil {
			return nil, err
		}
		creds = append(creds, c)
	}
	return creds, nil
}

// UpsertApiCredential atualiza ou insere uma credencial de API, criptografando a chave.
func UpsertApiCredential(db *sql.DB, provider ApiProvider, apiKey string, status string) error {
	encryptedKey, err := utils.Encrypt([]byte(apiKey))
	if err != nil {
		return err
	}

	query := `INSERT INTO api_credentials (provider, encrypted_key, status) VALUES ($1, $2, $3)
			   ON CONFLICT (provider) DO UPDATE SET encrypted_key = EXCLUDED.encrypted_key, status = EXCLUDED.status`
	_, err = db.Exec(query, provider, encryptedKey, status)
	return err
}

// GetDecryptedApiKey busca e descriptografa uma chave de API para uso interno.
func GetDecryptedApiKey(db *sql.DB, provider ApiProvider) (string, error) {
	var encryptedKey string
	query := `SELECT encrypted_key FROM api_credentials WHERE provider = $1 AND status = 'active'`
	err := db.QueryRow(query, provider).Scan(&encryptedKey)
	if err != nil {
		return "", err
	}

	decryptedKey, err := utils.Decrypt(encryptedKey)
	if err != nil {
		return "", err
	}

	return string(decryptedKey), nil
}

// UpdateApiTestResult atualiza o resultado de um teste de conexão de API.
func UpdateApiTestResult(db *sql.DB, provider ApiProvider, testErr error) error {
	var query string
	var err error

	if testErr != nil {
		query = `UPDATE api_credentials SET last_test_at = NOW(), last_error = $1 WHERE provider = $2`
		_, err = db.Exec(query, testErr.Error(), provider)
	} else {
		query = `UPDATE api_credentials SET last_test_at = NOW(), last_error = NULL WHERE provider = $1`
		_, err = db.Exec(query, provider)
	}
	return err
}
