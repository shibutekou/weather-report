package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		case "weather":
			err := rdb.Set(ctx, "state", "step-weather", 0).Err()
			if err != nil {
				log.Print(err.Error())
			}

			val, err := rdb.Get(ctx, "state").Result()
			log.Println(val)
			msg.Text = val
		case "help":
			msg.Text = "Отправь мне название города и код. Например, /Moscow RU"
		case "status":
			msg.Text = "Я в порядке :3"
		default:
			if len(strings.Fields(update.Message.Text)) != 2 { // validate message
				msg.Text = "Некорректные данные"
			} else {
				messages <- update.Message.Text
				time.Sleep(time.Second)

				weatherData := <-response // [description, temperature, cityName]
				prettyWeather := fmt.Sprintf("Сейчас в г. %s %s℃, %s",
					weatherData[2], weatherData[1], weatherData[0])

				msg.Text = prettyWeather
			}
		}

		if _, err = bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}

}

func currentState(userID string, rdb *redis.Client) string {
	val, err := rdb.Get(ctx, userID).Result()
	if err != nil {
		log.Fatal(err.Error())
	}

	return val
}

func setState(userID string, state, rdb *redis.Client) {
	err := rdb.Set(ctx, userID, state, 0).Err()
	if err != nil {
		log.Fatal(err.Error())
	}
}
