package domain

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"hosting-backend/internal/provisioning"
)

// CheckDomainRequest é a estrutura para o corpo da requisição de verificação de domínio.
type CheckDomainRequest struct {
	Domain string `json:"domain"`
}

// CheckDomainResponse é a estrutura para a resposta da verificação.
type CheckDomainResponse struct {
	Domain      string `json:"domain"`
	Available   bool   `json:"available"`
	Message     string `json:"message"`
}

// CheckDomainHandler lida com a verificação de disponibilidade de domínios via API da Hostinger.
func CheckDomainHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CheckDomainRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		domain := strings.ToLower(strings.TrimSpace(req.Domain))
		if domain == "" {
			http.Error(w, "O campo 'domain' não pode ser vazio", http.StatusBadRequest)
			return
		}

		// Inicializa o provisioner da Hostinger
		hostinger, err := provisioning.NewHostingerProvisioner()
		if err != nil {
			log.Printf("Erro ao inicializar o provisioner da Hostinger: %v", err)
			http.Error(w, "Erro interno do servidor - falha na configuração", http.StatusInternalServerError)
			return
		}

		// Chama a API da Hostinger para verificar a disponibilidade
		availability, err := hostinger.CheckDomainAvailability([]string{domain})
		if err != nil {
			log.Printf("Erro ao verificar disponibilidade do domínio '%s' na Hostinger: %v", domain, err)
			http.Error(w, "Erro ao consultar serviço de domínios", http.StatusServiceUnavailable)
			return
		}

		// Prepara a resposta
		response := CheckDomainResponse{Domain: domain}
		if result, ok := availability.Data[domain]; ok && result.Status == "available" {
			response.Available = true
			response.Message = fmt.Sprintf("Parabéns! O domínio %s está disponível para registro.", domain)
		} else {
			response.Available = false
			response.Message = fmt.Sprintf("Que pena. O domínio %s não está disponível.", domain)
		}
        
        w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
