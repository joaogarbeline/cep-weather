package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joaogarbeline/cep-weather/internal/client"
	"github.com/joaogarbeline/cep-weather/internal/handler"
	"github.com/joaogarbeline/cep-weather/internal/service"
)

func setupTestServers(t *testing.T, cepBody string, cepStatus int, weatherBody string, weatherStatus int) (*httptest.Server, *httptest.Server) {
	t.Helper()

	cepServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(cepStatus)
		fmt.Fprint(w, cepBody)
	}))

	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(weatherStatus)
		fmt.Fprint(w, weatherBody)
	}))

	return cepServer, weatherServer
}

func TestIntegration_Success(t *testing.T) {
	cepJSON := `{"cep":"01310-100","logradouro":"Avenida Paulista","bairro":"Bela Vista","localidade":"São Paulo","uf":"SP","erro":false}`
	weatherJSON := `{"current":{"temp_c":28.5}}`

	cepServer, weatherServer := setupTestServers(t, cepJSON, 200, weatherJSON, 200)
	defer cepServer.Close()
	defer weatherServer.Close()

	cepClient := client.NewViaCEPClientWithBaseURL(cepServer.URL)
	weatherClient := client.NewWeatherAPIClientWithBaseURL("test-key", weatherServer.URL)
	svc := service.NewWeatherServiceWithClients(cepClient, weatherClient)
	h := handler.NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/01310100", nil)
	w := httptest.NewRecorder()
	h.GetWeatherByCEP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp service.TemperatureResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if resp.TempC != 28.5 {
		t.Errorf("expected TempC=28.5, got %.2f", resp.TempC)
	}
	if resp.TempF != 83.3 {
		t.Errorf("expected TempF=83.3, got %.2f", resp.TempF)
	}
	if resp.TempK != 301.5 {
		t.Errorf("expected TempK=301.5, got %.2f", resp.TempK)
	}
}

func TestIntegration_InvalidCEP(t *testing.T) {
	cepServer, weatherServer := setupTestServers(t, "", 200, "", 200)
	defer cepServer.Close()
	defer weatherServer.Close()

	cepClient := client.NewViaCEPClientWithBaseURL(cepServer.URL)
	weatherClient := client.NewWeatherAPIClientWithBaseURL("test-key", weatherServer.URL)
	svc := service.NewWeatherServiceWithClients(cepClient, weatherClient)
	h := handler.NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/123abc", nil)
	w := httptest.NewRecorder()
	h.GetWeatherByCEP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestIntegration_CEPNotFound(t *testing.T) {
	cepJSON := `{"erro":true}`
	cepServer, weatherServer := setupTestServers(t, cepJSON, 200, "", 200)
	defer cepServer.Close()
	defer weatherServer.Close()

	cepClient := client.NewViaCEPClientWithBaseURL(cepServer.URL)
	weatherClient := client.NewWeatherAPIClientWithBaseURL("test-key", weatherServer.URL)
	svc := service.NewWeatherServiceWithClients(cepClient, weatherClient)
	h := handler.NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/99999999", nil)
	w := httptest.NewRecorder()
	h.GetWeatherByCEP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
