package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const weatherAPIBaseURL = "http://api.weatherapi.com/v1"

// WeatherAPIResponse represents the response from WeatherAPI
type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// WeatherAPIClient handles requests to the WeatherAPI
type WeatherAPIClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

// NewWeatherAPIClient creates a new WeatherAPI client
func NewWeatherAPIClient(apiKey string) *WeatherAPIClient {
	return &WeatherAPIClient{
		httpClient: &http.Client{},
		apiKey:     apiKey,
		baseURL:    weatherAPIBaseURL,
	}
}

// NewWeatherAPIClientWithBaseURL creates a new WeatherAPI client with a custom base URL (for testing)
func NewWeatherAPIClientWithBaseURL(apiKey, baseURL string) *WeatherAPIClient {
	return &WeatherAPIClient{
		httpClient: &http.Client{},
		apiKey:     apiKey,
		baseURL:    baseURL,
	}
}

// GetCurrentTemperature fetches the current temperature for the given city
func (c *WeatherAPIClient) GetCurrentTemperature(city string) (float64, error) {
	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("q", city)
	params.Set("aqi", "no")

	reqURL := fmt.Sprintf("%s/current.json?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather API returned status: %d", resp.StatusCode)
	}

	var result WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode weather response: %w", err)
	}

	return result.Current.TempC, nil
}
