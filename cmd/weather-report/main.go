package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bruma1994/weather-report/internal/owm"
	"github.com/bruma1994/weather-report/internal/telegram"
	"github.com/go-redis/redis/v8"
)

var apikey = os.Getenv("APIKEY")

func main() {
	var (
		weather  = make(chan []string, 1)
		response = make(chan []string, 1)
		messages = make(chan string, 1)
	)

	var rdb = InitRedisClient()

	log.Println("Service started...")

	go telegram.RunBot(messages, response, rdb)

	go func(weather, response chan []string, messages chan string) {
		for true {
			cityData := strings.Fields(<-messages) // ["/city", "code"]

			owm.Weather(cityData[0], cityData[1], apikey, weather)

			weatherData := <-weather
			response <- weatherData
		}
	}(weather, response, messages)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-ch

	log.Println("Service stopped...")
}

func InitRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}
