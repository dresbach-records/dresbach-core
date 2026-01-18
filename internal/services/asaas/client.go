package asaas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"hosting-backend/internal/config"
)

// AsaasClient é o cliente para interagir com a API do Asaas.
type AsaasClient struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

// NewAsaasClient cria uma nova instância do cliente Asaas.
func NewAsaasClient() *AsaasClient {
	return &AsaasClient{
		BaseURL: config.AsaasBaseURL,
		APIKey:  config.AsaasAPIKey,
		Client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// doRequest é uma função auxiliar para fazer requisições à API do Asaas.
func (a *AsaasClient) doRequest(method, path string, body any) ([]byte, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("falha ao encodar o corpo da requisição: %w", err)
		}
	}

	req, err := http.NewRequest(method, a.BaseURL+path, &buf)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar a requisição http: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+a.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar a requisição http: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler o corpo da resposta: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API Asaas retornou status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// --- Structs de Requisição ---

// CustomerRequest representa os dados para criar um cliente.
type CustomerRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	CpfCnpj string `json:"cpfCnpj"`
	Phone   string `json:"phone,omitempty"`
}

// PaymentRequest representa os dados para criar uma cobrança.
type PaymentRequest struct {
	Customer    string  `json:"customer"`
	BillingType string  `json:"billingType"`
	Value       float64 `json:"value"`
	DueDate     string  `json:"dueDate"`
	Description string  `json:"description"`
}

// --- Structs de Resposta ---

// CreateCustomerResponse é a resposta da API ao criar um cliente.
type CreateCustomerResponse struct {
	ID string `json:"id"`
	// Adicione outros campos se necessário
}

// CreatePaymentResponse é a resposta da API ao criar uma cobrança.
type CreatePaymentResponse struct {
	ID         string  `json:"id"`
	Status     string  `json:"status"`
	InvoiceURL string  `json:"invoiceUrl"`
}

// --- Funções da API ---

// CreateCustomer cria um novo cliente no Asaas.
func (a *AsaasClient) CreateCustomer(c CustomerRequest) (*CreateCustomerResponse, error) {
	resp, err := a.doRequest("POST", "/v3/customers", c)
	if err != nil {
		return nil, err
	}

	var customerResponse CreateCustomerResponse
	if err := json.Unmarshal(resp, &customerResponse); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta de CreateCustomer: %w", err)
	}

	return &customerResponse, nil
}

// CreatePayment cria uma nova cobrança no Asaas.
func (a *AsaasClient) CreatePayment(p PaymentRequest) (*CreatePaymentResponse, error) {
	resp, err := a.doRequest("POST", "/v3/payments", p)
	if err != nil {
		return nil, err
	}

	var paymentResponse CreatePaymentResponse
	if err := json.Unmarshal(resp, &paymentResponse); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta de CreatePayment: %w", err)
	}

	return &paymentResponse, nil
}
