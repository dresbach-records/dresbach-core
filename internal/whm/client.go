package whm

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// Client é o cliente da API para interagir com o WHM.
type Client struct {
	Host       string
	Username   string
	ApiToken   string
	HttpClient *http.Client
}

// NewClient cria e retorna um novo cliente WHM.
// As configurações são lidas de variáveis de ambiente.
func NewClient() (*Client, error) {
	host := os.Getenv("WHM_HOST")
	user := os.Getenv("WHM_USER")
	token := os.Getenv("WHM_API_TOKEN")

	if host == "" || user == "" || token == "" {
		return nil, fmt.Errorf("as variáveis de ambiente WHM_HOST, WHM_USER, e WHM_API_TOKEN devem ser definidas")
	}

	return &Client{
		Host:     host,
		Username: user,
		ApiToken: token,
		HttpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}, nil
}

// CreateAccount cria uma nova conta de hospedagem no cPanel.
// Esta é uma função de exemplo para demonstração.
func (c *Client) CreateAccount(domain, username, password, email, plan string) error {
	// A URL para a chamada de criação de conta na API do WHM1
	// Exemplo: https://your-whm-host:2087/json-api/createacct?api.version=1&...
	url := fmt.Sprintf("https://%s:2087/json-api/createacct?api.version=1&domain=%s&username=%s&password=%s&contactemail=%s&plan=%s",
		c.Host, domain, username, password, email, plan)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// A autenticação é feita via um header 'Authorization'
	authHeader := fmt.Sprintf("whm %s:%s", c.Username, c.ApiToken)
	req.Header.Set("Authorization", authHeader)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao executar requisição para o WHM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Aqui você adicionaria um tratamento de erro mais robusto,
		// lendo o corpo da resposta para obter a mensagem de erro exata do WHM.
		return fmt.Errorf("o WHM respondeu com status não OK: %s", resp.Status)
	}

	// Decodificar a resposta JSON do WHM para verificar se a operação foi bem-sucedida.
	// ... (lógica de tratamento da resposta)

	fmt.Println("Conta criada com sucesso para o domínio:", domain)
	return nil
}
