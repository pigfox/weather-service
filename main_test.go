package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mock NWS that returns 86°F and "Partly Cloudy"
func mockNWS(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/points/", func(w http.ResponseWriter, r *http.Request) {
		base := "http://" + r.Host
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"properties": map[string]any{
				"forecast": base + "/forecast",
			},
		})
	})
	mux.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"properties": map[string]any{
				"periods": []map[string]any{
					{
						"name":            "Today",
						"isDaytime":       true,
						"startTime":       "2025-08-12T09:00:00-07:00",
						"temperature":     86.0, // Fahrenheit
						"temperatureUnit": "F",
						"shortForecast":   "Partly Cloudy",
					},
				},
			},
		})
	})
	return httptest.NewServer(mux)
}

func TestWeatherHandler_Success(t *testing.T) {
	nws := mockNWS(t)
	defer nws.Close()

	srv := httptest.NewServer(weatherHandler(nws.URL))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/weather?lat=34.05&lon=-118.25")
	if err != nil {
		t.Fatalf("http get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d", res.StatusCode)
	}

	var out struct {
		ShortForecast    string  `json:"short_forecast"`
		TemperatureF     float64 `json:"temperature_f"`
		Characterization string  `json:"characterization"`
		Source           string  `json:"source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if out.ShortForecast != "Partly Cloudy" {
		t.Fatalf("short_forecast: got %q", out.ShortForecast)
	}
	if out.TemperatureF < 85.9 || out.TemperatureF > 86.1 {
		t.Fatalf("temperature_f: got %v", out.TemperatureF)
	}
	if out.Characterization != "hot" {
		t.Fatalf("characterization: got %q", out.Characterization)
	}
	if out.Source != "National Weather Service" {
		t.Fatalf("source: got %q", out.Source)
	}
}

func TestWeatherHandler_Validation(t *testing.T) {
	nws := mockNWS(t)
	defer nws.Close()

	srv := httptest.NewServer(weatherHandler(nws.URL))
	defer srv.Close()

	// Missing lat/lon
	res, err := http.Get(srv.URL + "/weather")
	if err != nil {
		t.Fatalf("http get: %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
	res.Body.Close()

	// Bad lat
	res, err = http.Get(srv.URL + "/weather?lat=999&lon=0")
	if err != nil {
		t.Fatalf("http get: %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
	res.Body.Close()
}

func TestWeatherHandler_NWSError(t *testing.T) {
	// NWS returns 500
	mux := http.NewServeMux()
	mux.HandleFunc("/points/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", 500)
	})
	badNWS := httptest.NewServer(mux)
	defer badNWS.Close()

	srv := httptest.NewServer(weatherHandler(badNWS.URL))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/weather?lat=34&lon=-118")
	if err != nil {
		t.Fatalf("http get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", res.StatusCode)
	}
}
