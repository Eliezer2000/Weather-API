package handler

import (
	"encoding/json"
	"net/http"

	"log"
	"github.com/Eliezer2000/weather-api/internal/model"
	"github.com/Eliezer2000/weather-api/internal/service"
	"github.com/gorilla/mux"
)

// WeatherHandler gerencia as requisições HTTP
type WeatherHandler struct {
	service *service.WeatherService
}

// NewWeatherHandler cria uma nova instância do WeatherHandler
func NewWeatherHandler(service *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{service: service}
}

// GetWeather lida com a requisição GET /weather/{cep}
func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	vars := mux.Vars(r)
	cep := vars["cep"]

	weather, err := h.service.GetWeatherByCEP(cep)
	if err != nil {
		log.Printf("Error from service for CEP %s: %v", cep, err) // Log para depuração
		var status int
		var message string
		switch err.Error() {
		case "invalid zipcode":
			status = http.StatusUnprocessableEntity // 422
			message = "invalid zipcode"
		case "can not find zipcode":
			status = http.StatusNotFound // 404
			message = "can not find zipcode"
		default:
			status = http.StatusInternalServerError // 500
			message = "internal server error"
		}
		log.Printf("Returning error: status=%d, message=%s", status, message) // Log para depuração
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: message})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)
}