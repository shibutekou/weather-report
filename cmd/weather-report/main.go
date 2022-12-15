package main

import (
	"github.com/bruma1994/weather-report/internal/owm"
	"log"
	"os"
)

var apikey = os.Getenv("APIKEY")

func main() {
	weather := owm.Weather("Moscow", "RU", apikey)

	log.Println(weather.Weather[0].Description)
}
