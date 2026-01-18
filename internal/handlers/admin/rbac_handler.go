package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"hosting-backend/internal/models"
	"github.com/go-chi/chi/v5"
)

// GetRolesHandler lista todas as funções de administrador.
// Rota: GET /admin/roles
func GetRolesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roles, err := models.GetAllRoles(db)
		if err != nil {
			http.Error(w, "Erro ao buscar funções", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(roles)
	}
}

// GetPermissionsForRoleHandler lista as permissões de uma função específica.
// Rota: GET /admin/roles/{id}/permissions
func GetPermissionsForRoleHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "ID de função inválido", http.StatusBadRequest)
			return
		}

		permissions, err := models.GetPermissionsForRole(db, roleID)
		if err != nil {
			http.Error(w, "Erro ao buscar permissões da função", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(permissions)
	}
}

// AssignPermissionsToRoleHandler atribui permissões a uma função.
// Rota: POST /admin/roles/{id}/permissions
func AssignPermissionsToRoleHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "ID de função inválido", http.StatusBadRequest)
			return
		}

		var payload struct {
			PermissionIDs []int `json:"permission_ids"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		if err := models.AssignPermissionsToRole(db, roleID, payload.PermissionIDs); err != nil {
			http.Error(w, "Erro ao atribuir permissões", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetPermissionsHandler lista todas as permissões disponíveis no sistema.
// Rota: GET /admin/permissions
func GetPermissionsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		permissions, err := models.GetAllPermissions(db)
		if err != nil {
			http.Error(w, "Erro ao buscar permissões", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(permissions)
	}
}
