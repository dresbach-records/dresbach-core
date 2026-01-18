package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"hosting-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

// ...(structs)... 

// LoginHandler processa o login com todas as checagens de segurança.
func LoginHandler(db *sql.DB) http.HandlerFunc {
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	return func(w http.ResponseWriter, r *http.Request) {
		// ...(código de decode e get IP existente)... 
		ipAddress := models.GetIP(r)

		// 1. CHECAGEM DE BLOQUEIO POR FORÇA BRUTA
		// ...(código existente)... 

		var clientID int
		var passwordHash string
		var twoFactorEnabled bool
		var enforceIPWhitelist bool // Nova variável
		query := "SELECT id, password_hash, two_factor_enabled, enforce_ip_whitelist FROM clients WHERE email = ? AND status = 'active'"
		err := db.QueryRow(query, req.Email).Scan(&clientID, &passwordHash, &twoFactorEnabled, &enforceIPWhitelist)
		if err != nil {
            // ...(lógica de log de falha existente)... 
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
            return
        }

        // *** 2. NOVA LÓGICA DE LOGIN POR IP ***
        if enforceIPWhitelist {
            isAllowed, err := models.IsIPAllowedForClient(db, clientID, ipAddress)
            if err != nil {
                http.Error(w, "Erro ao verificar permissão de IP", http.StatusInternalServerError)
                return
            }
            if !isAllowed {
                models.LogLoginAttempt(db, req.Email, ipAddress, r.UserAgent(), false) // Loga a falha
                http.Error(w, "Acesso negado para este endereço IP", http.StatusForbidden)
                return
            }
        }

		// 3. CHECAGEM DE SENHA
		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
            // ...(lógica de log de falha existente)... 
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		}

        // ...(lógica de log de sucesso e 2FA existente)... 
	}
}

// ...(outros handlers)... 
