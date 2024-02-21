package structs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Auth Credentials

func SetAPIKey() {
	Auth = getCredentials()
}

func getCredentials() Credentials {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
	apiKey := os.Getenv("API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("api key required")
	}
	return Credentials{ApiKEY: apiKey}
}

type Credentials struct {
	ApiKEY string
}

type Location struct {
	Lat  string
	Lon  string
	Name string
}

type Locations struct {
	Area []Location
}

type RequestParams struct {
	URL              string
	Method           string
	Headers          map[string]string
	Body             string
	ExpectedResponse interface{}
}

type Weather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}
