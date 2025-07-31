package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WeatherAPIKey string
	Port          string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	return &Config{
		WeatherAPIKey: os.Getenv("WEATHER_API_KEY"),
		Port:          os.Getenv("PORT"),
	}, nil
}
