package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Request struct {
	Action string `json:"action"`
	Input  string `json:"input"`
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

type Response struct {
	Output string         `json:"output"`
	Error  string         `json:"error"`
	Data   map[string]any `json:"data"`
}

type InputParams struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type WeatherResponse struct {
	Current      CurrentWeather `json:"current"`
	CurrentUnits CurrentUnits   `json:"current_units"`
	Daily        DailyWeather   `json:"daily"`
}

type CurrentWeather struct {
	Temperature      float64 `json:"temperature_2m"`
	WeatherCode      int     `json:"weather_code"`
	WindSpeed        float64 `json:"wind_speed_10m"`
	RelativeHumidity int     `json:"relative_humidity_2m"`
}

type CurrentUnits struct {
	Temperature      string `json:"temperature_2m"`
	WindSpeed        string `json:"wind_speed_10m"`
	RelativeHumidity string `json:"relative_humidity_2m"`
}

type DailyWeather struct {
	Time           []string  `json:"time"`
	TemperatureMax []float64 `json:"temperature_2m_max"`
	TemperatureMin []float64 `json:"temperature_2m_min"`
	WeatherCode    []int     `json:"weather_code"`
}

func weatherDescription(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code == 1:
		return "Mainly clear"
	case code == 2:
		return "Partly cloudy"
	case code == 3:
		return "Overcast"
	case code == 45 || code == 48:
		return "Fog"
	case code == 51 || code == 53 || code == 55:
		return "Drizzle"
	case code == 56 || code == 57:
		return "Freezing drizzle"
	case code == 61 || code == 63 || code == 65:
		return "Rain"
	case code == 66 || code == 67:
		return "Freezing rain"
	case code == 71 || code == 73 || code == 75:
		return "Snowfall"
	case code == 77:
		return "Snow grains"
	case code == 80 || code == 81 || code == 82:
		return "Rain showers"
	case code == 85 || code == 86:
		return "Snow showers"
	case code == 95:
		return "Thunderstorm"
	case code == 96 || code == 99:
		return "Thunderstorm with hail"
	default:
		return "Unknown"
	}
}

func main() {
	var req Request
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "decode error: %v\n", err)
		os.Exit(1)
	}

	var params InputParams
	if err := json.Unmarshal([]byte(req.Input), &params); err != nil {
		json.NewEncoder(os.Stdout).Encode(Response{Error: "invalid input: " + err.Error()})
		return
	}

	if params.Latitude == "" || params.Longitude == "" {
		json.NewEncoder(os.Stdout).Encode(Response{Error: "latitude and longitude are required"})
		return
	}

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=temperature_2m,weather_code,wind_speed_10m,relative_humidity_2m&daily=temperature_2m_max,temperature_2m_min,weather_code&timezone=auto&forecast_days=3",
		params.Latitude, params.Longitude,
	)

	resp, err := http.Get(url)
	if err != nil {
		json.NewEncoder(os.Stdout).Encode(Response{Error: "API request failed: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		json.NewEncoder(os.Stdout).Encode(Response{Error: fmt.Sprintf("API returned status %d", resp.StatusCode)})
		return
	}

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		json.NewEncoder(os.Stdout).Encode(Response{Error: "failed to parse API response: " + err.Error()})
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Current: %.1f%s, %s, Wind %.0f %s, Humidity %d%s",
		weather.Current.Temperature, weather.CurrentUnits.Temperature,
		weatherDescription(weather.Current.WeatherCode),
		weather.Current.WindSpeed, weather.CurrentUnits.WindSpeed,
		weather.Current.RelativeHumidity, weather.CurrentUnits.RelativeHumidity,
	))

	if len(weather.Daily.Time) > 0 {
		sb.WriteString("\n\nForecast:")
		for i, date := range weather.Daily.Time {
			sb.WriteString(fmt.Sprintf("\n- %s: %.0f°C / %.0f°C, %s",
				date,
				weather.Daily.TemperatureMin[i],
				weather.Daily.TemperatureMax[i],
				weatherDescription(weather.Daily.WeatherCode[i]),
			))
		}
	}

	json.NewEncoder(os.Stdout).Encode(Response{Output: sb.String()})
}
