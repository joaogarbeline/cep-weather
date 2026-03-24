package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joaogarbeline/cep-weather/internal/handler"
	"github.com/joaogarbeline/cep-weather/internal/service"
)

// mockWeatherService is a test double for WeatherServiceInterface
type mockWeatherService struct {
	result *service.TemperatureResponse
	err    error
}

func (m *mockWeatherService) GetTemperaturesByCEP(cep string) (*service.TemperatureResponse, error) {
	return m.result, m.err
}

func TestGetWeatherByCEP_Success(t *testing.T) {
	svc := &mockWeatherService{
		result: &service.TemperatureResponse{TempC: 25, TempF: 77, TempK: 298},
	}
	h := handler.NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/01310100", nil)
	w := httptest.NewRecorder()
	h.GetWeatherByCEP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var body service.TemperatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body.TempC != 25 {
		t.Errorf("expected TempC=25, got %.2f", body.TempC)
	}
}

func TestGetWeatherByCEP_InvalidZipcode(t *testing.T) {
	svc := &mockWeatherService{err: service.ErrInvalidZipcode}
	h := handler.NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/123", nil)
	w := httptest.NewRecorder()
	h.GetWeatherByCEP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestGetWeatherByCEP_NotFound(t *testing.T) {
	svc := &mockWeatherService{err: service.ErrZipcodeNotFound}
	h := handler.NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/99999999", nil)
	w := httptest.NewRecorder()
	h.GetWeatherByCEP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetWeatherByCEP_MethodNotAllowed(t *testing.T) {
	svc := &mockWeatherService{}
	h := handler.NewHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/01310100", nil)
	w := httptest.NewRecorder()
	h.GetWeatherByCEP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}
