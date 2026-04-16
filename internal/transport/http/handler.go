package httptransport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	corecountry "github.com/itsdarkhost/rbk-week2/internal/core/country"
	coreweather "github.com/itsdarkhost/rbk-week2/internal/core/weather"
	weatherservice "github.com/itsdarkhost/rbk-week2/internal/service/weather"
)

type WeatherService interface {
	CityWeather(ctx context.Context, city string) (weatherservice.CityWeatherResponse, error)
	CountryWeather(ctx context.Context, country string) (weatherservice.CountryWeatherResponse, error)
	TopWarmCities(ctx context.Context, country string) (weatherservice.CountryTopResponse, error)
}

type Handler struct {
	service WeatherService
}

func New(service WeatherService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()

	r.Route("/weather", func(r chi.Router) {
		r.Get("/country/{country}/top", h.countryTop)
		r.Get("/country/{country}", h.country)
		r.Get("/{city}", h.city)
	})

	return r
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"routes": []string{
			"GET /weather/{city}",
			"GET /weather/country/{country}",
			"GET /weather/country/{country}/top",
		},
	})
}

func (h *Handler) city(w http.ResponseWriter, r *http.Request) {
	city := strings.TrimSpace(chi.URLParam(r, "city"))
	if city == "" {
		writeError(w, http.StatusBadRequest, "city is required")
		return
	}

	data, err := h.service.CityWeather(r.Context(), city)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (h *Handler) country(w http.ResponseWriter, r *http.Request) {
	country := strings.TrimSpace(chi.URLParam(r, "country"))
	if country == "" {
		writeError(w, http.StatusBadRequest, "country is required")
		return
	}

	data, err := h.service.CountryWeather(r.Context(), country)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (h *Handler) countryTop(w http.ResponseWriter, r *http.Request) {
	country := strings.TrimSpace(chi.URLParam(r, "country"))
	if country == "" {
		writeError(w, http.StatusBadRequest, "country is required")
		return
	}

	data, err := h.service.TopWarmCities(r.Context(), country)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, coreweather.ErrLocationNotFound):
		writeError(w, http.StatusNotFound, "город не найден")
	case errors.Is(err, corecountry.ErrCountryNotFound):
		writeError(w, http.StatusNotFound, "страна не найдена")
	case errors.Is(err, weatherservice.ErrCountryHasNoCities):
		writeError(w, http.StatusNotFound, "для страны не найден список городов")
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
