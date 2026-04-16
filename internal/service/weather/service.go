package weather

import (
	"context"
	"errors"
	"sort"
	"strings"

	corecountry "github.com/itsdarkhost/rbk-week2/internal/core/country"
	coreweather "github.com/itsdarkhost/rbk-week2/internal/core/weather"
)

var ErrCountryHasNoCities = errors.New("country has no cities")

type WeatherClient interface {
	GeocodeCity(ctx context.Context, city, countryCode string) (coreweather.Location, error)
	Current(ctx context.Context, latitude, longitude float64, timezone string) (coreweather.CurrentWeather, error)
}

type CountryClient interface {
	FindCountry(ctx context.Context, country string) (corecountry.Country, error)
	CountryCode(ctx context.Context, country string) (string, error)
}

type Service struct {
	weather   WeatherClient
	countries CountryClient
}

type CityWeatherResponse struct {
	City                 string  `json:"city"`
	Country              string  `json:"country"`
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	Timezone             string  `json:"timezone"`
	TemperatureC         float64 `json:"temperature_c"`
	ApparentTemperatureC float64 `json:"apparent_temperature_c"`
	WindSpeedKMH         float64 `json:"wind_speed_kmh"`
	WeatherCode          int     `json:"weather_code"`
	UpdatedAt            string  `json:"updated_at"`
	Recommendation       string  `json:"recommendation"`
}

type CountryWeatherResponse struct {
	Country         string                `json:"country"`
	RequestedCities int                   `json:"requested_cities"`
	ResolvedCities  int                   `json:"resolved_cities"`
	Cities          []CityWeatherResponse `json:"cities"`
}

type CountryTopResponse struct {
	Country         string                `json:"country"`
	RequestedCities int                   `json:"requested_cities"`
	ResolvedCities  int                   `json:"resolved_cities"`
	Top             []CityWeatherResponse `json:"top"`
}

func New(weather WeatherClient, countries CountryClient) *Service {
	return &Service{
		weather:   weather,
		countries: countries,
	}
}

func (s *Service) CityWeather(ctx context.Context, city string) (CityWeatherResponse, error) {
	location, err := s.weather.GeocodeCity(ctx, city, "")
	if err != nil {
		return CityWeatherResponse{}, err
	}

	current, err := s.weather.Current(ctx, location.Latitude, location.Longitude, location.Timezone)
	if err != nil {
		return CityWeatherResponse{}, err
	}

	return buildCityResponse(location, current), nil
}

func (s *Service) CountryWeather(ctx context.Context, country string) (CountryWeatherResponse, error) {
	response, err := s.collectCountryWeather(ctx, country)
	if err != nil {
		return CountryWeatherResponse{}, err
	}

	sort.Slice(response.Cities, func(i, j int) bool {
		return strings.ToLower(response.Cities[i].City) < strings.ToLower(response.Cities[j].City)
	})

	return response, nil
}

func (s *Service) TopWarmCities(ctx context.Context, country string) (CountryTopResponse, error) {
	response, err := s.collectCountryWeather(ctx, country)
	if err != nil {
		return CountryTopResponse{}, err
	}

	sort.Slice(response.Cities, func(i, j int) bool {
		if response.Cities[i].TemperatureC == response.Cities[j].TemperatureC {
			return strings.ToLower(response.Cities[i].City) < strings.ToLower(response.Cities[j].City)
		}
		return response.Cities[i].TemperatureC > response.Cities[j].TemperatureC
	})

	if len(response.Cities) > 3 {
		response.Cities = response.Cities[:3]
	}

	return CountryTopResponse{
		Country:         response.Country,
		RequestedCities: response.RequestedCities,
		ResolvedCities:  response.ResolvedCities,
		Top:             response.Cities,
	}, nil
}

func (s *Service) collectCountryWeather(ctx context.Context, country string) (CountryWeatherResponse, error) {
	countryData, err := s.countries.FindCountry(ctx, country)
	if err != nil {
		return CountryWeatherResponse{}, err
	}

	cityNames := uniqueCities(countryData.Cities)
	if len(cityNames) == 0 {
		return CountryWeatherResponse{}, ErrCountryHasNoCities
	}

	countryCode, err := s.countries.CountryCode(ctx, countryData.Name)
	if err != nil {
		countryCode = ""
	}

	cities := make([]CityWeatherResponse, 0, len(cityNames))
	seen := map[string]struct{}{}

	for _, cityName := range cityNames {
		location, err := s.weather.GeocodeCity(ctx, cityName, countryCode)
		if err != nil {
			continue
		}

		key := strings.ToLower(location.Name) + "|" + location.CountryCode
		if _, ok := seen[key]; ok {
			continue
		}

		current, err := s.weather.Current(ctx, location.Latitude, location.Longitude, location.Timezone)
		if err != nil {
			continue
		}

		seen[key] = struct{}{}
		cities = append(cities, buildCityResponse(location, current))
	}

	if len(cities) == 0 {
		return CountryWeatherResponse{}, coreweather.ErrLocationNotFound
	}

	return CountryWeatherResponse{
		Country:         countryData.Name,
		RequestedCities: len(cityNames),
		ResolvedCities:  len(cities),
		Cities:          cities,
	}, nil
}

func uniqueCities(cities []string) []string {
	unique := make([]string, 0, len(cities))
	seen := map[string]struct{}{}

	for _, city := range cities {
		city = strings.TrimSpace(city)
		if city == "" {
			continue
		}

		key := strings.ToLower(city)
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		unique = append(unique, city)
	}

	return unique
}

func buildCityResponse(location coreweather.Location, current coreweather.CurrentWeather) CityWeatherResponse {
	return CityWeatherResponse{
		City:                 location.Name,
		Country:              location.Country,
		Latitude:             location.Latitude,
		Longitude:            location.Longitude,
		Timezone:             location.Timezone,
		TemperatureC:         current.TemperatureC,
		ApparentTemperatureC: current.ApparentTemperatureC,
		WindSpeedKMH:         current.WindSpeedKMH,
		WeatherCode:          current.WeatherCode,
		UpdatedAt:            current.Time,
		Recommendation:       clothingRecommendation(current.TemperatureC),
	}
}

func clothingRecommendation(temperature float64) string {
	switch {
	case temperature < 0:
		return "куртка"
	case temperature < 10:
		return "теплая одежда"
	default:
		return "легкая одежда"
	}
}
