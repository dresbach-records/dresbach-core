package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"hosting-backend/internal/models"
)

// GetSiteSettingsHandler busca as configurações institucionais do site.
// Rota: GET /admin/settings/site
func GetSiteSettingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Adicionar middleware de autenticação e autorização para admin

		settings, err := models.GetSiteSettings(db)
		if err != nil {
			http.Error(w, "Erro ao buscar configurações do site", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}
}

// UpdateSiteSettingsHandler atualiza as configurações institucionais do site.
// Rota: PUT /admin/settings/site
func UpdateSiteSettingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Adicionar middleware de autenticação e autorização para admin

		var settings models.SiteSettings
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		if err := models.UpdateSiteSettings(db, &settings); err != nil {
			// TODO: Logar este erro no audit_log
			http.Error(w, "Erro ao atualizar as configurações do site", http.StatusInternalServerError)
			return
		}

		// TODO: Logar a alteração bem-sucedida no audit_log

		w.WriteHeader(http.StatusNoContent) // Resposta de sucesso sem corpo
	}
}

// UpdateMaintenanceModeHandler controla o modo de manutenção do site.
// Rota: PUT /admin/settings/maintenance
func UpdateMaintenanceModeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Adicionar middleware de autenticação e autorização para admin (permissão: site.manage)

		var payload struct {
			Enabled bool   `json:"enabled"`
			Message string `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		if err := models.UpdateMaintenanceMode(db, payload.Enabled, payload.Message); err != nil {
			// TODO: Logar a falha no audit_log
			http.Error(w, "Erro ao atualizar o modo de manutenção", http.StatusInternalServerError)
			return
		}

		// TODO: Logar a alteração bem-sucedida no audit_log (quem ativou/desativou e quando)

		w.WriteHeader(http.StatusNoContent)
	}
}
