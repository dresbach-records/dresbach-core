package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"hosting-backend/internal/models"
)

// GetSystemSettingsHandler busca todas as configurações de sistema.
// Rota: GET /admin/settings/system
func GetSystemSettingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		settings, err := models.GetSystemSettingsAsMap(db)
		if err != nil {
			http.Error(w, "Erro ao buscar configurações do sistema", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}
}

// UpdateSystemSettingsHandler atualiza as configurações de sistema.
// Rota: PUT /admin/settings/system
func UpdateSystemSettingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var settings map[string]string
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		if err := models.UpdateSystemSettings(db, settings); err != nil {
			// TODO: Logar este erro no audit_log
			http.Error(w, "Erro ao atualizar as configurações do sistema", http.StatusInternalServerError)
			return
		}

		// TODO: Logar a alteração bem-sucedida no audit_log (quais chaves foram alteradas)

		w.WriteHeader(http.StatusNoContent)
	}
}
