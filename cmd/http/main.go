package main

import (
	"log"
	"net/http"
	"time"

	"github.com/itsdarkhost/rbk-week2/internal/config"
	"github.com/itsdarkhost/rbk-week2/internal/core/country"
	weatherclient "github.com/itsdarkhost/rbk-week2/internal/core/weather"
	weatherservice "github.com/itsdarkhost/rbk-week2/internal/service/weather"
	httptransport "github.com/itsdarkhost/rbk-week2/internal/transport/http"
	"github.com/itsdarkhost/rbk-week2/internal/utils"
)

func main() {
	cfg := config.Load()

	httpClient := utils.NewHTTPClient(time.Duration(cfg.Countries.RequestTimeoutSec) * time.Second)
	weatherAPI := weatherclient.NewClient(httpClient, cfg.OpenMeteo.GeocodingBaseURL, cfg.OpenMeteo.ForecastBaseURL)
	countryAPI := country.NewClient(httpClient, cfg.Countries.BaseURL, cfg.Countries.MetadataBaseURL)
	service := weatherservice.New(weatherAPI, countryAPI)

	httpHandler := httptransport.New(service)

	server := &http.Server{
		Addr:              ":" + cfg.HTTP.Port,
		Handler:           httpHandler.Router(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Listening on http://localhost:%s", cfg.HTTP.Port)
	log.Fatal(server.ListenAndServe())
}
