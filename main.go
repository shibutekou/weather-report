package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bruma1994/weather-report/models"
)

const apikey = "apikey"

func main() {
	coord := getCoordinates("Moscow", "RU")

	weatherResponse := getWeather(coord)
	log.Println(weatherResponse.Weather[0].Description)
}

func getCoordinates(city, countryCode string) models.Coordinates {
	res, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s,%s&limit=5&appid=%s", city, countryCode, apikey))
	if err != nil {
		log.Println(err.Error())
	}

	var coord []models.Coordinates

	resBody, err := io.ReadAll(res.Body)
	err = json.Unmarshal(resBody, &coord)
	if err != nil {
		log.Println(err.Error())
	}

	return coord[0]
}

func getWeather(coord models.Coordinates) models.WeatherResponse {
	res, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s", coord.Lat, coord.Lon, apikey))
	if err != nil {
		log.Println(err.Error())
	}

	var weatherResponse models.WeatherResponse

	resBody, err := io.ReadAll(res.Body)

	err = json.Unmarshal(resBody, &weatherResponse)
	if err != nil {
		log.Println(err.Error())
	}

	return weatherResponse
}
