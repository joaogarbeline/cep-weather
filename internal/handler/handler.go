package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/joaogarbeline/cep-weather/internal/service"
)

// WeatherServiceInterface defines the method needed from the service
type WeatherServiceInterface interface {
	GetTemperaturesByCEP(cep string) (*service.TemperatureResponse, error)
}

// Handler holds the HTTP handler dependencies
type Handler struct {
	svc WeatherServiceInterface
}

// NewHandler creates a new Handler
func NewHandler(svc WeatherServiceInterface) *Handler {
	return &Handler{svc: svc}
}

var onlyDigits = regexp.MustCompile(`^\d+$`)

// GetWeatherByCEP handles GET /{cep}
func (h *Handler) GetWeatherByCEP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract CEP from URL path: "/{cep}"
	cep := r.URL.Path[1:] // strip leading "/"

	result, err := h.svc.GetTemperaturesByCEP(cep)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidZipcode):
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		case errors.Is(err, service.ErrZipcodeNotFound):
			http.Error(w, "can not find zipcode", http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
