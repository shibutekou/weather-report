package main

import (
	"github.com/bruma1994/weather-report/internal/owm"
	"github.com/bruma1994/weather-report/internal/telegram"
	"log"
	"os"
	"strings"
	"time"
)

var apikey = os.Getenv("APIKEY")

func main() {
	weather := make(chan string, 1)
	messages := make(chan string, 1)
	response := make(chan string, 1)

	go telegram.Bot(messages, response)

	data := <-messages
	cityData := strings.Fields(data)
	log.Println("Service started...")

	owm.Weather(cityData[0], cityData[1], apikey, weather)

	weatherData := <-weather
	response <- weatherData
	time.Sleep(3 * time.Second)
}
