package admin

import (
	"encoding/json"
	"net/http"
)

// DashboardData representa os dados a serem exibidos no dashboard do admin.

// DashboardData representa os dados a serem exibidos no dashboard do admin.
type DashboardData struct {
	Message string `json:"message"`
}

// DashboardHandler processa as solicitações para o endpoint /admin/dashboard.
// DashboardHandler processa as solicitações para o endpoint /admin/dashboard.
func DashboardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DashboardData{Message: "Bem-vindo ao dashboard do admin!"})
	}
}
