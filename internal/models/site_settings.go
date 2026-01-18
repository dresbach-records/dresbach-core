package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// SiteSettings representa as configurações globais do site.
type SiteSettings struct {
	ID                  int             `json:"-"` // Oculto no JSON de resposta
	CompanyName         sql.NullString  `json:"company_name"`
	Slogan              sql.NullString  `json:"slogan"`
	Description         sql.NullString  `json:"description"`
	PhoneNumbers        json.RawMessage `json:"phone_numbers"`
	Whatsapp            sql.NullString  `json:"whatsapp"`
	InstitutionalEmail  sql.NullString  `json:"institutional_email"`
	Address             sql.NullString  `json:"address"`
	SocialLinks         json.RawMessage `json:"social_links"`
	LogoURL             sql.NullString  `json:"logo_url"`
	FaviconURL          sql.NullString  `json:"favicon_url"`
	MaintenanceEnabled  bool            `json:"maintenance_enabled"`
	MaintenanceMessage  sql.NullString  `json:"maintenance_message"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

// GetSiteSettings busca as configurações do site no banco de dados.
func GetSiteSettings(db *sql.DB) (*SiteSettings, error) {
	var s SiteSettings
	query := `SELECT 
        company_name, slogan, description, phone_numbers, whatsapp, 
        institutional_email, address, social_links, logo_url, favicon_url, 
        maintenance_enabled, maintenance_message, updated_at 
    FROM site_settings WHERE id = 1`

	err := db.QueryRow(query).Scan(
		&s.CompanyName, &s.Slogan, &s.Description, &s.PhoneNumbers, &s.Whatsapp,
		&s.InstitutionalEmail, &s.Address, &s.SocialLinks, &s.LogoURL, &s.FaviconURL,
		&s.MaintenanceEnabled, &s.MaintenanceMessage, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// UpdateSiteSettings atualiza as configurações do site.
func UpdateSiteSettings(db *sql.DB, s *SiteSettings) error {
	query := `UPDATE site_settings SET
		company_name = ?, slogan = ?, description = ?, phone_numbers = ?, whatsapp = ?,
		institutional_email = ?, address = ?, social_links = ?, logo_url = ?, favicon_url = ?
	WHERE id = 1`

	_, err := db.Exec(query,
		s.CompanyName, s.Slogan, s.Description, s.PhoneNumbers, s.Whatsapp,
		s.InstitutionalEmail, s.Address, s.SocialLinks, s.LogoURL, s.FaviconURL,
	)
	return err
}

// UpdateMaintenanceMode atualiza especificamente o status do modo de manutenção.
func UpdateMaintenanceMode(db *sql.DB, enabled bool, message string) error {
	query := `UPDATE site_settings SET maintenance_enabled = ?, maintenance_message = ? WHERE id = 1`
	_, err := db.Exec(query, enabled, message)
	return err
}
