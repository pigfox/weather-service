package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultAddr    = ":8080"
	defaultNWSBase = "https://api.weather.gov"
	httpTimeout    = 7 * time.Second
)

type Period struct {
	Name            string  `json:"name"`
	IsDaytime       bool    `json:"isDaytime"`
	StartTime       string  `json:"startTime"`
	Temperature     float64 `json:"temperature"`
	TemperatureUnit string  `json:"temperatureUnit"`
	ShortForecast   string  `json:"shortForecast"`
}

type nwsPointsResp struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

type nwsForecastResp struct {
	Properties struct {
		Periods []Period `json:"periods"`
	} `json:"properties"`
}

type weatherResponse struct {
	ShortForecast    string  `json:"short_forecast"`
	TemperatureF     float64 `json:"temperature_f"`
	Characterization string  `json:"characterization"`
	Source           string  `json:"source"`
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = defaultAddr
	}
	nwsBase := os.Getenv("NWS_BASE")
	if nwsBase == "" {
		nwsBase = defaultNWSBase
	}
	nwsBase = strings.TrimRight(nwsBase, "/")

	mux := http.NewServeMux()
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	mux.Handle("/weather", weatherHandler(nwsBase))

	s := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("weather-service listening on %s (NWS_BASE=%s)", addr, nwsBase)
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}

func weatherHandler(nwsBase string) http.Handler {
	client := &http.Client{Timeout: httpTimeout}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		latStr := r.URL.Query().Get("lat")
		lonStr := r.URL.Query().Get("lon")

		lat, lon, err := parseLatLon(latStr, lonStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), httpTimeout)
		defer cancel()

		forecastURL, err := getForecastURL(ctx, client, nwsBase, lat, lon)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": fmt.Sprintf("failed to resolve forecast URL: %v", err)})
			return
		}

		period, err := fetchCurrentPeriod(ctx, client, forecastURL)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": fmt.Sprintf("failed to fetch forecast: %v", err)})
			return
		}

		tempF := toFahrenheit(period.Temperature, period.TemperatureUnit)
		resp := weatherResponse{
			ShortForecast:    period.ShortForecast,
			TemperatureF:     round1(tempF),
			Characterization: characterizeTempF(tempF),
			Source:           "National Weather Service",
		}

		_ = json.NewEncoder(w).Encode(resp)
	})
}

func parseLatLon(latStr, lonStr string) (float64, float64, error) {
	if latStr == "" || lonStr == "" {
		return 0, 0, errors.New("missing lat or lon query params")
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid lat: %w", err)
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid lon: %w", err)
	}
	if lat < -90 || lat > 90 {
		return 0, 0, errors.New("lat must be between -90 and 90")
	}
	if lon < -180 || lon > 180 {
		return 0, 0, errors.New("lon must be between -180 and 180")
	}
	return lat, lon, nil
}

func getForecastURL(ctx context.Context, client *http.Client, base string, lat, lon float64) (string, error) {
	url := fmt.Sprintf("%s/points/%.4f,%.4f", base, lat, lon)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("points endpoint returned %d", res.StatusCode)
	}

	var pr nwsPointsResp
	if err := json.NewDecoder(res.Body).Decode(&pr); err != nil {
		return "", err
	}
	if pr.Properties.Forecast == "" {
		return "", errors.New("no forecast URL in response")
	}
	return pr.Properties.Forecast, nil
}

func fetchCurrentPeriod(ctx context.Context, client *http.Client, forecastURL string) (Period, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, forecastURL, nil)

	res, err := client.Do(req)
	if err != nil {
		return Period{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return Period{}, fmt.Errorf("forecast endpoint returned %d", res.StatusCode)
	}

	var fr nwsForecastResp
	if err := json.NewDecoder(res.Body).Decode(&fr); err != nil {
		return Period{}, err
	}
	if len(fr.Properties.Periods) == 0 {
		return Period{}, errors.New("no periods in forecast")
	}

	return fr.Properties.Periods[0], nil
}

func toFahrenheit(val float64, unit string) float64 {
	switch strings.ToUpper(strings.TrimSpace(unit)) {
	case "C":
		return (val * 9.0 / 5.0) + 32.0
	default:
		return val
	}
}

func characterizeTempF(tempF float64) string {
	if tempF < 50 {
		return "cold"
	}
	if tempF <= 77 {
		return "moderate"
	}
	return "hot"
}

func round1(x float64) float64 {
	return math.Round(x*10) / 10
}
