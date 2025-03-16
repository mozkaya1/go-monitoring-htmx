package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type WeatherBucket struct {
	Status           int    `json:"status"`
	UpdateTime       string `json:"updatetime"`
	Location         string `json:"location"`
	Temp             string `json:"temp"`
	WeatherDesc      string `json:"weatherDesc"`
	Humidity         string `json:"humidity"`
	FeelsLikeC       string `json:"feelsLikeC"`
	WindspeedKm      string `json:"windspeedKm"`
	AreaName         string `json:"areaName"`
	Latitude         string `json:"latitude"`
	Longitude        string `json:"longitude"`
	Country          string `json:"country"`
	Sunrise          string `json:"sunrise"`
	Sunset           string `json:"sunset"`
	MoonIllumination string `json:"moon_illumination"`
	MoonPhase        string `json:"moon_phase"`
	Moonrise         string `json:"moonrise"`
	Moonset          string `json:"moonset"`
}

type Currency struct {
	Status int                `json:"status"`
	Assets map[string]float64 `json:"assets"`
}

type CryptoAsset struct {
	Symbol             string `json:"symbol"`
	LastPrice          string `json:"lastPrice"`
	PriceChangePercent string `json:"priceChangePercent"`
}

type Crypto struct {
	Status int                    `json:"status"`
	Asset  map[string]CryptoAsset `json:"asset"`
}

type APIResponse struct {
	Time          string        `json:"time"`
	WeatherBucket WeatherBucket `json:"weatherbucket"`
	Currency      Currency      `json:"currency"`
	Crypto        Crypto        `json:"crypto"`
}

func GetApi() (APIResponse, error) {
	// change API url according to your requirement as location, asset, coin etc..
	// Details : https://github.com/mozkaya1/go-api#
	resp, err := http.Get("http://localhost:8080/api?location=Kudelstaart&")
	if err != nil {
		return APIResponse{}, err
	}
	defer resp.Body.Close()
	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}
	var v APIResponse
	err = json.Unmarshal(r, &v)
	if err != nil {
		return APIResponse{}, err
	}
	return v, nil
}
