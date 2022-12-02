package owm

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bruma1994/weather-report/internal/owm/models"
	"gopkg.in/ini.v1"
)

func GetCoordinates(city, countryCode string, cfg *ini.File) (models.Coordinates, error) {
	res, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s,%s&limit=5&appid=%s",
		city,
		countryCode,
		cfg.Section("OWMConfiguration").Key("apikey").String()))

	if err != nil {
		return models.Coordinates{}, err
	}

	var coord []models.Coordinates

	resBody, err := io.ReadAll(res.Body)
	err = json.Unmarshal(resBody, &coord)
	if err != nil {
		return models.Coordinates{}, err
	}

	if len(coord) < 1 {
		return models.Coordinates{}, err
	} else {
		return coord[0], nil
	}
}

func GetWeather(coord models.Coordinates, cfg *ini.File) models.WeatherResponse {
	res, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s",
		coord.Lat,
		coord.Lon,
		cfg.Section("OWMConfiguration").Key("apikey").String()))

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
