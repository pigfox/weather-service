package request

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	"weather-service/constants"
	"weather-service/structs"
)

type Resp struct {
	Code int         `json:"code,omitempty"` //nolint:gofmt
	Body interface{} `json:"body,omitempty"`
}

func Make(params structs.RequestParams) (structs.Weather, error) {
	var err error
	var req *http.Request

	var weather structs.Weather
	if params.Body == "" {
		req, err = http.NewRequest(params.Method, params.URL, http.NoBody)
	} else {
		req, err = http.NewRequest(params.Method, params.URL, bytes.NewBuffer([]byte(params.Body)))
	}

	if err != nil {
		return weather, err
	}

	for k, v := range params.Headers {
		req.Header.Add(k, v)
	}

	length := len(params.Body)
	if 0 < length {
		req.Header.Add("Content-Length", strconv.Itoa(length))
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.TIMEOUT*time.Millisecond)
	defer cancel()

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return weather, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return weather, err
	}

	err = json.Unmarshal(resBody, &weather)
	if err != nil {
		return weather, err
	}

	return weather, nil
}

func SetUpParams(lat string, lon string) structs.RequestParams {
	if len(lat) == 0 || len(lon) == 0 {
		log.Fatal("invalid lat/lon")
	}
	params := structs.RequestParams{}
	params.URL = "https://api.openweathermap.org/data/2.5/weather?units=imperial&lat=" + lat + "&lon=" + lon + "&appid=" + structs.Auth.ApiKEY
	params.Method = "GET"
	return params
}
