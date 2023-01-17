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

var usersWithAssignedCities map[int64]string = make(map[int64]string)

var ctx = context.Background()

func RunBot(messages chan string, response chan []string, rdb *redis.Client) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGTOKEN"))
	if err != nil {
		log.Panic(err)
	}

	initMenuButtonCommands(bot) // menu button commands list

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
		case "start":
			setState(update.SentFrom().ID, start, rdb)
			msg.Text = "Если хочешь узнать погоду в твоем городе, " +
				"сперва выбери город командой /city, а затем отправь команду /weather"
		case "status":
			msg.Text = "pong"
		case "city":
			msg.Text = "Отправь мне название города и код. Например, /Moscow RU"
			setState(update.SentFrom().ID, weather, rdb)
		case "weather":
			msg.Text = sendWeather(update, messages, response)
		default:
			if currentState(update.SentFrom().ID, rdb) == weather {
				if len(strings.Fields(update.Message.Text)) == 2 {
					msg.Text = assignCity(update.SentFrom().ID, update.Message.Text)
				}
			} else {
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

func assignCity(userID int64, city string) string {
	if len(strings.Fields(city)) != 2 { // validate message
		return "Некорректные данные. Если хочешь вернуться в начало, отправь /start"
	}
	usersWithAssignedCities[userID] = city
	return "Город выбран. Теперь можешь смотреть погоду когда захочешь командой /weather!"
}

func sendWeather(update tgbotapi.Update, messages chan string, response chan []string) string {
	if _, ok := usersWithAssignedCities[update.SentFrom().ID]; ok {
		messages <- usersWithAssignedCities[update.SentFrom().ID]
		time.Sleep(time.Second)

		weatherData := <-response // [description, temperature, cityName]
		prettyWeather := fmt.Sprintf("Сейчас в г. %s %s℃, %s",
			weatherData[cityName], weatherData[temperature], weatherData[description])

		return prettyWeather
	} else {
		return "Сперва выбери город. Отправь команду /city"
	}
}

func initMenuButtonCommands(bot *tgbotapi.BotAPI) {
	config := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "Начало работы с Weather Report",
		},
		tgbotapi.BotCommand{
			Command:     "city",
			Description: "Выбор города",
		},
		tgbotapi.BotCommand{
			Command:     "weather",
			Description: "Хочешь узнать погоду в любом городе мира?",
		},
		tgbotapi.BotCommand{
			Command:     "ping",
			Description: "Проверь, в порядке ли я!",
		},
	)

	_, _ = bot.Request(config)
}
