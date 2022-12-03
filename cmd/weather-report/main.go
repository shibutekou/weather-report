package main

import (
	"github.com/bruma1994/weather-report/internal/owm"
	"log"
	"os"
)

var apikey = os.Getenv("APIKEY")

func main() {
	coord, err := owm.GetCoordinates("Cheboksary", "RU", apikey)
	if err != nil {
		log.Println("failed to get coordinates!")
	}

	weatherResponse := owm.GetWeather(coord, apikey)
	log.Println(weatherResponse.Weather[0].Description)
}
