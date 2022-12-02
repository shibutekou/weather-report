package main

import (
	"github.com/bruma1994/weather-report/internal/owm"
	"gopkg.in/ini.v1"
	"log"
)

func main() {
	cfg, err := ini.Load("owmconfig.ini")
	if err != nil {
		log.Fatalln(err.Error())
	}

	coord := owm.GetCoordinates("Moscow", "RU", cfg)

	weatherResponse := owm.GetWeather(coord, cfg)
	log.Println(weatherResponse.Weather[0].Icon)
}
