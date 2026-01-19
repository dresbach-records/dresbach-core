package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"hosting-backend/internal/models"
)

// --- Funções de Acesso ao Banco de Dados (integradas no handler) ---

func getFiscalSettings(db *sql.DB) (*models.FiscalSettings, error) {
	var settings models.FiscalSettings
	query := "SELECT id, provider, company_name, cnpj, municipal_registration, city, state, iss_rate, environment, created_at, updated_at FROM fiscal_settings WHERE id = 1"
	err := db.QueryRow(query).Scan(
		&settings.ID, &settings.Provider, &settings.CompanyName, &settings.CNPJ, 
		&settings.MunicipalRegistration, &settings.City, &settings.State, 
		&settings.ISSRate, &settings.Environment, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func updateFiscalSettings(db *sql.DB, settings *models.FiscalSettings) error {
	query := `UPDATE fiscal_settings SET 
		provider = ?, company_name = ?, cnpj = ?, municipal_registration = ?, 
		city = ?, state = ?, iss_rate = ?, environment = ? 
		WHERE id = 1`
	_, err := db.Exec(query, settings.Provider, settings.CompanyName, settings.CNPJ, settings.MunicipalRegistration, settings.City, settings.State, settings.ISSRate, settings.Environment)
	return err
}


// --- Handlers da API ---

// GetFiscalSettingsHandler retorna as configurações fiscais atuais.
func GetFiscalSettingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		settings, err := getFiscalSettings(db)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Configurações fiscais ainda não definidas", http.StatusNotFound)
				return
			}
			http.Error(w, "Erro ao buscar configurações fiscais: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}
}

// UpdateFiscalSettingsHandler atualiza as configurações fiscais.
func UpdateFiscalSettingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var settings models.FiscalSettings
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			http.Error(w, "Payload inválido: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := updateFiscalSettings(db, &settings); err != nil {
			http.Error(w, "Erro ao atualizar as configurações fiscais: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Configurações fiscais atualizadas com sucesso!"})
	}
}
