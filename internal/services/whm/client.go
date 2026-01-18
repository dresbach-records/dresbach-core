package whm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// Client é um cliente para a API do WHM.
type Client struct {
	Host       string
	Username   string
	ApiToken   string
	HttpClient *http.Client
}

// NewClient cria um novo cliente WHM a partir de variáveis de ambiente.
func NewClient() (*Client, error) {
	host := os.Getenv("WHM_HOST")
	user := os.Getenv("WHM_USER")
	token := os.Getenv("WHM_API_TOKEN")

	if host == "" || user == "" || token == "" {
		return nil, fmt.Errorf("as variáveis de ambiente WHM_HOST, WHM_USER, e WHM_API_TOKEN são obrigatórias")
	}

	return &Client{
		Host:       host,
		Username:   user,
		ApiToken:   token,
		HttpClient: &http.Client{},
	}, nil
}

// CreateAccountResponse define a estrutura da resposta da API do WHM.
// Apenas os campos que nos interessam estão mapeados.
type CreateAccountResponse struct {
	Metadata struct {
		Command string `json:"command"`
		Reason  string `json:"reason"`
		Result  int    `json:"result"`
	} `json:"metadata"`
}

// CreateAccount cria uma nova conta de hospedagem cPanel.
func (c *Client) CreateAccount(domain, userName, password, planName string) (*CreateAccountResponse, error) {
	// Constrói a URL da API
	apiURL := fmt.Sprintf("https://%s:2087/json-api/createacct?api.version=1", c.Host)

	// Adiciona os parâmetros da conta na URL
	params := url.Values{}
	params.Add("domain", domain)
	params.Add("username", userName)
	params.Add("password", password)
	params.Add("plan", planName)

	req, err := http.NewRequest("GET", apiURL+"&"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Adiciona o cabeçalho de autorização
	authHeader := fmt.Sprintf("WHM %s:%s", c.Username, c.ApiToken)
	req.Header.Add("Authorization", authHeader)

	// Executa a requisição
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição para o WHM: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do WHM: %w", err)
	}

	// Decodifica a resposta JSON
	var apiResp CreateAccountResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON do WHM: %w. Resposta: %s", err, string(body))
	}

	// Verifica se a operação foi bem-sucedida (Result == 1)
	if apiResp.Metadata.Result != 1 {
		return &apiResp, fmt.Errorf("falha ao criar conta no WHM: %s", apiResp.Metadata.Reason)
	}

	return &apiResp, nil
}
