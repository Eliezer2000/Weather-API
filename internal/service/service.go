package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"github.com/Eliezer2000/weather-api/internal/config"
	"github.com/Eliezer2000/weather-api/internal/model"
)

// WeatherService gerencia a lógica de negócios
type WeatherService struct {
	config *config.Config
	client *http.Client
}

// NewWeatherService cria uma nova instância do WeatherService
func NewWeatherService(config *config.Config) *WeatherService {
	return &WeatherService{
		config: config,
		client: &http.Client{},
	}
}

// GetWeatherByCEP retorna as temperaturas para um CEP
func (s *WeatherService) GetWeatherByCEP(cep string) (*model.WeatherResponse, error) {
	if !isValidCEP(cep) {
		log.Printf("Invalid CEP: %s", cep)
		return nil, fmt.Errorf("invalid zipcode")
	}

	viaCEPResp, err := s.getLocationByCEP(cep)
	if err != nil {
		log.Printf("Error getting location for CEP %s: %v", cep, err)
		return nil, err
	}
	if viaCEPResp.Erro == "true" {
		log.Printf("CEP %s not found", cep)
		return nil, fmt.Errorf("can not find zipcode")
	}

	weatherResp, err := s.getWeatherByLocation(viaCEPResp.Localidade)
	if err != nil {
		log.Printf("Error getting weather for location %s: %v", viaCEPResp.Localidade, err)
		return nil, err
	}

	tempC := weatherResp.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	return &model.WeatherResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}

// isValidCEP valida se o CEP tem 8 dígitos numéricos
func isValidCEP(cep string) bool {
	matched, _ := regexp.MatchString(`^\d{8}$`, cep)
	return matched
}

// getLocationByCEP consulta a API ViaCEP
func (s *WeatherService) getLocationByCEP(cep string) (*model.ViaCEPResponse, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	log.Printf("ViaCEP request URL: %s", url)

	resp, err := s.client.Get(url)
	if err != nil {
		log.Printf("ViaCEP request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("ViaCEP response status: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading ViaCEP response body: %v", err)
		return nil, err
	}
	log.Printf("ViaCEP response body: %s", string(body))

	var viaCEPResp model.ViaCEPResponse
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&viaCEPResp); err != nil {
		log.Printf("ViaCEP JSON decode error: %v", err)
		return nil, err
	}

	return &viaCEPResp, nil
}

// getWeatherByLocation consulta a API WeatherAPI
func (s *WeatherService) getWeatherByLocation(location string) (*model.WeatherAPIResponse, error) {
	encodedLocation := url.QueryEscape(location)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", s.config.WeatherAPIKey, encodedLocation)
	log.Printf("WeatherAPI request URL: %s", url)

	resp, err := s.client.Get(url)
	if err != nil {
		log.Printf("WeatherAPI request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("WeatherAPI response status: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading WeatherAPI response body: %v", err)
		return nil, err
	}
	log.Printf("WeatherAPI response body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&errorResp); err == nil && errorResp.Error.Message != "" {
			log.Printf("WeatherAPI error response: %s", errorResp.Error.Message)
			return nil, fmt.Errorf("WeatherAPI error: %s", errorResp.Error.Message)
		}
		log.Printf("WeatherAPI returned unexpected status: %d", resp.StatusCode)
		return nil, fmt.Errorf("WeatherAPI returned status %d: %s", resp.StatusCode, string(body))
	}

	var weatherResp model.WeatherAPIResponse
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&weatherResp); err != nil {
		log.Printf("WeatherAPI JSON decode error: %v", err)
		return nil, err
	}
	return &weatherResp, nil
}