package main

import (
	"fmt"
	"log"
	"time"
	"weather-service/request"
	"weather-service/structs"
)

func main() {
	structs.SetAPIKey()
	locations := getLocations()
	for _, location := range locations.Area {
		requestParams := request.SetUpParams(location.Lat, location.Lon)
		data, err := request.Make(requestParams)
		if err != nil {
			log.Fatal(err)
		}

		weather(data)
		time.Sleep(1 * time.Second) //free plan only allows 60 reqs/min
	}
}

/*
Write an http server that uses the Open Weather API that exposes an endpoint that takes in lat/long coordinates.
This endpoint should return what the weather condition is outside in that area (snow, rain, etc),
whether it’s hot, cold, or moderate outside (use your own discretion on what temperature equates to each type).
*/

func weather(data structs.Weather) {
	fmt.Print("Weather in " + data.Name + " on ")
	t := time.Unix(int64(data.Dt), 0)
	readableTime := t.Format("2006-01-02 15:04:05")
	fmt.Println(readableTime)

	fmt.Println("Wind direction is at", data.Wind.Deg, "degrees")
	fmt.Println("Wind speed is", data.Wind.Speed, "mph")
	fmt.Println("Humidity is", data.Main.Humidity, "%")
	fmt.Println("The weather condition outside in that area is", data.Weather[0].Description)      //(snow, rain, etc)
	fmt.Println("The temperature is", int(data.Main.Temp), "F, which is", mapTemp(data.Main.Temp)) //hot, cold, or moderate outside
	fmt.Println("\n")
}

func mapTemp(temp float64) string {
	switch {
	case temp < 32:
		return "freezing"
	case temp > 32 && temp < 50:
		return "cold"
	case temp >= 50 && temp < 65:
		return "moderate"
	case temp >= 65 && temp < 80:
		return "pleasant"
	case temp >= 80:
		return "hot"
	}
	return ""
}
