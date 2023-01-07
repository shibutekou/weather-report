package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// constants for weather data from OWM
const (
	description = iota
	temperature
	cityName
)

// constants for fixing the state of the user
const (
	start = iota
	weather
)

var ctx = context.Background()

func RunBot(messages chan string, response chan []string, rdb *redis.Client) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGTOKEN"))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "state":
			msg.Text = strconv.Itoa(currentState(update.SentFrom().ID, rdb))
		case "start":
			setState(update.SentFrom().ID, start, rdb)
			msg.Text = "Если хочешь узнать погоду в твоем городе, отправь команду /weather"
		case "status":
			msg.Text = "Я в порядке :3"
		case "weather":
			msg.Text = "Отправь мне название города и код. Например, /Moscow RU"
			setState(update.SentFrom().ID, weather, rdb)
		default:
			if currentState(update.SentFrom().ID, rdb) == weather {
				msg.Text = sendWeather(update, messages, response)
			} else if currentState(update.SentFrom().ID, rdb) == start {
				msg.Text = "Не знаю такой команды. Отправь /start, чтобы начать диалог"
			}
		}

		if _, err = bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}

}

func currentState(userID int64, rdb *redis.Client) int {
	val, err := rdb.Get(ctx, strconv.Itoa(int(userID))).Result()
	if err != nil {
		log.Fatal(err.Error())
	}

	id, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(err.Error())
	}

	return id
}

func setState(userID int64, state int, rdb *redis.Client) {
	err := rdb.Set(ctx, strconv.Itoa(int(userID)), state, 0).Err()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func sendWeather(update tgbotapi.Update, messages chan string, response chan []string) string {
	if len(strings.Fields(update.Message.Text)) != 2 { // validate message
		return "Некорректные данные. Если хочешь вернуться в начало, отправь /start"
	} else {
		messages <- update.Message.Text
		time.Sleep(time.Second)

		weatherData := <-response // [description, temperature, cityName]
		prettyWeather := fmt.Sprintf("Сейчас в г. %s %s℃, %s",
			weatherData[cityName], weatherData[temperature], weatherData[description])

		return prettyWeather
	}
}
