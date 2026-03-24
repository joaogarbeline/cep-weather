package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const viaCEPBaseURL = "https://viacep.com.br/ws"

// ViaCEPResponse represents the response from ViaCEP API
type ViaCEPResponse struct {
	CEP        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Erro       bool   `json:"erro"`
}

// ViaCEPClient handles requests to the ViaCEP API
type ViaCEPClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewViaCEPClient creates a new ViaCEP client
func NewViaCEPClient() *ViaCEPClient {
	return &ViaCEPClient{
		httpClient: &http.Client{},
		baseURL:    viaCEPBaseURL,
	}
}

// NewViaCEPClientWithBaseURL creates a new ViaCEP client with a custom base URL (for testing)
func NewViaCEPClientWithBaseURL(baseURL string) *ViaCEPClient {
	return &ViaCEPClient{
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

// GetLocationByCEP fetches location data for the given CEP
func (c *ViaCEPClient) GetLocationByCEP(cep string) (*ViaCEPResponse, error) {
	url := fmt.Sprintf("%s/%s/json/", c.baseURL, cep)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CEP data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("invalid CEP format")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Erro {
		return nil, nil // CEP not found
	}

	return &result, nil
}
