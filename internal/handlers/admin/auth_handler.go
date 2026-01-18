package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"hosting-backend/internal/models"
	"hosting-backend/internal/utils"
)

// LoginPayload é a estrutura de dados esperada no corpo da requisição de login.
type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler gerencia a autenticação do usuário administrador.
// Rota: POST /admin/login
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload LoginPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		// 1. Buscar usuário pelo nome de usuário
		// Esta query é um exemplo. Em uma aplicação real, você teria uma função no modelo como GetAdminUserByUsername.
		var user models.AdminUser
		query := `SELECT id, username, password, role_id FROM admin_users WHERE username = ?`
		row := db.QueryRow(query, payload.Username)
		if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.RoleID); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Erro no servidor", http.StatusInternalServerError)
			return
		}

		// 2. Verificar a senha
		if !user.CheckPassword(payload.Password) {
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		}

		// 3. Buscar o nome da função (role)
		var roleName string
		db.QueryRow(`SELECT name FROM roles WHERE id = ?`, user.RoleID).Scan(&roleName)

		// 4. Gerar o token JWT
		token, err := utils.GenerateToken(user.ID, roleName)
		if err != nil {
			http.Error(w, "Erro ao gerar o token", http.StatusInternalServerError)
			return
		}

		// 5. Retornar o token
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}
