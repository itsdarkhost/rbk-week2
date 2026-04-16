package weather

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/itsdarkhost/rbk-week2/internal/utils"
)

var ErrLocationNotFound = errors.New("location not found")

type Client struct {
	geocoding *utils.Client
	forecast  *utils.Client
}

func NewClient(httpClient *http.Client, geocodingBaseURL, forecastBaseURL string) *Client {
	return &Client{
		geocoding: utils.NewClient(geocodingBaseURL, httpClient, nil),
		forecast:  utils.NewClient(forecastBaseURL, httpClient, nil),
	}
}

func (c *Client) GeocodeCity(ctx context.Context, city, countryCode string) (Location, error) {
	query := map[string]string{
		"name":     city,
		"count":    "1",
		"language": "en",
		"format":   "json",
	}
	if countryCode != "" {
		query["countryCode"] = strings.ToUpper(countryCode)
	}

	resp, err := c.geocoding.Get(ctx, "/v1/search", query)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return Location{}, utils.ReadErrorResponse(resp)
	}

	var payload geocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return Location{}, err
	}

	if len(payload.Results) == 0 {
		return Location{}, ErrLocationNotFound
	}

	result := payload.Results[0]

	return Location{
		Name:        result.Name,
		Country:     result.Country,
		CountryCode: result.CountryCode,
		Latitude:    result.Latitude,
		Longitude:   result.Longitude,
		Timezone:    result.Timezone,
	}, nil
}

func (c *Client) Current(ctx context.Context, latitude, longitude float64, timezone string) (CurrentWeather, error) {
	resp, err := c.forecast.Get(ctx, "/v1/forecast", map[string]string{
		"latitude":         formatCoordinate(latitude),
		"longitude":        formatCoordinate(longitude),
		"current":          "temperature_2m,apparent_temperature,wind_speed_10m,weather_code",
		"temperature_unit": "celsius",
		"wind_speed_unit":  "kmh",
		"timezone":         fallbackTimezone(timezone),
		"forecast_days":    "1",
	})
	if err != nil {
		return CurrentWeather{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return CurrentWeather{}, utils.ReadErrorResponse(resp)
	}

	var payload forecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return CurrentWeather{}, err
	}

	return CurrentWeather{
		Time:                 payload.Current.Time,
		TemperatureC:         payload.Current.TemperatureC,
		ApparentTemperatureC: payload.Current.ApparentTemperatureC,
		WindSpeedKMH:         payload.Current.WindSpeedKMH,
		WeatherCode:          payload.Current.WeatherCode,
	}, nil
}

func formatCoordinate(value float64) string {
	return strconv.FormatFloat(value, 'f', 6, 64)
}

func fallbackTimezone(timezone string) string {
	if strings.TrimSpace(timezone) == "" {
		return "auto"
	}
	return timezone
}
