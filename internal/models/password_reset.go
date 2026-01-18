package models

import (
	"database/sql"
	"crypto/rand"
	"encoding/hex"
	"time"
)

// PasswordReset representa um token de redefinição de senha no banco de dados.
type PasswordReset struct {
	ID        int
	ClientID  int
	Token     string
	ExpiresAt time.Time
}

// CreatePasswordResetToken gera um novo token para um cliente e o salva no banco de dados.
func CreatePasswordResetToken(db *sql.DB, clientID int) (string, error) {
	// Gerar um token aleatório e seguro
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	// Definir o tempo de expiração (ex: 1 hora a partir de agora)
	expiresAt := time.Now().Add(1 * time.Hour)

	stmt, err := db.Prepare("INSERT INTO password_resets (client_id, token, expires_at) VALUES (?, ?, ?)")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(clientID, token, expiresAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetClientIDByToken valida um token e retorna o ID do cliente associado.
// Ele também deleta o token após o uso bem-sucedido para evitar reutilização.
func GetClientIDByToken(db *sql.DB, token string) (int, error) {
	var pr PasswordReset

	query := "SELECT id, client_id, expires_at FROM password_resets WHERE token = ?"
	err := db.QueryRow(query, token).Scan(&pr.ID, &pr.ClientID, &pr.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, sql.ErrNoRows // sql.ErrNoRows é um erro esperado (token inválido)
		}
		return 0, err
	}

	// Verificar se o token expirou
	if time.Now().After(pr.ExpiresAt) {
		// Token expirado, vamos deletá-lo
		_ = deleteToken(db, token)
		return 0, sql.ErrNoRows // Tratar como se não existisse
	}

	// Token é válido e foi usado, deletar para segurança
	_ = deleteToken(db, token)

	return pr.ClientID, nil
}

// deleteToken remove um token do banco de dados.
func deleteToken(db *sql.DB, token string) error {
	stmt, err := db.Prepare("DELETE FROM password_resets WHERE token = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token)
	return err
}
