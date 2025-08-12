# Weather Service

## Coding Assessment

This is the code exam for new Golang roles — **Weather Service Assignment**.

### Assignment

Write an HTTP server that serves the current weather.  
Your server should expose an endpoint that:

1. Accepts **latitude** and **longitude** coordinates.
2. Returns the **short forecast** for that area for **today** (e.g., "Partly Cloudy").
3. Returns a classification of whether the temperature is **hot**, **cold**, or **moderate** (use your discretion to define thresholds).
4. Uses the **National Weather Service API** as a data source.

---

### Things to Consider

- The purpose of this exercise is to provide a sample of your work for discussion in the **technical interview**.
- We respect your time — spend as long as you need, but we expect it to take around **1 hour**.
- We do **not** expect a production-ready service, but you should comment on any shortcuts you take.
- The submitted project should **build** and include brief instructions so we can verify it works.
- You may use **any language or stack** you’re most comfortable in.

---
# From the project root
go mod tidy

Automatically - using a premade script to run everything
./run.sh

Manually
go test -v
go run main.go

# Los Angeles, CA
curl "http://localhost:8080/weather?lat=34.05&lon=-118.25"
{"short_forecast":"Sunny","temperature_f":84,"characterization":"hot","source":"National Weather Service"}

# New York, NY
curl "http://localhost:8080/weather?lat=40.7128&lon=-74.0060"
{"short_forecast":"Mostly Clear","temperature_f":76,"characterization":"moderate","source":"National Weather Service"}

# Phoenix, AZ
curl "http://localhost:8080/weather?lat=33.4484&lon=-112.0740"
{"short_forecast":"Slight Chance Showers And Thunderstorms","temperature_f":110,"characterization":"hot","source":"National Weather Service"}
