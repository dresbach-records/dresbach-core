package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"hosting-backend/internal/auth"
	"hosting-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest define a estrutura para uma requisição de login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token,omitempty"` // Para 2FA
}

// LoginResponse define a estrutura para uma resposta de login bem-sucedida.
type LoginResponse struct {
	Token        string `json:"token"`
	TwoFARequired bool   `json:"two_factor_required"`
}

// LoginHandler processa o login com todas as checagens de segurança.
func LoginHandler(db *sql.DB) http.HandlerFunc {
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Pedido inválido", http.StatusBadRequest)
			return
		}

		ipAddress := models.GetIP(r)

		// 1. CHECAGEM DE BLOQUEIO POR FORÇA BRUTA (simulação)
		isBlocked, _ := models.CheckAndBlockIP(db, ipAddress, 10, 5*time.Minute)
		if isBlocked {
			http.Error(w, "Muitas tentativas de login. Tente novamente mais tarde.", http.StatusTooManyRequests)
			return
		}

		var clientID int
		var passwordHash string
		var twoFactorEnabled bool
		var twoFactorSecret sql.NullString
		var enforceIPWhitelist bool
		query := "SELECT id, password_hash, two_factor_enabled, two_factor_secret, enforce_ip_whitelist FROM clients WHERE email = ? AND status = 'active'"
		err := db.QueryRow(query, req.Email).Scan(&clientID, &passwordHash, &twoFactorEnabled, &twoFactorSecret, &enforceIPWhitelist)

		if err != nil {
			if err == sql.ErrNoRows {
				models.LogLoginAttempt(db, req.Email, ipAddress, r.UserAgent(), false)
				http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Erro no servidor", http.StatusInternalServerError)
			return
		}

        // 2. LÓGICA DE LOGIN POR IP
        if enforceIPWhitelist {
            isAllowed, err := models.IsIPAllowedForClient(db, clientID, ipAddress)
            if err != nil {
                http.Error(w, "Erro ao verificar permissão de IP", http.StatusInternalServerError)
                return
            }
            if !isAllowed {
                models.LogLoginAttempt(db, req.Email, ipAddress, r.UserAgent(), false)
                http.Error(w, "Acesso negado para este endereço IP", http.StatusForbidden)
                return
            }
        }

		// 3. CHECAGEM DE SENHA
		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
			models.LogLoginAttempt(db, req.Email, ipAddress, r.UserAgent(), false)
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		}

		// 4. CHECAGEM DE 2FA
		if twoFactorEnabled {
			if !twoFactorSecret.Valid {
				http.Error(w, "2FA ativado mas não configurado corretamente.", http.StatusInternalServerError)
				return
			}
			if req.Token == "" {
				// 2FA é necessário, mas o token não foi fornecido. Peça ao frontend.
				json.NewEncoder(w).Encode(LoginResponse{TwoFARequired: true})
				return
			}
			// Validar o token TOTP
			valid := totp.Validate(req.Token, twoFactorSecret.String)
			if !valid {
				models.LogLoginAttempt(db, req.Email, ipAddress, r.UserAgent(), false)
				http.Error(w, "Token 2FA inválido", http.StatusUnauthorized)
				return
			}
		}

		// 5. GERAÇÃO DO TOKEN JWT
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &auth.Claims{
			UserID: clientID,
			Email:  req.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
			return
		}

		models.LogLoginAttempt(db, req.Email, ipAddress, r.UserAgent(), true) // Log de sucesso
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{Token: tokenString, TwoFARequired: false})
	}
}
