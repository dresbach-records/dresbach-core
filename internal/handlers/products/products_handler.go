package products

import (
	"encoding/json"
	"log"
	"net/http"

	"hosting-backend/internal/provisioning"
)

// GetVpsProductsHandler busca e retorna os planos de VPS disponíveis para revenda.
func GetVpsProductsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Inicializa o provisioner da Hostinger
		hostinger, err := provisioning.NewHostingerProvisioner()
		if err != nil {
			log.Printf("Erro ao inicializar o provisioner da Hostinger: %v", err)
			http.Error(w, "Erro interno do servidor - falha na configuração do provedor", http.StatusInternalServerError)
			return
		}

		// Busca os planos de VPS do catálogo da Hostinger
		vpsPlans, err := hostinger.GetVpsCatalog()
		if err != nil {
			log.Printf("Erro ao buscar catálogo de VPS da Hostinger: %v", err)
			http.Error(w, "Não foi possível obter a lista de produtos no momento", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(vpsPlans); err != nil {
			http.Error(w, "Erro ao formatar a lista de produtos", http.StatusInternalServerError)
		}
	}
}
