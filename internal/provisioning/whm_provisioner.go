package provisioning

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// WhmProvisioner implementa a comunicação com a API do WHM/cPanel.
type WhmProvisioner struct {
	whmHost    string
	whmUser    string
	apiToken   string
	httpClient *http.Client
}

// --- Estruturas da API WHM ---

type CreateAccountResponse struct {
	Metadata struct {
		Result  int    `json:"result"`
		Reason  string `json:"reason"`
		Command string `json:"command"`
	} `json:"metadata"`
}

// --- Métodos do Provisioner ---

// NewWhmProvisioner cria uma nova instância do WhmProvisioner.
func NewWhmProvisioner() (*WhmProvisioner, error) {
	whmHost := os.Getenv("WHM_HOST")
	whmUser := os.Getenv("WHM_USER")
	apiToken := os.Getenv("WHM_API_TOKEN")

	if whmHost == "" || whmUser == "" || apiToken == "" {
		return nil, fmt.Errorf("as variáveis de ambiente WHM_HOST, WHM_USER e WHM_API_TOKEN devem estar definidas")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Cuidado em produção
	}

	return &WhmProvisioner{
		whmHost:    whmHost,
		whmUser:    whmUser,
		apiToken:   apiToken,
		httpClient: &http.Client{Transport: tr},
	}, nil
}

// makeWhmAPIRequest constrói e executa uma chamada para a API do WHM.
func (p *WhmProvisioner) makeWhmAPIRequest(function string, params url.Values) ([]byte, error) {
	apiURL := fmt.Sprintf("https://%s:2087/json-api/%s?api.version=1&%s", p.whmHost, function, params.Encode())

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição para o WHM: %w", err)
	}

	authHeader := fmt.Sprintf("whm %s:%s", p.whmUser, p.apiToken)
	req.Header.Set("Authorization", authHeader)

	log.Printf("[WHM Provisioner] Executando chamada para a função: %s", function)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao enviar requisição para o WHM: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler resposta do WHM: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API do WHM retornou status não-OK: %d - %s", resp.StatusCode, string(respBody))
	}

	log.Printf("[WHM Provisioner] Resposta da API WHM recebida com sucesso para a função: %s", function)
	return respBody, nil
}

// CreateAccount cria uma nova conta de hospedagem no servidor cPanel via API do WHM.
func (p *WhmProvisioner) CreateAccount(username, domain, plan, password, email string) (*CreateAccountResponse, error) {
	params := url.Values{}
	params.Set("username", username)
	params.Set("domain", domain)
	params.Set("plan", plan)
	params.Set("password", password)
	params.Set("contactemail", email)

	log.Printf("[WHM Provisioner] Enviando solicitação para criar conta '%s' com domínio '%s' no plano '%s'", username, domain, plan)

	respBody, err := p.makeWhmAPIRequest("createacct", params)
	if err != nil {
		return nil, fmt.Errorf("falha na chamada da API para criar conta: %w", err)
	}

	var response CreateAccountResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta da criação da conta: %w - Resposta: %s", err, string(respBody))
	}

	// A API do WHM retorna 'result: 1' em caso de sucesso.
	if response.Metadata.Result != 1 {
		log.Printf("[WHM Provisioner] Erro ao criar conta '%s': %s", username, response.Metadata.Reason)
		return nil, fmt.Errorf("WHM API retornou um erro: %s", response.Metadata.Reason)
	}

	log.Printf("[WHM Provisioner] Conta '%s' criada com sucesso no cPanel.", username)
	return &response, nil
}
