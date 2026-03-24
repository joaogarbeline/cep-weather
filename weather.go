package service

import (
	"fmt"
	"math"

	"github.com/joaogarbeline/cep-weather/internal/client"
)

// TemperatureResponse holds temperatures in all three units
type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// CEPClient interface for fetching location by CEP
type CEPClient interface {
	GetLocationByCEP(cep string) (*client.ViaCEPResponse, error)
}

// WeatherClient interface for fetching temperature
type WeatherClient interface {
	GetCurrentTemperature(city string) (float64, error)
}

// WeatherService handles the business logic
type WeatherService struct {
	cepClient     CEPClient
	weatherClient WeatherClient
}

// NewWeatherService creates a new WeatherService with default clients
func NewWeatherService(weatherAPIKey string) *WeatherService {
	return &WeatherService{
		cepClient:     client.NewViaCEPClient(),
		weatherClient: client.NewWeatherAPIClient(weatherAPIKey),
	}
}

// NewWeatherServiceWithClients creates a WeatherService with custom clients (for testing)
func NewWeatherServiceWithClients(cepClient CEPClient, weatherClient WeatherClient) *WeatherService {
	return &WeatherService{
		cepClient:     cepClient,
		weatherClient: weatherClient,
	}
}

// ErrInvalidZipcode is returned when the CEP format is invalid
var ErrInvalidZipcode = fmt.Errorf("invalid zipcode")

// ErrZipcodeNotFound is returned when the CEP is not found
var ErrZipcodeNotFound = fmt.Errorf("can not find zipcode")

// GetTemperaturesByCEP fetches temperatures for the city corresponding to the given CEP
func (s *WeatherService) GetTemperaturesByCEP(cep string) (*TemperatureResponse, error) {
	if !isValidCEP(cep) {
		return nil, ErrInvalidZipcode
	}

	location, err := s.cepClient.GetLocationByCEP(cep)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	if location == nil {
		return nil, ErrZipcodeNotFound
	}

	city := location.Localidade
	if city == "" {
		return nil, ErrZipcodeNotFound
	}

	tempC, err := s.weatherClient.GetCurrentTemperature(city)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather: %w", err)
	}

	return &TemperatureResponse{
		TempC: round(tempC, 2),
		TempF: round(celsiusToFahrenheit(tempC), 2),
		TempK: round(celsiusToKelvin(tempC), 2),
	}, nil
}

// isValidCEP checks if the CEP has exactly 8 numeric digits
func isValidCEP(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	for _, c := range cep {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// celsiusToFahrenheit converts Celsius to Fahrenheit
func celsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

// celsiusToKelvin converts Celsius to Kelvin
func celsiusToKelvin(c float64) float64 {
	return c + 273
}

// round rounds a float64 to a given number of decimal places
func round(val float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(val*p) / p
}
