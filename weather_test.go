package service_test

import (
	"errors"
	"testing"

	"github.com/joaogarbeline/cep-weather/internal/client"
	"github.com/joaogarbeline/cep-weather/internal/service"
)

// mockCEPClient is a test double for CEPClient
type mockCEPClient struct {
	response *client.ViaCEPResponse
	err      error
}

func (m *mockCEPClient) GetLocationByCEP(cep string) (*client.ViaCEPResponse, error) {
	return m.response, m.err
}

// mockWeatherClient is a test double for WeatherClient
type mockWeatherClient struct {
	tempC float64
	err   error
}

func (m *mockWeatherClient) GetCurrentTemperature(city string) (float64, error) {
	return m.tempC, m.err
}

func TestGetTemperaturesByCEP_Success(t *testing.T) {
	cepClient := &mockCEPClient{
		response: &client.ViaCEPResponse{
			CEP:        "01310100",
			Localidade: "São Paulo",
			UF:         "SP",
		},
	}
	weatherClient := &mockWeatherClient{tempC: 28.5}

	svc := service.NewWeatherServiceWithClients(cepClient, weatherClient)

	result, err := svc.GetTemperaturesByCEP("01310100")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.TempC != 28.5 {
		t.Errorf("expected TempC=28.5, got %.2f", result.TempC)
	}
	if result.TempF != 83.3 {
		t.Errorf("expected TempF=83.3, got %.2f", result.TempF)
	}
	if result.TempK != 301.5 {
		t.Errorf("expected TempK=301.5, got %.2f", result.TempK)
	}
}

func TestGetTemperaturesByCEP_InvalidCEP(t *testing.T) {
	tests := []struct {
		name string
		cep  string
	}{
		{"too short", "0131010"},
		{"too long", "013101000"},
		{"has letters", "0131010A"},
		{"empty", ""},
		{"with dash", "01310-100"},
	}

	svc := service.NewWeatherServiceWithClients(&mockCEPClient{}, &mockWeatherClient{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.GetTemperaturesByCEP(tt.cep)
			if !errors.Is(err, service.ErrInvalidZipcode) {
				t.Errorf("expected ErrInvalidZipcode, got %v", err)
			}
		})
	}
}

func TestGetTemperaturesByCEP_NotFound(t *testing.T) {
	cepClient := &mockCEPClient{response: nil}
	svc := service.NewWeatherServiceWithClients(cepClient, &mockWeatherClient{})

	_, err := svc.GetTemperaturesByCEP("99999999")
	if !errors.Is(err, service.ErrZipcodeNotFound) {
		t.Errorf("expected ErrZipcodeNotFound, got %v", err)
	}
}

func TestGetTemperaturesByCEP_CEPClientError(t *testing.T) {
	cepClient := &mockCEPClient{err: errors.New("network error")}
	svc := service.NewWeatherServiceWithClients(cepClient, &mockWeatherClient{})

	_, err := svc.GetTemperaturesByCEP("01310100")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetTemperaturesByCEP_WeatherClientError(t *testing.T) {
	cepClient := &mockCEPClient{
		response: &client.ViaCEPResponse{Localidade: "São Paulo"},
	}
	weatherClient := &mockWeatherClient{err: errors.New("API error")}
	svc := service.NewWeatherServiceWithClients(cepClient, weatherClient)

	_, err := svc.GetTemperaturesByCEP("01310100")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestTemperatureConversions tests the conversion formulas directly
func TestTemperatureConversions(t *testing.T) {
	tests := []struct {
		tempC    float64
		expF     float64
		expK     float64
	}{
		{0, 32, 273},
		{100, 212, 373},
		{-40, -40, 233},
		{28.5, 83.3, 301.5},
	}

	for _, tt := range tests {
		cepClient := &mockCEPClient{
			response: &client.ViaCEPResponse{Localidade: "TestCity"},
		}
		weatherClient := &mockWeatherClient{tempC: tt.tempC}
		svc := service.NewWeatherServiceWithClients(cepClient, weatherClient)

		result, err := svc.GetTemperaturesByCEP("01310100")
		if err != nil {
			t.Fatalf("unexpected error for tempC=%.1f: %v", tt.tempC, err)
		}

		if result.TempC != tt.tempC {
			t.Errorf("TempC: expected %.2f, got %.2f", tt.tempC, result.TempC)
		}
		if result.TempF != tt.expF {
			t.Errorf("TempF (for %.1f°C): expected %.2f, got %.2f", tt.tempC, tt.expF, result.TempF)
		}
		if result.TempK != tt.expK {
			t.Errorf("TempK (for %.1f°C): expected %.2f, got %.2f", tt.tempC, tt.expK, result.TempK)
		}
	}
}
