package auth

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"

	"hosting-backend/internal/models"
	"github.com/pquerna/otp/totp"
)

// Setup2FAResponse define a resposta para a solicitação de configuração do 2FA.
type Setup2FAResponse struct {
	Secret    string `json:"secret"`     // O segredo TOTP, para backup manual
	QRCode    string `json:"qr_code"`    // O QR code como uma string de imagem em base64
	Issuer    string `json:"issuer"`     // O nome do serviço (seu aplicativo)
	AccountName string `json:"account_name"` // O e-mail do usuário
}

// Enable2FARequest define o corpo da requisição para ativar o 2FA.
type Enable2FARequest struct {
	Token string `json:"token"`
}

// Setup2FAHandler gera um segredo TOTP e um QR Code para o usuário.
func Setup2FAHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Obter o ID do cliente e o e-mail do token JWT (contexto)
		claims, ok := GetClaims(r.Context())
		if !ok {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}
		clientID := claims.UserID
		clientEmail := claims.Email

		// 2. Gerar uma nova chave TOTP
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "Seu Hosting Inc",
			AccountName: clientEmail,
		})
		if err != nil {
			http.Error(w, "Erro ao gerar chave 2FA", http.StatusInternalServerError)
			return
		}

		// 3. Gerar o QR Code
		var buf bytes.Buffer
		img, err := key.Image(200, 200)
		if err != nil {
			http.Error(w, "Erro ao gerar QR code", http.StatusInternalServerError)
			return
		}
		png.Encode(&buf, img)
		qrCodeString := base64.StdEncoding.EncodeToString(buf.Bytes())

		// 4. Salvar o segredo *temporariamente* no banco para verificação posterior
		// Ou, em uma abordagem stateless, incluir o segredo em um JWT de curta duração.
		// Por simplicidade aqui, vamos salvar na tabela de usuário, mas não ativado.
		query := "UPDATE clients SET two_factor_secret = ? WHERE id = ?"
		if _, err := db.Exec(query, key.Secret(), clientID); err != nil {
			http.Error(w, "Erro ao salvar segredo 2FA", http.StatusInternalServerError)
			return
		}

		// 5. Responder com o segredo e o QR code
		resp := Setup2FAResponse{
			Secret:    key.Secret(),
			QRCode:    "data:image/png;base64," + qrCodeString,
			Issuer:      key.Issuer(),
			AccountName: key.AccountName(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// Enable2FAHandler valida o token TOTP e ativa o 2FA para o usuário.
func Enable2FAHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetClaims(r.Context())
		if !ok {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}
		clientID := claims.UserID

		var req Enable2FARequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Pedido inválido", http.StatusBadRequest)
			return
		}

		// 1. Buscar o segredo do banco de dados
		var secret string
		query := "SELECT two_factor_secret FROM clients WHERE id = ?"
		if err := db.QueryRow(query, clientID).Scan(&secret); err != nil {
			http.Error(w, "Não foi possível encontrar o segredo 2FA", http.StatusNotFound)
			return
		}

		// 2. Validar o token TOTP
		valid := totp.Validate(req.Token, secret)
		if !valid {
			http.Error(w, "Token 2FA inválido", http.StatusUnauthorized)
			return
		}

		// 3. Marcar o 2FA como ativado no banco
		updateQuery := "UPDATE clients SET two_factor_enabled = TRUE WHERE id = ?"
		if _, err := db.Exec(updateQuery, clientID); err != nil {
			http.Error(w, "Erro ao ativar o 2FA", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
