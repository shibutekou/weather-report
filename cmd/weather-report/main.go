package main

import (
	"github.com/bruma1994/weather-report/internal/owm"
	"github.com/bruma1994/weather-report/internal/telegram"
	"log"
	"os"
	"strings"
)

var apikey = os.Getenv("APIKEY")

func main() {
	weather := make(chan string, 1)
	messages := make(chan string, 1)
	response := make(chan string, 1)

	log.Println("Service started...")

	go telegram.Bot(messages, response)

	runApp(weather, messages, response)
}

func runApp(weather, messages, response chan string) {
	for true {
		cityData := strings.Fields(<-messages) // ["/city", "code"]

		go owm.Weather(cityData[0], cityData[1], apikey, weather)

		weatherData := <-weather
		response <- weatherData
	}
}
