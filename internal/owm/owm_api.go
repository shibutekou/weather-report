package owm

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/bruma1994/weather-report/internal/owm/models"
)

func GetCoordinates(city, countryCode string, apikey string) (models.Coordinates, error) {
	res, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s,%s&limit=5&appid=%s",
		city,
		countryCode,
		apikey))

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

func Weather(city, code string, apikey string, weather chan []string) {
	coordiantes, err := GetCoordinates(city, code, apikey)
	if err != nil {
		log.Fatal(err.Error())
	}

	res, err := http.Get(
		fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lang=ru&lat=%f&lon=%f&appid=%s&units=metric",
			coordiantes.Lat,
			coordiantes.Lon,
			apikey))
	if err != nil {
		log.Println(err.Error())
	}

	var weatherResponse models.WeatherResponse

	resBody, err := io.ReadAll(res.Body)

	err = json.Unmarshal(resBody, &weatherResponse)
	if err != nil {
		log.Println(err.Error())
	}
	weatherData := []string{
		weatherResponse.Weather[0].Description,
		strconv.Itoa(int(weatherResponse.Main.Temp)),
		weatherResponse.Name,
	}

	weather <- weatherData
}
