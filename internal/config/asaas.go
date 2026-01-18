package config

import (
	"os"
)

var (
	AsaasAPIKey string
	AsaasBaseURL string
)

// InitAsaas inicializa a configuração da API do Asaas.
func InitAsaas() {
	AsaasAPIKey = os.Getenv("ASAAS_API_KEY")
	AsaasBaseURL = os.Getenv("ASAAS_BASE_URL")
}
