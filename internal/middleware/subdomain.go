
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// SiteSettings armazena as configurações globais do site.
// Em um cenário real, isso viria de um banco de dados ou de um sistema de configuração.
var SiteSettings = struct {
	SiteMode      string
	ThemeMode     string
	CountdownEnd  string
}{
	SiteMode:     "online", // "online", "countdown", "maintenance"
	ThemeMode:    "default",  // "default", "christmas"
	CountdownEnd: "2025-12-31T23:59:59-03:00",
}

// SubdomainMiddleware analisa o subdomínio e aplica as regras de negócio.
func SubdomainMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		ctx := context.WithValue(r.Context(), "host", host)

		// Adiciona o header Vary para indicar que a resposta depende do Host.
		w.Header().Add("Vary", "Host")

		switch host {
		case "www.dresbachhosting.com.br":
			if SiteSettings.SiteMode == "countdown" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{
					"mode":          SiteSettings.SiteMode,
					"countdown_end": SiteSettings.CountdownEnd,
				})
				return
			}
			// Adicione aqui a lógica para outros modos, como "maintenance".

		case "area-do-cliente.dresbachhosting.com.br":
			// A área do cliente pode ter um tema, mas não é afetada pelo countdown.
			// A informação do tema pode ser passada para o frontend de várias formas.
			// Uma delas é através de um header customizado.
			if SiteSettings.ThemeMode == "christmas" {
				w.Header().Set("X-Theme-Mode", "christmas")
			}

		case "checkout.dresbachhosting.com.br":
			// Nenhuma regra especial de tema ou modo se aplica aqui.
			// O request segue o fluxo normal.
			break

		case "admin.dresbachhosting.com.br":
			// As rotas de admin já são protegidas por outros middlewares (auth, rbac).
			// Nenhuma lógica adicional de subdomínio é necessária aqui.
			break
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORSMiddleware adiciona os headers de CORS necessários.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := []string{
			"https://www.dresbachhosting.com.br",
			"https://area-do-cliente.dresbachhosting.com.br",
			"https://checkout.dresbachhosting.com.br",
			"https://admin.dresbachhosting.com.br",
		}
		origin := r.Header.Get("Origin")

		for _, allowed := range allowedOrigins {
			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Theme-Mode")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Se for uma requisição OPTIONS, apenas retorne os headers.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// UpdateSettingsHandler permite que o admin altere as configurações do site.
// Esta função deve ser protegida e acessível apenas por administradores.
func UpdateSettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var newSettings struct {
		SiteMode     string `json:"site_mode"`
		ThemeMode    string `json:"theme_mode"`
		CountdownEnd string `json:"countdown_end"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	// Validação básica
	if newSettings.SiteMode != "" {
		SiteSettings.SiteMode = newSettings.SiteMode
	}
	if newSettings.ThemeMode != "" {
		SiteSettings.ThemeMode = newSettings.ThemeMode
	}
	if newSettings.CountdownEnd != "" {
		// Validar o formato da data seria uma boa prática
		_, err := time.Parse(time.RFC3339, newSettings.CountdownEnd)
		if err != nil {
			http.Error(w, "Formato de data inválido para countdown_end", http.StatusBadRequest)
			return
		}
		SiteSettings.CountdownEnd = newSettings.CountdownEnd
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SiteSettings)
}
