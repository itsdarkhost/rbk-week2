package weather

type Location struct {
	Name        string
	Country     string
	CountryCode string
	Latitude    float64
	Longitude   float64
	Timezone    string
}

type CurrentWeather struct {
	Time                 string
	TemperatureC         float64
	ApparentTemperatureC float64
	WindSpeedKMH         float64
	WeatherCode          int
}

type geocodingResponse struct {
	Results []struct {
		Name        string  `json:"name"`
		Country     string  `json:"country"`
		CountryCode string  `json:"country_code"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		Timezone    string  `json:"timezone"`
	} `json:"results"`
}

type forecastResponse struct {
	Current struct {
		Time                 string  `json:"time"`
		TemperatureC         float64 `json:"temperature_2m"`
		ApparentTemperatureC float64 `json:"apparent_temperature"`
		WindSpeedKMH         float64 `json:"wind_speed_10m"`
		WeatherCode          int     `json:"weather_code"`
	} `json:"current"`
}
