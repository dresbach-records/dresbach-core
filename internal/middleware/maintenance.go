package middleware

import (
	"database/sql"
	"net/http"
	"sync"
	"time"

	"hosting-backend/internal/models"
)

var (
	maintenanceStatus   bool
	maintenanceMessage  string
	lastChecked         time.Time
	mu                  sync.RWMutex
	checkInterval       = 1 * time.Minute // Faz cache do status por 1 minuto para não sobrecarregar o DB
)

// checkAndUpdateMaintenanceStatus verifica o status de manutenção no DB se o cache expirou.
func checkAndUpdateMaintenanceStatus(db *sql.DB) {
	mu.RLock()
	isCacheValid := time.Since(lastChecked) < checkInterval
	mu.RUnlock()

	if isCacheValid {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Double-check para evitar race conditions
	if time.Since(lastChecked) < checkInterval {
		return
	}

	settings, err := models.GetSiteSettings(db)
	if err == nil {
		maintenanceStatus = settings.MaintenanceEnabled
		maintenanceMessage = settings.MaintenanceMessage.String
	} else {
		// Em caso de erro, desativa a manutenção como medida de segurança para não travar o site
		maintenanceStatus = false
	}
	lastChecked = time.Now()
}

// MaintenanceMiddleware verifica se o site está em modo manutenção.
func MaintenanceMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			checkAndUpdateMaintenanceStatus(db)

			mu.RLock()
			isEnabled := maintenanceStatus
			message := maintenanceMessage
			mu.RUnlock()

			if isEnabled {
				// TODO: Implementar lógica para verificar se o usuário é um administrador logado
				// if session.IsAdmin(r.Context()) { next.ServeHTTP(w, r); return }

				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusServiceUnavailable)
				// TODO: Idealmente, servir um template HTML aqui
				w.Write([]byte(message))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
