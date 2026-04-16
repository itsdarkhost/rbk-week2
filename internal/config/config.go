package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTP      HTTPConfig
	OpenMeteo OpenMeteoConfig
	Countries CountriesConfig
}

type HTTPConfig struct {
	Port string `envconfig:"HTTP_PORT" default:"8080"`
}

type OpenMeteoConfig struct {
	GeocodingBaseURL string `envconfig:"OPEN_METEO_GEOCODING_BASE_URL" default:"https://geocoding-api.open-meteo.com"`
	ForecastBaseURL  string `envconfig:"OPEN_METEO_FORECAST_BASE_URL" default:"https://api.open-meteo.com"`
}

type CountriesConfig struct {
	BaseURL           string `envconfig:"COUNTRIES_API_BASE_URL" default:"https://countriesnow.space/api/v0.1"`
	MetadataBaseURL   string `envconfig:"COUNTRY_METADATA_API_BASE_URL" default:"https://restcountries.com"`
	RequestTimeoutSec int    `envconfig:"REQUEST_TIMEOUT_SEC" default:"10"`
}

func Load() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}
