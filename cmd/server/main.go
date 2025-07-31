package main

import (
	"log"
	"net/http"

	"github.com/Eliezer2000/weather-api/internal/config"
	"github.com/Eliezer2000/weather-api/internal/handler"
	"github.com/Eliezer2000/weather-api/internal/service"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	weatherService := service.NewWeatherService(cfg)

	weatherHandler := handler.NewWeatherHandler(weatherService)

	router := mux.NewRouter()
	router.HandleFunc("/weather/{cep}", weatherHandler.GetWeather).Methods("GET")

	log.Printf("Server running on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
