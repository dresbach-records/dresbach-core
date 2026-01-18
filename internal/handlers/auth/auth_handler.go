package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest define a estrutura do corpo da solicitação de login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse define a estrutura da resposta de login bem-sucedido.
type LoginResponse struct {
	Token string `json:"token"`
}

// Claims define as reivindicações personalizadas para o token JWT.
type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// LoginHandler processa as solicitações de login.
func LoginHandler(db *sql.DB) http.HandlerFunc {
	// A chave JWT é lida uma vez no início, para eficiência.
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Pedido inválido", http.StatusBadRequest)
			return
		}

		var id int
		var passwordHash, role string
		err := db.QueryRow("SELECT id, password_hash, role FROM users WHERE email = ? AND is_active = TRUE", req.Email).Scan(&id, &passwordHash, &role)
		if err == sql.ErrNoRows {
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Email: req.Email,
			Role:  role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Erro ao gerar o token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
	}
}
