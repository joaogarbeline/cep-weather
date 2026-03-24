package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joaogarbeline/cep-weather/internal/handler"
	"github.com/joaogarbeline/cep-weather/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	svc := service.NewWeatherService(weatherAPIKey)
	h := handler.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.GetWeatherByCEP)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
