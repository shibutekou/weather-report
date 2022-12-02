package main

import (
	"github.com/bruma1994/weather-report/internal/owm"
	"gopkg.in/ini.v1"
	"log"
)

func main() {
	cfg := loadConfig("//")

	coord, err := owm.GetCoordinates("Cheboksary", "RU", cfg)
	if err != nil {
		log.Println("failed to get coordinates!")
	}

	weatherResponse := owm.GetWeather(coord, cfg)
	log.Println(weatherResponse.Weather[0].Description)
}

func loadConfig(filepath string) *ini.File {
	cfg, err := ini.Load(filepath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return cfg
}
