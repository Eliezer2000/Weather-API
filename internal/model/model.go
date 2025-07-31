package model

// ViaCEPResponse representa a resposta da API ViaCEP
type ViaCEPResponse struct {
	Cep        string `json:"cep"`
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"` 
}

// WeatherAPIResponse representa a resposta da API WeatherAPI
type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// WeatherResponse representa a resposta da API interna
type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Message string `json:"message"`
}