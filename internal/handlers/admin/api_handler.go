package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"hosting-backend/internal/models"
	"github.com/go-chi/chi/v5"
)

// GetApiCredentialsHandler lista as credenciais de API (sem as chaves).
// Rota: GET /admin/apis
func GetApiCredentialsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creds, err := models.GetApiCredentialsForAdmin(db)
		if err != nil {
			http.Error(w, "Erro ao buscar credenciais de API", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(creds)
	}
}

// UpdateApiCredentialHandler atualiza uma credencial de API.
// Rota: PUT /admin/apis/{provider}
func UpdateApiCredentialHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")

		var payload struct {
			ApiKey string `json:"api_key"`
			Status string `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		if err := models.UpsertApiCredential(db, models.ApiProvider(provider), payload.ApiKey, payload.Status); err != nil {
			http.Error(w, "Erro ao salvar a credencial de API", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// TestApiConnectionHandler testa a conexão com uma API externa.
// Rota: POST /admin/apis/{provider}/test
func TestApiConnectionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := models.ApiProvider(chi.URLParam(r, "provider"))

		// apiKey, err := models.GetDecryptedApiKey(db, provider)
		// if err != nil {
		// 	models.UpdateApiTestResult(db, provider, err)
		// 	http.Error(w, "Chave de API não encontrada ou inativa", http.StatusNotFound)
		// 	return
		// }

		// TODO: Implementar a lógica de teste real para cada provedor.
		// Por exemplo, para o Stripe, fazer uma chamada simples como `stripe.Balance.Get()`
		// Para o WHM, uma chamada como `whm.GetVersion()`.
		// Por enquanto, vamos simular um sucesso.
		var testErr error = nil // Simulação

		if err := models.UpdateApiTestResult(db, provider, testErr); err != nil {
			http.Error(w, "Erro ao salvar resultado do teste", http.StatusInternalServerError)
			return
		}

		if testErr != nil {
			w.WriteHeader(http.StatusBadRequest) // 400 se o teste falhou
			json.NewEncoder(w).Encode(map[string]string{"status": "failed", "error": testErr.Error()})
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		}
	}
}
