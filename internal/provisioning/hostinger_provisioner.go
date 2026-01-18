package provisioning

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const hostingerAPIBaseURL = "https://developers.hostinger.com/api"

// HostingerProvisioner implementa a comunicação com a API da Hostinger.
type HostingerProvisioner struct {
	apiToken   string
	httpClient *http.Client
}

// --- Estruturas Comuns e de Domínio ---
type HostingerAvailabilityRequest struct {
	Domains []string `json:"domains"`
}

type HostingerAvailabilityResponse struct {
	Data map[string]struct {
		Status string `json:"status"`
	} `json:"data"`
}

type WhoisProfile struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Country  string `json:"country"`
	State    string `json:"state"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Postcode string `json:"postcode"`
	Phone    string `json:"phone"`
	TaxID    string `json:"tax_id,omitempty"`
}

type WhoisProfileResponse struct {
	ID int `json:"id"`
}

type DomainPurchaseRequest struct {
	Domain         string `json:"domain"`
	WhoisProfileID int    `json:"whois_profile_id"`
	Period         int    `json:"period"`
}

type DomainPurchaseResponse struct {
	OrderID string `json:"order_id"`
}

// --- Estruturas de Catálogo e VPS ---

type CatalogItemPrice struct {
	Period int `json:"period"`
	PeriodType string `json:"period_type"`
	Price int `json:"price"`
	Currency string `json:"currency"`
}

type CatalogItem struct {
	ID           string             `json:"id"`
	ProductGroup string             `json:"product_group"`
	Name         string             `json:"name"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	Pricing      []CatalogItemPrice `json:"pricing"`
}

type CatalogResponse struct {
	Data []CatalogItem `json:"data"`
}

type VpsPlan struct {
	PlanID      string `json:"plan_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Currency    string `json:"currency"`
}

type VpsPurchaseRequest struct {
	PlanID   string `json:"plan_id"`
	Location string `json:"location"`
	Template string `json:"template"`
	Hostname string `json:"hostname"`
	Password string `json:"password"`
	SshKeyID string `json:"ssh_key_id,omitempty"`
}

type VpsPurchaseResponse struct {
	VmID string `json:"id"`
}

// --- Métodos do Provisioner ---

func NewHostingerProvisioner() (*HostingerProvisioner, error) {
	apiToken := os.Getenv("HOSTINGER_API_TOKEN")
	if apiToken == "" {
		return nil, fmt.Errorf("variável de ambiente HOSTINGER_API_TOKEN não definida")
	}
	return &HostingerProvisioner{apiToken: apiToken, httpClient: &http.Client{}}, nil
}

func (p *HostingerProvisioner) makeHostingerAPIRequest(method, path string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", hostingerAPIBaseURL, path)
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("falha ao serializar corpo: %w", err)
		}
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.apiToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao enviar requisição: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler resposta: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API retornou status não-OK: %d - %s", resp.StatusCode, string(respBody))
	}
	log.Printf("[Hostinger Provisioner] Resposta da API para '%s %s' recebida.", method, path)
	return respBody, nil
}

func (p *HostingerProvisioner) CheckDomainAvailability(domains []string) (*HostingerAvailabilityResponse, error) { /* ... */ return nil, nil }
func (p *HostingerProvisioner) CreateWhoisProfile(profile WhoisProfile) (*WhoisProfileResponse, error) { /* ... */ return nil, nil }
func (p *HostingerProvisioner) RegisterDomain(domain string, period int, whoisID int) (*DomainPurchaseResponse, error) { /* ... */ return nil, nil }

func (p *HostingerProvisioner) GetVpsCatalog() ([]VpsPlan, error) {
	log.Println("[Hostinger Provisioner] Buscando catálogo de produtos VPS...")
	respBody, err := p.makeHostingerAPIRequest("GET", "/billing/v1/catalog", nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar catálogo: %w", err)
	}
	var catalogResponse CatalogResponse
	if err := json.Unmarshal(respBody, &catalogResponse); err != nil {
		return nil, fmt.Errorf("falha ao decodificar catálogo: %w", err)
	}
	var vpsPlans []VpsPlan
	for _, item := range catalogResponse.Data {
		if item.ProductGroup == "vps" && strings.HasPrefix(item.ID, "vps_") {
			var monthlyPrice int
			var currency string
			for _, priceInfo := range item.Pricing {
				if priceInfo.Period == 1 && priceInfo.PeriodType == "month" {
					monthlyPrice = priceInfo.Price
					currency = priceInfo.Currency
					break
				}
			}
			if monthlyPrice > 0 {
				vpsPlans = append(vpsPlans, VpsPlan{PlanID: item.ID, Name: item.Title, Description: item.Description, Price: monthlyPrice, Currency: currency})
			}
		}
	}
	log.Printf("[Hostinger Provisioner] Encontrados %d planos de VPS.", len(vpsPlans))
	return vpsPlans, nil
}

// PurchaseVps provisiona um novo servidor virtual privado.
func (p *HostingerProvisioner) PurchaseVps(req VpsPurchaseRequest) (*VpsPurchaseResponse, error) {
	log.Printf("[Hostinger Provisioner] Iniciando provisionamento de VPS com plano '%s' em '%s'", req.PlanID, req.Location)

	respBody, err := p.makeHostingerAPIRequest("POST", "/vps/v1/virtual-machines", req)
	if err != nil {
		return nil, fmt.Errorf("falha ao provisionar VPS na Hostinger: %w", err)
	}

	var purchaseResponse VpsPurchaseResponse
	if err := json.Unmarshal(respBody, &purchaseResponse); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta do provisionamento de VPS: %w", err)
	}

	log.Printf("[Hostinger Provisioner] VPS provisionado com sucesso. VM ID: %s", purchaseResponse.VmID)
	return &purchaseResponse, nil
}
